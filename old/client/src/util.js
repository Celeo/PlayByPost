export const formatRoll = (roll) => {
  if (roll.crit) {
    return `${roll.string} => ${roll.value} (crit!)`
  }
  return `${roll.string} => ${roll.value}`
}
