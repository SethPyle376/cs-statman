import Vue from 'vue'
import VueRouter from 'vue-router'
import App from '../components/Layout'
import Upload from '../components/Upload'

Vue.use(VueRouter)

const routes = [
  {
    path: '/app',
    name: 'App',
    component: App,
    children : [
      {
        path: ':userID/dashboard',
        component: () => import('../components/Dashboard')
      },
      {
        path: 'home',
        component: () => import('../components/Home')
      },
      {
        path: 'upload',
        component: Upload
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router