<template lang="pug">
  div
    h3.display-1.header.mb-1 Dice roller 
    p.mb-2 For help with this, check the 
      router-link(:to="{ name: 'help' }") help
      |  page.
    v-layout(row)
      v-flex(sm12 xs6)
        v-text-field(label="What for" v-model="what")
    v-layout(row v-for="(die, index) in dice" :key="index")
      v-flex.pl-2(sm3 xs2)
        v-text-field(label="Count" v-model="die.count" type="number" min="1")
      v-flex.pl-2(sm3 xs2)
        v-text-field(label="Sides" v-model="die.sides" type="number" min="2" max="20")
      v-flex.pl-2(sm3 xs2)
        v-text-field(label="Modifier" v-model="die.mod" type="number")
      v-flex
        v-btn(fab small color="red" title="Remove this" @click="removeRow(index)" :disabled="dice.length === 1")
          v-icon(dark) fa-minus-circle
        v-btn(fab small color="green" title="Add another dice" @click="addRow")
          v-icon(dark) fa-plus
    v-layout(row)
      v-flex(xs8)
        v-btn(block color="green" title="Roll the dice" @click="roll" :disabled="rolling || !valid") Roll!
    div(v-if="rolled.length > 0")
      v-list.pt-0
        v-list-tile(v-for="roll in rolled" :key="roll.id")
          v-list-tile-content.pt-0 {{ roll.string }} =&gt; {{ roll.value }}
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") An error occurred with the dice roller widget
</template>

<script>
import Vue from 'vue'
import axios from 'axios'

const formatDie = (die) => {
  let s = `${die.count}d${die.sides}`
  if (die.mod === 0) {
    return s
  }
  if (die.mod < 0) {
    return `${s} - ${Math.abs(die.mod)}`
  }
  return `${s} + ${die.mod}`
}

export default {
  data() {
    return {
      rolling: false,
      error: false,
      what: '',
      dice: [
        {
          count: 1,
          sides: 20,
          mod: 0
        }
      ],
      handler: axios.create({ headers: { Authorization: this.$store.getters.uuid } })
    }
  },
  computed: {
    valid() {
      if (this.what === '') {
        return false
      }
      for (let die of this.dice) {
        if (die.count < 1 || die.sides < 2 || die.sides > 20) {
          return false
        }
      }
      return true
    },
    rolled() {
      return this.$store.getters.pendingRolls
    }
  },
  methods: {
    addRow() {
      this.dice.push({ count: 1, sides: 20, mod: 0 })
    },
    removeRow(index) {
      this.dice.splice(index, 1)
    },
    async roll() {
      if (this.rolling) {
        return
      }
      this.rolling = true
      try {
        const roll = this.what + ': ' + this.dice.map(e => formatDie(e)).join(', ')
        await this.handler.post(`${Vue.config.SERVER_URL}roll`, { roll })
        await this.loadPendingDice()
      } catch (err) {
        console.error(err)
        this.error = true
      }
      this.what = ''
      this.dice = [{ count: 1, sides: 20, mod: 0 }]
      this.rolling = false
    },
    async loadPendingDice() {
      try {
        const response = await this.handler.get(`${Vue.config.SERVER_URL}roll`)
        this.$store.commit('SET_PENDING_ROLLS', response.data)
      } catch (err) {
        console.error(err)
        this.error = true
      }
    }
  },
  async created() {
    await this.loadPendingDice()
  }
}
</script>
