import Vue from 'vue'
import Router from 'vue-router'

import Posts from '@/components/Posts'
import Login from '@/components/Login'
import Logout from '@/components/Logout'
import Register from '@/components/Register'

import store from './store'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/',
      component: Posts,
      name: 'posts'
    },
    {
      path: '/login',
      component: Login,
      name: 'login'
    },
    {
      path: '/logout',
      component: Logout,
      name: 'logout'
    },
    {
      path: '/register',
      component: Register,
      name: 'register'
    }
  ],
  mode: 'history'
})

router.beforeEach((to, from, next) => {
  if (['login', 'logout'].indexOf(to.name) === -1) {
    if (!store.getters.isLoggedIn) {
      next({ name: 'login' })
      return
    }
  }
  next()
})

export default router
