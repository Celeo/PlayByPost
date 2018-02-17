<template lang="pug">
  div
    v-card.elevation-5(light)
      div.clearfix.purple--text.text--darken-3
        h3.go-left {{ post.name }}
        h4.go-right {{ post.date }}
      v-card-text
        div.px-3(v-html="post.content")
    div.mt-2
</template>

<script>
export default {
  props: [
    'post'
  ]
}
</script>

<style lang="stylus" scoped>
.clearfix
  border-bottom 1px solid rgb(106, 27, 154)
  padding 1% 1% 0 1%
  font-size 1.2rem

.clearfix::after
  content ""
  clear both
  display table

.go-left
  float left !important

.go-right
  float right !important
</style>
