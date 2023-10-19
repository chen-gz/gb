<script setup lang="ts">

// import {getPhoto, getPhotoIds, PhotoItem, UpdatePhoto, uploadPhotos} from "@/photo_api";
import {getPhotoById, getPhotoIds, PhotoItemV2, uploadPhotos} from "@/photo_api_v2";
import {ref} from "vue";
import {showSuccess} from "@/apiv4";

const fileUploadArea = document.documentElement;
const dialog = ref(false);
const dialogImageSrc = ref("");
var page = ref(1);
// var selectedPhotos = ref([] as PhotoWithUrl[]);

// function isSelected(photo: PhotoWithUrl) {
//     return selectedPhotos.value.includes(photo);
// }

// function toggleSelect(photo: PhotoWithUrl) {
// if (isSelected(photo)) {
//     selectedPhotos.value = selectedPhotos.value.filter((element) => element !== photo);
// } else {
//     selectedPhotos.value.push(photo);
// }
// }

var selectMode = ref(false);

// function toggleSelectMode() {
//     selectedPhotos.value = [];
//     selectMode.value = !selectMode.value;
// }

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
  photo: PhotoItemV2;
  thum_url: string;
  jpg_url: string;
}

//
const number_of_pages = ref(0);
const number_of_photos = ref(0);
const page_size = 24;

async function fetchPage(page: number) {
  console.log("fetch page: " + page)
  let photoIds = (await getPhotoIds()).ids;
  // let photoIds = []
  // generate a list from 1 to 100
  for (let i = 4528; i <= 4530; i++) {
    photoIds.push(i)
  }
  // sort by id
  photoIds.sort((a, b) => b - a);
  // get id base on page and page size
  const start = (page - 1) * page_size;
  const end = page * page_size;
  photoIds = photoIds.slice(start, end);

  const promises = photoIds.map(async (id) => {
    // return await getPhoto({"id": id});
    return await getPhotoById(id)
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
  // let photoIds = []
  // generate a list from 1 to 100
  // for (let i = 4528; i <= 4530; i++) {
  //     photoIds.push(i)
  // }
  number_of_photos.value = photoIds.length;
  number_of_pages.value = Math.ceil(photoIds.length / page_size);
  console.log("number of pages: " + number_of_pages.value);
  console.log("number of photos: " + number_of_photos.value);
}

fetchData();
fetchPage(1);

// function deletePhoto(photo: PhotoItem) {
//     photo.deleted = true;
//     UpdatePhoto(photo);
//     // remove correspond tags
//     showElements.value = showElements.value.filter((element) => element.photo.id !== photo.id);
// }

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
<template>
  <div></div>
  <v-container fluid>
    <!--    <v-row>-->
    <div class="d-flex flex-wrap">

      <template v-for="(element, index) in showElements" :key="index">
<!--        <v-col sm="12" md="3" lg="2">-->
        <div class="flex-wrap flex-grow-1 d-flex">
<!--          <v-card :key="element.photo.id" >-->
            <v-img :src="element.thum_url"
                   :key="element.photo.id"
                   style="height: 40vh"
                   @click="openDialog(element)"
                   width="auto"
                   class="fill-height flex-grow-1"
            >
            </v-img>
<!--            <v-card-actions>-->
<!--              <v-spacer/>-->
<!--              <v-btn icon="mdi mdi-download" @click="downloadPhoto(element)"></v-btn>-->
<!--&lt;!&ndash;              <v-btn icon="mdi mdi-delete" @click="deletePhoto(element.photo)"></v-btn>&ndash;&gt;-->
<!--              <v-btn icon="mdi mdi-share-variant" @click="sharePhoto(element)"></v-btn>-->
<!--            </v-card-actions>-->
<!--          </v-card>-->
        </div>
<!--        </v-col>-->
      </template>

<!--    </v-row>-->
      </div>
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
