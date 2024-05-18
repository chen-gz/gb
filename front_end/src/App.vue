<script setup>
import Sidebar from "@/components/sidebar.vue";
import PostPage from "@/components/main-wrapper/content/post.vue";
import {watch} from "vue";
// nextTick call tocbot
import {nextTick} from 'vue'
import Category from "@/components/main-wrapper/content/category.vue";
import MainWrapper from "@/components/main-wrapper.vue";
import Sidebar_bottom from "@/components/sidebar/sidebar_bottom.vue";

// add copy button to code block
watch(() => {
	nextTick(() => {
  		document.querySelectorAll('pre').forEach((block) => {
		  if (!block.querySelector('.copy-button')) {
			const copyButton = document.createElement('button');
			copyButton.className = 'copy-button';
			copyButton.textContent = 'Copy';
			copyButton.addEventListener('click', () => {
			  navigator.clipboard.writeText(block.textContent).then(() => {
				copyButton.textContent = 'Copied!';
				setTimeout(() => {
				  copyButton.textContent = 'Copy';
				}, 2000);
			  });
			});
			block.style.position = 'relative';
			block.appendChild(copyButton);
		  }
		});

	});
});



</script>

<template>
    <aside id="app-sidebar">
        <sidebar></sidebar>
    </aside>
    <main>
        <main-wrapper></main-wrapper>
    </main>
</template>

<style lang="sass">
@import "assets/variables"
main
    z-index: 2
    flex-grow: 0
    flex-shrink: 0
    overflow-y: auto
    overflow-x: hidden
    word-wrap: anywhere
    width: calc(100% - $sidebar-width)
    height: 100%

#app-sidebar
    width: $sidebar-width
    height: 100%
    flex-grow: 1
    flex-shrink: 0
    overflow: hidden
@media (max-width: 768px)
    main
        width: 100%
    #app-sidebar
        width: 0
        display: none
    //.sidebar-bottom
    //    display: none
    //.sidebar-top
    //    display: none
    //.sidebar-top
</style>
