<script setup>

</script>

<template>
  <header class="profile-wrapper">
    <a href="/" id="avatar" class="rounded-circle">
      <img src="https://minio.ggeta.com/blog-public-data/mine_square.jpg"
           width="112" height="112" alt="avatar"
           onerror="this.style.display='none'"></a>
    <h1 class="site-title"><a href="/">GGETA</a></h1>
    <p class="site-subtitle fst-italic mb-0">Knowledge and Live</p>
  </header>

</template>

<style scoped lang="sass">

.profile-wrapper
  margin: 3rem 1.5rem 2.5rem 2.5rem

.rounded-circle
  border-radius: 50%
  overflow: hidden
  display: block
  width: 112px
  height: 112px

  img
    width: 100%
    height: 100%
    object-fit: cover
    object-position: center
    border-radius: 50%
    //border: 1px solid #e9e9e9
    //box-shadow: 0 0 0 1px #e9e9e9
    //transition: box-shadow 0.3s

    &:hover
      box-shadow: 0 0 0 2px #007bff

.site-title
  opacity: 0.7
  width: 100%
  font-size: 1.5rem
  font-weight: 600
  margin-top: 1rem
  margin-bottom: 0.5rem
  font-family: 'Lato', sans-serif

  a
    color: #212529
    text-decoration: none

    &:hover
      text-decoration: underline

.site-subtitle
  font-size: 95%
  color: var(--site-subtitle-color)
  margin-bottom: 0
  margin-top: 0

  a
    color: #6c757d
    text-decoration: none

    &:hover
      text-decoration: underline



</style>