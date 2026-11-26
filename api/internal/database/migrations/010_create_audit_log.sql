-- Migration: 010_create_audit_log.sql
-- Tabel jejak audit aktivitas pengguna

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NULL,
    aksi VARCHAR(50) NOT NULL,
    entitas VARCHAR(50) NOT NULL,
    entitas_id BIGINT UNSIGNED NULL,
    data_lama JSON NULL,
    data_baru JSON NULL,
    ip_address VARCHAR(45) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audit_log_user_id (user_id),
    INDEX idx_audit_log_aksi (aksi),
    INDEX idx_audit_log_entitas (entitas, entitas_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
