-- Migration: 004_alter_users_role.sql
-- Perluas enum role dari 2 ('super-admin','admin-arsip') ke 4 role kanonik ('administrator','operator','viewer','approver')

ALTER TABLE users MODIFY COLUMN role VARCHAR(50) NOT NULL DEFAULT 'operator';

UPDATE users SET role = 'administrator' WHERE role = 'super-admin';
UPDATE users SET role = 'operator' WHERE role = 'admin-arsip';

ALTER TABLE users MODIFY COLUMN role ENUM('administrator', 'operator', 'viewer', 'approver') NOT NULL DEFAULT 'operator';
