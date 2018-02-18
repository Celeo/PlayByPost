<template lang="pug">
  div
    v-container
      v-flex(xs6 offset-xs3)
        h1.display-2.mb-4 Login
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
        v-btn(@click="submit" :disabled="!isValid || isLoading" color="primary") Submit
        span.grey--text(v-if="isLoading") Loading...
        v-alert.mt-3.black--text(type="warning" :value="true" v-if="error") Login unsuccessful
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

export default {
  data() {
    return {
      name: '',
      password: '',
      valid: false,
      isLoading: false,
      error: false
    }
  },
  computed: {
    isValid() {
      return this.name !== '' && this.password !== ''
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      const data = { name: this.name, password: this.password }
      try {
        const response = await axios.post(`${Vue.config.SERVER_URL}login`, data)
        this.$store.commit('LOG_IN', response.data)
        window.localStorage.setItem('login', JSON.stringify(response.data))
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
