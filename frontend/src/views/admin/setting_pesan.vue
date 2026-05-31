<template>
  <div class="max-w-3xl mx-auto">
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 md:p-8">
      
      <!-- HEADER -->
      <div class="mb-8">
        <h2 class="text-xl font-bold text-slate-800 flex items-center gap-3">
          <i class="fas fa-envelope-open-text text-indigo-500"></i>
          Pengaturan Pesan Kelulusan
        </h2>
        <p class="text-sm text-slate-500 mt-2">Atur pesan yang muncul saat siswa membuka amplop digital</p>
      </div>

      <!-- FORM -->
      <div class="space-y-6">
        
        <!-- PESAN LULUS -->
        <div>
          <label class="flex items-center gap-2 font-bold text-slate-700 mb-2">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            Pesan untuk Siswa LULUS
          </label>
          <textarea 
            v-model="form.pesan_lulus"
            rows="4"
            class="w-full border border-slate-200 rounded-xl p-4 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            placeholder="Tulis pesan kelulusan..."
          ></textarea>
          <p class="text-xs text-slate-400 mt-1">
            <i class="fas fa-info-circle mr-1"></i>
            Gunakan <code class="bg-slate-100 px-1 rounded">{NAMA_INSTANSI}</code> untuk nama sekolah otomatis
          </p>
        </div>

        <!-- PESAN TIDAK LULUS -->
        <div>
          <label class="flex items-center gap-2 font-bold text-slate-700 mb-2">
            <span class="w-2 h-2 rounded-full bg-red-500"></span>
            Pesan untuk Siswa TIDAK LULUS
          </label>
          <textarea 
            v-model="form.pesan_tidak_lulus"
            rows="4"
            class="w-full border border-slate-200 rounded-xl p-4 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            placeholder="Tulis pesan untuk yang tidak lulus..."
          ></textarea>
          <p class="text-xs text-slate-400 mt-1">
            <i class="fas fa-info-circle mr-1"></i>
            Gunakan <code class="bg-slate-100 px-1 rounded">{NAMA_INSTANSI}</code> untuk nama sekolah otomatis
          </p>
        </div>

        <!-- PREVIEW -->
        <div class="bg-slate-50 rounded-xl p-5 border border-slate-100">
          <h3 class="font-bold text-slate-700 mb-4 flex items-center gap-2">
            <i class="fas fa-eye text-slate-500"></i>
            Preview
          </h3>
          
          <div class="space-y-4">
            <!-- Preview Lulus -->
            <div class="bg-emerald-50 border border-emerald-200 rounded-lg p-4">
              <p class="text-[10px] font-bold text-emerald-700 uppercase mb-1">✅ LULUS</p>
              <p class="text-sm text-slate-700 leading-relaxed">{{ previewLulus || '(Belum diisi)' }}</p>
            </div>
            
            <!-- Preview Tidak Lulus -->
            <div class="bg-red-50 border border-red-200 rounded-lg p-4">
              <p class="text-[10px] font-bold text-red-700 uppercase mb-1">❌ TIDAK LULUS</p>
              <p class="text-sm text-slate-700 leading-relaxed">{{ previewTidakLulus || '(Belum diisi)' }}</p>
            </div>
          </div>
        </div>

        <!-- TOMBOL -->
        <div class="flex justify-end pt-4 border-t border-slate-100">
          <button 
            @click="saveSettings"
            :disabled="loading"
            class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-300 text-white px-8 py-3 rounded-xl text-sm font-bold shadow-lg shadow-indigo-100 transition-all flex items-center gap-2"
          >
            <i v-if="loading" class="fas fa-spinner fa-spin"></i>
            <i v-else class="fas fa-save"></i>
            {{ loading ? 'Menyimpan...' : 'Simpan Pengaturan' }}
          </button>
        </div>

      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import Swal from 'sweetalert2'

const route = useRoute()
const slug = route.params.slug
const loading = ref(false)

const form = ref({
  pesan_lulus: '',
  pesan_tidak_lulus: ''
})

// Preview dengan nama sekolah dummy
const previewLulus = computed(() => {
  if (!form.value.pesan_lulus) return ''
  return form.value.pesan_lulus.replace('{NAMA_INSTANSI}', 'SMA Negeri 1 Contoh')
})

const previewTidakLulus = computed(() => {
  if (!form.value.pesan_tidak_lulus) return ''
  return form.value.pesan_tidak_lulus.replace('{NAMA_INSTANSI}', 'SMA Negeri 1 Contoh')
})

const loadSettings = async () => {
  try {
    const res = await api.get(`/${slug}/admin/setting/pesan`)
    if (res.data?.data) {
      form.value.pesan_lulus = res.data.data.pesan_lulus || ''
      form.value.pesan_tidak_lulus = res.data.data.pesan_tidak_lulus || ''
    }
  } catch (err) {
    console.error('Gagal load settings:', err)
  }
}

const saveSettings = async () => {
  loading.value = true
  try {
    const res = await api.put(`/${slug}/admin/setting/pesan`, form.value)
    Swal.fire('Berhasil!', 'Pesan kelulusan berhasil disimpan', 'success')
  } catch (err) {
    Swal.fire('Error', err.response?.data?.message || 'Gagal menyimpan', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(loadSettings)
</script>