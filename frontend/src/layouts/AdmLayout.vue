<template>
  <div v-if="isImpersonating" class="fixed top-0 left-0 right-0 z-50 bg-amber-100 border-b border-amber-200 text-amber-900 py-2 shadow-sm">
    <div class="container mx-auto px-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <i class="fas fa-user-shield text-amber-600"></i>
        <span class="text-xs md:text-sm font-medium">
          Anda sedang mengakses: <span class="font-bold">{{ impersonatedInstansi }}</span>
        </span>
      </div>
      <button 
        @click="exitImpersonate" 
        class="text-xs font-bold px-3 py-1 bg-white border border-amber-300 rounded hover:bg-amber-50 transition-all flex items-center gap-2"
      >
        <i class="fas fa-sign-out-alt"></i>
        Keluar
      </button>
    </div>
  </div>

  <div :class="{ 'pt-12': isImpersonating }" class="min-h-screen bg-[#F4F6F9] flex font-sans text-slate-700">
    
    <!-- SIDEBAR -->
    <aside class="w-64 bg-[#343A40] flex flex-col sticky top-0 h-screen shadow-xl z-30">
      
      <!-- NAMA SEKOLAH + LOGO (DINAMIS) -->
      <div class="p-5 border-b border-slate-700 flex items-center gap-3">
        <!-- LOGO SEKOLAH -->
        <div class="w-10 h-10 rounded-full flex items-center justify-center shadow-lg overflow-hidden bg-white">
          <img 
            v-if="userProfile?.logo_instansi"
            :src="`http://localhost:3000/uploads/instansi/${userProfile.logo_instansi}`"
            class="w-full h-full object-cover"
            alt="Logo Sekolah"
          />
          <div v-else class="w-full h-full bg-emerald-500 flex items-center justify-center">
            <i class="fas fa-graduation-cap text-white text-sm"></i>
          </div>
        </div>
        <span class="text-white font-bold text-lg tracking-tight truncate">{{ auth.nama_instansi || 'NAMA SEKOLAH' }}</span>
      </div>

      <!-- NAVIGATION MENU -->
      <nav class="flex-1 overflow-y-auto py-4 custom-scrollbar">
        <p class="px-6 text-[10px] font-bold text-slate-500 uppercase mb-2 tracking-widest">Main Menu</p>
        
        <router-link :to="`/${route.params.slug}/dashboard`" 
                     class="flex items-center gap-3 px-6 py-3 text-slate-300 hover:bg-slate-700 hover:text-white transition-all border-l-4 border-transparent" 
                     active-class="active-nav">
          <i class="fas fa-tachometer-alt w-5"></i>
          <span class="text-sm">Dashboard</span>
        </router-link>

        <div v-for="(item, index) in menuItems" :key="index">
          <router-link :to="item.to" 
                       class="flex items-center gap-3 px-6 py-3 text-slate-300 hover:bg-slate-700 hover:text-white transition-all border-l-4 border-transparent" 
                       active-class="active-nav">
            <i :class="[item.icon, 'w-5']"></i>
            <span class="text-sm">{{ item.label }}</span>
          </router-link>
        </div>

        <p class="px-6 text-[10px] font-bold text-slate-500 uppercase mt-6 mb-2 tracking-widest">Pengaturan</p>
        <button @click="auth.logout()" class="w-full flex items-center gap-3 px-6 py-3 text-red-400 hover:bg-red-900/20 transition-all border-l-4 border-transparent">
          <i class="fas fa-power-off w-5"></i>
          <span class="text-sm font-bold italic">Logout</span>
        </button>
      </nav>
    </aside>

    <!-- MAIN CONTENT -->
    <main class="flex-1 flex flex-col overflow-hidden relative">
      
      <!-- HEADER -->
      <header class="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-6 shadow-sm z-20">
        <div class="flex items-center gap-4">
          <button class="text-slate-500 hover:text-slate-800 p-2 rounded-lg hover:bg-slate-100 transition-colors">
            <i class="fas fa-bars text-lg"></i>
          </button>
          <h1 class="text-lg font-bold text-slate-800 hidden md:block capitalize">{{ route.meta.title || 'Dashboard' }}</h1>
        </div>
        
        <div class="flex items-center gap-6">
          
          <!-- PILIH ANGKATAN -->
          <div class="relative group">
            <button class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-bold text-slate-600 hover:bg-indigo-50 hover:text-indigo-600 transition-all border border-transparent hover:border-indigo-100">
              <i class="fas fa-calendar-alt text-indigo-500"></i>
              <span class="hidden sm:inline">Lulusan - {{ auth.selectedYear }}</span>
              <span class="sm:hidden">{{ auth.selectedYear }}</span>
              <i class="fas fa-chevron-down text-[10px] text-slate-400 ml-1"></i>
            </button>

            <div class="absolute right-0 top-full mt-2 w-48 bg-white rounded-xl shadow-xl border border-slate-100 hidden group-hover:block z-50 animate-fade-in-down origin-top-right overflow-hidden">
              <div class="py-1">
                <div class="px-4 py-2 text-[10px] font-bold text-slate-400 uppercase bg-slate-50 border-b border-slate-50">Pilih Angkatan</div>
                <button 
                  v-for="year in [2024, 2025, 2026, 2027]" 
                  :key="year"
                  @click="changeYear(year)"
                  class="w-full text-left px-4 py-3 text-sm text-slate-600 hover:bg-indigo-50 hover:text-indigo-600 flex justify-between items-center transition-colors"
                  :class="{ 'bg-indigo-50 text-indigo-700 font-bold': auth.selectedYear == year }"
                >
                  <span>Angkatan {{ year }}</span>
                  <i v-if="auth.selectedYear == year" class="fas fa-check text-xs text-indigo-600"></i>
                </button>
              </div>
            </div>
          </div>

          <div class="h-8 w-[1px] bg-slate-200 hidden sm:block"></div>

          <!-- DROPDOWN PROFIL -->
          <div class="relative">
            <button @click="toggleProfileDropdown" class="flex items-center gap-3 focus:outline-none px-4 py-2 rounded-lg hover:bg-slate-100 transition-colors border border-transparent hover:border-slate-200">
              <div class="w-8 h-8 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-xs font-bold overflow-hidden">
                <img 
                  v-if="userProfile?.foto_profile"
                  :src="`http://localhost:3000/uploads/admin/${userProfile.foto_profile}`"
                  class="w-full h-full object-cover"
                  alt="Profile"
                />
                <i v-else class="fas fa-user"></i>
              </div>
              
              <div class="hidden md:flex flex-col items-start">
                <span class="text-xs font-bold text-slate-700 leading-none">{{ userProfile?.email ? userProfile.email.split('@')[0] : 'Admin' }}</span>
                <span class="text-[10px] text-slate-400 font-medium">Lihat Profil</span>
              </div>
              
              <i class="fas fa-chevron-down text-[10px] text-slate-400 transition-transform" :class="{ 'rotate-180': isProfileOpen }"></i>
            </button>

            <!-- DROPDOWN CONTENT -->
            <div v-if="isProfileOpen" class="absolute right-0 top-14 w-72 bg-white rounded-xl shadow-2xl border border-slate-100 p-5 animate-fade-in-down origin-top-right z-50">
              
              <!-- HEADER DENGAN FOTO PROFIL -->
              <div class="flex items-center gap-4 mb-4 pb-4 border-b border-slate-100">
                <div class="relative group">
                  <div class="w-16 h-16 rounded-full bg-indigo-50 text-indigo-600 flex items-center justify-center text-xl font-black border-2 border-indigo-200 shadow-sm overflow-hidden">
                    <img 
                      v-if="userProfile?.foto_profile"
                      :src="`http://localhost:3000/uploads/admin/${userProfile.foto_profile}`"
                      class="w-full h-full object-cover"
                      alt="Profile"
                    />
                    <span v-else>{{ profileInitials }}</span>
                  </div>
                  
                  <!-- Tombol Upload -->
                  <label class="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                    <i class="fas fa-camera text-white text-sm"></i>
                    <input type="file" accept="image/jpeg,image/png" class="hidden" @change="handleUploadFoto" />
                  </label>
                </div>
                
                <div class="overflow-hidden flex-1">
                  <h4 class="text-sm font-bold text-slate-800 truncate">{{ userProfile?.email || 'Loading...' }}</h4>
                  <p class="text-[10px] text-slate-400 font-medium uppercase">{{ userProfile?.role === 'super_admin' ? 'Super Administrator' : 'Admin Sekolah' }}</p>
                  
                  <button 
                    v-if="userProfile?.foto_profile"
                    @click="handleDeleteFoto"
                    class="text-[9px] text-red-500 hover:text-red-700 mt-1"
                  >
                    <i class="fas fa-trash-alt mr-1"></i>Hapus Foto
                  </button>
                </div>
              </div>

              <!-- INFO SEKOLAH (TANPA LOGO - BERSIH) -->
              <div class="space-y-3 mb-4">
                <div class="flex justify-between items-center p-2 rounded-lg bg-slate-50">
                  <span class="text-[10px] font-bold text-slate-400 uppercase">Sekolah</span>
                  <span class="text-xs font-bold text-slate-700 truncate max-w-[150px]">{{ userProfile?.nama_instansi || '-' }}</span>
                </div>
                <div class="flex justify-between items-center p-2 rounded-lg bg-slate-50">
                  <span class="text-[10px] font-bold text-slate-400 uppercase">Slug Login</span>
                  <span class="text-xs font-mono bg-white px-2 py-0.5 rounded text-indigo-600 border border-indigo-100 shadow-sm">
                    /{{ userProfile?.slug || '-' }}
                  </span>
                </div>
              </div>

              <!-- TOMBOL LOGOUT -->
              <div class="pt-4 border-t border-slate-50">
                <button @click="auth.logout()" class="w-full py-2 rounded-lg bg-red-50 text-red-600 text-[10px] font-bold hover:bg-red-100 transition-colors border border-red-100">
                  <i class="fas fa-sign-out-alt mr-1"></i>Logout
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- CONTENT AREA -->
      <div class="p-6 md:p-8 overflow-y-auto flex-1 bg-[#F4F6F9]">
        <router-view /> 
      </div>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import Swal from 'sweetalert2'

const auth = useAuthStore()
const route = useRoute()

const isProfileOpen = ref(false)
const userProfile = ref(null)
const isImpersonating = ref(false)
const impersonatedInstansi = ref('')

const profileInitials = computed(() => {
  if (!userProfile.value?.email) return 'AD'
  return userProfile.value.email.substring(0, 2).toUpperCase()
})

const toggleProfileDropdown = () => {
  isProfileOpen.value = !isProfileOpen.value
  if (!userProfile.value && isProfileOpen.value) {
    fetchUserProfile()
  }
}

const fetchUserProfile = async () => {
  try {
    const res = await api.get('/me')
    if (res.data.status === 'success') {
      userProfile.value = res.data.data
    }
  } catch (err) {
    console.error("Gagal ambil profil:", err)
  }
}

const handleUploadFoto = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 2 * 1024 * 1024) {
    Swal.fire('Error', 'Ukuran file maksimal 2MB', 'error')
    return
  }

  const formData = new FormData()
  formData.append('foto', file)

  try {
    const res = await api.post('/upload-foto', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    if (res.data.status === 'success') {
      Swal.fire('Berhasil!', 'Foto profil & logo sekolah berhasil diupload', 'success')
      fetchUserProfile()
    }
  } catch (err) {
    Swal.fire('Error', err.response?.data?.message || 'Gagal upload', 'error')
  }
}

const handleDeleteFoto = async () => {
  const result = await Swal.fire({
    title: 'Hapus Foto?',
    text: 'Foto profil dan logo sekolah akan dihapus permanen',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    confirmButtonText: 'Ya, Hapus',
    cancelButtonText: 'Batal'
  })

  if (result.isConfirmed) {
    try {
      const res = await api.delete('/delete-foto')
      if (res.data.status === 'success') {
        Swal.fire('Berhasil!', 'Foto profil & logo sekolah dihapus', 'success')
        fetchUserProfile()
      }
    } catch (err) {
      Swal.fire('Error', 'Gagal menghapus foto', 'error')
    }
  }
}

const changeYear = (year) => {
  auth.setSelectedYear(year.toString())
}

const exitImpersonate = async () => {
  try {
    const res = await api.post('/restore-super-admin')
    if (res.data.status === 'success') {
      sessionStorage.removeItem('isImpersonating')
      sessionStorage.removeItem('impersonatedInstansi')
      window.location.href = '/super-admin'
    }
  } catch (err) {
    console.error("Restore error:", err)
    window.location.href = '/login'
  }
}

const menuItems = computed(() => {
  const slug = route.params.slug
  if (auth.role === 'super_admin') {
    return [
      { label: 'Dashboard Utama', icon: 'fas fa-home', to: `/super/dashboard` },
      { label: 'Kelola Instansi', icon: 'fas fa-school', to: `/super/instansi` },
    ]
  }

  return [
    { label: 'Daftar Mapel', icon: 'fas fa-book-open', to: `/${slug}/admin/mapel` },
    { label: 'Kelola Data Siswa', icon: 'fas fa-users', to: `/${slug}/admin/siswa` },
    { label: 'Foto Siswa', icon: 'fas fa-camera-retro', to: `/${slug}/admin/foto-siswa` },
    { label: 'Input Nilai Siswa', icon: 'fas fa-pen-nib', to: `/${slug}/admin/nilai` },
    { label: 'Rekap Nilai Siswa', icon: 'fas fa-file-invoice', to: `/${slug}/admin/leger` },
    { label: 'Setting SKL', icon: 'fas fa-sliders-h', to: `/${slug}/admin/template-skl` },
    { label: 'Jadwal Pengumuman', icon: 'fas fa-clock', to: `/${slug}/admin/setting-waktu` },
  ]
})

onMounted(() => {
  isImpersonating.value = sessionStorage.getItem('isImpersonating') === 'true'
  impersonatedInstansi.value = sessionStorage.getItem('impersonatedInstansi') || ''
  fetchUserProfile()
})
</script>

<style scoped>
.active-nav {
  background-color: #1a1e21;
  border-left-color: #10b981 !important;
  color: white !important;
}

@keyframes fadeInDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-down {
  animation: fadeInDown 0.2s ease-out forwards;
}
</style>