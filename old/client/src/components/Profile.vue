<template lang="pug">
  div
    h1.display-2.mb-3 Profile
    v-container(fluid grid-list-md)
      v-layout(row wrap)
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              change-standard-information
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              change-password
    v-container(fluid grid-list-md)
      v-layout(row wrap)
        v-flex(d-flex lg6 md12 sm12)
          v-card.elevation-5(light)
            v-card-text
              invalidate-logins
</template>

<script>
import ChangeStandardInformation from '@/components/ChangeStandardInformation'
import ChangePassword from '@/components/ChangePassword'
import InvalidateLogins from '@/components/InvalidateLogins'

export default {
  components: {
    ChangeStandardInformation,
    ChangePassword,
    InvalidateLogins
  }
}
</script>
