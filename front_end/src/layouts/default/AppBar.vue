<template>
  <v-app-bar flat>
    <v-app-bar-title class="d-none d-sm-flex text-decoration-none" >
      <router-link to="/" class="text-decoration-none">
        GGETA
      </router-link>
    </v-app-bar-title>
    <v-btn class="d-sm-none"><i class="fa fa-home fa-lg" aria-hidden="true"></i></v-btn>
    <v-spacer class="d-sm-none"/>
    <v-text-field
        flat
        ref="searchText"
        v-show="showSearchTextField"
        v-model="values"
        prepend-inner-icon="mdi-magnify"
        placeholder="Search"
        single-line
        density="compact"
        hide-details
        variant="solo-filled"
        class="mr-4"
        @input="debouncedEmit(values)"
        @focusout="showSearchTextField=!showSearchTextField"
        autofocus
    ></v-text-field>
    <v-app-bar-nav-icon class="" icon="mdi-magnify" @click="showSearchTextField=!showSearchTextField;"/>

    <v-btn text="Photos" to="/photos" class="d-none d-sm-flex"/>
    <v-btn text="Posts" to="/posts" class="d-none d-sm-flex"/>
    <v-btn text="Tags" to="/tags" class="d-none d-sm-flex"/>
    <v-btn text="About" to="/about" class="d-none d-sm-flex"/>

    <v-app-bar-nav-icon class="d-sm-none" to="/photos"><i class="fa fa-camera fa-lg"/></v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-sm-none" to="/posts"><i class="fa fa-sticky-note fa-lg"/></v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-sm-none" to="/tags"><i class="fa fa-tags fa-lg"/></v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-sm-none" to="/about"><i class="fa fa-info fa-lg"/></v-app-bar-nav-icon>
    <v-app-bar-nav-icon @click="login()" v-if="!is_logined" :key="is_logined.toString()"><i class="fa fa-user"/></v-app-bar-nav-icon>
    <v-app-bar-nav-icon class="d-none d-sm-flex" @click="newpost()" v-if="is_logined" :key="is_logined.toString()" ><i class="fa fa-plus"/>
    </v-app-bar-nav-icon>
  </v-app-bar>
</template>

<script lang="ts" setup>
import {createCommentVNode, ref, watch} from "vue";
import router from "@/router";
import {logined, newPostV4, is_logined} from "@/apiv4";

var showSearchTextField = ref(false)
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
// function debounce(func, delay) {
//   let timeoutId = null;
//   return (...args) => {
//     clearTimeout(timeoutId);
//     timeoutId = setTimeout(() => func.apply(this, args), delay);
//   };
// }
// const debouncedEmit = debounce((value) => {
//   $emit('search', value);
// }, 500); // Adjust the delay as needed

const emit = defineEmits([ 'search'])
function debounce<T extends (...args: any[]) => any>(func: T, delay: number): (...funcArgs: Parameters<T>) => void {
      let timeoutId: number | null = null;
      return (...args: Parameters<T>) => {
        if (timeoutId !== null) {
          clearTimeout(timeoutId);
        }
        timeoutId = window.setTimeout(() => {
          func(...args);
        }, delay);
      };
    }

    // Debounced version of the emit function
    const debouncedEmit = debounce((value: string) => {

      emit('search', value);
    }, 500); // Adjust the delay as needed


</script>
