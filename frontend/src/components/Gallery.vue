<template>
  <div>
    <div v-if="gallery">
      <h1 class="uk-heading-divider">
        {{ gallery.title }}
      </h1>
      <h2>
        {{ gallery.description }}
      </h2>
      <div
        class="uk-position-relative uk-visible-toggle"
        tabindex="-1"
        uk-slideshow="animation: push"
      >
        <ul class="uk-slideshow-items">
          <li  v-for="(art, index) in gallery.art" :key="index">
            <img :src="art.uri" alt="" uk-cover />
            <div
              class="uk-position-center uk-position-small uk-text-center uk-light"
            >
              <h2 class="uk-margin-remove">{{art.title}}</h2>
              <p class="uk-margin-remove">
                {{art.description}}
              </p>
            </div>
          </li>
        </ul>

        <div class="uk-light">
          <a
            class="uk-position-center-left uk-position-small uk-hidden-hover"
            href="#"
            uk-slidenav-previous
            uk-slideshow-item="previous"
          ></a>
          <a
            class="uk-position-center-right uk-position-small uk-hidden-hover"
            href="#"
            uk-slidenav-next
            uk-slideshow-item="next"
          ></a>
        </div>
      </div>
    </div>
    <div v-else>
      <h1>Oops, something went wrong when collecting your gallery</h1>
    </div>
  </div>
</template>

<script>
export default {
  name: "Gallery",
  computed: {
    profile() {
      console.log("Gallery: For profile", this.$store.getters.profile);
      return this.$store.getters.profile;
    },
    profileUrl() {
      return "https://github.com/" + this.profile["login"];
    },
    gallery() {
      return this.$store.getters.gallery;
    },
  },
  mounted() {
    this.$store.dispatch("refreshGallery");
  },
};
</script>

<style scoped>
</style>