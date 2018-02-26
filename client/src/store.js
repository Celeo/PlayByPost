import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  name: null,
  postsPerPage: null,
  newestAtTop: null,
  uuid: null,
  pendingRolls: []
}

const mutations = {
  LOG_IN(state, payload) {
    const { uuid, name, postsPerPage, newestAtTop } = payload
    state.name = name
    state.postsPerPage = postsPerPage
    state.newestAtTop = newestAtTop
    state.uuid = uuid
    state.pendingRolls = []
  },

  LOG_OUT(state) {
    state.name = null
    state.postsPerPage = null
    state.newestAtTop = null
    state.uuid = null
    state.pendingRolls = []
  },

  UPDATE_DATA(state, payload) {
    const { name, postsPerPage, newestAtTop } = payload
    state.name = name
    state.postsPerPage = postsPerPage
    state.newestAtTop = newestAtTop
  },

  SET_PENDING_ROLLS(state, rolls) {
    state.pendingRolls = rolls
  },

  CLEAR_PENDING_ROLLS(state) {
    state.pendingRolls = []
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
  },

  postsPerPage(state) {
    return state.postsPerPage
  },

  newestAtTop(state) {
    return state.newestAtTop
  },

  pendingRolls(state) {
    return state.pendingRolls
  }
}

/* eslint-disable no-new */
export default new Vuex.Store({
  state,
  mutations,
  getters
})
