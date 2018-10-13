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

const loadRoller = (campaignId) => {
  const api = window.location.origin + `/campaign/${campaignId}/roll`
  new Vue({
    el: '#rollerApp',
    data: {
      rolling: false,
      error: false,
      action: '',
      dice: [
        {
          count: 1,
          sides: 20,
          mod: 0
        }
      ],
      pending: []
    },
    computed: {
      valid () {
        if (this.action === '') {
          return false
        }
        for (let die of this.dice) {
          if (die.count < 1 || die.sides < 2 || die.sides > 100) {
            return false
          }
        }
        return true
      }
    },
    methods: {
      addRow () {
        this.dice.push({ count: 1, sides: 20, mod: 0 })
      },
      removeRow (index) {
        this.dice.splice(index, 1)
      },
      roll () {
        if (this.rolling) {
          return
        }
        this.rolling = true
        fetch(api, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ roll: this.action + ': ' + this.dice.map(e => formatDie(e)).join(', ') })
          })
          .then(response => {
            this.error = false
            return response.json()
          })
          .then(json => {
            this.pending = json
            this.action = ''
            this.dice = [
              {
                count: 1,
                sides: 20,
                mod: 0
              }
            ]
          })
          .catch(error => {
            console.error(error)
            this.error = true
          })
          .finally(() => this.rolling = false)
      },
      formatRoll (roll) {
        if (roll.is_crit) {
          return `${roll.string} => ${roll.value} (crit!)`
        }
        return `${roll.string} => ${roll.value}`
      }
    },
    mounted () {
      fetch(api)
        .then(response => {
          this.error = false
          return response.json()
        })
        .then(json => {
          this.pending = json
        })
        .catch(error => {
          console.error(error)
          this.error = true
        })
    }
  })
}

window.loadRoller = loadRoller
