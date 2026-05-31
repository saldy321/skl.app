import { defineStore } from 'pinia'
import api from '../api'
import Swal from 'sweetalert2'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        isLoggedIn: localStorage.getItem('isLoggedIn') === 'true',
        role: localStorage.getItem('role') || null,
        nama_instansi: localStorage.getItem('nama_instansi') || null,
        instansi_id: localStorage.getItem('instansi_id') || null,
        slug: localStorage.getItem('slug') || null,
        tingkat:localStorage.getItem('tingkat') || null,
        selectedYear: localStorage.getItem('selectedYear') || new Date().getFullYear().toString(),
        sessionTimerId: null,
        isSessionEnding: false
    }),
    actions: {
        async login(email, password) {
            try {
                const res = await api.post('/login', { email, password });
                
                if (res.data.status === 'otp_required') {
                    return { success: true, status: 'otp_required' };
                }
                
                this.setAuthData(res.data);
                return { success: true, ...res.data };
            } catch(err) {
                return {
                    success: false,
                    message: err.response?.data?.message || 'Login gagal!'
                };
            }
        },

        async verifyOTP(email, otp) {
            try {
                const res = await api.post('/verify-otp', { email, otp });
                this.setAuthData(res.data);
                return { success: true, ...res.data };
            } catch (err) {
                return {
                    success: false,
                    message: err.response?.data?.message || 'OTP Salah!'
                };
            }
        },

        setAuthData(data) {
            const role = data.role;
            const slug = data.slug;
            const nama_instansi = data.nama_instansi;
            const instansi_id = data.instansi_id;
            const tingkat = data.tingkat;

            this.isLoggedIn = true;
            this.role = role; 
            this.slug = slug;
            this.nama_instansi = nama_instansi;
            this.instansi_id = instansi_id;
            this.tingkat = tingkat || null;
            this.isSessionEnding = false;

            localStorage.setItem('isLoggedIn', 'true');
            localStorage.setItem('role', role || '');
            localStorage.setItem('slug', slug || '');
            localStorage.setItem('nama_instansi', nama_instansi || '');
            localStorage.setItem('instansi_id', instansi_id || '');
              localStorage.setItem('tingkat', tingkat || '');

            this.startSessionTimer();
        },

        async LoginSiswa(slug, nisn, tanggal_lahir) {
            try {
                const res = await api.post(`/${slug}/login-siswa`, { nisn, tanggal_lahir });
                
                console.log('[STORE] Login response:', res.data);
                
              
                if (res.data && res.data.status === 'success') {
       
                    const responseSlug = res.data.slug || slug;
                    
                    this.isLoggedIn = true;
                    this.role = 'siswa';
                    this.slug = responseSlug;
                    this.nama_instansi = res.data.nama_instansi || '';
                    this.instansi_id = res.data.instansi_id || '';
                    this.tingkat = res.data.tingkat || '';
                    
                    localStorage.setItem('isLoggedIn', 'true');
                    localStorage.setItem('role', 'siswa');
                    localStorage.setItem('slug', responseSlug);
                    localStorage.setItem('nama_instansi', this.nama_instansi);
                    localStorage.setItem('instansi_id', this.instansi_id);
                    localStorage.setItem('tingkat', this.tingkat);
                    
                    console.log('[STORE] Data tersimpan:', {
                        slug: this.slug,
                        role: this.role
                    });
                    
                    this.startSessionTimer();
                    
                    return { success: true, message: res.data.message, slug: responseSlug };
                }
                
                return { 
                    success: false, 
                    message: res.data?.message || 'Login gagal' 
                };
            } catch (err) {
                console.error('[STORE] Login error:', err);
                return { 
                    success: false, 
                    message: err.response?.data?.message || 'Data Siswa tidak ditemukan!' 
                };
            }
        },

        startSessionTimer() {
            if (this.sessionTimerId) {
                clearTimeout(this.sessionTimerId);
            }
            
            const timeToNotify = 23 * 60 * 60 * 1000;

            this.sessionTimerId = setTimeout(() => {
                this.isSessionEnding = true;

                Swal.fire({
                    icon: 'warning',
                    title: 'Sesi Hampir Berakhir',
                    text: 'Sesi Anda akan segera habis dalam 1 jam. Klik OK untuk login ulang.',
                    confirmButtonText: 'OK, Login Ulang',
                    allowOutsideClick: false
                }).then((result) => {
                    if (result.isConfirmed) {
                        this.logout(); 
                    }
                });
            }, timeToNotify); 
        },

        async checkAuth() {
  try {
    // Cek session ke backend
    const res = await api.get('/me')
    
    if (res.data && res.data.id) {
      // Session valid, update store
      this.isLoggedIn = true
      this.role = res.data.role
      this.nama_instansi = res.data.nama_instansi
      this.instansi_id = res.data.instansi_id
      this.slug = res.data.slug
      this.tingkat = res.data.tingkat
      
      // Update localStorage (sinkronisasi)
      localStorage.setItem('isLoggedIn', 'true')
      localStorage.setItem('role', res.data.role || '')
      localStorage.setItem('slug', res.data.slug || '')
      localStorage.setItem('nama_instansi', res.data.nama_instansi || '')
      localStorage.setItem('instansi_id', res.data.instansi_id || '')
      localStorage.setItem('tingkat', res.data.tingkat || '')
      
      return true
    } else {
      // Data tidak lengkap, logout
      this.logout()
      return false
    }
  } catch (err) {
    // Error (401 Unauthorized), logout
    this.logout()
    return false
  }
},

        setSelectedYear(year) {
            this.selectedYear = year.toString();
            localStorage.setItem('selectedYear', year.toString());
        },

        async logout() {
            if (this.sessionTimerId) {
                clearTimeout(this.sessionTimerId);
                this.sessionTimerId = null;
            }
            
            try { 
                await api.get('/logout'); 
            } catch(e) {}
            
            localStorage.clear();
            this.$reset();
            window.location.href = '/login';
        }
    }
})