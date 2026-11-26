# 📁 Arsip Record Center

> **Sistem Informasi Manajemen Arsip**  
> **Kementerian Pemberdayaan Perempuan dan Perlindungan Anak Republik Indonesia (Kemen PPPA)**  
> Sistem pengelolaan dan pendataan daur hidup arsip berbasis web secara terintegrasi, akuntabel, dan aman.

---

## 📌 Daftar Isi

1. [Tentang Sistem](#-tentang-sistem)
2. [Fitur Utama](#-fitur-utama)
3. [Arsitektur & Tech Stack](#-arsitektur--tech-stack)
4. [Struktur Repositori](#-struktur-repositori)
5. [Prasyarat Sistem (Prerequisites)](#-prasyarat-sistem-prerequisites)
6. [Panduan Instalasi & Menjalankan](#-panduan-instalasi--menjalankan)
   - [A. Setup Database](#a-setup-database)
   - [B. Setup Backend (API)](#b-setup-backend-api)
   - [C. Setup Frontend (UI)](#c-setup-frontend-ui)
7. [Akun Pengguna Default (Seeder)](#-akun-pengguna-default-seeder)
8. [Matriks Peran & Hak Akses (RBAC)](#-matriks-peran--hak-akses-rbac)
9. [Daftar Endpoint Utama (API)](#-daftar-endpoint-utama-api)
10. [Perintah Pengembangan (CLI & Makefile)](#-perintah-pengembangan-cli--makefile)
11. [Lisensi & Kebijakan](#-lisensi--kebijakan)

---

## 📖 Tentang Sistem

**Arsip Record Center** adalah aplikasi terintegrasi yang dirancang untuk mendukung operasional pengelolaan kearsipan di lingkungan Kementerian Pemberdayaan Perempuan dan Perlindungan Anak Republik Indonesia (Kemen PPPA). Sistem ini mendigitalisasi pengelolaan Daftar Arsip Inaktif (DAI), penataan fisik di gudang/depot arsip, pencatatan sirkulasi peminjaman, serta penyusutan arsip (pemusnahan dan penyerahan) secara akuntabel dan aman.

Sistem menerapkan prinsip **Role-Based Access Control (RBAC)** dan **Segregation of Duties (SoD)** yang dilengkapi dengan perlindungan autentikasi multi-faktor (MFA / TOTP) serta jejak audit (*audit trail*).

---

## ✨ Fitur Utama

- **Pengelolaan Daftar Arsip Inaktif (DAI):**
  - Standar metadata 20 kolom informasi kearsipan.
  - Klasifikasi arsip otomatis berbasis kode klasifikasi kearsipan.
  - Auto-fill Jadwal Retensi Arsip (JRA): menghitung masa aktif, inaktif, nasib akhir, dan tingkat keamanan secara otomatis.
  - Manajemen lokasi fisik simpan bertingkat (Depot, Rak, Baris, Boks, Folder).
  - Integrasi QR Code untuk setiap berkas dan boks arsip.
- **Sirkulasi & Operasional Arsip:**
  - Peminjaman dan pengembalian berkas fisik dengan pelacakan jatuh tempo dan status transaksi.
  - Penyusutan arsip (usul musnah dan usul serah) berjenjang dengan berita acara.
- **Label & Barcode Generator:**
  - Cetak label identitas boks dan berkas berstandar fisik dengan kode QR unik.
- **Laporan & Rekapitulasi:**
  - Export rekapitulasi Daftar Arsip Inaktif (DAI) ke format spreadsheet/Excel.
  - Statistik volume arsip, grafik retensi, dan pemanfaatan kapasitas ruang penyimpanan.
- **Keamanan & Tata Kelola Pengguna:**
  - RBAC dengan 4 peran kanonik (`administrator`, `operator`, `approver`, `viewer`).
  - Autentikasi JWT (Access Token & Refresh Token) dengan mekanisme *Blacklisted Tokens*.
  - Autentikasi Dua Faktor (TOTP/MFA via Authenticator App).
  - Fitur *Step-Up Verification* untuk reset MFA darurat oleh Administrator.
  - Proteksi Rate Limiting & Anti-Brute Force.

---

## 🛠️ Arsitektur & Tech Stack

Proyek ini dibangun menggunakan arsitektur monorepo terpisah antara API dan UI:

```
┌─────────────────────────────────┐       HTTP / REST (JSON)       ┌─────────────────────────────────┐
│          FRONTEND (UI)          │ ◄────────────────────────────► │          BACKEND (API)          │
│  Nuxt 4 + Vue 3 + Tailwind CSS  │                                │  Golang 1.25 + Gin + GORM       │
│  Port: 3000                     │                                │  Port: 8000                     │
└─────────────────────────────────┘                                └─────────────────────────────────┘
                                                                                    │
                                                                                    ▼
                                                                   ┌─────────────────────────────────┐
                                                                   │            DATABASE             │
                                                                   │        MySQL 8.0 / MariaDB      │
                                                                   └─────────────────────────────────┘
```

### 1. Backend (`api/`)
- **Bahasa:** Go (Golang) versi 1.25+
- **HTTP Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/) (MySQL Driver)
- **Autentikasi & Keamanan:** `golang-jwt/jwt/v5`, `golang.org/x/crypto` (Bcrypt/Argon2), AES-256 GCM untuk secret MFA.
- **Pola Arsitektur:** Clean Layered Architecture (`handler` ➔ `service` ➔ `repository` ➔ `model` / `dto`).

### 2. Frontend (`ui/`)
- **Framework:** [Nuxt 4](https://nuxt.com/) (Vue 3, Composition API, TypeScript)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **Icons:** [Lucide Vue Next](https://lucide.dev/)
- **HTTP Client:** Axios dengan interceptor token refresh otomatis
- **QR Code:** `qrcode` library
- **Pola Arsitektur:** 3-Tier Pattern (`pages/` thin routing ➔ `modules/` feature implementation ➔ `components/` reusable UI).

---

## 📂 Struktur Repositori

```
arsip-record-center/
├── README.md               # Dokumentasi utama proyek
├── .gitignore              # Konfigurasi git ignore global
│
├── api/                    # 🚀 Backend Service (Golang)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go     # Entry point aplikasi & CLI flags
│   ├── internal/
│   │   ├── config/         # Konfigurasi environment (.env loader)
│   │   ├── database/       # Inisialisasi DB, migrasi, dan seeders
│   │   ├── dto/            # Data Transfer Object (Request/Response)
│   │   ├── handler/        # HTTP Handlers (Controller layer)
│   │   ├── middleware/     # Auth, RBAC, Security, Audit log, Rate Limit
│   │   ├── model/          # GORM Entities / Database Schema
│   │   ├── repository/     # Data Access Layer
│   │   ├── router/         # Rute API & middleware bindings
│   │   └── service/        # Business Logic Layer & Unit Tests
│   ├── pkg/                # Reusable utilities (Helper, Hash, dsb.)
│   ├── storage/            # Tempat penyimpanan dokumen fisik/lampiran
│   ├── .env.example        # Template konfigurasi environment backend
│   ├── Makefile            # Task runner backend (run, build, test, seed)
│   ├── go.mod              # Definisi modul & dependencies Golang
│   └── go.sum              # Checksum dependencies Golang
│
└── ui/                     # 🎨 Frontend Web App (Nuxt 4)
    ├── app/
    │   ├── assets/         # Static styles & Tailwind CSS base
    │   ├── components/     # Reusable UI (template, form primitives, c-modal, dll.)
    │   ├── composables/    # Vue composables (useArsip, useSession, useToast)
    │   ├── constants/      # Kamus data, menu navigasi, opsi enum
    │   ├── layouts/        # Layout template (default, login, error)
    │   ├── middleware/     # Route navigation guards (auth, role)
    │   ├── modules/        # Implementasi fitur (Grid, Form, Detail per fitur)
    │   ├── pages/          # Thin routing Nuxt (memanggil module dengan RBAC)
    │   ├── plugins/        # Axios interceptor, client plugins
    │   ├── types/          # TypeScript interface & type definitions
    │   └── utils/          # Helper date, formatter, dll.
    ├── public/             # Static public assets (logo, favicon)
    ├── .env.example        # Template konfigurasi environment frontend
    ├── nuxt.config.ts      # Konfigurasi Nuxt 4
    ├── package.json        # Dependencies & NPM scripts
    ├── tailwind.config.js  # Konfigurasi tema Tailwind CSS
    └── tsconfig.json       # Konfigurasi TypeScript
```

---

## 📋 Prasyarat Sistem (Prerequisites)

Pastikan dependensi berikut sudah terpasang di sistem operasi Anda:

- **Go:** Versi 1.25.0 atau lebih tinggi ([Download Go](https://go.dev/dl/))
- **Node.js:** Versi 18.x atau 20.x LTS ([Download Node.js](https://nodejs.org/))
- **Package Manager:** `yarn` (direkomendasikan) atau `npm`
- **Database:** MySQL 8.0+ atau MariaDB 10.5+
- **Git:** Versi terbaru

---

## 🚀 Panduan Instalasi & Menjalankan

### A. Setup Database

Buat database baru pada server MySQL Anda:

```sql
CREATE DATABASE arsip_record_center CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

### B. Setup Backend (API)

1. Masuk ke direktori `api`:
   ```bash
   cd api
   ```

2. Salin template `.env.example` ke `.env`:
   ```bash
   cp .env.example .env
   ```

3. Sesuaikan konfigurasi di file `.env`:
   ```env
   APP_NAME=RecordCenter-API
   APP_PORT=8000

   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_DATABASE=arsip_record_center
   DB_USERNAME=root
   DB_PASSWORD=password_db_anda

   JWT_SECRET=ganti-dengan-kunci-rahasia-minimal-32-karakter-acak
   JWT_ACCESS_EXPIRATION_HOURS=7
   JWT_REFRESH_EXPIRATION_HOURS=30

   CORS_ALLOWED_ORIGINS=http://localhost:3000

   STORAGE_PATH=./storage/dokumen
   MAX_UPLOAD_MB=10

   MFA_REQUIRED=false
   MFA_ISSUER=RecordCenter
   MFA_ENCRYPTION_KEY=ganti-dengan-32-byte-hex-atau-base64-key
   ```

4. Unduh dependencies Golang:
   ```bash
   go mod download
   # atau
   make deps
   ```

5. Jalankan migrasi dan seeder database awal:
   ```bash
   # Menjalankan migrasi database
   go run ./cmd/server/main.go --migrate

   # Menjalankan seeder (Kode Klasifikasi, Master Pendukung, dan 4 User Default)
   go run ./cmd/server/main.go --seed
   ```
   *Atau jika menggunakan `make`:*
   ```bash
   make migrate
   make seed
   ```

6. Jalankan server API backend:
   ```bash
   make run
   # atau
   go run ./cmd/server/main.go
   ```
   API akan berjalan di `http://localhost:8000`.

---

### C. Setup Frontend (UI)

1. Buka terminal baru dan masuk ke direktori `ui`:
   ```bash
   cd ui
   ```

2. Salin file konfigurasi environment:
   ```bash
   cp .env.example .env
   ```

3. Pastikan konfigurasi URL API mengarah ke backend:
   ```env
   NUXT_PUBLIC_API_BASE_URL=http://localhost:8000/api
   ```

4. Install dependencies Node.js:
   ```bash
   yarn install
   # atau
   npm install
   ```

5. Jalankan development server Nuxt:
   ```bash
   yarn dev
   # atau
   npm run dev
   ```
   Aplikasi web dapat diakses melalui browser di `http://localhost:3000`.

---

## 👥 Akun Pengguna Default (Seeder)

Setelah menjalankan perintah seeder (`--seed`), sistem menyediakan 4 akun default untuk masing-masing peran:

| No | Role | Username | Password Default | NIP | Keterangan |
|---|---|---|---|---|---|
| 1 | **Administrator** | `administrator` | `superadmin123` | `198501012010011001` | Administrator Sistem (Kelola user, role, master data, audit) |
| 2 | **Operator** | `operator1` | `adminarsip123` | `199002022010012002` | Operator Record Center (Input DAI, boks, peminjaman, cetak label) |
| 3 | **Approver** | `approver1` | `approver123` | `198804042010121004` | Pejabat Penilai / Pimpinan (Persetujuan peminjaman & usul penyusutan) |
| 4 | **Viewer** | `viewer1` | `viewer123` | `199503032015031003` | Unit Kerja / Pengolah (Pencarian dan katalog arsip) |

> ⚠️ **Catatan Keamanan:** Segera ubah kata sandi default dan aktifkan Multi-Factor Authentication (MFA) pada saat implementasi di lingkungan staging maupun produksi.

---

## 🛡️ Matriks Peran & Hak Akses (RBAC)

| Modul / Fitur | Administrator | Operator | Approver | Viewer |
|---|:---:|:---:|:---:|:---:|
| **Dashboard Statistik** | ✅ Full | ✅ Full | ✅ Full | 👁️ Terbatas |
| **Pencarian & Detail Arsip** | ✅ | ✅ | ✅ | ✅ |
| **Input / Edit / Hapus Arsip** | ✅ | ✅ | ❌ | ❌ |
| **Cetak Label & QR Code** | ✅ | ✅ | ❌ | ❌ |
| **Pengajuan Peminjaman** | ✅ | ✅ | ✅ | ✅ |
| **Persetujuan Peminjaman** | ✅ | ❌ | ✅ | ❌ |
| **Pengembalian Berkas** | ✅ | ✅ | ❌ | ❌ |
| **Usul Penyusutan Arsip** | ✅ | ✅ | ❌ | ❌ |
| **Persetujuan Penyusutan** | ✅ | ❌ | ✅ | ❌ |
| **Laporan & Rekap DAI** | ✅ Export | ✅ Export | ✅ Export | 👁️ Lihat |
| **Kelola Master Data** | ✅ Full | 👁️ Lihat | 👁️ Lihat | ❌ |
| **Kelola Pengguna & Role** | ✅ Full | ❌ | ❌ | ❌ |
| **Reset MFA & Audit Log** | ✅ Full | ❌ | ❌ | ❌ |

---

## 📡 Daftar Endpoint Utama (API)

Base URL: `http://localhost:8000/api`

### Autentikasi & Profil
- `POST /auth/login` — Login pengguna (mendapatkan JWT access & refresh token).
- `POST /auth/refresh` — Refresh access token yang telah kedaluwarsa.
- `POST /auth/logout` — Logout dan mendaftarkan token ke blacklist.
- `GET /auth/profile` — Mendapatkan profil pengguna yang sedang login.
- `POST /auth/change-password` — Mengubah kata sandi akun sendiri.

### Autentikasi Dua Faktor (MFA)
- `POST /mfa/setup` — Inisialisasi setup MFA (mendapatkan URI QR Code & Secret).
- `POST /mfa/enable` — Mengaktifkan MFA setelah verifikasi token pertama kali.
- `POST /mfa/verify` — Memverifikasi kode TOTP saat login.
- `POST /mfa/disable` — Menonaktifkan MFA.
- `POST /mfa/admin-reset` — Reset paksa MFA user oleh Administrator (Step-Up).

### Pengelolaan Arsip (DAI)
- `GET /arsip` — Menampilkan daftar arsip (pagination, filter unit, filter klasifikasi, search).
- `POST /arsip` — Menambahkan berkas arsip baru (auto-fill JRA & SKKAD).
- `GET /arsip/:id` — Detail lengkap berkas arsip (20 kolom metadata).
- `PUT /arsip/:id` — Memperbarui data berkas arsip.
- `DELETE /arsip/:id` — Soft-delete berkas arsip.
- `POST /arsip/import` — Import data arsip massal dari template format standar.
- `GET /arsip/export` — Export data arsip ke format Excel/Spreadsheet.

### Operasional (Peminjaman & Penyusutan)
- `GET /peminjaman` — Daftar transaksi peminjaman arsip.
- `POST /peminjaman` — Pengajuan peminjaman arsip baru.
- `PUT /peminjaman/:id/approve` — Persetujuan peminjaman oleh Approver.
- `PUT /peminjaman/:id/return` — Konfirmasi pengembalian berkas fisik.
- `GET /penyusutan` — Daftar usulan penyusutan arsip (Pemusnahan/Penyerahan).
- `POST /penyusutan` — Pengajuan usulan penyusutan baru.
- `PUT /penyusutan/:id/status` — Persetujuan/eksekusi status penyusutan.

### Master Data & Pengaturan
- `GET /master/kode-klasifikasi` — Daftar kode klasifikasi kearsipan.
- `GET /master/unit-kerja` — Daftar unit pengolah / pencipta arsip.
- `GET /master/lokasi-simpan` — Master data depot, rak, baris, dan boks.
- `GET /users` — Daftar akun pengguna (khusus Administrator).
- `GET /roles` — Matriks konfigurasi peran dan permission.

---

## 💻 Perintah Pengembangan (CLI & Makefile)

### Backend (`api/`)
```bash
make help        # Menampilkan seluruh perintah make yang tersedia
make run         # Menjalankan backend server development
make build       # Kompilasi binary Go ke ./bin/
make test        # Menjalankan unit test dengan race detector dan coverage
make test-cover  # Menghasilkan visualisasi HTML coverage test
make lint        # Menjalankan golangci-lint dan go vet
make fmt         # Memformat kode program menggunakan gofmt
make migrate     # Menjalankan migrasi database
make seed        # Menjalankan seeder data default
make clean       # Membersihkan file binary dan artefak build
```

### Frontend (`ui/`)
```bash
yarn dev         # Menjalankan development server (Hot Module Replacement)
yarn build       # Membangun aplikasi untuk produksi (output di .output/)
yarn preview     # Meninjau hasil build produksi secara lokal
yarn generate    # Melakukan static site generation
```

---

## 🔒 Lisensi & Kebijakan

Dokumen dan kode sumber ini merupakan properti internal **Kementerian Pemberdayaan Perempuan dan Perlindungan Anak Republik Indonesia (Kemen PPPA)**. Hak cipta dilindungi undang-undang. Dilarang mendistribusikan, menyalin, atau mempublikasikan kode tanpa izin tertulis dari pihak berwenang.
