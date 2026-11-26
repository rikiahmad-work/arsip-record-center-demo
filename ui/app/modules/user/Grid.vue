<template>
  <div class="dt-responsive table-responsive user-table-wrapper">
    <table
      id="tabel-akun"
      class="table table-bordered table-hover align-middle mb-0 nowrap"
    >
      <thead class="table-light">
        <tr>
          <th class="text-nowrap">NAMA</th>
          <th class="text-nowrap">USERNAME</th>
          <th class="text-nowrap">NIP</th>
          <th class="text-nowrap">JABATAN</th>
          <th class="text-nowrap">UNOR</th>
          <th class="text-nowrap">EMAIL</th>
          <th class="text-nowrap">ROLE</th>
          <th class="text-center text-nowrap">STATUS</th>
          <th class="text-center text-nowrap">MFA</th>
          <th class="text-center text-nowrap">SSO</th>
          <th class="text-center text-nowrap sticky-action-th">AKSI</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.id">
          <td class="fw-semibold text-dark">
            <div class="d-flex align-items-center gap-2">
              <UiAppAvatar :name="row.nama" :size="28" />
              <span>{{ row.nama }}</span>
            </div>
          </td>
          <td class="text-nowrap font-monospace text-primary small">{{ row.username || '-' }}</td>
          <td class="text-nowrap font-monospace small">{{ row.nip }}</td>
          <td>{{ row.jabatan || '-' }}</td>
          <td>{{ row.unitKerja || '-' }}</td>
          <td>{{ row.email || '-' }}</td>
          <td class="text-nowrap">
            <span class="badge bg-light-primary text-primary border border-primary-subtle px-2 py-1">
              {{ row.roleLabel || row.role }}
            </span>
          </td>
          <td class="text-center text-nowrap">
            <span v-if="row.deletedAt" class="badge bg-secondary">Terhapus</span>
            <span
              v-else
              :class="[
                'badge',
                row.isActive
                  ? 'bg-light-success text-success border border-success-subtle'
                  : 'bg-light-danger text-danger border border-danger-subtle',
              ]"
            >
              {{ row.isActive ? 'Aktif' : 'Nonaktif' }}
            </span>
          </td>
          <td class="text-center text-nowrap">
            <span
              :class="[
                'badge',
                row.mfaEnabled
                  ? 'bg-light-success text-success border border-success-subtle'
                  : 'bg-light-secondary text-secondary border border-secondary-subtle',
              ]"
            >
              <i :class="row.mfaEnabled ? 'ti ti-shield-check me-1' : 'ti ti-shield-x me-1'"></i>
              {{ row.mfaEnabled ? 'Aktif' : 'Nonaktif' }}
            </span>
          </td>
          <td class="text-center">
            <span v-if="row.isSso" class="badge bg-light-info text-info border border-info-subtle">
              SSO
            </span>
            <span v-else class="text-muted small">Lokal</span>
          </td>
          <td class="text-center text-nowrap sticky-action-td">
            <div class="d-inline-flex gap-1 justify-content-center flex-nowrap">
              <template v-if="row.deletedAt">
                <button
                  type="button"
                  class="btn btn-sm btn-light-success d-inline-flex align-items-center gap-1"
                  title="Pulihkan akun"
                  @click="emit('restore', row)"
                >
                  <i class="ti ti-restore"></i> Pulihkan
                </button>
              </template>
              <template v-else>
                <button
                  type="button"
                  class="btn btn-sm btn-light-primary"
                  title="Edit Akun"
                  @click="emit('edit', row)"
                >
                  <i class="ti ti-edit"></i>
                </button>
                <button
                  type="button"
                  class="btn btn-sm btn-light-warning"
                  title="Reset Kata Sandi"
                  @click="emit('resetPassword', row)"
                >
                  <i class="ti ti-key"></i>
                </button>
                <template v-if="!currentUserId || row.id !== currentUserId">
                  <button
                    v-if="row.isActive"
                    type="button"
                    class="btn btn-sm btn-light-danger"
                    title="Nonaktifkan Akun"
                    @click="emit('nonaktif', row)"
                  >
                    <i class="ti ti-lock"></i>
                  </button>
                  <button
                    v-else
                    type="button"
                    class="btn btn-sm btn-light-success"
                    title="Aktifkan Akun"
                    @click="emit('aktif', row)"
                  >
                    <i class="ti ti-lock-open"></i>
                  </button>
                </template>
                <button
                  v-if="row.mfaEnabled"
                  type="button"
                  class="btn btn-sm btn-light-info"
                  title="Reset MFA Pengguna"
                  @click="emit('resetMfa', row)"
                >
                  <i class="ti ti-shield-off"></i>
                </button>
                <button
                  v-if="!currentUserId || row.id !== currentUserId"
                  type="button"
                  class="btn btn-sm btn-light-danger"
                  title="Hapus Akun"
                  @click="emit('hapus', row)"
                >
                  <i class="ti ti-trash"></i>
                </button>
              </template>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { UserRow } from '~/types/index'

interface Props {
  rows: UserRow[]
  currentUserId?: number
}

const props = withDefaults(defineProps<Props>(), {
  rows: () => [],
  currentUserId: undefined,
})

const emit = defineEmits<{
  edit: [row: UserRow]
  nonaktif: [row: UserRow]
  aktif: [row: UserRow]
  hapus: [row: UserRow]
  resetPassword: [row: UserRow]
  resetMfa: [row: UserRow]
  restore: [row: UserRow]
}>()

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
  if (win.$ && win.$.fn && win.$.fn.DataTable && win.$.fn.DataTable.isDataTable('#tabel-akun')) {
    try {
      win.$('#tabel-akun').DataTable().destroy()
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

  win.$('#tabel-akun').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [10] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ pengguna',
      infoEmpty: 'Menampilkan 0 pengguna',
      infoFiltered: '(disaring dari _MAX_ total pengguna)',
      zeroRecords: 'Tidak ada data pengguna yang sesuai filter',
      emptyTable: 'Tidak ada data pengguna yang ditemukan',
      paginate: {
        first: 'Pertama',
        last: 'Terakhir',
        next: 'Selanjutnya',
        previous: 'Sebelumnya',
      },
    },
  })

  // Cegah browser password manager mengisi otomatis (autofill) ke kolom pencarian DataTables
  // saat modal ubah/reset kata sandi dibuka
  const searchInput = document.querySelector<HTMLInputElement>('#tabel-akun_filter input')
  if (searchInput) {
    searchInput.setAttribute('autocomplete', 'off')
    searchInput.setAttribute('name', 'dt-search-tabel-akun')
  }
}

onMounted(() => {
  initOrRefreshDataTable()
})

// Ketika data berubah karena filter atau aksi CRUD, hanya perbarui isi DataTable tanpa me-reload halaman/card
watch(
  () => props.rows,
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
