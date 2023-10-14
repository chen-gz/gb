<template>
  <v-container fluid>
    <v-row>
      <template v-for="(element, _) in showElements" :key="index">
        <v-col sm="12" md="3" lg="2">
          <v-card>
            <v-img :src="element.thum_url"
                   style="min-height: 200px; min-width: 200px; position: relative;"
                   @click="openDialog(element)">
            </v-img>
            <v-card-actions style="height: 30px;">
              <v-spacer/>
              <v-btn icon="mdi mdi-download" @click="downloadPhoto(element)"></v-btn>
              <v-btn icon="mdi mdi-delete" @click="deletePhoto(element.photo)"></v-btn>
              <v-btn icon="mdi mdi-share-variant" @click="sharePhoto(element)"></v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </template>

    </v-row>
    <v-dialog v-model="dialog">
      <v-card height=95vh>
        <v-card-title>Dialog Title</v-card-title>
        <v-card variant="flat">
          <v-img :src="dialogImageSrc" width="100%" style="margin-right: auto; margin-left: auto;"></v-img>
        </v-card>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeDialog">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
  <v-pagination v-model="page" :length="number_of_pages" @update:model-value="fetchPage(page)"></v-pagination>
</template>
<script setup lang="ts">

import {getPhoto, getPhotoIds, PhotoItem, UpdatePhoto, uploadPhotos} from "@/photo_api";
import {ref} from "vue";
import {showSuccess} from "@/apiv4";

const fileUploadArea = document.documentElement;
const dialog = ref(false);
const dialogImageSrc = ref("");
var page = ref(1);
var selectedPhotos = ref([] as PhotoWithUrl[]);

function isSelected(photo: PhotoWithUrl) {
  return selectedPhotos.value.includes(photo);
}

function toggleSelect(photo: PhotoWithUrl) {
  if (isSelected(photo)) {
    selectedPhotos.value = selectedPhotos.value.filter((element) => element !== photo);
  } else {
    selectedPhotos.value.push(photo);
  }
}

var selectMode = ref(false);

function toggleSelectMode() {
  selectedPhotos.value = [];
  selectMode.value = !selectMode.value;
}

function closeDialog() {
  dialog.value = false;
}

function openDialog(photo: PhotoWithUrl) {
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
});
// get photo 6
var showElements = ref([] as PhotoWithUrl[]);

interface PhotoWithUrl {
  photo: PhotoItem;
  thum_url: string;
  jpg_url: string;
}

const number_of_pages = ref(0);
const number_of_photos = ref(0);
const page_size = 24;

async function fetchPage(page: number) {
  console.log("fetch page: " + page)
  let photoIds = (await getPhotoIds()).ids;
  // sort by id
  photoIds.sort((a, b) => b - a);
  // get id base on page and page size
  const start = (page - 1) * page_size;
  const end = page * page_size;
  photoIds = photoIds.slice(start, end);

  const promises = photoIds.map(async (id) => {
    return await getPhoto({"id": id});
  });
  const photos = await Promise.all(promises);
  photos.sort((a, b) => b.photo.id - a.photo.id);
  console.log("photos: " + photos[0].photo.id);
  showElements.value = [] as PhotoWithUrl[];
  for (const photo of photos) {
    showElements.value.push({photo: photo.photo, thum_url: photo.thum_url, jpg_url: photo.jpeg_url});
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

function deletePhoto(photo: PhotoItem) {
  photo.deleted = true;
  UpdatePhoto(photo);
  // remove correspond tags
  showElements.value = showElements.value.filter((element) => element.photo.id !== photo.id);
}

function downloadPhoto(photo: PhotoWithUrl) {
  window.open(photo.jpg_url, "_blank");
}

function sharePhoto(photo: PhotoWithUrl) {
  // copy to clipboard
  navigator.clipboard.writeText(photo.jpg_url).then(() => {
    showSuccess("Copied to clipboard")
  });

}
</script>
<style scoped>
.select-icon-container {
  position: absolute;
  top: 10px; /* Adjust the top position as needed */
  right: 10px; /* Adjust the right position as needed */
  z-index: 1; /* Ensure the icon is above the image */
  opacity: 0.1; /* Adjust the opacity value as needed (0.0 for fully transparent, 1.0 for fully opaque) */

}


</style>
