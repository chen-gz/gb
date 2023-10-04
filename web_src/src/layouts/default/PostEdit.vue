<template>
    <v-app>
        <v-navigation-drawer v-model="drawer" permanent floating location="right" width="300">
            <v-list h-75>
                <v-list-item title="Post Url">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.url"/> </v-list-item>
                <v-list-item title="Tags">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.tags"/>
                </v-list-item>
                <v-list-item title="Categories">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.category"/>
                </v-list-item>
                <v-list-item title="Author">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.author"/>
                </v-list-item>
                <v-list-item title="summary">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.summary"/>
                </v-list-item>
                <v-list-item title="draft">
                    <v-checkbox v-model="post.is_draft" dense/>
                </v-list-item>
                <v-list-item title="Cover Image">
                    <v-text-field h-75 density="compact" variant="outlined" v-model="post.cover_image"/>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar>
            <v-btn icon="mdi-arrow-left" @click="route.back()"/>
            <v-toolbar-title>
                <v-text-field
                    prepend-icon="mdi-pencil"
                    v-model="post.title"
                    flat
                    hide-details
                    variant="solo"
                    style="font-family: 'JetBrains Mono', monospace;"
                    class="post-input-area"
                ></v-text-field>

            </v-toolbar-title>

            <v-spacer></v-spacer>

            <v-tooltip text="Preview" location="bottom">
<!--                <template v-slot:activator="{ props }">-->
<!--                    <v-btn v-bind="props" icon="mdi-eye" @click="toggleRenderedPreview"/>-->
<!--                </template>-->
            </v-tooltip>

            <v-tooltip text="Save" location="bottom">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" icon="mdi-content-save"/>
                </template>
            </v-tooltip>
            <v-tooltip text="delete" location="bottom">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" color="red" icon="mdi-delete" @click="deletePostBtn"/>
                </template>
            </v-tooltip>
            <v-tooltip text="Post Settings" location="bottom">
                <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" icon @click="drawer = !drawer">
                        <v-icon>mdi-cog</v-icon>
                    </v-btn>
                </template>
            </v-tooltip>
        </v-app-bar>

        <v-main>
            <v-container style="height: 100%">
                <v-row justify="center" class="fill-height">
                    <v-col cols="{{editor_cols}}" class="fill-height">
                        <v-card class="mt-5 fill-height" ref="textarea">
                            <textarea class="fill-height" style="width: 100%" v-model="post.content">

                            </textarea>
                        </v-card>
                    </v-col>

                </v-row>

            </v-container>

        </v-main>

    </v-app>
</template>

<script lang="ts" setup>
import {nextTick, onMounted, onUnmounted, ref, watch} from "vue";
import {useRouter} from "vue-router";


const drawer = ref(false)
const showRendered = ref(false)
const editor_cols = ref(6) // this will be change when rendered preview is toggled
const editor_rows = ref(20)
// @ts-ignore
const textarea = ref(null)
// console.log(textarea)
onMounted(() => {
// @ts-ignore
    console.log(textarea.value.$el.clientHeight)
// @ts-ignore
    let heigh = textarea.value.$el.clientHeight
// @ts-ignore
    console.log(textarea.value.$el.clientWidth)

    console.log(heigh)
    editor_rows.value = Math.floor((heigh - 100) / 24)
    // editor_rows.value = 20
    console.log(editor_rows.value)
    console.log(heigh)

});
const route = useRouter();
let url = route.currentRoute.value.params.url as string || "";

let post = ref({} as V4PostData)
// let meta = ref({} as PostDataV3Meta);
// let content = ref({} as PostDataV3Content);
// let comment = ref({} as PostDataV3Comment);

getPostV4(url, false).then(
    (response) => {
        post.value = response.post
        // meta.value = response.post.meta
        // content.value = response.post.content
        // comment.value = response.post.comment
    }
)


// ctrl + s to save post
window.addEventListener('keydown', handleKeyDown)

function handleKeyDown(event: KeyboardEvent) {
    if (event.ctrlKey && event.key === 's') {
        event.preventDefault()
        savePost(true)
    }
}

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
})


async function savePost(redirect: boolean) {
    // ps.updated_at = new Date()
    let params = {} as V4PostData;
    params.id = post.value.id
    console.log(params)
    console.log(params)
    // params.content = content.value
    // params.content_update = true
    // params.comment = comment.value
    // params.comment_update = true
    updatePostV4(params).then(
        (response) => {
            if (response.status == "success") {
                showSuccess("Post saved")
                if (redirect)
                    route.push({path: '/posts/edit/' + response.url})
            } else {
                showError("Failed to save post")
            }
        }
    )
}

// let myInt = setInterval(savePost, 10000)
// onUnmounted(() => {
//     clearInterval(myInt);
// })


import {getPostV4, showError, showSuccess, updatePostV4, V4PostData} from "@/apiv4";

function deletePostBtn() {
    // meta.value.is_deleted = true
    post.value.is_deleted = true
    savePost(false)
    route.go(-1)
}


</script>

<style>
textarea {
    font-family: "JetBrains Mono", monospace;
    font-size: 16px;
line-height: 30px;
}
</style>

