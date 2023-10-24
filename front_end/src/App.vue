<template>

  <head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/a11y-light.min.css" integrity="sha512-WDk6RzwygsN9KecRHAfm9HTN87LQjqdygDmkHSJxVkVI7ErCZ8ZWxP6T8RvBujY1n2/E4Ac+bn2ChXnp5rnnHA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css" integrity="sha512-0aPQyyeZrWj9sCA46UlmWgKOP0mUipLQ6OZXu8l4IcAmD2u31EPEy9VcIMvl7SoAaKe8bLXZhYoMaE/in+gcgA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
  </head>
  <v-snackbar v-model="alert.show"
              :timeout="500"
              location="top right"
              position="relative"
              :color="alert.color">
    {{ alert.message }}
  </v-snackbar>

  <router-view/>
</template>

<script lang="ts" setup>
import {alert, verifyToken} from "@/apiv4";

// verify token
verifyToken().then(() => {
  console.log("token verified")
}).catch(() => {
  console.log("token not verified")
})


const script1 = document.createElement("script")
script1.src = "https://cdn.jsdelivr.net/npm/hash-wasm"
script1.type = "text/javascript"
document.head.appendChild(script1)
const script2 = document.createElement("script")
script2.src = "https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.7/MathJax.js?config=TeX-MML-AM_CHTML"
script2.async = false
script2.type = 'text/javascript'
document.head.appendChild(script2)
// const script2 = document.createElement('script')
// script2.src = "//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"
// script2.type = 'text/javascript'
// document.head.appendChild(script2)
const script3 = document.createElement('script')
script3.src = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"
// script3.src = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.3.1/highlight.min.js"
// script3.src = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"
script3.type = 'text/javascript'
document.head.appendChild(script3)
const script4 = document.createElement('script') // clipboard
script4.src = "https://cdnjs.cloudflare.com/ajax/libs/clipboard.js/2.0.8/clipboard.min.js"
script4.type = 'text/javascript'
document.head.appendChild(script4)

// const script5 = document.createElement('script')
// script5.src = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/css.min.js"
// script5.type = 'text/javascript'
// document.head.appendChild(script5)


</script>
<style>
ol, ul {
    margin-left: 20px;
}
/* Apply styles to nested ol and ul */
ol ol, ul ul {
    margin-left: 30px; /* Adjust the left margin for nested lists */
}

/* Apply styles to nested li elements */
li li {
    margin-left: 30px; /* Adjust the left margin for nested list items */
}

/* remove all color and style changes for link */
a {
    color: inherit;
    text-decoration: inherit;
}
/* soft light block color wrapper */
pre > code {
  background-color: #f6f8fa !important;
}
</style>
