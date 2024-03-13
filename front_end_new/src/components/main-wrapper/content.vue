<script setup>

</script>

<template>
    <div id="main-content-wrapper-inner">
        <router-view></router-view>
    </div>
</template>

<style scoped lang="sass">
#main-content-wrapper-inner
    padding-left: 3rem
    padding-right: 1rem
    padding-top: 1rem
    overflow-x: scroll
    overflow-y: scroll
    height: 100%
    width: calc(100% - 3rem - 1rem)
    //width: 100%
    scrollbar-width: none
#main-content-wrapper-inner::-webkit-scrollbar
    display: none

</style>