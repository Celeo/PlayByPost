<template lang="pug">
  div(v-if="!isLoading")
    div(v-if="!error")
      div(v-if="postIDs && postIDs.length > 0")
        div.text-xs-center(v-if="postIDs.length > postsPerPage")
          v-pagination.mb-2(:length="paginationLength" v-model="currentPage" :total-visible="paginationButtonCount" circle)
        post(v-for="id in postsInThisPage" :key="id" :id="id")
      div(v-else)
        h1 No posts have been made
      div.mt-5
        v-layout(wrap v-if="this.$store.getters.isLoggedIn")
          v-flex(lg8 xs12)
            editor(:func="save" v-model="newContent")
          v-flex.pl-3(lg4 xs12)
            roller.mt-3
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") {{ error }}
</template>

<script>
import API from '@/api'
import debounce from 'lodash/debounce'
import Editor from '@/components/Editor'
import Roller from '@/components/Roller'
import Post from '@/components/Post'

export default {
  components: {
    Editor,
    Roller,
    Post
  },
  data () {
    return {
      newContent: '',
      postIDs: [],
      isLoading: true,
      error: null,
      currentPage: 1,
      postsPerPage: null,
      newestAtTop: null
    }
  },
  methods: {
    async loadPostIDs () {
      try {
        const response = await new API(this).getAllPostIDs()
        this.postIDs = response.data
        this.error = null
      } catch (err) {
        console.error(err)
        this.error = 'Error loading posts'
      } finally {
        this.isLoading = false
      }
    },
    async save () {
      try {
        await new API(this).saveNewPost(this.newContent)
        this.error = null
        this.newContent = ''
        window.localStorage.removeItem('post')
        this.$store.commit('CLEAR_PENDING_ROLLS')
        await this.loadPostIDs()
        if (this.$store.getters.newestAtTop) {
          this.currentPage = 1
        } else {
          this.currentPage = this.paginationLength
        }
      } catch (err) {
        console.error(err)
        this.error = 'Error saving new post'
      }
    }
  },
  computed: {
    paginationLength () {
      return Math.ceil(this.postIDs.length / this.postsPerPage)
    },
    paginationButtonCount () {
      return screen.width < 1000 ? 5 : 7
    }
  },
  asyncComputed: {
    async postsInThisPage () {
      let retIds = Object.values(this.postIDs).slice()
      if (this.newestAtTop) {
        retIds = retIds.reverse()
      }
      retIds = retIds.slice((this.currentPage - 1) * this.postsPerPage, this.currentPage * this.postsPerPage)
      return retIds
    }
  },
  async created () {
    if (this.$store.getters.isLoggedIn) {
      this.postsPerPage = this.$store.getters.postsPerPage
      this.newestAtTop = this.$store.getters.newestAtTop
    } else {
      this.postsPerPage = 25
      this.newestAtTop = false
    }
  },
  async mounted () {
    await this.loadPostIDs()
    const postAttempt = window.localStorage.getItem('post')
    if (postAttempt) {
      this.newContent = postAttempt
    }
    if (this.$store.getters.goToLastPage) {
      this.$store.commit('SET_GO_TO_LAST_PAGE', false)
      this.currentPage = this.paginationLength
    }
  },
  watch: {
    newContent: debounce((val) => {
      window.localStorage.setItem('post', val)
    }, 1000)
  }
}
</script>
