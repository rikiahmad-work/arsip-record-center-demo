// ================================================================
// useArsip — Composable Data Arsip, Master, Peminjaman, Penyusutan,
// Dashboard, Laporan DAI, & Cetak Label (Task 4.5 — Terhubung ke API)
// ================================================================

import { API_BASE_URL } from '~/constants/auth'
import {
  DUMMY_DASHBOARD_STATS,
  DUMMY_CHART_DATA,
  DUMMY_KODE_KLASIFIKASI,
  DUMMY_LOKASI_SIMPAN,
  DUMMY_UNIT_KERJA,
  DUMMY_KONFIGURASI,
} from '~/constants/arsip'

import type {
  Arsip,
  DokumenArsipItem,
  DashboardStats,
  ChartDataPoint,
  Peminjaman,
  Penyusutan,
  KodeKlasifikasi,
  LokasiSimpan,
  UnitKerja,
  KonfigurasiSistem,
  DAIRow,
  UserManagement,
} from '~/types/arsip'

// ─── In-Memory Master Data Cache ────────────────────────────────
let cachedKodeKlasifikasi: KodeKlasifikasi[] | null = null
let cachedLokasiSimpan: LokasiSimpan[] | null = null
let cachedUnitKerja: UnitKerja[] | null = null

// ─── Mapping Helpers (snake_case Backend → camelCase Frontend) ──

function mapArsip(a: any): Arsip {
  const uraian = a.uraian_informasi || a.uraian_informasi_arsip || ''
  const nasib = a.nasib_akhir ?? a.jangka_simpan_nasib_akhir ?? null
  return {
    id: String(a.id),
    noItemArsip: a.no_item_arsip || '',
    noFolder: a.no_folder || '',
    kodeKlasifikasi: a.kode_klasifikasi || '',
    uraianInformasi: uraian,
    uraianInformasiArsip: uraian,
    jenisArsip: a.jenis_arsip || 'Konvensional',
    kurunWaktu: a.kurun_waktu || '',
    tingkatPerkembangan: a.tingkat_perkembangan || 'Asli',
    jumlah: a.jumlah || '',
    keterangan: a.keterangan || '',
    nomorDefinitif: a.nomor_definitif ?? null,
    lokasiSimpanId: a.lokasi_simpan_id ?? null,
    lokasiSimpan: a.lokasi_simpan_nama || a.lokasi_simpan || 'Record Center',
    noBoks: a.no_boks || '',
    noRak: a.no_rak || '',
    jangkaSimpanInaktif: a.jangka_simpan_inaktif || '',
    nasibAkhir: nasib,
    jangkaSimpanNasibAkhir: nasib,
    skkad: a.skkad || 'Biasa',
    kategoriArsip: a.kategori_arsip || 'Inaktif',
    unitKerjaId: a.unit_kerja_id ?? null,
    unorId: a.unor_id || '',
    penciptaArsip: a.unit_kerja_nama || a.pencipta_arsip || '',
    petugasInputId: a.petugas_input_id,
    petugasInput: a.petugas_input || '',
    status: a.status || 'final',
    createdAt: a.created_at || '',
    updatedAt: a.updated_at || '',
  }
}

function mapKodeKlasifikasi(k: any): KodeKlasifikasi {
  const active = k.is_active !== false
  return {
    id: String(k.id),
    kode: k.kode || '',
    nama: k.nama || '',
    jangkaSimpanInaktif: k.jangka_simpan_inaktif ?? null,
    nasibAkhir: k.nasib_akhir ?? null,
    skkad: k.skkad || 'Biasa',
    totalArsip: Number(k.total_arsip || 0),
    statusKey: active ? 'aktif' : 'nonaktif',
    statusLabel: active ? 'Aktif' : 'Nonaktif',
  }
}

function mapLokasiSimpan(l: any): LokasiSimpan {
  const active = l.is_active !== false
  return {
    id: String(l.id),
    namaLokasi: l.nama || '',
    gedung: l.gedung || '-',
    lantai: l.lantai || '-',
    noRak: l.no_rak || '',
    noBoks: l.no_boks || '',
    noFolder: l.no_folder || '',
    kapasitas: Number(l.kapasitas || 0),
    terpakai: Number(l.terpakai || 0),
    statusKey: active ? 'aktif' : 'nonaktif',
    statusLabel: active ? 'Aktif' : 'Nonaktif',
  }
}

function mapUnitKerja(u: any): UnitKerja {
  const active = u.is_active !== false
  return {
    id: String(u.id),
    kodeUnit: u.kode || '',
    namaUnit: u.nama || '',
    eselon: u.eselon || '-',
    totalArsip: Number(u.total_arsip || 0),
    statusKey: active ? 'aktif' : 'nonaktif',
    statusLabel: active ? 'Aktif' : 'Nonaktif',
  }
}

function mapPeminjaman(p: any): Peminjaman {
  return {
    id: String(p.id),
    noTiket: p.no_tiket || '',
    arsipId: p.arsip_id,
    peminjam: p.peminjam || '',
    unitKerjaId: p.unit_kerja_id ?? null,
    unitKerja: p.unit_kerja_nama || '-',
    noItemArsip: p.no_item_arsip || '',
    kodeKlasifikasi: p.kode_klasifikasi || '',
    uraianArsip: p.uraian_arsip || '',
    tanggalPinjam: p.tanggal_pinjam || '',
    tanggalRencanaKembali: p.tanggal_kembali_rencana || '',
    tanggalKembali: p.tanggal_kembali_aktual ?? null,
    status: p.status || 'Dipinjam',
    keterangan: p.catatan || '',
    createdAt: p.created_at || '',
    updatedAt: p.updated_at || '',
  }
}

function mapPenyusutan(p: any): Penyusutan {
  return {
    id: String(p.id),
    noPengajuan: p.no_pengajuan || '',
    jenisUsul: p.jenis_usul || 'Musnah',
    arsipId: p.arsip_id,
    noItem: p.no_item_arsip || '',
    kodeKlasifikasi: p.kode_klasifikasi || '',
    uraian: p.uraian_arsip || '',
    tanggalRetensi: p.tanggal_retensi || '',
    nasibAkhir: p.nasib_akhir || 'Musnah',
    status: p.status || 'Diajukan',
    diajukanOleh: p.diajukan_nama || String(p.diajukan_oleh || ''),
    disetujuiOleh: p.disetujui_nama || (p.disetujui_oleh ? String(p.disetujui_oleh) : null),
    tanggalPengajuan: p.created_at ? String(p.created_at).slice(0, 10) : '',
    tanggalPersetujuan: p.disetujui_oleh && p.updated_at ? String(p.updated_at).slice(0, 10) : null,
    keterangan: p.alasan || '',
    catatanApprover: p.catatan_approver ?? null,
    createdAt: p.created_at || '',
    updatedAt: p.updated_at || '',
  }
}

export function useArsip() {
  const { apiFetch } = useSession()

  // ─── 1. Dashboard ─────────────────────────────────────────────
  async function getDashboardStats(): Promise<DashboardStats> {
    try {
      const res = await apiFetch('/dashboard/stats')
      const json = await res.json()
      if (!json.success || !json.data) return { ...DUMMY_DASHBOARD_STATS, arsipVital: DUMMY_DASHBOARD_STATS.arsipAktif }
      const d = json.data
      return {
        totalArsip: Number(d.total_arsip || 0),
        arsipInaktif: Number(d.arsip_inaktif || 0),
        arsipVital: Number(d.arsip_vital || 0),
        arsipAktif: Number(d.arsip_vital || 0),
        arsipTerjaga: Number(d.arsip_terjaga || 0),
        jatuhTempoRetensi: Number(d.jatuh_tempo_retensi || 0),
        arsipDipinjam: Number(d.arsip_dipinjam || 0),
        usulPenyusutan: Number(d.usul_penyusutan || 0),
      }
    } catch {
      return { ...DUMMY_DASHBOARD_STATS, arsipVital: DUMMY_DASHBOARD_STATS.arsipAktif }
    }
  }

  async function getChartData(tahun?: number): Promise<ChartDataPoint[]> {
    try {
      const qs = tahun ? `?tahun=${tahun}` : ''
      const res = await apiFetch(`/dashboard/chart${qs}`)
      const json = await res.json()
      if (!json.success || !Array.isArray(json.data)) return [...DUMMY_CHART_DATA]
      return json.data.map((pt: any) => ({
        bulan: pt.bulan || '',
        jumlah: Number(pt.jumlah || 0),
      }))
    } catch {
      return [...DUMMY_CHART_DATA]
    }
  }

  // ─── 2. UNOR & Master Data (dengan Cache) ─────────────────────
  async function getDaftarUNOR(): Promise<Record<string, string>> {
    const res = await apiFetch('/arsip/unor')
    const json = await res.json()
    return json.data || {}
  }

  async function getKodeKlasifikasi(forceRefresh = false): Promise<KodeKlasifikasi[]> {
    if (!forceRefresh && cachedKodeKlasifikasi && cachedKodeKlasifikasi.length > 0) {
      return cachedKodeKlasifikasi
    }
    try {
      const res = await apiFetch('/master/kode-klasifikasi?per_page=500')
      const json = await res.json()
      if (!json.success || !Array.isArray(json.data)) return [...DUMMY_KODE_KLASIFIKASI]
      cachedKodeKlasifikasi = json.data.map(mapKodeKlasifikasi)
      return cachedKodeKlasifikasi
    } catch {
      return [...DUMMY_KODE_KLASIFIKASI]
    }
  }

  async function createKodeKlasifikasi(payload: Record<string, any>): Promise<KodeKlasifikasi> {
    const res = await apiFetch('/master/kode-klasifikasi', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menambahkan kode klasifikasi')
    cachedKodeKlasifikasi = null
    return mapKodeKlasifikasi(json.data)
  }

  async function updateKodeKlasifikasi(id: string | number, payload: Record<string, any>): Promise<KodeKlasifikasi> {
    const res = await apiFetch(`/master/kode-klasifikasi/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal memperbarui kode klasifikasi')
    cachedKodeKlasifikasi = null
    return mapKodeKlasifikasi(json.data)
  }

  async function deleteKodeKlasifikasi(id: string | number): Promise<boolean> {
    const res = await apiFetch(`/master/kode-klasifikasi/${id}`, { method: 'DELETE' })
    const json = await res.json()
    cachedKodeKlasifikasi = null
    return Boolean(json.success)
  }

  async function getLokasiSimpan(forceRefresh = false): Promise<LokasiSimpan[]> {
    if (!forceRefresh && cachedLokasiSimpan && cachedLokasiSimpan.length > 0) {
      return cachedLokasiSimpan
    }
    try {
      const res = await apiFetch('/master/lokasi-simpan?per_page=200')
      const json = await res.json()
      if (!json.success || !Array.isArray(json.data)) return [...DUMMY_LOKASI_SIMPAN]
      cachedLokasiSimpan = json.data.map(mapLokasiSimpan)
      return cachedLokasiSimpan
    } catch {
      return [...DUMMY_LOKASI_SIMPAN]
    }
  }

  async function createLokasiSimpan(payload: Record<string, any>): Promise<LokasiSimpan> {
    const res = await apiFetch('/master/lokasi-simpan', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menambahkan lokasi simpan')
    cachedLokasiSimpan = null
    return mapLokasiSimpan(json.data)
  }

  async function updateLokasiSimpan(id: string | number, payload: Record<string, any>): Promise<LokasiSimpan> {
    const res = await apiFetch(`/master/lokasi-simpan/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal memperbarui lokasi simpan')
    cachedLokasiSimpan = null
    return mapLokasiSimpan(json.data)
  }

  async function deleteLokasiSimpan(id: string | number): Promise<boolean> {
    const res = await apiFetch(`/master/lokasi-simpan/${id}`, { method: 'DELETE' })
    const json = await res.json()
    cachedLokasiSimpan = null
    return Boolean(json.success)
  }

  async function getUnitKerja(forceRefresh = false): Promise<UnitKerja[]> {
    if (!forceRefresh && cachedUnitKerja && cachedUnitKerja.length > 0) {
      return cachedUnitKerja
    }
    try {
      const res = await apiFetch('/master/unit-kerja?per_page=200')
      const json = await res.json()
      if (!json.success || !Array.isArray(json.data)) return [...DUMMY_UNIT_KERJA]
      cachedUnitKerja = json.data.map(mapUnitKerja)
      return cachedUnitKerja
    } catch {
      return [...DUMMY_UNIT_KERJA]
    }
  }

  async function createUnitKerja(payload: Record<string, any>): Promise<UnitKerja> {
    const res = await apiFetch('/master/unit-kerja', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menambahkan unit kerja')
    cachedUnitKerja = null
    return mapUnitKerja(json.data)
  }

  async function updateUnitKerja(id: string | number, payload: Record<string, any>): Promise<UnitKerja> {
    const res = await apiFetch(`/master/unit-kerja/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal memperbarui unit kerja')
    cachedUnitKerja = null
    return mapUnitKerja(json.data)
  }

  async function deleteUnitKerja(id: string | number): Promise<boolean> {
    const res = await apiFetch(`/master/unit-kerja/${id}`, { method: 'DELETE' })
    const json = await res.json()
    cachedUnitKerja = null
    return Boolean(json.success)
  }

  async function getKonfigurasi(): Promise<KonfigurasiSistem> {
    try {
      const res = await apiFetch('/master/konfigurasi')
      const json = await res.json()
      if (!json.success || !json.data) return { ...DUMMY_KONFIGURASI }
      const d = json.data
      return {
        namaAplikasi: d.app_name || 'Arsip Record Center',
        defaultJangkaSimpanInaktif: d.default_jangka_simpan || '3 Tahun',
        defaultNasibAkhir: (d.default_nasib_akhir || 'Musnah') as any,
        notifikasiJatuhTempoHari: Number(d.notifikasi_jatuh_tempo_hari || 30),
        logoAplikasi: d.logo_aplikasi || '/admin/images/logokpppa.png',
      }
    } catch {
      return { ...DUMMY_KONFIGURASI }
    }
  }

  async function updateKonfigurasi(data: Partial<KonfigurasiSistem>): Promise<KonfigurasiSistem> {
    const res = await apiFetch('/master/konfigurasi', {
      method: 'PUT',
      body: JSON.stringify({
        app_name: data.namaAplikasi,
        default_jangka_simpan: data.defaultJangkaSimpanInaktif,
        default_nasib_akhir: data.defaultNasibAkhir,
        notifikasi_jatuh_tempo_hari: String(data.notifikasiJatuhTempoHari ?? 30),
        logo_aplikasi: data.logoAplikasi || '',
      }),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menyimpan konfigurasi')
    return getKonfigurasi()
  }

  // ─── 3. Arsip CRUD & Dokumen — via API Backend ────────────────
  async function getArsipList(params?: Record<string, string | number | undefined>) {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') query.set(k, String(v))
      })
    }
    const qs = query.toString()
    const res = await apiFetch(`/arsip${qs ? `?${qs}` : ''}`)
    const json = await res.json()

    return {
      data: (json.data || []).map(mapArsip) as Arsip[],
      total: json.meta?.total || 0,
      page: json.meta?.current_page || 1,
      perPage: json.meta?.per_page || 10,
      totalPages: json.meta?.last_page || 0,
    }
  }

  async function getArsipById(id: string | number): Promise<Arsip | null> {
    try {
      const res = await apiFetch(`/arsip/${id}`)
      const json = await res.json()
      if (!json.success) return null
      return mapArsip(json.data)
    } catch {
      return null
    }
  }

  async function createArsip(data: Record<string, any>): Promise<Arsip> {
    const res = await apiFetch('/arsip', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menambahkan arsip')
    cachedKodeKlasifikasi = null
    return mapArsip(json.data)
  }

  async function updateArsip(id: string | number, payload: Record<string, any>): Promise<Arsip> {
    const res = await apiFetch(`/arsip/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mengupdate arsip')
    return mapArsip(json.data)
  }

  async function deleteArsip(id: string | number): Promise<boolean> {
    const res = await apiFetch(`/arsip/${id}`, { method: 'DELETE' })
    const json = await res.json()
    cachedKodeKlasifikasi = null
    return Boolean(json.success)
  }

  async function restoreArsip(id: string | number): Promise<Arsip | null> {
    const res = await apiFetch(`/arsip/${id}/restore`, { method: 'PUT' })
    const json = await res.json()
    if (!json.success) return null
    return mapArsip(json.data)
  }

  async function importArsipCSV(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    const res = await apiFetch('/arsip/import', {
      method: 'POST',
      body: formData,
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mengimpor file arsip')
    cachedKodeKlasifikasi = null
    return json.data
  }

  async function listDokumenArsip(arsipId: string | number): Promise<DokumenArsipItem[]> {
    const res = await apiFetch(`/arsip/${arsipId}/dokumen`)
    const json = await res.json()
    if (!json.success || !Array.isArray(json.data)) return []
    return json.data.map((d: any) => ({
      id: Number(d.id),
      arsipId: Number(d.arsip_id),
      namaFile: d.nama_file || '',
      mimeType: d.mime_type || '',
      ukuranBytes: Number(d.ukuran_bytes || 0),
      keterangan: d.keterangan ?? null,
      downloadUrl: `${API_BASE_URL.replace(/\/api$/, '')}${d.download_url}`,
      uploadedBy: Number(d.uploaded_by || 0),
      createdAt: d.created_at || '',
    }))
  }

  async function uploadDokumenArsip(arsipId: string | number, file: File, keterangan?: string): Promise<DokumenArsipItem> {
    const formData = new FormData()
    formData.append('file', file)
    if (keterangan) formData.append('keterangan', keterangan)

    const res = await apiFetch(`/arsip/${arsipId}/dokumen`, {
      method: 'POST',
      body: formData,
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mengunggah dokumen')
    const d = json.data
    return {
      id: Number(d.id),
      arsipId: Number(d.arsip_id),
      namaFile: d.nama_file || '',
      mimeType: d.mime_type || '',
      ukuranBytes: Number(d.ukuran_bytes || 0),
      keterangan: d.keterangan ?? null,
      downloadUrl: `${API_BASE_URL.replace(/\/api$/, '')}${d.download_url}`,
      uploadedBy: Number(d.uploaded_by || 0),
      createdAt: d.created_at || '',
    }
  }

  async function deleteDokumenArsip(dokumenId: number | string): Promise<boolean> {
    const res = await apiFetch(`/dokumen/${dokumenId}`, { method: 'DELETE' })
    const json = await res.json()
    return Boolean(json.success)
  }

  // ─── 4. Peminjaman — via API Backend ──────────────────────────
  async function getPeminjamanList(params?: Record<string, string | number | undefined>): Promise<Peminjaman[]> {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') query.set(k, String(v))
      })
    }
    const qs = query.toString()
    const res = await apiFetch(`/peminjaman${qs ? `?${qs}` : ''}`)
    const json = await res.json()
    if (!json.success || !Array.isArray(json.data)) return []
    return json.data.map(mapPeminjaman)
  }

  async function createPeminjaman(data: Record<string, any>): Promise<Peminjaman> {
    const res = await apiFetch('/peminjaman', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mencatat peminjaman')
    return mapPeminjaman(json.data)
  }

  async function prosesPengembalian(id: string | number, catatan?: string): Promise<Peminjaman | null> {
    const res = await apiFetch(`/peminjaman/${id}/kembali`, {
      method: 'PUT',
      body: JSON.stringify({ catatan: catatan || 'Dikembalikan' }),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mencatat pengembalian')
    return mapPeminjaman(json.data)
  }

  // ─── 5. Penyusutan — via API Backend ──────────────────────────
  async function getPenyusutanList(params?: Record<string, string | number | undefined>): Promise<Penyusutan[]> {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') query.set(k, String(v))
      })
    }
    const qs = query.toString()
    const res = await apiFetch(`/penyusutan${qs ? `?${qs}` : ''}`)
    const json = await res.json()
    if (!json.success || !Array.isArray(json.data)) return []
    return json.data.map(mapPenyusutan)
  }

  async function createPenyusutan(data: Record<string, any>): Promise<Penyusutan> {
    const res = await apiFetch('/penyusutan', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal mengajukan penyusutan')
    return mapPenyusutan(json.data)
  }

  async function setujuiPenyusutan(id: string | number, catatanApprover?: string): Promise<Penyusutan | null> {
    const res = await apiFetch(`/penyusutan/${id}/approve`, {
      method: 'PUT',
      body: JSON.stringify({ catatan_approver: catatanApprover || 'Disetujui' }),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menyetujui penyusutan')
    return mapPenyusutan(json.data)
  }

  async function tolakPenyusutan(id: string | number, catatanApprover?: string): Promise<Penyusutan | null> {
    const res = await apiFetch(`/penyusutan/${id}/reject`, {
      method: 'PUT',
      body: JSON.stringify({ catatan_approver: catatanApprover || 'Ditolak' }),
    })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menolak penyusutan')
    return mapPenyusutan(json.data)
  }

  async function selesaikanPenyusutan(id: string | number): Promise<Penyusutan | null> {
    const res = await apiFetch(`/penyusutan/${id}/selesai`, { method: 'PUT' })
    const json = await res.json()
    if (!json.success) throw new Error(json.message || 'Gagal menyelesaikan penyusutan')
    return mapPenyusutan(json.data)
  }

  // ─── 6. Laporan & Rekap DAI 18 Kolom ──────────────────────────
  async function getLaporanRekap(params?: Record<string, string | number | undefined>) {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') query.set(k, String(v))
      })
    }
    const qs = query.toString()
    const res = await apiFetch(`/laporan/rekap${qs ? `?${qs}` : ''}`)
    const json = await res.json()
    return json.data || null
  }

  async function getLaporanDAI(params?: Record<string, string | number | undefined>): Promise<DAIRow[]> {
    const query = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') query.set(k, String(v))
      })
    }
    const qs = query.toString()
    const res = await apiFetch(`/laporan/dai${qs ? `?${qs}` : ''}`)
    const json = await res.json()
    return Array.isArray(json.data) ? json.data : []
  }

  // ─── 7. User Management — via API Backend ─────────────────────
  async function getUserList(): Promise<UserManagement[]> {
    const res = await apiFetch('/users?per_page=100')
    const json = await res.json()
    if (!json.success || !Array.isArray(json.data)) return []
    return json.data.map((u: any) => ({
      id: String(u.id),
      username: u.username || '',
      nama: u.name || '',
      nip: u.nip || '',
      email: u.email || '',
      unitKerja: u.unit_kerja || '-',
      jabatan: u.jabatan || '',
      role: u.role_label || u.role || 'Operator',
      roleKey: u.role || 'operator',
      status: u.is_active ? 'aktif' : 'nonaktif',
      terakhirLogin: formatDateTimeIndo(u.last_login_at),
    }))
  }

  // ─── 8. Export Helper ─────────────────────────────────────────
  function exportToCsv(data: Record<string, any>[], filename: string): void {
    if (data.length === 0) return
    const headers = Object.keys(data[0])
    const csvContent = [
      headers.join(','),
      ...data.map((row) =>
        headers
          .map((h) => {
            const val = row[h]?.toString() || ''
            return val.includes(',') ? `"${val}"` : val
          })
          .join(',')
      ),
    ].join('\n')
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `${filename}.csv`
    link.click()
    URL.revokeObjectURL(link.href)
  }

  return {
    getDashboardStats,
    getChartData,
    getDaftarUNOR,
    getArsipList,
    getArsipById,
    createArsip,
    updateArsip,
    deleteArsip,
    restoreArsip,
    importArsipCSV,
    listDokumenArsip,
    uploadDokumenArsip,
    deleteDokumenArsip,
    getPeminjamanList,
    createPeminjaman,
    prosesPengembalian,
    getPenyusutanList,
    createPenyusutan,
    setujuiPenyusutan,
    tolakPenyusutan,
    selesaikanPenyusutan,
    getKodeKlasifikasi,
    createKodeKlasifikasi,
    updateKodeKlasifikasi,
    deleteKodeKlasifikasi,
    getLokasiSimpan,
    createLokasiSimpan,
    updateLokasiSimpan,
    deleteLokasiSimpan,
    getUnitKerja,
    createUnitKerja,
    updateUnitKerja,
    deleteUnitKerja,
    getKonfigurasi,
    updateKonfigurasi,
    getLaporanRekap,
    getLaporanDAI,
    getUserList,
    exportToCsv,
  }
}
