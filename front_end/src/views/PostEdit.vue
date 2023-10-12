<template>
  <v-container class="d-flex fill-height">
<!--    <v-app-bar class="d-flex">-->
<!--      <v-btn icon="mdi-arrow-left" @click="route.back()"/>-->
<!--      <v-toolbar-title>-->
<!--        <v-text-field-->
<!--            prepend-icon="mdi-pencil"-->
<!--            v-model="post.title"-->
<!--            flat-->
<!--            hide-details-->
<!--            variant="solo"-->
<!--            style="font-family: 'JetBrains Mono', monospace;"-->
<!--            class="post-input-area"-->
<!--        ></v-text-field>-->
<!--      </v-toolbar-title>-->
<!--      <v-spacer></v-spacer>-->
<!--      <v-tooltip text="Preview" location="bottom">-->
<!--        &lt;!&ndash;                <template v-slot:activator="{ props }">&ndash;&gt;-->
<!--        &lt;!&ndash;                    <v-btn v-bind="props" icon="mdi-eye" @click="toggleRenderedPreview"/>&ndash;&gt;-->
<!--        &lt;!&ndash;                </template>&ndash;&gt;-->
<!--      </v-tooltip>-->
<!--      <v-tooltip text="Save" location="bottom">-->
<!--        <template v-slot:activator="{ props }">-->
<!--          <v-btn @click="savePost(post)" v-bind="props" icon="mdi-content-save"/>-->
<!--        </template>-->
<!--      </v-tooltip>-->
<!--      <v-tooltip text="delete" location="bottom">-->
<!--        <template v-slot:activator="{ props }">-->
<!--          <v-btn v-bind="props" color="red" icon="mdi-delete" @click="deletePost(post)"/>-->
<!--        </template>-->
<!--      </v-tooltip>-->
<!--      <v-tooltip text="Post Settings" location="bottom">-->
<!--        <template v-slot:activator="{ props }">-->
<!--          <v-btn v-bind="props" icon @click="drawer = !drawer">-->
<!--            <v-icon>mdi-chevron-down</v-icon>-->
<!--          </v-btn>-->
<!--        </template>-->
<!--      </v-tooltip>-->
<!--    </v-app-bar>-->
    <v-container v-show="drawer" >
      <v-row>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Tags" v-model="post.tags" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Categories" v-model="post.category" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Post Url" v-model="post.url" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Author" v-model="post.author" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Summary" v-model="post.summary" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field label="Cover Image" v-model="post.cover_image" variant="solo"/>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-checkbox label="Draft" v-model="post.is_draft" variant="solo"/>
        </v-col>

      </v-row>

    </v-container>
    <v-row class=" d-flex fill-height">
      <v-col cols="12">
        <div id="code-editor" style="width: 100%;"
        class="fill-height"/>
      </v-col>
    </v-row>


  </v-container>
</template>

<script lang="ts" setup>
import {onMounted, onUnmounted, ref} from "vue";
import {useRouter} from "vue-router";
import {deletePost, getPostV4, savePost, showSuccess, UploadFile, V4PostData} from "@/apiv4";
import * as monaco from "monaco-editor";

const drawer = ref(false)
const route = useRouter();
let url = route.currentRoute.value.params.url as string || "";
let post = ref({} as V4PostData)


// ctrl + s to save post
window.addEventListener('keydown', handleKeyDown)

function handleKeyDown(event: KeyboardEvent) {
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault()
    post.value.content = editor.getValue()
    savePost(post.value)
    console.log("save post: ", editor.getValue())
  }
}

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})



// function deletePostBtn() {
//   // meta.value.is_deleted = true
//   post.value.is_deleted = true
//   savePost(false)
//   route.go(-2)
// }


const fileUploadArea = document.documentElement;

// console.log("fileList: " + fileList)
fileUploadArea.addEventListener("dragover", (event) => {
  event.preventDefault();
  fileUploadArea.classList.add("dragover");
});
fileUploadArea.addEventListener("dragleave", () => {
  fileUploadArea.classList.remove("dragover");
});
fileUploadArea.addEventListener("drop", (event) => {
  const fileList = document.querySelector("#file-list ul");

  event.preventDefault();
  fileUploadArea.classList.remove("dragover");

  // @ts-ignore
  const files = event.dataTransfer.files;
  for (const file of Array.from(files)) {
    const li = document.createElement("li");
    li.textContent = file.webkitRelativePath || file.name;
    // @ts-ignore
    fileList.appendChild(li);
    // uploadFileToServer(file);
    UploadFile(file, post.value.id);
  }
});


const code = ref("console.log('Hello, world!');");
let editor : monaco.editor.IStandaloneCodeEditor;

// const initializeEditor = (editorElement) => {
function initializeEditor(editorElement: HTMLElement) {
  editor = monaco.editor.create(editorElement, {
    value: post.value.content,
    language: "markdown",
    theme: "vs-light",
    wordWrap: "on",

  });

  editor.onDidChangeModelContent(() => {
    code.value = editor.getValue();
  });

}

onMounted(async () => {
  await getPostV4(url, false).then(
      (response) => {
        post.value = response.post
      }
  )
  const editorElement = document.getElementById("code-editor");
  if (editorElement) {
    initializeEditor(editorElement);
    editor.setValue(post.value.content);
  }
});


</script>

<style scoped>
#code-editor {
  width: 100%;
  border: 1px solid #ccc;
}
</style>

