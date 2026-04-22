<template>
  <div class="min-h-screen bg-[#F3F4F6] font-sans text-slate-800">
    
    <!-- Top Bar Corporate -->
    <main class="max-w-7xl mx-auto p-6 space-y-6">
      
      <!-- Stats Cards (Clean Corporate Style) -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
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

      <!-- Button Tambah Instansi -->
      <div class="flex justify-end">
        <button 
          @click="openCreateModal"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-semibold flex items-center gap-2 transition-colors shadow-sm"
        >
          <i class="fas fa-plus text-xs"></i>
          Tambah Instansi
        </button>
      </div>

      <!-- Main Table Section (Professional Data Grid) -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
        <!-- Table Header -->
        <div class="px-6 py-5 border-b border-gray-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-gray-50/50">
          <div>
            <h2 class="text-base font-bold text-gray-900">Daftar Instansi Sekolah</h2>
            <p class="text-xs text-gray-500 mt-0.5">
              Menampilkan <span class="font-semibold text-blue-600">{{ recentInstansi.length }}</span> data terdaftar.
            </p>
          </div>
        </div>

        <!-- Table Content -->
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-200">
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider">Informasi Sekolah</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider text-center">Jenjang</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider">Slug / URL</th>
                <th class="px-6 py-3 text-[11px] font-bold text-gray-500 uppercase tracking-wider text-right">Aksi</th>
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
      <div class="flex items-center justify-end gap-2">
        <button 
          @click="openEditModal(item)" 
          title="Edit Data" 
          class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
        >
          <i class="fas fa-pen text-xs"></i>
        </button>
        
        <button 
          @click="openResetModal(item)" 
          title="Reset Password Admin" 
          class="p-1.5 text-gray-400 hover:text-orange-600 hover:bg-orange-50 rounded transition-colors"
        >
          <i class="fas fa-key text-xs"></i>
        </button>

        <button 
          @click="openDeleteModal(item)" 
          title="Hapus Instansi" 
          class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
        >
          <i class="fas fa-trash text-xs"></i>
        </button>
      </div>
    </td>
  </tr>

  <tr v-if="recentInstansi.length === 0">
    <td colspan="4" class="px-6 py-12 text-center">
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
        
        <!-- Footer Table -->
        <div class="px-6 py-3 border-t border-gray-200 bg-gray-50 flex justify-between items-center">
          <p class="text-[10px] text-gray-500">Showing all records</p>
        </div>
      </div>
    </main>

    <!-- MODAL CREATE / EDIT INSTANSI -->
    <div v-if="showInstansiModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden border border-gray-200">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <h3 class="text-sm font-bold text-gray-900">
            {{ isEditMode ? 'Edit Instansi' : 'Tambah Instansi Baru' }}
          </h3>
          <button @click="closeInstansiModal" class="text-gray-400 hover:text-gray-600">
            <i class="fas fa-times text-sm"></i>
          </button>
        </div>
        
        <div class="p-6 space-y-4">
          <!-- Nama Instansi -->
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Nama Instansi *</label>
            <input 
              v-model="instansiForm.nama_instansi" 
              type="text" 
              placeholder="Contoh: SMK Negeri 1 Jakarta"
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
          </div>

          <!-- Kode Instansi -->
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Kode Instansi *</label>
            <input 
              v-model="instansiForm.kode_instansi" 
              type="text" 
              placeholder="Contoh: SMKN1JKT"
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
          </div>

          <!-- Tingkat Sekolah -->
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Tingkat Sekolah *</label>
            <select 
              v-model="instansiForm.tingkat_sekolah"
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Pilih Tingkat</option>
              <option value="SD">SD</option>
              <option value="SMP">SMP</option>
              <option value="SMA">SMA</option>
              <option value="SMK">SMK</option>
            </select>
          </div>

          <!-- Alamat -->
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Alamat</label>
            <textarea 
              v-model="instansiForm.alamat" 
              rows="2"
              placeholder="Alamat lengkap sekolah..."
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            ></textarea>
          </div>

          <!-- Email Admin (hanya untuk CREATE) -->
          <div v-if="!isEditMode">
            <label class="block text-xs font-semibold text-gray-700 mb-2">Email Admin *</label>
            <input 
              v-model="instansiForm.email" 
              type="email" 
              placeholder="admin@sekolah.sch.id"
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
          </div>

          <!-- Password Admin (hanya untuk CREATE) -->
          <div v-if="!isEditMode">
            <label class="block text-xs font-semibold text-gray-700 mb-2">Password Admin *</label>
            <input 
              v-model="instansiForm.password" 
              type="password" 
              placeholder="Minimal 6 karakter"
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
          </div>

          <div class="flex gap-3 pt-2">
            <button @click="closeInstansiModal" class="flex-1 px-4 py-2.5 rounded-md border border-gray-300 text-gray-700 text-xs font-semibold hover:bg-gray-50 transition-colors">
              Batal
            </button>
            <button @click="submitInstansi" :disabled="isSubmitting" class="flex-1 px-4 py-2.5 rounded-md bg-blue-600 text-white text-xs font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2">
              <i v-if="isSubmitting" class="fas fa-circle-notch fa-spin"></i>
              {{ isSubmitting ? 'Memproses...' : (isEditMode ? 'Update' : 'Simpan') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL RESET PASSWORD -->
    <div v-if="showResetModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden border border-gray-200">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <h3 class="text-sm font-bold text-gray-900">Reset Password Admin</h3>
          <button @click="closeResetModal" class="text-gray-400 hover:text-gray-600">
            <i class="fas fa-times text-sm"></i>
          </button>
        </div>
        
        <div class="p-6">
          <div class="mb-6 p-4 bg-blue-50 rounded-lg border border-blue-100 flex items-center gap-3">
            <div class="w-8 h-8 bg-white text-blue-600 rounded flex items-center justify-center text-xs font-bold shadow-sm border border-blue-100">
              {{ selectedInstansi?.nama_instansi?.substring(0,2) }}
            </div>
            <div>
              <p class="text-sm font-semibold text-gray-900">{{ selectedInstansi?.nama_instansi }}</p>
              <p class="text-[10px] text-gray-500">Target Reset Password</p>
            </div>
          </div>

          <div class="mb-6">
            <label class="block text-xs font-semibold text-gray-700 mb-2">Password Baru</label>
            <input 
              v-model="newPassword" 
              type="text" 
              placeholder="Masukkan password baru..."
              class="w-full px-3 py-2.5 rounded-md border border-gray-300 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
            >
            <p class="text-[10px] text-gray-500 mt-1.5">Minimal 6 karakter. Perubahan bersifat permanen.</p>
          </div>

          <div class="flex gap-3">
            <button @click="closeResetModal" class="flex-1 px-4 py-2.5 rounded-md border border-gray-300 text-gray-700 text-xs font-semibold hover:bg-gray-50 transition-colors">
              Batal
            </button>
            <button @click="submitResetPassword" :disabled="isResetting" class="flex-1 px-4 py-2.5 rounded-md bg-blue-600 text-white text-xs font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2">
              <i v-if="isResetting" class="fas fa-circle-notch fa-spin"></i>
              {{ isResetting ? 'Memproses...' : 'Simpan Password' }}
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

// Data
const recentInstansi = ref([])
const stats = ref([
  { label: 'Total Instansi', value: '0', icon: 'fas fa-school', color: 'text-blue-600', bgColor: 'bg-blue-50' },
  { label: 'Admin Aktif', value: '0', icon: 'fas fa-user-shield', color: 'text-emerald-600', bgColor: 'bg-emerald-50' },
  { label: 'Siswa Nasional', value: '0', icon: 'fas fa-user-graduate', color: 'text-orange-600', bgColor: 'bg-orange-50' },
  { label: 'Status Server', value: 'Online', icon: 'fas fa-server', color: 'text-indigo-600', bgColor: 'bg-indigo-50' }
])

// Modal Instansi (Create/Edit)
const showInstansiModal = ref(false)
const isEditMode = ref(false)
const isSubmitting = ref(false)
const instansiForm = ref({
  id: '',
  nama_instansi: '',
  kode_instansi: '',
  tingkat_sekolah: '',
  alamat: '',
  email: '',
  password: ''
})

// Modal Reset Password
const showResetModal = ref(false)
const selectedInstansi = ref(null)
const newPassword = ref('')
const isResetting = ref(false)

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

// ========== CREATE / EDIT INSTANSI ==========
const openCreateModal = () => {
  isEditMode.value = false
  instansiForm.value = {
    id: '',
    nama_instansi: '',
    kode_instansi: '',
    tingkat_sekolah: '',
    alamat: '',
    email: '',
    password: ''
  }
  showInstansiModal.value = true
}

const openEditModal = (instansi) => {
  isEditMode.value = true
  instansiForm.value = {
    id: instansi.id,
    nama_instansi: instansi.nama_instansi,
    kode_instansi: instansi.kode_instansi,
    tingkat_sekolah: instansi.tingkat_sekolah,
    alamat: instansi.alamat || '',
    email: '',
    password: ''
  }
  showInstansiModal.value = true
}

const closeInstansiModal = () => {
  showInstansiModal.value = false
  instansiForm.value = {
    id: '',
    nama_instansi: '',
    kode_instansi: '',
    tingkat_sekolah: '',
    alamat: '',
    email: '',
    password: ''
  }
}

const submitInstansi = async () => {
  // Validasi
  if (!instansiForm.value.nama_instansi) {
    alert("Nama instansi wajib diisi!")
    return
  }
  if (!instansiForm.value.kode_instansi) {
    alert("Kode instansi wajib diisi!")
    return
  }
  if (!instansiForm.value.tingkat_sekolah) {
    alert("Tingkat sekolah wajib dipilih!")
    return
  }
  
  // Validasi khusus create
  if (!isEditMode.value) {
    if (!instansiForm.value.email) {
      alert("Email admin wajib diisi!")
      return
    }
    if (!instansiForm.value.password || instansiForm.value.password.length < 6) {
      alert("Password admin minimal 6 karakter!")
      return
    }
  }

  isSubmitting.value = true

  try {
    let res
    if (isEditMode.value) {
      // UPDATE
      res = await api.put(`/super/instansi/${instansiForm.value.id}`, {
        nama_instansi: instansiForm.value.nama_instansi,
        kode_instansi: instansiForm.value.kode_instansi,
        tingkat_sekolah: instansiForm.value.tingkat_sekolah,
        alamat: instansiForm.value.alamat
      })
    } else {
      // CREATE
      res = await api.post('/super/instansi', {
        nama_instansi: instansiForm.value.nama_instansi,
        kode_instansi: instansiForm.value.kode_instansi,
        tingkat_sekolah: instansiForm.value.tingkat_sekolah,
        alamat: instansiForm.value.alamat,
        email: instansiForm.value.email,
        password: instansiForm.value.password
      })
    }

    if (res.data.status === 'success') {
      alert(res.data.message)
      closeInstansiModal()
      fetchData() // Refresh data
    }
  } catch (err) {
    console.error(err)
    alert("❌ Gagal: " + (err.response?.data?.message || "Terjadi kesalahan"))
  } finally {
    isSubmitting.value = false
  }
}

// ========== DELETE INSTANSI ==========
const openDeleteModal = (instansi) => {
  if (confirm(`⚠️ Yakin ingin menghapus instansi "${instansi.nama_instansi}"?\n\nSemua data siswa, nilai, dan admin sekolah ini akan ikut terhapus PERMANEN!`)) {
    deleteInstansi(instansi.id)
  }
}

const deleteInstansi = async (id) => {
  try {
    const res = await api.delete(`/super/instansi/${id}`)
    if (res.data.status === 'success') {
      alert("✅ " + res.data.message)
      fetchData() // Refresh data
    }
  } catch (err) {
    console.error(err)
    alert("❌ Gagal hapus instansi: " + (err.response?.data?.message || "Terjadi kesalahan"))
  }
}

// ========== RESET PASSWORD ==========
const openResetModal = (instansi) => {
  selectedInstansi.value = instansi
  newPassword.value = ''
  showResetModal.value = true
}

const closeResetModal = () => {
  showResetModal.value = false
  selectedInstansi.value = null
  newPassword.value = ''
}

const submitResetPassword = async () => {
  if (!newPassword.value || newPassword.value.length < 6) {
    alert("Password minimal 6 karakter!")
    return
  }

  if (!confirm(`Yakin ingin mereset password admin untuk ${selectedInstansi.value.nama_instansi}?`)) {
    return
  }

  isResetting.value = true
  
  try {
    const res = await api.post(`/super/instansi/${selectedInstansi.value.id}/reset-password`, {
      new_password: newPassword.value
    })

    if (res.data.status === 'success') {
      alert("✅ " + res.data.message)
      closeResetModal()
    }
  } catch (err) {
    console.error(err)
    alert("❌ Gagal reset password: " + (err.response?.data?.message || "Terjadi kesalahan"))
  } finally {
    isResetting.value = false
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