<template>
  <div class="space-y-8 animate-fade-in">
    
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Dashboard Admin</h1>
        <p class="text-sm text-slate-500 mt-1">Ringkasan data dan aktivitas terkini.</p>
      </div>
    </div>

    <!-- Stats Grid (4 Cards) -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      
      <!-- Card 1: Total Siswa -->
      <router-link :to="`/${slug}/admin/siswa`" class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-indigo-200 transition-all duration-300 relative overflow-hidden">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Total Siswa</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-indigo-600 transition-colors">{{ stats.totalSiswa || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-indigo-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-users"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Terdaftar aktif</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-indigo-500 transition-colors"></i>
        </div>
      </router-link>

      <!-- Card 2: Total Mapel -->
      <router-link :to="`/${slug}/admin/mapel`" class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-teal-200 transition-all duration-300 relative overflow-hidden">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Total Mapel</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-teal-600 transition-colors">{{ stats.totalMapel || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-teal-50 text-teal-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-teal-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-book"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Mata pelajaran</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-teal-500 transition-colors"></i>
        </div>
      </router-link>

      <!-- Card 3: Nilai Terinput -->
      <router-link :to="`/${slug}/admin/nilai`" class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-blue-200 transition-all duration-300 relative overflow-hidden">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Nilai Terinput</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-blue-600 transition-colors">{{ stats.totalNilai || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-blue-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-pen"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Data nilai masuk</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-blue-500 transition-colors"></i>
        </div>
      </router-link>

      <!-- Card 4: Siswa Tidak Lulus -->
      <router-link :to="`/${slug}/admin/siswa`" class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-rose-200 transition-all duration-300 relative overflow-hidden">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Siswa Tidak Lulus</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-rose-600 transition-colors">{{ stats.siswaTidakLulus || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-rose-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-times-circle"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Perlu tindakan</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-rose-500 transition-colors"></i>
        </div>
      </router-link>

    </div>

    <div class="mt-6 text-center text-xs text-slate-400">
      Data diperbarui secara otomatis setiap kali halaman dimuat.
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import api from '@/api'
import router from '@/router';

const route = useRoute();
const slug = ref(route.params.slug);

const stats = ref({
  totalSiswa: 0,
  totalMapel: 0,
  totalNilai: 0,
  siswaTidakLulus: 0,
})

const fetchDashboardData = async (targetSlug) => {
  try {
    const response = await api.get(`/${targetSlug}/dashboard`)
    if (response.data) {
      stats.value = {
        totalSiswa: response.data.totalSiswa || 0,
        totalMapel: response.data.totalMapel || 0,
        totalNilai: response.data.totalNilai || 0,
        siswaTidakLulus: response.data.siswaTidakLulus || 0,
      }
    }
    slug.value = targetSlug
  } catch (err) {
    handleError(err)
  }
}

const handleError = (err) => {
  if (!err.response) {
    alert("Gagal terhubung ke server. Pastikan backend aktif.");
    return;
  }

  const status = err.response.status;

  switch (status) {
    case 401:
      localStorage.clear();
      router.push('/login'); 
      break;
    case 403:
      router.push({ name: 'AdminErrorPage' });
      break;
    case 404:
      router.push({ name: 'AdminErrorPage' });
      break;
    case 500:
      alert("Server lagi pusing (Internal Server Error).");
      break;
    default:
      console.error("Terjadi kesalahan sistem.");
  }
}

onMounted(() => {
  fetchDashboardData(route.params.slug);
});

onBeforeRouteUpdate(async (to, from) => {
  if (to.params.slug !== from.params.slug) {
    await fetchDashboardData(to.params.slug);
  }
});
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>