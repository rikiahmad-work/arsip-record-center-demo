<template>
  <TemplatePageTitle title="Pengelolaan Arsip" :breadcrumbs="[{ label: 'Pengelolaan Arsip' }]">
    <template #action>
      <div class="d-flex gap-2 flex-wrap">
        <button
          v-if="canWrite()"
          type="button"
          class="btn btn-outline-primary d-inline-flex align-items-center gap-1"
          @click="showImportModal = true"
        >
          <i class="ti ti-upload"></i>
          <span>Import Excel/CSV</span>
        </button>
        <button
          type="button"
          class="btn btn-success d-inline-flex align-items-center gap-1"
          title="Unduh file Excel (.xlsx) dengan format baku 18 kolom DAI KemenPPPA"
          @click="handleExport('xlsx')"
        >
          <i class="ti ti-file-spreadsheet"></i>
          <span>Export Excel (DAI)</span>
        </button>
        <button
          type="button"
          class="btn btn-outline-success d-inline-flex align-items-center gap-1"
          title="Unduh format CSV"
          @click="handleExport('csv')"
        >
          <i class="ti ti-download"></i>
          <span>CSV</span>
        </button>
        <button
          v-if="canWrite()"
          type="button"
          class="btn btn-primary d-inline-flex align-items-center gap-1"
          @click="$router.push('/pengelolaan-arsip/tambah')"
        >
          <i class="ti ti-plus"></i>
          <span>Tambah Arsip</span>
        </button>
      </div>
    </template>
  </TemplatePageTitle>

  <!-- 4 Stat Cards Ringkas -->
  <div class="row mb-4">
    <div v-for="s in statCards" :key="s.label" class="col-sm-6 col-xl-3">
      <div class="card mb-0" :class="s.bg">
        <div class="card-body py-3">
          <div class="d-flex justify-content-between align-items-center">
            <div>
              <h4 class="mb-0" :class="s.textClass">{{ s.value }}</h4>
              <small :class="s.textClass + '-50'">{{ s.label }}</small>
            </div>
            <i :class="[s.icon, s.textClass]" style="font-size: 1.75rem;"></i>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="card">
    <div class="card-body">
      <!-- Nav Tabs Sesuai Pola Manajemen User -->
      <ul class="nav nav-tabs mb-4" role="tablist">
        <li v-for="tab in outputTabs" :key="tab.key" class="nav-item" role="presentation">
          <button
            class="nav-link d-inline-flex align-items-center"
            :class="{ active: activeOutputTab === tab.key }"
            type="button"
            @click="selectOutputTab(tab.key)"
          >
            <i :class="tab.icon" class="me-1"></i>
            <span>{{ tab.label }}</span>
            <span :class="['badge rounded-pill ms-2', tab.badgeClass]">{{ tab.count }}</span>
          </button>
        </li>
      </ul>
          <!-- Bar Filter Dedicated (Unit Kerja, Kode Klasifikasi, Kategori, SKKAD, & Status Operasional) -->
          <div class="p-3 bg-light rounded border mb-3">
            <div class="row g-2 align-items-end">
              <div class="col-12 col-md-3">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-building me-1"></i>Unit Kerja / Pencipta Arsip
                </label>
                <FormSelect
                  v-model="filter.unitKerja"
                  :options="unitSelectOptions"
                  placeholder="Semua Unit Kerja"
                  search-placeholder="Cari unit kerja..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-12 col-md-3">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-tag me-1"></i>Kode Klasifikasi (JRA)
                </label>
                <FormSelect
                  v-model="filter.kode"
                  :options="kodeSelectOptions"
                  placeholder="Semua Kode Klasifikasi"
                  search-placeholder="Cari kode atau nama klasifikasi..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-6 col-md-2">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-category me-1"></i>Kategori Arsip
                </label>
                <FormSelect
                  v-model="filter.kategori"
                  :options="kategoriFilterOptions"
                  placeholder="Semua Kategori"
                  search-placeholder="Cari kategori..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-6 col-md-2">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-activity me-1"></i>Status Retensi / Pinjam
                </label>
                <FormSelect
                  v-model="filter.statusOperasional"
                  :options="statusOperasionalOptions"
                  placeholder="Semua Status"
                  size="sm"
                  clearable
                  no-margin
                />
              </div>

              <div class="col-12 col-md-2 d-flex gap-1">
                <button
                  type="button"
                  class="btn btn-sm w-100 d-inline-flex align-items-center justify-content-center gap-1"
                  :class="hasActiveFilter ? 'btn-outline-danger' : 'btn-outline-secondary opacity-50'"
                  :disabled="!hasActiveFilter"
                  @click="resetFilter"
                >
                  <i class="ti ti-filter-off"></i>
                  <span>Reset Filter</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Initial Loading Skeleton (Hanya Saat Pertama Kali Buka Halaman) -->
          <div v-if="initialLoading" class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Memuat...</span>
            </div>
            <p class="mt-2 text-muted">Memuat data arsip dari Record Center...</p>
          </div>

          <!-- DataTables Tanpa Reload Halaman + Sticky Kolom Aksi -->
          <ModulesPengelolaanArsipGrid
            v-else
            :list="filteredArsipList"
            :can-edit="canWrite()"
            :can-delete="isAdmin()"
            @detail="lihatDetail"
            @edit="editArsip"
            @hapus="konfirmasiHapus"
          />
        </div>
      </div>

  <!-- Modal Import CSV/Excel -->
  <div v-if="showImportModal" class="modal fade show" tabindex="-1" style="display:block;background:rgba(0,0,0,0.5);">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Import Data Arsip (CSV / Excel)</h5>
          <button type="button" class="btn-close" aria-label="Close" @click="showImportModal = false"></button>
        </div>
        <div class="modal-body py-4">
          <div class="text-center mb-3">
            <i class="ti ti-file-spreadsheet" style="font-size: 48px; color: #16a34a;"></i>
            <p class="mt-2 mb-1 fw-semibold">Unggah file Excel (.xlsx) atau CSV Daftar Arsip Inaktif (DAI)</p>
            <small class="text-muted d-block mb-2">
              Mendukung langsung format <strong>18 Kolom DAI Record Center</strong>. Jangka Simpan Inaktif, Nasib Akhir, dan SKKAD akan otomatis terisi dari Master Kode Klasifikasi.
            </small>
            <div class="d-flex justify-content-center gap-2">
              <button type="button" class="btn btn-sm btn-success" @click="handleExport('xlsx')">
                <i class="ti ti-file-spreadsheet me-1"></i>Unduh Template Excel (.xlsx)
              </button>
              <button type="button" class="btn btn-sm btn-light-success" @click="downloadTemplateCSV">
                <i class="ti ti-file-download me-1"></i>Template CSV
              </button>
            </div>
          </div>
          <input
            type="file"
            class="form-control"
            accept=".xlsx,.csv"
            @change="onImportFileChange"
          >
        </div>
        <div class="modal-footer justify-content-end">
          <button class="btn btn-outline-secondary" @click="showImportModal = false">Batal</button>
          <button class="btn btn-primary" :disabled="!importFile || importing" @click="submitImport">
            <span v-if="importing" class="spinner-border spinner-border-sm me-1"></span>
            Import Sekarang
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Konfirmasi Hapus -->
  <div v-if="hapusItem" class="modal fade show" tabindex="-1" style="display:block;background:rgba(0,0,0,0.5);">
    <div class="modal-dialog modal-dialog-centered modal-sm">
      <div class="modal-content">
        <div class="modal-header bg-light-danger">
          <h5 class="modal-title text-danger">Konfirmasi Hapus</h5>
          <button type="button" class="btn-close" aria-label="Close" @click="hapusItem = null"></button>
        </div>
        <div class="modal-body text-center p-4">
          <i class="ti ti-trash text-danger mb-3" style="font-size: 50px;"></i>
          <h6>Hapus arsip <strong>{{ hapusItem.noItemArsip }}</strong>?</h6>
          <p class="text-muted small mb-0">Data arsip akan dipindahkan ke tempat sampah (soft delete) dan dapat dipulihkan oleh Administrator.</p>
        </div>
        <div class="modal-footer justify-content-center">
          <button class="btn btn-secondary" @click="hapusItem = null">Batal</button>
          <button class="btn btn-danger" @click="hapusArsip">Ya, Hapus</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Arsip, KodeKlasifikasi, UnitKerja } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Pengelolaan Arsip | Arsip Record Center KemenPPPA' })

const route = useRoute()
const router = useRouter()
const { canWrite, isAdmin, apiFetch } = useSession()
const { getArsipList, deleteArsip, getKodeKlasifikasi, getUnitKerja, importArsipCSV } = useArsip()
const { showToast } = useToast()

const allArsipList = ref<Arsip[]>([])
const kodeList = ref<KodeKlasifikasi[]>([])
const unitList = ref<UnitKerja[]>([])
const initialLoading = ref(true)
const refreshing = ref(false)
const activeOutputTab = ref<'all' | 'inaktif' | 'vital' | 'usul_musnah' | 'usul_serah' | 'jatuh_tempo'>('all')

const filter = reactive({
  unitKerja: (route.query.unit_kerja as string) || '',
  kode: (route.query.kode as string) || '',
  kategori: (route.query.kategori as string) || '',
  skkad: (route.query.skkad as string) || '',
  statusOperasional: (route.query.status as string) || '',
})

watch(
  () => route.query,
  (q) => {
    if (typeof q.unit_kerja === 'string') filter.unitKerja = q.unit_kerja
    if (typeof q.kode === 'string') filter.kode = q.kode
    if (typeof q.kategori === 'string') filter.kategori = q.kategori
    if (typeof q.skkad === 'string') filter.skkad = q.skkad
    if (typeof q.status === 'string') filter.statusOperasional = q.status
  }
)

const outputTabs = computed(() => [
  {
    key: 'all' as const,
    label: 'Semua Arsip',
    icon: 'ti ti-files',
    badgeClass: 'bg-light-primary text-primary',
    count: allArsipList.value.length,
  },
  {
    key: 'inaktif' as const,
    label: 'Arsip Inaktif',
    icon: 'ti ti-archive',
    badgeClass: 'bg-light-success text-success',
    count: allArsipList.value.filter((a) => a.kategoriArsip.toUpperCase() === 'INAKTIF').length,
  },
  {
    key: 'vital' as const,
    label: 'Vital & Terjaga',
    icon: 'ti ti-shield-lock',
    badgeClass: 'bg-light-warning text-warning',
    count: allArsipList.value.filter((a) => ['VITAL', 'TERJAGA'].includes(a.kategoriArsip.toUpperCase())).length,
  },
  {
    key: 'usul_musnah' as const,
    label: 'Usul Musnah',
    icon: 'ti ti-trash',
    badgeClass: 'bg-light-danger text-danger',
    count: allArsipList.value.filter(
      (a) =>
        a.kategoriArsip.toUpperCase() === 'USUL MUSNAH' ||
        (a.isJatuhTempo && !String(a.nasibAkhir || a.jangkaSimpanNasibAkhir || '').toUpperCase().includes('PERMANEN'))
    ).length,
  },
  {
    key: 'usul_serah' as const,
    label: 'Usul Serah',
    icon: 'ti ti-transfer-in',
    badgeClass: 'bg-light-info text-info',
    count: allArsipList.value.filter(
      (a) =>
        a.kategoriArsip.toUpperCase() === 'USUL SERAH' ||
        (a.isJatuhTempo && String(a.nasibAkhir || a.jangkaSimpanNasibAkhir || '').toUpperCase().includes('PERMANEN'))
    ).length,
  },
  {
    key: 'jatuh_tempo' as const,
    label: 'Jatuh Tempo (31 Des)',
    icon: 'ti ti-clock-exclamation',
    badgeClass: 'bg-danger text-white',
    count: allArsipList.value.filter((a) => a.isJatuhTempo).length,
  },
])

function selectOutputTab(tabKey: typeof activeOutputTab.value) {
  activeOutputTab.value = tabKey
}

const hasActiveFilter = computed(() =>
  Boolean(
    activeOutputTab.value !== 'all' ||
      filter.unitKerja ||
      filter.kode ||
      filter.kategori ||
      filter.skkad ||
      filter.statusOperasional
  )
)

const unitSelectOptions = computed(() =>
  unitList.value.map((u) => ({
    value: u.namaUnit,
    label: `${u.kodeUnit} — ${u.namaUnit}`,
  }))
)

const kodeSelectOptions = computed(() =>
  kodeList.value.map((k) => ({
    value: k.kode,
    label: `${k.kode} — ${k.nama}`,
  }))
)

const kategoriFilterOptions = [
  { value: 'Inaktif', label: 'Inaktif' },
  { value: 'Vital', label: 'Vital' },
  { value: 'Terjaga', label: 'Terjaga' },
  { value: 'Usul Musnah', label: 'Usul Musnah' },
  { value: 'Usul Serah', label: 'Usul Serah' },
]

const statusOperasionalOptions = [
  { value: 'jatuh_tempo', label: 'Jatuh Tempo Retensi (31 Des)' },
  { value: 'dipinjam', label: 'Sedang Dipinjam' },
  { value: 'penyusutan', label: 'Dalam Proses Penyusutan' },
  { value: 'final', label: 'Status: Final' },
  { value: 'draft', label: 'Status: Draft' },
]

// Filter instan di sisi client/DataTables tanpa berkedip atau mereload seluruh halaman
const filteredArsipList = computed(() => {
  return allArsipList.value.filter((item) => {
    const katUpper = item.kategoriArsip.toUpperCase()
    const isPerm = String(item.nasibAkhir || item.jangkaSimpanNasibAkhir || '').toUpperCase().includes('PERMANEN')

    // Filter dari Tab Pengelompokan ANRI
    if (activeOutputTab.value === 'inaktif' && katUpper !== 'INAKTIF') return false
    if (activeOutputTab.value === 'vital' && !['VITAL', 'TERJAGA'].includes(katUpper)) return false
    if (activeOutputTab.value === 'usul_musnah' && !(katUpper === 'USUL MUSNAH' || (item.isJatuhTempo && !isPerm))) return false
    if (activeOutputTab.value === 'usul_serah' && !(katUpper === 'USUL SERAH' || (item.isJatuhTempo && isPerm))) return false
    if (activeOutputTab.value === 'jatuh_tempo' && !item.isJatuhTempo) return false

    // Filter Dropdown
    if (filter.unitKerja && item.penciptaArsip.toLowerCase() !== filter.unitKerja.toLowerCase()) return false
    if (filter.kode && item.kodeKlasifikasi !== filter.kode) return false
    if (filter.kategori && item.kategoriArsip.toLowerCase() !== filter.kategori.toLowerCase()) return false
    if (filter.skkad && (item.skkad || 'Biasa').toLowerCase() !== filter.skkad.toLowerCase()) return false

    if (filter.statusOperasional === 'jatuh_tempo' && !item.isJatuhTempo) return false
    if (filter.statusOperasional === 'dipinjam' && !item.sedangDipinjam) return false
    if (filter.statusOperasional === 'penyusutan' && !item.dalamPenyusutan) return false
    if (filter.statusOperasional === 'final' && item.status !== 'final') return false
    if (filter.statusOperasional === 'draft' && item.status !== 'draft') return false

    return true
  })
})

async function fetchData(isBackground = false) {
  if (isBackground) {
    refreshing.value = true
  } else {
    initialLoading.value = true
  }
  try {
    const result = await getArsipList({
      page: 1,
      per_page: 500,
      sort_by: 'id',
      sort_order: 'desc',
    })
    allArsipList.value = result.data
  } catch (err: any) {
    showToast(err?.message || 'Gagal memuat data arsip', 'danger')
  } finally {
    initialLoading.value = false
    refreshing.value = false
  }
}

onMounted(async () => {
  const [, kodes, units] = await Promise.all([fetchData(false), getKodeKlasifikasi(), getUnitKerja()])
  kodeList.value = kodes
  unitList.value = units
})

function resetFilter() {
  activeOutputTab.value = 'all'
  filter.unitKerja = ''
  filter.kode = ''
  filter.kategori = ''
  filter.skkad = ''
  filter.statusOperasional = ''
  if (Object.keys(route.query).length > 0) {
    router.replace({ path: '/pengelolaan-arsip', query: {} })
  }
}

const statCards = computed(() => [
  {
    label: 'Total Arsip Terdaftar',
    value: allArsipList.value.length,
    bg: 'bg-primary text-white',
    textClass: 'text-white',
    icon: 'ti ti-files',
  },
  {
    label: 'Arsip Inaktif',
    value: allArsipList.value.filter((a) => a.kategoriArsip.toUpperCase() === 'INAKTIF').length,
    bg: 'bg-success text-white',
    textClass: 'text-white',
    icon: 'ti ti-archive',
  },
  {
    label: 'Arsip Vital & Terjaga',
    value: allArsipList.value.filter((a) => ['VITAL', 'TERJAGA'].includes(a.kategoriArsip.toUpperCase())).length,
    bg: 'bg-warning text-white',
    textClass: 'text-white',
    icon: 'ti ti-shield-lock',
  },
  {
    label: 'Jatuh Tempo Retensi (31 Des)',
    value: allArsipList.value.filter((a) => a.isJatuhTempo).length,
    bg: 'bg-danger text-white',
    textClass: 'text-white',
    icon: 'ti ti-clock-exclamation',
  },
])

const showImportModal = ref(false)
const importFile = ref<File | null>(null)
const importing = ref(false)
const hapusItem = ref<Arsip | null>(null)

function onImportFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  importFile.value = input.files?.[0] || null
}

function downloadTemplateCSV() {
  const header = 'NO.,NO. ITEM ARSIP,KODE KLASIFIKASI,URAIAN INFORMASI ARSIP,JENIS ARSIP,KURUN WAKTU,TINGKAT PERKEMBANGAN,JUMLAH,KETERANGAN,NOMOR DEFINITIF,LOKASI SIMPAN,NO.BOKS,NO. RAK,JANGKA SIMPAN INAKTIF,NASIB AKHIR,SKKAD,KATEGORI ARSIP,PENCIPTA ARSIP\n'
  const sample = '1,1,DT.10,"Balasan Surat Permohonan Informasi Anak Penyandang Disabilitas pada SIMFONI PPPA Nomor B-1436/Setmen.Birodatin/DT.10/9/2021",KONVENSIONAL,2021 - 2022,ASLI,1 Lembar,KONDISI BAIK,,RECORDS CENTER,,,3 Tahun,Musnah,Biasa,INAKTIF,Biro Data dan Informasi\n'
  const blob = new Blob([header + sample], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'template-dai-18-kolom-record-center.csv'
  link.click()
  URL.revokeObjectURL(url)
}

async function submitImport() {
  if (!importFile.value) return
  importing.value = true
  try {
    const res = await importArsipCSV(importFile.value)
    showToast(`Import selesai: ${res.imported_rows} berhasil, ${res.failed_rows} gagal`, 'success')
    showImportModal.value = false
    importFile.value = null
    await fetchData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal mengimpor file arsip', 'danger')
  } finally {
    importing.value = false
  }
}

function lihatDetail(item: Arsip) {
  router.push(`/pengelolaan-arsip/${item.id}`)
}

function editArsip(item: Arsip) {
  router.push(`/pengelolaan-arsip/${item.id}/edit`)
}

function konfirmasiHapus(item: Arsip) {
  hapusItem.value = item
}

async function hapusArsip() {
  if (!hapusItem.value) return
  try {
    const ok = await deleteArsip(hapusItem.value.id)
    if (ok) {
      showToast('Arsip berhasil dihapus', 'success')
      await fetchData(true)
    } else {
      showToast('Gagal menghapus arsip', 'danger')
    }
  } catch (err: any) {
    showToast(err?.message || 'Gagal menghapus arsip', 'danger')
  } finally {
    hapusItem.value = null
  }
}

async function handleExport(format: 'xlsx' | 'csv' = 'xlsx') {
  try {
    const qs = new URLSearchParams()
    qs.set('format', format)
    if (filter.unitKerja) qs.set('unit_kerja', filter.unitKerja)
    if (filter.kode) qs.set('kode', filter.kode)
    if (filter.kategori) {
      qs.set('kategori', filter.kategori)
    } else if (activeOutputTab.value === 'inaktif') {
      qs.set('kategori', 'Inaktif')
    } else if (activeOutputTab.value === 'usul_musnah') {
      qs.set('kategori', 'Usul Musnah')
    } else if (activeOutputTab.value === 'usul_serah') {
      qs.set('kategori', 'Usul Serah')
    }
    if (filter.skkad) qs.set('skkad', filter.skkad)

    const res = await apiFetch(`/arsip/export?${qs.toString()}`)
    if (!res.ok) {
      let errMsg = 'Gagal mengekspor data arsip'
      try {
        const errJson = await res.json()
        if (errJson.message) errMsg = errJson.message
      } catch (_) {}
      showToast(errMsg, 'danger')
      return
    }

    const blob = await res.blob()
    const mimeType = format === 'xlsx'
      ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      : 'text/csv;charset=utf-8;'
    const fileBlob = new Blob([blob], { type: mimeType })
    const url = URL.createObjectURL(fileBlob)
    const link = document.createElement('a')
    link.href = url
    const ext = format === 'xlsx' ? 'xlsx' : 'csv'
    link.download = `DAFTAR-ARSIP-INAKTIF-KEMENPPPA-${new Date().toISOString().slice(0, 10)}.${ext}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    showToast(
      format === 'xlsx'
        ? 'Daftar Arsip Inaktif (Excel 18 Kolom DAI) berhasil diunduh'
        : 'Daftar Arsip Inaktif (CSV) berhasil diunduh',
      'success'
    )
  } catch (err: any) {
    showToast(err?.message || 'Gagal mengekspor data arsip', 'danger')
  }
}
</script>
