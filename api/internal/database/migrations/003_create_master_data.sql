-- Migration: 003_create_master_data.sql
-- Tabel master data: kode_klasifikasi (435 baris), unit_kerja, lokasi_simpan

CREATE TABLE IF NOT EXISTS kode_klasifikasi (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    jangka_simpan_inaktif VARCHAR(20) NULL,
    nasib_akhir ENUM('Musnah', 'Permanen') NULL,
    skkad ENUM('Biasa', 'Terbatas', 'Rahasia') NOT NULL DEFAULT 'Biasa',
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_kode_klasifikasi_skkad (skkad),
    INDEX idx_kode_klasifikasi_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS unit_kerja (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    kode VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    eselon VARCHAR(20) NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_unit_kerja_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lokasi_simpan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    gedung VARCHAR(255) NULL,
    lantai VARCHAR(20) NULL,
    no_rak VARCHAR(50) NULL,
    no_boks VARCHAR(50) NULL,
    no_folder VARCHAR(50) NULL,
    kapasitas INT NULL,
    terpakai INT NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_lokasi_simpan_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
