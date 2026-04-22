<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-bold text-gray-900">Manajemen NISN</h1>
          <p class="text-sm text-gray-500 mt-1">Fitur darurat untuk menangani NISN bermasalah (typo/fiktif/duplikat)</p>
        </div>
        <div class="bg-red-50 border border-red-200 rounded-lg px-4 py-2">
          <span class="text-red-600 text-xs font-bold uppercase tracking-wider">
            <i class="fas fa-exclamation-triangle mr-1"></i>Emergency Access
          </span>
        </div>
      </div>
    </div>

    <!-- Statistik Cepat -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-bold text-gray-500 uppercase tracking-wider">Total Siswa</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.totalSiswa || 0 }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
            <i class="fas fa-users"></i>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-bold text-gray-500 uppercase tracking-wider">Total Sekolah</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.totalSekolah || 0 }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-green-50 text-green-600 flex items-center justify-center">
            <i class="fas fa-school"></i>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-bold text-gray-500 uppercase tracking-wider">Aksi Hari Ini</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.aksiHariIni || 0 }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center">
            <i class="fas fa-history"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- Search Box -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
      <h3 class="text-sm font-bold text-gray-700 mb-4 flex items-center gap-2">
        <i class="fas fa-search text-blue-500"></i>
        Cari NISN
      </h3>
      
      <div class="flex gap-3">
        <input 
          v-model="searchNISN" 
          @keyup.enter="searchHandler"
          placeholder="Masukkan NISN (10 digit)..." 
          maxlength="10"
          class="flex-1 border border-gray-300 rounded-lg px-4 py-3 text-lg font-mono focus:ring-2 focus:ring-blue-500 outline-none"
        />
        <button 
          @click="searchHandler" 
          :disabled="searching"
          class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-bold transition-colors disabled:opacity-50"
        >
          <i v-if="searching" class="fas fa-spinner animate-spin mr-2"></i>
          {{ searching ? 'Mencari...' : 'Cari NISN' }}
        </button>
      </div>
      
      <p class="text-xs text-gray-400 mt-3">
        <i class="fas fa-info-circle mr-1"></i>
        Masukkan NISN yang ingin dicek. Jika ditemukan, Anda dapat menghapusnya permanen.
      </p>
    </div>

    <!-- Hasil Pencarian -->
    <div v-if="searchResult" class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b bg-gray-50">
        <h3 class="font-bold text-gray-700">Hasil Pencarian</h3>
      </div>
      
      <!-- NISN Tidak Ditemukan -->
      <div v-if="!searchResult.found" class="p-10 text-center">
        <div class="w-16 h-16 rounded-full bg-green-50 text-green-600 flex items-center justify-center mx-auto mb-4">
          <i class="fas fa-check-circle text-3xl"></i>
        </div>
        <p class="text-gray-600 text-lg font-semibold">NISN tidak ditemukan di sistem</p>
        <p class="text-gray-400 text-sm mt-1">NISN ini tersedia untuk didaftarkan.</p>
      </div>
      
      <!-- NISN Ditemukan -->
      <div v-else class="p-6">
        <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6 flex items-start gap-3">
          <i class="fas fa-exclamation-triangle text-yellow-600 text-xl"></i>
          <div>
            <p class="text-yellow-800 font-bold">NISN Sudah Terdaftar!</p>
            <p class="text-sm text-yellow-700">Jika ini kesalahan (typo/fiktif), Anda dapat menghapusnya.</p>
          </div>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">NISN</label>
            <p class="font-mono text-xl font-bold text-gray-900">{{ searchResult.data.nisn }}</p>
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Siswa</label>
            <p class="font-bold text-gray-900">{{ searchResult.data.nama_siswa }}</p>
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Sekolah</label>
            <p class="text-gray-700">{{ searchResult.data.nama_instansi }}</p>
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Kelas / Jurusan</label>
            <p class="text-gray-700">{{ searchResult.data.kelas }} / {{ searchResult.data.jurusan || '-' }}</p>
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Status</label>
            <span :class="searchResult.data.status_lulus ? 'bg-green-100 text-green-700 border-green-200' : 'bg-blue-100 text-blue-700 border-blue-200'" class="px-3 py-1 rounded-full text-xs font-bold border">
              {{ searchResult.data.status_lulus ? 'Lulus' : 'Aktif' }}
            </span>
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Terdaftar Sejak</label>
            <p class="text-gray-500 text-sm">{{ formatDate(searchResult.data.created_at) }}</p>
          </div>
        </div>
        
        <!-- Foto -->
        <div v-if="searchResult.data.foto_siswa" class="mt-6">
          <label class="block text-xs font-bold text-gray-500 uppercase mb-2">Foto Siswa</label>
          <img 
            :src="`http://localhost:3000/uploads/siswa/${searchResult.data.foto_siswa}`" 
            class="w-24 h-32 object-cover border rounded-lg shadow"
            alt="Foto Siswa"
          />
        </div>
        
        <!-- Tombol Hapus -->
        <div class="mt-6 pt-6 border-t flex justify-end">
          <button 
            @click="openDeleteModal" 
            class="bg-red-600 hover:bg-red-700 text-white px-6 py-3 rounded-lg font-bold transition-colors flex items-center gap-2"
          >
            <i class="fas fa-trash"></i>
            Hapus Permanen NISN Ini
          </button>
        </div>
      </div>
    </div>

    <!-- Info / Panduan -->
    <div v-if="!searchResult" class="bg-blue-50 border border-blue-200 rounded-lg p-6">
      <h3 class="text-sm font-bold text-blue-900 mb-3 flex items-center gap-2">
        <i class="fas fa-lightbulb"></i>
        Kapan Menggunakan Fitur Ini?
      </h3>
      <ul class="space-y-2 text-sm text-blue-800">
        <li class="flex items-start gap-2">
          <i class="fas fa-check-circle text-blue-500 mt-1 text-xs"></i>
          <span>Admin sekolah komplain tidak bisa mendaftarkan NISN karena sudah terpakai</span>
        </li>
        <li class="flex items-start gap-2">
          <i class="fas fa-check-circle text-blue-500 mt-1 text-xs"></i>
          <span>NISN yang terdaftar ternyata typo (salah input) dan perlu dihapus</span>
        </li>
        <li class="flex items-start gap-2">
          <i class="fas fa-check-circle text-blue-500 mt-1 text-xs"></i>
          <span>NISN fiktif/palsu yang perlu dibersihkan dari sistem</span>
        </li>
        <li class="flex items-start gap-2">
          <i class="fas fa-exclamation-triangle text-red-500 mt-1 text-xs"></i>
          <span class="text-red-700 font-semibold">Peringatan: Tindakan ini TIDAK DAPAT DIBATALKAN!</span>
        </li>
      </ul>
    </div>

    <!-- Modal Konfirmasi Hapus -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b bg-red-50 flex items-center gap-3">
          <i class="fas fa-exclamation-triangle text-red-600"></i>
          <h3 class="font-bold text-red-900">Konfirmasi Penghapusan Permanen</h3>
        </div>
        
        <div class="p-6">
          <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-4 text-sm">
            <p class="font-bold mb-2 text-red-900">Data yang akan dihapus:</p>
            <ul class="list-disc list-inside space-y-1 text-red-800">
              <li>NISN: {{ searchResult?.data?.nisn }}</li>
              <li>Nama: {{ searchResult?.data?.nama_siswa }}</li>
              <li>Sekolah: {{ searchResult?.data?.nama_instansi }}</li>
              <li>Semua nilai & foto akan ikut terhapus</li>
            </ul>
            <p class="mt-3 text-red-600 font-bold">⚠️ Tindakan ini TIDAK DAPAT DIBATALKAN!</p>
          </div>
          
          <label class="block mb-2 text-xs font-bold text-gray-700 uppercase">Alasan Penghapusan <span class="text-red-500">*</span></label>
          <select v-model="deleteReason" class="w-full border border-gray-300 rounded-lg p-3 text-sm mb-3 focus:ring-2 focus:ring-red-500">
            <option value="">-- Pilih Alasan --</option>
            <option value="NISN Typo (Salah Input)">NISN Typo (Salah Input)</option>
            <option value="NISN Fiktif / Palsu">NISN Fiktif / Palsu</option>
            <option value="Duplikat - Pindah ke Sekolah Lain">Duplikat - Pindah ke Sekolah Lain</option>
            <option value="Permintaan Admin Sekolah">Permintaan Admin Sekolah</option>
          </select>
          
          <textarea 
            v-if="deleteReason"
            v-model="deleteNote"
            placeholder="Detail tambahan (opsional)..."
            class="w-full border border-gray-300 rounded-lg p-3 text-sm"
            rows="2"
          ></textarea>
          
          <div class="flex gap-3 mt-6">
            <button @click="showDeleteModal = false" class="flex-1 px-4 py-3 border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 transition-colors">
              Batal
            </button>
            <button 
              @click="executeDelete" 
              :disabled="!deleteReason || deleting"
              class="flex-1 bg-red-600 hover:bg-red-700 text-white px-4 py-3 rounded-lg font-bold transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <i v-if="deleting" class="fas fa-spinner animate-spin"></i>
              {{ deleting ? 'Menghapus...' : 'Hapus Permanen' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import Swal from 'sweetalert2'

const searchNISN = ref('')
const searching = ref(false)
const searchResult = ref(null)
const showDeleteModal = ref(false)
const deleteReason = ref('')
const deleteNote = ref('')
const deleting = ref(false)

const stats = ref({
  totalSiswa: 0,
  totalSekolah: 0,
  aksiHariIni: 0
})

const fetchStats = async () => {
  try {
    const res = await api.get('/super/dashboard')
    if (res.data.status === 'success') {
      stats.value.totalSiswa = res.data.data.totalSiswaNasional || 0
      stats.value.totalSekolah = res.data.data.totalSekolah || 0
    }
  } catch (err) {
    console.error('Gagal fetch stats:', err)
  }
}

const searchHandler = async () => {
  if (!searchNISN.value || searchNISN.value.length !== 10) {
    Swal.fire('Perhatian', 'NISN harus 10 digit', 'warning')
    return
  }
  
  searching.value = true
  searchResult.value = null
  
  try {
    const res = await api.get('/super/nisn/search', { params: { nisn: searchNISN.value } })
    searchResult.value = res.data
  } catch (err) {
    Swal.fire('Error', 'Gagal mencari NISN', 'error')
  } finally {
    searching.value = false
  }
}

const openDeleteModal = () => {
  deleteReason.value = ''
  deleteNote.value = ''
  showDeleteModal.value = true
}

const executeDelete = async () => {
  if (!deleteReason.value) {
    Swal.fire('Perhatian', 'Alasan penghapusan wajib dipilih', 'warning')
    return
  }
  
  const confirm = await Swal.fire({
    title: 'Yakin Hapus Permanen?',
    text: 'Data siswa, nilai, dan foto akan dihapus selamanya!',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: 'Batal'
  })
  
  if (!confirm.isConfirmed) return
  
  deleting.value = true
  
  try {
    const alasan = deleteNote.value ? `${deleteReason.value} - ${deleteNote.value}` : deleteReason.value
    
    await api.post('/super/nisn/force-delete', {
      siswa_id: searchResult.value.data.id,
      alasan: alasan
    })
    
    showDeleteModal.value = false
    searchResult.value = null
    searchNISN.value = ''
    
    stats.value.aksiHariIni++
    
    Swal.fire('Berhasil!', 'NISN telah dihapus permanen dari sistem.', 'success')
  } catch (err) {
    Swal.fire('Gagal', err.response?.data?.message || 'Terjadi kesalahan', 'error')
  } finally {
    deleting.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' })
}

onMounted(() => {
  fetchStats()
})
</script>