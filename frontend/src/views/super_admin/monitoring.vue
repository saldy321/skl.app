<template>
  <div class="space-y-6">
    
    <!-- HEADER -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 class="text-lg font-semibold text-gray-900">Backup & Restore Data</h1>
          <p class="text-sm text-gray-500 mt-0.5">Kelola backup otomatis dan pemulihan data instansi</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="w-2 h-2 rounded-full" :class="setting.auto_backup ? 'bg-emerald-500 animate-pulse' : 'bg-gray-300'"></span>
          <span class="text-xs font-medium text-gray-600">{{ setting.auto_backup ? 'Auto Backup Aktif' : 'Auto Backup Nonaktif' }}</span>
          <span v-if="setting.auto_backup" class="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded-full">
            ⏰ Backup jam {{ formatDuaDigit(jamBackup) }}:{{ formatDuaDigit(menitBackup) }} WIB
          </span>
        </div>
      </div>
    </div>

    <!-- METRICS -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center"><i class="fas fa-school text-blue-600 text-sm"></i></div>
          <div><p class="text-2xl font-bold text-gray-900">{{ instansiList.length }}</p><p class="text-xs text-gray-500 font-medium">Total Instansi</p></div>
        </div>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-amber-50 flex items-center justify-center"><i class="fas fa-clock text-amber-600 text-sm"></i></div>
          <div>
            <p class="text-2xl font-bold text-gray-900">{{ formatDuaDigit(jamBackup) }}:{{ formatDuaDigit(menitBackup) }}</p>
            <p class="text-xs text-gray-500 font-medium">WIB • Setiap Hari</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm p-5">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-indigo-50 flex items-center justify-center"><i class="fas fa-history text-indigo-600 text-sm"></i></div>
          <div><p class="text-2xl font-bold text-gray-900">{{ setting.retensi_hari }}</p><p class="text-xs text-gray-500 font-medium">Hari Retensi</p></div>
        </div>
      </div>
      <div @click="showSetting = true" class="bg-white rounded-lg border border-gray-200 shadow-sm p-5 cursor-pointer hover:border-blue-300 hover:shadow transition-all">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-gray-50 flex items-center justify-center"><i class="fas fa-sliders-h text-gray-500 text-sm"></i></div>
          <div><p class="text-sm font-semibold text-gray-700">Pengaturan</p><p class="text-xs text-gray-400">Ubah Konfigurasi</p></div>
        </div>
      </div>
    </div>

    <!-- DAFTAR INSTANSI -->
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-100 bg-gray-50/50"><h2 class="text-sm font-semibold text-gray-700">Daftar Instansi</h2><p class="text-xs text-gray-400 mt-0.5">Klik untuk melihat riwayat backup</p></div>
      <div class="divide-y divide-gray-100">
        <div v-for="item in instansiList" :key="item.id" @click="pilihInstansi(item)" class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 cursor-pointer transition-colors" :class="{ 'bg-blue-50/50 border-l-2 border-blue-500': selectedInstansi?.id === item.id }">
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 rounded-full overflow-hidden bg-white border border-gray-200 flex-shrink-0">
              <img v-if="item.logo_instansi" :src="`http://localhost:3000/uploads/instansi/${item.logo_instansi}`" class="w-full h-full object-cover" />
              <div v-else class="w-full h-full bg-gray-50 text-gray-400 flex items-center justify-center text-xs font-bold">{{ item.nama_instansi?.substring(0,2).toUpperCase() }}</div>
            </div>
            <div><h3 class="text-sm font-semibold text-gray-900">{{ item.nama_instansi }}</h3><p class="text-xs text-gray-400">{{ item.tingkat_sekolah }} • <code class="text-[11px]">{{ item.slug }}</code></p></div>
          </div>
          <i class="fas fa-chevron-right text-gray-300 text-xs"></i>
        </div>
        <div v-if="instansiList.length === 0" class="px-6 py-10 text-center text-gray-400"><i class="fas fa-school text-3xl mb-2 opacity-30"></i><p class="text-sm">Belum ada instansi terdaftar</p></div>
      </div>
    </div>

    <!-- PANEL BACKUP -->
    <div v-if="selectedInstansi" class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
        <div><h3 class="text-sm font-semibold text-gray-900">{{ selectedInstansi.nama_instansi }}</h3><p class="text-xs text-gray-400">Riwayat Backup • {{ setting.retensi_hari }} Hari Terakhir</p></div>
        <div class="flex gap-2">
          <!-- TOMBOL IMPORT (DI HEADER) -->
          <button @click="triggerImport" class="inline-flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg text-xs font-semibold transition-colors">
            <i class="fas fa-upload"></i> Import Backup
          </button>
          <button @click="backupSekarang" :disabled="backupLoading" class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-xs font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
            <i v-if="backupLoading" class="fas fa-circle-notch animate-spin"></i>
            <i v-else class="fas fa-camera"></i>
            {{ backupLoading ? 'Memproses...' : 'Backup Sekarang' }}
          </button>
        </div>
      </div>
      <div v-if="loading" class="p-10 text-center"><i class="fas fa-circle-notch animate-spin text-xl text-blue-500 mb-2"></i><p class="text-sm text-gray-400">Memuat riwayat...</p></div>
      <div v-else-if="Object.keys(backupHistory).length > 0" class="divide-y divide-gray-100">
        <template v-for="(backups, tanggal) in backupHistory" :key="tanggal">
          <div class="px-6 py-2.5 bg-gray-50/80 text-xs font-semibold text-gray-500 uppercase tracking-wide">{{ tanggal }}</div>

          <div v-for="backup in backups" :key="backup.id" class="px-6 py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-gray-50/50 transition-colors">
            <div>
              <p class="text-sm font-medium text-gray-800">{{ backup.nama_backup }}</p>
              <div class="flex items-center gap-3 mt-1 text-xs text-gray-400">
                <span><i class="fas fa-user-graduate mr-1"></i>{{ backup.jumlah_siswa }} siswa</span>
                <span><i class="fas fa-pen mr-1"></i>{{ backup.jumlah_nilai }} nilai</span>
                <span><i class="fas fa-book mr-1"></i>{{ backup.jumlah_mapel }} mapel</span>
              </div>
            </div>
            
            <div class="flex flex-wrap items-center gap-2">
              <!-- DOWNLOAD -->
              <button @click="downloadBackup(backup)" class="inline-flex items-center gap-1.5 bg-blue-500 hover:bg-blue-600 text-white px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors">
                <i class="fas fa-download"></i> Download
              </button>
              
              <!-- RESTORE -->
              <button @click="pulihkanBackup(backup)" class="inline-flex items-center gap-1.5 bg-amber-500 hover:bg-amber-600 text-white px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors">
                <i class="fas fa-undo"></i> Kembalikan
              </button>
            </div>
          </div>    
        </template>
      </div>
      <div v-else class="px-6 py-12 text-center text-gray-400"><i class="fas fa-archive text-4xl mb-3 opacity-30"></i><p class="text-sm">Belum ada riwayat backup</p><p class="text-xs mt-1">Backup akan dibuat otomatis sesuai jadwal</p></div>
    </div>

    <input type="file" ref="fileInput" class="hidden" accept=".json" @change="handleImportFile" />

    <!-- MODAL SETTING -->
    <div v-if="showSetting" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 bg-gray-50/50"><h3 class="text-sm font-semibold text-gray-900">Konfigurasi Auto Backup</h3><p class="text-xs text-gray-400 mt-0.5">Atur jadwal dan retensi data</p></div>
        <div class="p-6 space-y-5">
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-100">
            <div><span class="text-sm font-medium text-gray-700">Aktifkan Auto Backup</span><p class="text-xs text-gray-400 mt-0.5">Backup otomatis setiap hari</p></div>
            <label class="relative inline-flex items-center cursor-pointer"><input type="checkbox" v-model="setting.auto_backup" class="sr-only peer" /><div class="w-9 h-5 bg-gray-200 rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-blue-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all"></div></label>
          </div>

          <!-- INPUT JAM & MENIT (LANGSUNG PAKE NUMBER) -->
          <div>
            <label class="block text-xs font-semibold text-gray-600 mb-1.5">Waktu Backup (WIB)</label>
            <div class="flex items-center gap-2">
              <input 
                v-model.number="jamBackup" 
                @change="validateJam"
                type="number" 
                min="0" 
                max="23"
                class="w-16 border border-gray-300 rounded-lg px-3 py-2.5 text-sm text-center font-mono focus:ring-2 focus:ring-blue-500 outline-none"
              />
              <span class="text-gray-400 font-bold">:</span>
              <input 
                v-model.number="menitBackup" 
                @change="validateMenit"
                type="number" 
                min="0" 
                max="59"
                class="w-16 border border-gray-300 rounded-lg px-3 py-2.5 text-sm text-center font-mono focus:ring-2 focus:ring-blue-500 outline-none"
              />
              <span class="text-xs text-gray-400">WIB</span>
            </div>
            <p class="text-[11px] text-gray-400 mt-1">Format 00 - 23 untuk jam dan 00 - 59 untuk menit</p>
          </div>

          <div>
            <label class="block text-xs font-semibold text-gray-600 mb-1.5">Retensi Data</label>
            <div class="flex items-center gap-2">
              <input 
                v-model.number="setting.retensi_hari" 
                type="number" 
                min="1" 
                max="365"
                class="w-24 border border-gray-300 rounded-lg px-3 py-2.5 text-sm text-center font-mono focus:ring-2 focus:ring-blue-500 outline-none"
              />
              <span class="text-xs text-gray-400">Hari</span>
            </div>
            <p class="text-[11px] text-gray-400 mt-1">Data lebih lama akan dihapus otomatis</p>
          </div>

        </div>
        <div class="px-6 py-4 border-t border-gray-100 bg-gray-50/50 flex gap-3 justify-end">
          <button @click="closeSetting" class="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg">Batal</button>
          <button @click="simpanSetting" :disabled="savingSetting" class="px-5 py-2 text-sm font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg disabled:opacity-50 inline-flex items-center gap-2"><i v-if="savingSetting" class="fas fa-circle-notch animate-spin"></i>{{ savingSetting ? 'Menyimpan...' : 'Simpan Konfigurasi' }}</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '@/api'
import Swal from 'sweetalert2'

// ========== DATA REAKTIF ==========
const instansiList = ref([])
const selectedInstansi = ref(null)
const backupHistory = ref({})
const loading = ref(false)
const backupLoading = ref(false)
const showSetting = ref(false)
const savingSetting = ref(false)

// Setting backup (sesuai backend)
const setting = ref({ 
  auto_backup: false, 
  retensi_hari: 60 
})

// Jam dan menit (int, sesuai backend)
const jamBackup = ref(2)
const menitBackup = ref(0)

const nextBackupTime = ref('')
const timeoutId = ref(null)

// ========== HELPER FUNCTIONS ==========
const formatDuaDigit = (angka) => {
  return String(angka).padStart(2, '0')
}

const validateJam = () => {
  let j = jamBackup.value
  if (isNaN(j)) j = 0
  if (j < 0) j = 0
  if (j > 23) j = 23
  jamBackup.value = j
}

const validateMenit = () => {
  let m = menitBackup.value
  if (isNaN(m)) m = 0
  if (m < 0) m = 0
  if (m > 59) m = 59
  menitBackup.value = m
}

const updateTimer = () => {
  if (!setting.value.auto_backup) { 
    nextBackupTime.value = '' 
    return 
  }
  
  const now = new Date()
  let next = new Date(now)
  next.setHours(jamBackup.value, menitBackup.value, 0, 0)
  
  let prefix = 'Hari ini'
  if (next <= now) {
    next.setDate(next.getDate() + 1)
    prefix = 'Besok'
  }
  
  nextBackupTime.value = `${prefix} ${next.toLocaleTimeString('id-ID', { hour:'2-digit', minute:'2-digit' })}`
}

// ========== SCHEDULE NOTIFICATION (SETTIMEOUT) ==========
const scheduleNextBackup = () => {
  const now = new Date()
  const target = new Date()
  target.setHours(jamBackup.value, menitBackup.value, 0, 0)
  
  if (target <= now) {
    target.setDate(target.getDate() + 1)
  }
  
  const msUntilTarget = target.getTime() - now.getTime()
  
  timeoutId.value = setTimeout(() => {
    // TAMPILKAN NOTIFIKASI
    Swal.fire({
      icon: 'success',
      title: 'Auto Backup Berhasil!',
      text: `Backup otomatis untuk semua instansi berhasil dilakukan.`,
      timer: 3000,
      showConfirmButton: false
    })
    
    // REFRESH RIWAYAT OTOMATIS (PILIH INSTANSI PERTAMA JIKA BELUM ADA)
    if (instansiList.value.length > 0) {
      const instansi = instansiList.value[0]
      selectedInstansi.value = instansi
      pilihInstansi(instansi)
    }
    
    scheduleNextBackup()
  }, msUntilTarget)
}

// ========== API CALLS ==========
const fetchSetting = async () => {
  try {
    const res = await api.get('/super/backup-setting')
    if (res.data.data) {
      setting.value = {
        auto_backup: !!res.data.data.auto_backup,
        retensi_hari: parseInt(res.data.data.retensi_hari) || 60
      }
      jamBackup.value = res.data.data.jam_backup ?? 2
      menitBackup.value = res.data.data.menit_backup ?? 0
      updateTimer()
    }
  } catch (err) { 
    console.error('Fetch setting error:', err)
  }
}

const simpanSetting = async () => {
  validateJam()
  validateMenit()
  savingSetting.value = true
  
  try {
    await api.put('/super/backup-setting', {
      jam_backup: jamBackup.value,
      menit_backup: menitBackup.value,
      retensi_hari: setting.value.retensi_hari,
      auto_backup: setting.value.auto_backup
    })
    
    // HAPUS SEMUA KEY BACKUP DI LOCALSTORAGE
    Object.keys(localStorage).forEach(key => {
      if (key.startsWith('backup_done_')) {
        localStorage.removeItem(key)
      }
    })
    
    Swal.fire({ icon:'success', title:'Konfigurasi Tersimpan', timer:2000, showConfirmButton:false })
    showSetting.value = false
    updateTimer()
    scheduleNextBackup()
  } catch (err) { 
    console.error('Save setting error:', err)
    Swal.fire('Error', 'Gagal menyimpan pengaturan', 'error') 
  } finally { 
    savingSetting.value = false 
  }
}

const closeSetting = () => {
  showSetting.value = false
  // Reset ke nilai semula (biar gak berubah kalau batal)
  fetchSetting()
}

const fetchInstansi = async () => {
  try { 
    const res = await api.get('/super/instansi')
    instansiList.value = res.data.data 
  } catch (err) { 
    console.error(err) 
  }
}

const pilihInstansi = async (instansi) => {
  selectedInstansi.value = instansi
  loading.value = true
  try { 
    const res = await api.get(`/super/instansi/${instansi.id}/backups`)
    backupHistory.value = res.data.data 
  } catch (err) { 
    console.error(err) 
  } finally { 
    loading.value = false 
  }
}

const backupSekarang = async () => {
  if (!selectedInstansi.value) return
  
  backupLoading.value = true
  try { 
    await api.post(`/super/instansi/${selectedInstansi.value.id}/backup`)
    Swal.fire({ icon:'success', title:'Backup Berhasil', timer:2000, showConfirmButton:false })
    await pilihInstansi(selectedInstansi.value)
  } catch (err) { 
    Swal.fire('Error', 'Gagal backup data', 'error') 
  } finally { 
    backupLoading.value = false 
  }
}

const pulihkanBackup = async (backup) => {
  const result = await Swal.fire({
    title: 'Konfirmasi Pemulihan', 
    icon: 'warning', 
    showCancelButton: true, 
    confirmButtonColor: '#d97706',
    html: `<div style="text-align:left"><p>Kembalikan ke:</p><strong>${backup.nama_backup}</strong><p style="color:#dc2626;margin-top:12px">⚠️ Semua data saat ini akan DIGANTI!</p></div>`
  })
  
  if (result.isConfirmed) {
    try { 
      await api.post(`/super/instansi/restore/${backup.id}`)
      Swal.fire({ icon:'success', title:'Data Dipulihkan', timer:2000, showConfirmButton:false })
      await pilihInstansi(selectedInstansi.value)
    } catch (err) { 
      Swal.fire('Error', 'Gagal restore data', 'error') 
    }
  }
}

const downloadBackup = async (backup) => {
  try {
    // Buka link download di tab baru
    window.open(
      `http://localhost:3000/api/super/instansi/${selectedInstansi.value.id}/backups/download/${backup.id}`, 
      '_blank'
    )
  } catch (err) {
    console.error('Download error:', err)
    Swal.fire('Error', 'Gagal mendownload backup', 'error')
  }
}

const fileInput = ref(null)

const triggerImport = () => {
  fileInput.value.click()
}

const handleImportFile = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // Validasi ekstensi
  if (!file.name.endsWith('.json')) {
    Swal.fire('Error', 'File harus berformat .json', 'error')
    return
  }

  // Baca file
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const content = JSON.parse(e.target.result)
      
      // Validasi struktur JSON
      if (!content.siswa || !content.nilai || !content.mapel) {
        Swal.fire('Error', 'Format JSON tidak valid. Harus memiliki key: siswa, nilai, mapel', 'error')
        return
      }
      
      // Konfirmasi
      const result = await Swal.fire({
        title: 'Import Backup',
        text: `Data akan menggantikan semua data ${selectedInstansi.value.nama_instansi}. Lanjutkan?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#10B981',
        confirmButtonText: 'Ya, Import!'
      })
      
      if (result.isConfirmed) {
        // Kirim ke backend
        await api.post(`/super/instansi/${selectedInstansi.value.id}/import`, content)
        Swal.fire('Berhasil!', 'Data berhasil di-import', 'success')
        await pilihInstansi(selectedInstansi.value)
      }
    } catch (err) {
      Swal.fire('Error', 'File JSON tidak valid: ' + err.message, 'error')
    }
  }
  reader.readAsText(file)
}

// ========== LIFECYCLE ==========
onMounted(async () => {
  await fetchInstansi()
  await fetchSetting()
  scheduleNextBackup()
})

onUnmounted(() => { 
  if (timeoutId.value) {
    clearTimeout(timeoutId.value)
  }
})
</script>