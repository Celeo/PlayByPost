<template lang="pug">
  div
    h1.display-3 Logging out ...
</template>

<script>
export default {
  created() {
    this.$store.commit('LOG_OUT')
    window.sessionStorage.removeItem('token')
    setInterval(() => {
      window.location.href = '/'
    }, 2000)
  }
}
</script>
