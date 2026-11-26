package model

import "time"

// DokumenArsip merepresentasikan file dokumen digitalisasi arsip.
type DokumenArsip struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArsipID     uint      `gorm:"column:arsip_id;not null;index" json:"arsip_id"`
	NamaFile    string    `gorm:"column:nama_file;type:varchar(255);not null" json:"nama_file"`
	Path        string    `gorm:"column:path;type:varchar(500);not null" json:"path"`
	MimeType    string    `gorm:"column:mime_type;type:varchar(100);not null" json:"mime_type"`
	UkuranBytes int64     `gorm:"column:ukuran_bytes;not null" json:"ukuran_bytes"`
	Keterangan  *string   `gorm:"column:keterangan;type:varchar(255)" json:"keterangan"`
	UploadedBy  uint      `gorm:"column:uploaded_by;not null" json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func (DokumenArsip) TableName() string {
	return "dokumen_arsip"
}

// Peminjaman merepresentasikan transaksi peminjaman dan pengembalian arsip.
type Peminjaman struct {
	ID                    uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	NoTiket               string     `gorm:"column:no_tiket;type:varchar(50);uniqueIndex;not null" json:"no_tiket"`
	ArsipID               uint       `gorm:"column:arsip_id;not null;index" json:"arsip_id"`
	Peminjam              string     `gorm:"column:peminjam;type:varchar(255);not null" json:"peminjam"`
	UnitKerjaID           *uint      `gorm:"column:unit_kerja_id;index" json:"unit_kerja_id"`
	TanggalPinjam         time.Time  `gorm:"column:tanggal_pinjam;type:date;not null" json:"tanggal_pinjam"`
	TanggalKembaliRencana time.Time  `gorm:"column:tanggal_kembali_rencana;type:date;not null" json:"tanggal_kembali_rencana"`
	TanggalKembaliAktual  *time.Time `gorm:"column:tanggal_kembali_aktual;type:date" json:"tanggal_kembali_aktual"`
	Status                string     `gorm:"column:status;type:enum('Dipinjam','Dikembalikan','Terlambat');not null;default:'Dipinjam'" json:"status"`
	Catatan               *string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedBy             uint       `gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`

	// Relasi
	Arsip     Arsip      `gorm:"foreignKey:ArsipID" json:"arsip,omitempty"`
	UnitKerja *UnitKerja `gorm:"foreignKey:UnitKerjaID" json:"unit_kerja,omitempty"`
}

func (Peminjaman) TableName() string {
	return "peminjaman"
}

// Penyusutan merepresentasikan pengajuan usul musnah atau usul serah arsip.
type Penyusutan struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NoPengajuan     string    `gorm:"column:no_pengajuan;type:varchar(50);uniqueIndex;not null" json:"no_pengajuan"`
	JenisUsul       string    `gorm:"column:jenis_usul;type:enum('Musnah','Serah');not null" json:"jenis_usul"`
	ArsipID         uint      `gorm:"column:arsip_id;not null;index" json:"arsip_id"`
	TanggalRetensi  time.Time `gorm:"column:tanggal_retensi;type:date;not null" json:"tanggal_retensi"`
	Alasan          *string   `gorm:"column:alasan;type:text" json:"alasan"`
	Status          string    `gorm:"column:status;type:enum('Diajukan','Disetujui','Ditolak','Selesai');not null;default:'Diajukan'" json:"status"`
	DiajukanOleh    uint      `gorm:"column:diajukan_oleh;not null" json:"diajukan_oleh"`
	DisetujuiOleh   *uint     `gorm:"column:disetujui_oleh" json:"disetujui_oleh"`
	CatatanApprover *string   `gorm:"column:catatan_approver;type:text" json:"catatan_approver"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relasi
	Arsip    Arsip `gorm:"foreignKey:ArsipID" json:"arsip,omitempty"`
	Pengaju  User  `gorm:"foreignKey:DiajukanOleh" json:"pengaju,omitempty"`
	Approver *User `gorm:"foreignKey:DisetujuiOleh" json:"approver,omitempty"`
}

func (Penyusutan) TableName() string {
	return "penyusutan"
}

// KonfigurasiSistem merepresentasikan key-value pengaturan aplikasi.
type KonfigurasiSistem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"column:key;type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     *string   `gorm:"column:value;type:text" json:"value"`
	UpdatedBy *uint     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (KonfigurasiSistem) TableName() string {
	return "konfigurasi_sistem"
}

// AuditLog merepresentasikan jejak audit aktivitas pengguna dan perubahan data.
type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint     `gorm:"column:user_id;index" json:"user_id"`
	Aksi      string    `gorm:"column:aksi;type:varchar(50);not null;index" json:"aksi"`
	Entitas   string    `gorm:"column:entitas;type:varchar(50);not null;index" json:"entitas"`
	EntitasID *uint     `gorm:"column:entitas_id;index" json:"entitas_id"`
	DataLama  *string   `gorm:"column:data_lama;type:json" json:"data_lama"`
	DataBaru  *string   `gorm:"column:data_baru;type:json" json:"data_baru"`
	IPAddress *string   `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
