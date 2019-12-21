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
    path: '/expenses/:id',
    name: 'expense',
    component: () => import('../views/Expense.vue'),
    props: true
  },
  {
    path: '/documents',
    name: 'documents',
    component: () => import('../views/Documents.vue')
  },
  {
    path: '/documents/:id',
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
