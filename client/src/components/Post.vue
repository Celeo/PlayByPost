<template lang="pug">
  div
    v-card.elevation-5(light)
      div.clearfix.purple--text.text--darken-3
        h3.go-left {{ post.name }}
          span.tag(v-if="post.tag.length > 0")  ({{ post.tag }})
        h4.go-right
          span(title="This is UTC time") {{ post.date }}
          span.ml-2.pb-3(v-if="post.name === name && post.canEdit")
            router-link.no-deco(:to="{ name: 'edit', params: { id: post.id } }")
              v-icon(color="black" title="Edit this post") fa-edit
      v-card-text
        div.px-1(v-html="post.content")
        div.px-1(v-if='post.rolls.length > 0')
          div.pt-4
            div.green--text.text--darken-1(v-for="roll in post.rolls" :key="roll.id") {{ roll | filterRoll }}
    div.mt-2
</template>

<script>
import { formatRoll } from '@/util.js'

export default {
  props: [
    'post'
  ],
  computed: {
    name () {
      return this.$store.getters.name
    }
  },
  filters: {
    filterRoll (roll) {
      return formatRoll(roll)
    }
  }
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

.no-deco
  text-decoration none

.tag
  font-size 75%
</style>
