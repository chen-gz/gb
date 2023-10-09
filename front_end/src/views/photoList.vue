
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
fileUploadArea.addEventListener("drop", async (event) => {
    const fileList = document.querySelector("#file-list ul");

    event.preventDefault();
    fileUploadArea.classList.remove("dragover");

    // @ts-ignore
    const files = event.dataTransfer.files;
    // sort file base one name
    // set string -> file
    var ori_files = new Map()
    var jpeg_files = new Map()

    // for Nikon format
    for (var i = 0; i < files.length; i++) {
        var file = files[i];
        var filename_without_ext = file.name.split(".")[0]
        console.log("filename_without_ext: " + filename_without_ext)
        if (file.name.endsWith(".NEF")) {
            ori_files.set(filename_without_ext, file)
        } else if (file.name.endsWith(".JPG")) {
            jpeg_files.set(filename_without_ext, file)
        }
    }
    // only original file is not accept
    for (var [key, value] of ori_files) {
        if (jpeg_files.has(key) == false) {
            alert("Please upload the original file and the corresponding jpeg file together. The server does not process the original file alone.")
            return
        }
    }

    // upload file to server
    for (var [key, value] of jpeg_files) {
        var ori_file = ori_files.get(key)
        // addPhoto(ori_file, value)
        await addPhoto(value, ori_file)
    }


    // seperate file by fileExtension


});

</script>
