// Composables
import {createRouter, createWebHistory, useRouter} from 'vue-router'
// import {SearchPostsRequestV4, searchPostsV4} from "@/apiv2";

// let searchParams = {} as SearchPostsRequestV4

const routes = [
    {
        path: '/',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Home', component: () => import( '@/layouts/Users/HomePage.vue'),},],
    },
    {
        path: '/posts',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Posts', component: () => import('@/layouts/Users/Lists.vue'),},],
    },
    {
        path: '/posts/:url',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Post', component: () => import('@/layouts/Users/PostPage.vue'),},],
    },
    {
        path: '/posts/edit/:url',
        component: () => import('@/layouts/default/PostEdit.vue')
    },
    {
        path: '/tag',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Tag', component: () => import('@/layouts/Users/Tag.vue'),},
            {
                path: ':id',
                name: 'TagItem',
                props: true,
                component: () => import('@/layouts/Users/Tag.vue'),

            }]
    },
    // {
    //     path: '/tag/:id',
    //     component: () => import('@/layouts/default/Default.vue'),
    //     children: [{path: '', name: 'Tag_with_name', component: () => import('@/layouts/Users/Lists.vue'),},],
    //     props: true,
    // },
    // {
    //     path: '/categories',
    //     component: () => import('@/layouts/default/Default.vue'),
    //     children: [{path: '', name: 'categories', component: () => import('@/layouts/Users/Category.vue'),},],
    // },
    {
        path: '/login',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'Login', component: () => import('@/layouts/Users/Login.vue'),}],
    },
    {
        path: '/admin',
        component: () => import('@/layouts/default/Admin.vue'),
    },
    {
        path: '/admin/posts/:type',
        component: () => import('@/layouts/default/Admin.vue'),
        children: [{path: '', name: 'AdminPosts', component: () => import('@/layouts/Admin/PostList.vue'),}],
    },
    {
        path: '/about',
        component: () => import('@/layouts/default/Default.vue'),
        children: [{path: '', name: 'About', component: () => import('@/layouts/Users/About.vue'),}],
    }

]

const router = createRouter({
    // history: createWebHistory(process.env.BASE_URL),
    history: createWebHistory(),
    // history: createWebHashHistory("/"),
    routes,

})

export default router
