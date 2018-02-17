<template lang="pug">
  div
    h3.display-1.mb-4 Change name
    v-form
      v-text-field(
        label="New name"
        v-model="name"
        required
      )
    v-btn(@click="submit" :disabled="!isValid || isLoading" color="primary") Submit
    span.grey--text(v-if="isLoading") Loading...
    v-alert.mt-3.black--text(type="warning" :value="true" v-if="!!error") {{ error }}
    v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Name changed successfully
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

export default {
  data() {
    return {
      name: '',
      valid: false,
      isLoading: false,
      error: null,
      successful: false,
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  computed: {
    isValid() {
      return this.name !== ''
    }
  },
  methods: {
    async submit() {
      this.isLoading = true
      try {
        await this.handler.post(`${Vue.config.SERVER_URL}profile/name`, { name: this.name })
        this.successful = true
        this.error = null
        this.$store.commit('UPDATE_NAME', this.name)
      } catch (err) {
        console.error(err)
        this.error = 'Could not change name'
      } finally {
        this.isLoading = false
      }
    }
  },
  created() {
    this.name = this.$store.getters.name
  }
}
</script>
