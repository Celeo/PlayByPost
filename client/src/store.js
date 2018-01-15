import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  name: null,
  uuid: null
}

const mutations = {
  LOG_IN(state, payload) {
    const { uuid, name } = payload
    state.uuid = uuid
    state.name = name
  },

  LOG_OUT(state) {
    state.uuid = null
    state.name = null
  }
}

const getters = {
  isLoggedIn(state) {
    return !!state.uuid && !!state.name
  },

  uuid(state) {
    return state.uuid
  },

  name(state) {
    return state.name
  }
}

/* eslint-disable no-new */
export default new Vuex.Store({
  state,
  mutations,
  getters
})
