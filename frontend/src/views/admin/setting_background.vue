<template>
  <div class="max-w-3xl mx-auto">
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 md:p-8">
      
      <div class="mb-8">
        <h2 class="text-xl font-bold text-slate-800 flex items-center gap-3">
          <i class="fas fa-image text-indigo-500"></i>
          Pengaturan Background Halaman Siswa
        </h2>
        <p class="text-sm text-slate-500 mt-2">Atur background yang tampil saat siswa membuka amplop</p>
      </div>

      <div class="space-y-6">
        
   
        <div>
          <label class="font-bold text-slate-700 mb-2 block">Preview Background</label>
          <div class="relative w-full h-48 rounded-xl overflow-hidden border-2 border-slate-200 bg-slate-100">
            <img 
              v-if="currentBackground" 
              :src="getBackgroundUrl(currentBackground)" 
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-500 to-purple-600">
              <p class="text-white text-sm">Background Default</p>
            </div>
          </div>
        </div>

        <!-- Upload Background Baru -->
        <div>
          <label class="font-bold text-slate-700 mb-2 block">Upload Background Baru</label>
          <div class="border-2 border-dashed border-slate-200 rounded-xl p-6 text-center hover:border-indigo-300 transition">
            <input 
              type="file" 
              ref="fileInput"
              @change="handleFileUpload"
              accept="image/jpeg,image/png,image/webp"
              class="hidden"
            />
            <button 
              @click="$refs.fileInput.click()"
              class="bg-indigo-50 hover:bg-indigo-100 text-indigo-600 px-6 py-3 rounded-xl text-sm font-bold transition"
            >
              <i class="fas fa-upload mr-2"></i>
              Pilih Gambar
            </button>
            <p class="text-xs text-slate-400 mt-3">
              Format: JPG, PNG, WEBP. Maksimal 2MB
            </p>
          </div>
        </div>

        <!-- Tombol Hapus -->
        <div v-if="currentBackground" class="flex justify-end">
          <button 
            @click="deleteBackground"
            :disabled="loading"
            class="bg-red-500 hover:bg-red-600 text-white px-6 py-3 rounded-xl text-sm font-bold transition"
          >
            <i class="fas fa-trash mr-2"></i>
            Hapus Background
          </button>
        </div>

      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import Swal from 'sweetalert2'

const route = useRoute()
const slug = route.params.slug
const loading = ref(false)
const currentBackground = ref('')
const fileInput = ref(null)

const getBackgroundUrl = (filename) => {
  if (!filename) return ''
  return `http://localhost:3000/uploads/backgrounds/${filename}`
}

const loadBackground = async () => {
  try {
    const res = await api.get(`/${slug}/admin/setting/background`)
    if (res.data?.data) {
      currentBackground.value = res.data.data.background || ''
    }
  } catch (err) {
    console.error('Gagal load background:', err)
  }
}

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // Validasi size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    Swal.fire('Error', 'Ukuran file maksimal 2MB', 'error')
    return
  }

  const formData = new FormData()
  formData.append('background', file)

  loading.value = true
  try {
    const res = await api.post(`/${slug}/admin/setting/background/upload`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    if (res.data.success) {
      currentBackground.value = res.data.data.background
      Swal.fire('Berhasil!', 'Background berhasil diupload', 'success')
    }
  } catch (err) {
    Swal.fire('Error', err.response?.data?.message || 'Gagal upload', 'error')
  } finally {
    loading.value = false
    fileInput.value.value = ''
  }
}

const deleteBackground = async () => {
  const confirm = await Swal.fire({
    title: 'Yakin?',
    text: 'Background akan dihapus dan kembali ke default',
    icon: 'warning',
    showCancelButton: true
  })
  
  if (!confirm.isConfirmed) return

  loading.value = true
  try {
    const res = await api.delete(`/${slug}/admin/setting/background`)
    if (res.data.success) {
      currentBackground.value = ''
      Swal.fire('Berhasil!', 'Background berhasil dihapus', 'success')
    }
  } catch (err) {
    Swal.fire('Error', err.response?.data?.message || 'Gagal hapus', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(loadBackground)
</script>