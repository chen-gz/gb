<!--<script setup lang="ts">-->
<!--import {useRouter} from "vue-router";-->
<!--import loader from "@monaco-editor/loader";-->

<!--import {getPostV4, V4PostData} from "../../../../apiv4";-->
<!--import {onMounted} from "vue";-->
<!--// onMounted(async () => {-->
<!--// let router = useRouter()-->
<!--// let url = router.currentRoute.value.params.id as string-->
<!--// let post = {} as V4PostData-->
<!--// await getPostV4(url, false).then(-->
<!--//     (response) => {-->
<!--//         post = response.post-->
<!--//     }-->
<!--// )-->
<!--console.log("post" + post.content)-->
<!--// onMounted(-->
<!--//     async () => {-->
<!--const monaco = await loader.init()-->
<!--const editor = monaco.editor.create(document.getElementById('code_editor'), {-->
<!--    value: post.content,-->
<!--    language: 'markdown',-->
<!--    theme: 'one-light',-->
<!--});-->
<!--console.log("loaded editor")-->
<!--//     }-->
<!--// )-->

<!--// loader.init().then(monaco => {-->
<!--//     const editor = monaco.editor.create(document.getElementById('code_editor'), {-->
<!--//         value: `function x() {-->
<!--//   console.log("Hello world!");-->
<!--// }`,-->
<!--//         // value: post.content,-->
<!--//         language: 'markdown',-->
<!--//         theme: 'one-light',-->
<!--//     });-->
<!--// });-->

<!--</script>-->


<script setup lang="ts">
import loader from "@monaco-editor/loader";
import {useRouter} from "vue-router";
import {getPostV4} from "../../../../apiv4";
import {ref} from "vue";

let route = useRouter()
// let url = router.currentRoute.value.params.id;
// getPostV4(url, false).then(
//     (response) => {
//         console.log("post" + response.post.content)
//         loader.init().then(monaco => {
//             const editor = monaco.editor.create(document.getElementById('code_editor'), {
//                 value: response.post.content,
//                 language: 'markdown',
//                 theme: 'one-light',
//             });
//         });
//     }
// )

let url = route.currentRoute.value.params.id
console.log(url)

let post = ref({} as V4PostData);
// let post = ref();
let post_content = ref("");
let post_toc = ref("");
console.log(url)

getPostV4(url, true).then((response) => {
    console.log(response)
    post.value = response.post
    loader.init().then(monaco => {
        const editor = monaco.editor.create(document.getElementById('code_editor'), {
            value: post.value.content,
            language: 'markdown',
            theme: 'one-light',
            wrappingColumn: 80,
            wordWrap: 'on',
        });
    });
})

</script>
<template>
            <div id="code_editor"></div>
</template>


<style scoped lang="sass">

#code_editor
    width: calc(100% - 2px)
    height: calc(100% - 2px)
    border: 1px solid black
    overflow: hidden
    position: relative
    flex-grow: 1

</style>
