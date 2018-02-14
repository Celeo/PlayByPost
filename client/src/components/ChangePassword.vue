<template lang="pug">
  div
    v-container
      v-flex(xs6 offset-xs3)
        h1.display-2.mb-4 Change password
        v-form
          v-text-field(
            label="Old password"
            v-model="oldPassword"
            type="password"
            required
          )
          v-text-field(
            label="New password"
            v-model="newPassword"
            type="password"
            required
          )
        v-btn(@click="submit" :disabled="!isValid || isLoading" color="primary") Submit
        span.grey--text(v-if="isLoading") Loading...
        v-alert.mt-3.black--text(type="warning" :value="true" v-if="!!error") {{ error }}
        v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Password changed successfully
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

const mismatch = 'Old password does not match'

export default {
  data() {
    return {
      oldPassword: '',
      newPassword: '',
      valid: false,
      isLoading: false,
      error: null,
      successful: false,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  computed: {
    isValid() {
      return this.oldPassword !== '' && this.newPassword !== ''
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      const data = { old: this.oldPassword, new: this.newPassword }
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}changePassword`, data)
        this.successful = true
        this.error = null
        this.oldPassword = ''
        this.newPassword = ''
      } catch (err) {
        console.error(err)
        if (err.response && err.response.data.message === mismatch) {
          this.error = mismatch
        } else {
          this.error = 'Could not change password'
        }
      } finally {
        this.isLoading = false
      }
    }
  }
}
</script>
