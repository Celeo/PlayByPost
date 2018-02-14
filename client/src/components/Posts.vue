<template lang="pug">
  div
    div(v-if="!error")
      div(v-if="posts && posts.length > 0")
        post(v-for="post in posts" :key="post.id" :post="post")
      div(v-else)
        h1 No posts have been made
      editor(:func="save" v-model="newContent" v-if="this.$store.getters.isLoggedIn")
    div(v-else)
      v-alert.mt-3.black--text(type="warning" :value="true") {{ error }}
</template>

<script>
import Vue from 'vue'
import axios from 'axios'
import Editor from '@/components/Editor'
import Post from '@/components/Post'

export default {
  components: {
    Editor,
    Post
  },
  data() {
    return {
      newContent: '',
      posts: [],
      isLoading: true,
      error: null,
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
      }
    },
    async save() {
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}post`, { content: this.newContent })
        this.error = null
        this.newContent = ''
        this.loadData()
      } catch (err) {
        console.error(err)
        this.error = 'Error saving new post'
      }
    }
  },
  mounted() {
    this.loadData()
  }
}
</script>
