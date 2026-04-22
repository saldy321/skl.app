<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden" 
       :class="isEnvelopeOpen ? 'bg-gradient-to-br from-slate-800 via-slate-900 to-black' : 'bg-gradient-to-br from-blue-50 via-indigo-50 to-sky-50'">
    
    <!-- Background Decoration -->
    <div class="absolute inset-0 z-0 opacity-20 pointer-events-none">
      <div class="absolute top-10 left-10 w-40 h-40 bg-white/30 rounded-full blur-3xl"></div>
      <div class="absolute bottom-10 right-10 w-60 h-60 bg-indigo-200/20 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-sky-100/10 rounded-full blur-3xl"></div>
    </div>

    <!-- STATE 1: AMPLOP TERTUTUP -->
    <div v-if="!isEnvelopeOpen && !loading && siswa" 
         class="relative z-20 flex flex-col items-center justify-center animate-fade-in">
      
      <div class="text-center mb-8">
        <div class="inline-flex items-center gap-2 bg-white/80 backdrop-blur-sm px-4 py-2 rounded-full shadow-sm mb-4">
          <i class="fas fa-graduation-cap text-indigo-500"></i>
          <span class="text-xs font-bold text-slate-600 uppercase tracking-wider">{{ namaInstansi }}</span>
        </div>
        <h1 class="text-3xl md:text-4xl font-light text-slate-800 tracking-wide">
          Selamat, <span class="font-bold text-indigo-700">{{ siswa?.nama_siswa?.split(' ')[0] }}</span>
        </h1>
        <p class="text-slate-500 text-sm mt-2">Kamu telah menyelesaikan pendidikan dengan baik</p>
      </div>

      <!-- AMPLOP CONTAINER -->
      <div 
        @click="openEnvelopeWithAnimation" 
        class="cursor-pointer group relative w-80 md:w-96 h-56 perspective-1000 transition-transform duration-300 hover:scale-[1.02]"
      >
        <div class="absolute -bottom-4 left-2 right-2 h-8 bg-black/10 blur-xl rounded-full transition-all duration-300 group-hover:bg-black/20"></div>
        <div class="absolute inset-0 bg-gradient-to-br from-amber-100 to-amber-50 rounded-2xl shadow-xl border border-amber-200/50"></div>
        <div class="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-amber-200 to-amber-100 rounded-b-2xl shadow-inner"></div>
        <div class="absolute top-0 left-0 right-0 h-32 bg-gradient-to-b from-amber-200 to-amber-100 rounded-t-2xl origin-top transition-all duration-700 ease-out z-10"
             :class="{ 'rotate-x-180': isOpening }">
          <div class="absolute bottom-3 left-1/2 -translate-x-1/2 w-0 h-0 border-l-[12px] border-l-transparent border-r-[12px] border-r-transparent border-t-[12px] border-t-amber-300"></div>
        </div>
        <div class="absolute top-1/2 left-0 right-0 h-[2px] bg-amber-300/50 -translate-y-1/2 z-0"></div>
        <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-16 h-16 z-20">
          <div class="w-full h-full rounded-full bg-red-500/80 flex items-center justify-center shadow-lg border-2 border-red-600 animate-pulse">
            <i class="fas fa-lock text-white text-xl"></i>
          </div>
        </div>
        <div class="absolute bottom-8 left-0 right-0 text-center z-0">
          <span class="text-[10px] font-black text-amber-600/40 uppercase tracking-[0.3em]">CONFIDENTIAL</span>
        </div>
      </div>

      <div class="mt-10 flex flex-col items-center animate-bounce">
        <p class="text-slate-500 text-xs font-bold uppercase tracking-widest mb-2">Klik Amplop untuk Membuka</p>
        <i class="fas fa-fingerprint text-2xl text-indigo-400"></i>
      </div>
      
      <div class="absolute inset-0 pointer-events-none">
        <div v-for="i in 6" :key="i" 
             class="absolute w-1 h-1 bg-yellow-300 rounded-full animate-ping"
             :style="{ top: `${20 + i * 10}%`, left: `${10 + i * 15}%`, animationDelay: `${i * 0.2}s` }">
        </div>
      </div>
    </div>

    <!-- STATE 2: ANIMASI SURAT KELUAR -->
    <transition name="letter-emerge">
      <div v-if="isOpening && !showDocument" 
           class="relative z-30 flex flex-col items-center justify-center">
        <div class="w-72 h-96 bg-white shadow-2xl rounded-lg p-6 transform rotate-1 animate-float">
          <div class="w-full h-full border border-dashed border-indigo-200 rounded flex flex-col items-center justify-center">
            <div class="w-16 h-16 rounded-full bg-indigo-50 flex items-center justify-center mb-4">
              <i class="fas fa-envelope-open-text text-3xl text-indigo-500"></i>
            </div>
            <p class="text-indigo-700 font-serif text-lg text-center">Membuka Surat Kelulusan...</p>
            <div class="mt-4 flex gap-1">
              <span class="w-2 h-2 bg-indigo-400 rounded-full animate-bounce"></span>
              <span class="w-2 h-2 bg-indigo-400 rounded-full animate-bounce delay-100"></span>
              <span class="w-2 h-2 bg-indigo-400 rounded-full animate-bounce delay-200"></span>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- STATE 3: DOKUMEN SKL LENGKAP -->
    <transition name="document-appear">
      <div v-if="showDocument && siswa && template" 
           class="relative z-40 w-full max-w-5xl mx-auto">
        
        <!-- Tombol Admin -->
        <div class="flex justify-between items-center mb-4 print:hidden">
          <button @click="resetToEnvelope" class="bg-white/80 backdrop-blur-sm text-slate-700 px-4 py-2 rounded-full text-xs font-bold shadow-lg hover:bg-white transition-all">
            <i class="fas fa-arrow-left mr-2"></i>Kembali ke Amplop
          </button>
          <div class="flex gap-2">
            <button @click="$router.back()" class="bg-slate-600 hover:bg-slate-700 text-white px-4 py-2 rounded-full text-xs font-bold shadow-lg transition-all">
              <i class="fas fa-arrow-left mr-2"></i>Dashboard
            </button>
            <button @click="printDocument" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-full text-xs font-bold shadow-lg shadow-indigo-200 transition-all">
              <i class="fas fa-print mr-2"></i>Cetak / PDF
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

              <div class="grid grid-cols-[140px_20px_1fr] gap-y-2 mb-8 ml-4 text-sm">
                <span class="font-medium">Nama Lengkap</span><span>:</span><span class="font-bold uppercase">{{ siswa.nama_siswa }}</span>
                <span class="font-medium">Tempat, Tgl Lahir</span><span>:</span><span>{{ siswa.tempat_lahir }}, {{ formatTanggal(siswa.tanggal_lahir) }}</span>
                <span class="font-medium">NISN</span><span>:</span><span>{{ siswa.nisn }}</span>
                <span class="font-medium">Program Keahlian</span><span>:</span><span>{{ siswa.jurusan }}</span>
                <span class="font-medium">Tahun Ajaran</span><span>:</span><span>{{ siswa.tahun_lulus || '2025/2026' }}</span>
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

              <div class="flex justify-between items-start">
                <div v-if="template.pakai_foto" class="w-[3cm] h-[4cm] border-2 border-black flex items-center justify-center overflow-hidden bg-gray-50">
                  <img v-if="siswa.foto_siswa" :src="getFullImageUrl(siswa.foto_siswa, true)" class="w-full h-full object-cover" />
                  <div v-else class="text-center p-2 text-[8px] text-gray-400">Pas Foto 3x4</div>
                </div>
                <div v-else class="w-[3cm]"></div>

                <div class="w-[250px] text-left">
                  <p class="mb-1">{{ template.kota_surat || 'Bandung' }}, {{ template.tanggal_surat }}</p>
                  <p class="mb-16 font-medium">Kepala Sekolah,</p>
                  
                  <div class="relative">
                    <img v-if="template.pakai_stempel && template.stempel" :src="getFullImageUrl(template.stempel)" class="absolute -top-16 -left-8 w-24 opacity-80 mix-blend-multiply -rotate-12" />
                    <img v-if="template.pakai_ttd && template.ttd_kepsek" :src="getFullImageUrl(template.ttd_kepsek)" class="absolute -top-12 left-0 w-28" />
                  </div>
                  
                  <p class="font-bold underline uppercase text-lg">{{ template.nama_kepsek }}</p>
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
      <div class="w-16 h-16 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
      <p class="font-bold text-indigo-900 tracking-widest animate-pulse">MEMUAT DATA...</p>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const siswa = ref(null)
const template = ref(null)
const namaInstansi = ref('')

const isEnvelopeOpen = ref(false)
const isOpening = ref(false)
const showDocument = ref(false)

const getFullImageUrl = (filename, isStudentPhoto = false) => {
  if (!filename) return ''
  if (filename.startsWith('http') || filename.startsWith('data:')) return filename
  return isStudentPhoto 
    ? `http://localhost:3000/uploads/siswa/${filename}`
    : `http://localhost:3000/uploads/${filename}`
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

const openEnvelopeWithAnimation = () => {
  isOpening.value = true
  setTimeout(() => {
    isOpening.value = false
    isEnvelopeOpen.value = true
    showDocument.value = true
  }, 1500)
}

const resetToEnvelope = () => {
  showDocument.value = false
  isEnvelopeOpen.value = false
  isOpening.value = false
}

const printDocument = () => {
  window.print()
}

const loadDocument = async () => {
  try {
    const { slug, id } = route.params
    const res = await api.get(`/${slug}/admin/cetak/${id}`)
    
    if (res.data.status === 'success') {
      siswa.value = res.data.data.siswa
      template.value = res.data.data.template
    }
    
    const instansiRes = await api.get(`/check-instansi/${slug}`)
    if (instansiRes.data.success) {
      namaInstansi.value = instansiRes.data.data.nama_instansi
    }
  } catch (err) {
    console.error("Gagal memuat data:", err)
  } finally {
    loading.value = false
  }
}

onMounted(loadDocument)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Crimson+Pro:wght@400;600;700&display=swap');

.document-canvas { font-family: 'Crimson Pro', serif; }
.perspective-1000 { perspective: 1000px; }
.rotate-x-180 { transform: rotateX(180deg); }
.animate-fade-in { animation: fadeIn 0.6s ease-out; }
.animate-float { animation: float 2s ease-in-out infinite; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(1deg); }
  50% { transform: translateY(-10px) rotate(-1deg); }
}

.letter-emerge-enter-active { animation: letterEmerge 1.5s cubic-bezier(0.16, 1, 0.3, 1); }
.letter-emerge-leave-active { animation: letterEmerge 0.3s reverse; }

@keyframes letterEmerge {
  0% { opacity: 0; transform: translateY(100px) scale(0.8) rotate(-5deg); }
  30% { opacity: 1; transform: translateY(-20px) scale(1.02) rotate(2deg); }
  60% { transform: translateY(10px) scale(0.99) rotate(-1deg); }
  100% { opacity: 1; transform: translateY(0) scale(1) rotate(0); }
}

.document-appear-enter-active { animation: docAppear 0.8s ease-out; }

@keyframes docAppear {
  from { opacity: 0; transform: scale(0.95) translateY(20px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

@media print {
  @page { margin: 1.5cm; size: A4; }
  body { background: white !important; }
  .print\:hidden { display: none !important; }
  .document-canvas { box-shadow: none !important; padding: 0 !important; }
  .mix-blend-multiply { mix-blend-mode: multiply; }
  table { border-collapse: collapse; }
  th, td { border: 1px solid black !important; }
}
</style>