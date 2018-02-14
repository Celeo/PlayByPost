import Vue from 'vue'
import Vuetify from 'vuetify'

import App from './App'
import router from './router'
import store from './store'

import 'vuetify/src/stylus/main.styl'

Vue.use(Vuetify)

Vue.config.productionTip = false

let url = ''
if (process.env.NODE_ENV === 'development') {
  url = 'http://localhost:5000/'
} else {
  url = window.location.origin + '/api/'
}
Vue.config.SERVER_URL = url

const token = window.localStorage.getItem('token')
if (token) {
  store.commit('LOG_IN', JSON.parse(token))
}

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
