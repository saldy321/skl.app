<template>
  <div class="min-h-screen flex items-center justify-center bg-[#111111] p-4 font-sans relative overflow-hidden">
    
    <div class="absolute inset-0 bg-gradient-to-br from-cyan-900 via-slate-900 to-black z-0"></div>
    <div class="absolute top-0 left-0 w-full h-full bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-5 pointer-events-none"></div>

    <div v-if="statusCheck === 'loading'" class="relative z-10 text-center text-white">
      <div class="animate-spin rounded-full h-16 w-16 border-4 border-cyan-500 border-t-transparent mx-auto mb-4"></div>
      <p class="font-bold tracking-widest uppercase text-sm animate-pulse">Menghubungkan ke Server...</p>
    </div>

    <div v-else-if="statusCheck === 'valid'" class="relative z-10 w-full max-w-6xl flex flex-col lg:flex-row items-center justify-center gap-8 lg:gap-12 px-4">
      
      <!-- COUNTDOWN TIMER (KIRI) -->
      <div v-if="isClosed" class="w-full lg:w-auto flex-1 text-center lg:text-left animate-fade-in-up">
        <h2 class="text-cyan-400 text-lg md:text-xl font-bold uppercase tracking-[0.3em] mb-6 drop-shadow-lg">
          Pengumuman Dibuka Dalam
        </h2>
        
        <div class="flex justify-center lg:justify-start gap-3 md:gap-5">
          <div class="flex flex-col items-center">
            <div class="bg-[#1a1a1a] border border-white/10 text-white text-4xl md:text-5xl font-black py-4 px-3 md:px-5 rounded-2xl shadow-2xl min-w-[75px] md:min-w-[90px] text-center font-mono">
              {{ String(days).padStart(2, '0') }}
            </div>
            <span class="mt-2 text-[10px] md:text-xs font-bold uppercase tracking-widest text-gray-500">Hari</span>
          </div>

          <div class="flex flex-col items-center">
            <div class="bg-[#1a1a1a] border border-white/10 text-white text-4xl md:text-5xl font-black py-4 px-3 md:px-5 rounded-2xl shadow-2xl min-w-[75px] md:min-w-[90px] text-center font-mono">
              {{ String(hours).padStart(2, '0') }}
            </div>
            <span class="mt-2 text-[10px] md:text-xs font-bold uppercase tracking-widest text-gray-500">Jam</span>
          </div>

          <div class="flex flex-col items-center">
            <div class="bg-[#1a1a1a] border border-white/10 text-white text-4xl md:text-5xl font-black py-4 px-3 md:px-5 rounded-2xl shadow-2xl min-w-[75px] md:min-w-[90px] text-center font-mono">
              {{ String(minutes).padStart(2, '0') }}
            </div>
            <span class="mt-2 text-[10px] md:text-xs font-bold uppercase tracking-widest text-gray-500">Menit</span>
          </div>

          <div class="flex flex-col items-center">
            <div class="bg-[#1a1a1a] border border-cyan-500/30 text-cyan-400 text-4xl md:text-5xl font-black py-4 px-3 md:px-5 rounded-2xl shadow-2xl min-w-[75px] md:min-w-[90px] text-center font-mono animate-pulse">
              {{ String(seconds).padStart(2, '0') }}
            </div>
            <span class="mt-2 text-[10px] md:text-xs font-bold uppercase tracking-widest text-cyan-400">Detik</span>
          </div>
        </div>
      </div>

      <!-- FORM LOGIN (KANAN) - LEBIH LEBAR -->
      <div class="w-full lg:w-[480px] bg-[#1E1E1E] shadow-2xl p-6 md:p-8 rounded-3xl border border-white/10 transition-all duration-500 mx-auto" :class="{ 'opacity-75': isClosed }">
        
        <div class="mb-6">
          <!-- LOGO SEKOLAH (CONDITIONAL) -->
          <div v-if="tampilkanLogo && logoInstansi" class="flex justify-center mb-4">
            <div class="w-20 h-20 rounded-full overflow-hidden bg-white border-2 border-cyan-500/40 shadow-lg p-0.5">
              <img 
                :src="`http://localhost:3000/uploads/instansi/${logoInstansi}`"
                class="w-full h-full object-cover rounded-full"
                alt="Logo Sekolah"
              />
            </div>
          </div>

          <div class="flex items-center justify-center gap-3 mb-3">
            <div class="bg-cyan-500 text-black font-black italic px-3 py-1 rounded-full text-[10px] tracking-tighter shadow-lg shadow-cyan-500/20">
              sekolah 
            </div>
          </div>
          
          <h2 class="text-gray-300 font-bold text-sm tracking-wide text-center">{{ namaInstansi }}</h2>
          
          <h1 class="text-white text-xl md:text-2xl font-extrabold tracking-tight leading-tight uppercase text-center mt-2">
            PENGUMUMAN KELULUSAN SISWA
          </h1>
          <p class="text-gray-400 text-xs mt-2 text-center">
            Masukkan NISN dan Tanggal Lahir Anda
          </p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5" :class="{ 'pointer-events-none select-none': isClosed }">
         
          <div>
            <label class="block text-cyan-400 text-[11px] font-bold mb-1.5 uppercase tracking-wider">NISN</label>
            <input 
              v-model="form.nisn" 
              type="text" 
              maxlength="10"
              :disabled="isClosed"
              @input="form.nisn = form.nisn.replace(/\D/g, '')"
              class="w-full bg-[#2A2A2A] text-white p-3.5 rounded-xl outline-none border border-gray-700 focus:bg-[#333333] focus:border-cyan-500 transition-all duration-200 font-medium placeholder-gray-500 disabled:cursor-not-allowed disabled:opacity-60" 
              placeholder="Contoh: 0071234501"
              required 
            />
          </div>

          <div>
            <label class="block text-cyan-400 text-[11px] font-bold mb-1.5 uppercase tracking-wider">Tanggal Lahir</label>
            <div class="grid grid-cols-3 gap-3">
              <input 
                v-model="dob.day" 
                type="text" 
                maxlength="2"
                :disabled="isClosed"
                placeholder="Tgl"
                class="bg-[#2A2A2A] text-white p-3.5 rounded-xl outline-none border border-gray-700 focus:bg-[#333333] focus:border-cyan-500 transition-all font-medium text-center placeholder-gray-500 disabled:cursor-not-allowed disabled:opacity-60" 
                required
              />
              <input 
                v-model="dob.month" 
                type="text" 
                maxlength="2"
                :disabled="isClosed"
                placeholder="Bln"
                class="bg-[#2A2A2A] text-white p-3.5 rounded-xl outline-none border border-gray-700 focus:bg-[#333333] focus:border-cyan-500 transition-all font-medium text-center placeholder-gray-500 disabled:cursor-not-allowed disabled:opacity-60" 
                required
              />
              <input 
                v-model="dob.year" 
                type="text" 
                maxlength="4"
                :disabled="isClosed"
                placeholder="Thn"
                class="bg-[#2A2A2A] text-white p-3.5 rounded-xl outline-none border border-gray-700 focus:bg-[#333333] focus:border-cyan-500 transition-all font-medium text-center placeholder-gray-500 disabled:cursor-not-allowed disabled:opacity-60" 
                required
              />
            </div>
            <p class="text-gray-500 text-[10px] mt-1.5">Format: Tanggal Bulan Tahun</p>
          </div>

          <div v-if="errorMsg" class="bg-red-500/10 border border-red-500/30 text-red-400 text-xs font-bold p-3 rounded-lg text-center">
            <i class="fas fa-exclamation-triangle mr-1"></i> {{ errorMsg }}
          </div>

          <button 
            type="submit" 
            :disabled="loading || isClosed" 
            class="w-full bg-cyan-600 hover:bg-cyan-500 text-white font-black py-3.5 rounded-xl text-xs tracking-[0.2em] uppercase transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-cyan-900/30"
          >
            <span v-if="loading">
              <i class="fas fa-circle-notch animate-spin mr-2"></i> MEMPROSES...
            </span>
            <span v-else-if="isClosed">
              <i class="fas fa-clock mr-2"></i> MENUNGGU WAKTU...
            </span>
            <span v-else>
              <i class="fas fa-search mr-2"></i> LIHAT HASIL SELEKSI
            </span>
          </button>
        </form>

        <p class="text-gray-500 text-[10px] text-center mt-5">
          &copy; {{ new Date().getFullYear() }} SKL Digital System
        </p>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const loading = ref(false)
const errorMsg = ref('')
const statusCheck = ref('loading')
const namaInstansi = ref('')

const logoInstansi = ref('')
const tampilkanLogo = ref(true)

const isClosed = ref(false)
const countdownSeconds = ref(0)
let timerInterval = null

const form = reactive({ nisn: '' })
const dob = reactive({ day: '', month: '', year: '' })

const days = computed(() => Math.floor(countdownSeconds.value / (3600 * 24)))
const hours = computed(() => Math.floor((countdownSeconds.value % (3600 * 24)) / 3600))
const minutes = computed(() => Math.floor((countdownSeconds.value % 3600) / 60))
const seconds = computed(() => Math.floor(countdownSeconds.value % 60))

const startTimer = () => {
  if (timerInterval) clearInterval(timerInterval)
  timerInterval = setInterval(() => {
    if (countdownSeconds.value > 0) {
      countdownSeconds.value--
    } else {
      clearInterval(timerInterval)
      checkStatus()
    }
  }, 1000)
}

const checkStatus = async () => {
  try {
    const slug = route.params.slug.toLowerCase()
    
    if (route.name === 'SchoolNotFound') {
      return
    }
    
    const res = await api.get(`/check-instansi/${slug}`)
    
    if (res.data && res.data.success === true) {
      namaInstansi.value = res.data.data.nama_instansi
      isClosed.value = res.data.data.is_closed
      countdownSeconds.value = res.data.data.countdown || 0
      
      logoInstansi.value = res.data.data.logo_instansi || ''
      tampilkanLogo.value = res.data.data.tampilkan_logo !== false
      
      if (isClosed.value && countdownSeconds.value > 0) {
        startTimer()
      }
      statusCheck.value = 'valid'
    } else {
      router.replace(`/${slug}/not-found`)
    }
  } catch (err) {
    router.replace(`/${route.params.slug}/not-found`)
  }
}

onMounted(() => {
  checkStatus()
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

const handleLogin = async () => {
  if (isClosed.value) return

  if (form.nisn.length !== 10) {
    errorMsg.value = 'NISN wajib diisi 10 digit'
    return
  }
  
  if (!dob.day || !dob.month || !dob.year) {
    errorMsg.value = 'Tanggal lahir harus diisi lengkap'
    return
  }
  
  loading.value = true
  errorMsg.value = ''
  const formattedDate = `${dob.year}-${dob.month.padStart(2, '0')}-${dob.day.padStart(2, '0')}`
  
  try {
    const slug = route.params.slug
    const res = await auth.LoginSiswa(slug, form.nisn, formattedDate)
    
    if (res.success) {
      router.push({ name: 'SiswaCetak', params: { slug: slug } })
    } else {
      errorMsg.value = res.message
    }
  } catch (err) {
    errorMsg.value = 'Gagal terhubung ke pusat data seleksi.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.8s ease-out;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>