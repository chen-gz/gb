<script setup lang="ts">
import axios from "axios";
import {blogBackendUrl} from "@/config";
import {ref} from "vue";
import {useRouter} from "vue-router";
import {getDistinct, GetDistinctResponse, SearchPostsRequestV3} from "@/apiv2";
import Lists from "@/layouts/Users/Lists.vue";

let props = defineProps<{
    id?: String
}>();
console.log("id: " + props.id)


const tags = ref({}as GetDistinctResponse)
getDistinct("tags").then((response) => {
    tags.value = response
    console.log(tags.value.values)
})

</script>

<template>
    <v-container v-if="props.id == ('' || undefined)">
        <v-chip-group>
            <v-chip v-for="(item, index) in tags.values"
                    :key="index"
                    :to="'/tag/' + item"
            >
                {{item}}
            </v-chip>
        </v-chip-group>
    </v-container>
    <Lists v-else  :searchParam="{tags: props.id} as SearchPostsRequestV3"></Lists>
</template>

<style scoped>

</style>
