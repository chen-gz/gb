<template>
  <v-container class="d-flex fill-height">
    <v-row class=" d-flex fill-height">
      <v-col cols="12">
        <div id="code-editor" style="width: 100%; border: 1px solid #ccc; "
             class="fill-height"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import {onMounted, onUnmounted, ref} from "vue";
import {useRouter} from "vue-router";
import {getPostV4, savePost, UploadFile, V4PostData} from "@/apiv4";
import * as monaco from "monaco-editor";

const route = useRouter();
let url = route.currentRoute.value.params.url as string || "";
let post = ref({} as V4PostData)
window.addEventListener('keydown', handleKeyDown)


var editor_show = ref("content")

function handleKeyDown(event: KeyboardEvent) {
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault()
    if (editor_show.value === "content") {
      post.value.content = editor.getValue()
    } else {
      post.value = JSON.parse(editor.getValue())
    }
    savePost(post.value)
    console.log("save post: ", editor.getValue())
  } else if (event.ctrlKey && event.key === 'e') {
    event.preventDefault()
    if (editor_show.value === "content") {
      editor_show.value = "meta"
      editor.setValue(JSON.stringify(post.value, null, 4))
      // @ts-ignore
      monaco.editor.setModelLanguage(editor.getModel(), "json")
      editor.updateOptions({wordWrap: "off"})
    } else {
      editor_show.value = "content"
      editor.setValue(post.value.content)
      // @ts-ignore
      monaco.editor.setModelLanguage(editor.getModel(), "markdown")
      editor.updateOptions({wordWrap: "on"})
    }
  }
}

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})


const fileUploadArea = document.documentElement;

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
    UploadFile(file, post.value.id);
  }
});


const code = ref("console.log('Hello, world!');");
let editor: monaco.editor.IStandaloneCodeEditor;

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

// import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'

onMounted(async () => {
  await getPostV4(url, false).then(
      (response) => {
        post.value = response.post
      }
  )
  const editorElement = document.getElementById("code-editor");
  if (editorElement === null) return

  self.MonacoEnvironment = {
    getWorker(_, label) {
      if (label === 'json') return new jsonWorker()
      if (label === 'css' || label === 'scss' || label === 'less') return new cssWorker()
      if (label === 'html' || label === 'handlebars' || label === 'razor') return new htmlWorker()
      if (label === 'typescript' || label === 'javascript') return new tsWorker()
      return new editorWorker()
    }
  }

  editor = monaco.editor.create(editorElement, {
    value: "function hello() {\n\talert('Hello world!');\n}",
    language: 'javascript',
    wordWrap: "on",
  })
  editor.setValue(post.value.content)
  // @ts-ignore
  monaco.editor.setModelLanguage(editor.getModel(), "markdown")
});


</script>

<style scoped>
#code-editor {
}
</style>

