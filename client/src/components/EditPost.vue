<template lang="pug">
  div
    v-btn.ml-0(:to="{ name: 'posts' }")
      v-icon(left) fa-chevron-left
      | Back to posts
    div.header Note: if you are writing about additional actions that require more rolls, just make another post.
    editor.mt-3(:func="save" v-model="newContent" title="Edit post")
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") {{ error }}
    div(v-show="saved")
      v-alert.mt-3.black--text(type="success" :value="true") Changed saved
</template>

<script>
import Vue from 'vue'
import axios from 'axios'
import Editor from '@/components/Editor'

const editWindowError = 'Post is outside of edit window and can no longer be edited'

export default {
  components: {
    Editor
  },
  data() {
    return {
      post: null,
      newContent: '',
      error: null,
      saved: false,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  methods: {
    async loadData() {
      try {
        const response = await this.handler.get(`${Vue.config.SERVER_URL}post/${this.$route.params.id}`)
        this.post = response.data
        this.newContent = this.post.content
        this.error = null
      } catch (err) {
        console.error(err)
        if (err.response && err.response.status === 403) {
          this.error = editWindowError
        } else {
          this.error = 'Could not fetch post data from the server'
        }
      }
    },
    async save() {
      try {
        await this.handler.put(`${Vue.config.SERVER_URL}post/${this.$route.params.id}`, { content: this.newContent })
        this.$router.push({})
        this.error = null
        this.saved = true
      } catch (err) {
        console.error(err)
        if (err.response && err.response.status === 403) {
          this.error = editWindowError
        } else {
          this.error = 'Could not save post'
        }
      }
    }
  },
  async created() {
    await this.loadData()
  }
}
</script>

<style lang="stylus" scoped>
.header
  font-size 1.5rem
</style>
