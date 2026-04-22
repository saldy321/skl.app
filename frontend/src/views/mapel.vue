<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center px-2">
      <h1 class="text-2xl font-bold text-slate-800">Master Mata Pelajaran</h1>
      <button @click="openModal()" class="px-5 py-2.5 bg-emerald-500 hover:bg-emerald-600 text-white rounded-xl font-bold">
        <i class="fas fa-plus"></i> Tambah Mapel
      </button>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden mx-2">
      <table class="w-full text-left">
        <thead class="bg-slate-50 border-b">
          <tr>
            <th class="px-6 py-4 text-xs font-black text-slate-400 uppercase">Nama Mapel</th>
            <th class="px-6 py-4 text-xs font-black text-slate-400 uppercase">Kelompok</th>
            <th class="px-6 py-4 text-xs font-black text-slate-400 uppercase text-center">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in mapels" :key="m.id" class="hover:bg-slate-50 border-b border-slate-100">
            <td class="px-6 py-4 font-semibold text-slate-700">{{ m.nama_mapel }}</td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ m.kelompok }}</td>
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2">
                <button @click="openModal(m)" class="text-blue-500 p-2 hover:bg-blue-50 rounded-lg">
                  <i class="fas fa-edit"></i>
                </button>
                <button @click="deleteMapel(m.id)" class="text-red-400 p-2 hover:bg-red-50 rounded-lg">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="mapels.length === 0">
            <td colspan="3" class="px-6 py-10 text-center text-slate-400 italic">Belum ada data.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 z-[9999]">
      <div class="bg-white rounded-3xl max-w-md w-full p-8 shadow-2xl">
        <h3 class="text-xl font-bold mb-6">{{ isEdit ? 'Update Mapel' : 'Tambah Mapel' }}</h3>
        <div class="space-y-4">
          <div>
            <label class="text-xs font-bold text-slate-400 uppercase">Nama Mapel</label>
            <input v-model="form.nama_mapel" type="text" class="w-full px-4 py-3 bg-slate-50 border rounded-xl outline-none focus:ring-2 focus:ring-emerald-500">
          </div>
          <div>
            <label class="text-xs font-bold text-slate-400 uppercase">Kelompok</label>
            <select v-model="form.kelompok" class="w-full px-4 py-3 bg-slate-50 border rounded-xl outline-none">
              <option value="A">Umum (A)</option>
              <option value="B">Umum (B)</option>
              <option value="Kejuruan">Kejuruan (C)</option>
            </select>
          </div>
        </div>
        <div class="mt-8 flex flex-col gap-2">
          <button @click="submitForm" class="w-full py-3 bg-emerald-500 text-white rounded-xl font-bold">
            {{ isEdit ? 'Simpan Perubahan' : 'Tambahkan Sekarang' }}
          </button>
          <button @click="showModal = false" class="w-full py-3 text-slate-400 font-semibold">Batal</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import Swal from 'sweetalert2'

const route = useRoute()
const slug = route.params.slug
const mapels = ref([])
const showModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ nama_mapel: '', kelompok: 'A' })

const fetchMapels = async () => {
  try {
    const res = await api.get(`/${slug}/admin/mapel`)
    mapels.value = res.data.data || []
  } catch (e) {
    console.error(e)
  }
}

const openModal = (data = null) => {
  if (data) {
    isEdit.value = true
    currentId.value = data.id 
    form.nama_mapel = data.nama_mapel
    form.kelompok = data.kelompok
  } else {
    isEdit.value = false
    currentId.value = null
    form.nama_mapel = ''
    form.kelompok = 'A'
  }
  showModal.value = true
}

const submitForm = async () => {
  if (!form.nama_mapel) return Swal.fire('Oops!', 'Nama mapel wajib diisi', 'warning')
  
  try {
    if (isEdit.value) {
      await api.put(`/${slug}/admin/mapel/${currentId.value}`, form)
    } else {
      await api.post(`/${slug}/admin/mapel`, form)
    }
    showModal.value = false
    await fetchMapels() // Refresh list data
    Swal.fire({ icon: 'success', title: 'Mantap!', timer: 1000, showConfirmButton: false })
  } catch (err) {
    Swal.fire('Error', 'Gagal memproses data', 'error')
  }
}

const deleteMapel = async (id) => {
  const { isConfirmed } = await Swal.fire({
    title: 'Hapus mapel ini?',
    text: "Data nilai yang terikat mungkin ikut bermasalah.",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#1e293b',
    confirmButtonText: 'Ya, Hapus!'
  })

  if (isConfirmed) {
    try {
      await api.delete(`/${slug}/admin/mapel/${id}`)
      await fetchMapels()
      Swal.fire('Terhapus!', '', 'success')
    } catch (e) {
      Swal.fire('Gagal!', 'Data sedang digunakan.', 'error')
    }
  }
}

onMounted(() => fetchMapels())
</script>