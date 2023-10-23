<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {formatDate, SearchPostsRequestV4, SearchPostsResponseV4, searchPostsV4} from "@/apiv4";


const props = defineProps<{
    searchParam: SearchPostsRequestV4
}>();
console.log(props.searchParam)

const route = useRouter();
console.log(route.currentRoute.value.params.id)
let params = {} as SearchPostsRequestV4
console.log("props.searchParam: ", props.searchParam)
if (props.searchParam != undefined) {
    params = props.searchParam as SearchPostsRequestV4
}
console.log(params)
// let params = {} as SearchPostsRequestV3
var res = ref({} as SearchPostsResponseV4)
params.sort = "created_at DESC"
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
    <v-table class="ml-5 mr-5">
        <tbody>
        <tr v-for="(item, index) in res.posts" :key="index" style="font-family: 'JetBrains Mono', monospace">
            <td>
                <router-link :to="'/posts/' + item.url">
                    <span> {{ item.title }} </span>
                </router-link>
            </td>
<!--            align right for data -->
            <td class="text-right"
            >{{ formatDate(item.created_at) }}</td>
        </tr>
        </tbody>
    </v-table>
</template>
