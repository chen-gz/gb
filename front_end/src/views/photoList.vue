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
<!--    <v-row>-->
<!--      <template-->
<!--          v-for="(element, index) in elements"-->
<!--          :key="index">-->
<!--        <v-col cols="auto" :sm="4" :md="3" :lg="2" :xl="1">-->
<!--          <v-card>-->
<!--            <v-img :src="element" height="100%" :cover=true />-->
<!--            <v-card-actions class="justify-center">-->
<!--              <v-spacer/>-->

<!--              <v-btn icon="mdi mdi-delete" class="mx-2">-->
<!--              </v-btn>-->
<!--            </v-card-actions>-->
<!--          </v-card>-->
<!--        </v-col>-->
<!--      </template>-->
<!--    </v-row>-->
    <v-row>
      <template v-for="(element, index) in elements" :key="index">
        <v-col cols="auto" :sm="4" :md="3" :lg="2" :xl="1">
          <v-click-outside @click="openDialog(element.id)">
            <v-card class="clickable-card">
              <v-img :src="element.thum_url" height="100%" :cover=true />
              <v-card-actions class="justify-center">
                <v-spacer />
                <v-btn icon="mdi mdi-delete" class="mx-2"></v-btn>
              </v-card-actions>
            </v-card>
          </v-click-outside>
        </v-col>
      </template>
    </v-row>
  </v-container>


</template>
<script setup lang="ts">

import {getPhoto, getPhotoIds, uploadPhotos} from "@/photo_api";
import {ref} from "vue";

const fileUploadArea = document.documentElement;
function openDialog(index: number) {
  console.log(index)
}

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
var elements = ref([] as photo[]);
interface photo {
  id: number;
  thum_url: string;
}


getPhotoIds().then(
    (response) => {
      console.log(response);
      var ids = response.ids;
      // sort ids in descending order
      ids.sort((a, b) => b - a);
      for (var i = 0; i < ids.length; i++) {
        getPhoto({id: ids[i]}).then((response) => {
          console.log(response);
          elements.value.push({id: response.photo.id, thum_url: response.thum_url});
        });
      }

    }
)

</script>
