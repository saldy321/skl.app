import './style.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import api from './api'
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'
import '@fortawesome/fontawesome-free/css/all.min.css'


const pinia = createPinia()
const app = createApp(App)


app.use(pinia)



async function validateSessionBeforeRender() {
  const authStore = useAuthStore()
  
  // Jika role siswa, langsung return (tidak perlu cek /me)
  if (authStore.role === 'siswa') {
    console.log('[MAIN] Role siswa, skip /me')
    return
  }
  
  if (localStorage.getItem('isLoggedIn') === 'true') {
    try {
      const res = await api.get('/me')
      console.log('[MAIN] Response /me:', res.data)
      
      if (res.data.status === 'success') {
        authStore.setAuthData(res.data.data)
      } else {
        await authStore.logout()
      }
    } catch (error) {
      console.log('[MAIN] Error fetch /me:', error)
      await authStore.logout()
    }
  } else {
    authStore.$reset()
  }
}

await validateSessionBeforeRender()

  
app.component('QuillEditor', QuillEditor)
app.use(router)
app.mount('#app')





