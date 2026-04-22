<template>
  <div class="space-y-8 animate-fade-in">
    
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Dashboard Admin</h1>
        <p class="text-sm text-slate-500 mt-1">Ringkasan data dan aktivitas terkini.</p>
      </div>
      <div class="flex items-center gap-3">
        <span class="px-3 py-1 bg-emerald-50 text-emerald-600 text-xs font-semibold rounded-full border border-emerald-100 flex items-center gap-2">
          <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse"></span>
          Data Real-time
        </span>
      </div>
    </div>

    <!-- Stats Grid (HANYA CARD YANG DATANYA ADA) -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      
      <!-- Card 1: Total Siswa (wajib ada) -->
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

      <!-- Card 2: Siswa Lulus (wajib ada) -->
      <router-link :to="`/${slug}/admin/leger`" class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-emerald-200 transition-all duration-300 relative overflow-hidden">
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Siswa Lulus</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-emerald-600 transition-colors">{{ stats.siswaLulus || 0 }}</h3>
          </div>
          <div class="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-emerald-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-graduation-cap"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Status kelulusan</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-emerald-500 transition-colors"></i>
        </div>
      </router-link>

      <!-- Card 3: Belum Input Nilai (hanya tampil jika ada datanya) -->
      <router-link 
        v-if="stats.belumNilai !== undefined && stats.belumNilai > 0"
        :to="`/${slug}/admin/nilai`" 
        class="group bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md hover:border-orange-200 transition-all duration-300 relative overflow-hidden"
      >
        <div class="flex justify-between items-start mb-4">
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-1">Perlu Review</p>
            <h3 class="text-3xl font-bold text-slate-800 group-hover:text-orange-600 transition-colors">{{ stats.belumNilai }}</h3>
          </div>
          <div class="w-12 h-12 bg-orange-50 text-orange-600 rounded-lg flex items-center justify-center text-xl group-hover:bg-orange-600 group-hover:text-white transition-all duration-300">
            <i class="fas fa-exclamation-circle"></i>
          </div>
        </div>
        <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400 font-medium">Menunggu input nilai</span>
          <i class="fas fa-arrow-right text-[10px] text-slate-300 group-hover:text-orange-500 transition-colors"></i>
        </div>
      </router-link>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Left Column: Announcements -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <div class="px-6 py-5 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
            <h3 class="text-sm font-bold text-slate-800 uppercase tracking-wide">Pengumuman Sistem</h3>
            <button class="text-xs font-semibold text-indigo-600 hover:text-indigo-700">Lihat Semua</button>
          </div>
          
          <div class="p-6">
            <div class="flex gap-5 items-start p-5 bg-white rounded-xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow">
              <div class="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center flex-shrink-0">
                <i class="fas fa-bullhorn text-lg"></i>
              </div>
              <div class="flex-1">
                <div class="flex justify-between items-start mb-2">
                  <h4 class="font-bold text-slate-900 text-sm">Update Sistem SKL Digital</h4>
                  <span class="text-[10px] font-medium text-slate-400 bg-slate-100 px-2 py-1 rounded">Hari Ini</span>
                </div>
                <p class="text-sm text-slate-600 leading-relaxed mb-3">
                  Sistem telah diperbarui. Pastikan Anda memeriksa menu 
                  <span class="font-semibold text-indigo-600">"Setting Waktu"</span> untuk mengatur jadwal pengumuman kelulusan siswa.
                </p>
                <a href="#" class="text-xs font-bold text-indigo-600 hover:text-indigo-800 flex items-center gap-1">
                  Baca Selengkapnya <i class="fas fa-chevron-right text-[10px]"></i>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Activity Feed -->
      <div class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden h-full">
          <div class="px-6 py-5 border-b border-slate-100 bg-slate-50/50">
            <h3 class="text-sm font-bold text-slate-800 uppercase tracking-wide">Aktivitas Terbaru</h3>
          </div>
          <div class="p-6">
             <div class="flex flex-col items-center justify-center py-12 text-center">
               <div class="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mb-4 border border-slate-100">
                 <i class="fas fa-history text-slate-300 text-2xl"></i>
               </div>
               <h4 class="text-sm font-bold text-slate-700">Belum Ada Aktivitas</h4>
               <p class="text-xs text-slate-400 mt-1 max-w-[200px]">Log aktivitas admin dan perubahan data akan muncul di sini.</p>
             </div>
          </div>
        </div>
      </div>

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
  siswaLulus: 0,
  belumNilai: 0,
})

const fetchDashboardData = async (targetSlug) => {
  try {
    const response = await api.get(`/${targetSlug}/dashboard`)
    if (response.data) {
      stats.value = {
        totalSiswa: response.data.totalSiswa || 0,
        siswaLulus: response.data.siswaLulus || 0,
        belumNilai: response.data.belumNilai || 0
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
/* Animasi halus saat halaman dimuat */
.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>