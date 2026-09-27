import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

import './assets/styles/main.css'

const app = createApp(App)

app.use(createPinia())

// Must run after createPinia() but before router/mount — the guard
// in router/index.ts checks authStore.isAuthenticated on the very
// first navigation, so the session needs to already be restored.
useAuthStore().initFromStorage()

app.use(router)

app.mount('#app')