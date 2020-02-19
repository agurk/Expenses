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
    },
    {
        path: '/analysis',
        name: 'analysis',
        component: () => import('../views/Analysis.vue')
    },
    {
        path: '/assets',
        name: 'assets',
        component: () => import('../views/Assets.vue')
    },
    {
        path: '/search',
        name: 'search',
        component: () => import('../views/Search.vue')
    },
    {
        path: '/config',
        name: 'config',
        component: () => import('../views/Config.vue')
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
