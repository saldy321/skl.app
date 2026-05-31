<template>
  <div class="min-h-screen bg-[#F3F4F6] font-sans text-slate-800">
    
    <main class="max-w-7xl mx-auto p-6 space-y-6">
      
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="(stat, idx) in stats" :key="idx" 
          class="bg-white p-5 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow duration-200">
          <div class="flex justify-between items-start mb-4">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide">{{ stat.label }}</p>
              <h3 class="text-2xl font-bold text-gray-900 mt-1">{{ stat.value }}</h3>
            </div>
            <div :class="`w-10 h-10 ${stat.bgColor} ${stat.color} rounded-lg flex items-center justify-center text-lg`">
              <i :class="stat.icon"></i>
            </div>
          </div>
          <div class="pt-4 border-t border-gray-100">
            <p class="text-[10px] text-gray-400 font-medium">Last updated: Just now</p>
          </div>
        </div>
      </div>

      <!-- Table Section -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
        <div class="px-6 py-5 border-b border-gray-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-gray-50/50">
          <div>
            <h2 class="text-base font-bold text-gray-900">Daftar Instansi Sekolah</h2>
            <p class="text-xs text-gray-500 mt-0.5">
              Menampilkan <span class="font-semibold text-blue-600">{{ recentInstansi.length }}</span> data terdaftar.
            </p>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-200">
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider">Informasi Sekolah</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider text-center">Kode Instansi</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider text-center">Jenjang</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider">Slug / URL</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider text-right">Terdaftar</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="item in recentInstansi" :key="item.id" class="hover:bg-gray-50 transition-colors duration-150">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <!-- LOGO SEKOLAH -->
                    <div class="w-10 h-10 rounded-full flex items-center justify-center overflow-hidden bg-white border border-gray-200">
                      <img 
                        v-if="item.logo_instansi"
                        :src="`http://localhost:3000/uploads/instansi/${item.logo_instansi}`"
                        class="w-full h-full object-cover"
                        alt="Logo"
                      />
                      <div v-else class="w-full h-full bg-blue-50 text-blue-600 flex items-center justify-center text-xs font-bold">
                        {{ item.nama_instansi ? item.nama_instansi.substring(0, 2).toUpperCase() : 'SC' }}
                      </div>
                    </div>
                    
                    <!-- INFO TEKS -->
                    <div>
                      <h4 class="text-sm font-semibold text-gray-900">{{ item.nama_instansi }}</h4>
                      <p class="text-[10px] text-gray-400 font-mono mt-0.5">ID: {{ item.id.toString().substring(0,8) }}...</p>
                    </div>
                  </div>
                </td>

                <td class="px-6 py-4 text-center">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full bg-gray-100 text-gray-700 text-[10px] font-semibold border border-gray-200">
                    {{ item.kode_instansi || '-' }}
                  </span>
                </td>

                <td class="px-6 py-4 text-center">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full bg-gray-100 text-gray-700 text-[10px] font-semibold border border-gray-200">
                    {{ item.tingkat_sekolah || 'Umum' }}
                  </span>
                </td>

                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <i class="fas fa-link text-gray-300 text-[10px]"></i>
                    <span class="text-xs text-gray-600 font-medium">/{{ item.slug }}</span>
                  </div>
                </td>

                <td class="px-6 py-4 text-right">
                  <span class="text-xs text-gray-500">{{ formatTimestamp(item.created_at) }}</span>
                </td>
              </tr>

              <tr v-if="recentInstansi.length === 0">
                <td colspan="5" class="px-6 py-12 text-center">
                  <div class="inline-flex flex-col items-center justify-center">
                    <div class="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center mb-3">
                      <i class="fas fa-school text-gray-400 text-lg"></i>
                    </div>
                    <h3 class="text-sm font-semibold text-gray-900">Belum Ada Data</h3>
                    <p class="text-xs text-gray-500 mt-1">Silakan tambahkan instansi sekolah baru.</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="px-6 py-3 border-t border-gray-200 bg-gray-50 flex justify-between items-center">
          <p class="text-[10px] text-gray-500">Showing all records</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

// Data
const recentInstansi = ref([])
const stats = ref([
  { label: 'Total Instansi', value: '0', icon: 'fas fa-school', color: 'text-blue-600', bgColor: 'bg-blue-50' },
  { label: 'Admin Aktif', value: '0', icon: 'fas fa-user-shield', color: 'text-emerald-600', bgColor: 'bg-emerald-50' },
  { label: 'Siswa Nasional', value: '0', icon: 'fas fa-user-graduate', color: 'text-orange-600', bgColor: 'bg-orange-50' }
])

// Helper format timestamp
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleDateString('id-ID', { 
    day: 'numeric', 
    month: 'short', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// ========== FETCH DATA ==========
const fetchData = async () => {
  try {
    const resStats = await api.get('/super/dashboard')
    if(resStats.data.status === 'success') {
      const d = resStats.data.data
      stats.value[0].value = d.totalSekolah.toString()
      stats.value[1].value = d.totalAdmin.toString()
      stats.value[2].value = d.totalSiswaNasional.toString()
    }

    const resSekolah = await api.get('/super/instansi')
    if(resSekolah.data.status === 'success') {
      recentInstansi.value = resSekolah.data.data
    }
  } catch (err) {
    console.error("Dashboard error:", err)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* Custom Scrollbar for Clean Look */
.overflow-x-auto::-webkit-scrollbar {
  height: 6px;
  width: 6px;
}
.overflow-x-auto::-webkit-scrollbar-track {
  background: transparent;
}
.overflow-x-auto::-webkit-scrollbar-thumb {
  background: #CBD5E1;
  border-radius: 10px;
}
.overflow-x-auto::-webkit-scrollbar-thumb:hover {
  background: #94A3B8;
}
</style>