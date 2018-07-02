import Vue from 'vue'
import axios from 'axios'

let axiosInstance = null

export const formatRoll = (roll) => {
  if (roll.crit) {
    return `${roll.string} => ${roll.value} (crit!)`
  }
  return `${roll.string} => ${roll.value}`
}

export const getAxios = (component) => {
  if (axiosInstance === null) {
    axiosInstance = axios.create({ headers: { Authorization: component.$store.getters.uuid } })
  }
  return axiosInstance
}

export const buildEndpoint = (path) => {
  return Vue.config.SERVER_URL + path
}
