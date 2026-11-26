package seeders

import (
	"log"

	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// SeedDefaultUsers membuat 4 akun default sesuai 4 role kanonik (idempoten).
func SeedDefaultUsers(db *gorm.DB) error {
	hashHelper := helper.NewHashHelper()

	defaultUsers := []struct {
		Username string
		Name     string
		NIP      string
		Email    string
		Password string
		Role     string
		Unit     string
	}{
		{
			Username: "administrator",
			Name:     "Administrator",
			NIP:      "198501012010011001",
			Email:    "superadmin@kemenpppa.go.id",
			Password: "superadmin123",
			Role:     string(model.RoleAdministrator),
			Unit:     "Biro Umum & SDM",
		},
		{
			Username: "operator1",
			Name:     "Operator 1",
			NIP:      "199002022010012002",
			Email:    "adminarsip@kemenpppa.go.id",
			Password: "adminarsip123",
			Role:     string(model.RoleOperator),
			Unit:     "Biro Umum & SDM",
		},
		{
			Username: "viewer1",
			Name:     "Viewer 1",
			NIP:      "199503032015031003",
			Email:    "viewer@kemenpppa.go.id",
			Password: "viewer123",
			Role:     string(model.RoleViewer),
			Unit:     "Sekretariat Kementerian",
		},
		{
			Username: "approver1",
			Name:     "Approver 1",
			NIP:      "198804042010121004",
			Email:    "approver@kemenpppa.go.id",
			Password: "approver123",
			Role:     string(model.RoleApprover),
			Unit:     "Biro Umum & SDM",
		},
	}

	for _, u := range defaultUsers {
		var existing model.User
		result := db.Unscoped().Where("nip = ?", u.NIP).First(&existing)
		if result.Error == nil {
			// Pastikan role dan status aktif ter-update ke role kanonik
			unitKerja := u.Unit
			_ = db.Unscoped().Model(&existing).Updates(map[string]interface{}{
				"role":       u.Role,
				"name":       u.Name,
				"unit_kerja": &unitKerja,
				"is_active":  true,
				"deleted_at": nil,
			}).Error
			log.Printf("✅ User dengan NIP %s diperbarui ke role %s", u.NIP, u.Role)
			continue
		}

		hashedPassword, err := hashHelper.Hash(u.Password)
		if err != nil {
			return err
		}

		unitKerja := u.Unit
		user := model.User{
			Username:  u.Username,
			Name:      u.Name,
			NIP:       u.NIP,
			Email:     u.Email,
			Password:  hashedPassword,
			Role:      u.Role,
			UnitKerja: &unitKerja,
			IsActive:  true,
			IsSSO:     false,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		log.Printf("✅ User %s (%s) berhasil dibuat", u.Name, u.Role)
	}

	return nil
}

// SeedMasterPendukung mengisi data awal unit_kerja, lokasi_simpan, dan konfigurasi_sistem (idempoten).
func SeedMasterPendukung(db *gorm.DB) error {
	eselon2 := "Eselon II"
	eselon1 := "Eselon I"
	units := []model.UnitKerja{
		{Kode: "BU.01", Nama: "Biro Umum & SDM", Eselon: &eselon2, IsActive: true},
		{Kode: "BD.02", Nama: "Biro Data dan Informasi", Eselon: &eselon2, IsActive: true},
		{Kode: "BH.03", Nama: "Biro Hukum dan Humas", Eselon: &eselon2, IsActive: true},
		{Kode: "BP.04", Nama: "Biro Perencanaan dan Keuangan", Eselon: &eselon2, IsActive: true},
		{Kode: "SK.05", Nama: "Sekretariat Kementerian", Eselon: &eselon1, IsActive: true},
		{Kode: "IN.06", Nama: "Inspektorat", Eselon: &eselon2, IsActive: true},
		{Kode: "KG.07", Nama: "Deputi Bidang Kesetaraan Gender", Eselon: &eselon1, IsActive: true},
		{Kode: "HA.08", Nama: "Deputi Bidang Pemenuhan Hak Anak", Eselon: &eselon1, IsActive: true},
		{Kode: "PP.09", Nama: "Deputi Bidang Perlindungan Hak Perempuan", Eselon: &eselon1, IsActive: true},
		{Kode: "PA.10", Nama: "Deputi Bidang Perlindungan Khusus Anak", Eselon: &eselon1, IsActive: true},
	}

	for _, uk := range units {
		var existing model.UnitKerja
		if err := db.Where("kode = ?", uk.Kode).First(&existing).Error; err != nil {
			if err := db.Create(&uk).Error; err != nil {
				return err
			}
		}
	}

	gedungUtama := "Gedung Utama KemenPPPA"
	lt2 := "2"
	lt3 := "3"
	rak1 := "R-01"
	boks1 := "B-001"
	folder1 := "F-01"
	kap500 := 500
	lokasiList := []model.LokasiSimpan{
		{Nama: "Record Center", Gedung: &gedungUtama, Lantai: &lt2, NoRak: &rak1, NoBoks: &boks1, NoFolder: &folder1, Kapasitas: &kap500, Terpakai: 3, IsActive: true},
		{Nama: "Ruang Arsip Lt. 3", Gedung: &gedungUtama, Lantai: &lt3, NoRak: &rak1, NoBoks: &boks1, NoFolder: &folder1, Kapasitas: &kap500, Terpakai: 0, IsActive: true},
	}

	for _, lok := range lokasiList {
		var existing model.LokasiSimpan
		if err := db.Where("nama = ?", lok.Nama).First(&existing).Error; err != nil {
			if err := db.Create(&lok).Error; err != nil {
				return err
			}
		}
	}

	configs := []struct {
		Key   string
		Value string
	}{
		{Key: "app_name", Value: "SiARSIP"},
		{Key: "default_jangka_simpan", Value: "3 Tahun"},
		{Key: "default_nasib_akhir", Value: "Musnah"},
		{Key: "notifikasi_jatuh_tempo_hari", Value: "30"},
	}

	for _, c := range configs {
		var existing model.KonfigurasiSistem
		if err := db.Where("`key` = ?", c.Key).First(&existing).Error; err != nil {
			val := c.Value
			if err := db.Create(&model.KonfigurasiSistem{Key: c.Key, Value: &val}).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("✅ Master pendukung (unit_kerja, lokasi_simpan, konfigurasi_sistem) berhasil di-seed")
	return nil
}
