// Define feature flags before importing Vue
window.__VUE_OPTIONS_API__ = true
window.__VUE_PROD_DEVTOOLS__ = false 
window.__VUE_PROD_HYDRATION_MISMATCH_DETAILS__ = false

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Buat instance aplikasi Vue
const app = createApp(App)

// Gunakan router
app.use(router)

// Mount aplikasi ke DOM
app.mount('#app')
