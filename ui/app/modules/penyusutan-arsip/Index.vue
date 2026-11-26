<template>
  <TemplatePageTitle title="Penyusutan Arsip" :breadcrumbs="[{ label: 'Penyusutan Arsip' }]">
    <template #action>
      <button
        v-if="canWrite()"
        type="button"
        class="btn btn-primary d-inline-flex align-items-center gap-1"
        @click="showTambah = true"
      >
        <i class="ti ti-plus"></i>
        <span>Ajukan Usul Penyusutan</span>
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
                  <i class="ti ti-category me-1"></i>Jenis Usulan Penyusutan
                </label>
                <FormSelect
                  v-model="filter.jenisUsul"
                  :options="jenisUsulOptions"
                  placeholder="Semua Jenis Usul (Musnah & Serah)"
                  search-placeholder="Cari jenis usul..."
                  size="sm"
                  searchable
                  clearable
                  no-margin
                />
              </div>

              <div class="col-12 col-md-5">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-status-change me-1"></i>Status Persetujuan
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
            <p class="mt-2 text-muted">Memuat data usulan penyusutan...</p>
          </div>

          <ModulesPenyusutanArsipGrid
            v-else
            :list="filteredList"
            :can-approve="canApprove()"
            :can-write="canWrite()"
            @detail="detailPenyusutan"
            @setujui="openDecisionModal($event, 'approve')"
            @tolak="openDecisionModal($event, 'reject')"
            @selesai="handleSelesai"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Pengajuan Penyusutan -->
  <ModulesPenyusutanArsipForm
    v-if="showTambah"
    :arsip-list="arsipList"
    :submitting="submitting"
    @close="showTambah = false"
    @save="simpanPenyusutan"
  />

  <!-- Modal Keputusan Approver (Setujui / Tolak) -->
  <UiAppModal
    v-if="decisionTarget"
    :title="decisionType === 'approve' ? 'Setujui Usulan Penyusutan Arsip' : 'Tolak Usulan Penyusutan Arsip'"
    size="md"
    @close="decisionTarget = null"
  >
    <div class="mb-3">
      <div class="p-3 bg-light rounded border mb-3">
        <div class="d-flex justify-content-between align-items-center mb-1">
          <span class="fw-bold text-primary">{{ decisionTarget.noPengajuan }}</span>
          <span :class="['badge', decisionTarget.jenisUsul === 'Serah' ? 'bg-info' : 'bg-dark']">
            Usul {{ decisionTarget.jenisUsul }}
          </span>
        </div>
        <div class="fw-semibold text-dark">{{ decisionTarget.noItem }} — {{ decisionTarget.uraian }}</div>
        <small class="text-muted">Diajukan oleh: {{ decisionTarget.diajukanOleh }} pada {{ formatDateIndo(decisionTarget.tanggalPengajuan) }}</small>
      </div>

      <label class="form-label fw-semibold">
        Catatan Pejabat Penilai / Approver <span v-if="decisionType === 'reject'" class="text-danger">*</span>
      </label>
      <textarea
        v-model="decisionNote"
        class="form-control"
        rows="3"
        :placeholder="decisionType === 'approve' ? 'Contoh: Telah diverifikasi sesuai JRA, disetujui untuk pelaksanaan penyusutan' : 'Masukkan alasan penolakan (berkas akan dikembalikan ke status Inaktif)'"
      ></textarea>
    </div>
    <template #footer>
      <button type="button" class="btn btn-outline-secondary" @click="decisionTarget = null">Batal</button>
      <button
        type="button"
        :class="['btn', decisionType === 'approve' ? 'btn-success' : 'btn-danger']"
        :disabled="submitting"
        @click="submitDecision"
      >
        <span v-if="submitting" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else :class="['ti me-1', decisionType === 'approve' ? 'ti-check' : 'ti-x']"></i>
        {{ decisionType === 'approve' ? 'Setujui Usulan' : 'Tolak Usulan' }}
      </button>
    </template>
  </UiAppModal>

  <!-- Modal Detail Penyusutan -->
  <UiAppModal v-if="detailData" title="Detail Usulan Penyusutan Arsip" size="lg" @close="detailData = null">
    <div class="row g-3">
      <div class="col-md-4">
        <label class="text-muted small">No. Pengajuan</label>
        <p class="fw-bold text-primary mb-0">{{ detailData.noPengajuan }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Jenis Usulan</label>
        <p class="mb-0">
          <span :class="['badge', detailData.jenisUsul === 'Serah' ? 'bg-info' : 'bg-dark']">
            Usul {{ detailData.jenisUsul }}
          </span>
        </p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Status</label>
        <p class="mb-0">
          <span :class="['badge', statusBadge(detailData.status)]">{{ detailData.status }}</span>
        </p>
      </div>
      <div class="col-md-6">
        <label class="text-muted small">No. Item Arsip</label>
        <p class="fw-semibold mb-0">{{ detailData.noItem }}</p>
      </div>
      <div class="col-md-6">
        <label class="text-muted small">Kode Klasifikasi &amp; Nasib Akhir JRA</label>
        <p class="mb-0">
          <span class="badge bg-light-primary text-primary me-1">{{ detailData.kodeKlasifikasi || '-' }}</span>
          <span class="badge bg-secondary">{{ detailData.nasibAkhir }}</span>
        </p>
      </div>
      <div class="col-12">
        <label class="text-muted small">Uraian Informasi Arsip</label>
        <p class="mb-0">{{ detailData.uraian }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Tanggal Jatuh Tempo Retensi</label>
        <p class="mb-0 fw-medium text-dark">{{ formatDateIndo(detailData.tanggalRetensi) }}</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Diajukan Oleh</label>
        <p class="mb-0 fw-medium text-dark">{{ detailData.diajukanOleh }} ({{ formatDateIndo(detailData.tanggalPengajuan) }})</p>
      </div>
      <div class="col-md-4">
        <label class="text-muted small">Disetujui / Diverifikasi Oleh</label>
        <p class="mb-0">{{ detailData.disetujuiOleh || '-' }}</p>
      </div>
      <div class="col-12">
        <label class="text-muted small">Alasan / Dasar Pengajuan</label>
        <p class="mb-0">{{ detailData.keterangan || '-' }}</p>
      </div>
      <div v-if="detailData.catatanApprover" class="col-12">
        <label class="text-muted small">Catatan Approver</label>
        <div class="p-2 bg-light rounded border small">{{ detailData.catatanApprover }}</div>
      </div>
    </div>
    <template #footer>
      <button class="btn btn-outline-secondary" @click="detailData = null">Tutup</button>
      <template v-if="canApprove() && detailData.status.toLowerCase() === 'diajukan'">
        <button class="btn btn-danger" @click="openDecisionFromDetail(detailData, 'reject')">
          <i class="ti ti-x me-1"></i>Tolak
        </button>
        <button class="btn btn-success" @click="openDecisionFromDetail(detailData, 'approve')">
          <i class="ti ti-check me-1"></i>Setujui
        </button>
      </template>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { Penyusutan, Arsip } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Penyusutan Arsip | Arsip Record Center KemenPPPA' })

const { canWrite, canApprove } = useSession()
const {
  getPenyusutanList,
  createPenyusutan,
  setujuiPenyusutan,
  tolakPenyusutan,
  selesaikanPenyusutan,
  getArsipList,
} = useArsip()
const { showToast } = useToast()

const allPenyusutan = ref<Penyusutan[]>([])
const arsipList = ref<Arsip[]>([])
const initialLoading = ref(true)
const submitting = ref(false)

const filter = reactive({
  jenisUsul: '',
  status: '',
})

const hasActiveFilter = computed(() => Boolean(filter.jenisUsul || filter.status))

const jenisUsulOptions = [
  { value: 'Musnah', label: 'Usul Musnah' },
  { value: 'Serah', label: 'Usul Serah (Permanen / ANRI)' },
]

const statusFilterOptions = [
  { value: 'Diajukan', label: 'Diajukan (Menunggu Persetujuan)' },
  { value: 'Disetujui', label: 'Disetujui' },
  { value: 'Ditolak', label: 'Ditolak' },
  { value: 'Selesai', label: 'Selesai Dieksekusi' },
]

const filteredList = computed(() =>
  allPenyusutan.value.filter((item) => {
    if (filter.jenisUsul && item.jenisUsul.toLowerCase() !== filter.jenisUsul.toLowerCase()) return false
    if (filter.status && item.status.toLowerCase() !== filter.status.toLowerCase()) return false
    return true
  })
)

const statCards = computed(() => [
  {
    label: 'Total Usulan Penyusutan',
    value: allPenyusutan.value.length,
    bg: 'bg-primary text-white',
    textClass: 'text-white',
    icon: 'ti ti-files',
  },
  {
    label: 'Menunggu Persetujuan',
    value: allPenyusutan.value.filter((p) => p.status.toLowerCase() === 'diajukan').length,
    bg: 'bg-warning text-white',
    textClass: 'text-white',
    icon: 'ti ti-clock',
  },
  {
    label: 'Disetujui / Selesai',
    value: allPenyusutan.value.filter((p) => ['disetujui', 'selesai'].includes(p.status.toLowerCase())).length,
    bg: 'bg-success text-white',
    textClass: 'text-white',
    icon: 'ti ti-circle-check',
  },
  {
    label: 'Usul Musnah vs Serah',
    value: `${allPenyusutan.value.filter((p) => p.jenisUsul === 'Musnah').length} / ${allPenyusutan.value.filter((p) => p.jenisUsul === 'Serah').length}`,
    bg: 'bg-danger text-white',
    textClass: 'text-white',
    icon: 'ti ti-trash',
  },
])

async function loadData(isBackground = false) {
  if (!isBackground) initialLoading.value = true
  try {
    const [disposals, arsipRes] = await Promise.all([
      getPenyusutanList({ per_page: 500 }),
      getArsipList({ per_page: 500 }),
    ])
    allPenyusutan.value = disposals
    arsipList.value = arsipRes.data
  } catch (err: any) {
    showToast(err?.message || 'Gagal memuat data penyusutan', 'danger')
  } finally {
    initialLoading.value = false
  }
}

onMounted(() => loadData(false))

function resetFilter() {
  filter.jenisUsul = ''
  filter.status = ''
}

const showTambah = ref(false)
const detailData = ref<Penyusutan | null>(null)
const decisionTarget = ref<Penyusutan | null>(null)
const decisionType = ref<'approve' | 'reject'>('approve')
const decisionNote = ref('')

function detailPenyusutan(item: Penyusutan) {
  detailData.value = item
}

function openDecisionModal(item: Penyusutan, type: 'approve' | 'reject') {
  decisionTarget.value = item
  decisionType.value = type
  decisionNote.value = type === 'approve'
    ? 'Disetujui sesuai Jadwal Retensi Arsip (JRA)'
    : ''
}

function openDecisionFromDetail(item: Penyusutan, type: 'approve' | 'reject') {
  detailData.value = null
  openDecisionModal(item, type)
}

async function simpanPenyusutan(payload: Record<string, any>) {
  submitting.value = true
  try {
    await createPenyusutan(payload)
    showToast('Usulan penyusutan berhasil diajukan & status arsip diperbarui', 'success')
    showTambah.value = false
    await loadData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal mengajukan penyusutan', 'danger')
  } finally {
    submitting.value = false
  }
}

async function submitDecision() {
  if (!decisionTarget.value) return
  if (decisionType.value === 'reject' && !decisionNote.value.trim()) {
    showToast('Catatan alasan penolakan wajib diisi', 'warning')
    return
  }
  submitting.value = true
  try {
    if (decisionType.value === 'approve') {
      await setujuiPenyusutan(decisionTarget.value.id, decisionNote.value)
      showToast(`Usulan ${decisionTarget.value.noPengajuan} berhasil disetujui`, 'success')
    } else {
      await tolakPenyusutan(decisionTarget.value.id, decisionNote.value)
      showToast(`Usulan ${decisionTarget.value.noPengajuan} ditolak & kategori arsip dikembalikan ke Inaktif`, 'info')
    }
    decisionTarget.value = null
    await loadData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal memproses keputusan penyusutan', 'danger')
  } finally {
    submitting.value = false
  }
}

async function handleSelesai(item: Penyusutan) {
  submitting.value = true
  try {
    await selesaikanPenyusutan(item.id)
    showToast(`Eksekusi penyusutan ${item.noPengajuan} ditandai Selesai`, 'success')
    await loadData(true)
  } catch (err: any) {
    showToast(err?.message || 'Gagal menyelesaikan penyusutan', 'danger')
  } finally {
    submitting.value = false
  }
}

function statusBadge(status: string) {
  const s = status.toLowerCase()
  if (s === 'disetujui') return 'bg-info'
  if (s === 'selesai') return 'bg-success'
  if (s === 'ditolak') return 'bg-danger'
  return 'bg-warning text-dark'
}
</script>
