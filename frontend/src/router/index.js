import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // ========== AUTH ==========
  { 
    path: '/login', 
    name: 'login', 
    component: () => import('../views/auth/login.vue') 
  },
  { 
    path: '/:slug/login-siswa', 
    name: 'LoginSiswa', 
    component: () => import('../views/auth/login_siswa.vue'),
    meta: { public: true }
  },
  
  // ========== ERRORS ==========
  { 
    path: '/:slug/not-found', 
    name: 'SchoolNotFound', 
    component: () => import('../components/ErrorSekolah.vue'),
    meta: { public: true }
  },
  { 
    path: '/admin-error', 
    name: 'AdminErrorPage', 
    component: () => import('../components/ErrorSekolah.vue'),
    meta: { public: true }
  },
  
  // ========== SUPER ADMIN ==========
  {
    path: '/super-admin',
    component: () => import('../layouts/SuperLayout.vue'),
    meta: { requiresAuth: true, role: 'super_admin' },
    children: [
      { 
        path: '', 
        name: 'super-dashboard', 
        component: () => import('../views/super_admin/dashboard.vue') 
      },
      { 
        path: 'instansi', 
        name: 'super-instansi', 
        component: () => import('../views/super_admin/instansi.vue') 
      },
      { 
        path: 'instansi/tambah', 
        name: 'super-instansi-tambah', 
        component: () => import('../views/super_admin/add_instansi.vue') 
      },
      { 
        path: 'instansi/edit/:id', 
        name: 'super-instansi-edit', 
        component: () => import('../views/super_admin/edit_instansi.vue') 
      },
      { 
        path: 'manage', 
        name: 'super-manage', 
        component: () => import('../views/super_admin/manage.admin.vue') 
      },
      { 
        path: 'manajemen-nisn', 
        name: 'super-manajemen-nisn', 
        component: () => import('../views/super_admin/manajemen_nisn.vue') 
      },
      {
        path: 'monitoring',
        name: 'super-monitoring',
        component: () => import('../views/super_admin/monitoring.vue')
      },
    ]
  },

  // ========== SISWA ==========
  {
    path: '/:slug/portal',
    component: () => import('../layouts/SiswaLayout.vue'),
    meta: { requiresAuth: true, role: 'siswa' },
    children: [
      { 
        path: 'skl', 
        name: 'SiswaCetak', 
        component: () => import('../views/siswa/dashboard.vue') 
      }
    ]
  },

  // ========== ADMIN ==========
  {
    path: '/:slug',
    component: () => import('../layouts/AdmLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { 
        path: 'dashboard', 
        name: 'admin-dashboard', 
        component: () => import('../views/admin/dashboard.vue'), 
        meta: { role: ['admin'] } 
      },
      { 
        path: 'admin/mapel', 
        name: 'AdminMapel', 
        component: () => import('../views/admin/mapel.vue'), 
        meta: { role: 'admin' } 
      },
      { 
        path: 'admin/siswa', 
        name: 'AdminSiswa', 
        component: () => import('../views/admin/siswa.vue'), 
        meta: { role: 'admin' } 
      },
      { 
        path: 'admin/foto-siswa', 
        name: 'FotoSiswa', 
        component: () => import('../views/admin/foto_siswa.vue'),
        meta: { role: 'admin' } 
      },
      { 
        path: 'admin/nilai', 
        name: 'AdminNilai', 
        component: () => import('../views/admin/nilai.vue'), 
        meta: { role: ['admin'] } 
      },
      { 
        path: 'admin/leger', 
        name: 'AdminLeger', 
        component: () => import('../views/admin/leger.vue'), 
        meta: { role: 'admin' } 
      },
      { 
        path: 'admin/template-skl', 
        name: 'AdminTemplateSKL', 
        component: () => import('../views/admin/TemplateSKL.vue'), 
        meta: { role: 'admin' } 
      },
      { 
        path: 'admin/setting-waktu', 
        name: 'SettingWaktu', 
        component: () => import('../views/admin/setting_waktu.vue'), 
        meta: { role: 'admin' } 
      },
      {
      path: '/:slug/admin/setting-pesan',
      name: 'SettingPesan',
     component: () => import('@/views/admin/setting_pesan.vue'),
      meta: { title: 'Pengaturan Pesan', requiresAuth: true }
      },
     
{
  path: '/:slug/admin/setting-background',
  name: 'AdminBackground',
  component: () => import('@/views/admin/setting_background.vue'),
  meta: { requiresAuth: true, role: 'admin' }
},
      { 
        path: 'admin/cetak/:id', 
        name: 'AdminCetakSKL', 
        component: () => import('../views/admin/adm_cetak.vue'), 
        meta: { role: ['admin'] } 
      },
    ]
  },

  // ========== DEFAULT ==========
  { path: '/', redirect: '/login' },
  { path: '/:pathMatch(.*)*', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  
  if (to.params.slug && to.params.slug !== to.params.slug.toLowerCase()) {
    const lowercaseSlug = to.params.slug.toLowerCase()
    const newPath = to.path.replace(to.params.slug, lowercaseSlug)
    console.log(`[ROUTER] Redirect slug: ${to.params.slug} → ${lowercaseSlug}`)
    return next(newPath)
  }
  
  console.log('[ROUTER] Navigasi ke:', to.path, 'dari:', from.path)
  
  // Route public selalu diizinkan
  if (to.meta.public || to.path === '/login' || to.path.includes('/login-siswa') || to.path.includes('/not-found')) {
    console.log('[ROUTER] Route public, diizinkan')
    return next()
  }

  const auth = useAuthStore()
  const isLoggedIn = auth.isLoggedIn || localStorage.getItem('isLoggedIn') === 'true'
  const role = auth.role || localStorage.getItem('role')
  const slug = auth.slug || localStorage.getItem('slug')

  console.log('[ROUTER] isLoggedIn:', isLoggedIn, 'role:', role, 'slug:', slug)

  // Belum login, redirect ke login
  if (!isLoggedIn) {
    console.log('[ROUTER] Belum login, redirect ke /login')
    return next('/login')
  }

  if (role === 'super_admin') {
    // Izinkan akses ke semua route /super-admin/*
    if (to.path.startsWith('/super-admin')) {
      console.log('[ROUTER] Super admin, route diizinkan:', to.path)
      return next()
    }
    // Redirect ke dashboard super admin
    console.log('[ROUTER] Super admin redirect ke /super-admin')
    return next('/super-admin')
  }

  // Khusus untuk role SISWA
  if (role === 'siswa') {
    if (to.path.includes('/portal') || to.path === `/${slug}/portal/skl`) {
      console.log('[ROUTER] Siswa, route diizinkan')
      return next()
    }
    console.log('[ROUTER] Siswa, redirect ke dashboard')
    return next(`/${slug}/portal/skl`)
  }

  // Proteksi route super admin dari user lain
  if (to.meta.role && Array.isArray(to.meta.role) && to.meta.role.includes('super_admin')) {
    if (role !== 'super_admin') {
      console.log('[ROUTER] Bukan super admin, redirect ke /admin-error')
      return next('/admin-error')
    }
  }

  //  Cek tenant untuk admin
  if (to.params.slug && role === 'admin' && to.params.slug !== slug) {
    console.log('[ROUTER] Slug tidak cocok untuk admin')
    return next('/admin-error')
  }

  // Cek role umum
  if (to.meta.role) {
    const allowedRoles = Array.isArray(to.meta.role) ? to.meta.role : [to.meta.role]
    if (!allowedRoles.includes(role)) {
      console.log('[ROUTER] Role tidak sesuai, redirect ke dashboard')
      if (role === 'admin') return next(`/${slug}/dashboard`)
      if (role === 'siswa') return next(`/${slug}/portal/skl`)
      if (role === 'super_admin') return next('/super-admin')
      return next('/login')
    }
  }

  console.log('[ROUTER] Route diizinkan')
  return next()
})

export default router