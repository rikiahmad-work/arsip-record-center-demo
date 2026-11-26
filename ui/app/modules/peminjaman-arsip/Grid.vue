<template>
  <div class="dt-responsive table-responsive user-table-wrapper">
    <table
      id="tabel-peminjaman"
      class="table table-bordered table-hover align-middle mb-0 nowrap"
    >
      <thead class="table-light">
        <tr>
          <th class="text-center" style="width: 45px;">NO</th>
          <th class="text-nowrap">NO. TIKET</th>
          <th class="text-nowrap">PEMINJAM &amp; UNIT KERJA</th>
          <th class="text-nowrap">NO. ITEM / KODE</th>
          <th style="min-width: 240px;">URAIAN INFORMASI ARSIP</th>
          <th class="text-center text-nowrap">TGL. PINJAM</th>
          <th class="text-center text-nowrap">RENCANA KEMBALI</th>
          <th class="text-center text-nowrap">TGL. KEMBALI</th>
          <th class="text-center text-nowrap">STATUS</th>
          <th class="text-center text-nowrap sticky-action-th">AKSI</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in list" :key="item.id">
          <td class="text-center">{{ idx + 1 }}</td>
          <td class="text-nowrap fw-semibold text-primary">{{ item.noTiket }}</td>
          <td class="text-nowrap">
            <div class="fw-semibold text-dark">{{ item.peminjam }}</div>
            <small class="text-muted">{{ item.unitKerja || '-' }}</small>
          </td>
          <td class="text-nowrap">
            <div class="fw-semibold">{{ item.noItemArsip }}</div>
            <span v-if="item.kodeKlasifikasi" class="badge bg-light-primary text-primary mt-1">
              {{ item.kodeKlasifikasi }}
            </span>
          </td>
          <td style="min-width: 240px; max-width: 320px; white-space: normal;">
            <div class="text-dark">{{ item.uraianArsip }}</div>
            <small v-if="item.keterangan" class="text-muted d-block mt-1">Catatan: {{ item.keterangan }}</small>
          </td>
          <td class="text-center text-nowrap">{{ formatDateIndo(item.tanggalPinjam) }}</td>
          <td class="text-center text-nowrap">{{ formatDateIndo(item.tanggalRencanaKembali) }}</td>
          <td class="text-center text-nowrap">{{ item.tanggalKembali ? formatDateIndo(item.tanggalKembali) : '-' }}</td>
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
                title="Lihat Detail Peminjaman"
                @click="emit('detail', item)"
              >
                <i class="ti ti-eye"></i>
              </button>
              <button
                v-if="canWrite && isMasihDipinjam(item.status)"
                type="button"
                class="btn btn-sm btn-light-success"
                title="Proses Pengembalian Arsip"
                @click="emit('kembali', item)"
              >
                <i class="ti ti-check"></i>
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { Peminjaman } from '~/types/arsip'

interface Props {
  list: Peminjaman[]
  canWrite?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  canWrite: true,
})

const emit = defineEmits<{
  detail: [item: Peminjaman]
  kembali: [item: Peminjaman]
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
  if (win.$ && win.$.fn?.DataTable && win.$.fn.DataTable.isDataTable('#tabel-peminjaman')) {
    try {
      win.$('#tabel-peminjaman').DataTable().destroy()
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

  win.$('#tabel-peminjaman').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [9] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ peminjaman',
      infoEmpty: 'Menampilkan 0 peminjaman',
      infoFiltered: '(disaring dari _MAX_ total peminjaman)',
      zeroRecords: 'Tidak ada data peminjaman yang sesuai filter',
      emptyTable: 'Belum ada data peminjaman arsip',
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

function isMasihDipinjam(status: string) {
  const s = status.toLowerCase()
  return s === 'dipinjam' || s === 'terlambat'
}

function statusBadge(status: string) {
  const s = status.toLowerCase()
  if (s === 'dikembalikan') return 'bg-success'
  if (s === 'terlambat') return 'bg-danger'
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
