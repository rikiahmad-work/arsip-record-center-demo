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
          class="btn btn-outline-success d-inline-flex align-items-center gap-1"
          @click="handleExport"
        >
          <i class="ti ti-download"></i>
          <span>Export Excel/CSV</span>
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

  <div class="row">
    <div class="col-sm-12">
      <div class="card">
        <div class="card-body">
          <!-- Bar Filter Dedicated di Bawah Header (Grid 12 Kolom Stabil + Searchable Dropdown Select2-style) -->
          <div class="p-3 bg-light rounded border mb-3">
            <div class="row g-2 align-items-end">
              <div class="col-12 col-md-4">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-tag me-1"></i>Kode Klasifikasi (JRA)
                </label>
                <FormSelect
                  v-model="filter.kode"
                  :options="kodeSelectOptions"
                  placeholder="Semua Kode Klasifikasi (435 Kode)"
                  search-placeholder="Cari kode atau nama klasifikasi..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-6 col-md-3">
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

              <div class="col-6 col-md-3">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-shield-lock me-1"></i>Klasifikasi Keamanan (SKKAD)
                </label>
                <FormSelect
                  v-model="filter.skkad"
                  :options="skkadFilterOptions"
                  placeholder="Semua SKKAD"
                  search-placeholder="Cari SKKAD..."
                  size="sm"
                  searchable
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
            <p class="mt-2 mb-1 fw-semibold">Unggah file CSV Daftar Arsip Record Center</p>
            <small class="text-muted d-block mb-2">
              Jangka Simpan Inaktif, Nasib Akhir, dan SKKAD akan otomatis terisi dari Master 435 Kode Klasifikasi.
            </small>
            <button type="button" class="btn btn-sm btn-light-success" @click="downloadTemplateCSV">
              <i class="ti ti-file-download me-1"></i>Unduh Template CSV Standar
            </button>
          </div>
          <input
            type="file"
            class="form-control"
            accept=".csv"
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
import type { Arsip, KodeKlasifikasi } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Pengelolaan Arsip | Arsip Record Center KemenPPPA' })

const route = useRoute()
const router = useRouter()
const { canWrite, isAdmin, apiFetch } = useSession()
const { getArsipList, deleteArsip, getKodeKlasifikasi, importArsipCSV } = useArsip()
const { showToast } = useToast()

const allArsipList = ref<Arsip[]>([])
const kodeList = ref<KodeKlasifikasi[]>([])
const initialLoading = ref(true)
const refreshing = ref(false)

const filter = reactive({
  kode: (route.query.kode as string) || '',
  kategori: (route.query.kategori as string) || '',
  skkad: (route.query.skkad as string) || '',
})

watch(
  () => route.query,
  (q) => {
    if (typeof q.kode === 'string') filter.kode = q.kode
    if (typeof q.kategori === 'string') filter.kategori = q.kategori
    if (typeof q.skkad === 'string') filter.skkad = q.skkad
  }
)

const hasActiveFilter = computed(() =>
  Boolean(filter.kode || filter.kategori || filter.skkad)
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

const skkadFilterOptions = [
  { value: 'Biasa', label: 'Biasa' },
  { value: 'Terbatas', label: 'Terbatas' },
  { value: 'Rahasia', label: 'Rahasia' },
]

// Filter instan di sisi client/DataTables tanpa berkedip atau mereload seluruh halaman
const filteredArsipList = computed(() => {
  return allArsipList.value.filter((item) => {
    if (filter.kode && item.kodeKlasifikasi !== filter.kode) return false
    if (filter.kategori && item.kategoriArsip.toLowerCase() !== filter.kategori.toLowerCase()) return false
    if (filter.skkad && (item.skkad || 'Biasa').toLowerCase() !== filter.skkad.toLowerCase()) return false
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
  const [, kodes] = await Promise.all([fetchData(false), getKodeKlasifikasi()])
  kodeList.value = kodes
})

function resetFilter() {
  filter.kode = ''
  filter.kategori = ''
  filter.skkad = ''
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
    label: 'Usul Musnah / Serah',
    value: allArsipList.value.filter((a) => ['USUL MUSNAH', 'USUL SERAH'].includes(a.kategoriArsip.toUpperCase())).length,
    bg: 'bg-danger text-white',
    textClass: 'text-white',
    icon: 'ti ti-clock',
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
  const header = 'no_item_arsip,kode_klasifikasi,uraian_informasi,kurun_waktu,tingkat_perkembangan,jumlah,keterangan,nomor_definitif,lokasi_simpan,no_rak,no_boks,no_folder,jenis_arsip,kategori_arsip,pencipta_arsip\n'
  const sample = 'ARS-2026-010,HK.00.00,Berkas Penyusunan Peraturan Menteri PPPA Tahun 2024,2024,Asli,1 Folder,Baik,DEF-PPPA-2026-010,Record Center Lt. 2,R-01,B-01,F-01,Konvensional,Inaktif,Biro Hukum dan Humas\n'
  const blob = new Blob([header + sample], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'template-import-arsip-record-center.csv'
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

async function handleExport() {
  try {
    const qs = new URLSearchParams()
    if (filter.kode) qs.set('kode', filter.kode)
    if (filter.kategori) qs.set('kategori', filter.kategori)
    if (filter.skkad) qs.set('skkad', filter.skkad)

    const res = await apiFetch(`/arsip/export${qs.toString() ? `?${qs.toString()}` : ''}`)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `daftar-arsip-record-center-${new Date().toISOString().slice(0, 10)}.csv`
    link.click()
    URL.revokeObjectURL(url)
    showToast('Daftar arsip berhasil diekspor', 'success')
  } catch {
    showToast('Gagal mengekspor data arsip', 'danger')
  }
}
</script>
