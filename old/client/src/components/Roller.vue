<template lang="pug">
  div
    h3.display-1.header.mb-1 Dice roller
    p.mb-2 For help with this, check the
      = " "
      router-link(:to="{ name: 'help' }") help
      = " "
      | page.
    v-text-field(label="What for" v-model="what")
    v-layout(row v-for="(die, index) in dice" :key="index")
      v-flex.pl-2(lg3 xs2)
        v-text-field(label="Count" v-model="die.count" type="number" min="1")
      v-flex.pl-2(lg3 xs2)
        v-text-field(label="Sides" v-model="die.sides" type="number" min="2" max="100")
      v-flex.pl-2(lg3 xs2)
        v-text-field(label="Mod" v-model="die.mod" type="number")
      v-flex
        v-btn.mx-1.px-0(fab small color="red" title="Remove this" @click="removeRow(index)" :disabled="dice.length === 1")
          v-icon(dark) fa-minus-circle
        v-btn.mx-1.px-0(fab small color="green" title="Add another dice" @click="addRow")
          v-icon(dark) fa-plus
    v-btn(block color="green" title="Roll the dice" @click="roll" :disabled="rolling || !valid") Roll!
    div(v-if="rolled.length > 0")
      v-list.pt-0
        v-list-tile(v-for="roll in rolled" :key="roll.id")
          v-list-tile-content.pt-0 {{ roll | filterRoll }}
    div(v-if="error")
      v-alert.mt-3.black--text(type="error" :value="true") An error occurred with the dice roller widget
</template>

<script>
import API from '@/api'
import { formatRoll } from '@/util.js'

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
  data () {
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
      ]
    }
  },
  computed: {
    valid () {
      if (this.what === '') {
        return false
      }
      for (let die of this.dice) {
        if (die.count < 1 || die.sides < 2 || die.sides > 100) {
          return false
        }
      }
      return true
    },
    rolled () {
      return this.$store.getters.pendingRolls
    }
  },
  methods: {
    addRow () {
      this.dice.push({ count: 1, sides: 20, mod: 0 })
    },
    removeRow (index) {
      this.dice.splice(index, 1)
    },
    async roll () {
      if (this.rolling) {
        return
      }
      this.rolling = true
      try {
        await new API(this).rollDice(this.what + ': ' + this.dice.map(e => formatDie(e)).join(', '))
        await this.loadPendingDice()
      } catch (err) {
        console.error(err)
        this.error = true
      }
      this.what = ''
      this.dice = [{ count: 1, sides: 20, mod: 0 }]
      this.rolling = false
    },
    async loadPendingDice () {
      try {
        const response = await new API(this).getPendingDice()
        this.$store.commit('SET_PENDING_ROLLS', response.data)
      } catch (err) {
        console.error(err)
        this.error = true
      }
    }
  },
  filters: {
    filterRoll (roll) {
      return formatRoll(roll)
    }
  },
  async created () {
    await this.loadPendingDice()
  }
}
</script>
