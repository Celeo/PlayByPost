<template lang="pug">
  div
    h3.display-1.header.mt-5.mb-1 {{ title }}
    vue-editor(v-model="innerValue" :editorToolbar="toolbar")
    v-btn(color="info" :disabled="!hasWrittenContent" @click="func") {{ buttonText }}
      v-icon(right dark) fa-floppy-o
    br
    span.help.ml-2 Need help with formatting? Check the 
      router-link(:to="{ name: 'help' }") help
      |  page.
    div.below-editor
</template>

<script>
import { VueEditor } from 'vue2-editor'

export default {
  components: {
    VueEditor
  },
  props: {
    value: String,
    title: {
      type: String,
      default: 'New post'
    },
    buttonText: {
      type: String,
      default: 'Save'
    },
    func: {
      type: Function,
      default: () => { console.error('Unimplemented editor save') }
    }
  },
  data() {
    return {
      innerValue: this.value,
      toolbar: [
        ['bold', 'italic'],
        [{ 'header': [false, 3, 6] }],
        [{ list: 'bullet' }],
        [{ color: ['blue', 'black'] }, 'clean']
      ]
    }
  },
  computed: {
    hasWrittenContent() {
      return this.innerValue.length > 0 && this.innerValue !== '<p><br></p>'
    }
  },
  watch: {
    innerValue(val) {
      this.$emit('input', val)
    },
    value(val) {
      this.innerValue = val
    }
  }
}
</script>

<style lang="stylus" scoped>
.help
  font-size 1rem !important

.below-editor
  height 5rem
</style>
