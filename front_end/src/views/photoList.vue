<template>
<!--  <div id="file-list">-->
<!--    <p>Uploaded files:</p>-->
<!--    <ul></ul>-->
<!--  </div>-->
<!--  <v-container>-->
<!--    <v-card v-for="(element, index) in elements" :key="index" style="width: 150px; height: 200px">-->
<!--      <v-img :src="element" />-->
<!--    </v-card>-->
<!--  </v-container>-->
  <v-container fluid>
    <v-row>
      <v-col
        v-for="(element, index) in elements"
          :key="index"
          cols="auto"
          :sm="4"
      :md="3"
      :lg="2"
      :xl="1"
      >
      <v-card class="pa-2" style="height: 300px;">
        <v-img :src="element" />
      </v-card>
      </v-col>
    </v-row>
  </v-container>


</template>
<script setup lang="ts">
import {addPhoto, getPhoto, getPhotoIds, uploadPhotos} from "@/photo_api";
import {ref} from "vue";

const fileUploadArea = document.documentElement;

// console.log("fileList: " + fileList)
fileUploadArea.addEventListener("dragover", (event) => {
  event.preventDefault();
  fileUploadArea.classList.add("dragover");
});
fileUploadArea.addEventListener("dragleave", () => {
  fileUploadArea.classList.remove("dragover");
});
fileUploadArea.addEventListener("drop", async (event) => {
  // const fileList = document.querySelector("#file-list ul");
  event.preventDefault();
  event.stopPropagation(); // Stop event propagation
  fileUploadArea.classList.remove("dragover");
  // @ts-ignore
  const files = event.dataTransfer.files;
  uploadPhotos(files);
  console.log("drop function called");
});
// get photo 6
var elements = ref([] as string[]);


getPhotoIds().then(
    (response) => {
      console.log(response);
      var ids = response.ids;
      // sort ids in descending order
      ids.sort((a, b) => b - a);
      for (var i = 0; i < ids.length; i++) {
        getPhoto({id: ids[i]}).then((response) => {
          console.log(response);
          elements.value.push(response.thum_url);
        });
      }

    }
)



</script>
