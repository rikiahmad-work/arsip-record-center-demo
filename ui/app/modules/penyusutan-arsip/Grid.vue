<template>
  <div class="dt-responsive table-responsive user-table-wrapper">
    <table
      id="tabel-penyusutan"
      class="table table-bordered table-hover align-middle mb-0 nowrap"
    >
      <thead class="table-light">
        <tr>
          <th class="text-center" style="width: 45px;">NO</th>
          <th class="text-nowrap">NO. PENGAJUAN</th>
          <th class="text-center text-nowrap">JENIS USUL</th>
          <th class="text-nowrap">NO. ITEM / KODE</th>
          <th style="min-width: 240px;">URAIAN INFORMASI ARSIP</th>
          <th class="text-center text-nowrap">TGL. RETENSI</th>
          <th class="text-center text-nowrap">NASIB AKHIR</th>
          <th class="text-nowrap">DIAJUKAN OLEH</th>
          <th class="text-center text-nowrap">STATUS</th>
          <th class="text-center text-nowrap sticky-action-th">AKSI</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in list" :key="item.id">
          <td class="text-center">{{ idx + 1 }}</td>
          <td class="text-nowrap fw-semibold text-primary">{{ item.noPengajuan }}</td>
          <td class="text-center text-nowrap">
            <span :class="['badge', item.jenisUsul === 'Serah' ? 'bg-info' : 'bg-dark']">
              Usul {{ item.jenisUsul }}
            </span>
          </td>
          <td class="text-nowrap">
            <div class="fw-semibold">{{ item.noItem }}</div>
            <span v-if="item.kodeKlasifikasi" class="badge bg-light-primary text-primary mt-1">
              {{ item.kodeKlasifikasi }}
            </span>
          </td>
          <td style="min-width: 240px; max-width: 320px; white-space: normal;">
            <div class="text-dark">{{ item.uraian }}</div>
            <small v-if="item.keterangan" class="text-muted d-block mt-1">Alasan: {{ item.keterangan }}</small>
          </td>
          <td class="text-center text-nowrap">{{ formatDateIndo(item.tanggalRetensi) }}</td>
          <td class="text-center text-nowrap">
            <span :class="['badge', item.nasibAkhir.toLowerCase() === 'permanen' ? 'bg-primary' : 'bg-secondary']">
              {{ item.nasibAkhir }}
            </span>
          </td>
          <td class="text-nowrap">
            <div class="fw-medium">{{ item.diajukanOleh }}</div>
            <small class="text-muted">{{ formatDateIndo(item.tanggalPengajuan) }}</small>
          </td>
          <td class="text-center text-nowrap">
            <span :class="['badge', statusBadge(item.status)]">
              {{ item.status }}
            </span>
          </td>
          <td class="text-center text-nowrap sticky-action-td">
            <div class="d-inline-flex gap-1 justify-content-center flex-nowrap">
              <button
                type="button"
                class="btn btn-sm btn-light-info"
                title="Lihat Detail Usulan"
                @click="emit('detail', item)"
              >
                <i class="ti ti-eye"></i>
              </button>
              <template v-if="canApprove && item.status.toLowerCase() === 'diajukan'">
                <button
                  type="button"
                  class="btn btn-sm btn-light-success"
                  title="Setujui Usulan Penyusutan"
                  @click="emit('setujui', item)"
                >
                  <i class="ti ti-check"></i>
                </button>
                <button
                  type="button"
                  class="btn btn-sm btn-light-danger"
                  title="Tolak Usulan Penyusutan"
                  @click="emit('tolak', item)"
                >
                  <i class="ti ti-x"></i>
                </button>
              </template>
              <button
                v-if="(canApprove || canWrite) && item.status.toLowerCase() === 'disetujui'"
                type="button"
                class="btn btn-sm btn-light-primary d-inline-flex align-items-center gap-1"
                title="Tandai Eksekusi Penyusutan Selesai"
                @click="emit('selesai', item)"
              >
                <i class="ti ti-circle-check"></i>
                <span>Selesai</span>
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { Penyusutan } from '~/types/arsip'

interface Props {
  list: Penyusutan[]
  canApprove?: boolean
  canWrite?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  canApprove: false,
  canWrite: true,
})

const emit = defineEmits<{
  detail: [item: Penyusutan]
  setujui: [item: Penyusutan]
  tolak: [item: Penyusutan]
  selesai: [item: Penyusutan]
}>()

async function ensureDataTablesLoaded() {
  const scripts = [
    '/admin/js/plugins/jquery.min.js',
    '/admin/js/plugins/jquery.dataTables.min.js',
    '/admin/js/plugins/dataTables.bootstrap5.min.js',
  ]
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
  await scripts.reduce((p, s) => p.then(() => loadScript(s)), Promise.resolve())
}

function destroyDataTable() {
  const win = window as any
  if (win.$ && win.$.fn?.DataTable && win.$.fn.DataTable.isDataTable('#tabel-penyusutan')) {
    try {
      win.$('#tabel-penyusutan').DataTable().destroy()
    } catch (_) {}
  }
}

async function initDataTable() {
  await ensureDataTablesLoaded()
  await nextTick()
  const win = window as any
  if (!win.$ || !win.$.fn?.DataTable) return

  destroyDataTable()
  await nextTick()

  win.$('#tabel-penyusutan').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [9] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ usulan',
      infoEmpty: 'Menampilkan 0 usulan',
      infoFiltered: '(disaring dari _MAX_ total usulan)',
      zeroRecords: 'Tidak ada data penyusutan yang sesuai filter',
      emptyTable: 'Belum ada data usulan penyusutan arsip',
      paginate: {
        first: 'Pertama',
        last: 'Terakhir',
        next: 'Selanjutnya',
        previous: 'Sebelumnya',
      },
    },
  })
}

onMounted(() => initDataTable())

watch(
  () => props.list,
  async () => {
    destroyDataTable()
    await nextTick()
    await initDataTable()
  },
  { flush: 'pre' }
)

onBeforeUnmount(() => destroyDataTable())

function statusBadge(status: string) {
  const s = status.toLowerCase()
  if (s === 'disetujui') return 'bg-info'
  if (s === 'selesai') return 'bg-success'
  if (s === 'ditolak') return 'bg-danger'
  return 'bg-warning text-dark'
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
