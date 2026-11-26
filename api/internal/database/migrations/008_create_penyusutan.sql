-- Migration: 008_create_penyusutan.sql
-- Tabel pengajuan penyusutan arsip: Usul Musnah & Usul Serah (Nota Dinas #5)

CREATE TABLE IF NOT EXISTS penyusutan (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    no_pengajuan VARCHAR(50) NOT NULL UNIQUE,
    jenis_usul ENUM('Musnah', 'Serah') NOT NULL,
    arsip_id BIGINT UNSIGNED NOT NULL,
    tanggal_retensi DATE NOT NULL,
    alasan TEXT NULL,
    status ENUM('Diajukan', 'Disetujui', 'Ditolak', 'Selesai') NOT NULL DEFAULT 'Diajukan',
    diajukan_oleh BIGINT UNSIGNED NOT NULL,
    disetujui_oleh BIGINT UNSIGNED NULL,
    catatan_approver TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_penyusutan_arsip_id (arsip_id),
    INDEX idx_penyusutan_jenis_usul (jenis_usul),
    INDEX idx_penyusutan_status (status),
    CONSTRAINT fk_penyusutan_arsip FOREIGN KEY (arsip_id) REFERENCES arsip(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_penyusutan_pengaju FOREIGN KEY (diajukan_oleh) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
