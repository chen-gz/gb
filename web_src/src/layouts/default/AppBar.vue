<template>
    <v-app-bar flat floating height="54" >
        <v-app-bar-title
            class="ml-16 d-none d-sm-flex"
            text="GGETA"
            style="cursor: pointer"
            @click="$router.push('/')"
        />
        <v-btn icon="mdi-home" to="/" class="d-flex d-sm-none"/>
        <v-spacer/>
        <v-btn text="Posts" to="/posts" class="d-none d-sm-flex"/>
        <v-btn text="Tags" to="/tag" class="d-none d-sm-flex"/>
        <v-btn text="About" to="/about" class="d-none d-sm-flex"/>

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
        <v-app-bar-nav-icon v-show="!showSearch"
                            icon="mdi-magnify"
                            @click="showSearch=!showSearch;"/>
        <v-menu>
            <template v-slot:activator="{ props: menu }">
                <v-app-bar-nav-icon icon="mdi-account" v-bind="menu"/>
            </template>
            <v-list v-show="logined">
                <v-list-item to="/admin" title="Admin" :active="false"/>
                <v-list-item to="/" @click="logout" :active="false" title="Log out"/>
            </v-list>
            <v-list v-show="!logined">
                <v-list-item to="/login">
                    <v-list-item-title>Log in</v-list-item-title>
                </v-list-item>
            </v-list>
        </v-menu>
        <v-app-bar-nav-icon @click="toggleTheme" icon="mdi-circle-slice-4"/>
    </v-app-bar>
</template>

<script lang="ts" setup>
import {useTheme} from "vuetify";
// import {logined, logout} from "@/apiv2";
import {ref} from "vue";
import {logined, logout} from "@/apiv4";

const theme = useTheme()
let values = ref('')
let showSearch = ref(false)
const searchText = ref(null)

function toggleTheme() {
    theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
    console.log("call toggleTheme")
}
</script>
<style>
@media only screen and (max-width: 600px) {


}

</style>
