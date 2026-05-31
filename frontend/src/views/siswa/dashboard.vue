<template>
<div v-if="state !== 'skl'" class="absolute inset-0 z-0">
  <!-- Pake background-image biar auto-resize halus -->
  <div 
    v-if="backgroundImage" 
    class="absolute inset-0 bg-cover bg-center bg-no-repeat"
    :style="{ backgroundImage: `url(${backgroundImage})` }"
  ></div>
  
  <div v-else class="absolute inset-0 bg-gradient-to-br from-blue-600 via-purple-600 to-pink-500"></div>
  
  <div class="absolute inset-0 bg-black/30"></div>
</div>

<div v-if="state === 'skl'" class="absolute inset-0 z-0 bg-white"></div>
<div v-if="!loading && siswa && state !== 'skl'" class="relative z-20 flex flex-col items-center justify-center min-h-[60vh] pt-8">
  
  <!-- TEKS DATA SISWA (DI ATAS) -->
<div class="text-center mb-8 animate-fade-in">
  <!-- BARIS 1: NAMA (DIKECILKAN) -->
  <h1 class="text-2xl md:text-4xl font-sans font-black text-white drop-shadow-md tracking-wide">
    {{ siswa?.nama_siswa }}
  </h1>
  
  <!-- BARIS 2: JURUSAN -->
  <p class="text-white text-base md:text-xl font-sans font-bold mt-3 tracking-wide">
    | {{ siswa?.jurusan }} |
  </p>
  
  <!-- BARIS 3: KELAS -->
  <p class="text-white text-base font-sans font-bold mt-2 tracking-wide">
    | {{ siswa?.kelas || 'XII RPL 2' }} |
  </p>
</div>


  <!-- STATE 1: AMPLOP TERTUTUP yg costume a -->
  <div v-if="state === 'closed'" class="flex flex-col items-center justify-center animate-fade-in">
    <div @click="openEnvelope" class="cursor-pointer transition-transform hover:scale-105">
      <img src="/src/assets/close.png" alt="Amplop Tertutup" class="w-64 md:w-80 h-auto" />
    </div>
    <div class="mt-6 text-center animate-bounce">
      <p class="text-white/80 text-xs font-bold uppercase tracking-widest">Klik Amplop untuk Membuka</p>
    </div>
  </div>

  <!-- STATE 2: LOADING  -->
  <div v-if="state === 'loading'" class="flex flex-col items-center justify-center animate-fade-in">
    <div class="flex items-center gap-3">
      <div class="w-4 h-4 rounded-full bg-blue-900 animate-bounce-dot"></div>
      <div class="w-4 h-4 rounded-full bg-emerald-500 animate-bounce-dot" style="animation-delay: 0.1s"></div>
      <div class="w-4 h-4 rounded-full bg-sky-400 animate-bounce-dot" style="animation-delay: 0.2s"></div>
      <div class="w-4 h-4 rounded-full bg-amber-400 animate-bounce-dot" style="animation-delay: 0.3s"></div>
    </div>
    <p class="mt-4 text-white/80 text-xs tracking-widest animate-pulse">Mempersiapkan surat...</p>
  </div>

  <!-- STATE 3: SURAT SELAMAT ANDA LULUS -->
<div v-if="state === 'lulus' && showLulusCard" 
     class="relative z-20 flex flex-col items-center justify-center animate-fade-in">
  
  <div class="text-center mb-6 max-w-2xl px-6">
    <div class="text-white text-xs md:text-sm font-sans font-bold leading-snug tracking-wide">
      <div v-if="siswa?.status_lulus">
        {{ pesanLulusDinamis }}
      </div>
      <div v-else>
        {{ pesanTidakLulusDinamis }}
      </div>
    </div>
  </div>

  <div class="relative w-64 md:w-80">
    <img src="/src/assets/open.png" class="w-full h-auto" />
    
    <!-- LOGO SEKOLAH (HANYA SATU, TIDAK DUPLIKAT) -->
    <div v-if="logoInstansi && tampilkanLogo" class="absolute top-[63%] left-1/2 -translate-x-[40%] z-10">
      <img :src="getFullImageUrl(logoInstansi)" 
           class="w-10 h-10 object-contain rounded-full bg-white/80 p-1 shadow-sm" />
    </div>
    
    <!-- TEKS DINAMIS -->
    <div class="absolute" 
         :style="{ top: posisiTeks.top, left: '50%', transform: 'translateX(-50%)', width: posisiTeks.width }"
         :class="['text-center', posisiTeks.fontSize]">
      <h2 class="font-serif font-bold text-black leading-tight tracking-wide"
          style="font-family: 'Times New Roman', serif; white-space: pre-line;">
        {{ teksLulus }}
      </h2>
    </div>
  </div>

  <!-- Tombol Cetak SKL -->
  <div class="mt-6">
    <button @click="showSkl" 
            class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-full text-sm font-bold shadow-lg transition-all">
      <i class="fas fa-print mr-2"></i> Cetak SKL
    </button>
  </div>
</div>

</div>


<transition name="document-appear">
  <div v-if="state === 'skl' && siswa && template" 
       class="relative z-40 w-full max-w-5xl mx-auto">
        
    <!-- Toolbar -->
    <div class="flex justify-end items-center mb-4 print:hidden">
      <div class="flex gap-2">
        <button @click="printDocument" class="bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2 rounded-full text-xs font-bold shadow-lg shadow-emerald-200 transition-all">
          <i class="fas fa-print mr-2"></i>Cetak / Simpan PDF
        </button>
        <button @click="handleLogout" class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-full text-xs font-bold shadow-lg transition-all">
          <i class="fas fa-sign-out-alt mr-2"></i>Keluar
        </button>
      </div>
    </div>

    <!-- DOKUMEN SKL -->
    <div class="bg-white shadow-2xl rounded-2xl overflow-hidden print:shadow-none print:rounded-none">
      <div class="document-canvas bg-white p-6 sm:p-10 print:p-0">
        
        <div v-if="template.file_header" class="w-full mb-6">
          <img :src="getFullImageUrl(template.file_header)" class="w-full h-auto object-contain" alt="Kop Surat" />
        </div>

        <div class="px-2 sm:px-8 py-4 text-[11pt] leading-relaxed text-black font-serif">
          
          <div class="text-center mb-8">
            <h1 class="text-xl font-bold underline uppercase decoration-2 underline-offset-4">{{ template.nama_surat || 'SURAT KETERANGAN LULUS' }}</h1>
            <p class="text-md font-medium mt-2">Nomor: {{ template.no_surat }}</p>
          </div>

          <div class="mb-6 ql-editor !p-0 text-justify" v-html="processedDasarSurat"></div>

          <!-- DATA SISWA - PAKAI FLEX BUKAN GRID (FIX RAPI) -->
          <div class="data-siswa mb-8" style="margin-left: 20px;">
            <div class="flex items-start mb-1">
              <span style="width: 140px; font-weight: 500;">Nama Lengkap</span>
              <span style="width: 20px;">:</span>
              <span style="flex: 1; font-weight: 700; text-transform: uppercase;">{{ siswa.nama_siswa }}</span>
            </div>
            <div class="flex items-start mb-1">
              <span style="width: 140px; font-weight: 500;">Tempat, Tgl Lahir</span>
              <span style="width: 20px;">:</span>
              <span style="flex: 1;">{{ siswa.tempat_lahir }}, {{ formatTanggal(siswa.tanggal_lahir) }}</span>
            </div>
            <div class="flex items-start mb-1">
              <span style="width: 140px; font-weight: 500;">NISN</span>
              <span style="width: 20px;">:</span>
              <span style="flex: 1;">{{ siswa.nisn }}</span>
            </div>
            <div class="flex items-start mb-1">
              <span style="width: 140px; font-weight: 500;">Program Keahlian</span>
              <span style="width: 20px;">:</span>
              <span style="flex: 1;">{{ siswa.jurusan }}</span>
            </div>
            <div class="flex items-start mb-1">
              <span style="width: 140px; font-weight: 500;">Tahun Ajaran</span>
              <span style="width: 20px;">:</span>
              <span style="flex: 1;">{{ siswa.tahun_lulus || '2025/2026' }}</span>
            </div>
          </div>

          <div class="mb-6 ql-editor !p-0 text-justify" v-html="processedIsiSurat"></div>

          <div v-if="template.tampilkan_nilai_siswa && siswa.nilai && siswa.nilai.length > 0" class="mb-8">
            <p class="mb-3 font-medium italic">Daftar nilai yang diperoleh:</p>
            <table class="w-full border-collapse border border-black text-[10pt]">
              <thead>
                <tr class="bg-gray-50">
                  <th class="border border-black px-2 py-1 w-10 text-center">No</th>
                  <th class="border border-black px-3 py-1 text-left">Mata Pelajaran</th>
                  <th class="border border-black px-2 py-1 w-20 text-center">Nilai</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(n, index) in siswa.nilai" :key="n.id">
                  <td class="border border-black px-2 py-1 text-center">{{ index + 1 }}</td>
                  <td class="border border-black px-3 py-1">{{ n.mapel?.nama_mapel || '-' }}</td>
                  <td class="border border-black px-2 py-1 text-center font-bold">{{ n.nilai_angka }}</td>
                </tr>
              </tbody>
              <tfoot>
                <tr class="font-bold bg-gray-100">
                  <td colspan="2" class="border border-black px-3 py-1 text-right">Rata-Rata</td>
                  <td class="border border-black px-2 py-1 text-center">{{ hitungRataRata(siswa.nilai) }}</td>
                </tr>
              </tfoot>
            </table>
          </div>

          <div class="mt-6 mb-10 ql-editor !p-0 italic text-sm" v-html="template.penutup_surat"></div>

<div class="signature-area" style="display: flex; justify-content: space-between; align-items: flex-end; margin-top: 40px;">
  <div class="signature-left" style="width: 100px;">
    <div v-if="template.pakai_foto" style="width: 30mm; height: 40mm; border: 2px solid black; display: flex; align-items: center; justify-content: center; overflow: hidden; background-color: #f9fafb;">
      <img v-if="siswa.foto_siswa" :src="getFullImageUrl(siswa.foto_siswa, true)" style="width: 100%; height: 100%; object-fit: cover;" />
      <div v-else class="text-center p-2 text-[8px] text-gray-400">Pas Foto 3x4</div>
    </div>
  </div>

  <div class="signature-right" style="text-align: center; width: 250px;">
    <p class="mb-1">{{ template.kota_surat || 'Bandung' }}, {{ template.tanggal_surat }}</p>
    <p style="margin-bottom: 30px;">Kepala Sekolah,</p>
    
    <div style="position: relative; min-height: 80px;">
      <img v-if="template.pakai_stempel && template.stempel" :src="getFullImageUrl(template.stempel)" style="position: absolute; top: -50px; left: -30px; width: 80px; opacity: 0.8; mix-blend-mode: multiply;" />
      <img v-if="template.pakai_ttd && template.ttd_kepsek" :src="getFullImageUrl(template.ttd_kepsek)" style="position: absolute; top: -40px; left: 70px; width: 100px;" />
    </div>
    
    <p style="font-weight: bold; text-decoration: underline; text-transform: uppercase; margin-top: 5px;">{{ template.nama_kepsek }}</p>
    <p>NIP. {{ template.nip_kepsek }}</p>
  </div>
</div>
            

        </div>
      </div>
    </div>
  </div>
</transition>

<!-- LOADING -->
<div v-if="loading" class="fixed inset-0 z-50 bg-white/90 backdrop-blur-sm flex flex-col items-center justify-center">
  <div class="w-16 h-16 border-4 border-emerald-100 border-t-emerald-600 rounded-full animate-spin mb-4"></div>
  <p class="font-bold text-emerald-900 tracking-widest animate-pulse">MEMUAT DATA...</p>
</div>

</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'
import bgImage from '@/assets/ini.png'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const siswa = ref(null)
const template = ref(null)
const namaInstansi = ref('')
const logoInstansi = ref('') 
const tampilkanLogo = ref(false)
const loadingAnim = ref(false)
const pesanLulus = ref('')
const pesanTidakLulus = ref('')
const backgroundImage = ref('')

const state = ref('closed')
const showLulusCard = ref(false)

const pesanLulusDinamis = computed(() => {
  if (!pesanLulus.value) return ''
  return pesanLulus.value.replace(/{NAMA_INSTANSI}/g, namaInstansi.value || 'Sekolah')
})

const pesanTidakLulusDinamis = computed(() => {
  if (!pesanTidakLulus.value) return ''
  return pesanTidakLulus.value.replace(/{NAMA_INSTANSI}/g, namaInstansi.value || 'Sekolah')
})

// --- HELPERS ---
const getFullImageUrl = (filename, isStudentPhoto = false) => {
  if (!filename) return ''
  if (filename.startsWith('http') || filename.startsWith('data:')) return filename
  if (isStudentPhoto) {
    return `http://localhost:3000/uploads/siswa/${filename}`
  }
  return `http://localhost:3000/uploads/instansi/${filename}`
}

const formatTanggal = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const hitungRataRata = (nilaiList) => {
  if (!nilaiList || nilaiList.length === 0) return '0.00'
  const total = nilaiList.reduce((acc, curr) => acc + (Number(curr.nilai_angka) || 0), 0)
  return (total / nilaiList.length).toFixed(2)
}

const processPlaceholders = (html) => {
  if (!html || !siswa.value) return ''
  const statusText = siswa.value.status_lulus ? '<strong>LULUS</strong>' : '<strong>TIDAK LULUS</strong>'
  return html
    .replace(/\[STATUS\]/g, statusText)
    .replace(/\{\{\s*STATUS\s*\}\}/g, statusText)
    .replace(/\[TAHUN_AJARAN\]/g, siswa.value.tahun_lulus || '2025/2026')
    .replace(/\[NAMA SEKOLAH\]/g, auth.nama_instansi || 'Sekolah Kita')
}

const processedDasarSurat = computed(() => processPlaceholders(template.value?.dasar_surat))
const processedIsiSurat = computed(() => processPlaceholders(template.value?.isi_surat))

const openEnvelope = () => {
  state.value = 'loading'
  
  setTimeout(() => {
    state.value = 'lulus'
    showLulusCard.value = true
  }, 5000)
}

const showSkl = () => {
  state.value = 'skl'
  showLulusCard.value = false
}


const printDocument = () => {
  // Ambil elemen SKL
  const sklElement = document.querySelector('.relative.z-40.w-full.max-w-5xl.mx-auto')
  if (!sklElement) return
  
  // Clone seluruh isi SKL
  const printContents = sklElement.cloneNode(true)
  
  // Hapus tombol dan toolbar dari clone
  const toolbar = printContents.querySelector('.flex.justify-end.items-center.mb-4')
  if (toolbar) toolbar.remove()
  
  // Buat iframe tersembunyi untuk cetak
  const iframe = document.createElement('iframe')
  iframe.style.position = 'absolute'
  iframe.style.width = '0'
  iframe.style.height = '0'
  iframe.style.border = 'none'
  document.body.appendChild(iframe)
  
  const iframeDoc = iframe.contentWindow.document
  iframeDoc.open()
  iframeDoc.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <title>Cetak SKL</title>
      <style>
        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }
        body {
          font-family: 'Times New Roman', 'Crimson Pro', serif;
          background: white;
          padding: 0;
          margin: 0 auto;
          width: 100%;
          max-width: 21cm;
        }
        @page {
          size: A4;
          margin: 0.5cm;
        }
        /* SEMUA KONTEN DALAM 1 HALAMAN */
        .document-canvas {
          width: 100%;
          max-width: 100%;
          background: white;
          font-size: 11pt;
          line-height: 1.3;
          page-break-after: avoid;
          page-break-inside: avoid;
          break-inside: avoid;
          overflow-x: hidden;
        }
        /* KOP SURAT - PASTIKAN TIDAK MELEBAR */
        .document-canvas img {
          max-width: 100%;
          height: auto;
          margin-bottom: 5px;
        }
        .document-canvas h1 {
          font-size: 14pt;
          margin: 5px 0;
        }
        /* JUDUL */
        .text-center {
          text-align: center;
        }
        /* DATA SISWA */
        .data-siswa {
          margin: 10px 0 10px 0;
          width: 100%;
        }
        .data-siswa > div {
          display: flex;
          flex-wrap: nowrap;
          margin-bottom: 4px;
        }
        .data-siswa span:first-child {
          width: 140px;
          flex-shrink: 0;
        }
        .data-siswa span:nth-child(2) {
          width: 20px;
          flex-shrink: 0;
        }
        .data-siswa span:last-child {
          flex: 1;
          word-wrap: break-word;
        }
        /* TABEL NILAI */
        table {
          width: 100%;
          border-collapse: collapse;
          margin: 10px 0;
          table-layout: fixed;
        }
        th, td {
          border: 1px solid black;
          padding: 6px 8px;
          text-align: left;
          vertical-align: top;
        }
        th:first-child, td:first-child {
          width: 10%;
          text-align: center;
        }
        th:last-child, td:last-child {
          width: 15%;
          text-align: center;
        }
        /* TANDA TANGAN */
        .signature-area {
          display: flex;
          justify-content: space-between;
          margin-top: 30px;
          width: 100%;
        }
        .signature-left {
          width: 100px;
        }
        .signature-right {
          text-align: center;
          width: 250px;
        }
        .ql-editor {
          margin: 8px 0;
          text-align: justify;
        }
        /* UTILITY */
        .mb-8 { margin-bottom: 1rem; }
        .mb-6 { margin-bottom: 0.75rem; }
        .mb-1 { margin-bottom: 0.25rem; }
        .mt-2 { margin-top: 0.5rem; }
        .mt-6 { margin-top: 1.5rem; }
        .italic { font-style: italic; }
        .font-bold { font-weight: bold; }
        .underline { text-decoration: underline; }
        .uppercase { text-transform: uppercase; }
        .w-full { width: 100%; }
        .border-collapse { border-collapse: collapse; }
        .border { border: 1px solid black; }
        .bg-white { background: white; }
        /* HAPUS SHADOW & ROUNDED */
        .shadow-2xl, .shadow-lg, .shadow {
          box-shadow: none !important;
        }
        .rounded-2xl, .rounded-xl, .rounded {
          border-radius: 0 !important;
        }
        /* SEMBUNYIKAN ELEMEN YANG TIDAK PERLU */
        .print\\:hidden {
          display: none !important;
        }
        /* PASTIKAN TEKS TIDAK KELUAR DARI BATAS */
        .px-2, .px-6, .sm\\:px-8, .sm\\:px-10 {
          padding-left: 0 !important;
          padding-right: 0 !important;
        }
        .document-canvas > .px-2 {
          padding: 0 !important;
        }
      </style>
    </head>
    <body>
      ${printContents.outerHTML}
    </body>
    </html>
  `)
  iframeDoc.close()
  
  // Cetak setelah iframe siap
  setTimeout(() => {
    iframe.contentWindow.print()
    setTimeout(() => {
      document.body.removeChild(iframe)
    }, 1000)
  }, 200)
}


const handleLogout = () => {
  auth.logout()
  router.push('/login')
}

const teksLulus = computed(() => {
  if (!siswa.value) return ''
  if (siswa.value.status_lulus) {
    return 'SELAMAT ANDA LULUS'
  } else {
    return 'SELAMAT\nANDA TIDAK LULUS'
  }
})

const posisiTeks = computed(() => {
  if (!siswa.value) return { top: '30%', fontSize: 'text-lg md:text-xl', width: '65%' }
  
  if (siswa.value.status_lulus) {
    return {
      top: '30%',
      fontSize: 'text-lg md:text-xl',
      width: '65%'
    }
  } else {
    return {
      top: '32%',
      fontSize: 'text-lg md:text-xl',
      width: '60%'
    }
  }
})  

const loadBackground = async () => {
  try {
    const slug = route.params.slug
    const res = await api.get(`/${slug}/portal/setting/background`)
    if (res.data?.data?.background) {
      backgroundImage.value = `http://localhost:3000/uploads/backgrounds/${res.data.data.background}`
    }
  } catch (err) {
    console.log('Gagal load background:', err)
  }
}

const loadDocument = async () => {
  try {
    const slug = route.params.slug
    const res = await api.get(`/${slug}/portal/dashboard`)
    
    console.log('[SISWA] Response:', res.data)
    
    const dataRes = res.data.data || res.data
    siswa.value = dataRes.siswa
    template.value = dataRes.template
    
    console.log('[SISWA] Siswa:', siswa.value)
    console.log('[SISWA] Template:', template.value)
    console.log('[SISWA] Status Lulus:', siswa.value.status_lulus)
    
    const instansiRes = await api.get(`/check-instansi/${slug}`)
    if (instansiRes.data.success) {
      namaInstansi.value = instansiRes.data.data.nama_instansi
      logoInstansi.value = instansiRes.data.data.logo_instansi
      tampilkanLogo.value = instansiRes.data.data.tampilkan_logo
    }
    
    try {
      const settingRes = await api.get(`/${slug}/portal/setting/pesan`)
      if (settingRes.data?.data) {
        pesanLulus.value = settingRes.data.data.pesan_lulus || ''
        pesanTidakLulus.value = settingRes.data.data.pesan_tidak_lulus || ''
      }
    } catch (err) {
      console.log('Setting pesan belum diisi atau error:', err)
    }
    
    await loadBackground()

    console.log('[SISWA] Status Lulus:', siswa.value.status_lulus)
    console.log('[SISWA] Pesan Lulus:', pesanLulus.value)
    console.log('[SISWA] Pesan Tidak Lulus:', pesanTidakLulus.value)
    
  } catch (err) {
    console.error("Gagal memuat data:", err)
  } finally {
    loading.value = false
  }
}

onMounted(loadDocument)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&display=swap');

.document-canvas { 
  font-family: 'Times New Roman', 'Crimson Pro', serif; 
}

.animate-fade-in { animation: fadeIn 0.6s ease-out; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes bounce-dot {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.animate-bounce-dot {
  animation: bounce-dot 0.6s ease-in-out infinite;
}

.document-appear-enter-active { animation: docAppear 0.8s ease-out; }
         
@keyframes docAppear {
  from { opacity: 0; transform: scale(0.95) translateY(20px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}


@media print {
  /* SEMBUNYIKAN TOMBOL & ANIMASI */
  .print\:hidden,
  button,
  .animate-bounce,
  .animate-fade-in {
    display: none !important;
  }
  
  /* UKURAN KERTAS A4 - MARGIN KECIL */
  @page {
    size: A4;
    margin: 0.3cm;
  }
  
  /* DOKUMEN UTAMA */
  .document-canvas {
    width: 100%;
    background: white;
    font-size: 10pt;
    line-height: 1.3;
  }
  
  /* HAPUS BACKGROUND GRADIEN */
  .absolute:not(.document-canvas .absolute) {
    background: transparent !important;
  }
  
  /* KOP SURAT */
  .document-canvas img {
    max-width: 100%;
    height: auto;
    margin-bottom: 5px;
  }
  
  /* JUDUL SURAT */
  .document-canvas h1 {
    font-size: 12pt;
    margin: 3px 0;
  }
  
  .document-canvas .text-center p {
    margin: 2px 0;
  }
  
  /* DATA SISWA - PERTAHANKAN FLEX */
  .data-siswa {
    margin: 5px 0;
  }
  
  .data-siswa .flex {
    display: flex;
    margin-bottom: 2px;
  }
  
  /* QUIL EDITOR (ISI SURAT) */
  .ql-editor {
    margin: 3px 0;
    padding: 0;
  }
  
  /* TABEL NILAI */
  .document-canvas table {
    width: 100%;
    border-collapse: collapse;
    margin: 5px 0;
  }
  
  .document-canvas th,
  .document-canvas td {
    border: 1px solid black;
    padding: 3px 5px;
  }
  
  
/* TANDA TANGAN */
.signature-area {
  margin-top: 40px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end !important;  /* RATA BAWAH AGAR SEJAJAR */
}

.signature-right p {
  margin: 2px 0;
}

.signature-right p:nth-child(2) {
  margin-bottom: 10px;
}

/* STEMPEL & TTD TETAP DI POSISINYA */
.signature-right div[style*="position: relative"] {
  min-height: 50px;
}


.signature-left {
  width: auto !important;
  display: flex !important;
  justify-content: flex-start !important;
}

.signature-left > div {
  width: 30mm !important;
  height: 40mm !important;
  border: 2px solid black !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  overflow: hidden !important;
  background-color: #f9fafb !important;
  margin: 0 !important;
  padding: 0 !important;
}

.signature-left img {
  width: 100% !important;
  height: 100% !important;
  object-fit: cover !important;
  margin: 0 !important;
  padding: 0 !important;
}

  
  .document-canvas,
  .document-canvas > div {
    page-break-inside: avoid;
    break-inside: avoid;
  }
  

  .mix-blend-multiply {
    mix-blend-mode: multiply;
  }
}
</style>