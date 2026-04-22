<template>
  <div class="max-w-5xl mx-auto">
    <!-- Header Section -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-900 tracking-tight">Tambah Instansi Baru</h2>
      <p class="text-sm text-gray-500 mt-1">Lengkapi data sekolah dan akun administrator utama untuk memulai.</p>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-8">
      
      <!-- SECTION 1: DATA SEKOLAH -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex items-center gap-3">
          <div class="w-8 h-8 bg-blue-100 text-blue-600 rounded-lg flex items-center justify-center">
            <i class="fas fa-school"></i>
          </div>
          <div>
            <h3 class="text-sm font-bold text-gray-900">Profil Instansi</h3>
            <p class="text-[10px] text-gray-500 uppercase tracking-wide">Informasi dasar sekolah</p>
          </div>
        </div>
        
        <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Nama Instansi -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Nama Instansi <span class="text-red-500">*</span></label>
            <input 
              v-model="form.nama_instansi" 
              @input="generateSlug"
              type="text" 
              required 
              class="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all placeholder-gray-400" 
              placeholder=""
            >
          </div>

          <!-- Slug / Kode -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Slug Sekolah <span class="text-red-500">*</span></label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm font-medium">/</span>
              <input 
                v-model="form.slug" 
                @input="form.slug = $event.target.value.toLowerCase()"
                type="text" 
                required 
                class="w-full pl-8 pr-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all placeholder-gray-400 font-mono" 
                placeholder="contoh: smk-negeri-1-jakarta"
              >
            </div>
            <p class="text-[10px] text-gray-500 mt-1">Digunakan sebagai URL login siswa (e.g., domain.com/slug). Auto-generate dari nama jika kosong.</p>
          </div>

          <!-- Tingkat Pendidikan -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Tingkat Sekolah <span class="text-red-500">*</span></label>
            <select 
              v-model="form.tingkat_sekolah" 
              class="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white cursor-pointer"
            >
              <option v-for="(val, key) in TINGKAT_SEKOLAH" :key="key" :value="val">{{ key }} - {{ val }}</option>
            </select>
          </div>

          <!-- Alamat -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Alamat Lengkap</label>
            <textarea 
              v-model="form.alamat" 
              rows="1"
              class="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all placeholder-gray-400 resize-none" 
              placeholder=""
            ></textarea>
          </div>
        </div>
      </div>

      <!-- SECTION 2: AKUN ADMIN -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex items-center gap-3">
          <div class="w-8 h-8 bg-emerald-100 text-emerald-600 rounded-lg flex items-center justify-center">
            <i class="fas fa-user-shield"></i>
          </div>
          <div>
            <h3 class="text-sm font-bold text-gray-900">Akun Administrator</h3>
            <p class="text-[10px] text-gray-500 uppercase tracking-wide">Kredensial login untuk admin sekolah</p>
          </div>
        </div>
        
        <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Email Admin -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Email Admin <span class="text-red-500">*</span></label>
            <input 
              v-model="form.email" 
              type="email" 
              required 
              class="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition-all placeholder-gray-400" 
              placeholder="masukkan email anda"
            >
            <p class="text-[10px] text-gray-500 mt-1">Email ini akan menerima kredensial otomatis.</p>
          </div>

          <!-- Password Admin -->
          <div class="col-span-2 md:col-span-1">
            <label class="block text-xs font-semibold text-gray-700 mb-1.5">Password Awal <span class="text-red-500">*</span></label>
            <input 
              v-model="form.password" 
              type="password" 
              required 
              class="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm text-gray-900 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition-all placeholder-gray-400" 
              placeholder="Minimal 6 karakter"
            >
          </div>
        </div>
      </div>

      <!-- ACTION BUTTONS -->
      <div class="flex flex-col sm:flex-row justify-end gap-4 pt-4 border-t border-gray-200">
        <button 
          type="button" 
          @click="resetForm" 
          class="px-6 py-2.5 rounded-lg border border-gray-300 text-gray-700 text-sm font-semibold hover:bg-gray-50 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
        >
          Batal
        </button>
        
        <button 
          type="submit" 
          :disabled="isSubmitting" 
          class="px-8 py-2.5 rounded-lg bg-blue-600 text-white text-sm font-semibold hover:bg-blue-700 transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 flex items-center justify-center gap-2 min-w-[160px]"
        >
          <i v-if="isSubmitting" class="fas fa-circle-notch animate-spin"></i>
          <span v-if="isSubmitting">Memproses...</span>
          <span v-else>Simpan & Deploy</span>
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '@/api'
import Swal from 'sweetalert2'

// Ambil enum dari BE (disinkronkan manual di FE)
const TINGKAT_SEKOLAH = {
  SD: "SD",
  SMP: "SMP",
  MTS: "MTS",
  SMA: "SMA",
  MA: "MA",
  SMK: "SMK",
}

const showModal = ref(false)
const isSubmitting = ref(false)
const loadingTable = ref(false)
const listInstansi = ref([])

// Fungsi untuk generate slug otomatis dari nama
const generateSlug = () => {
  if (form.nama_instansi && !form.slug) {
    form.slug = form.nama_instansi
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
  }
}

const form = reactive({
  nama_instansi: '',
  slug: '',
  kode_instansi: '',
  tingkat_sekolah: TINGKAT_SEKOLAH.SMK,
  alamat: '',
  email: '',
  password: ''
})

const resetForm = () => {
  Object.assign(form, {
    nama_instansi: '',
    slug: '',
    kode_instansi: '',
    tingkat_sekolah: TINGKAT_SEKOLAH.SMK,
    alamat: '',
    email: '',
    password: ''
  })
}

const openModal = () => {
  resetForm()
  showModal.value = true
}

const fetchInstansi = async () => {
  loadingTable.value = true
  try {
    const res = await api.get('/super/instansi')
    listInstansi.value = res.data.data
  } catch (err) {
    console.error(err)
  } finally {
    loadingTable.value = false
  }
}

const handleSubmit = async () => {
  isSubmitting.value = true
  
  // Validasi slug: generate otomatis jika kosong
  let slug = form.slug
  if (!slug) {
    slug = form.nama_instansi
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
  }
  
  const payload = {
    nama_instansi: form.nama_instansi,
    kode_instansi: form.kode_instansi,
    tingkat_sekolah: form.tingkat_sekolah,
    alamat: form.alamat,
    slug: slug,
    email: form.email,
    password: form.password
  }
  
  try {
    const res = await api.post('/super/instansi', payload)
    Swal.fire('Berhasil!', res.data.message, 'success')
    resetForm()
    fetchInstansi()
  } catch (err) {
    Swal.fire('Gagal!', err.response?.data?.message || "Cek koneksi", 'error')
  } finally {
    isSubmitting.value = false
  }
}

const confirmDelete = async (item) => {
  const result = await Swal.fire({
    title: 'Hapus Instansi?',
    text: `Semua data ${item.nama_instansi} bakal ilang!`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#0f172a',
    cancelButtonText: 'Batal',
    confirmButtonText: 'Ya, Hapus!'
  })

  if (result.isConfirmed) {
    try {
      await api.delete(`/super/instansi/${item.id}`)
      Swal.fire('Terhapus!', 'Instansi ludes.', 'success')
      fetchInstansi()
    } catch (err) {
      Swal.fire('Error!', 'Gagal hapus data.', 'error')
    }
  }
}

onMounted(() => fetchInstansi())

// Ekspor fungsi yang diperlukan (untuk parent component jika ada)
defineExpose({
  openModal,
  fetchInstansi
})
</script>