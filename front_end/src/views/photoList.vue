
<template>
  <div id="file-list">
    <p>Uploaded files:</p>
    <ul></ul>
  </div>

</template>

<style scoped>

</style>

<script setup lang="ts">
import {addPhoto} from "@/photo_api";


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
  // sort file base one name
  file.sort((a, b) => {
    if (a.name > b.name) {
      return 1;
    } else {
      return -1;
    }
  })

  // for (const file of Array.from(files)) {
    // const li = document.createElement("li");
    // li.textContent = file.webkitRelativePath || file.name;
    // // @ts-ignore
    // fileList.appendChild(li);
    // // uploadFileToServer(file);
    // UploadFile(file, post.value.id);
  //
  }
});

</script>
