<template>
    <v-app>
        <v-navigation-drawer v-model="drawer" :rail="rail" permanent @click="rail = false">
            <v-card class="mt-4 ml-2 mr-1" flat>
                <v-card-actions>
                    <v-avatar image="https://github.com/chen-gz/picBed/blob/master/mine_square.jpg?raw=true"/>
                    <v-spacer/>
                    <v-btn @click="showSearchLine=!showSearchLine" icon="mdi-magnify" flat/>
                </v-card-actions>
                <v-text-field
                    v-if="showSearchLine"
                    variant="solo"
                    label="Search"
                    density="compact"
                />
            </v-card>
            <v-list>
                <v-list-item prepend-icon="mdi-home" to="/" title="View site"/>
                <v-spacer />
                <v-divider/>
                <v-list-subheader title="Post Manager"/>
                <v-list-item prepend-icon="mdi-note-edit" title="Posts" to='/admin/posts/publish' />
                <v-list-item prepend-icon="mdi-note-text" title="Drafts" to='/admin/posts/draft' />
                <v-list-item prepend-icon="mdi-delete-empty" title="Trash" value="Trash" to="/admin/posts/deleted"/>

                <v-spacer />
                <v-divider/>
                <v-list-subheader title="User Manager"/>
                <v-list-item prepend-icon="mdi-account" title="Account" value="account"/>
                <v-list-item prepend-icon="mdi-account-group-outline" title="Users" value="users"/>
            </v-list>
            <v-divider></v-divider>


            <!--            </v-list>-->
        </v-navigation-drawer>
        <v-main>
            <v-app-bar>
                <v-btn icon="mdi-dots-vertical" @click="drawer = !drawer"/>
                <v-spacer/>
                <v-selection-control-group>

                </v-selection-control-group>
                <v-btn prepend-icon="mdi-plus"
                       text="New Post"
                       @click="createNewPost();"
                       variant="elevated"
                       color="primary"
                       class="ml-4"
                />
            </v-app-bar>
            <router-view/>
        </v-main>
    </v-app>
</template>

<script setup lang="ts">
import {ref} from "vue";

import router from "@/router";
import {newPostV4, showError, showSuccess} from "@/apiv4";

let drawer = ref(true)
const rail = ref(false)
let showSearchLine = ref(false)

async function createNewPost() {
    let response = await newPostV4();
    console.log(response)
    if (response.status === "success") {
        showSuccess("create post success")
        router.push('/posts/edit/' + response.url)
    } else {
        showError("create post failed")
    }
}

</script>
