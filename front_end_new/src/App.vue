<script setup>
import HelloWorld from './components/HelloWorld.vue'
import TheWelcome from './components/TheWelcome.vue'
import Sidebar from "@/components/sidebar.vue";
import PostPage from "@/components/PostPage.vue";
</script>

<template>
  <head>
            <link rel="stylesheet"
                  href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">

  </head>
<!--  <header>-->
<!--    <img alt="Vue logo" class="logo" src="./assets/logo.svg" width="125" height="125" />-->

<!--    <div class="wrapper">-->
<!--      <HelloWorld msg="You did it!" />-->
<!--    </div>-->
<!--  </header>-->

<!--  <main>-->
<!--    <TheWelcome />-->
<!--  </main>-->
  <sidebar></sidebar>
  <main>
<!--    <hello-world msg="hello"/>-->
<!--    <h1>Helle world</h1>-->
    <post-page/>
  </main>

</template>

<style lang="sass">
main
  z-index: 99
  padding-left: 3rem


//header {
//  line-height: 1.5;
//}
//.logo {
//  display: block;
//  margin: 0 auto 2rem;
//}
//@media (min-width: 1024px) {
//  header {
//    display: flex;
//    place-items: center;
//    padding-right: calc(var(--section-gap) / 2);
//  }
//
//  .logo {
//    margin: 0 2rem 0 0;
//  }
//
//
//  header .wrapper {
//    display: flex;
//    place-items: flex-start;
//    flex-wrap: wrap;
//  }
//
//}
</style>
