<template>
  <TemplatePageTitle title="Peminjaman Arsip" :breadcrumbs="[{ label: 'Peminjaman Arsip' }]">
    <template #action>
      <button
        v-if="canWrite()"
        type="button"
        class="btn btn-primary d-inline-flex align-items-center gap-1"
        @click="openTambah"
      >
        <i class="ti ti-plus"></i>
        <span>Catat Peminjaman</span>
      </button>
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
          <!-- Bar Filter Dedicated di Bawah Header -->
          <div class="p-3 bg-light rounded border mb-3">
            <div class="row g-2 align-items-end">
              <div class="col-12 col-md-5">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-building me-1"></i>Unit Kerja Peminjam
                </label>
                <FormSelect
                  v-model="filter.unitKerja"
                  :options="unitFilterOptions"
                  placeholder="Semua Unit Kerja"
                  search-placeholder="Cari unit kerja..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-12 col-md-5">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-status-change me-1"></i>Status Peminjaman
                </label>
                <FormSelect
                  v-model="filter.status"
                  :options="statusFilterOptions"
                  placeholder="Semua Status"
                  search-placeholder="Cari status..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-12 col-md-2">
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

          <div v-if="initialLoading" class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Memuat...</span>
            </div>
            <p class="mt-2 text-muted">Memuat data peminjaman arsip...</p>
          </div>

          <ModulesPeminjamanArsipGrid
            v-else
            :list="filteredList"
            :can-write="canWrite()"
            @detail="detailPeminjaman"
            @kembali="openKembaliModal"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Tambah Peminjaman -->
  <ModulesPeminjamanArsipForm
    v-if="showTambah"
    :arsip-list="arsipList"
    :unit-list="unitList"
    :submitting="submitting"
    @close="showTambah = false"
    @save="simpanPeminjaman"
  />

  <!-- Modal Proses Pengembalian -->
  <UiAppModal v-if="kembaliTarget" title="Konfirmasi Pengembalian Arsip" size="md" @close="kembaliTarget = null">
    <div class="mb-3">
      <div class="p-3 bg-light rounded border mb-3">
        <div class="small text-muted">No. Tiket / Peminjam</div>
        <div class="fw-bold text-dark">{{ kembaliTarget.noTiket }} — {{ kembaliTarget.peminjam }}</div>
        <div class="small text-muted mt-2">Berkas Arsip</div>
        <div class="fw-medium">{{ kembaliTarget.noItemArsip }} ({{ kembaliTarget.kodeKlasifikasi }})</div>
        <div class="small text-muted">{{ kembaliTarget.uraianArsip }}</div>
      </div>

      <label class="form-label fw-semibold">Catatan Kondisi Saat Dikembalikan</label>
      <textarea
        v-model="catatanKembali"
        class="form-control"
        rows="2"
        placeholder="Contoh: Berkas fisik lengkap dan kondisi baik"
      ></textarea>
    </div>
    <template #footer>
      <button type="button" class="btn btn-outline-secondary" @click="kembaliTarget = null">Batal</button>
      <button type="button" class="btn btn-success" :disabled="submitting" @click="submitPengembalian">
        <span v-if="submitting" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="ti ti-check me-1"></i>Konfirmasi Dikembalikan
      </button>
    </template>
  </UiAppModal>

  <!-- Modal Detail Peminjaman -->
  <UiAppModal v-if="detailData" title="Detail Peminjaman Arsip" size="lg" @close="detailData = null">
    <div class="row g-3">
      <div class="col-md-6">
        <label class="text-muted small">No. Tiket</label>
        <p class="fw-bold text-primary mb-0">{{ detailData.noTiket }}</p>
      </div>
      <div class="col-md-6">
        <label class="text-muted small">Status</label>
        <p class="mb-0">
          <span :class="['badge', statusBadgeClass(detailData.status)]">{{ detailData.status }}</span>
        </p>
      </div>
      <div class="col-md-6">
        <label class="text-muted small">Nama Peminjam</label>
        <p class="fw-semibold mb-0">{{ detailData.peminjam }}</p>
      </div>
      <div class="col-md-6">
        <label class="text-muted small">Unit Kerja</label>
        <p class="mb-0">{{ detailData.unitKerja || '-' }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">No. Item Arsip</label>
        <p class="fw-semibold mb-0">{{ detailData.noItemArsip }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Kode Klasifikasi</label>
        <p class="mb-0"><span class="badge bg-light-primary text-primary">{{ detailData.kodeKlasifikasi || '-' }}</span></p>
      </div>
      <div class="col-md-12">
        <label class="text-muted small">Uraian Informasi Arsip</label>
        <p class="mb-0">{{ detailData.uraianArsip }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Tanggal Pinjam</label>
        <p class="mb-0 fw-medium text-dark">{{ formatDateIndo(detailData.tanggalPinjam) }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Rencana Kembali</label>
        <p class="mb-0 fw-medium text-dark">{{ formatDateIndo(detailData.tanggalRencanaKembali) }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Tanggal Kembali Aktual</label>
        <p class="mb-0 fw-medium text-dark">{{ detailData.tanggalKembali ? formatDateIndo(detailData.tanggalKembali) : '-' }}</p>
      </div>
      <div class="col-12">
        <label class="text-muted small">Catatan / Keperluan</label>
        <p class="mb-0">{{ detailData.keterangan || '-' }}</p>
      </div>
    </div>
    <template #footer>
      <button class="btn btn-outline-secondary" @click="detailData = null">Tutup</button>
      <button
        v-if="canWrite() && isMasihDipinjam(detailData.status)"
        class="btn btn-success"
        @click="openKembaliFromDetail(detailData)"
      >
        <i class="ti ti-check me-1"></i>Proses Pengembalian
      </button>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { Peminjaman, Arsip, UnitKerja } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Peminjaman Arsip | Arsip Record Center KemenPPPA' })

const { canWrite } = useSession()
const { getPeminjamanList, createPeminjaman, prosesPengembalian, getArsipList, getUnitKerja } = useArsip()
const { showToast } = useToast()

const allPeminjaman = ref<Peminjaman[]>([])
const arsipList = ref<Arsip[]>([])
const unitList = ref<UnitKerja[]>([])
const initialLoading = ref(true)
const submitting = ref(false)

const filter = reactive({
  unitKerja: '',
  status: '',
})

const hasActiveFilter = computed(() => Boolean(filter.unitKerja || filter.status))

const statusFilterOptions = [
  { value: 'Dipinjam', label: 'Dipinjam' },
  { value: 'Terlambat', label: 'Terlambat (Jatuh Tempo)' },
  { value: 'Dikembalikan', label: 'Dikembalikan' },
]

const unitFilterOptions = computed(() =>
  unitList.value.map((u) => ({
    value: u.namaUnit,
    label: `${u.kodeUnit} — ${u.namaUnit}`,
  }))
)

const filteredList = computed(() =>
  allPeminjaman.value.filter((p) => {
    if (filter.unitKerja && p.unitKerja !== filter.unitKerja) return false
    if (filter.status && p.status.toLowerCase() !== filter.status.toLowerCase()) return false
    return true
  })
)

const statCards = computed(() => [
  {
    label: 'Total Transaksi Pinjam',
    value: allPeminjaman.value.length,
    bg: 'bg-primary text-white',
    textClass: 'text-white',
    icon: 'ti ti-book',
  },
  {
    label: 'Sedang Dipinjam',
    value: allPeminjaman.value.filter((p) => p.status.toLowerCase() === 'dipinjam').length,
    bg: 'bg-warning text-white',
    textClass: 'text-white',
    icon: 'ti ti-clock',
  },
  {
    label: 'Sudah Dikembalikan',
    value: allPeminjaman.value.filter((p) => p.status.toLowerCase() === 'dikembalikan').length,
    bg: 'bg-success text-white',
    textClass: 'text-white',
    icon: 'ti ti-circle-check',
  },
  {
    label: 'Terlambat Kembali',
    value: allPeminjaman.value.filter((p) => p.status.toLowerCase() === 'terlambat').length,
    bg: 'bg-danger text-white',
    textClass: 'text-white',
    icon: 'ti ti-alert-triangle',
  },
])

async function loadData(isBackground = false) {
  if (!isBackground) initialLoading.value = true
  try {
    const [loans, arsipRes, units] = await Promise.all([
      getPeminjamanList({ per_page: 500 }),
      getArsipList({ per_page: 500 }),
      getUnitKerja(),
    ])
    allPeminjaman.value = loans
    arsipList.value = arsipRes.data
    unitList.value = units
  } catch (err: any) {
    showToast(err?.message || 'Gagal memuat data peminjaman', 'danger')
  } finally {
    initialLoading.value = false
  }
}

onMounted(() => loadData(false))

function resetFilter() {
  filter.unitKerja = ''
  filter.status = ''
}

const showTambah = ref(false)
const detailData = ref<Peminjaman | null>(null)
const kembaliTarget = ref<Peminjaman | null>(null)
const catatanKembali = ref('')

function openTambah() {
  showTambah.value = true
}

function detailPeminjaman(item: Peminjaman) {
  detailData.value = item
}

function openKembaliModal(item: Peminjaman) {
  kembaliTarget.value = item
  catatanKembali.value = item.keterangan || 'Berkas dikembalikan lengkap'
}

function openKembaliFromDetail(item: Peminjaman) {
  detailData.value = null
  openKembaliModal(item)
}

async function simpanPeminjaman(payload: Record<string, any>) {
  submitting.value = true
  try {
    await createPeminjaman(payload)
    showToast('Peminjaman arsip berhasil dicatat', 'success')
    showTambah.value = false
    await loadData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal mencatat peminjaman', 'danger')
  } finally {
    submitting.value = false
  }
}

async function submitPengembalian() {
  if (!kembaliTarget.value) return
  submitting.value = true
  try {
    await prosesPengembalian(kembaliTarget.value.id, catatanKembali.value)
    showToast(`Arsip ${kembaliTarget.value.noItemArsip} berhasil dikembalikan`, 'success')
    kembaliTarget.value = null
    await loadData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal memproses pengembalian arsip', 'danger')
  } finally {
    submitting.value = false
  }
}

function isMasihDipinjam(status: string) {
  const s = status.toLowerCase()
  return s === 'dipinjam' || s === 'terlambat'
}

function statusBadgeClass(status: string) {
  const s = status.toLowerCase()
  if (s === 'dikembalikan') return 'bg-success'
  if (s === 'terlambat') return 'bg-danger'
  return 'bg-warning text-dark'
}
</script>
