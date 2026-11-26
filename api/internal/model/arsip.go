package model

import (
	"time"

	"gorm.io/gorm"
)

// Konstanta kompatibilitas untuk kode existing
const (
	JenisKonvensional string = "KONVENSIONAL"
	JenisElektronik   string = "ELEKTRONIK"

	TingkatAsli string = "ASLI"
	TingkatCopy string = "COPY"

	NasibMusnah   string = "MUSNAH"
	NasibPermanen string = "PERMANEN"

	KategoriInaktif string = "INAKTIF"
	KategoriVital   string = "VITAL"
	KategoriTerjaga string = "TERJAGA"

	StatusDraft string = "draft"
	StatusFinal string = "final"
)

// Arsip merepresentasikan entitas arsip inaktif (20 kolom DAI + No. Folder).
type Arsip struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	NoItemArsip            string         `gorm:"column:no_item_arsip;type:varchar(100);not null;uniqueIndex:idx_arsip_unor_item,priority:2" json:"no_item_arsip"`
	NoFolder               *string        `gorm:"column:no_folder;type:varchar(50)" json:"no_folder"`
	KodeKlasifikasi        string         `gorm:"column:kode_klasifikasi;type:varchar(50);not null;index" json:"kode_klasifikasi"`
	UraianInformasiArsip   string         `gorm:"column:uraian_informasi_arsip;type:text;not null" json:"uraian_informasi_arsip"`
	JenisArsip             string         `gorm:"column:jenis_arsip;type:varchar(30);not null" json:"jenis_arsip"`
	KurunWaktu             string         `gorm:"column:kurun_waktu;type:varchar(50);not null" json:"kurun_waktu"`
	TingkatPerkembangan    string         `gorm:"column:tingkat_perkembangan;type:varchar(30);not null" json:"tingkat_perkembangan"`
	Jumlah                 string         `gorm:"column:jumlah;type:varchar(100);not null" json:"jumlah"`
	Keterangan             string         `gorm:"column:keterangan;type:varchar(255)" json:"keterangan"`
	NomorDefinitif         *string        `gorm:"column:nomor_definitif;type:varchar(100)" json:"nomor_definitif"`
	LokasiSimpanID         *uint          `gorm:"column:lokasi_simpan_id;index" json:"lokasi_simpan_id"`
	LokasiSimpan           string         `gorm:"column:lokasi_simpan;type:varchar(255);not null" json:"lokasi_simpan"`
	NoBoks                 string         `gorm:"column:no_boks;type:varchar(50);not null" json:"no_boks"`
	NoRak                  string         `gorm:"column:no_rak;type:varchar(50);not null" json:"no_rak"`
	JangkaSimpanInaktif    string         `gorm:"column:jangka_simpan_inaktif;type:varchar(100);not null" json:"jangka_simpan_inaktif"`
	JangkaSimpanNasibAkhir *string        `gorm:"column:jangka_simpan_nasib_akhir;type:varchar(30)" json:"jangka_simpan_nasib_akhir"`
	NasibAkhir             *string        `gorm:"column:nasib_akhir;type:enum('Musnah','Permanen')" json:"nasib_akhir"`
	SKKAD                  string         `gorm:"column:skkad;type:enum('Biasa','Terbatas','Rahasia');not null;default:'Biasa'" json:"skkad"`
	KategoriArsip          string         `gorm:"column:kategori_arsip;type:varchar(100);not null;default:'Inaktif'" json:"kategori_arsip"`
	UnitKerjaID            *uint          `gorm:"column:unit_kerja_id;index" json:"unit_kerja_id"`
	UnorID                 string         `gorm:"column:unor_id;type:varchar(36);not null;uniqueIndex:idx_arsip_unor_item,priority:1" json:"unor_id"`
	PenciptaArsip          string         `gorm:"column:pencipta_arsip;type:varchar(255);not null" json:"pencipta_arsip"`
	PetugasInputID         uint           `gorm:"column:petugas_input_id;not null" json:"petugas_input_id"`
	Status                 string         `gorm:"column:status;type:enum('draft','final');not null;default:'final'" json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	PetugasInput    User           `gorm:"foreignKey:PetugasInputID" json:"petugas_input,omitempty"`
	LokasiSimpanRef *LokasiSimpan  `gorm:"foreignKey:LokasiSimpanID" json:"lokasi_simpan_ref,omitempty"`
	UnitKerjaRef    *UnitKerja     `gorm:"foreignKey:UnitKerjaID" json:"unit_kerja_ref,omitempty"`
	DokumenList     []DokumenArsip `gorm:"foreignKey:ArsipID" json:"dokumen_list,omitempty"`
	PeminjamanList  []Peminjaman   `gorm:"foreignKey:ArsipID" json:"peminjaman_list,omitempty"`
	PenyusutanList  []Penyusutan   `gorm:"foreignKey:ArsipID" json:"penyusutan_list,omitempty"`
}

// TableName override default table name ke "arsip".
func (Arsip) TableName() string {
	return "arsip"
}
