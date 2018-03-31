<template lang="pug">
  div
    h4.display-1.mb-4 Log out everywhere else
    p.
      If you wish to log out of all <i>other</i> devices, click this button.
    v-btn(@click="invalidate" color="warning") Do it
    span.grey--text(v-if="isLoading") Loading...
    v-alert.mt-3.black--text(type="error" :value="true" v-if="!!error") {{ error }}
    v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Logins invalidated successfully
    
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

export default {
  data() {
    return {
      isLoading: null,
      error: null,
      successful: false,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  methods: {
    async invalidate() {
      this.isLoading = true
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}profile/invalidate`)
        this.successful = true
        this.error = null
      } catch (err) {
        console.error(err)
        this.error = 'Could not invalidate other logins'
        this.successful = false
      } finally {
        this.isLoading = false
      }
    }
  }
}
</script>
