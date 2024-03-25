<script lang="ts" setup>
import {searchPostsV4, SearchPostsRequestV4} from "../../../apiv4.js";
import {ref} from "vue";

let param = {} as SearchPostsRequestV4
param.limit = {start: 0, size: 5}
param.sort = "updated_at DESC"
let posts = ref([])
searchPostsV4(param).then((res) => {
    posts.value = res.posts
})

</script>

<template>
    <div id="right-panel-wrapper-inner">
        <div id="recently-updated">
            <h4>Recently Updated</h4>
            <ul>
                <li v-for
                        ="post in posts" :key="post.id">
                    <a :href="'/post/' + post.url">{{ post.title }}</a>
                </li>
            </ul>
        </div>
        <div id="trending Tags" >
            <h4>Trending Tags</h4>
        </div>
        <div id="toc-wrapper">
            <h4> Contents</h4>
            <div id="toc"></div>
        </div>
    </div>
</template>

<style lang="sass" scoped>
#right-panel-wrapper-inner
    display: flex
    flex-grow: 0
    flex-shrink: 0
    padding-top: 48px
    flex-direction: column
    // add vertical spacing between each items
    > div
        margin-bottom: 20px

#recently-updated
    margin-bottom: 20px
    ul
        margin:  0 0 10px 0
        //margin-bottom: 10px
        font-size: 13.6px
        line-height: 17px
        list-style: none
        li
            margin-bottom: 6.4px
            margin-top: 6.4px
            text-align: left
            box-sizing: content-box
        a
            padding-left: 17px
            padding-bottom: 3.2px
            padding-top: 3.2px




</style>