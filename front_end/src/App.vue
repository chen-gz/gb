<template>

  <head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
  </head>
  <v-snackbar v-model="alert.show"
              :timeout="500"
              location="top right"
              position="relative"
              :color="alert.color">
    {{ alert.message }}
  </v-snackbar>

  <router-view/>
</template>

<script lang="ts" setup>
import {alert, verifyToken} from "@/apiv4";

document.title = "EtaBlog";
// verify token
verifyToken().then(() => {
  console.log("token verified")
}).catch(() => {
  console.log("token not verified")
})



const script1 = document.createElement("script")
script1.src = "https://cdn.jsdelivr.net/npm/hash-wasm"
script1.type = "text/javascript"
document.head.appendChild(script1)

</script>

