<template lang="pug">
  div
    h3.display-1.mt-1.mb-2 Posts
    v-card.elevation-3(light v-for="post in posts" :key="post.id")
      v-card-title(primary-title)
        a(:href="'#' + post.id").mr-1.not-link
          v-icon(color="grey lighten-1") fa-link
        h3.headline {{ post.name }} at {{ post.timestamp }}
      v-card-text
        div.px-3(v-html="post.content")
    h3.display-1.header.mt-5.mb-2 New post
    vue-editor(v-model="newContent" :editorToolbar="toolbar")
    v-btn(color="info" :disabled="!hasWrittenContent" @click="save") Save
      v-icon(right dark) fa-floppy-o
</template>

<script>
import { VueEditor } from 'vue2-editor'

export default {
  components: {
    VueEditor
  },
  data() {
    return {
      newContent: '', // this.$store.getters.posts[0].content,
      toolbar: [
        ['bold', 'italic'],
        [{ list: 'bullet' }],
        [{ color: ['blue'] }, 'clean']
      ]
    }
  },
  methods: {
    save() {
      console.log('save')
      // TODO
    }
  },
  computed: {
    posts() {
      return this.$store.getters.posts
    },
    hasWrittenContent() {
      return this.newContent.length > 0 && this.newContent !== '<p><br></p>'
    }
  }
}
</script>
