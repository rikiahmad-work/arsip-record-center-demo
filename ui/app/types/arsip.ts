// ================================================================
// SiARSIP — Type Definitions Kanonik (20 Field DAI + Enum Kanonik)
// ================================================================

// ─── Enum Kanonik (selaras dengan api/internal/model/enums.go) ──

export type JenisArsip = 'Konvensional' | 'Elektronik' | 'KONVENSIONAL' | 'ELEKTRONIK'
export type TingkatPerkembangan = 'Asli' | 'Copy' | 'Salinan' | 'Tembusan' | 'ASLI' | 'COPY'
export type NasibAkhir = 'Musnah' | 'Permanen' | 'MUSNAH' | 'PERMANEN'
export type SKKAD = 'Biasa' | 'Terbatas' | 'Rahasia' | 'TERBUKA' | 'TERTUTUP'
export type KategoriArsip =
  | 'Inaktif'
  | 'Vital'
  | 'Terjaga'
  | 'Usul Musnah'
  | 'Usul Serah'
  | 'INAKTIF'
  | 'VITAL'
  | 'TERJAGA'
  | 'AKTIF'

// ─── Arsip (20 Field DAI + No. Folder) ──────────────────────────

export interface Arsip {
  id: string
  noItemArsip: string
  noFolder: string
  kodeKlasifikasi: string
  uraianInformasi: string
  uraianInformasiArsip: string
  jenisArsip: JenisArsip
  kurunWaktu: string
  tingkatPerkembangan: TingkatPerkembangan
  jumlah: string
  keterangan: string
  nomorDefinitif: string | null
  lokasiSimpanId?: number | null
  lokasiSimpan: string
  noBoks: string
  noRak: string
  jangkaSimpanInaktif: string
  nasibAkhir: NasibAkhir | null
  jangkaSimpanNasibAkhir: NasibAkhir | null
  skkad: SKKAD
  kategoriArsip: KategoriArsip
  unitKerjaId?: number | null
  unorId?: string
  penciptaArsip: string
  petugasInputId?: number
  petugasInput: string
  status: 'draft' | 'final'
  createdAt: string
  updatedAt: string
}

export interface DokumenArsipItem {
  id: number
  arsipId: number
  namaFile: string
  mimeType: string
  ukuranBytes: number
  keterangan: string | null
  downloadUrl: string
  uploadedBy: number
  createdAt: string
}

// ─── Peminjaman ─────────────────────────────────────────────────

export type StatusPeminjaman =
  | 'Dipinjam'
  | 'Dikembalikan'
  | 'Terlambat'
  | 'dipinjam'
  | 'dikembalikan'
  | 'terlambat'

export interface Peminjaman {
  id: string
  noTiket: string
  arsipId?: number
  peminjam: string
  unitKerjaId?: number | null
  unitKerja: string
  noItemArsip: string
  kodeKlasifikasi?: string
  uraianArsip: string
  tanggalPinjam: string
  tanggalRencanaKembali: string
  tanggalKembali: string | null
  status: StatusPeminjaman
  keterangan: string
  createdAt: string
  updatedAt: string
}

// ─── Penyusutan (Usul Musnah & Usul Serah) ──────────────────────

export type JenisUsulPenyusutan = 'Musnah' | 'Serah' | 'musnah' | 'serah'
export type StatusPenyusutan =
  | 'Diajukan'
  | 'Disetujui'
  | 'Ditolak'
  | 'Selesai'
  | 'diajukan'
  | 'disetujui'
  | 'ditolak'
  | 'selesai'

export interface Penyusutan {
  id: string
  noPengajuan?: string
  jenisUsul?: JenisUsulPenyusutan
  arsipId?: number
  noItem: string
  kodeKlasifikasi?: string
  uraian: string
  tanggalRetensi: string
  nasibAkhir: NasibAkhir
  status: StatusPenyusutan
  diajukanOleh: string
  disetujuiOleh: string | null
  tanggalPengajuan: string
  tanggalPersetujuan: string | null
  keterangan: string
  catatanApprover?: string | null
  createdAt: string
  updatedAt: string
}

// ─── Master Data ────────────────────────────────────────────────

export interface KodeKlasifikasi {
  id: string
  kode: string
  nama: string
  jangkaSimpanInaktif?: string | null
  nasibAkhir?: string | null
  skkad?: string
  totalArsip: number
  statusKey: 'aktif' | 'nonaktif'
  statusLabel: string
}

export interface LokasiSimpan {
  id: string
  namaLokasi: string
  gedung: string
  lantai: string
  noRak?: string
  noBoks?: string
  noFolder?: string
  kapasitas: number
  terpakai: number
  statusKey: 'aktif' | 'nonaktif'
  statusLabel: string
}

export interface UnitKerja {
  id: string
  kodeUnit: string
  namaUnit: string
  eselon: string
  totalArsip: number
  statusKey: 'aktif' | 'nonaktif'
  statusLabel: string
}

export interface KonfigurasiSistem {
  namaAplikasi: string
  defaultJangkaSimpanInaktif: string
  defaultNasibAkhir: NasibAkhir
  notifikasiJatuhTempoHari: number
  logoAplikasi: string | null
}

// ─── Dashboard & Laporan DAI ────────────────────────────────────

export interface DashboardStats {
  totalArsip: number
  arsipVital: number
  arsipAktif: number // alias kompatibilitas untuk komponen lama
  arsipInaktif: number
  arsipTerjaga?: number
  jatuhTempoRetensi: number
  arsipDipinjam?: number
  usulPenyusutan?: number
}

export interface ChartDataPoint {
  bulan: string
  jumlah: number
}

export interface DAIRow {
  no: number
  no_item_arsip: string
  kode_klasifikasi: string
  uraian_informasi: string
  jenis_arsip: string
  kurun_waktu: string
  tingkat_perkembangan: string
  jumlah: string
  keterangan: string
  nomor_definitif: string
  lokasi_simpan: string
  no_boks: string
  no_rak: string
  jangka_simpan_inaktif: string
  nasib_akhir: string
  skkad: string
  kategori_arsip: string
  pencipta_arsip: string
}

// ─── Filter / Query ─────────────────────────────────────────────

export interface ArsipFilter {
  search: string
  kodeKlasifikasi: string
  kategori: string
  skkad?: string
  tahun?: string
  unitKerja?: string
  lokasi: string
  page: number
  perPage: number
  sortBy: string
  sortOrder: 'asc' | 'desc'
}

// ─── Role & User Management ─────────────────────────────────────

export type UserStatus = 'aktif' | 'nonaktif'

export interface UserManagement {
  id: string
  username?: string
  nama: string
  nip: string
  email: string
  unitKerja: string
  jabatan?: string
  role: string
  roleKey?: string
  status: UserStatus
  terakhirLogin: string
}

export interface SelectOption {
  value: string
  label: string
}
