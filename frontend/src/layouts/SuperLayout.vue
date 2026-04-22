<template>
  <div class="min-h-screen bg-[#F8FAFC] flex font-sans text-slate-900">
    
    <!-- SIDEBAR -->
    <aside class="w-20 lg:w-72 bg-white border-r border-slate-200 flex flex-col transition-all duration-300 z-30 shadow-[4px_0_24px_rgba(0,0,0,0.02)]">
      
   
      <div class="p-6 lg:p-8 border-b border-slate-50/50 mb-2">
        <div class="flex items-center gap-4">
          <!-- Avatar -->
          <div class="w-11 h-11 rounded-xl bg-indigo-600 flex-shrink-0 flex items-center justify-center shadow-lg shadow-indigo-200 overflow-hidden">
            <img 
              v-if="userProfile?.email"
              :src="`https://ui-avatars.com/api/?name=${userProfile.email}&background=4f46e5&color=fff&bold=true`" 
              alt="Admin Profile" 
              class="w-full h-full object-cover"
            >
            <span v-else class="text-white font-black text-lg tracking-tighter">SA</span>
          </div>
          
          <!-- Info User -->
          <div class="hidden lg:block overflow-hidden">
            <h2 class="font-black text-sm tracking-tight leading-none text-slate-900 truncate">
              {{ userProfile?.email ? userProfile.email.split('@')[0] : 'Super Admin' }}
            </h2>
            <p class="text-[9px] text-indigo-500 font-bold uppercase tracking-[0.1em] mt-0.5">Super Admin</p>
            <div class="flex items-center gap-1 mt-1">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
              <span class="text-[8px] text-emerald-500 font-bold uppercase tracking-tighter">Online</span>
            </div>
          </div>
        </div>
        
        <!-- Mobile View (Hanya Avatar) -->
        <div class="lg:hidden flex justify-center">
          <div class="w-11 h-11 rounded-xl bg-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-200 overflow-hidden">
            <img 
              v-if="userProfile?.email"
              :src="`https://ui-avatars.com/api/?name=${userProfile.email}&background=4f46e5&color=fff&bold=true`" 
              alt="Admin Profile" 
              class="w-full h-full object-cover"
            >
            <span v-else class="text-white font-black text-lg tracking-tighter">SA</span>
          </div>
        </div>
      </div>

      <!-- NAVIGATION MENU -->
      <nav class="flex-1 px-4 space-y-1.5 mt-2">
        <p class="hidden lg:block px-4 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-4">Main Menu</p>
        
        <router-link to="/super-admin" class="nav-item" active-class="active-link" exact-active-class="active-link">
          <div class="nav-icon"><i class="fas fa-th-large"></i></div>
          <span class="nav-text">Dashboard</span>
        </router-link>
        
        <router-link to="/super-admin/instansi" class="nav-item" active-class="active-link">
          <div class="nav-icon"><i class="fas fa-school"></i></div>
          <span class="nav-text">Instansi</span>
        </router-link>

        <router-link to="/super-admin/manage" class="nav-item" active-class="active-link">
          <div class="nav-icon"><i class="fas fa-user-shield"></i></div>
          <span class="nav-text">Admin Sekolah</span>
        </router-link>

        <router-link to="/super-admin/manajemen-nisn" class="nav-item" active-class="active-link">
          <div class="nav-icon"><i class="fas fa-id-card"></i></div>
          <span class="nav-text">Manajemen NISN</span>
          <span class="hidden lg:inline-block ml-auto bg-red-500 text-white text-[9px] font-bold px-2 py-0.5 rounded-full">!</span>
        </router-link>
      </nav>  

      <!-- LOGOUT BUTTON -->
      <div class="p-4 lg:p-6 border-t border-slate-50">
        <button @click="handleLogout" class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all duration-300 w-full group">
          <div class="w-8 h-8 flex items-center justify-center rounded-xl bg-slate-50 text-slate-400 group-hover:bg-red-500 group-hover:text-white transition-all">
            <i class="fas fa-sign-out-alt"></i>
          </div>
          <span class="hidden lg:block font-bold text-[11px] uppercase tracking-widest">Keluar System</span>
        </button>
      </div>
    </aside>

    <!-- MAIN CONTENT AREA -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- ✅ HEADER SEKARANG SIMPLE (TANPA AVATAR) -->
      <header class="h-24 bg-white/70 backdrop-blur-xl border-b border-slate-100 px-10 flex items-center z-20">
        <div>
          <h1 class="text-xl font-black text-slate-900 tracking-tight">
            {{ getGreeting() }}, {{ userProfile?.email ? userProfile.email.split('@')[0] : 'Admin' }}! 👋
          </h1>
          <div class="flex items-center gap-2 mt-1">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <p class="text-[10px] text-slate-400 font-bold uppercase tracking-[0.15em]">Super Admin Panel • SKL System</p>
          </div>
        </div>
      </header>

      <div class="flex-1 overflow-y-auto p-8 lg:p-12 scroll-smooth bg-slate-50/30">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const router = useRouter()
const auth = useAuthStore()

const userProfile = ref(null)

const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Selamat Pagi'
  if (hour < 18) return 'Selamat Siang'
  return 'Selamat Malam'
}

const fetchSuperProfile = async () => {
  try {
    const res = await api.get('/me')
    if (res.data.status === 'success') {
      userProfile.value = res.data.data
    }
  } catch (err) {
    console.error("Gagal mengambil profil super admin:", err)
  }
}

const handleLogout = () => {
  auth.logout()
  router.push('/login')
}

onMounted(() => {
  fetchSuperProfile()
})
</script>

<style scoped>
/* Base Link Style */
.nav-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.875rem 1rem;
  border-radius: 1rem;
  color: #94a3b8;
  transition: all 0.3s ease;
  text-decoration: none;
}

.nav-item:not(.active-link):hover {
  background-color: #f8fafc;
  color: #475569;
}

.nav-icon {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
}

.nav-text {
  display: none;
  font-weight: 700;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

@media (min-width: 1024px) {
  .nav-text { display: block; }
}

.active-link {
  background-color: #0f172a !important;
  color: #ffffff !important;
  box-shadow: 0 20px 25px -5px rgba(15, 23, 42, 0.1), 0 8px 10px -6px rgba(15, 23, 42, 0.1);
}

.active-link .nav-icon, .active-link .nav-text {
  color: #ffffff !important;
}

::-webkit-scrollbar { width: 5px; }
::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
</style>