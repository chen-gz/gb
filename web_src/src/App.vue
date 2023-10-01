<template>
    <v-snackbar v-model="alert.show"
                :timeout="500"
                location="top right"
                position="relative"
                :color="alert.color"
    >
        {{ alert.message }}
    </v-snackbar>
    <router-view/>
</template>

<script lang="ts" setup>
// import {alert} from "@/apiv2";
import {alert} from "@/apiv4";


// add MathJax script to the page
const script = document.createElement('script')
// script.src = "https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"
script.src = "https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.7/MathJax.js?config=TeX-MML-AM_CHTML"
script.async = false
script.type = 'text/javascript'
// script.text = 'MathJax.Hub.Config({tex2jax: {inlineMath: [["$","$"]}});'
document.head.appendChild(script)

const script2 = document.createElement('script')
script2.src = "//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"
script2.type = 'text/javascript'
document.head.appendChild(script2)
const script3 = document.createElement('script')
script3.src = "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
script3.type = 'text/stylesheets'
document.head.appendChild(script3)

const sheet = document.createElement('link')
sheet.rel = "stylesheet"
sheet.href = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/default.min.css"
document.head.appendChild(sheet)




</script>


<style>


</style>
