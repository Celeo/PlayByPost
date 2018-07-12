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
import API from '@/api'
import Editor from '@/components/Editor'

const editWindowError = 'Post is outside of edit window and can no longer be edited'

export default {
  components: {
    Editor
  },
  data () {
    return {
      post: null,
      newContent: '',
      error: null,
      saved: false
    }
  },
  methods: {
    async loadData () {
      try {
        const response = await new API(this).getSinglePost(this.$route.params.id)
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
    async save () {
      try {
        await new API(this).setPostContent(this.$route.params.id, this.newContent)
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
  async created () {
    await this.loadData()
  }
}
</script>

<style lang="stylus" scoped>
.header
  font-size 1.5rem
</style>
