<template>
  <!--        element inside this block from left to right -->
  <div class="post_page">
    <v-container>
      <v-row>
        <v-col cols="10" sm="12" md="10">
          <div style="justify-content: center">
            <h1 class="post_title" v-html="post.title" style="font-size: 40px; font-family: 'Noto Serif SC', serif;"/>
            <div class="post_content">
              <!--              close to title far away from content -->
              <div style="font-size: 15px; font-family: 'Noto Serif SC', serif; color: #6c6c6c;
                margin-top: -10px; margin-bottom: 20px;">
                <span>Author：{{ post.author }}</span>
                <span class="ml-3">Update: {{ formatDate(post.updated_at) }}</span>
                <span class="ml-1">|</span>

                <span class="ml-3" style="color: #000000; cursor: pointer;" @click="sharepost()">Share</span>
                <span class="ml-3" style="color: #000000;  cursor: pointer;" @click="deletePost(post)">Delete</span>
                <router-link :to="'/posts/edit/' + post.url">
                        <span class="ml-3" style="color: #000000;">
                        Edit
                        </span>
                </router-link>


              </div>
              <div v-html="post_content"></div>
            </div>
          </div>
        </v-col>
        <v-col cols="2" sm="0" md="2">
          <div v-if="post_toc.length >0 " class="post_toc" v-html="post_toc" style="top: 10%; wrap-option: wrap;"></div>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {deletePost, formatDate, getPostV4, showSuccess, V4PostData} from "@/apiv4";

const route = useRouter();
console.log(route.currentRoute.value.params.url)
let url: string = "";
if (typeof route.currentRoute.value.params.url !== "undefined") {
  url = route.currentRoute.value.params.url as string
}

let post = ref({} as V4PostData);
let post_content = ref("");
let post_toc = ref("");

function sharepost() {
  let url = window.location.href
  navigator.clipboard.writeText(url).then(() => {
    showSuccess("url copied to clipboard")
  })
}


getPostV4(url, true).then((response) => {
  console.log(response)
  post.value = response.post
  post_content.value = response.html
  // get toc from post content they are surrended by <nav> tag
  // find the first <nav> tag
  let nav_start = post_content.value.indexOf("<nav")
  let nav_end = post_content.value.indexOf("</nav>")
  if (nav_start == -1 || nav_end == -1) {
    post_toc.value = ""
  } else {
    post_toc.value = post_content.value.substring(nav_start, nav_end + 6)
    post_content.value = post_content.value.substring(nav_end + 6)

  }
  post_toc.value = post_toc.value.substring(16, post_toc.value.length - 19)
  console.log(post_toc.value)
})

// watch(post_content, (old, newe) => {
//   console.log("post changed")
//   nextTick(() => {
//     // @ts-ignore
//     window.MathJax.Hub.Queue(["Typeset", window.MathJax.Hub]);
//     // @ts-ignore
//     window.hljs.highlightAll()
//   })
// });

</script>


<style scoped>
.post_content {
  font-size: 20px;
  line-height: 40px;
  font-family: "Noto Serif SC", serif;
  align-self: center;
  align-content: center;
  justify-content: center;
  justify-self: center;
  wrap-option: wrap;
}
</style>

