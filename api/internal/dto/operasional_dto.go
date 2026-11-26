package dto

import "time"

// ================================================================
// 1. Peminjaman DTO (Nota Dinas #4)
// ================================================================

type PeminjamanListRequest struct {
	Page        int    `form:"page"          binding:"omitempty,min=1"`
	PerPage     int    `form:"per_page"      binding:"omitempty,min=1,max=500"`
	Q           string `form:"q"`
	Search      string `form:"search"`
	Status      string `form:"status"        binding:"omitempty,oneof=Dipinjam Dikembalikan Terlambat dipinjam dikembalikan terlambat"`
	UnitKerjaID uint   `form:"unit_kerja_id"`
	ArsipID     uint   `form:"arsip_id"`
}

func (r PeminjamanListRequest) ResolvedSearch() string {
	if r.Q != "" {
		return r.Q
	}
	return r.Search
}

type StorePeminjamanRequest struct {
	NoTiket               string  `json:"no_tiket"                binding:"omitempty,max=50"`
	ArsipID               uint    `json:"arsip_id"                binding:"required"`
	Peminjam              string  `json:"peminjam"                binding:"required,max=255"`
	UnitKerjaID           *uint   `json:"unit_kerja_id"`
	TanggalPinjam         string  `json:"tanggal_pinjam"          binding:"required"` // YYYY-MM-DD
	TanggalKembaliRencana string  `json:"tanggal_kembali_rencana" binding:"required"` // YYYY-MM-DD
	Catatan               *string `json:"catatan"`
}

type UpdatePeminjamanRequest struct {
	ArsipID               uint    `json:"arsip_id"                binding:"required"`
	Peminjam              string  `json:"peminjam"                binding:"required,max=255"`
	UnitKerjaID           *uint   `json:"unit_kerja_id"`
	TanggalPinjam         string  `json:"tanggal_pinjam"          binding:"required"`
	TanggalKembaliRencana string  `json:"tanggal_kembali_rencana" binding:"required"`
	Catatan               *string `json:"catatan"`
}

type KembaliPeminjamanRequest struct {
	TanggalKembaliAktual string  `json:"tanggal_kembali_aktual"` // YYYY-MM-DD (opsional, default hari ini)
	Catatan              *string `json:"catatan"`
}

type PeminjamanResponse struct {
	ID                    uint      `json:"id"`
	NoTiket               string    `json:"no_tiket"`
	ArsipID               uint      `json:"arsip_id"`
	NoItemArsip           string    `json:"no_item_arsip"`
	KodeKlasifikasi       string    `json:"kode_klasifikasi"`
	UraianArsip           string    `json:"uraian_arsip"`
	Peminjam              string    `json:"peminjam"`
	UnitKerjaID           *uint     `json:"unit_kerja_id"`
	UnitKerjaNama         string    `json:"unit_kerja_nama"`
	TanggalPinjam         string    `json:"tanggal_pinjam"`
	TanggalKembaliRencana string    `json:"tanggal_kembali_rencana"`
	TanggalKembaliAktual  *string   `json:"tanggal_kembali_aktual"`
	Status                string    `json:"status"`
	Catatan               *string   `json:"catatan"`
	CreatedBy             uint      `json:"created_by"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ================================================================
// 2. Penyusutan DTO (Usul Musnah & Usul Serah — Nota Dinas #5)
// ================================================================

type PenyusutanListRequest struct {
	Page      int    `form:"page"       binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page"   binding:"omitempty,min=1,max=500"`
	Q         string `form:"q"`
	Search    string `form:"search"`
	JenisUsul string `form:"jenis_usul" binding:"omitempty,oneof=Musnah Serah musnah serah"`
	Status    string `form:"status"     binding:"omitempty,oneof=Diajukan Disetujui Ditolak Selesai diajukan disetujui ditolak selesai"`
	ArsipID   uint   `form:"arsip_id"`
}

func (r PenyusutanListRequest) ResolvedSearch() string {
	if r.Q != "" {
		return r.Q
	}
	return r.Search
}

type StorePenyusutanRequest struct {
	NoPengajuan    string  `json:"no_pengajuan"    binding:"omitempty,max=50"`
	JenisUsul      string  `json:"jenis_usul"      binding:"required"` // Musnah / Serah
	ArsipID        uint    `json:"arsip_id"        binding:"required"`
	TanggalRetensi string  `json:"tanggal_retensi" binding:"required"` // YYYY-MM-DD
	Alasan         *string `json:"alasan"`
}

type DecisionPenyusutanRequest struct {
	CatatanApprover *string `json:"catatan_approver"`
}

type PenyusutanResponse struct {
	ID              uint      `json:"id"`
	NoPengajuan     string    `json:"no_pengajuan"`
	JenisUsul       string    `json:"jenis_usul"`
	ArsipID         uint      `json:"arsip_id"`
	NoItemArsip     string    `json:"no_item_arsip"`
	KodeKlasifikasi string    `json:"kode_klasifikasi"`
	UraianArsip     string    `json:"uraian_arsip"`
	NasibAkhir      string    `json:"nasib_akhir"`
	TanggalRetensi  string    `json:"tanggal_retensi"`
	Alasan          *string   `json:"alasan"`
	Status          string    `json:"status"`
	DiajukanOleh    uint      `json:"diajukan_oleh"`
	DiajukanNama    string    `json:"diajukan_nama"`
	DisetujuiOleh   *uint     `json:"disetujui_oleh"`
	DisetujuiNama   string    `json:"disetujui_nama"`
	CatatanApprover *string   `json:"catatan_approver"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ================================================================
// 3. Dokumen Arsip DTO (Digitalisasi — Nota Dinas #8)
// ================================================================

type DokumenArsipResponse struct {
	ID          uint      `json:"id"`
	ArsipID     uint      `json:"arsip_id"`
	NamaFile    string    `json:"nama_file"`
	MimeType    string    `json:"mime_type"`
	UkuranBytes int64     `json:"ukuran_bytes"`
	Keterangan  *string   `json:"keterangan"`
	DownloadURL string    `json:"download_url"`
	UploadedBy  uint      `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// ================================================================
// 4. Audit Log DTO
// ================================================================

type AuditLogListRequest struct {
	Page    int    `form:"page"     binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=200"`
	Aksi    string `form:"aksi"`
	Entitas string `form:"entitas"`
	UserID  uint   `form:"user_id"`
}
