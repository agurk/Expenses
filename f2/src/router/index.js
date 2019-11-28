import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/expenses'
  },
  {
    path: '/expenses',
    name: 'expenses',
    component: () => import('../views/Expenses.vue')
  },
  {
    path: '/expense/:id',
    name: 'expense',
    component: () => import('../views/Expense.vue'),
    props: true
  },
  {
    path: '/document/:id',
    name: 'document',
    component: () => import('../views/Document.vue'),
    props: true
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
