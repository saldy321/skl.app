<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-xl font-bold text-gray-900">Daftar Instansi Sekolah</h1>
          <p class="text-sm text-gray-500 mt-1">Kelola dan akses dashboard setiap sekolah</p>
        </div>
        <router-link to="/super-admin/instansi/tambah" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-semibold flex items-center gap-2">
          <i class="fas fa-plus"></i>
          Tambah Instansi
        </router-link>
      </div>
    </div>

    <!-- Tabel Instansi -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase">Informasi Sekolah</th>
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase text-center">Jenjang</th>
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase">Slug / URL</th>
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase text-right">Aksi</th>
            </tr>
          </thead>
        <tbody class="divide-y divide-gray-100">
  <tr v-for="item in instansiList" :key="item.id" class="hover:bg-gray-50">
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
            {{ item.nama_instansi?.substring(0,2).toUpperCase() }}
          </div>
        </div>
        <div>
          <h4 class="text-sm font-semibold text-gray-900">{{ item.nama_instansi }}</h4>
          <p class="text-[10px] text-gray-400 font-mono">ID: {{ item.id?.substring(0,8) }}...</p>
        </div>
      </div>
    </td>
    <td class="px-6 py-4 text-center">
      <span class="inline-flex px-2.5 py-0.5 rounded-full bg-gray-100 text-gray-700 text-[10px] font-semibold">
        {{ item.tingkat_sekolah || 'Umum' }}
      </span>
    </td>
    <td class="px-6 py-4">
      <div class="flex items-center gap-2">
        <i class="fas fa-link text-gray-300 text-xs"></i>
        <span class="text-xs text-gray-600 font-medium">/{{ item.slug }}</span>
      </div>
    </td>
    <td class="px-6 py-4 text-right">
      <div class="flex items-center justify-end gap-2">
        <button 
          @click="impersonateAsAdmin(item)" 
          title="Masuk sebagai Admin Sekolah" 
          class="p-1.5 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded transition-colors"
        >
          <i class="fas fa-user-secret text-xs"></i>
        </button>
        
        <router-link :to="`/super-admin/instansi/edit/${item.id}`" title="Edit Data" class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors">
          <i class="fas fa-pen text-xs"></i>
        </router-link>

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
  <tr v-if="instansiList.length === 0">
    <td colspan="4" class="px-6 py-12 text-center text-gray-500">
      Belum ada data instansi
    </td>
  </tr>
</tbody>
        </table>
      </div>
    </div>

    <!-- MODAL RESET PASSWORD -->
    <div v-if="showResetModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b bg-gray-50 flex justify-between items-center">
          <h3 class="text-sm font-bold text-gray-900">Reset Password Admin</h3>
          <button @click="closeResetModal" class="text-gray-400 hover:text-gray-600">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div class="p-3 bg-blue-50 rounded-lg">
            <p class="text-sm font-semibold">{{ selectedInstansi?.nama_instansi }}</p>
          </div>
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Password Baru</label>
            <input v-model="newPassword" type="password" placeholder="Minimal 6 karakter" class="w-full px-3 py-2 rounded-md border border-gray-300">
          </div>
          <div class="flex gap-3">
            <button @click="closeResetModal" class="flex-1 px-4 py-2 border rounded-md">Batal</button>
            <button @click="submitResetPassword" :disabled="isResetting" class="flex-1 px-4 py-2 bg-orange-600 text-white rounded-md">
              {{ isResetting ? 'Memproses...' : 'Reset Password' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import Swal from 'sweetalert2'

const router = useRouter()
const instansiList = ref([])
const showResetModal = ref(false)
const selectedInstansi = ref(null)
const newPassword = ref('')
const isResetting = ref(false)

const fetchInstansi = async () => {
  try {
    const res = await api.get('/super/instansi')
    if (res.data.status === 'success') {
      instansiList.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  }
}

const impersonateAsAdmin = async (instansi) => {
  try {
    const result = await Swal.fire({
      title: 'Masuk sebagai Admin?',
      text: `Anda akan login sebagai admin untuk ${instansi.nama_instansi}`,
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#8b5cf6',
      confirmButtonText: 'Ya, Lanjutkan',
      cancelButtonText: 'Batal'
    })

    if (result.isConfirmed) {
      const res = await api.post(`/super/impersonate/${instansi.id}`)
      if (res.data.status === 'success') {
        sessionStorage.setItem('isImpersonating', 'true')
        sessionStorage.setItem('impersonatedInstansi', instansi.nama_instansi)
        window.location.href = `/${instansi.slug}/dashboard`
      }
    }
  } catch (err) {
    console.error(err)
    Swal.fire('Gagal!', err.response?.data?.message || 'Tidak bisa masuk sebagai admin', 'error')
  }
}

const openDeleteModal = (instansi) => {
  Swal.fire({
    title: 'Hapus Instansi?',
    text: `Semua data ${instansi.nama_instansi} (siswa, nilai, admin) akan ikut terhapus permanen!`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#ef4444',
    confirmButtonText: 'Ya, Hapus!',
    cancelButtonText: 'Batal'
  }).then(async (result) => {
    if (result.isConfirmed) {
      try {
        await api.delete(`/super/instansi/${instansi.id}`)
        Swal.fire('Terhapus!', 'Instansi berhasil dihapus.', 'success')
        fetchInstansi()
      } catch (err) {
        Swal.fire('Error!', err.response?.data?.message || 'Gagal hapus data.', 'error')
      }
    }
  })
}

const openResetModal = (instansi) => {
  selectedInstansi.value = instansi
  showResetModal.value = true
}

const closeResetModal = () => {
  showResetModal.value = false
  selectedInstansi.value = null
  newPassword.value = ''
}

const submitResetPassword = async () => {
  if (!newPassword.value || newPassword.value.length < 6) {
    Swal.fire('Error', 'Password minimal 6 karakter!', 'error')
    return
  }

  isResetting.value = true
  try {
    const res = await api.post(`/super/instansi/${selectedInstansi.value.id}/reset-password`, {
      new_password: newPassword.value
    })
    if (res.data.status === 'success') {
      Swal.fire('Berhasil!', res.data.message, 'success')
      closeResetModal()
    }
  } catch (err) {
    Swal.fire('Gagal!', err.response?.data?.message || 'Terjadi kesalahan', 'error')
  } finally {
    isResetting.value = false
  }
}

onMounted(() => {
  fetchInstansi()
})
</script>