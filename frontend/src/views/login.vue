<template>
  <div class="min-h-screen bg-[#F3F4F6] flex items-center justify-center p-4 md:p-10 font-sans text-slate-900">
    <div class="bg-white w-full max-w-7xl rounded-[2rem] shadow-2xl flex overflow-hidden min-h-[85vh]">
      
      <!-- Left Side: Branding -->
      <div class="hidden md:flex md:w-1/2 bg-gradient-to-br from-[#E0E7FF] via-[#F8FAFC] to-[#E0E7FF] p-12 items-center justify-center relative overflow-hidden">
      

        <div class="bg-white/40 backdrop-blur-3xl border border-white/50 p-12 rounded-[3rem] shadow-xl w-full max-w-sm text-center relative z-10">
          <div class="w-24 h-24 bg-white rounded-3xl mx-auto mb-8 flex items-center justify-center shadow-lg overflow-hidden p-2">
            <img src="@/assets/sko.png" alt="Logo Sekolah" class="w-full h-full object-contain">
          </div>
          <h3 class="text-3xl font-semibold text-slate-900">Form Login Admin</h3>
          <p class="text-slate-500 mt-4 text-sm font-bold uppercase tracking-widest">Surat keterangan lulis</p>
        </div>
      </div>

      <!-- Right Side: Form -->
      <div class="w-full md:w-1/2 bg-white p-8 md:p-24 flex flex-col justify-center">
        <div class="mb-12 text-center md:text-left">
          <h2 class="text-3xl font-semibold tracking-tight mb-4 text-slate-900 transition-all duration-500">
            {{ showOTPForm ? 'Verifikasi OTP' : 'Welcome back' }}
          </h2>
          <p class="text-slate-400 font-normal text-lg text-pretty">
            {{ showOTPForm ? 'Masukkan kode 6 digit yang dikirim ke email Anda' : 'silahkan untuk mmengisi' }}
          </p>
        </div>

        <!-- Error Message Display -->
        <transition name="fade">
          <div v-if="errorMsg" class="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700 text-sm font-bold rounded-r-xl">
            {{ errorMsg }}
          </div>
        </transition>

        <!-- LOGIN FORM -->
        <form v-if="!showOTPForm" @submit.prevent="handleLogin" class="space-y-7 transition-all duration-500">
          <div class="space-y-3">
            <label class="block text-[10px] font-black uppercase tracking-[0.2em] ml-1 text-slate-700">Email Address</label>
            <input v-model="email" type="email" required class="w-full px-6 py-5 rounded-2xl border-2 border-slate-100 bg-slate-50 focus:outline-none focus:ring-4 focus:ring-slate-900/5 focus:border-slate-900 focus:bg-white transition-all duration-300 text-sm font-bold text-slate-900" placeholder="super@gmail.com">
          </div>

          <div class="space-y-3">
            <div class="flex justify-between items-center ml-1">
              <label class="block text-[10px] font-black uppercase tracking-[0.2em] text-slate-700">Password</label>
              <a href="#" class="text-[10px] font-black text-slate-400 hover:text-slate-900 transition-colors uppercase">Lupa?</a>
            </div>
            <input v-model="password" type="password" required class="w-full px-6 py-5 rounded-2xl border-2 border-slate-100 bg-slate-50 focus:outline-none focus:ring-4 focus:ring-slate-900/5 focus:border-slate-900 focus:bg-white transition-all duration-300 text-sm font-bold text-slate-900" placeholder="********">
          </div>

          <!-- Slider Captcha -->
          <div class="space-y-3">
            <label class="block text-[10px] font-black uppercase tracking-[0.2em] ml-1 text-slate-700">Verification</label>
            <div ref="sliderContainer" class="relative w-full h-16 bg-slate-100 rounded-2xl border-2 border-slate-100 overflow-hidden flex items-center justify-center select-none" :class="{ 'bg-emerald-50 border-emerald-100': isVerified }">
              <span class="text-xs font-bold transition-opacity duration-300 pointer-events-none" :class="isVerified ? 'text-emerald-600' : 'text-slate-400'">
                {{ isVerified ? 'VERIFIKASI BERHASIL' : 'Geser slider ke target' }}
              </span>
              <div class="absolute right-2 w-12 h-12 border-2 border-dashed border-slate-300 rounded-xl" :class="{ 'hidden': isVerified }"></div>
              <div @mousedown="startDrag" @touchstart="startDrag" :style="{ transform: `translateX(${sliderPos}px)` }" class="absolute left-2 w-12 h-12 bg-white rounded-xl shadow-md flex items-center justify-center border border-slate-200 touch-none cursor-default" :class="{ 'bg-emerald-500 border-emerald-400 shadow-emerald-200': isVerified, 'transition-transform duration-300 ease-out': !isDragging }">
                <div v-if="!isVerified" class="grid grid-cols-2 gap-1 pointer-events-none">
                  <div v-for="i in 4" :key="i" class="w-1.5 h-1.5 bg-slate-300 rounded-full"></div>
                </div>
                <svg v-else class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M5 13l4 4L19 7"></path>
                </svg>
              </div>
            </div>
          </div>

          <button type="submit" :disabled="loading || !isVerified" class="w-full bg-slate-900 hover:bg-black text-white font-black py-5 rounded-2xl transition-all duration-300 shadow-2xl shadow-slate-900/20 active:scale-[0.98] uppercase text-xs tracking-[0.3em] disabled:bg-slate-300">
            {{ loading ? 'MEMPROSES...' : 'MASUK SEKARANG' }}
          </button>
        </form>

        <!-- OTP FORM -->
        <form v-else @submit.prevent="handleVerifyOTP" class="space-y-7 transition-all duration-500">
          <div class="space-y-6">
            <div class="space-y-4">
              <label class="block text-[10px] font-black uppercase tracking-[0.2em] text-slate-700 text-center">Kode OTP (6 Digit)</label>
              <input 
                v-model="otp" 
                type="text" 
                maxlength="6" 
                required 
                autofocus
                class="w-full px-6 py-8 rounded-3xl border-2 border-slate-100 bg-slate-50 focus:outline-none focus:border-slate-900 text-center text-4xl font-black tracking-[0.5em] text-slate-900 transition-all shadow-inner"
                placeholder="******"
              >
            </div>

            <div class="text-center space-y-2">
              <p class="text-sm font-bold uppercase tracking-widest transition-colors duration-300" :class="timeLeft < 30 ? 'text-red-500 animate-pulse' : 'text-slate-400'">
                <span v-if="timerActive">KODE BERAKHIR DALAM: {{ formattedTime }}</span>
                <span v-else class="text-red-600">KODE TELAH KEDALUWARSA</span>
              </p>
              
              <button v-if="!timerActive" @click="resendOTP" type="button" class="text-xs font-black text-blue-600 hover:text-blue-800 uppercase tracking-widest underline decoration-2 underline-offset-4">
                Kirim Ulang Kode OTP
              </button>
            </div>

            <p class="text-center text-[10px] text-slate-400 font-bold uppercase tracking-widest cursor-pointer hover:text-slate-900 transition-colors" @click="resetForm">
              Salah Email? Kembali ke Login
            </p>
          </div>

          <button type="submit" :disabled="loading || otp.length < 6 || !timerActive" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-black py-5 rounded-2xl transition-all duration-300 shadow-xl shadow-emerald-500/20 active:scale-[0.98] uppercase text-xs tracking-[0.3em] disabled:bg-slate-300 disabled:shadow-none">
            {{ loading ? 'KONFIRMASI OTP' : 'VERIFIKASI SEKARANG' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ErrorInstansi from '@/components/ErrorSekolah.vue'

const auth = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const otp = ref('')
const errorMsg = ref('')
const loading = ref(false)
const showOTPForm = ref(false)

// Timer Logic
const timeLeft = ref(300)
const timerActive = ref(true)
let timerInterval = null

const formattedTime = computed(() => {
  const minutes = Math.floor(timeLeft.value / 60)
  const seconds = timeLeft.value % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
})

const startTimer = () => {
  if (timerInterval) clearInterval(timerInterval)
  timeLeft.value = 300
  timerActive.value = true
  
  timerInterval = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      timerActive.value = false
      clearInterval(timerInterval)
    }
  }, 1000)
}

// Reset timer untuk kirim ulang OTP
const resetTimer = () => {
  if (timerInterval) clearInterval(timerInterval)
  timeLeft.value = 300
  timerActive.value = true
  
  timerInterval = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      timerActive.value = false
      clearInterval(timerInterval)
    }
  }, 1000)
}

watch(showOTPForm, (newVal) => {
  if (newVal) startTimer()
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

// Slider Captcha Logic
const isVerified = ref(false)
const sliderPos = ref(0)
const sliderContainer = ref(null)
let isDragging = false
let startX = 0

const startDrag = (e) => {
  if (isVerified.value) return
  isDragging = true
  startX = (e.type === 'mousedown' ? e.clientX : e.touches[0].clientX) - sliderPos.value
  window.addEventListener('mousemove', onDrag)
  window.addEventListener('touchmove', onDrag)
  window.addEventListener('mouseup', stopDrag)
  window.addEventListener('touchend', stopDrag)
}

const onDrag = (e) => {
  if (!isDragging) return
  const clientX = e.type === 'mousemove' ? e.clientX : e.touches[0].clientX
  const maxMove = sliderContainer.value.offsetWidth - 48 - 16
  let move = clientX - startX
  if (move < 0) move = 0
  if (move > maxMove) move = maxMove
  sliderPos.value = move
  
  if (move >= maxMove - 5) {
    isVerified.value = true
    sliderPos.value = maxMove
    stopDrag()
  }
}

const stopDrag = () => {
  isDragging = false
  if (!isVerified.value) sliderPos.value = 0
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('touchmove', onDrag)
  window.removeEventListener('mouseup', stopDrag)
  window.removeEventListener('touchend', stopDrag)
}

// --- AUTH ACTIONS ---

const handleLogin = async () => {
  if (!isVerified.value) return
  
  errorMsg.value = ''
  loading.value = true
  
  try {
    const res = await auth.login(email.value, password.value)
    
    // Cek apakah login gagal
    if (!res.success) {
      // Tampilkan pesan error dari backend
      errorMsg.value = res.message || 'Email atau password salah!'
      resetCaptcha()
      loading.value = false
      return
    }

    if (res.status === 'otp_required') {
      showOTPForm.value = true
      loading.value = false
      return
    } 
    
    if (res.status === 'success') {
      triggerRedirect(res.role, res.slug)
      loading.value = false
      return
    }

  } catch (err) {
    // Tangkap error dari axios
    console.error("Login error:", err)
    errorMsg.value = err.response?.data?.message || 'Terjadi kesalahan, coba lagi nanti'
    resetCaptcha()
    loading.value = false
  }
}

const handleVerifyOTP = async () => {
  if (!timerActive.value) {
    errorMsg.value = "Kode OTP sudah kedaluwarsa!";
    return;
  }
  
  loading.value = true;
  errorMsg.value = ''; 
  
  try {
    const res = await auth.verifyOTP(email.value, otp.value);

    if (!res.success) {
      errorMsg.value = res.message || 'Kode OTP salah!';
      loading.value = false;
      return;
    }

    triggerRedirect(res.role, res.slug)
    loading.value = false

  } catch (err) {
    errorMsg.value = err.response?.data?.message || 'Kode OTP salah!';
    loading.value = false;
  }
}

// Fungsi Redirect
const triggerRedirect = (role, slug) => {
  console.log('[Redirect] Role:', role, 'Slug:', slug)
  
  if (role === 'super_admin') {
    router.replace({ name: 'super-dashboard' })
  } 
  else if (role === 'admin') {
    if (slug && slug !== 'null' && slug !== '') {
      router.replace(`/${slug}/dashboard`)
    } else {
      console.error('[Redirect] Slug tidak ditemukan untuk admin')
      router.replace('/login')
    }
  } 
  else if (role === 'siswa') {
    if (slug && slug !== 'null' && slug !== '') {
      router.replace(`/${slug}/portal/skl`)
    } else {
      router.replace('/login')
    }
  } 
  else {
    router.replace('/login')
  }

  setTimeout(() => {
    if (window.location.pathname.includes('/login')) {
      if (role === 'super_admin') {
        window.location.href = '/super-admin'
      } 
      else if (role === 'admin' && slug) {
        window.location.href = `/${slug}/dashboard`
      } 
      else if (role === 'siswa' && slug) {
        window.location.href = `/${slug}/portal/skl`
      } 
      else {
        window.location.href = '/login'
      }
    }
  }, 1000)
}

// ✅ PERBAIKAN: Kirim ulang OTP dengan reset timer
const resendOTP = async () => {
  otp.value = ''
  resetTimer()  // Reset timer dulu
  
  loading.value = true
  errorMsg.value = ''
  
  try {
    const res = await auth.login(email.value, password.value)
    
    if (res.status === 'otp_required') {
      loading.value = false
      alert("Kode OTP baru telah dikirim ke email Anda")
    } else {
      errorMsg.value = res.message || 'Gagal mengirim ulang OTP'
      loading.value = false
    }
  } catch (err) {
    errorMsg.value = err.response?.data?.message || 'Gagal mengirim ulang OTP'
    loading.value = false
  }
}

const resetCaptcha = () => {
  isVerified.value = false
  sliderPos.value = 0
}

const resetForm = () => {
  showOTPForm.value = false
  otp.value = ''
  resetCaptcha()
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>