<template lang="pug">
  div
    h4.display-1.mb-4 Change general info
    p.
      This is your current information. If you want to change something, change that box/picker/checkbox and click Submit.
      You can change multiple things at once. Your <strong>name</strong> here is both your <i>login</i> name and what
      appears on your posts as your character's name, so if you change it, don't forget that the next time you login,
      you'll be using the new name.
    v-form
      v-text-field(
        label="New name"
        v-model="name"
      )
      v-text-field(
        label="New email"
        v-model="email"
      )
      v-text-field(
        label="Posts per page"
        v-model="postsPerPage"
        type="number"
        min="0"
      )
      v-checkbox(
        label="Sort by newest posts at the top"
        v-model="newestAtTop"
        color="blue"
      )
      p Your 'tag' is what shows next to your name on your posts from now on.
      v-text-field(
        label="Tag"
        v-model="tag"
      )
    v-btn(
      @click="submit"
      :disabled="isLoading"
      color="primary"
    ) Submit
    span.grey--text(v-if="isLoading") Loading...
    v-alert.mt-3.black--text(type="error" :value="true" v-if="!!error") {{ error }}
    v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Info changed successfully
</template>

<script>
import { getAxios, buildEndpoint } from '@/util'

export default {
  data() {
    return {
      name: this.$store.getters.name,
      email: '',
      postsPerPage: this.$store.getters.postsPerPage,
      newestAtTop: this.$store.getters.newestAtTop,
      tag: this.$store.getters.tag,
      isLoading: false,
      error: null,
      successful: false
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      try {
        let updateData = {
          name: this.name,
          email: this.email,
          postsPerPage: this.postsPerPage.toString(),
          newestAtTop: this.newestAtTop,
          tag: this.tag
        }
        const response = await getAxios(this).put(buildEndpoint('profile'), updateData)
        this.successful = true
        this.error = null
        this.$store.commit('UPDATE_DATA', response.data)
        updateData['uuid'] = this.$store.getters.uuid
        window.localStorage.setItem('login', JSON.stringify(updateData))
      } catch (err) {
        console.error(err)
        this.error = 'Could not change information'
      } finally {
        this.isLoading = false
      }
    }
  },
  async created() {
    try {
      const response = await getAxios(this).get(buildEndpoint('profile'))
      this.email = response.data.email
    } catch (err) {
      console.error(err)
    }
  }
}
</script>
