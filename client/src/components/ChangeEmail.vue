<template lang="pug">
  div
    h3.display-1.mb-4 Change email
    v-form
      v-text-field(
        label="New email"
        v-model="email"
        required
      )
    v-btn(@click="submit" :disabled="!isValid || isLoading" color="primary") Submit
    span.grey--text(v-if="isLoading") Loading...
    v-alert.mt-3.black--text(type="warning" :value="true" v-if="!!error") {{ error }}
    v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Email changed successfully
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

export default {
  data() {
    return {
      email: '',
      valid: false,
      isLoading: false,
      error: null,
      successful: false,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  computed: {
    isValid() {
      return this.email !== ''
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}profile/email`, { email: this.email })
        this.successful = true
        this.error = null
      } catch (err) {
        console.error(err)
        this.error = 'Could not change email'
      } finally {
        this.isLoading = false
      }
    }
  },
  async created() {
    try {
      const response = await this.handler.get(`${Vue.config.SERVER_URL}profile`)
      this.email = response.data.email
    } catch (err) {
      console.error(err)
    }
  }
}
</script>
