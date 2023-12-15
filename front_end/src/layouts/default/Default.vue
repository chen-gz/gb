<!--this is the home page -->
<template>
  <v-app style="height: 100vh">
    <default-bar @search="handleSearchEvent"/>
    <v-main class="main-container align-self-center" style=" max-width: 1400px; width: 100vw;">
      <Suspense>
   <router-view v-if="!search"></router-view>
    <div v-else>
      <!-- show search result -->
      <Lists :searchParam="searchParam" :key="list_key"/>
    </div>
      </Suspense>
    </v-main>


  </v-app>
</template>

<script lang="ts" setup>
import DefaultBar from './AppBar.vue'
import {SearchPostsRequestV4, searchPostsV4, V4PostData} from "@/apiv4";
import {ref, watch} from "vue";
import Lists from '@/views/Lists.vue';
import { useRoute } from 'vue-router';
// import {SearchPostsRequestV4} from "@/apiv4";

const itemPerPage = 10
let blog_list = ref([] as V4PostData[])
let search = ref(false)
let searchParam = ref({} as SearchPostsRequestV4)
let list_key = ref(0)

let len = ref(0)

async function getPostPage(page: number) {
  let pa2: SearchPostsRequestV4 = {} as SearchPostsRequestV4
  pa2.sort = "updated_at DESC"
  pa2.rendered = false
  pa2.limit = {start: (page - 1) * itemPerPage, size: itemPerPage}
  searchPostsV4(pa2).then((response) => {
    blog_list.value = response.posts
    console.log("get post page")
    console.log(blog_list.value)
    // console.log(blog_list.value)
  })
}
function handleSearchEvent(str: string) {
  console.log("search function called with param, " + str)
  search.value = true
  searchParam.value.content = str
  list_key.value += 1
}

// add shortcut when esc is pressed
window.addEventListener('keydown', function (e) {
  if (e.key === 'Escape') {
    search.value = false
  }
})
// when the router is changed, reset the search

// Watch for route changes
const route = useRoute();
watch(route, () => {
  // Reset search when the route changes
  search.value = false;
  searchParam.value = {} as SearchPostsRequestV4;
  list_key.value += 1; // Optional: reset the key if needed
});

</script>

<style scoped>
.main-container {
  height: calc(100vh - 64px);
}
</style>
