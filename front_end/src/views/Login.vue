<template>
  <v-container style="height: 100%;">
    <v-row justify="center" align="center" align-content-sm="stretch">
      <v-col cols="12" sm="8" md="6" lg="4" style="min-width: 400px;">
        <v-card class="elevation-12">
          <v-toolbar dark color="primary">
            <v-toolbar-title>Login to GGETA</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form>
              <v-text-field
                v-model="email"
                prepend-icon="mdi-email"
                name="email"
                label="Email"
                placeholder="Hello@World.com" type="text"
                required
              >
              </v-text-field>
              <v-text-field
                v-model="password"
                prepend-icon="mdi-lock"
                name="password"
                label="Password"
                :type="showPassword ? 'text' : 'password'"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                required
              />
              <v-btn color="primary" block @click="login_submit">Login</v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import {ref} from "vue";
import router from "@/router";

import {logined, loginV4, logout, showError, showSuccess} from "@/apiv4";
var email = ref('')
var password = ref('')
var showPassword = ref(false)


async function login_submit() {
  loginV4(email.value, password.value).then((response) => {
    // localStorage.setItem('token', response.token)
    // localStorage.setItem('userEmail', response.email)
    // localStorage.setItem('userName', response.name)
    showSuccess("login success")
    router.push('/')
  }).catch(() => {
    logout()
    showError("login failed. please check your email and password.")
  })
}
</script>

<style scoped>
.align-center {
  height: 100%;
}
</style>
