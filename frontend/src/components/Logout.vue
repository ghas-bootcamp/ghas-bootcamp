<template>
         <div id="spinner" uk-spinner="ratio: 3"></div>
</template>

<script>
export default {
    name: 'Logout',
    mounted() {
        this.$store.dispatch('logout').then(() => {
            this.$router.replace('/').catch(failure => failure) // ignore failure because we will be redirected
        })
    }
}
</script>