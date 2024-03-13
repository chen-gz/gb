<template>
  <div class="post_page">
    <h1> {{ post.title }}</h1>
    <div class="post-meta text-muted">
      <div class="post-header">
        <span> Posted on {{
            new Date(post.created_at).toLocaleDateString("en-US", {
              year: 'numeric',
              month: 'short',
              day: 'numeric'
            })
          }}</span>
        <span class="before_dot"> Updated on {{
            new Date(post.updated_at).toLocaleDateString("en-US", {
              year: 'numeric',
              month: 'short',
              day: 'numeric'
            })
          }}</span>
        <!--        <span>Updated {{post.updated_at}}</span>-->
        <div class="sub-meta">
          <span class="author">{{ post.author }}</span>
<!--          add edit button-->
          <button>edit</button>
          <button>delete</button>
        </div>
      </div>
    </div>
    <div id="markdown_content" v-html="post_content"></div>
  </div>
</template>

<style lang="sass" scoped>

//@import "@/assets/variables"
//button
//  margin-left: 1rem
//  //margin-right: 2rem
//  //background-color: inherit
//  border: none
//  text-decoration: underline
//  text-underline: #0d6efd
//.post-header
//  margin-bottom: 2rem
//  margin-top: 2rem
//
//.before_dot::before
//  content: "•"
//  padding-left: .25rem
//  padding-right: .25rem
//
//span
//  display: inline-block
//  color: var(--text-muted-color)
//
//.author
//  color: var(--text-muted-color)

</style>

<script lang="ts" setup>
import {nextTick, onMounted, ref, watch, watchEffect} from "vue";
import {getPostV4, V4PostData} from "../../../../apiv4";

let url = "253"
let post = ref({} as V4PostData);
// let post = ref();
let post_content = ref("");
let post_toc = ref("");

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

watch(post_content, () => {
  console.log("post changed")
  nextTick(() => {
    // @ts-ignore
    // window.MathJax.Hub.Queue(["Typeset", window.MathJax.Hub]);
    window.MathJax.startup.defaultReady();

// @ts-ignore
    window.hljs.highlightAll()
    if (typeof window.tocbot !== 'undefined') {
      window.tocbot.init({
        tocSelector: '#toc',
        contentSelector: '#markdown_content',
        ignoreSelector: '[data-toc-skip]',
        headingSelector: 'h2, h3, h4',
        orderedList: false,
        scrollSmooth: false
      });
      // tocbot.run();
      window.tocbot.refresh();
      console.log('tocbot is defined');
    } else {
      console.error('tocbot is not defined');
    }
  })
});
// watch for scroll  (event)
// const handleScroll = () => {
//   // Your scroll handling logic
//   console.log(window.scrollY);
//   window.tocbot.refresh();
// };


// onMounted(() => {
//   window.addEventListener('scroll', handleScroll);
// });


</script>



