import Vue from 'vue'
import Vuex from 'vuex'
import rawData from '@/rawData'

Vue.use(Vuex)

const state = {
  name: null,
  token: null,
  posts: rawData.posts
}

const mutations = {
  LOG_IN(state, payload) {
    const { token, tokenData } = payload
    state.token = token
    state.name = tokenData.name
  },

  LOG_OUT(state) {
    state.token = null
    state.name = null
  },

  NEW_POST(state, payload) {
    // TODO
  },

  EDIT_POST(state, payload) {
    // TODO
  }
}

const getters = {
  isLoggedIn(state) {
    return true // TODO
  },

  token(state) {
    return state.token
  },

  name(state) {
    return state.name
  },

  posts(state) {
    return state.posts
  }
}

/* eslint-disable no-new */
export default new Vuex.Store({
  state,
  mutations,
  getters
})
