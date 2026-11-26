package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQL membuka koneksi ke MySQL dan mengembalikan *gorm.DB.
// Jika database belum ada, akan dibuat otomatis.
func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	// Konfigurasi GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Buat database dulu jika belum ada
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return nil, fmt.Errorf("gagal membuat database: %w", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi database: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// createDatabaseIfNotExists membuat database MySQL jika belum ada.
func createDatabaseIfNotExists(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.Name)
	if err := db.Exec(sql).Error; err != nil {
		return err
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	log.Printf("✅ Database `%s` siap", cfg.Database.Name)
	return nil
}

// AutoMigrate menjalankan migrasi tabel secara otomatis dan menormalisasi data role/enum legacy.
func AutoMigrate(db *gorm.DB) error {
	// 1. Jika tabel users sudah ada dengan enum lama ('super-admin','admin-arsip'),
	//    perlebar sementara ke VARCHAR(50) dan migrasikan nilai lama ke kanonik sebelum GORM mengubah ENUM.
	if db.Migrator().HasTable(&model.User{}) {
		_ = db.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(50) NOT NULL DEFAULT 'operator'").Error
		_ = db.Exec("UPDATE users SET role = 'administrator' WHERE role = 'super-admin'").Error
		_ = db.Exec("UPDATE users SET role = 'operator' WHERE role = 'admin-arsip'").Error
	}

	// 2. Jika tabel arsip sudah ada dengan nilai SKKAD lama, normalisasi ke Biasa/Terbatas/Rahasia
	if db.Migrator().HasTable(&model.Arsip{}) {
		_ = db.Exec("ALTER TABLE arsip MODIFY COLUMN skkad VARCHAR(50) NOT NULL DEFAULT 'Biasa'").Error
		_ = db.Exec("UPDATE arsip SET skkad = 'Biasa' WHERE UPPER(skkad) = 'TERBUKA' OR skkad = ''").Error
		_ = db.Exec("UPDATE arsip SET skkad = 'Terbatas' WHERE UPPER(skkad) = 'TERTUTUP'").Error
		_ = db.Exec("UPDATE arsip SET kategori_arsip = 'Inaktif' WHERE UPPER(kategori_arsip) IN ('INAKTIF', 'AKTIF')").Error
		_ = db.Exec("UPDATE arsip SET kategori_arsip = 'Vital' WHERE UPPER(kategori_arsip) = 'VITAL'").Error
		_ = db.Exec("UPDATE arsip SET kategori_arsip = 'Terjaga' WHERE UPPER(kategori_arsip) = 'TERJAGA'").Error
	}

	// 3. AutoMigrate seluruh tabel entitas
	if err := db.AutoMigrate(
		&model.User{},
		&model.BlacklistedToken{},
		&model.RolePermission{},
		&model.UserMfa{},
		&model.MfaChallenge{},
		&model.MfaAuditLog{},
		&model.KodeKlasifikasi{},
		&model.UnitKerja{},
		&model.LokasiSimpan{},
		&model.Arsip{},
		&model.DokumenArsip{},
		&model.Peminjaman{},
		&model.Penyusutan{},
		&model.KonfigurasiSistem{},
		&model.AuditLog{},
	); err != nil {
		return err
	}

	// 4. Pastikan kode klasifikasi DT.10 (Pengumpulan dan Pengolahan Data) tersedia di master
	ensureDT10KodeKlasifikasi(db)

	// 5. Sinkronisasi unit_kerja_id, unor_id, dan normalisasi no_item_arsip menjadi nomor urut (1, 2, 3...) per Unit Kerja
	normalizeArsipNoItemPerUnitKerja(db)
	return nil
}

func ensureDT10KodeKlasifikasi(db *gorm.DB) {
	var count int64
	if err := db.Model(&model.KodeKlasifikasi{}).Where("kode = ?", "DT.10").Count(&count).Error; err == nil && count == 0 {
		inaktif := "3 Tahun"
		nasib := "Musnah"
		_ = db.Create(&model.KodeKlasifikasi{
			Kode:                "DT.10",
			Nama:                "Pengumpulan dan Pengolahan Data (DT.10)",
			JangkaSimpanInaktif: &inaktif,
			NasibAkhir:          &nasib,
			SKKAD:               "Biasa",
			IsActive:            true,
		}).Error
	}
}

// normalizeArsipNoItemPerUnitKerja memastikan seluruh data arsip memiliki unit_kerja_id & unor_id yang selaras
// dengan master unit_kerja, serta menormalisasi no_item_arsip yang masih berupa kode khusus (misal ARS/...)
// menjadi nomor urut biasa mulai dari 1 sesuai Unit Kerja.
func normalizeArsipNoItemPerUnitKerja(db *gorm.DB) {
	var units []model.UnitKerja
	if err := db.Find(&units).Error; err != nil {
		return
	}

	unitByID := make(map[uint]model.UnitKerja, len(units))
	unitByKode := make(map[string]model.UnitKerja, len(units))
	unitByNama := make(map[string]model.UnitKerja, len(units))
	for _, u := range units {
		unitByID[u.ID] = u
		unitByKode[strings.ToUpper(strings.TrimSpace(u.Kode))] = u
		unitByNama[strings.ToUpper(strings.TrimSpace(u.Nama))] = u
	}

	var allArsip []model.Arsip
	if err := db.Unscoped().Order("id ASC").Find(&allArsip).Error; err != nil || len(allArsip) == 0 {
		return
	}

	// Bebaskan sementara no_item_arsip pada baris soft-deleted agar tidak bentrok dengan UNIQUE KEY (unor_id, no_item_arsip)
	for i := range allArsip {
		a := &allArsip[i]
		if a.DeletedAt.Valid && !strings.Contains(a.NoItemArsip, "_del_") {
			newDelNo := fmt.Sprintf("%s_del_%d", a.NoItemArsip, a.ID)
			_ = db.Unscoped().Model(&model.Arsip{}).Where("id = ?", a.ID).Update("no_item_arsip", newDelNo).Error
			a.NoItemArsip = newDelNo
		}
	}

	// Kelompokkan baris aktif berdasarkan Unit Kerja
	grouped := make(map[string][]*model.Arsip)
	for i := range allArsip {
		a := &allArsip[i]
		if a.DeletedAt.Valid {
			continue
		}

		var matched *model.UnitKerja
		if a.UnitKerjaID != nil {
			if u, ok := unitByID[*a.UnitKerjaID]; ok {
				uCopy := u
				matched = &uCopy
			}
		}
		if matched == nil {
			if u, ok := unitByKode[strings.ToUpper(strings.TrimSpace(a.UnorID))]; ok {
				uCopy := u
				matched = &uCopy
			}
		}
		if matched == nil {
			if u, ok := unitByNama[strings.ToUpper(strings.TrimSpace(a.PenciptaArsip))]; ok {
				uCopy := u
				matched = &uCopy
			}
		}

		if matched != nil {
			uid := matched.ID
			a.UnitKerjaID = &uid
			a.UnorID = matched.Kode
			a.PenciptaArsip = matched.Nama
		} else if strings.TrimSpace(a.UnorID) == "" {
			a.UnorID = "BU.01"
		}

		grouped[a.UnorID] = append(grouped[a.UnorID], a)
	}

	for unorID, list := range grouped {
		usedNum := make(map[int]bool, len(list))
		needsAssign := make([]*model.Arsip, 0)

		for _, a := range list {
			trimmed := strings.TrimSpace(a.NoItemArsip)
			if n, err := strconv.Atoi(trimmed); err == nil && n > 0 && !usedNum[n] {
				usedNum[n] = true
				a.NoItemArsip = strconv.Itoa(n)
			} else {
				needsAssign = append(needsAssign, a)
			}
		}

		// Jika ada arsip yang memakai kode khusus (misal ARS/...) atau duplikat, beri nomor urut mulai dari 1
		nextCounter := 1
		for _, a := range needsAssign {
			for usedNum[nextCounter] {
				nextCounter++
			}
			usedNum[nextCounter] = true
			a.NoItemArsip = strconv.Itoa(nextCounter)
			nextCounter++
		}

		for _, a := range list {
			_ = db.Model(&model.Arsip{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
				"unit_kerja_id":  a.UnitKerjaID,
				"unor_id":        unorID,
				"pencipta_arsip": a.PenciptaArsip,
				"no_item_arsip":  a.NoItemArsip,
			}).Error
		}
	}
}

// RunSQLMigration menjalankan file SQL migration dari folder migrations/.
func RunSQLMigration(db *gorm.DB, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("gagal membaca folder migrations: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("gagal membaca file %s: %w", file.Name(), err)
		}

		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("gagal menjalankan migrasi %s: %w", file.Name(), err)
		}

		log.Printf("Migrasi %s berhasil dijalankan", file.Name())
	}

	return nil
}
