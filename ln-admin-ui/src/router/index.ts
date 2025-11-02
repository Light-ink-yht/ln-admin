import { createRouter, createWebHistory } from 'vue-router'
import ThemeDemo from '@/views/ThemeDemo.vue'
import Login from '@/views/login/Login.vue'
import Layout from '@/views/layout/layout.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'Layout',
            component: Layout,
        },
        {
            path: '/login',
            name: 'Login',
            component: Login,
        },
        {
            path: '/theme-demo',
            name: 'ThemeDemo',
            component: ThemeDemo,
        },
    ],
})

export default router
