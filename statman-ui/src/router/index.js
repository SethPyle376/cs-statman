import Vue from 'vue'
import VueRouter from 'vue-router'
import App from '../components/Layout'

Vue.use(VueRouter)

const routes = [
  {
    path: '/app',
    name: 'App',
    component: App,
    children : [
      {
        path: 'dashboard/:userID',
        component: () => import('../components/Dashboard')
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router