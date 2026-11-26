-- Migration: 006_create_dokumen_arsip.sql
-- Tabel dokumen pendukung (digitalisasi arsip — Nota Dinas #8)

CREATE TABLE IF NOT EXISTS dokumen_arsip (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    arsip_id BIGINT UNSIGNED NOT NULL,
    nama_file VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    ukuran_bytes BIGINT NOT NULL,
    keterangan VARCHAR(255) NULL,
    uploaded_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_dokumen_arsip_id (arsip_id),
    CONSTRAINT fk_dokumen_arsip FOREIGN KEY (arsip_id) REFERENCES arsip(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_dokumen_uploader FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
