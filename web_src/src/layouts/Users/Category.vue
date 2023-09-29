<script setup lang="ts">
import {ref} from "vue";
import {getDistinct, GetDistinctResponse, PostDataV3Meta, SearchPostsRequestV3, searchPostsV3} from "@/apiv2";

export interface Category {
    name: string
    count: number
    posts: PostDataV3Meta[]
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
        let param = {} as SearchPostsRequestV3
        param.limit = {start: 0, size: 5}
        param.categories = cate.name
        searchPostsV3(param).then(
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
