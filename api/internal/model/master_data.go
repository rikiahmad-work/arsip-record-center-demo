package model

import (
	"time"

	"gorm.io/gorm"
)

// KodeKlasifikasi merepresentasikan master 435 kode klasifikasi (JRA & SKKAD).
type KodeKlasifikasi struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode                string         `gorm:"column:kode;type:varchar(20);uniqueIndex;not null" json:"kode"`
	Nama                string         `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	JangkaSimpanInaktif *string        `gorm:"column:jangka_simpan_inaktif;type:varchar(20)" json:"jangka_simpan_inaktif"`
	NasibAkhir          *string        `gorm:"column:nasib_akhir;type:enum('Musnah','Permanen')" json:"nasib_akhir"`
	SKKAD               string         `gorm:"column:skkad;type:enum('Biasa','Terbatas','Rahasia');not null;default:'Biasa'" json:"skkad"`
	IsActive            bool           `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (KodeKlasifikasi) TableName() string {
	return "kode_klasifikasi"
}

// UnitKerja merepresentasikan unit kerja / pencipta arsip.
type UnitKerja struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"column:kode;type:varchar(20);uniqueIndex;not null" json:"kode"`
	Nama      string         `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	Eselon    *string        `gorm:"column:eselon;type:varchar(20)" json:"eselon"`
	IsActive  bool           `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UnitKerja) TableName() string {
	return "unit_kerja"
}

// LokasiSimpan merepresentasikan lokasi penyimpanan fisik/digital (rak, boks, folder).
type LokasiSimpan struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama      string         `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	Gedung    *string        `gorm:"column:gedung;type:varchar(255)" json:"gedung"`
	Lantai    *string        `gorm:"column:lantai;type:varchar(20)" json:"lantai"`
	NoRak     *string        `gorm:"column:no_rak;type:varchar(50)" json:"no_rak"`
	NoBoks    *string        `gorm:"column:no_boks;type:varchar(50)" json:"no_boks"`
	NoFolder  *string        `gorm:"column:no_folder;type:varchar(50)" json:"no_folder"`
	Kapasitas *int           `gorm:"column:kapasitas" json:"kapasitas"`
	Terpakai  int            `gorm:"column:terpakai;default:0;not null" json:"terpakai"`
	IsActive  bool           `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LokasiSimpan) TableName() string {
	return "lokasi_simpan"
}
