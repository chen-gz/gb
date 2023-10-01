<template>
    <v-text-field
        class="mt-5 mr-5"
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
        density="compact"
        variant="outlined"
    ></v-text-field>
    <v-data-table :items="posts" :headers="headers" :search="search">
        <template v-slot:[`item.actions`]="{ item }">
            <v-icon size="small" class="me-2" @click="editItem(item.raw)">
                mdi-pencil
            </v-icon>
            <v-icon color="red" size="small" @click="deleteItem(item.raw)"> mdi-delete</v-icon>
        </template>

    </v-data-table>
</template>

<style scoped>

</style>

<script setup lang="ts">
import {ref, watch} from "vue";
import {useRouter} from "vue-router";
import router from "@/router";
import {SearchPostsRequestV4, searchPostsV4, V4PostData} from "@/apiv4";

let route = useRouter();
let type = ref("");
if (typeof route.currentRoute.value.params.type !== "undefined") {
    type.value = route.currentRoute.value.params.type as string
}
const headers = ref([
    {title: 'Title', align: 'start', key: 'title'},
    {title: 'Author', align: 'end', key: 'author'},
    {title: 'Tags', align: 'end', key: 'tags'},
    {title: 'Category', align: 'end', key: 'category'},
    {title: 'Actions', align: 'end', key: 'actions', sortable: false}
] as any[])

let posts = ref([] as V4PostData[])
let search = ref("")


function editItem(item: V4PostData) {
    router.push("/posts/edit/" + item.url)
}

function deleteItem(item: V4PostData) {
    console.log(item)
}

async function init() {
    let pa = {} as SearchPostsRequestV4
    pa.sort = "update_time desc, id desc"
    if (type.value == "publish") {
        pa.is_deleted = false
        pa.is_draft = false
    }
    if (type.value == "deleted") {
        pa.is_deleted = true
    }
    if (type.value == "draft") {
        pa.is_draft = true
    }
    searchPostsV4(pa).then((response) => {
        // pagination_len.value = Math.ceil(response.number_of_posts / itemsPerPage.value)
        posts.value = response.posts
        // console.log(response)
        // console.log(response.number_of_posts)
    })
}

init()

watch(router.currentRoute, (old, newe) => {
    type = ref("");
    if (typeof route.currentRoute.value.params.type !== "undefined") {
        type.value = route.currentRoute.value.params.type as string
    }
    init()
    console.log("type changed")
    // update the page num
})

</script>
