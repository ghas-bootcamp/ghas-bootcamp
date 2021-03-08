import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'

import Axios from 'axios'

import App from './App.vue'
import Login from './components/Login.vue'
import Logout from './components/Logout.vue'
import NotFound from './components/NotFound.vue'
import Gallery from './components/Gallery.vue'
import AuthorizationCallback from './components/AuthorizationCallback.vue'

Vue.config.productionTip = false

Vue.use(Vuex)
Vue.use(VueRouter)

Vue.prototype.$http = Axios

const someUnusedProperty = new Vue({
  data: {
    myTestProperty: 2020
  },
  created: () => {
    console.log('myTestProperty is: ' + this.myTestProperty);
  }
});

const jwt = {
  decode(token) {
    if (!token) return {}
    const claimset = token.split('.', 3)[1]
    return JSON.parse(atob(claimset))
  },
  isExpired(token){
    const claimset = this.decode(token)
    console.log("Claimset:", claimset)
    const exp = claimset['exp']
    if (exp === undefined) {
      return false
    }
    console.log("Token expiration time since epoch:", exp)
    const nowLocal = new Date()
    const nowLocalInSecondsSincEpoch = Math.floor(nowLocal.getTime()/1000)
    console.log("Local time since epoch:", nowLocal.getTime())
    const nowUTCInSecondsSincEpoch = nowLocalInSecondsSincEpoch + nowLocal.getTimezoneOffset() * 60
    console.log("UTC time since epoch:", nowUTCInSecondsSincEpoch)
    console.log("Expiration delta: ", nowUTCInSecondsSincEpoch - exp)
    return exp <= nowUTCInSecondsSincEpoch
  }
}

const store = new Vuex.Store({
  state: {
    token: localStorage.getItem("token") || '',
    nonce: null,
    gallery: null
  },
  mutations: {
    token(state, token) {
      state.token = token
    },
    nonce(state, nonce) {
      state.nonce = nonce
    },
    gallery(state, gallery) {
      state.gallery = gallery
    }
  },
  actions: {
    authenticate({ commit }, {code, state}) {
      return new Promise((resolve, reject) => {
        Axios.get(`http://localhost:5000/authenticate/${code}`).then((response) => {
          if (response.data.token) {
            const token = response.data.token
            console.log("Authenticated with token", token)
            commit('token', token )
            localStorage.setItem('token', token)
  
            const [nonce, returnUrl] = atob(state).split(':', 2)
            console.log("Received authentication state with nonce", nonce, "and return url", returnUrl)
            resolve(returnUrl)
          } else {
            reject(response.data.error)
          }
        }).catch((e) => {
          reject(e)
        })
      })
    },
    logout({ commit }) {
      return new Promise((resolve) => {
        console.log("Logout")
        commit('token', null)
        localStorage.removeItem('token')
        resolve()
      })
    },
    nonce({commit}) {
      return new Promise((resolve) => {
        var stateArray = new Uint8Array(16);
        window.crypto.getRandomValues(stateArray);
        // First convert a typed array to normal array before transforming into a hex string.
        // Otherwise, we will loose randomness due to the wrapping of Uint8.
        const nonce = Array.from(stateArray)
          .map((i) => ("0" + i.toString(16)).slice(-2))
          .join("");

        commit('nonce', nonce)
        resolve(nonce)
      })
    },
    refreshGallery({commit}) {
      console.log("Triggering gallery refresh")
      return new Promise((resolve, reject) => {
        Axios.get("http://localhost:8081/gallery").then((response) => {
          console.log("Refreshed gallery with:", response.data)
          const gallery = response.data
          Axios.get("http://localhost:8081/gallery/art").then((response) => {
            gallery.art = response.data
            commit('gallery', gallery)
            resolve()
          }).catch((e) => {
            console.log("Failed to refresh gallery art with error:", e)
            reject()
          })
        }).catch((e) => {
          console.log("Failed to refresh gallery with error:", e)
          reject()
        })
      })
    }
  },
  getters: {
    isLoggedIn: state => !!state.token,
    profile: state => jwt.decode(state.token)['profile'],
    token: state => state.token,
    nonce: state => state.nonce,
    gallery: state => state.gallery
  }
})

Axios.interceptors.request.use(function (config) {
  const token = store.getters.token
  if (token != '') {
    config.headers.Authorization = `Bearer ${token}`
  } else {
    config.headers.Authorization = ''
  }
    
  return config;
});

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', redirect: '/gallery'},
    { path: '/login', name: 'Login', component: Login },
    { path: '/login/callback', name: "AuthorizationCallback", component: AuthorizationCallback },
    { path: '/logout', name: 'Logout', component: Logout},
    { path: '/gallery', name: 'Gallery', component: Gallery, meta: { requiresAuth: true } },
    { path: '*', component: NotFound }
  ]
})

router.beforeEach((to, from, next) => {
  console.log('beforeEach', to, from)
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters.isLoggedIn && !jwt.isExpired(store.getters.token)) {
      next()
      return
    }
    next({ path: '/login', query: { returnUrl: to.path } })
  }
  else {
    next()
  }
})

new Vue({
  render: h => h(App),
  store,
  router,
  jwt
}).$mount('#app')
