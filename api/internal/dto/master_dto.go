package dto

import "time"

// ================================================================
// 1. Kode Klasifikasi DTO
// ================================================================

type KodeKlasifikasiListRequest struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	PerPage  int    `form:"per_page"  binding:"omitempty,min=1,max=500"`
	Search   string `form:"search"`
	Q        string `form:"q"`
	SKKAD    string `form:"skkad"     binding:"omitempty,oneof=Biasa Terbatas Rahasia"`
	IsActive *bool  `form:"is_active"`
}

type StoreKodeKlasifikasiRequest struct {
	Kode                string  `json:"kode"                  binding:"required,max=20"`
	Nama                string  `json:"nama"                  binding:"required,max=255"`
	JangkaSimpanInaktif *string `json:"jangka_simpan_inaktif" binding:"omitempty,max=20"`
	NasibAkhir          *string `json:"nasib_akhir"           binding:"omitempty,oneof=Musnah Permanen"`
	SKKAD               string  `json:"skkad"                 binding:"required,oneof=Biasa Terbatas Rahasia"`
	IsActive            *bool   `json:"is_active"`
}

type UpdateKodeKlasifikasiRequest struct {
	Kode                string  `json:"kode"                  binding:"required,max=20"`
	Nama                string  `json:"nama"                  binding:"required,max=255"`
	JangkaSimpanInaktif *string `json:"jangka_simpan_inaktif" binding:"omitempty,max=20"`
	NasibAkhir          *string `json:"nasib_akhir"           binding:"omitempty,oneof=Musnah Permanen"`
	SKKAD               string  `json:"skkad"                 binding:"required,oneof=Biasa Terbatas Rahasia"`
	IsActive            *bool   `json:"is_active"`
}

type KodeKlasifikasiResponse struct {
	ID                  uint      `json:"id"`
	Kode                string    `json:"kode"`
	Nama                string    `json:"nama"`
	JangkaSimpanInaktif *string   `json:"jangka_simpan_inaktif"`
	NasibAkhir          *string   `json:"nasib_akhir"`
	SKKAD               string    `json:"skkad"`
	TotalArsip          int64     `json:"total_arsip"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ================================================================
// 2. Lokasi Simpan DTO
// ================================================================

type LokasiSimpanListRequest struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	PerPage  int    `form:"per_page"  binding:"omitempty,min=1,max=500"`
	Search   string `form:"search"`
	Q        string `form:"q"`
	IsActive *bool  `form:"is_active"`
}

type StoreLokasiSimpanRequest struct {
	Nama      string  `json:"nama"      binding:"required,max=255"`
	Gedung    *string `json:"gedung"    binding:"omitempty,max=255"`
	Lantai    *string `json:"lantai"    binding:"omitempty,max=20"`
	NoRak     *string `json:"no_rak"    binding:"omitempty,max=50"`
	NoBoks    *string `json:"no_boks"   binding:"omitempty,max=50"`
	NoFolder  *string `json:"no_folder" binding:"omitempty,max=50"`
	Kapasitas *int    `json:"kapasitas"`
	Terpakai  *int    `json:"terpakai"`
	IsActive  *bool   `json:"is_active"`
}

type UpdateLokasiSimpanRequest struct {
	Nama      string  `json:"nama"      binding:"required,max=255"`
	Gedung    *string `json:"gedung"    binding:"omitempty,max=255"`
	Lantai    *string `json:"lantai"    binding:"omitempty,max=20"`
	NoRak     *string `json:"no_rak"    binding:"omitempty,max=50"`
	NoBoks    *string `json:"no_boks"   binding:"omitempty,max=50"`
	NoFolder  *string `json:"no_folder" binding:"omitempty,max=50"`
	Kapasitas *int    `json:"kapasitas"`
	Terpakai  *int    `json:"terpakai"`
	IsActive  *bool   `json:"is_active"`
}

type LokasiSimpanResponse struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama"`
	Gedung    *string   `json:"gedung"`
	Lantai    *string   `json:"lantai"`
	NoRak     *string   `json:"no_rak"`
	NoBoks    *string   `json:"no_boks"`
	NoFolder  *string   `json:"no_folder"`
	Kapasitas *int      `json:"kapasitas"`
	Terpakai  int       `json:"terpakai"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================================================
// 3. Unit Kerja DTO
// ================================================================

type UnitKerjaListRequest struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	PerPage  int    `form:"per_page"  binding:"omitempty,min=1,max=500"`
	Search   string `form:"search"`
	Q        string `form:"q"`
	IsActive *bool  `form:"is_active"`
}

type StoreUnitKerjaRequest struct {
	Kode     string  `json:"kode"      binding:"required,max=20"`
	Nama     string  `json:"nama"      binding:"required,max=255"`
	Eselon   *string `json:"eselon"    binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

type UpdateUnitKerjaRequest struct {
	Kode     string  `json:"kode"      binding:"required,max=20"`
	Nama     string  `json:"nama"      binding:"required,max=255"`
	Eselon   *string `json:"eselon"    binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active"`
}

type UnitKerjaResponse struct {
	ID         uint      `json:"id"`
	Kode       string    `json:"kode"`
	Nama       string    `json:"nama"`
	Eselon     *string   `json:"eselon"`
	TotalArsip int64     `json:"total_arsip"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ================================================================
// 4. Konfigurasi Sistem DTO
// ================================================================

type UpdateKonfigurasiRequest struct {
	AppName                  string            `json:"app_name"`
	DefaultJangkaSimpan      string            `json:"default_jangka_simpan"`
	DefaultNasibAkhir        string            `json:"default_nasib_akhir"`
	NotifikasiJatuhTempoHari string            `json:"notifikasi_jatuh_tempo_hari"`
	LogoAplikasi             string            `json:"logo_aplikasi"`
	Items                    map[string]string `json:"items"`
}
