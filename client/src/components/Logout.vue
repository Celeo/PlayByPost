<template lang="pug">
  div
    h1.display-2 Logging out ...
</template>

<script>
export default {
  created() {
    this.$store.commit('LOG_OUT')
    window.localStorage.removeItem('login')
    setInterval(() => {
      window.location.href = '/'
    }, 1000)
  }
}
</script>
