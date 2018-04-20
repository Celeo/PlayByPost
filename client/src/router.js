import Vue from 'vue'
import Router from 'vue-router'

import Posts from '@/components/Posts'
import Login from '@/components/Login'
import Logout from '@/components/Logout'
import Profile from '@/components/Profile'
import Register from '@/components/Register'
import Help from '@/components/Help'
import EditPost from '@/components/EditPost'

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
      path: '/profile',
      component: Profile,
      name: 'profile'
    },
    {
      path: '/register',
      component: Register,
      name: 'register'
    },
    {
      path: '/help',
      component: Help,
      name: 'help'
    },
    {
      path: '/edit/:id',
      component: EditPost,
      name: 'edit'
    },
    {
      path: '/last',
      redirect: to => {
        store.commit('SET_GO_TO_LAST_PAGE', true)
        return { name: 'posts' }
      }
    }
  ],
  mode: 'history'
})

router.beforeEach((to, from, next) => {
  if (to.name === 'logout' && !store.getters.isLoggedIn) {
    next({ name: 'posts' })
    return
  }
  if (to.name === 'login' && store.getters.isLoggedIn) {
    next({ name: 'posts' })
    return
  }
  next()
})

export default router
