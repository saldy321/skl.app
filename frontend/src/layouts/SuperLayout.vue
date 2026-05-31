<template>
  <div class="min-h-screen bg-[#F5F7FA] flex font-sans text-slate-800">
    
    <!-- SIDEBAR -->
    <aside class="w-72 bg-white border-r border-slate-200 flex flex-col transition-all duration-300 z-30 fixed top-0 left-0 h-full">
      
      <!-- LOGO AREA -->
      <div class="p-6 border-b border-slate-100">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-full bg-red-500 flex items-center justify-center">
            <i class="fas fa-store text-white text-xs"></i>
          </div>
          <span class="text-xl font-bold text-slate-800 tracking-tight">SKL<span class="text-red-500">System</span></span>
        </div>
      </div>

      <!-- NAVIGATION MENU -->
      <nav class="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
        <div class="px-3 text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-4">Main Menu</div>
        
        <router-link to="/super-admin" class="nav-item" active-class="active-link" exact-active-class="active-link">
          <div class="nav-icon"><i class="fas fa-home"></i></div>
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
          <span class="inline-block ml-auto bg-red-500 text-white text-[9px] font-bold px-2 py-0.5 rounded-full">!</span>
        </router-link>

        <router-link to="/super-admin/monitoring" class="nav-item" active-class="active-link">
          <div class="nav-icon"><i class="fas fa-cloud-upload-alt"></i></div>
          <span class="nav-text">Backup & Restore</span>
        </router-link>
      </nav>  

      <!-- LOGOUT -->
      <div class="p-4 border-t border-slate-100">
        <button @click="handleLogout" class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-500 hover:text-red-500 hover:bg-red-50 transition-all w-full">
          <i class="fas fa-sign-out-alt w-5 text-center"></i>
          <span class="text-sm font-semibold">Keluar</span>
        </button>
      </div>
    </aside>

    <!-- MAIN CONTENT -->
    <div class="ml-72 flex-1 flex flex-col min-h-screen">
      
      <!-- TOP BAR -->
      <header class="h-16 bg-white border-b border-slate-200 px-8 flex items-center justify-between sticky top-0 z-20">
        
        <!-- Left: Breadcrumb / Greeting -->
        <div>
          <h1 class="text-lg font-bold text-slate-800">
            {{ getGreeting() }}, {{ userProfile?.email?.split('@')[0] || 'Admin' }}!
          </h1>
          <p class="text-xs text-slate-400">Super Admin Panel</p>
        </div>

        <!-- Right: Profile (Clickable) -->
        <div class="flex items-center gap-3 pl-4 border-l border-slate-200">
          <button 
            @click="showProfileInfo" 
            class="flex items-center gap-3 hover:bg-slate-50 rounded-xl px-3 py-2 transition-colors cursor-pointer"
          >
            <div class="w-10 h-10 rounded-full bg-indigo-600 flex items-center justify-center text-white font-bold shadow-sm overflow-hidden">
              <img 
                v-if="userProfile?.email"
                :src="`https://ui-avatars.com/api/?name=${userProfile.email}&background=4f46e5&color=fff&bold=true`" 
                alt="Profile" 
                class="w-full h-full object-cover"
              >
              <span v-else class="text-sm">SA</span>
            </div>
            <div class="hidden sm:block text-left">
              <p class="text-sm font-semibold text-slate-800">{{ userProfile?.email?.split('@')[0] || 'Super Admin' }}</p>
              <p class="text-[10px] text-slate-400">Super Admin</p>
            </div>
          </button>
        </div>

      </header>

      <!-- PAGE CONTENT -->
      <main class="flex-1 p-8 overflow-y-auto bg-[#F5F7FA]">
        <router-view />
      </main>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import Swal from 'sweetalert2'

const router = useRouter()
const auth = useAuthStore()

const userProfile = ref(null)

const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Selamat Pagi'
  if (hour < 18) return 'Selamat Siang'
  return 'Selamat Malam'
}

const showProfileInfo = () => {
  if (userProfile.value?.email) {
    Swal.fire({
      title: 'Super Admin',
      html: `
        <div class="text-left">
          <p><strong>Email:</strong> ${userProfile.value.email}</p>
          <p><strong>Role:</strong> ${userProfile.value.role || 'Super Admin'}</p>
        </div>
      `,
      icon: 'info',
      confirmButtonColor: '#4f46e5',
      confirmButtonText: 'OK'
    })
  } else {
    Swal.fire('Info', 'Data profil belum tersedia.', 'info')
  }
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
.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 1rem;
  border-radius: 0.75rem;
  color: #64748b;
  transition: all 0.2s ease;
  text-decoration: none;
}

.nav-item:not(.active-link):hover {
  background-color: #f1f5f9;
  color: #334155;
}

.nav-icon {
  width: 1.75rem;
  height: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.nav-text {
  font-weight: 600;
  font-size: 0.875rem;
}

.active-link {
  background-color: #0f172a !important;
  color: #ffffff !important;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.15);
}

.active-link .nav-icon, .active-link .nav-text {
  color: #ffffff !important;
}
</style>