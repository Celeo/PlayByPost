<template lang="pug">
  v-app(light)
    v-toolbar(app)
      v-toolbar-title Play By Post
      v-spacer
      v-toolbar-items
        v-btn(flat :exact="true" :to="{ name: 'posts' }") Posts
        v-btn(flat :exact="true" :to="{name: 'login'}" v-if="!loggedIn") Log in
        v-btn(flat :exact="true" :to="{name: 'register'}" v-if="!loggedIn") Register
        v-btn(flat :exact="true" :to="{name: 'logout'}" v-if="loggedIn") Log out
    v-content
      v-container
        transition(name="fade" mode="out-in")
          router-view
</template>

<script>
export default {
  computed: {
    loggedIn() {
      return this.$store.getters.isLoggedIn
    }
  }
}
</script>

<style lang="stylus">
.fade-enter-active, .fade-leave-active
  transition opacity .2s

.fade-enter, .fade-leave-active
  opacity 0

.ql-syntax
  background-color inherit !important
  color blue !important
  padding 0 !important
  margin 0 !important

.card__text p
  margin-bottom 0

.card__text ul
  margin-left 1rem

a.not-link
  text-decoration none
</style>
