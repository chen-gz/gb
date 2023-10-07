<template>
  <v-container class="d-flex fill-height">
    <v-app-bar class="d-flex">
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
          <v-btn @click="savePost(post)" v-bind="props" icon="mdi-content-save"/>
        </template>
      </v-tooltip>

      <v-tooltip text="delete" location="bottom">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" color="red" icon="mdi-delete" @click="deletePost(post)"/>
        </template>
      </v-tooltip>

      <v-tooltip text="Post Settings" location="bottom">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon @click="drawer = !drawer">
            <v-icon>mdi-chevron-down</v-icon>
          </v-btn>
        </template>
      </v-tooltip>
    </v-app-bar>
    <!--    new textfield that for post settings, click the icon show the fileds -->
    <v-container
      v-show="drawer"
    >
      <v-row>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Tags"
            v-model="post.tags"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Categories"
            v-model="post.category"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Post Url"
            v-model="post.url"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Author"
            v-model="post.author"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Summary"
            v-model="post.summary"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            label="Cover Image"
            v-model="post.cover_image"
            variant="solo"
          />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-checkbox
            label="Draft"
            v-model="post.is_draft"
            variant="solo"
          />
        </v-col>

      </v-row>

    </v-container>

    <div id="file-list">
      <p>Uploaded files:</p>
      <ul></ul>
    </div>


    <v-container style="width: 100%" class="flex-grow fill-height">
      <!--      <v-row justify="center" class="fill-height" style="width: 100%; height: 100%;">-->
      <!--                    <v-col cols="{{editor_cols}}" class="fill-height">-->
      <!--                        <v-card class="mt-5 fill-height" ref="textarea" style="width: 100vw">-->
      <!--                            <textarea class="fill-height" style="width: 100%" v-model="post.content">-->
      <!--                            </textarea>-->
      <!--                          {{post.content}}-->
      <!--      font-family: "JetBrains Mono", monospace;-->
      <!--      font-size: 16px;"-->
      <v-textarea v-model="post.content"
                  ref="textarea"
                  class="fill-height"
                  style="width: 100%; font-family: 'JetBrains Mono', monospace; font-size: 16px;"
                  :rows="editor_rows"
                  :row-height="editor_row_height"
      ></v-textarea>
      <!--                        </v-card>-->
      <!--                    </v-col>-->
      <!--      </v-row>-->
    </v-container>

  </v-container>
</template>

<script lang="ts" setup>
import {nextTick, onMounted, onUnmounted, ref, watch} from "vue";
import {useRouter} from "vue-router";
import {deletePost, savePost} from "@/apiv4";


import {getPostV4, showError, showSuccess, updatePostV4, UploadFile, V4PostData} from "@/apiv4";

const drawer = ref(false)
const showRendered = ref(false)
const editor_cols = ref(6) // this will be change when rendered preview is toggled
const editor_rows = ref(25)
const editor_row_height = ref(24)
// @ts-ignore
const textarea = ref(null)
console.log(textarea)
onMounted(() => {
// @ts-ignore
  console.log(textarea.value.$el.clientHeight)
// @ts-ignore
  let heigh = textarea.value.$el.clientHeight
// @ts-ignore
  console.log(textarea.value.$el.clientWidth)

  console.log(heigh)
  editor_rows.value = Math.floor((heigh) / 25) - 1
  // editor_rows.value = 20
  // console.log(editor_rows.value)
  console.log("editor_rows.value: " + editor_rows.value)
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


//async function savePost(redirect: boolean) {
//  // ps.updated_at = new Date()
//  // get token
//  var token = localStorage.getItem('token')
//  console.log("token: " + token)
//  let params = {} as V4PostData;
//  params = post.value
//
//  updatePostV4(params).then(
//    (response) => {
//      if (response.status == "success") {
//        post.value = response.post
//        showSuccess("Post saved")
//        if (redirect)
//          route.push({path: '/posts/edit/' + post.value.url})
//      } else {
//        showError("Failed to save post")
//      }
//    }
//  )
//}

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

  const files = event.dataTransfer.files;
  for (const file of files) {
    const li = document.createElement("li");
    li.textContent = file.webkitRelativePath || file.name;
    fileList.appendChild(li);
    // uploadFileToServer(file);
    UploadFile(file, post.value.id);

  }
});

</script>

<style>
</style>

