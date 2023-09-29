/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import App from './App.vue'
// Composables
import {createApp, ref} from 'vue'
import 'font-awesome/css/font-awesome.css'
// Plugins
import {registerPlugins} from '@/plugins'

import {verifyToken, logout, logined} from "@/apiv2";
// import BootstrapVue, {IconsPlugin} from "bootstrap-vue";

const app = createApp(App)
registerPlugins(app)
app.mount('#app')


verifyToken().then(
    response => {
        if (response.status != "success") logout()
        else {
            localStorage.setItem("userName", response.name)
            localStorage.setItem("userEmail", response.email)
            logined.value = true
        }
    })
    .catch(() => {
        logout()
    })


