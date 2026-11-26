// ================================================================
// SiARSIP — Data Statis & Dummy untuk Modul Arsip Record Center
// ================================================================

import type {
  Arsip,
  Peminjaman,
  Penyusutan,
  KodeKlasifikasi,
  LokasiSimpan,
  UnitKerja,
  KonfigurasiSistem,
  DashboardStats,
  ChartDataPoint,
  SelectOption,
  UserManagement,
} from '~/types/arsip'

// ─── Opsi Select Kanonik (Fase 4) ────────────────────────────────

export const JENIS_ARSIP_OPTIONS: SelectOption[] = [
  { value: 'Konvensional', label: 'Konvensional' },
  { value: 'Elektronik', label: 'Elektronik' },
]

export const TINGKAT_PERKEMBANGAN_OPTIONS: SelectOption[] = [
  { value: 'Asli', label: 'Asli' },
  { value: 'Copy', label: 'Copy' },
  { value: 'Salinan', label: 'Salinan' },
  { value: 'Tembusan', label: 'Tembusan' },
]

export const NASIB_AKHIR_OPTIONS: SelectOption[] = [
  { value: 'Musnah', label: 'Musnah' },
  { value: 'Permanen', label: 'Permanen' },
]

export const SKKAD_OPTIONS: SelectOption[] = [
  { value: 'Biasa', label: 'Biasa' },
  { value: 'Terbatas', label: 'Terbatas' },
  { value: 'Rahasia', label: 'Rahasia' },
]

export const KATEGORI_ARSIP_OPTIONS: SelectOption[] = [
  { value: 'Inaktif', label: 'Inaktif' },
  { value: 'Vital', label: 'Vital' },
  { value: 'Terjaga', label: 'Terjaga' },
  { value: 'Usul Musnah', label: 'Usul Musnah' },
  { value: 'Usul Serah', label: 'Usul Serah' },
]

export const STATUS_ARSIP_OPTIONS: SelectOption[] = [
  { value: 'draft', label: 'Draft' },
  { value: 'final', label: 'Final' },
]

export const STATUS_PEMINJAMAN_OPTIONS: SelectOption[] = [
  { value: 'Dipinjam', label: 'Dipinjam' },
  { value: 'Dikembalikan', label: 'Dikembalikan' },
  { value: 'Terlambat', label: 'Terlambat' },
]

export const JENIS_USUL_OPTIONS: SelectOption[] = [
  { value: 'Musnah', label: 'Usul Musnah' },
  { value: 'Serah', label: 'Usul Serah' },
]

export const STATUS_PENYUSUTAN_OPTIONS: SelectOption[] = [
  { value: 'Diajukan', label: 'Diajukan' },
  { value: 'Disetujui', label: 'Disetujui' },
  { value: 'Ditolak', label: 'Ditolak' },
  { value: 'Selesai', label: 'Selesai' },
]

export const ROLE_USER_OPTIONS: SelectOption[] = [
  { value: 'administrator', label: 'Administrator' },
  { value: 'operator', label: 'Operator' },
  { value: 'viewer', label: 'Viewer' },
  { value: 'approver', label: 'Approver' },
]

export const UNIT_KERJA_OPTIONS: SelectOption[] = [
  { value: 'biro-umum-sdm', label: 'Biro Umum & SDM' },
  { value: 'biro-data-statistik', label: 'Biro Data dan Statistik' },
  { value: 'biro-hukum', label: 'Biro Hukum' },
  { value: 'biro-perencanaan', label: 'Biro Perencanaan' },
  { value: 'sekretariat-jenderal', label: 'Sekretariat Jenderal' },
  { value: 'inspektorat', label: 'Inspektorat' },
]

// ─── Kode Klasifikasi (SIMFONI) ─────────────────────────────────

export const DUMMY_KODE_KLASIFIKASI: KodeKlasifikasi[] = [
  { id: 'KK-001', kode: 'DT.10', nama: 'Keputusan Menteri', totalArsip: 45, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-002', kode: 'KP.02', nama: 'Nota Dinas', totalArsip: 120, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-003', kode: 'HK.05', nama: 'Peraturan Perundang-undangan', totalArsip: 32, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-004', kode: 'TU.00', nama: 'Tata Usaha', totalArsip: 78, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-005', kode: 'KU.01', nama: 'Keuangan', totalArsip: 95, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-006', kode: 'PR.03', nama: 'Pengadaan Barang/Jasa', totalArsip: 55, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-007', kode: 'KP.04', nama: 'Kepegawaian', totalArsip: 67, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-008', kode: 'DT.12', nama: 'Surat Edaran', totalArsip: 40, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'KK-009', kode: 'KL.01', nama: 'Kebijakan Program', totalArsip: 23, statusKey: 'nonaktif', statusLabel: 'Nonaktif' },
  { id: 'KK-010', kode: 'LP.01', nama: 'Laporan', totalArsip: 88, statusKey: 'aktif', statusLabel: 'Aktif' },
]

// ─── Lokasi Simpan ──────────────────────────────────────────────

export const DUMMY_LOKASI_SIMPAN: LokasiSimpan[] = [
  { id: 'LS-001', namaLokasi: 'Ruang Arsip Lt. 2', gedung: 'Gedung Utama', lantai: '2', kapasitas: 500, terpakai: 340, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'LS-002', namaLokasi: 'Ruang Arsip Lt. 3', gedung: 'Gedung Utama', lantai: '3', kapasitas: 400, terpakai: 280, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'LS-003', namaLokasi: 'Gudang Arsip Lama', gedung: 'Gedung Annex', lantai: '1', kapasitas: 300, terpakai: 298, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'LS-004', namaLokasi: 'Ruang Data Center', gedung: 'Gedung Utama', lantai: '1', kapasitas: 200, terpakai: 120, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'LS-005', namaLokasi: 'Ruang Arsip Digital', gedung: 'Gedung Annex', lantai: '2', kapasitas: 100, terpakai: 45, statusKey: 'aktif', statusLabel: 'Aktif' },
]

// ─── Unit Kerja ─────────────────────────────────────────────────

export const DUMMY_UNIT_KERJA: UnitKerja[] = [
  { id: 'UK-001', kodeUnit: 'BU.01', namaUnit: 'Biro Umum & SDM', eselon: 'Eselon II', totalArsip: 320, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'UK-002', kodeUnit: 'BD.02', namaUnit: 'Biro Data dan Statistik', eselon: 'Eselon II', totalArsip: 195, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'UK-003', kodeUnit: 'BH.03', namaUnit: 'Biro Hukum', eselon: 'Eselon II', totalArsip: 145, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'UK-004', kodeUnit: 'BP.04', namaUnit: 'Biro Perencanaan', eselon: 'Eselon II', totalArsip: 210, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'UK-005', kodeUnit: 'SJ.05', namaUnit: 'Sekretariat Jenderal', eselon: 'Eselon I', totalArsip: 430, statusKey: 'aktif', statusLabel: 'Aktif' },
  { id: 'UK-006', kodeUnit: 'IN.06', namaUnit: 'Inspektorat', eselon: 'Eselon II', totalArsip: 88, statusKey: 'aktif', statusLabel: 'Aktif' },
]

// ─── Konfigurasi Sistem ─────────────────────────────────────────

export const DUMMY_KONFIGURASI: KonfigurasiSistem = {
  namaAplikasi: 'Arsip Record Center',
  defaultJangkaSimpanInaktif: '5 Tahun',
  defaultNasibAkhir: 'musnah',
  notifikasiJatuhTempoHari: 30,
  logoAplikasi: '/admin/images/logokpppa.png',
}

// ─── Dashboard ──────────────────────────────────────────────────

export const DUMMY_DASHBOARD_STATS: DashboardStats = {
  totalArsip: 2847,
  arsipAktif: 1124,
  arsipInaktif: 1723,
  jatuhTempoRetensi: 47,
}

export const DUMMY_CHART_DATA: ChartDataPoint[] = [
  { bulan: 'Jan', jumlah: 45 },
  { bulan: 'Feb', jumlah: 52 },
  { bulan: 'Mar', jumlah: 38 },
  { bulan: 'Apr', jumlah: 61 },
  { bulan: 'Mei', jumlah: 47 },
  { bulan: 'Jun', jumlah: 55 },
  { bulan: 'Jul', jumlah: 42 },
]

// ─── Dummy Arsip ────────────────────────────────────────────────

export const DUMMY_ARSIP: Arsip[] = [
  {
    id: 'ARS-001',
    noItemArsip: 'ARS/2026/DT.10/001',
    kodeKlasifikasi: 'DT.10',
    uraianInformasiArsip: 'Keputusan Menteri tentang Pedoman Pengelolaan Arsip',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2024-2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '1 Berkas',
    keterangan: 'Kondisi baik',
    nomorDefinitif: 'DEF/DT.10/001',
    lokasiSimpan: 'Ruang Arsip Lt. 2',
    noBoks: 'B-001',
    noRak: 'R-01',
    jangkaSimpanInaktif: '5 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERBUKA',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Biro Umum & SDM',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-01-15',
    updatedAt: '2026-01-15',
  },
  {
    id: 'ARS-002',
    noItemArsip: 'ARS/2026/KP.02/002',
    kodeKlasifikasi: 'KP.02',
    uraianInformasiArsip: 'Nota Dinas tentang Rapat Koordinasi Bulanan',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '1 Berkas',
    keterangan: 'Lengkap',
    nomorDefinitif: 'DEF/KP.02/002',
    lokasiSimpan: 'Ruang Arsip Lt. 2',
    noBoks: 'B-001',
    noRak: 'R-01',
    jangkaSimpanInaktif: '2 Tahun',
    jangkaSimpanNasibAkhir: 'MUSNAH',
    skkad: 'TERTUTUP',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Sekretariat Jenderal',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-02-20',
    updatedAt: '2026-02-20',
  },
  {
    id: 'ARS-003',
    noItemArsip: 'ARS/2026/HK.05/003',
    kodeKlasifikasi: 'HK.05',
    uraianInformasiArsip: 'Peraturan Menteri No. 12 Tahun 2026',
    jenisArsip: 'ELEKTRONIK',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'COPY',
    jumlah: '1 File PDF',
    keterangan: 'Softcopy tersedia',
    nomorDefinitif: 'DEF/HK.05/003',
    lokasiSimpan: 'Ruang Data Center',
    noBoks: '-',
    noRak: 'R-Digital',
    jangkaSimpanInaktif: '10 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERBUKA',
    kategoriArsip: 'VITAL',
    penciptaArsip: 'Biro Hukum',
    petugasInput: 'Operator Arsip',
    status: 'final',
    createdAt: '2026-03-10',
    updatedAt: '2026-03-10',
  },
  {
    id: 'ARS-004',
    noItemArsip: 'ARS/2026/TU.00/004',
    kodeKlasifikasi: 'TU.00',
    uraianInformasiArsip: 'Surat Masuk dan Keluar Bulan Januari-Maret 2026',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '3 Berkas',
    keterangan: 'Terjilid rapi',
    nomorDefinitif: 'DEF/TU.00/004',
    lokasiSimpan: 'Ruang Arsip Lt. 3',
    noBoks: 'B-002',
    noRak: 'R-02',
    jangkaSimpanInaktif: '2 Tahun',
    jangkaSimpanNasibAkhir: 'MUSNAH',
    skkad: 'TERTUTUP',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Biro Umum & SDM',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-04-05',
    updatedAt: '2026-04-05',
  },
  {
    id: 'ARS-005',
    noItemArsip: 'ARS/2026/KU.01/005',
    kodeKlasifikasi: 'KU.01',
    uraianInformasiArsip: 'Laporan Keuangan Triwulan I Tahun 2026',
    jenisArsip: 'ELEKTRONIK',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '1 File Excel',
    keterangan: 'Sudah diaudit',
    nomorDefinitif: 'DEF/KU.01/005',
    lokasiSimpan: 'Ruang Data Center',
    noBoks: '-',
    noRak: 'R-Digital',
    jangkaSimpanInaktif: '10 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERTUTUP',
    kategoriArsip: 'VITAL',
    penciptaArsip: 'Biro Perencanaan',
    petugasInput: 'Operator Arsip',
    status: 'final',
    createdAt: '2026-04-15',
    updatedAt: '2026-04-15',
  },
  {
    id: 'ARS-006',
    noItemArsip: 'ARS/2026/PR.03/006',
    kodeKlasifikasi: 'PR.03',
    uraianInformasiArsip: 'Dokumen Pengadaan Peralatan Kantor Tahun 2026',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '2 Berkas',
    keterangan: 'Lengkap dengan BAST',
    nomorDefinitif: 'DEF/PR.03/006',
    lokasiSimpan: 'Ruang Arsip Lt. 2',
    noBoks: 'B-003',
    noRak: 'R-01',
    jangkaSimpanInaktif: '5 Tahun',
    jangkaSimpanNasibAkhir: 'MUSNAH',
    skkad: 'TERBUKA',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Biro Umum & SDM',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-05-12',
    updatedAt: '2026-05-12',
  },
  {
    id: 'ARS-007',
    noItemArsip: 'ARS/2026/KP.04/007',
    kodeKlasifikasi: 'KP.04',
    uraianInformasiArsip: 'SK Kenaikan Pangkat Pegawai Periode April 2026',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '1 Berkas',
    keterangan: 'Rahasia',
    nomorDefinitif: 'DEF/KP.04/007',
    lokasiSimpan: 'Ruang Arsip Lt. 3',
    noBoks: 'B-004',
    noRak: 'R-02',
    jangkaSimpanInaktif: '20 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERTUTUP',
    kategoriArsip: 'VITAL',
    penciptaArsip: 'Biro Umum & SDM',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-05-20',
    updatedAt: '2026-05-20',
  },
  {
    id: 'ARS-008',
    noItemArsip: 'ARS/2026/LP.01/008',
    kodeKlasifikasi: 'LP.01',
    uraianInformasiArsip: 'Laporan Kinerja Semester I Tahun 2026',
    jenisArsip: 'ELEKTRONIK',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'ASLI',
    jumlah: '1 File PDF',
    keterangan: 'Draft review',
    nomorDefinitif: null,
    lokasiSimpan: 'Ruang Arsip Digital',
    noBoks: '-',
    noRak: 'R-Digital',
    jangkaSimpanInaktif: '5 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERBUKA',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Biro Data dan Statistik',
    petugasInput: 'Operator Arsip',
    status: 'draft',
    createdAt: '2026-07-01',
    updatedAt: '2026-07-01',
  },
  {
    id: 'ARS-009',
    noItemArsip: 'ARS/2026/DT.12/009',
    kodeKlasifikasi: 'DT.12',
    uraianInformasiArsip: 'Surat Edaran tentang Libur Nasional',
    jenisArsip: 'KONVENSIONAL',
    kurunWaktu: '2026',
    tingkatPerkembangan: 'COPY',
    jumlah: '1 Lembar',
    keterangan: 'Informasi umum',
    nomorDefinitif: null,
    lokasiSimpan: 'Ruang Arsip Lt. 2',
    noBoks: 'B-005',
    noRak: 'R-01',
    jangkaSimpanInaktif: '1 Tahun',
    jangkaSimpanNasibAkhir: 'MUSNAH',
    skkad: 'TERBUKA',
    kategoriArsip: 'INAKTIF',
    penciptaArsip: 'Sekretariat Jenderal',
    petugasInput: 'Admin Arsip',
    status: 'final',
    createdAt: '2026-06-10',
    updatedAt: '2026-06-10',
  },
  {
    id: 'ARS-010',
    noItemArsip: 'ARS/2026/HK.05/010',
    kodeKlasifikasi: 'HK.05',
    uraianInformasiArsip: 'Peraturan Pemerintah tentang Kearsipan',
    jenisArsip: 'ELEKTRONIK',
    kurunWaktu: '2025-2026',
    tingkatPerkembangan: 'COPY',
    jumlah: '1 File PDF',
    keterangan: 'Dari sistem JDIH',
    nomorDefinitif: 'DEF/HK.05/010',
    lokasiSimpan: 'Ruang Data Center',
    noBoks: '-',
    noRak: 'R-Digital',
    jangkaSimpanInaktif: '10 Tahun',
    jangkaSimpanNasibAkhir: 'PERMANEN',
    skkad: 'TERBUKA',
    kategoriArsip: 'AKTIF',
    penciptaArsip: 'Biro Hukum',
    petugasInput: 'Operator Arsip',
    status: 'final',
    createdAt: '2026-07-15',
    updatedAt: '2026-07-15',
  },
]

// ─── Dummy Peminjaman ───────────────────────────────────────────

export const DUMMY_PEMINJAMAN: Peminjaman[] = [
  {
    id: 'PMJ-001',
    noTiket: 'PIN/2026/07/001',
    peminjam: 'Drs. Agus Wijaya',
    unitKerja: 'Biro Hukum',
    noItemArsip: 'ARS/2026/HK.05/003',
    uraianArsip: 'Peraturan Menteri No. 12 Tahun 2026',
    tanggalPinjam: '2026-07-10',
    tanggalRencanaKembali: '2026-07-24',
    tanggalKembali: null,
    status: 'dipinjam',
    keterangan: 'Untuk bahan rapat',
    createdAt: '2026-07-10',
    updatedAt: '2026-07-10',
  },
  {
    id: 'PMJ-002',
    noTiket: 'PIN/2026/07/002',
    peminjam: 'Siti Nurhaliza, S.Sos',
    unitKerja: 'Biro Perencanaan',
    noItemArsip: 'ARS/2026/KU.01/005',
    uraianArsip: 'Laporan Keuangan Triwulan I Tahun 2026',
    tanggalPinjam: '2026-07-05',
    tanggalRencanaKembali: '2026-07-19',
    tanggalKembali: '2026-07-18',
    status: 'dikembalikan',
    keterangan: 'Sudah dikembalikan',
    createdAt: '2026-07-05',
    updatedAt: '2026-07-18',
  },
  {
    id: 'PMJ-003',
    noTiket: 'PIN/2026/06/003',
    peminjam: 'Dr. H. Mahmud, M.Si',
    unitKerja: 'Inspektorat',
    noItemArsip: 'ARS/2026/KP.02/002',
    uraianArsip: 'Nota Dinas tentang Rapat Koordinasi Bulanan',
    tanggalPinjam: '2026-06-20',
    tanggalRencanaKembali: '2026-07-05',
    tanggalKembali: null,
    status: 'terlambat',
    keterangan: 'Belum dikembalikan, sudah melewati jatuh tempo',
    createdAt: '2026-06-20',
    updatedAt: '2026-06-20',
  },
]

// ─── Dummy Penyusutan ───────────────────────────────────────────

export const DUMMY_PENYUSUTAN: Penyusutan[] = [
  {
    id: 'PST-001',
    noItem: 'ARS/2021/TU.00/015',
    uraian: 'Surat Masuk Tahun 2021',
    tanggalRetensi: '2026-08-01',
    nasibAkhir: 'musnah',
    status: 'diajukan',
    diajukanOleh: 'Admin Arsip',
    disetujuiOleh: null,
    tanggalPengajuan: '2026-07-15',
    tanggalPersetujuan: null,
    keterangan: 'Jangka simpan 5 tahun habis',
    createdAt: '2026-07-15',
    updatedAt: '2026-07-15',
  },
  {
    id: 'PST-002',
    noItem: 'ARS/2020/LP.01/010',
    uraian: 'Laporan Tahunan 2020',
    tanggalRetensi: '2026-06-01',
    nasibAkhir: 'permanen',
    status: 'disetujui',
    diajukanOleh: 'Operator Arsip',
    disetujuiOleh: 'Super Admin SiARSIP',
    tanggalPengajuan: '2026-05-20',
    tanggalPersetujuan: '2026-06-05',
    keterangan: 'Arsip bernilai historis',
    createdAt: '2026-05-20',
    updatedAt: '2026-06-05',
  },
  {
    id: 'PST-003',
    noItem: 'ARS/2022/KP.02/008',
    uraian: 'Nota Dinas Internal Tahun 2022',
    tanggalRetensi: '2026-07-01',
    nasibAkhir: 'musnah',
    status: 'selesai',
    diajukanOleh: 'Admin Arsip',
    disetujuiOleh: 'Super Admin SiARSIP',
    tanggalPengajuan: '2026-06-01',
    tanggalPersetujuan: '2026-06-15',
    keterangan: 'Telah dimusnahkan sesuai prosedur',
    createdAt: '2026-06-01',
    updatedAt: '2026-06-20',
  },
]

// ─── Dummy Manajemen User ───────────────────────────────────────

export const DUMMY_USER_MANAGEMENT: UserManagement[] = [
  { id: 'USR-001', nama: 'Super Admin SiARSIP', nip: '196805152005011001', email: 'superadmin@kemenpppa.go.id', unitKerja: 'Biro Umum & SDM', role: 'Super Admin', status: 'aktif', terakhirLogin: '2026-07-21 08:30' },
  { id: 'USR-002', nama: 'Admin Arsip', nip: '198705102011012001', email: 'admin.arsip@kemenpppa.go.id', unitKerja: 'Biro Umum & SDM', role: 'Administrator', status: 'aktif', terakhirLogin: '2026-07-21 07:45' },
  { id: 'USR-003', nama: 'Operator Arsip', nip: '199104182015031004', email: 'operator.arsip@kemenpppa.go.id', unitKerja: 'Biro Data dan Statistik', role: 'Operator', status: 'aktif', terakhirLogin: '2026-07-20 16:20' },
  { id: 'USR-004', nama: 'Viewer Arsip', nip: '199203152016021005', email: 'viewer@kemenpppa.go.id', unitKerja: 'Sekretariat Jenderal', role: 'Viewer', status: 'aktif', terakhirLogin: '2026-07-19 10:15' },
  { id: 'USR-005', nama: 'Dina Maharani', nip: '198912222014022003', email: 'dina.maharani@kemenpppa.go.id', unitKerja: 'Inspektorat', role: 'Operator', status: 'nonaktif', terakhirLogin: '2026-06-15 09:00' },
  { id: 'USR-006', nama: 'Budi Santoso', nip: '199006152017011006', email: 'budi.santoso@kemenpppa.go.id', unitKerja: 'Biro Hukum', role: 'Viewer', status: 'aktif', terakhirLogin: '2026-07-18 13:30' },
]

// ─── Helper: Map status ke warna badge ──────────────────────────

export const STATUS_BADGE_COLOR: Record<string, string> = {
  // Arsip
  asli: 'bg-success',
  salinan: 'bg-warning',
  tembusan: 'bg-info',
  aktif: 'bg-success',
  inaktif: 'bg-secondary',
  vital: 'bg-warning',
  draft: 'bg-warning',
  final: 'bg-success',
  // Peminjaman
  dipinjam: 'bg-primary',
  dikembalikan: 'bg-success',
  terlambat: 'bg-danger',
  // Penyusutan
  diajukan: 'bg-info',
  disetujui: 'bg-success',
  ditolak: 'bg-danger',
  selesai: 'bg-secondary',
  // User
  'super-admin': 'bg-warning',
  admin: 'bg-primary',
  operator: 'bg-info',
  viewer: 'bg-secondary',
  // SIKD
  terbuka: 'bg-success',
  tertutup: 'bg-danger',
  // Kode Klasifikasi
  'bg-label-primary': 'bg-light-primary',
  'bg-label-info': 'bg-light-info',
}

export const STATUS_BADGE_LABEL: Record<string, string> = {
  asli: 'Asli',
  salinan: 'Salinan',
  tembusan: 'Tembusan',
  aktif: 'Aktif',
  inaktif: 'Inaktif',
  vital: 'Vital',
  draft: 'Draft',
  final: 'Final',
  dipinjam: 'Dipinjam',
  dikembalikan: 'Dikembalikan',
  terlambat: 'Terlambat',
  diajukan: 'Diajukan',
  disetujui: 'Disetujui',
  ditolak: 'Ditolak',
  selesai: 'Selesai',
  'super-admin': 'Super Admin',
  admin: 'Admin',
  operator: 'Operator',
  viewer: 'Viewer',
  terbuka: 'Terbuka',
  tertutup: 'Tertutup',
  musnah: 'Musnah',
  permanen: 'Permanen',
  konvensional: 'Konvensional',
  digital: 'Digital',
}
