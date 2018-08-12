<template lang="pug">
  div
    h1.display-2.pad-bottom Search
    p Type the words you're searching for in the box below (minimum 3 letters), and click the button.
    v-form(@submit.prevent="submit")
      v-text-field(
        label="Search term"
        v-model="term"
        required
      )
      v-btn(@click="submit" :disabled="!validSearch || isLoading" color="primary") Submit
    div.results(v-show="hasResults")
      post(v-for="id in matchingIds" :key="id" :id="id")
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") {{ error }}
</template>

<script>
import API from '@/api'
import Post from '@/components/Post'

export default {
  components: {
    Post
  },
  data () {
    return {
      term: '',
      matchingIds: [],
      postMap: {},
      isLoading: false,
      error: null
    }
  },
  methods: {
    async submit () {
      try {
        this.isLoading = true
        const response = await new API(this).search(this.term)
        this.error = null
        this.matchingIds = response.data
      } catch (err) {
        console.error(err)
        this.error = 'Could not complete search'
      } finally {
        this.isLoading = false
      }
    },
    getPost (id) {
      if (id in this.postMap) {
        return this.postMap[id]
      }
      return new API(this).getSinglePost(id)
    }
  },
  computed: {
    validSearch () {
      return this.term.length > 2
    },
    hasResults () {
      return this.matchingIds.length > 0
    }
  }
}
</script>

<style lang="stylus" scoped>
.results
  margin-top 2rem
</style>
