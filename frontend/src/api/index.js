import axios from 'axios';
import Cookies from 'js-cookie';
import { useAuthStore } from '@/stores/auth';
import Swal from 'sweetalert2';

const api = axios.create({
    baseURL: 'http://localhost:3000/api',
    withCredentials: true,
});

api.interceptors.request.use(
    (config) => {
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

api.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error) => {
      
        const url = error.config?.url || '';
        const isLoginEndpoint = url === '/login' || 
                                url === '/verify-otp' || 
                                url.includes('/login-siswa');  // ← TAMBAHKAN INI
        
        // Jika error 401 dan BUKAN endpoint login, baru trigger logout
        if (error.response && error.response.status === 401 && !isLoginEndpoint) {
            const authStore = useAuthStore();

            if (authStore.isSessionEnding) {
                authStore.logout();
                return Promise.reject(error);
            }
            if (!Swal.isVisible()) {
                Swal.fire({
                    icon: 'warning',
                    title: 'Sesi Berakhir',
                    text: 'Sesi anda telah habis, silahkan login kembali',
                    timer: 2000,
                    showConfirmButton: false,
                    backdrop: true,
                    allowOutsideClick: false,
                });
            }
            authStore.logout();
            window.location.href = '/login';
        }
        
        // Untuk error 401 di endpoint login, biarkan saja biar pesan error dari backend tampil
        return Promise.reject(error);
    }
);

export default api;