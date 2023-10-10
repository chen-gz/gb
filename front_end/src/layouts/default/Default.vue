<!--this is the home page -->
<template>
  <v-app>
    <default-bar/>
    <v-main style="max-width: 1400px; width: 100vw;" class="align-self-center">
      <Suspense>

      <router-view/>
      </Suspense>
    </v-main>


  </v-app>
</template>

<script lang="ts" setup>
import DefaultBar from './AppBar.vue'
import {SearchPostsRequestV4, searchPostsV4, V4PostData} from "@/apiv4";
import {ref} from "vue";
// import {SearchPostsRequestV4} from "@/apiv4";

const itemPerPage = 10
let blog_list = ref([] as V4PostData[])

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

// function init() {
//   let pa = {} as SearchPostsRequestV4
//   pa.rendered = true
//   pa.counts_only = true
//   pa.limit = {start: 0, size: 10}
//   searchPostsV4(pa).then((response) => {
//     len.value = Math.ceil(response.number_of_posts / itemPerPage)
//   })
//   getPostPage(1)
// }
//
// init()
</script>
