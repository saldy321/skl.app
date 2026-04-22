<template>
  <div class="p-6">
    <!-- Header Actions -->
    <div class="flex flex-wrap gap-2 mb-4 items-center justify-between">
      <div class="flex gap-2">
        <input
          type="file"
          ref="fileInput"
          class="hidden"
          accept=".xlsx, .xls"
          @change="handleFileUpload"
        />
        <button
          @click="triggerFileInput"
          :disabled="uploading"
          class="bg-[#FF6B4A] hover:bg-[#e85a3d] text-white px-3 py-1.5 rounded text-xs font-bold transition-colors disabled:opacity-50 flex items-center gap-2"
        >
          <i v-if="uploading" class="fas fa-spinner animate-spin"></i>
          {{ uploading ? 'Uploading...' : 'Import Excel' }}
        </button>

        <!-- TOMBOL BARU: PROSES LULUS MASSAL -->
        <button 
          @click="openPromoteModal"
          class="bg-emerald-600 hover:bg-emerald-700 text-white px-3 py-1.5 rounded text-xs font-bold transition-colors flex items-center gap-2 shadow-sm border border-emerald-500"
        >
          <i class="fas fa-graduation-cap"></i> Proses Lulus Massal
        </button>

        <button class="bg-slate-700 hover:bg-slate-800 text-white px-3 py-1.5 rounded text-xs font-bold transition-colors">
          <i class="fas fa-print mr-1"></i> Cetak SKL Massal
        </button>
      </div>

      <button
        @click="tambahSiswa"
        class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded text-[11px] font-bold transition-colors flex items-center gap-2 shadow-sm"
      >
        <i class="fas fa-plus"></i> Tambah Siswa Manual
      </button>
    </div>

    <!-- Table Data Siswa -->
    <div class="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
      <div class="p-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
        <h2 class="text-lg font-bold text-slate-800">Data Siswa</h2>
        <span class="text-xs text-slate-500">Total: {{ dataSiswa.length }} Siswa</span>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-[13px]">
          <thead>
            <tr class="border-b border-slate-200 bg-slate-50 text-slate-600 uppercase tracking-wider text-[11px]">
              <th class="p-3 font-bold text-center w-10">#</th>
              <th class="p-3 font-bold">NISN</th>
              <th class="p-3 font-bold">Nama Siswa</th>
              <th class="p-3 font-bold text-center">
                Kelas<span v-if="showJurusanField">/Jurusan</span>
              </th>
              <th class="p-3 font-bold text-center">Angkatan</th>
              <th class="p-3 font-bold text-center">Status</th>
              <th class="p-3 font-bold text-center w-32">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="7" class="p-10 text-center">
                <div class="flex flex-col items-center justify-center gap-3">
                  <i class="fas fa-circle-notch animate-spin text-2xl text-blue-500"></i>
                  <p class="text-slate-500 text-sm">Memuat data siswa...</p>
                </div>
              </td>
            </tr>

            <tr v-else-if="dataSiswa.length === 0">
              <td colspan="7" class="p-10 text-center text-slate-500">
                Belum ada data siswa. Silakan import atau tambah manual.
              </td>
            </tr>

            <tr v-for="(item, index) in dataSiswa" :key="item.id" class="hover:bg-slate-50 transition-colors group">
              <td class="p-3 text-center text-slate-400">{{ index + 1 }}</td>
              <td class="p-3 font-mono text-slate-600">{{ item.nisn }}</td>
              <td class="p-3">
                <div class="font-bold text-slate-800">{{ item.nama_siswa }}</div>
                <div class="text-[10px] text-slate-400">{{ item.tempat_lahir }}, {{ item.tanggal_lahir }}</div>
              </td>
              <td class="p-3 text-center">
                <div class="text-slate-700">{{ item.kelas }}</div>
                <div v-if="showJurusanField" class="text-[10px] text-slate-500">{{ item.jurusan || '-' }}</div>
              </td>
              
              <td class="p-3 text-center">
                <span class="inline-block px-2 py-1 bg-indigo-50 text-indigo-700 rounded text-[11px] font-bold border border-indigo-100">
                  {{ item.tahun_lulus || '-' }}
                </span>
              </td>

              <td class="p-3 text-center">
                <span :class="item.status_lulus ? 'bg-emerald-100 text-emerald-700 border-emerald-200' : 'bg-red-100 text-red-700 border-red-200'" 
                  class="px-2 py-1 rounded text-[10px] font-bold uppercase border">
                  {{ item.status_lulus ? 'Lulus' : 'Belum Lulus' }}
                </span>
              </td>
              <td class="p-3 text-center">
                <div class="flex justify-center gap-1 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                  <router-link 
                    :to="`/${$route.params.slug}/admin/cetak/${item.id}`" 
                    target="_blank"
                    class="w-7 h-7 bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white rounded flex items-center justify-center transition-all" title="Cetak SKL">
                    <i class="fas fa-print text-[11px]"></i>
                  </router-link>

                  <button @click="editSiswa(item)" class="w-7 h-7 bg-orange-50 text-orange-600 hover:bg-orange-600 hover:text-white rounded flex items-center justify-center transition-all" title="Edit">
                    <i class="fas fa-pencil-alt text-[11px]"></i>
                  </button>
                  
                  <button @click="deleteSiswa(item.id)" class="w-7 h-7 bg-red-50 text-red-600 hover:bg-red-600 hover:text-white rounded flex items-center justify-center transition-all" title="Hapus">
                    <i class="fas fa-trash-alt text-[11px]"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL FORM TAMBAH/EDIT SISWA -->
    <div v-if="isEditing || isAdding" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="resetForm"></div>
      
      <div class="relative bg-white w-full max-w-lg rounded-xl shadow-2xl overflow-hidden transform transition-all scale-100">
        <div class="bg-slate-50 p-4 border-b border-slate-200 flex justify-between items-center">
          <h3 class="font-bold text-slate-800 text-lg">{{ isAdding ? 'Tambah Siswa Baru' : 'Edit Data Siswa' }}</h3>
          <button @click="resetForm" class="text-slate-400 hover:text-red-500 transition-colors">
            <i class="fas fa-times text-lg"></i>
          </button>
        </div>

        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
          <!-- Baris 1: Nama & NISN -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">NISN <span class="text-red-500">*</span></label>
              <input v-model="form.nisn" type="text" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Nama Lengkap <span class="text-red-500">*</span></label>
              <input v-model="form.nama_siswa" type="text" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
          </div>

          <!-- Baris 2: TTL -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Tempat Lahir</label>
              <input v-model="form.tempat_lahir" type="text" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Tanggal Lahir</label>
              <input v-model="form.tanggal_lahir" type="date" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
          </div>

          <!-- Baris 3: Kelas, Jurusan, Angkatan -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Kelas</label>
              <input v-model="form.kelas" type="text" placeholder="XII" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
            
            <!-- FIELD JURUSAN - CONDITIONAL -->
            <div v-if="showJurusanField">
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Jurusan</label>
              <input v-model="form.jurusan" type="text" placeholder="RPL" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
            
            <div :class="showJurusanField ? '' : 'md:col-span-2'">
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Angkatan/Lulus</label>
              <select v-model="form.tahun_lulus" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white">
                <option value="" disabled>Pilih Tahun</option>
                <option value="2024">2024</option>
                <option value="2025">2025</option>
                <option value="2026">2026</option>
                <option value="2027">2027</option>
                <option value="2028">2028</option>
              </select>
            </div>
          </div>

          <!-- Baris 4: Status & Gender -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Jenis Kelamin</label>
              <select v-model="form.jenis_kelamin" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white">
                <option value="L">Laki-laki</option>
                <option value="P">Perempuan</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Status Kelulusan</label>
              <select v-model="form.status_lulus" class="w-full border border-slate-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white">
                <option :value="true">LULUS</option>
                <option :value="false">BELUM LULUS / AKTIF</option>
              </select>
            </div>
          </div>
        </div>

        <div class="p-4 bg-slate-50 border-t border-slate-200 flex gap-3 justify-end">
          <button @click="resetForm" class="px-4 py-2 text-sm font-bold text-slate-600 hover:bg-slate-200 rounded-lg transition-colors">
            Batal
          </button>
          <button @click="saveSiswa" class="px-6 py-2 text-sm font-bold text-white bg-blue-600 hover:bg-blue-700 rounded-lg shadow-md transition-all flex items-center gap-2">
            <i class="fas fa-save"></i> Simpan Data
          </button>
        </div>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- MODAL BARU: PROSES KELULUSAN MASSAL (SMART) -->
    <!-- ========================================== -->
    <div v-if="showPromoteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-900/80 backdrop-blur-sm" @click="closePromoteModal"></div>
      
      <div class="relative bg-white w-full max-w-lg rounded-xl shadow-2xl overflow-hidden p-0 transform transition-all">
        
        <!-- Header Modal -->
        <div class="bg-emerald-600 p-4 flex justify-between items-center text-white">
          <h3 class="font-bold text-lg flex items-center gap-2"><i class="fas fa-graduation-cap"></i> Proses Kelulusan Massal</h3>
          <button @click="closePromoteModal" class="text-white/80 hover:text-white"><i class="fas fa-times"></i></button>
        </div>

        <div class="p-6">
          
          <!-- STEP 1: INPUT TAHUN & CEK -->
          <div v-if="promotionStep === 1" class="space-y-4">
            <p class="text-sm text-slate-600">Masukkan tahun angkatan lulus untuk memproses siswa yang belum lulus.</p>
            
            <div>
              <label class="block text-xs font-bold text-slate-500 uppercase mb-1">Tahun Angkatan Lulus</label>
              <select v-model="promoteForm.tahun_lulus" class="w-full border border-slate-300 rounded-lg p-3 text-sm focus:ring-2 focus:ring-emerald-500 outline-none">
                <option value="2024">2024</option>
                <option value="2025">2025</option>
                <option value="2026">2026</option>
                <option value="2027">2027</option>
              </select>
            </div>

            <button @click="checkEligibility" :disabled="checking" class="w-full bg-slate-800 hover:bg-slate-900 text-white py-3 rounded-lg font-bold flex justify-center items-center gap-2">
              <i v-if="checking" class="fas fa-spinner animate-spin"></i>
              {{ checking ? 'Menganalisa Data...' : 'Cek Kelayakan Dulu' }}
            </button>
          </div>

          <!-- STEP 2: HASIL ANALISA & KONFIRMASI -->
          <div v-if="promotionStep === 2" class="space-y-4">
            <div class="bg-slate-50 p-4 rounded-lg border border-slate-200 mb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="text-xs font-bold text-slate-500 uppercase">Hasil Pengecekan (KKM: {{ analysisResult.kkm_used }})</span>
                <span class="text-xs text-slate-400">Total Aktif: {{ analysisResult.total_cek }}</span>
              </div>
              
              <div class="grid grid-cols-2 gap-3 mt-2">
                <div class="bg-emerald-50 p-3 rounded border border-emerald-100 text-center">
                  <div class="text-2xl font-bold text-emerald-600">{{ analysisResult.layak_count }}</div>
                  <div class="text-[10px] uppercase font-bold text-emerald-700">Layak Lulus</div>
                </div>
                <div class="bg-red-50 p-3 rounded border border-red-100 text-center">
                  <div class="text-2xl font-bold text-red-600">{{ analysisResult.invalid_count }}</div>
                  <div class="text-[10px] uppercase font-bold text-red-700">Belum Memenuhi</div>
                </div>
              </div>
            </div>

            <!-- Opsi Tindakan -->
            <div v-if="analysisResult.invalid_count > 0" class="space-y-3">
              <p class="text-xs font-bold text-slate-600">Pilih Tindakan:</p>
              
              <label class="flex items-start gap-3 p-3 border rounded cursor-pointer hover:bg-emerald-50 transition-colors" :class="{ 'border-emerald-500 bg-emerald-50 ring-1 ring-emerald-500': promoteMode === 'eligible_only' }">
                <input type="radio" v-model="promoteMode" value="eligible_only" class="mt-1 text-emerald-600 focus:ring-emerald-500">
                <div>
                  <span class="text-sm font-bold text-slate-800">Luluskan Hanya Yang Layak (Rekomendasi)</span>
                  <p class="text-[10px] text-slate-500">Sistem akan melewatkan {{ analysisResult.invalid_count }} siswa bermasalah. Edit manual nanti jika perlu.</p>
                </div>
              </label>

              <label class="flex items-start gap-3 p-3 border rounded cursor-pointer hover:bg-red-50 transition-colors" :class="{ 'border-red-500 bg-red-50 ring-1 ring-red-500': promoteMode === 'force_all' }">
                <input type="radio" v-model="promoteMode" value="force_all" class="mt-1 text-red-600 focus:ring-red-500">
                <div>
                  <span class="text-sm font-bold text-slate-800">Paksa Luluskan Semua</span>
                  <p class="text-[10px] text-slate-500">Abaikan syarat nilai. Wajib mengisi alasan administrasi.</p>
                  <textarea v-if="promoteMode === 'force_all'" v-model="forceReason" placeholder="Contoh: Keputusan Rapat Dewan Guru No. XX..." class="w-full mt-2 text-xs p-2 border rounded focus:ring-1 focus:ring-red-500 bg-white"></textarea>
                </div>
              </label>
            </div>

            <div class="flex gap-3 pt-2">
              <button @click="promotionStep = 1" class="flex-1 px-4 py-2.5 text-sm font-bold text-slate-600 hover:bg-slate-100 rounded-lg">Kembali</button>
              <button @click="executePromotion" :disabled="processing" class="flex-1 px-4 py-2.5 text-sm font-bold text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-md flex justify-center items-center gap-2">
                <i v-if="processing" class="fas fa-circle-notch animate-spin"></i>
                {{ processing ? 'Memproses...' : 'Konfirmasi & Proses' }}
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import Swal from 'sweetalert2'
import { TINGKAT_SEKOLAH } from '@/utils/constants'

const route = useRoute()
const auth = useAuthStore()
const fileInput = ref(null) 
const uploading = ref(false)
const dataSiswa = ref([])
const loading = ref(true)

const isEditing = ref(false)
const isAdding = ref(false)
const selectedId = ref(null)

// State Fitur Massal
const showPromoteModal = ref(false)
const promotionStep = ref(1) // 1: Input, 2: Result
const checking = ref(false)
const processing = ref(false)
const promoteMode = ref('eligible_only') // default
const forceReason = ref('')
const analysisResult = ref({})
const promoteForm = ref({
  tahun_lulus: new Date().getFullYear().toString()
})

// ✅ COMPUTED PROPERTY UNTUK CONDITIONAL UI JURUSAN
const showJurusanField = computed(() => {
    const tingkat = auth.tingkat?.toUpperCase() || ''
    return [TINGKAT_SEKOLAH.SMK, TINGKAT_SEKOLAH.SMA, TINGKAT_SEKOLAH.MA].includes(tingkat)
})

// Default Form State Siswa
const defaultForm = {
  nisn: '', nama_siswa: '', tempat_lahir: '', tanggal_lahir: '',
  kelas: '', jurusan: '', status_lulus: true, jenis_kelamin: 'L', tahun_lulus: ''
}
const form = ref({ ...defaultForm })

// --- LOGIC SISWA BIASA ---
const resetForm = () => { isEditing.value = false; isAdding.value = false; selectedId.value = null; form.value = { ...defaultForm } }
const tambahSiswa = () => { resetForm(); isAdding.value = true }

const fetchSiswa = async () => {
  loading.value = true
  try {
    const res = await api.get(`/${route.params.slug}/admin/siswa`, {
      params: {
        tahun_lulus: auth.selectedYear
      }
    })
    dataSiswa.value = res.data.data
  } catch (err) {
    console.error("Gagal tarik data:", err)
  } finally {
    loading.value = false
  }
}

watch(() => auth.selectedYear, () => {
    fetchSiswa();
})

const editSiswa = (item) => { selectedId.value = item.id; form.value = { ...item, tahun_lulus: item.tahun_lulus || '' }; isEditing.value = true }
const saveSiswa = async () => {
  if (!form.value.nisn || !form.value.nama_siswa) return Swal.fire('Warning', 'NISN/Nama Wajib', 'warning')
  try {
    const fd = new FormData()
    Object.keys(form.value).forEach(key => fd.append(key, form.value[key]))
    if (isAdding.value) await api.post(`/${route.params.slug}/admin/siswa`, fd)
    else await api.put(`/${route.params.slug}/admin/siswa/${selectedId.value}`, fd)
    Swal.fire('Success', 'Disimpan!', 'success'); fetchSiswa(); resetForm()
  } catch (err) { Swal.fire('Error', err.response?.data?.message, 'error') }
}
const deleteSiswa = async (id) => {
  const res = await Swal.fire({ title: 'Hapus?', icon: 'warning', showCancelButton: true, confirmButtonColor: '#E74C3C' })
  if (res.isConfirmed) { await api.delete(`/${route.params.slug}/admin/siswa/${id}`); Swal.fire('Deleted!', '', 'success'); fetchSiswa() }
}
const triggerFileInput = () => fileInput.value.click()
const handleFileUpload = async (event) => {
  const file = event.target.files[0]; if (!file) return
  const formData = new FormData(); formData.append('file', file)
  uploading.value = true
  try { await api.post(`/${route.params.slug}/admin/siswa/import`, formData); Swal.fire('Mantap!', 'Import Berhasil', 'success'); fetchSiswa() }
  catch (err) { Swal.fire('Error', 'Import Gagal', 'error') }
  finally { uploading.value = false; event.target.value = '' }
}

// --- LOGIC FITUR MASSAL BARU ---

const openPromoteModal = () => {
  showPromoteModal.value = true
  promotionStep.value = 1
  analysisResult.value = {}
  promoteMode.value = 'eligible_only'
  forceReason.value = ''
}

const closePromoteModal = () => {
  showPromoteModal.value = false
}

const checkEligibility = async () => {
  checking.value = true
  try {
    const slug = route.params.slug
    const res = await api.post(`/${slug}/admin/check-graduation`, {
      tahun_lulus: promoteForm.value.tahun_lulus
    })
    
    analysisResult.value = res.data
    
    if (res.data.total_cek === 0) {
      Swal.fire('Info', 'Tidak ada siswa aktif untuk diproses.', 'info')
      closePromoteModal()
      return
    }

    promotionStep.value = 2
  } catch (err) {
    Swal.fire('Error', err.response?.data?.message || 'Gagal cek kelayakan', 'error')
  } finally {
    checking.value = false
  }
}

const executePromotion = async () => {
  if (promoteMode.value === 'force_all' && !forceReason.value.trim()) {
    Swal.fire('Peringatan', 'Alasan administrasi wajib diisi jika memaksa kelulusan!', 'warning')
    return
  }

  const result = await Swal.fire({
    title: 'Konfirmasi Akhir?',
    text: promoteMode.value === 'eligible_only' 
      ? `Hanya ${analysisResult.value.layak_count} siswa layak yang akan diproses.`
      : `SEMUA ${analysisResult.value.total_cek} siswa akan diproses (Termasuk yang bermasalah).`,
    icon: 'question',
    showCancelButton: true,
    confirmButtonColor: '#10B981',
    confirmButtonText: 'Ya, Proses Sekarang!'
  })

  if (result.isConfirmed) {
    processing.value = true
    try {
      const slug = route.params.slug
      
      const payload = {
        tahun_lulus: promoteForm.value.tahun_lulus,
        mode: promoteMode.value,
        reason: forceReason.value
      }

      if (promoteMode.value === 'eligible_only') {
        payload.student_ids = analysisResult.value.layak_ids
      }

      await api.post(`/${slug}/admin/execute-promotion`, payload)
      
      Swal.fire('Berhasil!', 'Proses kelulusan massal selesai.', 'success')
      closePromoteModal()
      fetchSiswa()
    } catch (err) {
      Swal.fire('Gagal!', err.response?.data?.message, 'error')
    } finally {
      processing.value = false
    }
  }
}

onMounted(fetchSiswa)
</script>

<style scoped>
/* Your custom styles here if needed */
</style>