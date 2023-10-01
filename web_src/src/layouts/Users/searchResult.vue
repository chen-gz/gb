<script setup lang="ts">
import {PropType} from "vue";
// import {PostDataV3Meta} from "@/apiv2";

defineProps(
    {
        blog_list: {
            type: Array as PropType<PostDataV3Meta[]>,
            required: true
        }
    }
)

</script>

<template>
    <v-container>
        <v-row>
            <v-col v-for="(item, index) in blog_list" :key="index" cols="12" sm="6" md="4" lg="3">
                <router-link :to="'/posts/' + item['url']">
                    <v-card class="card" height="300px">
                        <v-card-title> {{ item.title || 'No title' }}</v-card-title>
                        <v-card-text v-html="item.summary" class="post_content"/>
                    </v-card>
                </router-link>
            </v-col>
        </v-row>
    </v-container>


</template>

<style scoped>

</style>

<style scoped>

</style>
