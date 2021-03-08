<template>
  <div>
    <div v-if="hasError" class="uk-alert-danger" uk-alert>
      <a class="uk-alert-close" uk-close></a>
      <p v-html="error"></p>
    </div>
    <div class="uk-flex uk-flex-center uk-flex-middle uk-text-center">
      <div class="uk-card uk-card-default">
        <div class="uk-card-media-top">
          <img
            src="@/assets/images/gh-mark.png"
            class=""
            width="100px"
            height="100px"
          />
        </div>
        <div class="uk-card-body">
          <a href="#" @click="login" class="uk-button uk-button-text"
            >Login with GitHub</a
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Login",
  computed: {
    hasError() {
      const params = new URLSearchParams(window.location.search)
      return params.has('error')
    },
    error() {
      const params = new URLSearchParams(window.location.search)
      return params.get('error')
    },
  },
  methods: {
    login() {
      const returnUrl = this.$route.query.returnUrl || "/";

      this.$store.dispatch("nonce").then(nonce => {
        const state = btoa(`${nonce}:${returnUrl}`)
        window.location =
        `https://github.com/login/oauth/authorize?client_id=aad2f5434c1efdb25ce6&return_url=${window.location.origin}/login/callback&state=${state}`
      })
    },
  },
};
</script>

<style scoped>
</style>