<script setup>

import Topbar from "@/components/main-wrapper/topbar.vue";
import MainContent from "@/components/main-wrapper/content.vue";
import RightPanel from "@/components/main-wrapper/right-panel.vue";
</script>

<template>
    <div id="main-wrapper-topbar-wrapper-outer">
        <topbar/>
    </div>
    <div id="main-wrapper-main" class="d-flex flex-row">
        <div id="main-content-wrapper-outer">
            <main-content/>
        </div>
        <div id="main-wrapper-right-panel-outer">
            <right-panel/>
        </div>
    </div>
</template>

<style scoped lang="sass">
@import "src/assets/variables"
$topbar-height: 56px
$right-panel-width: 300px

#main-wrapper-topbar-wrapper-outer
    height: $topbar-height
    width: 100%
    background-color: #03b303
    z-index: 10

#main-wrapper-main
    top : $topbar-height
    overflow: hidden
    height: calc(100% - #{$topbar-height})
    width: 100%
#main-content-wrapper-outer
    width: calc(100% - #{$right-panel-width})
    height: 100%
    overflow: hidden
    background-color: #f5f5f5
    z-index: 1
    flex-grow: 0
#main-wrapper-right-panel-outer
    width: $right-panel-width
</style>