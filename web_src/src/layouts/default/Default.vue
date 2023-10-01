<template>
    <v-app>
        <default-bar @search="searchCallBack"/>
        <v-main>
            <div style="width: 98%; margin-left: auto;margin-right: auto">
                <!--                    <div style="width: 90%; margin-left: auto;margin-right: auto">-->
                <router-view v-if="!showSearch"/>
                <div v-else>
                    <search-result :blog_list="blog_list"/>
                </div>
            </div>
        </v-main>
        <footer>
            <v-col class="text-center" cols="12">
                Power by Eta Blog - &copy {{ new Date().getFullYear() + ' Guangzong' }}
                <br> All rights reserved.
            </v-col>
        </footer>

    </v-app>
</template>

<script lang="ts" setup>
import DefaultBar from './AppBar.vue'
import {onUnmounted, ref, watch} from "vue";
import searchResult from "@/layouts/Users/searchResult.vue";
import router from "@/router";
import {SearchPostsRequestV4, searchPostsV4, V4PostData} from "@/apiv4";

let showSearch = ref(false)
let blog_list = ref([] as V4PostData[])

function searchCallBack(search_text: string) {
    showSearch.value = true
    console.log(showSearch.value)
    console.log(search_text)
    let param = {} as SearchPostsRequestV4
    param.content = search_text
    param.limit = {start: 0, size: 20}
    searchPostsV4(param).then(
        (response) => {
            blog_list.value = response.posts || []
            console.log(blog_list.value)
        }
    )
}

// capture for esc key to close search
window.addEventListener('keydown', handleKeyDown)
// if router change, close the search

watch(() => router.currentRoute.value, (to, from) => {
    showSearch.value = false
})


function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
        showSearch.value = false
    }

}

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
})

</script>

<style>
#page {
    width: 100vw;
    min-height: 100vh;
}

#app {
    max-width: 1440px;
    margin: 0 auto;
}

a {
    color: #1976d2;
    text-decoration: none;
    transition: color 0.15s ease;
}

blockquote {
    border-left: 4px solid #e0e0e0;
    margin-left: 0;
    padding-left: 1em;
    color: #616161;
}

@media screen and (max-width: 600px) {
    .post_toc {
        display: none;
    }
}

ol, ul {
    padding-left: 2em;
}


</style>
