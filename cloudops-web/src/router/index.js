import { createRouter, createWebHistory } from 'vue-router';
const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/index.vue'),
    },
    {
        path: '/',
        component: () => import('@/components/Layout/AppLayout.vue'),
        children: [
            {
                path: '',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/index.vue'),
            },
        ],
    },
];
const router = createRouter({
    history: createWebHistory(),
    routes,
});
export default router;
