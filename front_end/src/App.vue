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
import {alert} from "@/apiv4";

document.title = "EtaBlog";
</script>
