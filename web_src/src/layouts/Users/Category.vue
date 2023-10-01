<script setup lang="ts">
import {ref} from "vue";
// import {GetDistinctResponse, PostDataV3Meta} from "@/apiv2";
import {searchPostsV4, SearchPostsRequestV4, getDistinct, V4PostData, GetDistinctResponse} from "@/apiv4";


export interface Category {
    name: string
    count: number
    posts: V4PostData[]
}

var cates = ref({} as GetDistinctResponse)
let catesShow = ref([] as Category[])
getDistinct("category").then((response) => {
    cates.value = response
    for (let i = 0; i < cates.value.values.length; i++) {
        let cate = cates.value.values[i]
        let cateShow = {} as Category
        if (cate === "") {
            continue
        }
        cateShow.name = cate
        cateShow.count = 0
        cateShow.posts = []
        catesShow.value.push(cateShow)
    }
    // get 5 posts for each category
    for (let i = 0; i < catesShow.value.length; i++) {
        let cate = catesShow.value[i]
        let param = {} as SearchPostsRequestV4
        param.limit = {start: 0, size: 5}
        param.categories = cate.name
        searchPostsV4(param).then(
            (response) => {
                cate.posts = response.posts || []
            }
        )
    }

    console.log(cates.value)
})


</script>

<template>
    <div v-for="(item, index) in catesShow" :key="index">
        <div>
            <h2>{{ item.name }}</h2>
            <ol>
                <li v-for="(post, index2) in item.posts" :key="index2">
                    {{ post.title }}
                </li>
            </ol>
        </div>
    </div>
</template>

<style scoped>

</style>
