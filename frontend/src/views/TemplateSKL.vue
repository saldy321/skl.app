<template>
  <div class="bg-[#f4f7f6] min-h-screen p-4 font-sans text-sm">
    <div class="max-w-7xl mx-auto bg-white shadow-sm border border-gray-200 p-6 rounded-sm">
      
      <h2 class="text-lg text-gray-600 mb-6 uppercase tracking-wide">Template SKL</h2>

    
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div>
          <label class="block text-gray-500 mb-1">Nama Surat</label>
          <input v-model="form.nama_surat" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>
        <div>
          <label class="block text-gray-500 mb-1">No Surat</label>
          <input v-model="form.no_surat" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>
        <div>
          <label class="block text-gray-500 mb-1">Tanggal Surat</label>
          <input v-model="form.tanggal_surat" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>

        <div>
          <label class="block text-blue-600 font-bold mb-1">Minimal Lulus (KKM)</label>
          <input v-model.number="form.minimal_kelulusan" type="number" step="0.01" class="w-full border border-blue-300 p-2 outline-none focus:border-blue-500 bg-blue-50">
        </div>

        <div>
          <label class="block text-gray-500 mb-1">Nama Kepala Sekolah</label>
          <input v-model="form.nama_kepsek" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>
        <div>
          <label class="block text-gray-500 mb-1">NIP Kepala Sekolah</label>
          <input v-model="form.nip_kepsek" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>
        <div>
          <label class="block text-gray-500 mb-1">Margin Top (Presisi Printer)</label>
          <input v-model.number="form.margin_top" type="number" class="w-full border border-gray-300 p-2 outline-none focus:border-blue-400">
        </div>
      </div>

      <!-- HEADER SURAT -->
      <div class="mb-8">
        <label class="block text-gray-500 mb-1">File Header (Kop Surat)</label>
        <div class="flex items-center gap-2 mb-4">
          <input type="file" @change="handleFileUpload($event, 'file_header')" class="text-xs border border-gray-300 p-1 bg-gray-50">
        </div>
        <div class="w-full border border-gray-200 py-10 flex items-center justify-center">
          <img v-if="form.file_header" :src="getFullImageUrl(form.file_header)" class="max-h-32">
          <h1 v-else class="text-4xl font-bold text-gray-800 uppercase opacity-20">KOP SURAT</h1>
        </div>
      </div>

      <!-- EDITOR TEXT -->
      <div class="space-y-6">
        <div v-for="field in ['dasar_surat', 'isi_surat']" :key="field">
          <label class="block text-gray-500 mb-1 capitalize">{{ field.replace('_', ' ') }}</label>
          <div class="border border-gray-300">
            <QuillEditor v-model:content="form[field]" content-type="html" theme="snow" :toolbar="fullToolbar" />
          </div>
        </div>

        <!-- TOGGLE TENGAH -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 py-4 border-y border-gray-100">
          <div v-for="opt in middleToggles" :key="opt.key">
            <label class="block text-gray-500 mb-1">{{ opt.label }}</label>
            <div class="flex items-center gap-2">
              <!-- KEY PROP PENTING UNTUK FORCE RE-RENDER -->
              <input type="checkbox" :key="'chk-' + opt.key" v-model="form[opt.key]">
              <span class="text-xs text-gray-400 font-bold italic">Pakai / Tidak</span>
            </div>
          </div>
        </div>

        <div>
          <label class="block text-gray-500 mb-1">Penutup Surat</label>
          <div class="border border-gray-300">
            <QuillEditor v-model:content="form.penutup_surat" content-type="html" theme="snow" :toolbar="fullToolbar" />
          </div>
        </div>
      </div>

      <!-- PENGATURAN STAMP & TTD -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-gray-500 mb-1">Width Stempel</label>
              <input v-model.number="form.width_stempel" type="number" class="w-full border border-gray-300 p-2">
            </div>
            <div>
              <label class="block text-gray-500 mb-1">Pakai Stempel</label>
              <div class="flex items-center gap-2 mt-2">
                <input type="checkbox" v-model="form.pakai_stempel">
                <span class="text-xs text-gray-400 font-bold italic">Pakai / Tidak</span>
              </div>
            </div>
          </div>
          <div>
            <label class="block text-gray-500 mb-1">Gambar Stempel</label>
            <input type="file" @change="handleFileUpload($event, 'stempel')" class="text-xs border border-gray-300 p-1 bg-gray-50 block mb-2">
            <img v-if="form.stempel" :src="getFullImageUrl(form.stempel)" :style="{ width: (form.width_stempel || 100) + 'px' }" class="border p-1 bg-white">
          </div>
        </div>

        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-gray-500 mb-1">Width TTD KEPSEK</label>
              <input v-model.number="form.width_ttd" type="number" class="w-full border border-gray-300 p-2">
            </div>
            <div>
              <label class="block text-gray-500 mb-1">Pakai TTD</label>
              <div class="flex items-center gap-2 mt-2">
                <input type="checkbox" v-model="form.pakai_ttd">
                <span class="text-xs text-gray-400 font-bold italic">Pakai / Tidak</span>
              </div>
            </div>
          </div>
          <div>
            <label class="block text-gray-500 mb-1">Gambar TTD KEPSEK</label>
            <input type="file" @change="handleFileUpload($event, 'ttd_kepsek')" class="text-xs border border-gray-300 p-1 bg-gray-50 block mb-2">
            <img v-if="form.ttd_kepsek" :src="getFullImageUrl(form.ttd_kepsek)" :style="{ width: (form.width_ttd || 100) + 'px' }" class="p-1">
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-8 pt-4 border-t border-gray-100 mb-10">
        <div v-for="foot in footerToggles" :key="foot.key">
          <label class="block text-gray-500 mb-1">{{ foot.label }}</label>
          <div class="flex items-center gap-2">
            <!-- KEY PROP PENTING -->
            <input type="checkbox" :key="'foot-' + foot.key" v-model="form[foot.key]">
            <span class="text-xs text-gray-400 font-bold italic">Pakai / Tidak</span>
          </div>
        </div>
      </div>

      <button @click="handleSave" :disabled="loading" class="bg-[#3498db] hover:bg-[#2980b9] text-white px-8 py-3 rounded text-sm transition-all shadow-md font-bold uppercase tracking-widest">
        {{ loading ? 'Sedang Menyimpan...' : 'Simpan Template' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { QuillEditor } from '@vueup/vue-quill'
import api from '../api/index.js'
import { useAuthStore } from '@/stores/auth'
import '@vueup/vue-quill/dist/vue-quill.snow.css'

const auth = useAuthStore()

// --- [UI CONFIG] ---
const fullToolbar = [
  ['bold', 'italic', 'underline', 'strike'],
  [{ 'list': 'ordered' }, { 'list': 'bullet' }],
  [{ 'align': [] }],
  ['clean']
]

const middleToggles = [
  { label: 'Nilai di Halaman Admin ?', key: 'tampilkan_nilai_admin' },
  { label: 'Nilai di Halaman Siswa ?', key: 'tampilkan_nilai_siswa' },
  { label: 'Gunakan Kelompok Mapel ?', key: 'pakai_kelompok_mapel' }
]

const footerToggles = [
  { label: 'Sertakan Foto Siswa ?', key: 'pakai_foto' },
  { label: 'Gunakan TTD QR-Code ?', key: 'pakai_ttd_qrcode' },
  { label: 'Tampilkan Nama Wali ?', key: 'pakai_nama_wali' }
]

// --- [STATE] ---
const loading = ref(false)
const form = ref({
  instansi_id: '',
  nama_surat: '',
  no_surat: '',
  tanggal_surat: '',
  nama_kepsek: '',
  nip_kepsek: '',
  margin_top: 0,
  minimal_kelulusan: 75,
  file_header: '',
  dasar_surat: '',
  isi_surat: '',
  penutup_surat: '',
  width_stempel: 100,
  pakai_stempel: false,
  stempel: '',
  width_ttd: 100,
  pakai_ttd: false,
  ttd_kepsek: '',
  tampilkan_nilai_admin: false,
  tampilkan_nilai_siswa: false,
  pakai_kelompok_mapel: false,
  pakai_foto: false, 
  pakai_ttd_qrcode: false,
  pakai_nama_wali: false
})

// --- [HELPER] ---
const getFullImageUrl = (filenameOrBase64) => {
  if (!filenameOrBase64) return ''
  if (filenameOrBase64.startsWith('http') || filenameOrBase64.startsWith('data:')) {
    return filenameOrBase64
  }
  return `http://localhost:3000/uploads/${filenameOrBase64}`
}

// --- [LOGIC] ---
const fetchTemplate = async () => {
  const userSlug = auth.slug 
  const instansiId = auth.instansi_id
  
  if (!userSlug || !instansiId) return

  try {
    const res = await api.get(`/${userSlug}/admin/template-skl`, {
      params: { instansi_id: instansiId }
    })
    
    if (res.data && res.data.data) {
      const data = res.data.data
      

      form.value.nama_surat = data.nama_surat || ''
      form.value.no_surat = data.no_surat || ''
      form.value.tanggal_surat = data.tanggal_surat || ''
      form.value.nama_kepsek = data.nama_kepsek || ''
      form.value.nip_kepsek = data.nip_kepsek || ''
      form.value.margin_top = data.margin_top || 0
      form.value.minimal_kelulusan = data.minimal_kelulusan || 75
      form.value.file_header = data.file_header || ''
      form.value.dasar_surat = data.dasar_surat || ''
      form.value.isi_surat = data.isi_surat || ''
      form.value.penutup_surat = data.penutup_surat || ''
      form.value.width_stempel = data.width_stempel || 100
      form.value.width_ttd = data.width_ttd || 150
      form.value.stempel = data.stempel || ''
      form.value.ttd_kepsek = data.ttd_kepsek || ''
      form.value.pakai_stempel = !!data.pakai_stempel
      form.value.pakai_ttd = !!data.pakai_ttd
      form.value.pakai_ttd_qrcode = !!data.pakai_ttd_qrcode
      form.value.pakai_nama_wali = !!data.pakai_nama_wali
      form.value.pakai_foto = !!data.pakai_foto       // <--- INI KUNCINYA
      form.value.pakai_kelompok_mapel = !!data.pakai_kelompok_mapel
      form.value.tampilkan_nilai_admin = !!data.tampilkan_nilai_admin
      form.value.tampilkan_nilai_siswa = !!data.tampilkan_nilai_siswa
      
    }
  } catch (err) {
    console.warn("Template data not found or error:", err.response?.data?.message)
  }
}

const handleFileUpload = (event, field) => {
  const file = event.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    form.value[field] = e.target.result 
  }
  reader.readAsDataURL(file)
}
const handleSave = async () => {
  if (!auth.isLoggedIn || !auth.slug) {
    alert("Sesi login hilang.")
    return
  }

  loading.value = true
  
  try {
    
    const payload = {

      nama_surat: form.value.nama_surat || "",
      no_surat: form.value.no_surat || "",
      tanggal_surat: form.value.tanggal_surat || "",
      nama_kepsek: form.value.nama_kepsek || "",
      nip_kepsek: form.value.nip_kepsek || "",
      file_header: form.value.file_header || "",
      stempel: form.value.stempel || "",
      ttd_kepsek: form.value.ttd_kepsek || "",
      dasar_surat: form.value.dasar_surat || "",
      isi_surat: form.value.isi_surat || "",
      penutup_surat: form.value.penutup_surat || "",


      width_stempel: parseInt(form.value.width_stempel) || 0,
      width_ttd: parseInt(form.value.width_ttd) || 0,
      margin_top: parseInt(form.value.margin_top) || 0,
      minimal_kelulusan: parseFloat(form.value.minimal_kelulusan) || 0,


      pakai_stempel: !!form.value.pakai_stempel,
      pakai_ttd: !!form.value.pakai_ttd,
      pakai_ttd_qrcode: !!form.value.pakai_ttd_qrcode,
      pakai_nama_wali: !!form.value.pakai_nama_wali,
      pakai_foto: !!form.value.pakai_foto,       
      pakai_kelompok_mapel: !!form.value.pakai_kelompok_mapel,
      tampilkan_nilai_admin: !!form.value.tampilkan_nilai_admin,
      tampilkan_nilai_siswa: !!form.value.tampilkan_nilai_siswa,
      
   
      instansi_id: auth.instansi_id
    }

    console.log("Payload Dikirim ke Backend:", payload)
    console.log("Tipe pakai_foto:", typeof payload.pakai_foto, "Nilai:", payload.pakai_foto)

 
    const res = await api.post(`/${auth.slug}/admin/template-skl`, payload)
    
    if (res.data.status === 'success') {
        alert("✅ Template Berhasil Disimpan!")
        
        await fetchTemplate() 
    } else {
        throw new Error(res.data.message || "Gagal menyimpan")
    }
    
  } catch (err) {
    console.error("Error Save Template:", err)
    const msg = err.response?.data?.message || err.message || "Terjadi kesalahan sistem"
    alert(" Gagal simpan: " + msg)
  } finally {
    loading.value = false
  }
}

onMounted(() => {

  setTimeout(fetchTemplate, 100)
})
</script>

<style>
.ql-toolbar.ql-snow { background: #f8f9fa !important; border-color: #ddd !important; }
.ql-container.ql-snow { border-color: #ddd !important; min-height: 120px; font-size: 14px; }
input:focus { transition: all 0.3s; }
</style>