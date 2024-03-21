<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {formatDate, SearchPostsRequestV4, SearchPostsResponseV4, searchPostsV4} from "/apiv4";

// let searchParam = SearchPostsRequestV4
const route = useRouter();
// url is tags/:id
let searchParam = {
    tags: route.currentRoute.value.query.tag,
    categories: route.currentRoute.value.query.cate,
    sort: "created_at DESC"
} as SearchPostsRequestV4;
let params = searchParam
// let params = {} as SearchPostsRequestV3
var res = ref({} as SearchPostsResponseV4)
// params.sort
searchPostsV4(params).then((response) => {
    res.value = response
    console.log(res.value)
    if (res.value.number_of_posts == 0) {
        res.value.posts = []
    }
    for (let i = 0; i < res.value.posts.length; i++) {
        res.value.posts[i].created_at = new Date(res.value.posts[i].created_at);
    }
    console.log(res.value)
})

</script>

<template>
    <h1>
        Tags Lists
    </h1>
    <ul>
        <li v-for="result in res.posts" :key="result.id">
            <a href="#" @click="route.push('/post/' + result.url)">
                {{result.title}}
            </a>
            <span class="dashs">
<!--            add dash to fill the space -->
             </span>

            <time>
                {{result.updated_at.toString()}}
            </time>
        </li>
    </ul>
</template>

<style scoped lang="sass">

.dashs
    margin-left: 10px
    flex-grow: 1
    border-bottom: 1px dotted currentColor // ; /* Use a dotted border bottom */

li
    display: inline-flex
    width: 100%
    margin-bottom: 20px
    a
        font-size: 20px
        color: #000
        text-decoration: none
        &:hover
            color: #000
            text-decoration: underline
h1
    margin-top: 50px
    margin-bottom: 20px
</style>