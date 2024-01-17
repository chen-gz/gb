<script setup lang="ts">

</script>

<template>
    <!--        element inside this block from left to right -->
    <div class="post_page">
        <div class="d-flex">
            <div class="scrolling-content ml-5">
                <div style="justify-content: center">
                </div>
            </div>

        </div>
    </div>
</template>

<script setup lang="ts">
import {nextTick, ref, watch} from "vue";
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

watch(post_content, () => {
    console.log("post changed")
    nextTick(() => {
        // @ts-ignore
        window.MathJax.Hub.Queue(["Typeset", window.MathJax.Hub]);
        // @ts-ignore
        window.hljs.highlightAll()
    })
});

</script>


<style scoped>
.post_content {
    /*font-size: 15px;*/
    line-height: 25px;
    font-family: "Noto Serif SC", serif;
    align-self: center;
    align-content: center;
    justify-content: center;
    justify-self: center;
    wrap-option: wrap;
}

.fixed-sidebar {
    display: flex;
    /*background-color: yellow;*/
    margin-right: 5px;
    min-width: 200px;

}

.scrolling-content {
    display: flex;
    overflow-y: auto;
    /*background-color: blue;*/
    height: calc(100vh - 64px);
    /*width: 100%;*/
    scrollbar-width: none;
}

/* hide scrollbar for chrome, safari and opera */
.scrolling-content::-webkit-scrollbar {
    display: none;
}


</style>

