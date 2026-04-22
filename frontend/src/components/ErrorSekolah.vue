<template>
  <div class="min-h-screen bg-white flex items-center justify-center p-6 font-sans">
    <div class="text-center max-w-md">
      <div class="w-20 h-20 bg-amber-50 rounded-full flex items-center justify-center mx-auto mb-6">
        <i class="fas fa-school text-amber-500 text-3xl"></i>
      </div>
      
      <div class="w-16 h-0.5 bg-slate-200 mx-auto mb-6"></div>
      
      <h2 class="text-2xl font-bold text-slate-800 mb-2">404</h2>
      <h3 class="text-xl font-semibold text-slate-700 mb-2">Sekolah Tidak Ditemukan</h3>
      <p class="text-slate-500 text-sm mb-2">
        Sekolah dengan kode <span class="font-mono text-amber-600 bg-amber-50 px-1 rounded">/{{ slug }}</span> tidak terdaftar.
      </p>
      <p class="text-slate-400 text-xs mb-8">
        Pastikan Anda menggunakan link yang benar dari sekolah Anda.
      </p>
      
      
    </div>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'

const route = useRoute()
const router = useRouter()

// ✅ Ambil slug dari params
const slug = route.params.slug || ''

const goHome = () => {
  // ✅ Gunakan replace untuk menghapus history
  router.replace('/login')
}
</script>