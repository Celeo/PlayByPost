<template lang="pug">
  div
    h1.display-2.mb-3 Profile
    v-container(fluid grid-list-md)
      v-layout(row wrap)
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              change-password
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              change-name
      v-layout(row wrap)
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              change-email
</template>

<script>
import ChangePassword from '@/components/ChangePassword'
import ChangeName from '@/components/ChangeName'
import ChangeEmail from '@/components/ChangeEmail'

export default {
  components: {
    ChangePassword,
    ChangeName,
    ChangeEmail
  }
}
</script>
