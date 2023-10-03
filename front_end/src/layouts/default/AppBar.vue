<template>
  <v-app-bar flat>
    <v-app-bar-title
      class="d-none d-sm-flex"
    >
      <v-icon icon="mdi-circle-slice-4" />
      GGETA
    </v-app-bar-title>
<!--    <v-btn icon="mdi-home" to="/" class="d-sm-none"/>-->
    <v-btn class="d-sm-none"><i class="fa fa-home fa-lg" aria-hidden="true"></i></v-btn>


    <v-btn text="Posts" to="/posts" class="d-none d-sm-flex"/>
    <v-btn text="Tags" to="/tag" class="d-none d-sm-flex"/>
    <v-btn text="About" to="/about" class="d-none d-sm-flex"/>

    <v-btn class="d-sm-none" > <i class="fa fa-tags fa-lg" aria-hidden="true"></i> </v-btn>
    <v-btn class="d-sm-none"> <i class="fa fa-sticky-note fa-lg" aria-hidden="true"></i> </v-btn>
    <v-btn class="d-sm-none"> <i class="fa fa-info fa-lg" aria-hidden="true"></i> </v-btn>


  </v-app-bar>
</template>

<script lang="ts" setup>
  //
</script>
