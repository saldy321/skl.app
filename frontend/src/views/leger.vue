<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">Leger & Rekap Nilai</h1>
        <p class="text-sm text-slate-500">Rekapitulasi otomatis berdasarkan data filter terbaru.</p>
      </div>
      
    </div>

    <div class="bg-white p-5 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end">
      <div class="flex-1 min-w-[200px]">
        <label class="block text-xs font-bold text-slate-500 uppercase mb-2">Pilih Kelas</label>
        <select v-model="filter.kelas" class="w-full bg-slate-50 border-slate-200 rounded-xl text-sm p-3 outline-none focus:ring-2 focus:ring-indigo-500/20">
          <option value="">-- Semua Kelas --</option>
          <option v-for="k in options.kelas" :key="k" :value="k">{{ k }}</option>
        </select>
      </div>
      <div class="flex-1 min-w-[200px]">
        <label class="block text-xs font-bold text-slate-500 uppercase mb-2">Pilih Jurusan</label>
        <select v-model="filter.jurusan" class="w-full bg-slate-50 border-slate-200 rounded-xl text-sm p-3 outline-none focus:ring-2 focus:ring-indigo-500/20">
          <option value="">-- Semua Jurusan --</option>
          <option v-for="j in options.jurusan" :key="j" :value="j">{{ j }}</option>
        </select>
      </div>
      <button @click="fetchLeger" :disabled="loading" class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-bold text-sm hover:bg-indigo-700 disabled:opacity-50 transition-all">
        <i class="fas fa-sync-alt mr-2" :class="{'animate-spin': loading}"></i>
        {{ loading ? 'Memuat...' : 'Tampilkan Data' }}
      </button>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      <div v-if="legerData.data && legerData.data.length > 0" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50">
              <th class="p-4 text-xs font-bold text-slate-400 uppercase border-b sticky left-0 bg-slate-50 z-20 min-w-[220px]">Identitas Siswa</th>
              <th v-for="m in legerData.mapels" :key="m.id" class="p-4 text-xs font-bold text-slate-400 uppercase border-b text-center min-w-[100px]">
                {{ m.nama_mapel }}
              </th>
              <th class="p-4 text-xs font-bold text-indigo-600 uppercase border-b text-center min-w-[100px] bg-indigo-50/50">
                Rata-Rata
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in legerData.data" :key="idx" class="hover:bg-slate-50/80 transition-all">
              <td class="p-4 border-b font-medium sticky left-0 bg-white z-10 shadow-[4px_0_8px_-4px_rgba(0,0,0,0.1)]">
                <div class="font-bold text-slate-800 uppercase">{{ row.nama }}</div>
                <div class="text-[10px] text-slate-400 font-mono">{{ row.nisn }}</div>
              </td>
              <td v-for="m in legerData.mapels" :key="m.id" class="p-4 border-b text-center font-mono text-sm">
                <span :class="(row[m.nama_mapel] !== '-' && parseFloat(row[m.nama_mapel]) < 75) ? 'text-rose-500 font-bold' : 'text-slate-600'">
                  {{ row[m.nama_mapel] }}
                </span>
              </td>
              <td class="p-4 border-b text-center font-bold bg-indigo-50/30">
                <span :class="parseFloat(row.rata_rata) < 75 ? 'text-rose-600' : 'text-indigo-600'">
                  {{ row.rata_rata ? parseFloat(row.rata_rata).toFixed(2) : '0.00' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else class="p-20 text-center">
        <div class="bg-slate-50 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4 text-slate-300 text-2xl">
          <i class="fas fa-database"></i>
        </div>
        <h3 class="text-slate-600 font-bold">Data Tidak Ditemukan</h3>
        <p class="text-slate-400 text-sm">Pastikan filter sudah benar atau data nilai sudah di-import.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import Swal from 'sweetalert2'

const route = useRoute()
const auth = useAuthStore()
const slug = route.params.slug

const loading = ref(false)
const filter = reactive({ kelas: '', jurusan: '' })

const options = reactive({ 
  kelas: [], 
  jurusan: [] 
})

const legerData = ref({
  mapels: [],
  data: []
})

const fetchFilters = async () => {
  try {
    const res = await api.get(`/${slug}/admin/nilai/filter-options`)
    options.kelas = res.data.kelas
    options.jurusan = res.data.jurusan
  } catch (e) {
    console.error("Gagal load filter", e)
  }
}

const fetchLeger = async () => {
  if (!filter.kelas || !filter.jurusan) {
    return Swal.fire('Info', 'Pilih Kelas & Jurusan dulu!', 'warning')
  }

  loading.value = true
  try {
    const res = await api.get(`/${slug}/admin/nilai/leger`, {
      params: {
        kelas: filter.kelas,
        jurusan: filter.jurusan,
        tahun_lulus: auth.selectedYear
      }
    })
    
    // Opsional: Urutin dari rata-rata tertinggi (ranking)
    if (res.data.data) {
        res.data.data.sort((a, b) => b.rata_rata - a.rata_rata)
    }
    
    legerData.value = res.data
  } catch (error) {
    Swal.fire('Error', 'Gagal ambil data rekap', 'error')
  } finally {
    loading.value = false
  }
}

watch(() => auth.selectedYear, () => {
    // Reset filter kelas/jurusan jika perlu, atau langsung fetch ulang
    fetchLeger();
});

onMounted(() => {
  fetchFilters()
})

const exportExcel = () => {
  Swal.fire('Info', 'Lagi diproses bro!', 'info')
}
</script>