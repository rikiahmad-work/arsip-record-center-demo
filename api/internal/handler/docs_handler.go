package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DocsHandler menyajikan spesifikasi OpenAPI 3.0.3 JSON dan halaman Swagger UI interaktif.
type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>SiARSIP API Documentation — Kementerian PPPA</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" />
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #fafbfc;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    }
    .topbar-header {
      background: linear-gradient(135deg, #0d6efd 0%, #0b5ed7 100%);
      color: #ffffff;
      padding: 14px 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    .topbar-title {
      font-size: 18px;
      font-weight: 700;
      letter-spacing: 0.3px;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .topbar-badge {
      font-size: 11px;
      font-weight: 600;
      background: rgba(255, 255, 255, 0.2);
      padding: 3px 8px;
      border-radius: 12px;
      letter-spacing: 0.5px;
    }
    .swagger-ui .topbar {
      display: none;
    }
    .swagger-ui .info {
      margin: 25px 0;
    }
    .swagger-ui .info .title {
      color: #212529;
      font-size: 28px;
    }
  </style>
</head>
<body>
  <header class="topbar-header">
    <div class="topbar-title">
      <span>SiARSIP Record Center</span>
      <span class="topbar-badge">OpenAPI 3.0</span>
    </div>
    <div style="font-size: 13px; opacity: 0.9;">
      Kementerian Pemberdayaan Perempuan dan Perlindungan Anak RI
    </div>
  </header>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: '/api/docs/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        persistAuthorization: true,
        defaultModelsExpandDepth: -1,
        docExpansion: 'list'
      });
    };
  </script>
</body>
</html>`

const openAPISpecJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "SiARSIP Record Center API",
    "version": "1.0.0",
    "description": "Dokumentasi resmi RESTful API Sistem Informasi Pengelolaan Arsip Dinamis Inaktif (Record Center) Kementerian Pemberdayaan Perempuan dan Perlindungan Anak Republik Indonesia. Dilengkapi autentikasi Multi-Factor Authentication (MFA / TOTP) dan otorisasi Role-Based Access Control (RBAC)."
  },
  "servers": [
    { "url": "http://localhost:8000", "description": "Server Pengembangan Lokal" }
  ],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "Masukkan Access Token JWT (Format: Bearer <token>)"
      }
    }
  },
  "security": [{ "bearerAuth": [] }],
  "paths": {
    "/api/auth/login": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Login pengguna (NIP / Username & Password)",
        "description": "Autentikasi tahap pertama. Mengembalikan token JWT langsung atau status mfa_required jika MFA aktif.",
        "security": [],
        "responses": {
          "200": { "description": "Login berhasil (JWT diterbitkan) atau tantangan MFA" },
          "401": { "description": "Kredensial salah / akun dinonaktifkan" },
          "429": { "description": "Batas percobaan login terlampaui (Rate Limit)" }
        }
      }
    },
    "/api/auth/mfa/verify": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Verifikasi TOTP / Recovery Code saat login",
        "description": "Autentikasi tahap kedua menggunakan kode 6-digit Authenticator atau 8-karakter Recovery Code.",
        "security": [],
        "responses": {
          "200": { "description": "Verifikasi MFA berhasil, token sesi penuh diterbitkan" },
          "401": { "description": "Kode verifikasi salah atau kedaluwarsa" },
          "429": { "description": "Batas percobaan MFA terlampaui" }
        }
      }
    },
    "/api/auth/mfa/status": {
      "get": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Cek status aktivasi MFA akun pengguna aktif",
        "responses": {
          "200": { "description": "Status MFA (aktif/nonaktif, sisa recovery codes)" }
        }
      }
    },
    "/api/auth/mfa/setup": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Inisialisasi konfigurasi MFA baru (Secret & QR URI)",
        "responses": {
          "200": { "description": "Secret key TOTP dan URI QR code untuk aplikasi Authenticator" }
        }
      }
    },
    "/api/auth/mfa/enable": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Aktivasi MFA akun setelah konfirmasi kode 6-digit pertama",
        "responses": {
          "200": { "description": "MFA berhasil diaktifkan dan 10 kode pemulihan darurat diterbitkan" }
        }
      }
    },
    "/api/auth/mfa/disable": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Nonaktifkan MFA akun pengguna aktif",
        "responses": {
          "200": { "description": "MFA berhasil dinonaktifkan" },
          "401": { "description": "Password konfirmasi tidak valid" }
        }
      }
    },
    "/api/auth/mfa/recovery-codes": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Buat ulang daftar kode pemulihan darurat (Recovery Codes)",
        "responses": {
          "200": { "description": "Daftar 10 recovery codes baru" }
        }
      }
    },
    "/api/auth/change-password": {
      "put": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Ubah kata sandi akun lokal pengguna aktif",
        "responses": {
          "200": { "description": "Kata sandi berhasil diperbarui" },
          "400": { "description": "Password lama salah atau konfirmasi password tidak sesuai" }
        }
      }
    },
    "/api/auth/profile": {
      "get": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Informasi profil dan identitas pegawai aktif",
        "responses": {
          "200": { "description": "Data profil pengguna dan hak akses" }
        }
      }
    },
    "/api/auth/logout": {
      "post": {
        "tags": ["Autentikasi & Akun"],
        "summary": "Logout & pencabutan (blacklist) token sesi JWT",
        "responses": {
          "200": { "description": "Sesi berhasil diakhiri" }
        }
      }
    },
    "/api/users": {
      "get": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Daftar pengguna terdaftar (Administrator)",
        "responses": {
          "200": { "description": "Daftar pengguna terpaginasi" }
        }
      },
      "post": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Tambah akun pengguna baru (Administrator)",
        "responses": {
          "201": { "description": "Akun pengguna berhasil dibuat" }
        }
      }
    },
    "/api/users/stats": {
      "get": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Statistik ringkasan pengguna per peran",
        "responses": {
          "200": { "description": "Statistik jumlah pengguna" }
        }
      }
    },
    "/api/users/{id}": {
      "get": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Detail akun pengguna",
        "responses": {
          "200": { "description": "Data detail pengguna" }
        }
      },
      "put": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Perbarui data pengguna (Administrator)",
        "responses": {
          "200": { "description": "Data pengguna berhasil diperbarui" }
        }
      },
      "delete": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Nonaktifkan pengguna (Soft delete)",
        "responses": {
          "200": { "description": "Pengguna berhasil dinonaktifkan" }
        }
      }
    },
    "/api/users/{id}/restore": {
      "put": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Pulihkan akun pengguna yang dinonaktifkan",
        "responses": {
          "200": { "description": "Akun pengguna berhasil dipulihkan" }
        }
      }
    },
    "/api/roles": {
      "get": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Daftar peran sistem dan matriks izin hak akses (RBAC)",
        "responses": {
          "200": { "description": "Daftar peran dan permission matrix" }
        }
      }
    },
    "/api/roles/{id}/permissions": {
      "put": {
        "tags": ["Manajemen Pengguna & Hak Akses"],
        "summary": "Perbarui alokasi hak akses untuk peran tertentu (Administrator)",
        "responses": {
          "200": { "description": "Izin hak akses berhasil diperbarui" }
        }
      }
    },
    "/api/arsip": {
      "get": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Daftar arsip inaktif (Filter pencarian, klasifikasi, SKKAD, kurun waktu)",
        "responses": {
          "200": { "description": "Daftar arsip terpaginasi (20 kolom)" }
        }
      },
      "post": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Registrasi arsip baru (Auto-fill JRA & SKKAD)",
        "responses": {
          "201": { "description": "Arsip berhasil didaftarkan" },
          "422": { "description": "Kode klasifikasi tidak ditemukan" }
        }
      }
    },
    "/api/arsip/stats": {
      "get": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Statistik arsip berdasarkan kategori dan tingkat keamanan SKKAD",
        "responses": {
          "200": { "description": "Statistik arsip" }
        }
      }
    },
    "/api/arsip/export": {
      "get": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Ekspor data berkas arsip ke CSV/Excel",
        "responses": {
          "200": { "description": "File CSV data arsip" }
        }
      }
    },
    "/api/arsip/import": {
      "post": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Impor massal berkas arsip dari file CSV/Excel",
        "responses": {
          "200": { "description": "Laporan hasil impor berkas arsip" }
        }
      }
    },
    "/api/arsip/{id}": {
      "get": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Detail lengkap berkas arsip beserta dokumen pendukung",
        "responses": {
          "200": { "description": "Detail berkas arsip" }
        }
      },
      "put": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Perbarui data berkas arsip (Auto-fill JRA)",
        "responses": {
          "200": { "description": "Arsip berhasil diperbarui" }
        }
      },
      "delete": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Hapus berkas arsip (Administrator)",
        "responses": {
          "200": { "description": "Arsip berhasil dihapus" }
        }
      }
    },
    "/api/arsip/{id}/restore": {
      "put": {
        "tags": ["Pengelolaan Arsip"],
        "summary": "Pulihkan berkas arsip yang terhapus (Administrator)",
        "responses": {
          "200": { "description": "Arsip berhasil dipulihkan" }
        }
      }
    },
    "/api/arsip/{id}/dokumen": {
      "get": {
        "tags": ["Dokumen Digital"],
        "summary": "Daftar berkas digitalisasi arsip",
        "responses": {
          "200": { "description": "Daftar file terlampir" }
        }
      },
      "post": {
        "tags": ["Dokumen Digital"],
        "summary": "Unggah berkas digitalisasi fisik (PDF/JPG/PNG)",
        "description": "Mendukung validasi magic bytes, batas ukuran berkas, dan sanitasi nama file.",
        "responses": {
          "201": { "description": "Berkas berhasil diunggah" },
          "422": { "description": "Format MIME atau ukuran file tidak diizinkan" }
        }
      }
    },
    "/api/dokumen/{id}/download": {
      "get": {
        "tags": ["Dokumen Digital"],
        "summary": "Unduh atau preview dokumen digital",
        "description": "Dilengkapi verifikasi izin akses terhadap klasifikasi keamanan SKKAD Rahasia/Sangat Rahasia.",
        "responses": {
          "200": { "description": "Stream file dokumen digital" },
          "403": { "description": "Akses dokumen ditolak (Hak akses tidak mencukupi)" }
        }
      }
    },
    "/api/dokumen/{id}": {
      "delete": {
        "tags": ["Dokumen Digital"],
        "summary": "Hapus berkas dokumen digital",
        "responses": {
          "200": { "description": "Dokumen berhasil dihapus" }
        }
      }
    },
    "/api/master/kode-klasifikasi": {
      "get": {
        "tags": ["Master Data"],
        "summary": "Daftar 435 kode klasifikasi arsip, retensi JRA, dan SKKAD",
        "responses": {
          "200": { "description": "Daftar master kode klasifikasi" }
        }
      },
      "post": {
        "tags": ["Master Data"],
        "summary": "Tambah kode klasifikasi arsip baru (Administrator)",
        "responses": {
          "201": { "description": "Kode klasifikasi berhasil ditambahkan" }
        }
      }
    },
    "/api/master/lokasi-simpan": {
      "get": {
        "tags": ["Master Data"],
        "summary": "Daftar lokasi fisik Record Center (Depot, Rak, Boks, Folder)",
        "responses": {
          "200": { "description": "Daftar lokasi simpan" }
        }
      },
      "post": {
        "tags": ["Master Data"],
        "summary": "Tambah lokasi penyimpanan fisik baru",
        "responses": {
          "201": { "description": "Lokasi simpan berhasil ditambahkan" }
        }
      }
    },
    "/api/master/unit-kerja": {
      "get": {
        "tags": ["Master Data"],
        "summary": "Daftar unit pengolah / pencipta arsip Kemen PPPA",
        "responses": {
          "200": { "description": "Daftar unit kerja" }
        }
      },
      "post": {
        "tags": ["Master Data"],
        "summary": "Tambah unit kerja baru",
        "responses": {
          "201": { "description": "Unit kerja berhasil ditambahkan" }
        }
      }
    },
    "/api/master/konfigurasi": {
      "get": {
        "tags": ["Master Data"],
        "summary": "Ambil konfigurasi sistem dan parameter aplikasi (Administrator)",
        "responses": {
          "200": { "description": "Data konfigurasi sistem" }
        }
      },
      "put": {
        "tags": ["Master Data"],
        "summary": "Perbarui konfigurasi sistem (Administrator)",
        "responses": {
          "200": { "description": "Konfigurasi sistem berhasil diperbarui" }
        }
      }
    },
    "/api/peminjaman": {
      "get": {
        "tags": ["Peminjaman Arsip"],
        "summary": "Daftar transaksi peminjaman fisik (Auto-sync status Terlambat)",
        "responses": {
          "200": { "description": "Daftar peminjaman arsip" }
        }
      },
      "post": {
        "tags": ["Peminjaman Arsip"],
        "summary": "Catat peminjaman fisik arsip baru",
        "responses": {
          "201": { "description": "Transaksi peminjaman berhasil dicatat" }
        }
      }
    },
    "/api/peminjaman/{id}/kembali": {
      "put": {
        "tags": ["Peminjaman Arsip"],
        "summary": "Catat pengembalian fisik arsip ke Record Center",
        "responses": {
          "200": { "description": "Pengembalian arsip berhasil diverifikasi" }
        }
      }
    },
    "/api/penyusutan": {
      "get": {
        "tags": ["Penyusutan Arsip"],
        "summary": "Daftar usulan penyusutan berkas (Usul Musnah & Usul Serah)",
        "responses": {
          "200": { "description": "Daftar pengajuan penyusutan" }
        }
      },
      "post": {
        "tags": ["Penyusutan Arsip"],
        "summary": "Ajukan usulan penyusutan berkas arsip",
        "responses": {
          "201": { "description": "Pengajuan penyusutan berhasil dibuat" }
        }
      }
    },
    "/api/penyusutan/{id}/approve": {
      "put": {
        "tags": ["Penyusutan Arsip"],
        "summary": "Setujui usulan penyusutan (Approver / Administrator)",
        "responses": {
          "200": { "description": "Usulan penyusutan disetujui" },
          "403": { "description": "Otorisasi ditolak (Bukan Approver)" }
        }
      }
    },
    "/api/penyusutan/{id}/reject": {
      "put": {
        "tags": ["Penyusutan Arsip"],
        "summary": "Tolak usulan penyusutan dengan catatan verifikasi",
        "responses": {
          "200": { "description": "Usulan penyusutan ditolak" }
        }
      }
    },
    "/api/penyusutan/{id}/selesai": {
      "put": {
        "tags": ["Penyusutan Arsip"],
        "summary": "Tandai penyusutan selesai dilaksanakan (Administrator)",
        "responses": {
          "200": { "description": "Status penyusutan diselesaikan" }
        }
      }
    },
    "/api/dashboard/stats": {
      "get": {
        "tags": ["Dashboard & Pelaporan"],
        "summary": "Indikator ringkasan statistik Record Center",
        "responses": {
          "200": { "description": "Data statistik dashboard" }
        }
      }
    },
    "/api/dashboard/chart": {
      "get": {
        "tags": ["Dashboard & Pelaporan"],
        "summary": "Grafik tren penambahan berkas arsip bulanan",
        "responses": {
          "200": { "description": "Data tren bulanan" }
        }
      }
    },
    "/api/laporan/rekap": {
      "get": {
        "tags": ["Dashboard & Pelaporan"],
        "summary": "Rekapitulasi arsip per kategori, keamanan, kode, dan unit kerja",
        "responses": {
          "200": { "description": "Rekapitulasi statistik arsip" }
        }
      }
    },
    "/api/laporan/dai": {
      "get": {
        "tags": ["Dashboard & Pelaporan"],
        "summary": "Laporan Daftar Arsip Inaktif (DAI) 18 kolom format ANRI",
        "responses": {
          "200": { "description": "Tabel 18 kolom DAI standar ANRI" }
        }
      }
    },
    "/api/laporan/export": {
      "get": {
        "tags": ["Dashboard & Pelaporan"],
        "summary": "Ekspor Laporan DAI 18 kolom ke file CSV/Excel",
        "responses": {
          "200": { "description": "File unduhan CSV DAI" }
        }
      }
    },
    "/api/label/box": {
      "get": {
        "tags": ["Cetak Label Fisik"],
        "summary": "Cetak dokumen PDF Label Boks Arsip Record Center",
        "responses": {
          "200": { "description": "Stream file PDF Label Boks" }
        }
      }
    },
    "/api/label/folder": {
      "get": {
        "tags": ["Cetak Label Fisik"],
        "summary": "Cetak dokumen PDF Label Folder Berkas Arsip",
        "responses": {
          "200": { "description": "Stream file PDF Label Folder" }
        }
      }
    },
    "/api/label/inaktif": {
      "get": {
        "tags": ["Cetak Label Fisik"],
        "summary": "Cetak dokumen PDF Label Berkas Arsip Inaktif",
        "responses": {
          "200": { "description": "Stream file PDF Label Inaktif" }
        }
      }
    }
  }
}`

func (h *DocsHandler) ServeSwaggerUI(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
}

func (h *DocsHandler) ServeOpenAPISpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(openAPISpecJSON))
}
