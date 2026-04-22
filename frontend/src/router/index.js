import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/login', name: 'login', component: () => import('../views/login.vue') },
  
  { 
    path: '/:slug/login-siswa', 
    name: 'LoginSiswa', 
    component: () => import('../views/LoginSiswa.vue'),
    meta: { public: true }
  },
  
  { 
    path: '/:slug/not-found', 
    name: 'SchoolNotFound', 
    component: () => import('../components/ErrorSekolah.vue'),
    meta: { public: true }
  },
  
  { 
    path: '/admin-error', 
    name: 'AdminErrorPage', 
    component: () => import('../components/AdminErrorSekolah.vue'),
    meta: { public: true }
  },
  
  {
    path: '/super-admin',
    component: () => import('../layouts/SuperLayout.vue'),
    meta: { requiresAuth: true, role: 'super_admin' },
    children: [
      { path: '', name: 'super-dashboard', component: () => import('../views/dashboard.vue') },
      { path: 'instansi', name: 'super-instansi', component: () => import('../views/super_admin/instansi.vue') },
      { path: 'instansi/tambah', name: 'super-instansi-tambah', component: () => import('../views/super_admin/add_instansi.vue') },
      { path: 'instansi/edit/:id', name: 'super-instansi-edit', component: () => import('../views/super_admin/add_instansi.vue') },
      { path: 'manage', name: 'super-manage', component: () => import('../views/super_admin/manage.admin.vue') },
        { path: 'manajemen-nisn', name: 'super-manajemen-nisn', component: () => import('../views/super_admin/manajemen_nisn.vue') },
    ]
  },

  {
    path: '/:slug/portal',
    component: () => import('../layouts/SiswaLayout.vue'),
    meta: { requiresAuth: true, role: 'siswa' },
    children: [
      { path: 'skl', name: 'SiswaCetak', component: () => import('../views/dashboard_siswa.vue') }
    ]
  },

  {
    path: '/:slug',
    component: () => import('../layouts/AdmLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: 'dashboard', name: 'admin-dashboard', component: () => import('../views/dashboard_admin.vue'), meta: { role: ['admin'] } },
      { path: 'admin/mapel', name: 'AdminMapel', component: () => import('../views/mapel.vue'), meta: { role: 'admin' } },
      { path: 'admin/siswa', name: 'AdminSiswa', component: () => import('../views/siswa.vue'), meta: { role: 'admin' } },
      { path: 'admin/foto-siswa', name: 'FotoSiswa', component: () => import('../views/foto_siswa.vue'), meta: { role: 'admin' } },
      { path: 'admin/nilai', name: 'AdminNilai', component: () => import('../views/nilai.vue'), meta: { role: ['admin'] } },
      { path: 'admin/leger', name: 'AdminLeger', component: () => import('../views/leger.vue'), meta: { role: 'admin' } },
      { path: 'admin/template-skl', name: 'AdminTemplateSKL', component: () => import('../views/TemplateSKL.vue'), meta: { role: 'admin' } },
      { path: 'admin/setting-waktu', name: 'SettingWaktu', component: () => import('../views/setting_waktu.vue'), meta: { role: 'admin' } },
      { path: 'admin/cetak/:id', name: 'AdminCetakSKL', component: () => import('../views/adm_cetak.vue'), meta: { role: ['admin'] } },
    

    ]
  },

  { path: '/', redirect: '/login' },
  { path: '/:pathMatch(.*)*', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// ✅ GUARD LENGKAP DENGAN LOGGING + PENANGANAN SISWA
router.beforeEach((to, from) => {
  console.log('[ROUTER] Navigasi ke:', to.path, 'dari:', from.path)
  
  // Route public selalu diizinkan
  if (to.meta.public || to.path === '/login' || to.path.includes('/login-siswa') || to.path.includes('/not-found')) {
    console.log('[ROUTER] Route public, diizinkan')
    return true
  }

  const auth = useAuthStore()
  const isLoggedIn = auth.isLoggedIn || localStorage.getItem('isLoggedIn') === 'true'
  const role = auth.role || localStorage.getItem('role')
  const slug = auth.slug || localStorage.getItem('slug')

  console.log('[ROUTER] isLoggedIn:', isLoggedIn, 'role:', role, 'slug:', slug)

  // Belum login, redirect ke login
  if (!isLoggedIn) {
    console.log('[ROUTER] Belum login, redirect ke /login')
    return '/login'
  }

  // ✅ TAMBAHKAN: Khusus untuk role SISWA
  if (role === 'siswa') {
    // Izinkan akses ke route portal siswa
    if (to.path.includes('/portal') || to.path === `/${slug}/portal/skl`) {
      console.log('[ROUTER] Siswa, route diizinkan')
      return true
    }
    // Selain itu, redirect ke dashboard siswa
    console.log('[ROUTER] Siswa, redirect ke dashboard')
    return `/${slug}/portal/skl`
  }

  // Cek role untuk super admin
  if (to.meta.role === 'super_admin' && role !== 'super_admin') {
    console.log('[ROUTER] Bukan super admin, redirect ke /admin-error')
    return '/admin-error'
  }

  // Cek tenant untuk admin
  if (to.params.slug && role === 'admin' && to.params.slug !== slug) {
    console.log('[ROUTER] Slug tidak cocok untuk admin')
    return '/admin-error'
  }

  // Cek role umum
  if (to.meta.role) {
    const allowedRoles = Array.isArray(to.meta.role) ? to.meta.role : [to.meta.role]
    if (!allowedRoles.includes(role)) {
      console.log('[ROUTER] Role tidak sesuai, redirect ke dashboard masing-masing')
      if (role === 'admin') return `/${slug}/dashboard`
      if (role === 'siswa') return `/${slug}/portal/skl`
      if (role === 'super_admin') return '/super-admin'
      return '/login'
    }
  }

  console.log('[ROUTER] Route diizinkan')
  return true
})

export default router