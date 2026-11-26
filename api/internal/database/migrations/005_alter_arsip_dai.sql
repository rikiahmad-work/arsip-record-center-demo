-- Migration: 005_alter_arsip_dai.sql
-- Penyesuaian tabel arsip untuk mendukung 20 kolom DAI + No. Folder + Enum Kanonik

ALTER TABLE arsip
    MODIFY COLUMN skkad ENUM('Biasa', 'Terbatas', 'Rahasia') NOT NULL DEFAULT 'Biasa',
    MODIFY COLUMN status ENUM('draft', 'final') NOT NULL DEFAULT 'final';
