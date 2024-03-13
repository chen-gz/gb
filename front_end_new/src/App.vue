<script setup>
import Sidebar from "@/components/sidebar.vue";
import PostPage from "@/components/main-wrapper/content/post.vue";
import {watch} from "vue";
// nextTick call tocbot
import {nextTick} from 'vue'
import Category from "@/components/main-wrapper/content/category.vue";
import MainWrapper from "@/components/main-wrapper.vue";
import Sidebar_bottom from "@/components/sidebar/sidebar_bottom.vue";

watch(() => {
});

</script>

<template>
    <aside id="app-sidebar">
        <sidebar></sidebar>
    </aside>
    <main>
        <main-wrapper></main-wrapper>
    </main>
</template>

<style lang="sass">
@import "assets/variables"
main
    z-index: 2
    flex-grow: 0
    flex-shrink: 0
    overflow-y: auto
    overflow-x: hidden
    word-wrap: anywhere
    width: calc(100% - $sidebar-width)
    height: 100%

#app-sidebar
    width: $sidebar-width
    height: 100%
    flex-grow: 1
    flex-shrink: 0
    overflow: hidden
</style>
