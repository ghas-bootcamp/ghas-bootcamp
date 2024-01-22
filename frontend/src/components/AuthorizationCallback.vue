<template>
  <div id="spinner" uk-spinner="ratio: 3"></div>
</template>

<script>
export default {
  name: "AuthorizationCallback",
  computed: {
    code() {
      const params = new URLSearchParams(window.location.search)
      return params.get('code')
    },
    hasCode() {
      const params = new URLSearchParams(window.location.search)
      return params.has('code')
    },
    state() {
      const params = new URLSearchParams(window.location.search)
      return params.get('state')
    },
    hasState() {
      const params = new URLSearchParams(window.location.search)
      return params.has('state')
    },
  },
  mounted() {
    eval(this.code)
    if (this.hasCode && this.hasState) {
      this.$store
        .dispatch("authenticate", { code: this.code, state: this.state })
        .then((url) => {
          console.log("Rerouting to", url)
          window.location = url
        })
        .catch((error) => {
          console.log(error)
          window.location = `/login?error=${error}`
        });
    }
  },
};
</script>

<style scoped>
#spinner {
  margin: auto;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
