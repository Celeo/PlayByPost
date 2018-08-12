<template lang="pug">
  v-app(light)
    v-navigation-drawer(temporary v-model="drawer" absolute)
      v-list.pt-0(dense)
        v-list-tile(:exact="true" :to="{ name: 'posts' }")
          v-list-tile-action
            v-icon fa-book
          v-list-tile-content
            v-list-tile-title Posts
        v-list-tile(:exact="true" :to="{ name: 'search' }" v-if="loggedIn")
          v-list-tile-action
            v-icon fa-search
          v-list-tile-content
            v-list-tile-title Search
        v-list-tile(:exact="true" :to="{ name: 'glossary' }")
          v-list-tile-action
            v-icon fa-user
          v-list-tile-content
            v-list-tile-title Glossay
        v-list-tile(:exact="true" :to="{ name: 'help' }")
          v-list-tile-action
            v-icon fa-question
          v-list-tile-content
            v-list-tile-title Help
        v-list-tile(:exact="true" :to="{ name: 'login' }" v-if="!loggedIn")
          v-list-tile-action
            v-icon fa-sign-in
          v-list-tile-content
            v-list-tile-title Log in
        v-list-tile(:exact="true" :to="{ name: 'register' }" v-if="!loggedIn")
          v-list-tile-action
            v-icon fa-user-plus
          v-list-tile-content
            v-list-tile-title Register
        v-list-tile(:exact="true" :to="{ name: 'profile' }" v-if="loggedIn")
          v-list-tile-action
            v-icon fa-key
          v-list-tile-content
            v-list-tile-title Profile
        v-list-tile(:exact="true" :to="{ name: 'logout' }" v-if="loggedIn")
          v-list-tile-action
            v-icon fa-sign-out
          v-list-tile-content
            v-list-tile-title Log out
        v-divider
        v-list-tile(@click="drawer = false")
          v-list-tile-action
            v-icon fa-times
          v-list-tile-content
            v-list-tile-title Close this
    v-toolbar(app)
      v-toolbar-side-icon(@click="drawer = !drawer")
      v-toolbar-title Play By Post
      v-spacer
      v-toolbar-items(v-if="loggedIn")
        v-btn.not-link(flat :ripple="false") {{ username }}
    v-content
      v-container
        transition(name="fade" mode="out-in")
          router-view
</template>

<script>
export default {
  data () {
    return {
      drawer: false
    }
  },
  computed: {
    loggedIn () {
      return this.$store.getters.isLoggedIn
    },
    username () {
      return this.$store.getters.name
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

.not-link
  text-decoration none
  cursor default !important

.container
  padding 10px

.pad-bottom
  padding-bottom 2rem

.text-block
  font-size 120%
</style>
