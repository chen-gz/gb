<script setup lang="ts">
import axios from "axios";
import {ref} from "vue";
import {useRouter} from "vue-router";
import {formatDate, SearchPostsRequestV3, SearchPostsResponseV3, searchPostsV3} from "@/apiv2";


const props = defineProps<{
    searchParam: SearchPostsRequestV3
}>();
console.log(props.searchParam)

const route = useRouter();
console.log(route.currentRoute.value.params.id)
let params = {} as SearchPostsRequestV3
console.log("props.searchParam: ", props.searchParam)
if (props.searchParam != undefined) {
    params = props.searchParam as SearchPostsRequestV3
}
console.log(params)
// let params = {} as SearchPostsRequestV3
var res = ref({} as SearchPostsResponseV3)
params.sort = "create_time DESC"
searchPostsV3(params).then((response) => {
    res.value = response
    console.log(res.value)
    if (res.value.number_of_posts == 0) {
        res.value.posts = []
    }
    for (let i = 0; i < res.value.posts.length; i++) {
        res.value.posts[i].create_time = new Date(res.value.posts[i].create_time);
    }
    console.log(res.value)
})


</script>

<template>
    <v-table>
        <tbody>
        <tr v-for="(item, index) in res.posts" :key="index">
            <td>
                <router-link :to="'/posts/' + item.url">
                    <span> {{ item.title }} </span>
                </router-link>
            </td>
            <td>{{ formatDate(item.create_time) }}</td>
            <td>
            </td>
        </tr>
        </tbody>
    </v-table>
</template>

<style scoped>
tr {
    font-family: "JetBrains Mono";
}

</style>
