<template lang="pug">
  div
    v-container
      v-flex(xs6 offset-xs3)
        h1.display-2.mb-4 Register
        v-form
          v-text-field(
            label="Name"
            v-model="name"
            required
          )
          v-text-field(
            label="Password"
            v-model="password"
            type="password"
            required
          )
          v-text-field(
            label="Join code"
            v-model="code"
            required
          )
          v-text-field(
            label="Email (optional)"
            v-model="email"
          )
        v-btn(@click="submit" :disabled="!isValid || isLoading" color="primary") Submit
        span.grey--text(v-if="isLoading") Loading...
        v-alert.mt-3(type="error" :value="true" v-if="error") An error occurred with registration.
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

export default {
  data() {
    return {
      name: '',
      password: '',
      email: '',
      code: '',
      valid: false,
      isLoading: false,
      error: false
    }
  },
  computed: {
    isValid() {
      return this.name !== '' && this.password !== '' && this.code !== ''
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      const data = { name: this.name, password: this.password, email: this.email, code: this.code }
      try {
        const response = await axios.post(`${Vue.config.SERVER_URL}register`, data)
        const { name, uuid } = response.data
        const payload = { name, uuid }
        this.$store.commit('LOG_IN', payload)
        window.localStorage.setItem('token', JSON.stringify(payload))
        this.$router.push({ name: 'posts' })
        this.error = false
      } catch (err) {
        console.error(err)
        this.error = true
      } finally {
        this.isLoading = false
      }
    }
  }
}
</script>
