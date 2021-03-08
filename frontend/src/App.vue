<template>
  <div id="app" class="uk-container uk-height-viewport">
    <nav class="uk-navbar-container uk-navbar-transparent" uk-navbar="mode: click">
      <div class="uk-navbar-left"></div>
      <div class="uk-navbar-right">
        <ul class="uk-navbar-nav">
          <li>
            <a class="uk-navbar-toggle" href="#" uk-navbar-toggle-icon></a>
            <div v-if="isLoggedIn" class="uk-navbar-dropdown">
              <ul class="uk-nav uk-navbar-dropdown-nav">
                <li><a href="#"><span class="uk-margin-small-right" uk-icon="icon: user"></span>Profile</a></li>
                <li class="uk-nav-divider"></li>
                <li><a href="#" @click="logout"><span class="uk-margin-small-right" uk-icon="icon: sign-out"></span>Logout</a></li>
              </ul>
            </div>
          </li>
        </ul>
      </div>
    </nav>
    <router-view></router-view>
  </div>
</template>

<script>
import UIKit from "uikit";
import Icons from "uikit/dist/js/uikit-icons";

UIKit.use(Icons);

export default {
  name: "App",
  computed: {
    isLoggedIn() {
      return this.$store.getters.isLoggedIn;
    },
    profile() {
      return this.$store.getters.profile;
    },
  },
  methods: {
    logout() {
      this.$router.push("/logout");
    },
  },
};
</script>

<style lang="less">
@import "../node_modules/uikit/src/less/uikit.less";
@import "./assets/less/theme.less";
</style>

<style>
#app {
  background-image: url("./assets/images/background.jpeg");
  background-repeat: no-repeat;
  background-size: cover;
  text-align: center;
}

#menu {
  color: #b4c2d1;
}

#menu-icon:hover {
  color: white;
}
</style>
