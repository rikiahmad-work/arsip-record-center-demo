package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	return db.AutoMigrate(
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
	)
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
