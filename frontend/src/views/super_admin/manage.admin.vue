<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
      <h1 class="text-xl font-bold text-gray-900">Manajemen Admin Sekolah</h1>
      <p class="text-sm text-gray-500 mt-1">Kelola email admin untuk setiap instansi</p>
    </div>

    <!-- Tabel Admin -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase">Instansi</th>
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase">Email Admin</th>
              <th class="px-6 py-3 text-xs font-bold text-gray-500 uppercase text-right">Aksi</th>
            </tr>
          </thead>
         <tbody class="divide-y divide-gray-100">
  <tr v-for="admin in adminList" :key="admin.id" class="hover:bg-gray-50">
    <td class="px-6 py-4">
      <div class="flex items-center gap-3">
        <!-- LOGO INSTANSI -->
        <div class="w-10 h-10 rounded-full flex items-center justify-center overflow-hidden bg-white border border-gray-200">
          <img 
            v-if="admin.instansi?.logo_instansi"
            :src="`http://localhost:3000/uploads/instansi/${admin.instansi.logo_instansi}`"
            class="w-full h-full object-cover"
            alt="Logo"
          />
          <div v-else class="w-full h-full bg-blue-50 text-blue-600 flex items-center justify-center text-xs font-bold">
            {{ admin.instansi?.nama_instansi?.substring(0,2) || 'SC' }}
          </div>
        </div>
        
        <!-- INFO TEKS -->
        <div>
          <h4 class="text-sm font-semibold text-gray-900">
            {{ admin.instansi?.nama_instansi || 'Tidak diketahui' }}
          </h4>
          <p class="text-[10px] text-gray-400 font-mono mt-0.5">{{ admin.email }}</p>
        </div>
      </div>
    </td>
    
    <!-- KOLOM EMAIL (BISA DIHAPUS KARENA SUDAH DI ATAS) -->
    <td class="px-6 py-4">
      <!-- Kosongkan atau hapus kolom ini -->
    </td>
    
    <td class="px-6 py-4 text-right">
      <div class="flex items-center justify-end gap-2">
        <button 
          @click="openEditEmailModal(admin)"
          title="Edit Email"
          class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
        >
          <i class="fas fa-pen text-xs"></i>
        </button>
      </div>
    </td>
  </tr>
  <tr v-if="adminList.length === 0">
    <td colspan="3" class="px-6 py-12 text-center text-gray-500">
      Belum ada data admin
    </td>
  </tr>
</tbody>
        </table>
      </div>
    </div>

    <!-- MODAL EDIT EMAIL -->
    <div v-if="showEmailModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b bg-gray-50 flex justify-between items-center">
          <h3 class="text-sm font-bold text-gray-900">Edit Email Admin</h3>
          <button @click="closeEmailModal" class="text-gray-400 hover:text-gray-600">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Instansi</label>
            <input :value="selectedAdmin?.instansi?.nama_instansi" disabled class="w-full px-3 py-2 rounded-md border border-gray-200 bg-gray-50 text-sm">
          </div>
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-2">Email Baru</label>
            <input v-model="newEmail" type="email" placeholder="admin@sekolah.sch.id" class="w-full px-3 py-2 rounded-md border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          </div>
          <div class="flex gap-3 pt-2">
            <button @click="closeEmailModal" class="flex-1 px-4 py-2 border rounded-md text-sm">Batal</button>
            <button @click="submitEditEmail" :disabled="updating" class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md text-sm disabled:opacity-50">
              {{ updating ? 'Menyimpan...' : 'Simpan' }}
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

const adminList = ref([])
const showEmailModal = ref(false)
const selectedAdmin = ref(null)
const newEmail = ref('')
const updating = ref(false)

const fetchAdmins = async () => {
  try {
    const res = await api.get('/super/admin-list')
    if (res.data.status === 'success') {
      adminList.value = res.data.data
    }
  } catch (err) {
    console.error("Gagal fetch admin:", err)
  }
}

const openEditEmailModal = (admin) => {
  selectedAdmin.value = admin
  newEmail.value = admin.email
  showEmailModal.value = true
}

const closeEmailModal = () => {
  showEmailModal.value = false
  selectedAdmin.value = null
  newEmail.value = ''
}

const submitEditEmail = async () => {
  if (!newEmail.value) {
    alert('Email tidak boleh kosong')
    return
  }
  
  updating.value = true
  try {
    const res = await api.put(`/super/admin/${selectedAdmin.value.id}/email`, {
      email: newEmail.value
    })
    if (res.data.status === 'success') {
      alert('✅ ' + res.data.message)
      closeEmailModal()
      fetchAdmins()
    }
  } catch (err) {
    console.error(err)
    alert('❌ Gagal update email: ' + (err.response?.data?.message || 'Terjadi kesalahan'))
  } finally {
    updating.value = false
  }
}

onMounted(() => {
  fetchAdmins()
})
</script>