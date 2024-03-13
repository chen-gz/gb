<script setup lang="ts">
import {onMounted, ref, watch} from "vue";
import {getDistinct, SearchPostsRequestV4, GetDistinctResponse} from "../../../../apiv4"
// import Lists from "@/views/Lists.vue";
let props = defineProps<{
  tag_name: String
}>();
let tags = ref({} as GetDistinctResponse)
let searchParam = {} as SearchPostsRequestV4
watch(() => props.tag_name, (old, newe) => {
  console.log("props.tag_name changed")
  searchParam.tags = props.tag_name as string
})
onMounted(() => {
  getDistinct("tags").then((response) => {
    tags.value = response
    // remove empty tag
    tags.value.values = tags.value.values.filter((item) => {
      return item != ""
    })
    if (tags.value.values.length > 0) {
      tags.value.values.sort()
    }
  });
})

</script>


<template>
    <h1>Tags </h1>
    <div>
        <div v-for="(item, index) in tags.values" :key="index">
            <a :to="'/tags/' + item">
                {{item}}
                <span>
                    1
                </span>
            </a>
        </div>
    </div>
</template>

