<template lang="pug">
  div(v-if="!isLoading")
    div(v-if="!error")
      div(v-if="posts && posts.length > 0")
        div.text-xs-center(v-if="posts.length > postsPerPage")
          v-pagination.mb-2(:length="paginationLength" v-model="page" :total-visible="paginationButtonCount" circle)
        post(v-for="post in postsThisPage()" :key="post.id" :post="post")
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
import debounce from 'lodash/debounce'
import Vue from 'vue'
import axios from 'axios'
import Editor from '@/components/Editor'
import Roller from '@/components/Roller'
import Post from '@/components/Post'

export default {
  components: {
    Editor,
    Roller,
    Post
  },
  data() {
    return {
      newContent: '',
      posts: [],
      isLoading: true,
      error: null,
      page: 1,
      postsPerPage: null,
      newestAtTop: null,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  methods: {
    async loadData() {
      try {
        const response = await axios.get(`${Vue.config.SERVER_URL}post`)
        this.posts = response.data
        this.error = null
      } catch (err) {
        console.error(err)
        this.error = 'Error loading posts'
      } finally {
        this.isLoading = false
      }
    },
    async save() {
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}post`, { content: this.newContent })
        this.error = null
        this.newContent = ''
        window.localStorage.removeItem('post')
        this.$store.commit('CLEAR_PENDING_ROLLS')
        await this.loadData()
        this.page = this.paginationLength
      } catch (err) {
        console.error(err)
        this.error = 'Error saving new post'
      }
    },
    postsThisPage() {
      let ret = this.posts.slice()
      if (this.newestAtTop) {
        ret = ret.reverse()
      }
      return ret.slice((this.page - 1) * this.postsPerPage, this.page * this.postsPerPage)
    }
  },
  computed: {
    paginationLength() {
      return Math.ceil(this.posts.length / this.postsPerPage)
    },
    paginationButtonCount() {
      if (screen.width < 1000) {
        return 5
      }
      return 7
    }
  },
  created() {
    if (this.$store.getters.isLoggedIn) {
      this.postsPerPage = this.$store.getters.postsPerPage
      this.newestAtTop = this.$store.getters.newestAtTop
    } else {
      this.postsPerPage = 25
      this.newestAtTop = false
    }
  },
  mounted() {
    this.loadData()
    const postAttempt = window.localStorage.getItem('post')
    if (postAttempt) {
      this.newContent = postAttempt
    }
  },
  watch: {
    newContent: debounce((val) => {
      window.localStorage.setItem('post', val)
    }, 1000)
  }
}
</script>
