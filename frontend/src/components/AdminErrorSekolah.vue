<template>
  <div class="min-h-screen bg-white flex items-center justify-center p-6 font-sans">
    <div class="text-center max-w-md">
      <!-- Icon Akses Ditolak -->
      <div class="w-20 h-20 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-6">
        <i class="fas fa-lock text-red-500 text-3xl"></i>
      </div>
      
      <div class="w-16 h-0.5 bg-slate-200 mx-auto mb-6"></div>
      
      <h2 class="text-xl font-semibold text-slate-700 mb-2">Akses Ditolak</h2>
      <p class="text-slate-500 text-sm mb-2">
        Anda tidak memiliki izin untuk mengakses halaman ini.
      </p>
      <p class="text-slate-400 text-xs mb-8">
        Halaman ini hanya dapat diakses oleh admin sekolah yang bersangkutan.
      </p>
      
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <button 
          @click="forceBackToDashboard"
          class="px-6 py-2.5 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-md transition-colors"
        >
          ← Kembali ke Dashboard
        </button>
        <button 
          @click="forceLogout"
          class="px-6 py-2.5 bg-white hover:bg-slate-50 text-slate-600 text-sm font-medium rounded-md border border-slate-200 transition-colors"
        >
          Logout
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';
import { computed } from 'vue';

const auth = useAuthStore();
const role = computed(() => auth.role || localStorage.getItem('role'));

const forceBackToDashboard = () => {
  const mySlug = auth.slug || localStorage.getItem('slug');
  if (mySlug && mySlug !== 'null' && mySlug !== '') {
    window.location.replace(`/${mySlug}/dashboard`);
  } else {
    window.location.replace('/login');
  }
}

const forceLogout = () => {
  localStorage.clear();
  window.location.replace('/login');
}
</script>