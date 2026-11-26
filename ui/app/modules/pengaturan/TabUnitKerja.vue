<template>
  <div>
    <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
      <div>
        <h6 class="mb-0">Daftar Unit Kerja</h6>
        <small class="text-muted">Total {{ list.length }} unit kerja / unit pengolah.</small>
      </div>
      <div>
        <button class="btn btn-sm btn-primary text-nowrap" type="button" @click="openTambah">
          <i class="ti ti-plus me-1"></i>Tambah Unit
        </button>
      </div>
    </div>

    <div class="dt-responsive table-responsive user-table-wrapper">
      <table
        id="tabel-unit-kerja"
        class="table table-bordered table-hover align-middle mb-0 nowrap"
      >
        <thead class="table-light">
          <tr>
            <th class="text-nowrap">KODE UNIT</th>
            <th class="text-nowrap">NAMA UNIT KERJA</th>
            <th class="text-center text-nowrap">ESELON</th>
            <th class="text-center text-nowrap">TOTAL ARSIP</th>
            <th class="text-center text-nowrap">STATUS</th>
            <th class="text-center text-nowrap sticky-action-th">AKSI</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td class="fw-semibold text-dark">{{ item.kodeUnit }}</td>
            <td>{{ item.namaUnit }}</td>
            <td class="text-center text-nowrap">{{ item.eselon || '-' }}</td>
            <td class="text-center text-nowrap">{{ item.totalArsip }}</td>
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
      :title="editingId ? 'Edit Unit Kerja' : 'Tambah Unit Kerja'"
      size="sm"
      @close="showForm = false"
    >
      <div class="row g-3">
        <div class="col-12">
          <label class="form-label">Kode Unit <span class="text-danger">*</span></label>
          <input v-model="form.kodeUnit" class="form-control" placeholder="Contoh: UK-001" required>
        </div>
        <div class="col-12">
          <label class="form-label">Nama Unit <span class="text-danger">*</span></label>
          <input v-model="form.namaUnit" class="form-control" placeholder="Contoh: Biro Umum" required>
        </div>
        <div class="col-12">
          <label class="form-label">Eselon</label>
          <input v-model="form.eselon" class="form-control" placeholder="Contoh: II">
        </div>
        <div class="col-12">
          <label class="form-label">Status</label>
          <select v-model="form.statusKey" class="form-select">
            <option value="aktif">Aktif</option>
            <option value="nonaktif">Nonaktif</option>
          </select>
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
import type { UnitKerja } from '~/types/index'

interface Props {
  list: UnitKerja[]
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  saving: false,
})

const emit = defineEmits<{
  save: [payload: Record<string, any>]
  hapus: [item: UnitKerja]
}>()

const showForm = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  kodeUnit: '',
  namaUnit: '',
  eselon: '',
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
  if (win.$ && win.$.fn && win.$.fn.DataTable && win.$.fn.DataTable.isDataTable('#tabel-unit-kerja')) {
    try {
      win.$('#tabel-unit-kerja').DataTable().destroy()
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

  win.$('#tabel-unit-kerja').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [5] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ unit',
      infoEmpty: 'Menampilkan 0 unit',
      infoFiltered: '(disaring dari _MAX_ total unit)',
      zeroRecords: 'Tidak ada unit kerja yang cocok',
      emptyTable: 'Tidak ada data unit kerja',
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
  form.kodeUnit = ''
  form.namaUnit = ''
  form.eselon = ''
  form.statusKey = 'aktif'
}

function openTambah() {
  resetForm()
  showForm.value = true
}

function openEdit(item: UnitKerja) {
  editingId.value = item.id
  form.kodeUnit = item.kodeUnit
  form.namaUnit = item.namaUnit
  form.eselon = item.eselon || ''
  form.statusKey = item.statusKey
  showForm.value = true
}

function simpan() {
  if (!form.kodeUnit.trim() || !form.namaUnit.trim()) return
  emit('save', {
    id: editingId.value,
    kode: form.kodeUnit.trim(),
    nama: form.namaUnit.trim(),
    eselon: form.eselon || null,
    is_active: form.statusKey === 'aktif',
  })
  showForm.value = false
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
