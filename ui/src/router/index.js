import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import Home from '@/components/Home.vue'
import Login from '@/components/Login.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '*',
    name: 'home',
    component: Home,
    beforeEnter: (to, from, next) => {
      if (store.state.user.logged) {
        next()
      } else (
        next('/login')
      )
    }
  }
]

const router = new VueRouter({
  routes,
  mode: 'history'
})

router.b

export default router
