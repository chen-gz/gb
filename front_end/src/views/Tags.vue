<script setup lang="ts">
import {ref} from "vue";
import {getDistinct, searchPostsV4, SearchPostsRequestV4, GetDistinctResponse} from "@/apiv4";
// import Lists from "@/layouts/Users/Lists.vue";

let props = defineProps<{
  id?: String
}>();
console.log("id: " + props.id)


const tags = ref({} as GetDistinctResponse)
getDistinct("tags").then((response) => {
  tags.value = response
  // remove empty tag
  tags.value.values = tags.value.values.filter((item) => {
    return item != ""
  })
  if (tags.value.values.length > 0) {
    tags.value.values.sort()
  }
  console.log(tags.value.values)
});

</script>

<template>
  <v-container v-if="props.id == ('' || undefined)">
    <v-chip-group>
      <v-chip v-for="(item, index) in tags.values"
              :key="index"
              :to="'/tag/' + item"
      >
        {{ item }}
      </v-chip>
    </v-chip-group>
  </v-container>
  <!--    <Lists v-else  :searchParam="{tags: props.id} as SearchPostsRequestV4"></Lists>-->
</template>

<style scoped>

</style>
