<template>
  <div>
    <div class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3">
      <div>
        <h6 class="mb-0">Daftar Kode Klasifikasi</h6>
        <small class="text-muted">Total {{ list.length }} kode — JRA &amp; SKKAD otomatis dipakai saat input arsip.</small>
      </div>
      <div>
        <button class="btn btn-sm btn-primary text-nowrap" type="button" @click="openTambah">
          <i class="ti ti-plus me-1"></i>Tambah Kode
        </button>
      </div>
    </div>

    <div class="dt-responsive table-responsive user-table-wrapper">
      <table
        id="tabel-kode-klasifikasi"
        class="table table-bordered table-hover align-middle mb-0 nowrap"
      >
        <thead class="table-light">
          <tr>
            <th class="text-nowrap">KODE</th>
            <th class="text-nowrap">NAMA KLASIFIKASI</th>
            <th class="text-center text-nowrap">JANGKA SIMPAN</th>
            <th class="text-center text-nowrap">NASIB AKHIR</th>
            <th class="text-center text-nowrap">SKKAD</th>
            <th class="text-center text-nowrap">TOTAL ARSIP</th>
            <th class="text-center text-nowrap">STATUS</th>
            <th class="text-center text-nowrap sticky-action-th">AKSI</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td class="fw-semibold text-dark">{{ item.kode }}</td>
            <td>{{ item.nama }}</td>
            <td class="text-center text-nowrap">{{ item.jangkaSimpanInaktif || '-' }}</td>
            <td class="text-center text-nowrap">{{ item.nasibAkhir || '-' }}</td>
            <td class="text-center text-nowrap">
              <span :class="['badge', skkadBadge(item.skkad)]">{{ item.skkad || '-' }}</span>
            </td>
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
      :title="editingId ? 'Edit Kode Klasifikasi' : 'Tambah Kode Klasifikasi'"
      size="lg"
      @close="showForm = false"
    >
      <div class="row g-3">
        <div class="col-md-4">
          <label class="form-label">Kode <span class="text-danger">*</span></label>
          <input v-model="form.kode" class="form-control" placeholder="Contoh: DT.001" required>
        </div>
        <div class="col-md-8">
          <label class="form-label">Nama Klasifikasi <span class="text-danger">*</span></label>
          <input v-model="form.nama" class="form-control" placeholder="Nama klasifikasi arsip" required>
        </div>
        <div class="col-md-4">
          <label class="form-label">Jangka Simpan Inaktif</label>
          <input v-model="form.jangkaSimpanInaktif" class="form-control" placeholder="Contoh: 5 Tahun">
        </div>
        <div class="col-md-4">
          <label class="form-label">Nasib Akhir</label>
          <select v-model="form.nasibAkhir" class="form-select">
            <option value="">— Pilih —</option>
            <option value="Musnah">Musnah</option>
            <option value="Permanen">Permanen</option>
            <option value="Dinilai Kembali">Dinilai Kembali</option>
          </select>
        </div>
        <div class="col-md-4">
          <label class="form-label">Klasifikasi Keamanan (SKKAD)</label>
          <select v-model="form.skkad" class="form-select">
            <option value="Biasa">Biasa / Terbuka</option>
            <option value="Terbatas">Terbatas</option>
            <option value="Rahasia">Rahasia</option>
            <option value="Sangat Rahasia">Sangat Rahasia</option>
          </select>
        </div>
        <div class="col-md-4">
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
import type { KodeKlasifikasi } from '~/types/index'

interface Props {
  list: KodeKlasifikasi[]
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  saving: false,
})

const emit = defineEmits<{
  save: [payload: Record<string, any>]
  hapus: [item: KodeKlasifikasi]
}>()

const showForm = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  kode: '',
  nama: '',
  jangkaSimpanInaktif: '',
  nasibAkhir: '',
  skkad: 'Biasa',
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
  if (win.$ && win.$.fn && win.$.fn.DataTable && win.$.fn.DataTable.isDataTable('#tabel-kode-klasifikasi')) {
    try {
      win.$('#tabel-kode-klasifikasi').DataTable().destroy()
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

  win.$('#tabel-kode-klasifikasi').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [7] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ kode',
      infoEmpty: 'Menampilkan 0 kode',
      infoFiltered: '(disaring dari _MAX_ total kode)',
      zeroRecords: 'Tidak ada kode klasifikasi yang cocok',
      emptyTable: 'Tidak ada data kode klasifikasi',
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
  form.kode = ''
  form.nama = ''
  form.jangkaSimpanInaktif = ''
  form.nasibAkhir = ''
  form.skkad = 'Biasa'
  form.statusKey = 'aktif'
}

function openTambah() {
  resetForm()
  showForm.value = true
}

function openEdit(item: KodeKlasifikasi) {
  editingId.value = item.id
  form.kode = item.kode
  form.nama = item.nama
  form.jangkaSimpanInaktif = item.jangkaSimpanInaktif || ''
  form.nasibAkhir = item.nasibAkhir || ''
  form.skkad = item.skkad || 'Biasa'
  form.statusKey = item.statusKey
  showForm.value = true
}

function simpan() {
  if (!form.kode.trim() || !form.nama.trim()) return
  emit('save', {
    id: editingId.value,
    kode: form.kode.trim(),
    nama: form.nama.trim(),
    jangka_simpan_inaktif: form.jangkaSimpanInaktif || null,
    nasib_akhir: form.nasibAkhir || null,
    skkad: form.skkad,
    is_active: form.statusKey === 'aktif',
  })
  showForm.value = false
}

function skkadBadge(skkad?: string) {
  switch (skkad) {
    case 'Rahasia':
      return 'bg-danger-subtle text-danger'
    case 'Terbatas':
      return 'bg-warning-subtle text-warning'
    default:
      return 'bg-success-subtle text-success'
  }
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
