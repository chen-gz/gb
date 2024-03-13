import { createRouter, createWebHistory } from 'vue-router';
import PostPage from "@/components/main-wrapper/content/post.vue";
import Tags from "@/components/main-wrapper/content/tags.vue";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: PostPage
    },
    {
        path: '/tags/',
        name: 'Tags',
        component: Tags
    },
];

const router = createRouter({
    // history: createWebHistory(process.env.BASE_URL),
    history: createWebHistory(),
    routes
});

export default router;
