<template>
<!--  <HelloWorld />-->

    <v-container>
      <v-row class="mt-10 mb-10 justify-center align-center" fluid>
        <v-card class="justify-center align-content-center" align="center"
                style="width: 100%"
                flat>
          <v-img>
            <v-avatar class="justify-center align-center" size="150">
              <v-img src="https://minio.ggeta.com/blog-public-data/mine_square.jpg"></v-img>
            </v-avatar>
          </v-img>
          <v-card-title class="justify-center align-center">Guangzong Chen</v-card-title>
          <v-card-actions class="justify-center align-center">
            <v-btn icon="mdi-github" href="https://github.com/chen-gz"/>
            <v-btn icon="mdi-email" href="mailto:chen-gz@outlook.com"/>
            <v-btn icon="mdi-key" href="https://keys.openpgp.org/search?q=chen-gz%40outlook.com"/>
            <v-btn icon href="https://t.me/Guangzong">
              <i class="fa fa-telegram fa-lg" aria-hidden="true"></i>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-row>
      <v-row v-for="(item, index) in blog_list" :key="index"
             class="mb-10 justify-center">
        <v-card style="width: 100%; max-width: 1200px" class="align-center">
          <!--            <router-link :to="'/posts/' + item['url']">-->
          <!--              <v-img-->
          <!--                v-show="item.cover_image!== ''"-->
          <!--                height="200px"-->
          <!--                :src="item.cover_image"-->
          <!--                cover-->
          <!--                max-height="200px"-->
          <!--              />-->
          <!--            </router-link>-->
          <v-card-title v-text="item.title || 'No title'"></v-card-title>
          <v-card-text>{{ item.summary }}</v-card-text>
          <v-card-actions>
            <v-btn :to="'/posts/' + item['url']" variant="text" text="Read more"/>
            <v-spacer/>
            <div v-show="item.tags!==''" class="mr-3">
              <v-icon class="mr-3">mdi-tag-multiple</v-icon>
              <!--                <router-link :to="'/tag/' + item.tags">-->
              <!--                  {{ item.tags }}-->
              <!--                </router-link>-->
            </div>
          </v-card-actions>
        </v-card>
      </v-row>
    </v-container>


</template>

<script lang="ts" setup>
  import HelloWorld from '@/components/HelloWorld.vue'
</script>
