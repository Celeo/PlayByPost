<template lang="pug">
  div
    h1.display-2.pad-bottom Glossary
    div(v-if="loaded && !error")
      div(v-if="glossary.content")
        div.text-block(v-html="glossary.content")
      p(v-else) There's nothing here.
      div.pad-bottom
      editor(:func="save" v-model="newContent" title="Edit glossary" v-if="isDM")
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") {{ error }}
</template>

<script>
import API from '@/api'
import Editor from '@/components/Editor'

export default {
  components: {
    Editor
  },
  data () {
    return {
      glossary: null,
      loaded: false,
      error: null,
      newContent: ''
    }
  },
  computed: {
    isDM () {
      return this.$store.getters.name === 'Dungeon Master'
    }
  },
  methods: {
    async save () {
      try {
        await new API(this).changeGlosary(this.newContent)
        await this.load()
      } catch (err) {
        console.error(err)
        this.error = 'Could not save changes'
      }
    },
    async load () {
      try {
        const response = await new API(this).getGlossary()
        this.glossary = response.data
        this.newContent = this.glossary.content
      } catch (error) {
        console.log(error)
        this.error = 'Could not load glossary'
      } finally {
        this.loaded = true
      }
    }
  },
  async created () {
    await this.load()
  }
}
</script>
