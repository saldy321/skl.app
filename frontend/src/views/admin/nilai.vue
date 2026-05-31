<template>
  <div class="space-y-6">
    <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-200 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <i class="fas fa-file-import text-emerald-500"></i> Import Nilai Siswa
        </h1>
        <p class="text-slate-500 text-sm">Kelola nilai massal via Excel dengan mudah dan cepat.</p>
      </div>
      
      <div class="flex items-center gap-3 bg-blue-50 border border-blue-100 p-4 rounded-xl">
        <div class="w-10 h-10 bg-blue-500 text-white rounded-full flex items-center justify-center animate-pulse">
          <i class="fas fa-lightbulb"></i>
        </div>
        <div class="text-xs text-blue-800 leading-tight">
          <b>Tips:</b> Cukup download template sekali, isi semua kolom mapel, lalu upload kembali. Sistem akan mendeteksi perubahan data secara otomatis.
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-1">
        <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-200 sticky top-6">
          <div class="flex items-center gap-2 mb-6">
            <span class="w-7 h-7 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center font-bold text-sm">1</span>
            <h2 class="font-bold text-slate-700">Download Template</h2>
          </div>
          
          <div class="space-y-5">
            <div>
              <label class="block text-[10px] font-bold text-slate-400 uppercase mb-2 tracking-widest">Pilih Kelas</label>
              <select v-model="filter.kelas" class="w-full bg-slate-50 border-slate-200 rounded-xl text-sm p-3 focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all">
                <option value="">-- Semua Kelas --</option>
                <option v-for="k in listKelas" :key="k" :value="k">{{ k }}</option>
              </select>
            </div>

            <div v-if="showJurusanField">
              <label class="block text-[10px] font-bold text-slate-400 uppercase mb-2 tracking-widest">Pilih Jurusan</label>
              <select v-model="filter.jurusan" class="w-full bg-slate-50 border-slate-200 rounded-xl text-sm p-3 focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all">
                <option value="">-- Semua Jurusan --</option>
                <option v-for="j in listJurusan" :key="j" :value="j">{{ j }}</option>
              </select>
            </div>

            <button 
              @click="handleDownload"
              :disabled="loadingDownload"
              class="w-full flex items-center justify-center gap-3 px-4 py-4 bg-emerald-500 hover:bg-emerald-600 text-white font-bold rounded-xl transition-all disabled:opacity-50 shadow-lg shadow-emerald-200"
            >
              <i v-if="loadingDownload" class="fas fa-circle-notch animate-spin"></i>
              <i v-else class="fas fa-file-download"></i>
              {{ loadingDownload ? 'Menyiapkan...' : 'Download Template' }}
            </button>
          </div>
        </div>
      </div>

      <div class="lg:col-span-2">
        <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
          <div class="flex items-center gap-2 mb-6">
            <span class="w-7 h-7 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center font-bold text-sm">2</span>
            <h2 class="font-bold text-slate-700">Upload & Sinkronkan</h2>
          </div>
          
          <div 
            class="relative border-2 border-dashed border-slate-200 rounded-2xl p-12 flex flex-col items-center justify-center bg-slate-50 hover:bg-emerald-50 hover:border-emerald-300 transition-all cursor-pointer group"
            @click="$refs.fileInput.click()"
            @dragover.prevent
            @drop.prevent="handleDrop"
          >
            <input type="file" ref="fileInput" class="hidden" @change="onFileChange" accept=".xlsx">
            
            <div class="w-20 h-20 bg-white rounded-2xl flex items-center justify-center shadow-sm mb-5 group-hover:scale-110 group-hover:rotate-3 transition-transform">
              <i class="fas fa-cloud-upload-alt text-3xl text-emerald-500"></i>
            </div>
            
            <h3 class="text-slate-700 font-bold text-lg" v-if="!selectedFile">Pilih file Excel (.xlsx)</h3>
            <h3 class="text-emerald-600 font-bold text-lg" v-else>{{ selectedFile.name }}</h3>
            <p class="text-slate-400 text-sm mt-1">Drag and drop file ke sini atau klik untuk mencari</p>
          </div>

          <div class="mt-8 flex flex-col sm:flex-row items-center justify-between gap-4">
            <div class="flex items-center gap-2 text-xs text-slate-400">
              <i class="fas fa-shield-alt text-emerald-500"></i>
              <span>Data aman & tervalidasi otomatis oleh sistem</span>
            </div>
            
            <button 
              @click="handleUpload"
              :disabled="!selectedFile || loadingUpload"
              class="w-full sm:w-auto px-10 py-4 bg-slate-800 hover:bg-slate-900 text-white font-bold rounded-xl transition-all disabled:opacity-50 flex items-center justify-center gap-3 shadow-xl"
            >
              <i v-if="loadingUpload" class="fas fa-spinner animate-spin"></i>
              <i v-else class="fas fa-rocket"></i>
              {{ loadingUpload ? 'Sedang Sinkronisasi...' : 'Sinkronkan Sekarang' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'
import Swal from 'sweetalert2'

const route = useRoute()
const slug = route.params.slug
const auth = useAuthStore()

const showJurusanField = computed(() => {
  const tingkat = auth.tingkat?.toUpperCase() || ''
  return ['SMK', 'SMA','MA'].includes(tingkat)
})

const loadingDownload = ref(false)
const loadingUpload = ref(false)
const selectedFile = ref(null)
const fileInput = ref(null)

const listKelas = ref([])
const listJurusan = ref([])
const filter = reactive({ kelas: '', jurusan: '' })

const fetchFilters = async () => {
  try {
    const res = await api.get(`/${slug}/admin/nilai/filters`)
    listKelas.value = res.data.kelas
    listJurusan.value = res.data.jurusan
  } catch (e) {
    console.error("Gagal load filter", e)
  }
}

onMounted(() => fetchFilters())

const onFileChange = (e) => {
  const file = e.target.files[0]
  if (file) selectedFile.value = file
}

const handleDrop = (e) => {
  const file = e.dataTransfer.files[0]
  if (file && file.name.endsWith('.xlsx')) {
    selectedFile.value = file
  } else {
    Swal.fire('Format Salah', 'Harap upload file dalam format .xlsx', 'warning')
  }
}

const handleDownload = async () => {
  loadingDownload.value = true
  try {
    const res = await api.get(`/${slug}/admin/nilai/template`, {
      params: filter,
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    const fileName = `Template_Nilai_${filter.kelas || 'Semua'}_${filter.jurusan || 'Semua'}.xlsx`.replace(/ /g, '_')
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    link.remove()
    
    Swal.fire({ icon: 'success', title: 'Siap diisi!', text: 'Silahkan isi nilai pada file Excel tersebut.', timer: 2000, showConfirmButton: false })
  } catch (err) {
    Swal.fire('Waduh!', 'Gagal generate template. Cek koneksi atau data siswa.', 'error')
  } finally {
    loadingDownload.value = false
  }
}

const handleUpload = async () => {
  if (!selectedFile.value) return
  
  loadingUpload.value = true
  const formData = new FormData()
  formData.append('file_excel', selectedFile.value)

  try {
    const res = await api.post(`/${slug}/admin/nilai/import`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    Swal.fire({
      title: 'Misi Sukses!',
      text: res.data.message,
      icon: 'success',
      confirmButtonColor: '#10B981'
    })
    
    selectedFile.value = null
    if (fileInput.value) fileInput.value.value = ''
  } catch (err) {
    Swal.fire('Gagal Sinkron', err.response?.data?.message || 'Terjadi kesalahan sistem.', 'error')
  } finally {
    loadingUpload.value = false
  }
}
</script>

<style scoped>

</style>