-- Migrasi 011: Tabel Multi-Factor Authentication (MFA TOTP & Recovery Codes)
-- Sesuai standar industri (RFC 6238, OWASP ASVS v4.0, NIST SP 800-63B) & pola SIMPEG Terintegrasi

CREATE TABLE IF NOT EXISTS `user_mfas` (
    `user_id` BIGINT UNSIGNED NOT NULL,
    `secret` TEXT NULL,
    `pending_secret` TEXT NULL,
    `recovery_codes` JSON NULL,
    `last_totp_counter` BIGINT UNSIGNED NULL,
    `failed_attempts` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `locked_until` DATETIME NULL,
    `enabled_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    CONSTRAINT `fk_user_mfas_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `mfa_challenges` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `challenge_hash` CHAR(64) NOT NULL,
    `purpose` VARCHAR(24) NOT NULL,
    `attempts` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `ip_address` VARCHAR(45) NULL,
    `user_agent_hash` CHAR(64) NULL,
    `expires_at` DATETIME NOT NULL,
    `token_expires_at` DATETIME NOT NULL,
    `consumed_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_mfa_challenges_hash` (`challenge_hash`),
    INDEX `idx_mfa_challenges_lookup` (`user_id`, `purpose`, `expires_at`),
    CONSTRAINT `fk_mfa_challenges_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `mfa_audit_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `actor_user_id` BIGINT UNSIGNED NULL,
    `target_user_id` BIGINT UNSIGNED NOT NULL,
    `action` VARCHAR(40) NOT NULL,
    `ip_address` VARCHAR(45) NULL,
    `user_agent` VARCHAR(255) NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_mfa_audit_actor` (`actor_user_id`),
    INDEX `idx_mfa_audit_target` (`target_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
