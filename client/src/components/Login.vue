<template lang="pug">
  div
    v-container
      v-flex(xs6 offset-xs3)
        h1.display-3 Login
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
            label="Email"
            v-model="email"
          )
        v-btn(@click="submit" :disabled="!isValid" color="primary") Submit
</template>

<script>
export default {
  data() {
    return {
      name: '',
      password: '',
      email: '',
      valid: false
    }
  },
  computed: {
    isValid() {
      return this.name !== '' && this.password !== ''
    }
  },
  methods: {
    async submit() {
      console.log('submit')
      // TODO
    }
  }
}
</script>
