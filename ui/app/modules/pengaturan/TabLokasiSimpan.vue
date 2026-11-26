<template>
  <div>
    <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
      <div>
        <h6 class="mb-0">Daftar Lokasi Simpan</h6>
        <small class="text-muted">Total {{ list.length }} lokasi penyimpanan arsip.</small>
      </div>
      <div>
        <button class="btn btn-sm btn-primary text-nowrap" type="button" @click="openTambah">
          <i class="ti ti-plus me-1"></i>Tambah Lokasi
        </button>
      </div>
    </div>

    <div class="dt-responsive table-responsive user-table-wrapper">
      <table
        id="tabel-lokasi-simpan"
        class="table table-bordered table-hover align-middle mb-0 nowrap"
      >
        <thead class="table-light">
          <tr>
            <th class="text-nowrap">NAMA LOKASI</th>
            <th class="text-center text-nowrap">NO. RAK</th>
            <th class="text-center text-nowrap">NO. BOKS</th>
            <th class="text-center text-nowrap">KAPASITAS</th>
            <th class="text-center text-nowrap">TERPAKAI</th>
            <th class="text-nowrap" style="min-width: 140px">UTILISASI</th>
            <th class="text-center text-nowrap">STATUS</th>
            <th class="text-center text-nowrap sticky-action-th">AKSI</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td class="fw-semibold text-dark">{{ item.namaLokasi }}</td>
            <td class="text-center text-nowrap">{{ item.noRak || '-' }}</td>
            <td class="text-center text-nowrap">{{ item.noBoks || '-' }}</td>
            <td class="text-center text-nowrap">{{ item.kapasitas }}</td>
            <td class="text-center text-nowrap">{{ item.terpakai }}</td>
            <td style="min-width: 140px">
              <div class="progress" style="height: 8px">
                <div
                  class="progress-bar"
                  :class="progressClass(item)"
                  :style="{ width: `${progressPct(item)}%` }"
                />
              </div>
              <small class="text-muted">{{ progressPct(item) }}%</small>
            </td>
            <td class="text-center text-nowrap">
              <UiAppBadge :status="item.statusKey" :label="item.statusLabel" />
            </td>
            <td class="text-center text-nowrap sticky-action-td">
              <div class="d-inline-flex gap-1 justify-content-center flex-nowrap">
                <button class="btn btn-sm btn-light-primary" title="Edit" type="button" @click="openEdit(item)">
                  <i class="ti ti-pencil"></i>
                </button>
                <button class="btn btn-sm btn-light-danger" title="Hapus" type="button" @click="emit('hapus', item)">
                  <i class="ti ti-trash"></i>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Form -->
    <UiAppModal
      v-if="showForm"
      :title="editingId ? 'Edit Lokasi Simpan' : 'Tambah Lokasi Simpan'"
      size="lg"
      @close="showForm = false"
    >
      <div class="row g-3">
        <div class="col-md-8">
          <label class="form-label">Nama Lokasi <span class="text-danger">*</span></label>
          <input v-model="form.namaSimpan" class="form-control" placeholder="Contoh: Rak A - Lantai 2" required>
        </div>
        <div class="col-md-4">
          <label class="form-label">Status</label>
          <select v-model="form.statusKey" class="form-select">
            <option value="aktif">Aktif</option>
            <option value="nonaktif">Nonaktif</option>
          </select>
        </div>
        <div class="col-md-4">
          <label class="form-label">Nomor Rak</label>
          <input v-model="form.noRak" class="form-control" placeholder="Contoh: R-01">
        </div>
        <div class="col-md-4">
          <label class="form-label">Nomor Boks</label>
          <input v-model="form.noBoks" class="form-control" placeholder="Contoh: B-01">
        </div>
        <div class="col-md-4">
          <label class="form-label">Kapasitas Maksimal</label>
          <input v-model.number="form.kapasitas" class="form-control" type="number" min="1" placeholder="50">
        </div>
      </div>

      <template #footer>
        <button class="btn btn-secondary" type="button" @click="showForm = false">Batal</button>
        <button class="btn btn-primary" type="button" :disabled="saving" @click="simpan">
          <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
          Simpan
        </button>
      </template>
    </UiAppModal>
  </div>
</template>

<script setup lang="ts">
import type { LokasiSimpan } from '~/types/index'

interface Props {
  list: LokasiSimpan[]
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  saving: false,
})

const emit = defineEmits<{
  save: [payload: Record<string, any>]
  hapus: [item: LokasiSimpan]
}>()

const showForm = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  namaSimpan: '',
  noRak: '',
  noBoks: '',
  kapasitas: 50,
  statusKey: 'aktif' as 'aktif' | 'nonaktif',
})

async function ensureDataTablesLoaded() {
  const scripts = [
    '/admin/js/plugins/jquery.min.js',
    '/admin/js/plugins/jquery.dataTables.min.js',
    '/admin/js/plugins/dataTables.bootstrap5.min.js',
  ]

  const loadLink = (href: string) => {
    if (document.querySelector(`link[href="${href}"]`)) return
    const l = document.createElement('link')
    l.rel = 'stylesheet'
    l.href = href
    document.head.appendChild(l)
  }

  const loadScript = (src: string) =>
    new Promise<void>((resolve) => {
      if (document.querySelector(`script[src="${src}"]`)) {
        resolve()
        return
      }
      const s = document.createElement('script')
      s.src = src
      s.onload = () => resolve()
      s.onerror = () => resolve()
      document.body.appendChild(s)
    })

  loadLink('/admin/css/plugins/dataTables.bootstrap5.min.css')
  await scripts.reduce((p, s) => p.then(() => loadScript(s)), Promise.resolve())
}

function destroyDataTable() {
  const win = window as any
  if (win.$ && win.$.fn && win.$.fn.DataTable && win.$.fn.DataTable.isDataTable('#tabel-lokasi-simpan')) {
    try {
      win.$('#tabel-lokasi-simpan').DataTable().destroy()
    } catch (_) {}
  }
}

async function initOrRefreshDataTable() {
  await ensureDataTablesLoaded()
  await nextTick()

  const win = window as any
  if (!win.$ || !win.$.fn || !win.$.fn.DataTable) return

  destroyDataTable()
  await nextTick()

  win.$('#tabel-lokasi-simpan').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [5, 7] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ lokasi',
      infoEmpty: 'Menampilkan 0 lokasi',
      infoFiltered: '(disaring dari _MAX_ total lokasi)',
      zeroRecords: 'Tidak ada lokasi simpan yang cocok',
      emptyTable: 'Tidak ada data lokasi simpan',
      paginate: {
        first: 'Pertama',
        last: 'Terakhir',
        next: 'Selanjutnya',
        previous: 'Sebelumnya',
      },
    },
  })
}

onMounted(() => {
  initOrRefreshDataTable()
})

watch(
  () => props.list,
  async () => {
    destroyDataTable()
    await nextTick()
    await initOrRefreshDataTable()
  },
  { flush: 'pre' }
)

onBeforeUnmount(() => {
  destroyDataTable()
})

function resetForm() {
  editingId.value = null
  form.namaSimpan = ''
  form.noRak = ''
  form.noBoks = ''
  form.kapasitas = 50
  form.statusKey = 'aktif'
}

function openTambah() {
  resetForm()
  showForm.value = true
}

function openEdit(item: LokasiSimpan) {
  editingId.value = item.id
  form.namaSimpan = item.namaLokasi
  form.noRak = item.noRak || ''
  form.noBoks = item.noBoks || ''
  form.kapasitas = item.kapasitas
  form.statusKey = item.statusKey
  showForm.value = true
}

function simpan() {
  if (!form.namaSimpan.trim()) return
  emit('save', {
    id: editingId.value,
    nama_simpan: form.namaSimpan.trim(),
    no_rak: form.noRak.trim() || null,
    no_boks: form.noBoks.trim() || null,
    kapasitas: form.kapasitas,
    is_active: form.statusKey === 'aktif',
  })
  showForm.value = false
}

function progressPct(item: LokasiSimpan) {
  if (!item.kapasitas) return 0
  return Math.min(100, Math.round((item.terpakai / item.kapasitas) * 100))
}

function progressClass(item: LokasiSimpan) {
  const pct = progressPct(item)
  if (pct >= 90) return 'bg-danger'
  if (pct >= 70) return 'bg-warning'
  return 'bg-success'
}
</script>

<style>
/* Matikan scroll di wrapper terluar agar kontrol DataTables (Length, Search, Info, Pagination) tetap diam */
.user-table-wrapper {
  overflow: visible !important;
}

/* Hanya baris pembungkus tabel (.dt-row) yang melakukan scroll horizontal */
.user-table-wrapper .dataTables_wrapper .dt-row {
  margin-left: 0;
  margin-right: 0;
}

.user-table-wrapper .dataTables_wrapper .dt-row > div {
  padding-left: 0;
  padding-right: 0;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

/* Kolom AKSI Fixed / Sticky di sebelah kanan saat tabel di-scroll horizontal */
.user-table-wrapper table.table thead > tr > th.sticky-action-th {
  position: sticky !important;
  right: 0 !important;
  z-index: 3 !important;
  background-color: #f8f9fa !important;
  box-shadow: -3px 0 6px rgba(0, 0, 0, 0.06);
}

.user-table-wrapper table.table tbody > tr > td.sticky-action-td {
  position: sticky !important;
  right: 0 !important;
  z-index: 2;
  background-color: #ffffff !important;
  box-shadow: -3px 0 6px rgba(0, 0, 0, 0.06);
}

.user-table-wrapper .table-hover > tbody > tr:hover > td.sticky-action-td {
  background-color: #eef2f6 !important;
}
</style>
