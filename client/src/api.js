import Vue from 'vue'
import axios from 'axios'

const buildEndpoint = (path) => {
  return Vue.config.SERVER_URL + path
}

export default class {
  constructor (component) {
    this._component = component
    let ax
    if (component !== null) {
      ax = axios.create({ headers: { Authorization: component.$store.getters.uuid } })
    } else {
      ax = axios.create()
    }
    this._axios = ax
  }

  setPassword (oldPassword, newPassword) {
    return this._axios.put(buildEndpoint('profile/password'), { old: oldPassword, new: newPassword })
  }

  getEmail () {
    return this._axios.get(buildEndpoint('profile'))
  }

  login (name, password) {
    return this._axios.post(buildEndpoint('login'), { name, password })
  }

  saveStandardInformation (name, email, postsPerPage, newestAtTop, tag) {
    return this._axios.put(buildEndpoint('profile'), { name, email, postsPerPage, newestAtTop, tag })
  }

  getSinglePost (id) {
    return this._axios.get(buildEndpoint(`post/${id}`))
  }

  setPostContent (id, content) {
    return this._axios.post(buildEndpoint(`post/${id}`), { content })
  }

  invalidateLogins () {
    return this._axios.post(buildEndpoint('profile/invalidate'))
  }

  getAllPostIDs () {
    return this._axios.get(buildEndpoint('posts'))
  }

  saveNewPost (content) {
    return this._axios.post(buildEndpoint('posts'), { content })
  }

  register (name, password, email, code) {
    return this._axios.post(buildEndpoint('register'), { name, password, email, code })
  }

  search (term) {
    return this._axios.post(buildEndpoint(`posts/search/${term}`))
  }

  getPendingDice () {
    return this._axios.get(buildEndpoint('roll'))
  }

  rollDice (roll) {
    return this._axios.post(buildEndpoint('roll'), { roll })
  }
}
