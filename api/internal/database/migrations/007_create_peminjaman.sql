-- Migration: 007_create_peminjaman.sql
-- Tabel transaksi peminjaman dan pengembalian arsip (Nota Dinas #4)

CREATE TABLE IF NOT EXISTS peminjaman (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    no_tiket VARCHAR(50) NOT NULL UNIQUE,
    arsip_id BIGINT UNSIGNED NOT NULL,
    peminjam VARCHAR(255) NOT NULL,
    unit_kerja_id BIGINT UNSIGNED NULL,
    tanggal_pinjam DATE NOT NULL,
    tanggal_kembali_rencana DATE NOT NULL,
    tanggal_kembali_aktual DATE NULL,
    status ENUM('Dipinjam', 'Dikembalikan', 'Terlambat') NOT NULL DEFAULT 'Dipinjam',
    catatan TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_peminjaman_arsip_id (arsip_id),
    INDEX idx_peminjaman_status (status),
    CONSTRAINT fk_peminjaman_arsip FOREIGN KEY (arsip_id) REFERENCES arsip(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_peminjaman_creator FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
