<template lang="pug">
  div
    v-card.elevation-3(light)
      v-card-title(primary-title)
        a(:href="'#' + post.id").mr-1.not-link
          v-icon(color="grey lighten-1") fa-link
        h3.headline(:id="post.id") {{ post.name }} at {{ post.date }}
      v-card-text
        div.px-3(v-html="post.content")
</template>

<script>
export default {
  props: [
    'post'
  ]
}
</script>
