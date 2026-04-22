<template>
  <div class="max-w-3xl mx-auto space-y-6">
    
    <!-- Header Section -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200 flex justify-between items-center">
      <div>
        <h2 class="text-xl font-black text-slate-800 tracking-tight">Pengaturan Sistem</h2>
        <p class="text-xs text-slate-500 font-medium mt-1">Atur waktu pengumuman & tampilan logo sekolah.</p>
      </div>
      <div class="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center text-xl">
        <i class="fas fa-cog"></i>
      </div>
    </div>

    <!-- ========== SETTING LOGO ========== -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-8 h-8 bg-blue-50 text-blue-600 rounded-lg flex items-center justify-center">
          <i class="fas fa-image text-sm"></i>
        </div>
        <h3 class="text-sm font-black text-slate-700 uppercase tracking-wide">Pengaturan Logo</h3>
      </div>

      <!-- Preview Logo -->
      <div class="mb-4 p-4 bg-slate-50 rounded-xl border border-slate-200">
        <p class="text-[10px] font-bold text-slate-400 uppercase mb-3">Preview Logo Saat Ini</p>
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-full overflow-hidden bg-white border-2 border-slate-200">
            <img 
              v-if="logoInstansi"
              :src="`http://localhost:3000/uploads/instansi/${logoInstansi}`"
              class="w-full h-full object-cover"
              alt="Logo"
            />
            <div v-else class="w-full h-full bg-indigo-50 text-indigo-600 flex items-center justify-center text-lg font-black">
              {{ namaInstansi?.substring(0,2).toUpperCase() || 'SC' }}
            </div>
          </div>
          <div>
            <p class="text-sm font-bold text-slate-800">{{ namaInstansi || 'Nama Sekolah' }}</p>
            <p class="text-[10px] text-slate-400 mt-0.5">Logo akan muncul di halaman login siswa</p>
          </div>
        </div>
      </div>

      <!-- Toggle Setting -->
      <div class="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-200">
        <div>
          <p class="text-sm font-bold text-slate-800">Tampilkan Logo di Halaman Login</p>
          <p class="text-xs text-slate-500">Aktifkan agar logo sekolah muncul di halaman login siswa</p>
        </div>
        
        <label class="relative inline-flex items-center cursor-pointer">
          <input 
            type="checkbox" 
            v-model="tampilkanLogo" 
            @change="saveLogoSetting"
            :disabled="savingLogo"
            class="sr-only peer"
          />
          <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
        </label>
      </div>

      <!-- Status Saving Logo -->
      <div v-if="savingLogo" class="mt-3 text-xs text-blue-600 flex items-center gap-2">
        <i class="fas fa-spinner animate-spin"></i> Menyimpan pengaturan logo...
      </div>
    </div>

    <!-- ========== SETTING WAKTU ========== -->
    <div class="bg-white p-8 rounded-2xl shadow-sm border border-slate-200">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-8 h-8 bg-amber-50 text-amber-600 rounded-lg flex items-center justify-center">
          <i class="fas fa-clock text-sm"></i>
        </div>
        <h3 class="text-sm font-black text-slate-700 uppercase tracking-wide">Pengaturan Waktu Pengumuman</h3>
      </div>

      <div class="mb-6">
        <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tanggal & Jam Pembukaan</label>
        <div class="flex flex-col sm:flex-row gap-4">
          <input 
            type="datetime-local" 
            v-model="localTime" 
            class="flex-1 bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-700 font-bold focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
          >
          <button 
            @click="saveTime" 
            :disabled="saving" 
            class="bg-slate-900 hover:bg-indigo-600 text-white px-8 py-3 rounded-xl text-xs font-black uppercase tracking-widest transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-slate-200"
          >
            <span v-if="saving"><i class="fas fa-circle-notch animate-spin mr-2"></i> Menyimpan...</span>
            <span v-else><i class="fas fa-save mr-2"></i> Simpan Jadwal</span>
          </button>
        </div>
        <p class="text-[10px] text-slate-400 mt-3 italic">
          * Siswa tidak akan bisa login atau melihat hasil sebelum waktu yang ditentukan tiba.
        </p>
      </div>

      <!-- Status Info -->
      <div class="bg-slate-50 rounded-xl p-4 border border-slate-100 flex items-start gap-4">
        <div class="w-10 h-10 bg-white rounded-lg border border-slate-200 flex items-center justify-center text-slate-400 flex-shrink-0">
          <i class="fas fa-info-circle"></i>
        </div>
        <div>
          <h4 class="text-sm font-bold text-slate-700 mb-1">Status Saat Ini</h4>
          <p class="text-xs text-slate-500 leading-relaxed">
            Jika waktu sudah disetting, halaman login siswa akan menampilkan <strong>Countdown Timer</strong>. 
            Pastikan waktu server backend sudah sesuai (WIB/WITA/WIT).
          </p>
        </div>
      </div>
    </div>

    <!-- Preview Card Jadwal -->
    <div v-if="lastSavedTime" class="bg-indigo-50/50 p-6 rounded-2xl border border-indigo-100 flex items-center justify-between">
      <div>
        <p class="text-[10px] font-black text-indigo-400 uppercase tracking-widest mb-1">Jadwal Terakhir Disimpan</p>
        <p class="text-lg font-bold text-indigo-900">{{ formatReadable(lastSavedTime) }}</p>
      </div>
      <div class="text-indigo-200 text-4xl">
        <i class="fas fa-calendar-check"></i>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import Swal from 'sweetalert2'
import { useRoute } from 'vue-router'

const route = useRoute()
const localTime = ref('')
const saving = ref(false)
const lastSavedTime = ref(null)

// State untuk setting logo
const tampilkanLogo = ref(true)
const logoInstansi = ref('')
const namaInstansi = ref('')
const savingLogo = ref(false)

// Fetch setting logo & instansi
const fetchLogoSetting = async () => {
  try {
    const slug = route.params.slug
    const res = await api.get(`/check-instansi/${slug}`)
    if (res.data.success) {
      logoInstansi.value = res.data.data.logo_instansi || ''
      namaInstansi.value = res.data.data.nama_instansi || ''
      tampilkanLogo.value = res.data.data.tampilkan_logo !== false
    }
  } catch (err) {
    console.error('Gagal fetch setting logo:', err)
  }
}

// Save setting logo
const saveLogoSetting = async () => {
  savingLogo.value = true
  
  try {
    const slug = route.params.slug
    const res = await api.post(`/${slug}/admin/setting/tampilkan-logo`, {
      tampilkan_logo: tampilkanLogo.value
    })
    
    if (res.data.status === 'success') {
      Swal.fire({
        icon: 'success',
        title: 'Berhasil!',
        text: 'Pengaturan logo telah disimpan.',
        timer: 1500,
        showConfirmButton: false
      })
    }
  } catch (err) {
    Swal.fire('Error', 'Gagal menyimpan pengaturan logo', 'error')
    tampilkanLogo.value = !tampilkanLogo.value // Revert
  } finally {
    savingLogo.value = false
  }
}

onMounted(() => {
  const now = new Date()
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset())
  localTime.value = now.toISOString().slice(0, 16)
  
  // Fetch setting logo
  fetchLogoSetting()
})

const saveTime = async () => {
  if (!localTime.value) {
    Swal.fire('Error!', 'Pilih tanggal dan waktu terlebih dahulu.', 'warning')
    return
  }
  
  saving.value = true
  try {
    const formatted = localTime.value.replace('T', ' ') + ':00'
    
    await api.post(`/${route.params.slug}/admin/setting/waktu-buka`, {
      waktu_buka: formatted
    })
    
    lastSavedTime.value = formatted
    Swal.fire({
      icon: 'success',
      title: 'Berhasil!',
      text: 'Jadwal pengumuman telah diperbarui. Siswa akan melihat countdown timer.',
      timer: 2000,
      showConfirmButton: false
    })
  } catch (err) {
    console.error(err)
    Swal.fire('Gagal!', err.response?.data?.message || 'Terjadi kesalahan saat menyimpan.', 'error')
  } finally {
    saving.value = false
  }
}

const formatReadable = (dateString) => {
  if (!dateString) return '-'
  const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' }
  return new Date(dateString.replace(' ', 'T')).toLocaleDateString('id-ID', options)
}
</script>