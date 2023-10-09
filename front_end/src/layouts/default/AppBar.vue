<template>
  <v-app-bar flat>
    <v-app-bar-title
      class="d-none d-sm-flex text-decoration-none"
    >
<!--      remove decoration for this link -->
      <router-link to="/" class="text-decoration-none">
        <v-icon icon="mdi-circle-slice-4"/>
        GGETA
      </router-link>
    </v-app-bar-title>
    <!--    <v-btn icon="mdi-home" to="/" class="d-sm-none"/>-->
    <v-btn class="d-sm-none"><i class="fa fa-home fa-lg" aria-hidden="true"></i></v-btn>
    <v-spacer class="d-sm-none"/>


    <v-text-field
      flat
      ref="searchText"
      v-show="showSearch"
      v-model="values"
      prepend-inner-icon="mdi-magnify"
      placeholder="Search"
      single-line
      density="compact"
      hide-details
      variant="solo-filled"
      class="mr-4"
      @keydown.enter="$emit('search', values)"
      @focusout="showSearch=!showSearch"
      autofocus
    ></v-text-field>
    <v-app-bar-nav-icon class="" icon="mdi-magnify" @click="showSearch=!showSearch;"/>

    <v-btn text="Photos" to="/photos" class="d-none d-sm-flex"/>
    <v-btn text="Posts" to="/posts" class="d-none d-sm-flex"/>
    <v-btn text="Tags" to="/tags" class="d-none d-sm-flex"/>
    <v-btn text="About" to="/about" class="d-none d-sm-flex"/>

    <v-app-bar-nav-icon class="d-sm-none"><i class="fa fa-tags fa-lg" aria-hidden="true"></i></v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-sm-none"><i class="fa fa-sticky-note fa-lg" aria-hidden="true"></i>
    </v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-sm-none"><i class="fa fa-info fa-lg" aria-hidden="true"></i></v-app-bar-nav-icon>
    <v-app-bar-nav-icon @click="login()" v-if="!is_logined"><i class="fa fa-user" aria-hidden="true"></i></v-app-bar-nav-icon>
    <v-app-bar-nav-icon @click="newpost()" v-if="is_logined"> <i class="fa fa-plus" aria-hidden="true"></i></v-app-bar-nav-icon>
  </v-app-bar>
</template>

<script lang="ts" setup>
import {createCommentVNode, ref, watch} from "vue";
import router from "@/router";
import {logined, newPostV4, is_logined} from "@/apiv4";

var showSearch = ref(false)
var values = ref('')
// var is_logined = ref(false)
// is_logined.value = logined()

function login() {
  // verify token
  if (!logined()) {
    router.push('/login')
  }
}
function newpost() {
  newPostV4().then(response => {
    router.push('/posts/edit/' + response.url)
  })
}
logined()
</script>
