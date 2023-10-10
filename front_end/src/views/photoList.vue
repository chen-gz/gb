<template>
  <v-container fluid>
    <v-row>
      <template v-for="(element, index) in elements" :key="index">
        <v-col cols="auto" :sm="4" :md="3" :lg="2">
          <v-click-outside @click="openDialog(element)">
            <v-card class="clickable-card">
              <v-img :src="element.thum_url" height="100%" :cover=true></v-img>
              <v-card-actions class="justify-center">
                <v-spacer/>
                <v-btn icon="mdi mdi-download" class="mx-2"></v-btn>
                <v-btn icon="mdi mdi-delete" class="mx-2"></v-btn>
              </v-card-actions>
            </v-card>
          </v-click-outside>
        </v-col>
      </template>
    </v-row>
    <v-dialog v-model="dialog">
      <v-card>
        <v-card-title>Dialog Title</v-card-title>
        <v-img :src="dialogImageSrc" height="100%" :cover=true></v-img>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
  <v-pagination v-model="page" :length="number_of_pages" @input="fetchPage(page)"></v-pagination>
</template>
<script setup lang="ts">

import {getPhoto, getPhotoIds, uploadPhotos} from "@/photo_api";
import {ref, watch} from "vue";

const fileUploadArea = document.documentElement;
const dialog = ref(false);
const dialogImageSrc = ref("");
const page = ref(1);

function closeDialog() {
  dialog.value = false;
}
function openDialog(photo: photo) {
  console.log(photo.jpg_url)
  dialogImageSrc.value = photo.jpg_url;
  dialog.value = true;
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
  jpg_url: string;
}

const number_of_pages = ref(0);
const number_of_photos = ref(0);
const page_size = 24;
async function fetchPage(page: number) {
  let photoIds = (await getPhotoIds()).ids;
  // get id base on page and page size
  const start = (page - 1) * page_size;
  const end = page * page_size;
  photoIds = photoIds.slice(start, end);

  const promises = photoIds.map(async (id) => {
    return await getPhoto({"id": id});
  });
  const photos = await Promise.all(promises);
  photos.sort((a, b) => b.photo.id - a.photo.id);
  elements.value = [] as photo[];
  for (const photo of photos) {
    elements.value.push({id: photo.photo.id, thum_url: photo.thum_url, jpg_url: photo.jpeg_url});
  }
}

async function fetchData() {
  const photoIds = (await getPhotoIds()).ids;
  number_of_photos.value = photoIds.length;
  number_of_pages.value = Math.ceil(photoIds.length / page_size);
  console.log("number of pages: " + number_of_pages.value);
  console.log("number of photos: " + number_of_photos.value);
}
fetchData();
fetchPage(1);
</script>
