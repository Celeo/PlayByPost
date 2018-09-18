import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  name: null,
  postsPerPage: null,
  newestAtTop: null,
  uuid: null,
  pendingRolls: [],
  goToLastPage: false,
  tag: null
}

const mutations = {
  LOG_IN (state, payload) {
    const { uuid, name, postsPerPage, newestAtTop, tag } = payload
    state.name = name
    state.postsPerPage = postsPerPage
    state.newestAtTop = newestAtTop
    state.uuid = uuid
    state.pendingRolls = []
    state.goToLastPage = false
    state.tag = tag
  },

  LOG_OUT (state) {
    state.name = null
    state.postsPerPage = null
    state.newestAtTop = null
    state.uuid = null
    state.pendingRolls = []
    state.goToLastPage = false
    state.tag = null
  },

  UPDATE_DATA (state, payload) {
    const { name, postsPerPage, newestAtTop, tag } = payload
    state.name = name
    state.postsPerPage = postsPerPage
    state.newestAtTop = newestAtTop
    state.tag = tag
  },

  SET_PENDING_ROLLS (state, rolls) {
    state.pendingRolls = rolls
  },

  CLEAR_PENDING_ROLLS (state) {
    state.pendingRolls = []
  },

  SET_GO_TO_LAST_PAGE (state, go) {
    state.goToLastPage = go
  }
}

const getters = {
  isLoggedIn (state) {
    return !!state.uuid && !!state.name
  },

  uuid (state) {
    return state.uuid
  },

  name (state) {
    return state.name
  },

  postsPerPage (state) {
    return state.postsPerPage
  },

  newestAtTop (state) {
    return state.newestAtTop
  },

  tag (state) {
    return state.tag
  },

  pendingRolls (state) {
    return state.pendingRolls
  },

  goToLastPage (state) {
    return state.goToLastPage
  }
}

/* eslint-disable no-new */
export default new Vuex.Store({
  state,
  mutations,
  getters
})
