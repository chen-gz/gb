<script lang="ts" setup>
import {ref} from 'vue'
import {SearchPostsRequestV4, searchPostsV4, V4PostData} from "/apiv4.js"
import {formatDate} from "/ui_utils";

let article = ref([] as V4PostData[])
// get 5 latest articles
// searchPostsV4

let param = {} as SearchPostsRequestV4;
// param.limit = {start: 0, count: 5};
// param.constructor();
// searchPostsV4()
param.limit = {start: 0, size: 10}
param.sort = "created_at DESC"
searchPostsV4(param).then((res) => {
    // console.log(res)
    article.value = res.posts
    for (let i = 0; i < article.value.length; i++) {
        if (article.value[i].summary.length == 0) {
            // get summary from content
            article.value[i].summary = article.value[i].content.substring(0, 300)
        }
    }
})

// sort: "created_at DESC"

</script>

<template>
    <div id="home-inner">
        <article v-for="post in article" :key="post.id" class="home-card" @click="$router.push('/post/' + post.url)">
            <h2 class="title">
                {{ post.title }}
                <!--                <a :href="'/post/' + post.url"> {{ post.title }} </a>-->
            </h2>
            <p class="home-summary"> {{ post.summary }}</p>
            <div class="home-card-bottom">
                <span>
                    <i class="fa-fw fas fa-calendar-alt"></i>
                    {{ formatDate(post.created_at) }}
                </span>
                <span v-if="post.category != ''"> <i class="fa-fw fas fa-folder"></i> {{ post.category }} </span>
                <span v-if="post.tags != ''"> <i class="fa-fw fas fa-tags"></i> {{ post.tags }}</span></div>
        </article>
    </div>
</template>

<style lang="sass" scoped>
.title
    margin-bottom: 10px
    font-weight: 500
    font-size: 1.2rem

article
    cursor: pointer
    display: flex
    //font-weight: 1rem
    line-height: 1.5rem
    flex-direction: column

#home-inner
    display: flex
    flex-direction: column
    width: 100%

.home-card
    height: 150px
    border: 1px solid #f6f6f6
    border-radius: 8px
    padding: 10px 10px 3px 16px
    width: 100%
    margin-bottom: 20px
    margin-top: 3px
    box-shadow: 1px 1px 5px lightgrey

    &:hover
        background-color: #f6f6f6

    .home-summary
        color: #666
        flex-grow: 1
        // use ... to hide the overflow
        text-overflow: ellipsis
        //white-space: nowrap
        overflow: hidden
        text-overflow: ellipsi
        display: -webkit-box
        -webkit-box-orient: vertical
        -webkit-line-clamp: 3
       
      



//    put at the bottom of the card
.home-card-bottom
    position: relative
    display: flex
    //justify-content: space-between
    height: 25px
    font-size: 0.8em
    color: #666
    align-items: center
    bottom: 1px

    opacity: 0.77

    span
        margin-right: 25px
        &:last-child
            margin-right: 0
        i 
            margin-right: 5px

</style>
