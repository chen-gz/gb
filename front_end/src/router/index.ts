// Composables
import {createRouter, createWebHistory} from 'vue-router'

const routes = [
    {
        path: '/',
        component: () => import('@/layouts/default/Default.vue'),
        children: [
            {path: '', name: 'Home', component: () => import(/* webpackChunkName: "home" */ '@/views/Home.vue'),},
            {path: 'about', name: 'About', component: () => import('@/views/About.vue'),},
            {path: 'tags', name: 'Tags', component: () => import('@/views/Tags.vue'),}
        ]
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Login', component: () => import('@/views/Login.vue'),}]
    },
    {
        path: '/posts',
        component: () => import('@/layouts/default/Default.vue'),
        children: [
            {path: '', name: 'Posts', component: () => import('@/views/Lists.vue'),},
            {path: ':url', name: 'Post', component: () => import('@/views/PostPage.vue'),},
            {path: 'edit/:url', name: 'PostEdit', component: () => import('@/views/PostEdit.vue'),},
        ]
    },
    {
        path: '/photos',
        component: () => import('@/layouts/default/Default.vue'),
        children: [
            {path: '', name: 'Photos', component: () => import('@/views/photoList.vue'),},
        ]
    },
    {
        path: '/tags',
        component: () => import('@/layouts/default/Default.vue'),
        children: [
            {path: ':tag_name', name: 'TagList', component: () => import('@/views/Lists.vue'),},
        ]
    }
]

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
})

export default router
