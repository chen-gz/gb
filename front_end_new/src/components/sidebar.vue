<script setup>

import Sidebar_header from "@/components/sidebar_components/sidebar_header.vue";
import Sidebar_nav from "@/components/sidebar_components/sidebar_nav.vue";
import Sidebar_bottom from "@/components/sidebar_components/sidebar_bottom.vue";
</script>

<template>
  <aside id="sidebar" class="d-flex flex-column">
    <sidebar_header/>
    <sidebar_nav/>
    <sidebar_bottom/>
  </aside>
</template>

<style scoped lang="sass">
// import the variables from the variables.saas file
@import "@/assets/variables.sass"
@import "@/assets/typography-light.scss"
%sidebar-link-hover
  &:hover
    color: var(--sidebar-active-color)
a
  @extend %sidebar-link-hover
#sidebar
  @include light-scheme
  width: $sidebar-width
  height: 100%
  background-color: var(--sidebar-bg)
  border-right: 1px solid var(--sidebar-border-color)
  font-family: $font-family-heading


</style>