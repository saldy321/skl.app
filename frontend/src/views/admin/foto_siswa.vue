<template>
  <div class="min-h-screen bg-gray-100 p-2 md:p-8">
    <div class="mb-4 flex flex-col md:flex-row md:items-center gap-3">
      <!-- Tombol Trigger Upload -->
      <button 
        @click="$refs.fileInput.click()" 
        :disabled="isLoading"
        class="bg-[#f05a40] hover:bg-[#d64a32] text-white px-4 py-2 rounded text-sm font-bold shadow-sm w-full md:w-auto transition-colors disabled:opacity-50"
      >
        <span v-if="isLoading">
          <i class="fas fa-spinner fa-spin mr-2"></i> Sedang Mengunggah ZIP...
        </span>
        <span v-else>
          <i class="fas fa-file-archive mr-2"></i> Unggah ZIP Foto (Massal)
        </span>
      </button>

      <input 
        type="file" 
        ref="fileInput" 
        class="hidden" 
        accept=".zip" 
        @change="handleMassUpload" 
      />
    </div>

    <!-- Tabel Siswa -->
    <div class="bg-white rounded shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200 text-gray-700 text-[12px] uppercase font-bold">
              <th class="px-6 py-4 w-32 text-center">Foto</th>
              <th class="px-4 py-4">NISN</th>
              <th class="px-4 py-4">Nama</th>
              <th class="px-4 py-4">Kelas</th>
              <th v-if="showJurusanField" class="px-4 py-4">Jurusan</th>
              <th class="px-4 py-4">Tahun Lulus</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="siswa in listSiswa" :key="siswa.id" class="hover:bg-blue-50 text-sm transition-colors">
            
              <td class="px-6 py-3 text-center">
                <div class="relative inline-block w-14 h-20 border border-gray-300 bg-white p-0.5 shadow-sm rounded-sm group">
                  <button 
                    v-if="siswa.foto_siswa"
                    @click="deleteFoto(siswa)"
                    class="absolute -top-2 -right-2 bg-red-600 text-white w-5 h-5 rounded-full flex items-center justify-center shadow-md hover:bg-red-800 transition-all z-20 border border-white opacity-0 group-hover:opacity-100"
                    title="Hapus Foto"
                  >
                    <i class="fas fa-times text-[10px]"></i>
                  </button>

                  <img v-if="siswa.foto_siswa" 
                       :src="`http://localhost:3000/uploads/siswa/${siswa.foto_siswa}?t=${new Date().getTime()}`" 
                       class="w-full h-full object-cover" />
                  
                  <div v-else class="w-full h-full bg-gray-100 flex items-center justify-center text-[10px] text-gray-400 font-bold uppercase text-center leading-tight">
                    NO<br>FOTO
                  </div>
                </div>
              </td>

            
              <td class="px-4 py-3 font-mono text-gray-600">{{ siswa.nisn }}</td>
              <td class="px-4 py-3 uppercase font-semibold text-gray-700">{{ siswa.nama_siswa }}</td>
              <td class="px-4 py-3 text-gray-600">{{ siswa.kelas || '-' }}</td>
              <td v-if= "showJurusanField" class="px-4 py-3 text-gray-600">{{ siswa.jurusan || '-' }}</td>
              <td class="px-4 py-3 font-bold text-blue-600">{{ siswa.tahun_lulus || '-' }}</td>
            </tr>
            
            <tr v-if="listSiswa.length === 0">
              <td :colspan="showJurusanField ? 7 : 6" class="p-10 text-center text-gray-500">
                Belum ada data siswa. Silakan tambah data melalui menu Data Siswa.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { TINGKAT_SEKOLAH } from '@/utils/constants.js'
import { useRoute } from 'vue-router'
import api from '../../api/index.js'
import Swal from 'sweetalert2'

const route = useRoute()
const listSiswa = ref([])
const slug = route.params.slug
const fileInput = ref(null)
const isLoading = ref(false)
const auth = useAuthStore(
)


const showJurusanField = computed(() => {
  const tingkat = auth.tingkat?.toUpperCase() || ''
    return ['SMK', 'SMA', 'MA'].includes(tingkat)
})

const fetchSiswa = async () => {
  try {
    const response = await api.get(`/${slug}/admin/siswa`)
    listSiswa.value = response.data.data
  } catch (error) {
    console.error("Gagal mengambil data:", error)
  }
}

const deleteFoto = async (siswa) => {
  const result = await Swal.fire({
    title: 'Konfirmasi Hapus Foto',
    text: `Anda akan menghapus foto siswa ${siswa.nama_siswa} (NISN: ${siswa.nisn}). Tindakan ini tidak dapat dibatalkan.`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    confirmButtonText: 'Ya, Hapus',
    cancelButtonText: 'Batal'
  })

  if (result.isConfirmed) {
    try {
      await api.delete(`/${slug}/admin/siswa/${siswa.id}/foto`) 
      
      await Swal.fire({
        icon: 'success',
        title: 'Berhasil',
        text: 'Foto siswa telah dihapus.',
        timer: 1500,
        showConfirmButton: false
      })
      await fetchSiswa() 
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'Gagal',
        text: 'Tidak dapat menghapus foto. Silakan coba lagi.',
      })
    }
  }
}

const handleMassUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('foto_zip', file)
  
  isLoading.value = true
  try {
    const res = await api.post(`/${slug}/admin/upload-siswa`, formData)
    
    // HANDLE RESPONSE PARTIAL (SEBAGIAN GAGAL)
    if (res.data.status === 'partial') {
      let failedMsg = ''
      if (res.data.failed_nisn && res.data.failed_nisn.length > 0) {
        const displayList = res.data.failed_nisn.slice(0, 10)
        failedMsg = displayList.join(', ')
        if (res.data.failed_nisn.length > 10) {
          failedMsg += ` dan ${res.data.failed_nisn.length - 10} lainnya`
        }
      }

      await Swal.fire({
        icon: 'warning',
        title: 'Upload Sebagian Berhasil',
        html: `
          <div style="text-align: left;">
            <p style="margin-bottom: 15px;">
              <strong>${res.data.success_count}</strong> foto berhasil diunggah, 
              <strong style="color: #dc2626;">${res.data.failed_count}</strong> foto gagal diproses.
            </p>
            ${failedMsg ? `
              <div style="background: #fef2f2; padding: 12px; border-radius: 8px; border: 1px solid #fecaca;">
                <p style="margin: 0 0 8px 0; font-weight: 600; color: #991b1b; font-size: 13px;">
                  <i class="fas fa-exclamation-triangle" style="margin-right: 6px;"></i>NISN Tidak Terdaftar:
                </p>
                <p style="margin: 0; font-family: monospace; font-size: 12px; color: #7f1d1d; word-break: break-all;">
                  ${failedMsg}
                </p>
              </div>
            ` : ''}
            <p style="margin-top: 15px; font-size: 12px; color: #6b7280;">
              <i class="fas fa-info-circle" style="margin-right: 5px;"></i>
              Pastikan NISN telah terdaftar di menu Data Siswa sebelum mengunggah foto.
            </p>
          </div>
        `,
        confirmButtonText: 'Mengerti',
        confirmButtonColor: '#4f46e5',
      })
      
      await fetchSiswa()
    }
    
    // HANDLE RESPONSE SUKSES SEMUA
    else if (res.data.status === 'success') {
      await Swal.fire({
        icon: 'success',
        title: 'Upload Berhasil',
        text: res.data.message,
        confirmButtonText: 'Selesai',
        confirmButtonColor: '#10b981',
      })
      await fetchSiswa()
    }
    
  } catch (error) {
    console.error('Upload error:', error)
    
    // HANDLE ERROR 400 (SEMUA GAGAL)
    if (error.response?.status === 400) {
      await Swal.fire({
        icon: 'error',
        title: 'Upload Gagal',
        text: error.response.data.message || 'Tidak ada data yang dapat diproses.',
        footer: 'Pastikan file ZIP berisi foto dengan nama sesuai NISN yang telah terdaftar.',
        confirmButtonText: 'Tutup',
        confirmButtonColor: '#6b7280',
      })
    } else {
      await Swal.fire({
        icon: 'error',
        title: 'Kesalahan Sistem',
        text: 'Tidak dapat terhubung ke server. Silakan coba lagi.',
        confirmButtonText: 'Tutup',
        confirmButtonColor: '#6b7280',
      })
    }
  } finally {
    isLoading.value = false
    event.target.value = ''
  }
}

onMounted(() => { fetchSiswa() })
</script>

<style scoped>
.group .opacity-0 {
  transition: opacity 0.2s ease-in-out;
}
</style>