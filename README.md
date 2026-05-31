# Aplikasi SKL Digital

Aplikasi SKL Digital adalah sistem pengelolaan Surat Keterangan Lulus (SKL) berbasis web dengan arsitektur multi-tenant. Platform ini memungkinkan satu instalasi pusat (Super Admin) mengelola banyak sekolah/instansi, sementara setiap sekolah memiliki Admin Sekolah yang mengelola data siswa, nilai, template SKL, dan pengumuman kelulusan. Siswa dapat mengakses portal kelulusan melalui URL unik sekolah (`/{slug}/login-siswa`), membuka amplop digital, dan mencetak atau menyimpan SKL dalam format PDF melalui fitur cetak browser.

Sistem dirancang untuk mendukung berbagai jenjang pendidikan (SD/MI, SMP/MTS, SMA/MA, SMK) dengan penyesuaian tampilan dan field data otomatis berdasarkan tingkat sekolah.

---

## Tech Stack Detail

### Backend

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v2.52.12-00ACD7?style=flat-square)
![GORM](https://img.shields.io/badge/GORM-v1.31.1-42B883?style=flat-square)
![MySQL](https://img.shields.io/badge/MySQL-Driver-4479A1?style=flat-square&logo=mysql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-v5.3.1-black?style=flat-square)

| Komponen | Library / Driver | Fungsi |
|---|---|---|
| Bahasa | Go 1.25 | Runtime backend |
| Web Framework | `github.com/gofiber/fiber/v2` | HTTP server, routing, middleware, CORS, static file |
| ORM | `gorm.io/gorm` | Abstraksi database dan relasi model |
| Database Driver | `gorm.io/driver/mysql`, `github.com/go-sql-driver/mysql` | Koneksi MySQL |
| Autentikasi | `github.com/golang-jwt/jwt/v5` | Token JWT (HS256) |
| Enkripsi Password | `golang.org/x/crypto` | Hashing dan verifikasi password (bcrypt) |
| Environment | `github.com/joho/godotenv` | Load variabel dari file `.env` |
| UUID | `github.com/google/uuid` | Primary key dan identitas entitas |
| Excel | `github.com/xuri/excelize/v2` | Import/export data siswa dan nilai |
| Email | `gopkg.in/gomail.v2`, SMTP native | Pengiriman OTP dan kredensial admin |
| Image Processing | `github.com/disintegration/imaging` (dependency) | Dukungan pemrosesan gambar background |

### Frontend

![Vue.js](https://img.shields.io/badge/Vue.js-3.5-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white)
![Pinia](https://img.shields.io/badge/Pinia-3.0-yellow?style=flat-square)
![Vite](https://img.shields.io/badge/Vite-7.3-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4.2-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)
![Axios](https://img.shields.io/badge/Axios-1.13-5A29E4?style=flat-square)

| Komponen | Library | Fungsi |
|---|---|---|
| Framework UI | `vue` ^3.5.30 | Komponen reaktif dan Single Page Application |
| State Management | `pinia` ^3.0.4 | Store autentikasi, sesi, dan filter tahun angkatan |
| Routing | `vue-router` ^5.0.3 | Navigasi halaman dan route guard berbasis role |
| Build Tool | `vite` ^7.3.1 | Development server dan production build |
| Vue Plugin | `@vitejs/plugin-vue` ^6.0.4 | Dukungan SFC Vue di Vite |
| DevTools | `vite-plugin-vue-devtools` ^8.0.7 | Debugging Vue di development |
| HTTP Client | `axios` ^1.13.6 | Komunikasi REST API dengan credentials |
| Cookie Helper | `js-cookie` ^3.0.5 | Manajemen cookie sisi klien |
| Styling | `tailwindcss` ^4.2.2, `@tailwindcss/postcss`, `autoprefixer`, `postcss` | Utility-first CSS |
| Icons | `@fortawesome/fontawesome-free` ^7.2.0 | Icon di seluruh antarmuka |
| Notifikasi | `sweetalert2` ^11.26.24 | Dialog konfirmasi, alert, dan notifikasi |
| Rich Text Editor | `@vueup/vue-quill` ^1.2.0 | Editor HTML untuk isi surat SKL |
| Spreadsheet Client | `xlsx` ^0.18.5 | Dukungan manipulasi file Excel di frontend |

---

## Breakdown Seluruh Fitur Aplikasi

### 1. Sistem Autentikasi dan Otorisasi

#### 1.1 Peran Pengguna (Role)

| Role | Kode | Akses |
|---|---|---|
| Super Admin | `super_admin` | Panel pusat: kelola instansi, admin sekolah, backup nasional, manajemen NISN darurat, impersonate admin |
| Admin Sekolah | `admin` | Panel sekolah: CRUD siswa, mapel, nilai, template SKL, pengaturan pengumuman |
| Siswa | `siswa` | Portal kelulusan: login, lihat status, cetak SKL |

Tidak terdapat role Guru terpisah dalam sistem ini. Seluruh operasional akademik di level sekolah ditangani oleh Admin Sekolah.

#### 1.2 Login Super Admin dan Admin Sekolah

- **Endpoint:** `POST /api/login`
- **Alur:**
  1. Pengguna memasukkan email dan password di halaman `/login`.
  2. Frontend memvalidasi slider captcha sebelum submit.
  3. Backend memverifikasi kredensial. Super Admin hardcoded dapat login langsung tanpa OTP (fallback development).
  4. Untuk Super Admin (database) dan Admin Sekolah: sistem generate OTP 6 digit, simpan ke database, kirim via email SMTP.
  5. Frontend menampilkan form OTP (`POST /api/verify-otp`).
  6. OTP valid selama 5 menit. Setelah verifikasi, JWT diterbitkan dan disimpan di cookie HttpOnly.
- **Response setelah login:** role, slug, nama instansi, tingkat sekolah, instansi_id.

#### 1.3 Login Siswa

- **Endpoint:** `POST /api/:slug/login-siswa`
- **Kredensial:** NISN + Tanggal Lahir (format sesuai data di database).
- **Validasi tenant:** Middleware `CheckTenant` memastikan slug URL valid dan siswa milik instansi tersebut.
- **Halaman frontend:** `/:slug/login-siswa`
- **Fitur tambahan di halaman login siswa:**
  - Pengecekan instansi via `GET /api/check-instansi/:slug`
  - Countdown timer jika pengumuman belum dibuka (`waktu_buka_pengumuman`)
  - Tampilan logo sekolah (opsional, dikontrol admin)
  - Form dinonaktifkan selama countdown aktif

#### 1.4 Keamanan JWT dan Cookie

- Token JWT ditandatangani dengan `JWT_SECRET` (HS256).
- Masa berlaku token: 24 jam.
- Cookie `token`: HttpOnly, Path `/`, SameSite `None`, MaxAge 86400 detik.
- Claims JWT: `id`, `role`, `tingkat_sekolah`, `nama_instansi`, `instansi_id`, `slug`, `exp`.
- Middleware `Auth`: validasi cookie di setiap request protected.
- Middleware `AuthSiswa`: validasi khusus role siswa di route portal.
- Middleware `Admin`: validasi role admin sekolah.
- Middleware `IsSuperAdmin`: validasi role super admin.
- Middleware `CheckTenant`: isolasi data per sekolah berdasarkan slug URL.

#### 1.5 Logout

- **Endpoint:** `GET /api/logout`
- Menghapus cookie token di server dan klien.
- Frontend Pinia store membersihkan localStorage dan redirect ke `/login`.

#### 1.6 Manajemen Sesi Frontend

- State autentikasi disimpan di Pinia + localStorage (`isLoggedIn`, `role`, `slug`, `nama_instansi`, `instansi_id`, `tingkat`).
- Timer sesi 23 jam: peringatan SweetAlert sebelum token expired.
- Interceptor Axios: auto-logout dan redirect jika response 401 (kecuali endpoint login).
- Route guard Vue Router: proteksi berdasarkan role, slug tenant, dan normalisasi slug ke lowercase.

#### 1.7 Profil Pengguna

- **Endpoint:** `GET /api/me`
- Mengembalikan profil Super Admin atau Admin Sekolah (email, role, instansi, foto profil, logo instansi).
- Upload foto profil admin: `POST /api/upload-foto` (maks. 2MB, JPG/PNG). Foto sekaligus dijadikan logo instansi.
- Hapus foto profil: `DELETE /api/delete-foto`

---

### 2. Modul Super Admin

#### 2.1 Dashboard Super Admin

- **Halaman:** `/super-admin`
- **Endpoint:** `GET /api/super/dashboard`
- **Statistik:** total sekolah (instansi), total admin sekolah, total siswa nasional.
- Menampilkan daftar instansi terbaru.

#### 2.2 Manajemen Instansi (CRUD Sekolah)

- **Halaman:** `/super-admin/instansi`, `/super-admin/instansi/tambah`, `/super-admin/instansi/edit/:id`

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/api/super/instansi` | Daftar sekolah baru sekaligus akun admin (RegisterFullSchool) |
| GET | `/api/super/instansi` | Daftar semua instansi |
| GET | `/api/super/instansi/:id` | Detail instansi by ID |
| PUT | `/api/super/instansi/:id` | Update data sekolah (nama, kode, slug, tingkat, alamat) |
| DELETE | `/api/super/instansi/:id` | Hard delete instansi beserta siswa, nilai, admin, mapel terkait |

- **Data instansi:** nama, kode instansi (unik), slug URL (unik), tingkat sekolah, alamat, logo, waktu buka pengumuman, toggle tampilkan logo.
- **Registrasi sekolah:** otomatis generate slug dari nama, hash password admin, kirim email kredensial via background goroutine.

#### 2.3 Manajemen Admin Sekolah

- **Halaman:** `/super-admin/manage`
- **Endpoint:**
  - `GET /api/super/admin-list` — daftar admin dengan relasi instansi
  - `POST /api/super/admin-register` — tambah admin baru ke instansi
  - `PUT /api/super/admin/:id/email` — ubah email admin
  - `POST /api/super/instansi/:instansi_id/reset-password` — reset password admin (min. 6 karakter)

#### 2.4 Impersonate (Akses Sebagai Admin Sekolah)

- **Endpoint:** `POST /api/super/impersonate/:instansi_id`
- Super Admin masuk ke dashboard sekolah dengan token role `admin`.
- Banner peringatan impersonate di layout admin.
- **Keluar impersonate:** `POST /api/restore-super-admin` — kembalikan sesi ke Super Admin.

#### 2.5 Manajemen NISN Darurat

- **Halaman:** `/super-admin/manajemen-nisn`
- **Endpoint:**
  - `GET /api/super/nisn/search?nisn=` — cari NISN lintas seluruh instansi (join siswa + instansi)
  - `POST /api/super/nisn/force-delete` — hapus permanen siswa beserta nilai (wajib isi alasan, audit log)

#### 2.6 Backup dan Restore Data Instansi

- **Halaman:** `/super-admin/monitoring`
- **Scheduler otomatis:** goroutine di `main.go` cek setiap 20 detik, backup semua instansi pada jam yang dikonfigurasi (WIB).

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/super/backup-setting` | Ambil konfigurasi auto backup |
| PUT | `/api/super/backup-setting` | Simpan konfigurasi (jam, menit, retensi hari, toggle auto) |
| POST | `/api/super/instansi/:id/backup` | Backup manual per instansi |
| GET | `/api/super/instansi/:id/backups` | Riwayat backup 60 hari (dikelompokkan per tanggal) |
| POST | `/api/super/instansi/restore/:backupId` | Pulihkan data dari snapshot backup |
| GET | `/api/super/instansi/:id/backups/download/:backupId` | Download file backup JSON |
| POST | `/api/super/instansi/:id/import` | Import backup JSON ke instansi |

- **Isi backup:** data siswa, nilai, dan mapel per instansi (JSON + file di folder `./backups/`).
- **Retensi:** backup lama dihapus otomatis (hard delete) sesuai `retensi_hari`.
- **Default setting:** jam 02:00 WIB, retensi 60 hari, auto backup aktif.

---

### 3. Modul Admin Sekolah

#### 3.1 Dashboard Admin

- **Halaman:** `/:slug/dashboard`
- **Endpoint:** `GET /api/:slug/dashboard`
- **Statistik:** total siswa, siswa lulus, siswa tidak lulus, total mapel, total nilai terinput.
- Kartu statistik dapat diklik untuk navigasi ke modul terkait.

#### 3.2 Filter Angkatan (Tahun Lulus)

- Selector angkatan di header layout admin (2024–2027).
- Disimpan di Pinia store (`selectedYear`) dan localStorage.
- Semua query data siswa difilter berdasarkan `tahun_lulus` query parameter.

#### 3.3 Manajemen Data Siswa (CRUD)

- **Halaman:** `/:slug/admin/siswa`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/siswa/` | Daftar siswa (filter `tahun_lulus`) |
| GET | `/api/:slug/admin/siswa/search?keyword=&tahun_lulus=` | Pencarian NISN atau nama siswa |
| GET | `/api/:slug/admin/siswa/:id` | Detail siswa dengan relasi nilai |
| POST | `/api/:slug/admin/siswa/` | Tambah siswa (form-data, opsional upload foto) |
| PUT | `/api/:slug/admin/siswa/:id` | Update siswa |
| DELETE | `/api/:slug/admin/siswa/:id` | Hapus permanen siswa + foto |
| DELETE | `/api/:slug/admin/siswa/all?tahun_lulus=` | Hapus massal semua siswa per angkatan |
| POST | `/api/:slug/admin/siswa/import` | Import massal dari Excel |

- **Field data siswa:** NISN (unik), nama, tempat/tanggal lahir, jenis kelamin (L/P), kelas, jurusan, tahun lulus, nama wali, status lulus, foto siswa.
- **Field jurusan:** ditampilkan otomatis untuk jenjang SMK, SMA, MA.
- **Form tambah/edit:** modal dengan validasi field wajib.
- **Pencarian real-time:** debounce input keyword NISN/nama.

#### 3.4 Import Data Siswa (Excel)

- **Endpoint:** `POST /api/:slug/admin/siswa/import`
- **Format kolom Excel (sheet pertama):** NISN, Nama, Tempat Lahir, Tanggal Lahir, Jenis Kelamin, Kelas, Jurusan, Tahun Lulus (opsional, default 2026).
- Duplikat NISN di-skip (`ON CONFLICT DO NOTHING`).
- Frontend: tombol Import Excel, accept `.xlsx` / `.xls`.

#### 3.5 Manajemen Mata Pelajaran (CRUD)

- **Halaman:** `/:slug/admin/mapel`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/mapel` | Daftar mapel instansi |
| POST | `/api/:slug/admin/mapel` | Tambah mapel |
| PUT | `/api/:slug/admin/mapel/:id` | Update mapel |
| DELETE | `/api/:slug/admin/mapel/:id` | Hapus permanen mapel |

- **Field mapel:** nama mapel, kelompok (A = Umum Wajib, B = Muatan Lokal, Kejuruan = Kelompok C untuk SMK).
- Slug mapel di-generate otomatis saat create.

#### 3.6 Manajemen Nilai

- **Halaman Input Nilai:** `/:slug/admin/nilai`
- **Halaman Leger:** `/:slug/admin/leger`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/nilai/filters` | Opsi filter kelas dan jurusan |
| GET | `/api/:slug/admin/nilai/filter-options` | Opsi filter kelas dan jurusan |
| GET | `/api/:slug/admin/nilai/template?kelas=&jurusan=` | Download template Excel nilai |
| POST | `/api/:slug/admin/nilai/import` | Import nilai dari Excel |
| GET | `/api/:slug/admin/nilai/leger?kelas=&jurusan=&tahun_lulus=` | Rekap leger nilai + rata-rata |

- **Template Excel nilai:** kolom ID_SISWA (hidden), NO, NISN, NAMA SISWA, plus kolom dinamis per mapel.
- **Import nilai:** upsert berdasarkan pasangan siswa_id + mapel_id, tahun ajaran default 2025/2026 semester 2.
- **Leger frontend:** filter kelas, jurusan, kelompok mapel (A/B/Kejuruan), highlight nilai di bawah 75.

#### 3.7 Status Kelulusan dan Proses Kelulusan Massal

- **Endpoint cek kelayakan:** `POST /api/:slug/admin/check-graduation`
  - Membaca KKM dari template SKL (`minimal_kelulusan`).
  - Mengecek setiap siswa belum lulus: rata-rata nilai >= KKM dan tidak ada nilai di bawah 60.
  - Mengembalikan daftar siswa layak (`layak_ids`) dan tidak layak beserta alasan.
- **Endpoint proses kelulusan:** `POST /api/:slug/admin/execute-promotion`
  - Mode `eligible_only`: update hanya siswa yang ID-nya layak.
  - Mode `force_all`: update semua siswa dengan status belum lulus.
  - Set `status_lulus = true` dan `tahun_lulus` sesuai input.
- **Frontend:** modal Proses Lulus di halaman siswa dengan preview hasil cek kelayakan.

#### 3.8 Manajemen Foto Siswa

- **Halaman:** `/:slug/admin/foto-siswa`

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/api/:slug/admin/upload-siswa` | Upload massal foto via ZIP (nama file = NISN) |
| DELETE | `/api/:slug/admin/siswa/:id/foto` | Hapus foto individual siswa |

- Upload ZIP diproses paralel (goroutine). Response: success, partial (sebagian gagal), atau error.
- Foto disimpan di `./public/uploads/siswa/`.
- Tabel preview foto per siswa dengan tombol hapus hover.

#### 3.9 Template SKL (Konfigurasi Surat)

- **Halaman:** `/:slug/admin/template-skl`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/template-skl?instansi_id=` | Ambil template SKL |
| POST | `/api/:slug/admin/template-skl` | Simpan/update template SKL |

- **Field template:**
  - Identitas surat: nama surat, nomor surat, tanggal surat
  - Kepala sekolah: nama, NIP
  - KKM: minimal kelulusan (float)
  - Kop surat: file header (upload gambar)
  - Isi surat: dasar surat, isi surat, penutup surat (rich text HTML via Quill Editor)
  - Stempel: gambar, lebar, toggle pakai stempel
  - Tanda tangan kepala sekolah: gambar, lebar, toggle pakai TTD
  - Margin top untuk presisi printer
  - Toggle: tampilkan nilai di halaman admin, tampilkan nilai di halaman siswa, gunakan kelompok mapel, sertakan foto siswa
  - Placeholder dinamis di isi surat: `{NAMA}`, `{NISN}`, `{KELAS}`, `{JURUSAN}`, `{TAHUN_LULUS}`, dll.


#### 3.11 Pengaturan Jadwal Pengumuman

- **Halaman:** `/:slug/admin/setting-waktu`
- **Endpoint:** `POST /api/:slug/admin/setting/waktu-buka`
- Set tanggal dan waktu pembukaan portal siswa (format `YYYY-MM-DD HH:MM:SS`, timezone Asia/Jakarta).
- Toggle tampilkan/sembunyikan logo di halaman login siswa: `POST /api/:slug/admin/setting/tampilkan-logo`

#### 3.12 Pengaturan Pesan Kelulusan

- **Halaman:** `/:slug/admin/setting-pesan`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/setting/pesan` | Ambil pesan (admin) |
| PUT | `/api/:slug/admin/setting/pesan` | Simpan pesan lulus dan tidak lulus |
| GET | `/api/:slug/portal/setting/pesan` | Ambil pesan (siswa/public) |

- Pesan kustom untuk siswa lulus dan tidak lulus saat membuka amplop digital.
- Placeholder `{NAMA_INSTANSI}` diganti otomatis di frontend.
- Preview pesan di halaman setting.

#### 3.13 Pengaturan Background Portal Siswa

- **Halaman:** `/:slug/admin/setting-background`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/:slug/admin/setting/background` | Ambil background |
| GET | `/api/:slug/portal/setting/background` | Ambil background (siswa) |
| POST | `/api/:slug/admin/setting/background/upload` | Upload background (JPG/JPEG/PNG/WEBP) |
| DELETE | `/api/:slug/admin/setting/background` | Hapus background |

- Background ditampilkan di halaman portal siswa sebelum membuka amplop.
- Default gradient jika belum diupload.

---

### 4. Modul Portal Siswa

#### 4.1 Dashboard / Portal SKL

- **Halaman:** `/:slug/portal/skl`
- **Endpoint:** `GET /api/:slug/portal/dashboard`
- Mengembalikan data siswa (dengan nilai dan mapel) + template SKL.

#### 4.2 Pengalaman Amplop Digital

- **State closed:** tampilan amplop tertutup, klik untuk membuka.
- **State loading:** animasi persiapan surat.
- **State lulus/tidak lulus:** kartu surat dengan pesan dinamis (dari setting pesan).
- Teks kelulusan disesuaikan tingkat sekolah (SD/MI, SMP/MTS, SMA/MA/SMK).
- Logo sekolah ditampilkan di amplop jika diaktifkan admin.

#### 4.3 Cetak dan Unduh SKL

- Tombol "Cetak SKL" menampilkan dokumen lengkap.
- Tombol "Cetak / Simpan PDF" memanggil dialog print browser (Save as PDF).
- CSS khusus `@media print` untuk layout A4, margin, dan penghapusan elemen UI.
- Konten SKL: kop surat, identitas siswa, isi surat HTML, tabel nilai (jika `tampilkan_nilai_siswa` aktif), foto siswa (jika `pakai_foto` aktif), stempel dan TTD.

#### 4.4 Validasi Akses Pengumuman

- Sebelum login, sistem cek `GET /api/check-instansi/:slug`.
- Jika `waktu_buka_pengumuman` belum lewat: countdown hari/jam/menit/detik, form login dinonaktifkan.
- Halaman error sekolah tidak ditemukan: `/:slug/not-found`.

---

### 5. Fitur Pendukung Lainnya

#### 5.1 Multi-Tenant Architecture

- Setiap sekolah memiliki slug unik sebagai segmen URL (`/:slug/...`).
- Data diisolasi per `instansi_id` di setiap query database.
- Super Admin dapat bypass tenant restriction.

#### 5.2 Static File Serving

- Backend melayani file upload di `/uploads` dari `./public/uploads/`.
- Subdirektori: `siswa/`, `admin/`, `instansi/`, `backgrounds/`.

#### 5.3 Penyesuaian Jenjang Sekolah

- Konstanta tingkat: SD, MI, SMP, MTS, SMA, MA, SMK, PUSAT.
- Field jurusan dan kelompok mapel kejuruan hanya aktif untuk jenjang SMK/SMA/MA.
- Teks pengumuman kelulusan disesuaikan per tingkat.

#### 5.4 Auto Migration Database

- Saat startup, GORM AutoMigrate menjalankan migrasi tabel: Admin, SuperAdmin, Instansi, TemplateSKL, Siswa, Mapel, Nilai, BackupInstansi, BackupSetting, Setting_pesan, SettingBackground.

#### 5.5 CORS

- Origin diizinkan: `http://localhost:5173`
- Credentials: true (cookie cross-origin).

#### 5.6 Logging Request

- Fiber logger middleware mencatat setiap HTTP request.

#### 5.7 Halaman Error

- `/admin-error` — akses ditolak (slug tidak cocok atau role tidak sesuai).
- `/:slug/not-found` — sekolah tidak terdaftar.

---

## Daftar Lengkap Endpoint API

### Public (Tanpa Auth)

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/login` | Login admin/super admin |
| POST | `/api/verify-otp` | Verifikasi OTP |
| POST | `/api/restore-super-admin` | Restore sesi super admin |
| GET | `/api/check-instansi/:slug` | Validasi slug sekolah + status pengumuman |

### Auth (Semua Role Login)

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/me` | Profil user |
| GET | `/api/logout` | Logout |
| POST | `/api/upload-foto` | Upload foto profil admin |
| DELETE | `/api/delete-foto` | Hapus foto profil admin |

### Super Admin (`/api/super/*`)

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/super/dashboard` | Statistik nasional |
| POST | `/api/super/instansi` | Daftar sekolah + admin |
| GET | `/api/super/instansi` | List instansi |
| GET | `/api/super/instansi/:id` | Detail instansi |
| PUT | `/api/super/instansi/:id` | Update instansi |
| DELETE | `/api/super/instansi/:id` | Hapus instansi |
| POST | `/api/super/admin-register` | Daftar admin sekolah |
| GET | `/api/super/admin-list` | List admin sekolah |
| POST | `/api/super/instansi/:instansi_id/reset-password` | Reset password admin |
| PUT | `/api/super/admin/:id/email` | Update email admin |
| POST | `/api/super/impersonate/:instansi_id` | Impersonate admin |
| GET | `/api/super/nisn/search` | Cari NISN nasional |
| POST | `/api/super/nisn/force-delete` | Hapus permanen NISN |
| POST | `/api/super/instansi/:id/backup` | Backup manual |
| GET | `/api/super/instansi/:id/backups` | Riwayat backup |
| POST | `/api/super/instansi/restore/:backupId` | Restore backup |
| GET | `/api/super/instansi/:id/backups/download/:backupId` | Download backup |
| POST | `/api/super/instansi/:id/import` | Import backup JSON |
| GET | `/api/super/backup-setting` | Setting auto backup |
| PUT | `/api/super/backup-setting` | Simpan setting backup |

### Tenant — Siswa

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/:slug/login-siswa` | Login siswa |
| GET | `/api/:slug/portal/dashboard` | Data SKL siswa |
| GET | `/api/:slug/portal/setting/pesan` | Pesan kelulusan (public) |
| GET | `/api/:slug/portal/setting/background` | Background portal (public) |

### Tenant — Admin

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/:slug/dashboard` | Dashboard admin |
| POST | `/api/:slug/admin/check-graduation` | Cek kelayakan lulus |
| POST | `/api/:slug/admin/execute-promotion` | Proses kelulusan massal |
| POST | `/api/:slug/admin/setting/waktu-buka` | Set jadwal pengumuman |
| POST | `/api/:slug/admin/setting/tampilkan-logo` | Toggle logo login siswa |
| GET/PUT | `/api/:slug/admin/setting/pesan` | CRUD pesan kelulusan |
| GET/POST/DELETE | `/api/:slug/admin/setting/background/*` | CRUD background |
| CRUD | `/api/:slug/admin/siswa/*` | Manajemen siswa + import |
| CRUD | `/api/:slug/admin/mapel/*` | Manajemen mapel |
| GET/POST | `/api/:slug/admin/template-skl` | Template SKL |
| GET | `/api/:slug/admin/cetak/:id` | Data cetak SKL |
| GET/POST | `/api/:slug/admin/nilai/*` | Template, import, leger nilai |
| POST/DELETE | `/api/:slug/admin/upload-siswa`, `/siswa/:id/foto` | Foto siswa |

---

## Daftar Halaman Frontend

| Path | Role | Deskripsi |
|---|---|---|
| `/login` | Public | Login admin dengan OTP |
| `/:slug/login-siswa` | Public | Login siswa + countdown |
| `/:slug/not-found` | Public | Sekolah tidak ditemukan |
| `/admin-error` | Public | Akses ditolak |
| `/super-admin` | Super Admin | Dashboard pusat |
| `/super-admin/instansi` | Super Admin | Daftar instansi |
| `/super-admin/instansi/tambah` | Super Admin | Tambah instansi |
| `/super-admin/instansi/edit/:id` | Super Admin | Edit instansi |
| `/super-admin/manage` | Super Admin | Manajemen admin sekolah |
| `/super-admin/manajemen-nisn` | Super Admin | Pencarian dan hapus NISN |
| `/super-admin/monitoring` | Super Admin | Backup dan restore |
| `/:slug/dashboard` | Admin | Dashboard sekolah |
| `/:slug/admin/mapel` | Admin | CRUD mata pelajaran |
| `/:slug/admin/siswa` | Admin | CRUD siswa + import + proses lulus |
| `/:slug/admin/foto-siswa` | Admin | Upload foto massal ZIP |
| `/:slug/admin/nilai` | Admin | Import nilai Excel |
| `/:slug/admin/leger` | Admin | Rekap leger nilai |
| `/:slug/admin/template-skl` | Admin | Konfigurasi template SKL |
| `/:slug/admin/setting-waktu` | Admin | Jadwal pengumuman + logo |
| `/:slug/admin/setting-pesan` | Admin | Pesan kelulusan |
| `/:slug/admin/setting-background` | Admin | Background portal siswa |
| `/:slug/admin/cetak/:id` | Admin | Preview dan cetak SKL |
| `/:slug/portal/skl` | Siswa | Portal amplop dan cetak SKL |

---

## Struktur Folder

```
project_v2/
├── backend/
│   ├── main.go                    # Entry point, CORS, static files, auto backup scheduler
│   ├── go.mod                     # Dependensi Go
│   ├── .env                       # Konfigurasi environment (DB, JWT, SMTP)
│   ├── backups/                   # File backup JSON instansi
│   ├── config/
│   │   └── database.go            # Koneksi MySQL, AutoMigrate
│   ├── routes/
│   │   └── api.go                 # Definisi seluruh route API
│   ├── middleware/
│   │   ├── auth.go                # JWT auth, tenant check, role guard
│   │   └── auth_siswa.go          # Auth khusus portal siswa
│   ├── models/
│   │   ├── admin.go               # Model admin sekolah
│   │   ├── super_admin.go         # Model super admin
│   │   ├── instansi.go            # Model sekolah/instansi
│   │   ├── siswa.go               # Model data siswa
│   │   ├── mapel.go               # Model mata pelajaran
│   │   ├── nilai.go               # Model nilai siswa
│   │   ├── template.go            # Model template SKL
│   │   ├── backup_instansi.go     # Model snapshot backup
│   │   ├── setting.go             # Model setting backup global
│   │   ├── setting_pesan.go       # Model pesan kelulusan
│   │   ├── setting_background.go  # Model background portal
│   │   └── const.go               # Konstanta role dan tingkat sekolah
│   ├── controllers/
│   │   ├── auth/
│   │   │   ├── login.go           # Login + OTP generation
│   │   │   ├── verify_otp.go      # Verifikasi OTP
│   │   │   ├── siswa_login.go     # Login siswa NISN + tanggal lahir
│   │   │   └── email.go           # Pengiriman email OTP
│   │   ├── admin/
│   │   │   ├── dashboard.go       # Statistik admin
│   │   │   ├── siswa.go           # CRUD siswa + import Excel
│   │   │   ├── mapel.go           # CRUD mapel
│   │   │   ├── nilai.go           # Template, import, leger nilai
│   │   │   ├── kelulusan.go       # Cek kelayakan + proses lulus massal
│   │   │   ├── skl.go             # Data cetak SKL admin
│   │   │   ├── template.go        # CRUD template SKL
│   │   │   ├── foto_siswa.go      # Upload ZIP foto massal
│   │   │   ├── profile.go         # Upload/hapus foto profil admin
│   │   │   ├── setting.go         # Jadwal pengumuman + toggle logo
│   │   │   ├── setting_pesan.go   # Pesan kelulusan
│   │   │   └── setting_background.go  # Background portal siswa
│   │   ├── super_admin/
│   │   │   ├── dashboard.go       # Statistik nasional
│   │   │   ├── instansi.go        # CRUD instansi + impersonate
│   │   │   ├── admin_sekolah.go   # CRUD admin sekolah
│   │   │   ├── management_nisn.go # Pencarian + force delete NISN
│   │   │   ├── backup.go          # Backup, restore, download, import
│   │   │   ├── backup_setting.go  # Konfigurasi auto backup
│   │   │   └── restore.go         # Restore sesi super admin
│   │   ├── siswa/
│   │   │   └── dashboard.go       # Data SKL untuk portal siswa
│   │   ├── check_instansi.go      # Validasi slug + countdown pengumuman
│   │   ├── profile.go             # Endpoint GET /me
│   │   ├── logout.go              # Logout
│   │   ├── cookie.go              # Helper set cookie JWT
│   │   └── helper.go              # Utility controller
│   ├── utils/
│   │   ├── jwt.go                 # Generate dan parse JWT
│   │   ├── password.go            # Hash dan verify password
│   │   ├── email.go               # Kirim email kredensial
│   │   ├── converter.go           # Konversi utility
│   │   └── constil.go             # Konstanta utility key
│   └── public/
│       └── uploads/
│           ├── siswa/             # Foto siswa
│           ├── admin/             # Foto profil admin
│           ├── instansi/          # Logo sekolah
│           └── backgrounds/       # Background portal siswa
│
└── frontend/
    ├── package.json               # Dependensi Node.js
    ├── vite.config.js             # Konfigurasi Vite + alias @
    ├── index.html                 # Entry HTML
    ├── public/                    # Asset statis publik
    └── src/
        ├── main.js                # Bootstrap Vue + Pinia + Router
        ├── App.vue                  # Root component
        ├── style.css                # Global styles + Tailwind
        ├── api/
        │   ├── index.js             # Axios instance + interceptors
        │   └── dashboard.js         # API helper dashboard
        ├── router/
        │   └── index.js             # Route definitions + guards
        ├── stores/
        │   └── auth.js              # Pinia auth store
        ├── utils/
        │   └── constants.js         # Konstanta tingkat sekolah dan role
        ├── layouts/
        │   ├── AdmLayout.vue        # Layout admin sekolah (sidebar + header)
        │   ├── SuperLayout.vue      # Layout super admin
        │   └── SiswaLayout.vue      # Layout portal siswa
        ├── components/
        │   └── ErrorSekolah.vue     # Halaman error sekolah
        ├── assets/                  # Gambar (logo, amplop, dll.)
        └── views/
            ├── auth/
            │   ├── login.vue        # Login admin + OTP + captcha slider
            │   └── login_siswa.vue  # Login siswa + countdown
            ├── admin/
            │   ├── dashboard.vue    # Dashboard statistik
            │   ├── siswa.vue        # CRUD siswa + import + proses lulus
            │   ├── mapel.vue        # CRUD mata pelajaran
            │   ├── nilai.vue        # Import nilai Excel
            │   ├── leger.vue        # Rekap leger nilai
            │   ├── foto_siswa.vue   # Upload foto ZIP massal
            │   ├── TemplateSKL.vue  # Editor template SKL
            │   ├── adm_cetak.vue    # Preview dan cetak SKL (admin)
            │   ├── setting_waktu.vue    # Jadwal pengumuman
            │   ├── setting_pesan.vue    # Pesan kelulusan
            │   └── setting_background.vue  # Background portal
            ├── siswa/
            │   └── dashboard.vue    # Portal amplop + cetak SKL
            └── super_admin/
                ├── dashboard.vue    # Dashboard nasional
                ├── instansi.vue     # Daftar instansi
                ├── add_instansi.vue # Form tambah instansi
                ├── edit_instansi.vue# Form edit instansi
                ├── manage.admin.vue # Manajemen admin sekolah
                ├── manajemen_nisn.vue  # Manajemen NISN darurat
                └── monitoring.vue   # Backup dan restore
```

---

## Panduan Instalasi dan Menjalankan Aplikasi

### Prasyarat

- **Go** 1.25 atau lebih baru
- **Node.js** ^20.19.0 atau >=22.12.0
- **MySQL** 8.x (atau MariaDB kompatibel)
- **Git**

### 1. Clone Repository

```bash
git clone <url-repository>
cd project_v2
```

### 2. Setup Backend

#### 2.1 Masuk ke direktori backend

```bash
cd backend
```

#### 2.2 Install dependensi Go

```bash
go mod tidy
```

#### 2.3 Buat database MySQL

Buat database kosong di MySQL, contoh:

```sql
CREATE DATABASE skl_digital CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 2.4 Konfigurasi environment

Salin atau buat file `.env` di folder `backend/`:

```env
DB_URL="user:password@tcp(127.0.0.1:3306)/skl_digital?charset=utf8mb4&parseTime=True&loc=Local"
JWT_SECRET="ganti_dengan_secret_key_yang_kuat"
SMTP_EMAIL=email_pengirim@gmail.com
SMTP_PASSWORD=app_password_smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

| Variabel | Deskripsi |
|---|---|
| `DB_URL` | DSN koneksi MySQL (user, password, host, port, nama database) |
| `JWT_SECRET` | Secret key untuk signing token JWT |
| `SMTP_EMAIL` | Email pengirim OTP |
| `SMTP_PASSWORD` | Password atau App Password SMTP |
| `SMTP_HOST` | Host server SMTP |
| `SMTP_PORT` | Port SMTP (587 untuk TLS) |

#### 2.5 Jalankan backend

```bash
go run main.go
```

Server backend berjalan di `http://localhost:3000`.

Saat pertama kali dijalankan, GORM akan otomatis membuat/migrate seluruh tabel database. Folder `./public/uploads/` dan `./backups/` akan digunakan untuk penyimpanan file.

### 3. Setup Frontend

#### 3.1 Masuk ke direktori frontend

```bash
cd ../frontend
```

#### 3.2 Install dependensi Node.js

```bash
npm install
```

#### 3.3 Jalankan development server

```bash
npm run dev
```

Frontend berjalan di `http://localhost:5173` (default Vite).

Pastikan `baseURL` di `frontend/src/api/index.js` mengarah ke backend:

```
http://localhost:3000/api
```

#### 3.4 Build production (opsional)

```bash
npm run build
npm run preview
```

### 4. Akses Aplikasi

| URL | Keterangan |
|---|---|
| `http://localhost:5173/login` | Login Admin / Super Admin |
| `http://localhost:5173/super-admin` | Panel Super Admin (setelah login) |
| `http://localhost:5173/{slug}/dashboard` | Dashboard Admin Sekolah |
| `http://localhost:5173/{slug}/login-siswa` | Portal login siswa |
| `http://localhost:5173/{slug}/portal/skl` | Portal cetak SKL siswa |

Ganti `{slug}` dengan slug instansi sekolah yang terdaftar (contoh: `smk-negeri-1`).

### 5. Urutan Menjalankan (Development)

1. Pastikan MySQL berjalan.
2. Jalankan backend: `cd backend && go run main.go`
3. Jalankan frontend: `cd frontend && npm run dev`
4. Buka browser ke `http://localhost:5173/login`

---

## Catatan 

- Backend module name: `skl-bakcend` (sesuai `go.mod`).
- CORS backend hanya mengizinkan origin `http://localhost:5173`. Sesuaikan jika deploy ke domain lain.
- Cookie `Secure` saat ini `false` (development). Set `true` saat production dengan HTTPS.
- Pencarian tracking SKL siswa menggunakan **NISN** dan **Tanggal Lahir**, bukan NIS.
- Auto backup berjalan sebagai background goroutine; pastikan backend tetap berjalan agar scheduler aktif.

---

