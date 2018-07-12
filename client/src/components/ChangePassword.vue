<template lang="pug">
  div
    h4.display-1.mb-4 Change password
    p.
      If you want to change your password, supply your <i>old</i> password in the first textbox, and your <i>new</i> password in the second textbox.
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
    v-alert.mt-3.black--text(type="error" :value="true" v-if="!!error") {{ error }}
    v-alert.mt-3.black--text(type="info" :value="true" v-if="successful") Password changed successfully
</template>

<script>
import API from '@/api'

const mismatch = 'Old password does not match'

export default {
  data () {
    return {
      oldPassword: '',
      newPassword: '',
      isLoading: false,
      error: null,
      successful: false
    }
  },
  computed: {
    isValid () {
      return this.oldPassword !== '' && this.newPassword !== ''
    }
  },
  methods: {
    async submit () {
      this.isLoading = true
      try {
        await new API(this).setPassword(this.oldPassword, this.newPassword)
        this.successful = true
        this.error = null
        this.oldPassword = ''
        this.newPassword = ''
      } catch (err) {
        console.error(err)
        if (err.response && err.response.data.error === mismatch) {
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
