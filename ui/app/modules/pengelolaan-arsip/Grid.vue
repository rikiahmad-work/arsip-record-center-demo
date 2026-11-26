<template>
  <div class="dt-responsive table-responsive user-table-wrapper">
    <table
      id="tabel-pengelolaan-arsip"
      class="table table-bordered table-hover align-middle mb-0 nowrap"
    >
      <thead class="table-light">
        <tr>
          <th class="text-center" style="width: 45px;">NO</th>
          <th class="text-nowrap">NO. ITEM</th>
          <th class="text-center text-nowrap">KODE</th>
          <th style="min-width: 260px;">URAIAN INFORMASI ARSIP</th>
          <th class="text-center text-nowrap">TAHUN</th>
          <th class="text-nowrap">LOKASI (RAK / BOKS / FOLDER)</th>
          <th class="text-nowrap">RETENSI &amp; NASIB AKHIR</th>
          <th class="text-center text-nowrap">KATEGORI &amp; SKKAD</th>
          <th class="text-nowrap">PENCIPTA ARSIP</th>
          <th class="text-center text-nowrap sticky-action-th">AKSI</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in list"
          :key="item.id"
        >
          <!-- 1. No. -->
          <td class="text-center">{{ index + 1 }}</td>

          <!-- 2. No. Item & Definitif -->
          <td class="text-nowrap">
            <div class="fw-semibold text-dark">{{ item.noItemArsip }}</div>
            <small v-if="item.nomorDefinitif" class="text-muted">Def: {{ item.nomorDefinitif }}</small>
          </td>

          <!-- 3. Kode Klasifikasi -->
          <td class="text-center text-nowrap">
            <span class="badge bg-light-primary text-primary fw-bold">{{ item.kodeKlasifikasi }}</span>
          </td>

          <!-- 4. Uraian Informasi Arsip + Pembeda Status Pinjam & Penyusutan -->
          <td style="min-width: 260px; max-width: 360px; white-space: normal;">
            <div class="fw-medium text-dark">
              {{ item.uraianInformasi || item.uraianInformasiArsip }}
            </div>
            <small class="text-muted d-block mt-1">
              {{ item.jenisArsip }} &bull; {{ item.tingkatPerkembangan }} &bull; {{ item.jumlah }}
              <template v-if="item.keterangan"> &bull; {{ item.keterangan }}</template>
            </small>
            <div class="d-flex flex-wrap gap-1 mt-1">
              <span v-if="item.status === 'draft'" class="badge bg-warning text-dark">
                <i class="ti ti-file-pencil me-1"></i>Draft
              </span>
              <span
                v-if="item.sedangDipinjam"
                :class="['badge', String(item.statusPinjam).toLowerCase() === 'terlambat' ? 'bg-danger' : 'bg-warning text-dark']"
              >
                <i class="ti ti-book-upload me-1"></i>{{ String(item.statusPinjam).toLowerCase() === 'terlambat' ? 'Dipinjam (Terlambat)' : 'Sedang Dipinjam' }}
              </span>
              <span
                v-if="item.dalamPenyusutan"
                class="badge bg-dark text-white"
              >
                <i class="ti ti-recycle me-1"></i>{{ item.statusPenyusutan || 'Proses Penyusutan' }}
              </span>
            </div>
          </td>

          <!-- 5. Kurun Waktu -->
          <td class="text-center text-nowrap">{{ item.kurunWaktu }}</td>

          <!-- 6. Lokasi Simpan (Rak / Boks / Folder) -->
          <td class="text-nowrap">
            <div class="fw-medium text-dark">{{ item.lokasiSimpan || 'Record Center' }}</div>
            <small class="text-muted">
              Rak {{ item.noRak || '-' }} &bull; Boks {{ item.noBoks || '-' }} &bull; Folder {{ item.noFolder || '-' }}
            </small>
          </td>

          <!-- 7. Retensi & Nasib Akhir (Standar ANRI Cut-off 31 Desember) -->
          <td class="text-nowrap">
            <div class="d-flex align-items-center gap-1 mb-1">
              <span class="small text-dark fw-medium">{{ item.jangkaSimpanInaktif || '-' }}</span>
              <span :class="['badge', isPermanen(item.nasibAkhir || item.jangkaSimpanNasibAkhir) ? 'bg-primary' : 'bg-secondary']">
                {{ item.nasibAkhir || item.jangkaSimpanNasibAkhir || 'Musnah' }}
              </span>
            </div>
            <div v-if="item.isJatuhTempo" class="mt-1">
              <span
                :class="[
                  'badge',
                  isPermanen(item.nasibAkhir || item.jangkaSimpanNasibAkhir)
                    ? 'bg-light-info text-info border border-info'
                    : 'bg-light-danger text-danger border border-danger'
                ]"
                :title="`Cut-off Retensi ANRI: ${item.tanggalJatuhTempo || ''}`"
              >
                <i class="ti ti-alert-triangle me-1"></i>{{ item.rekomendasiRetensi || 'Jatuh Tempo Retensi' }}
              </span>
            </div>
            <small v-else-if="item.tanggalJatuhTempo" class="text-muted d-block" style="font-size: 11px;">
              Batas: {{ item.tanggalJatuhTempo }}
            </small>
          </td>

          <!-- 8. SKKAD & Kategori -->
          <td class="text-center text-nowrap">
            <div class="d-flex flex-column gap-1 align-items-center">
              <span :class="['badge', kategoriBadge(item.kategoriArsip)]">
                {{ item.kategoriArsip }}
              </span>
              <span :class="['badge', skkadBadge(item.skkad)]">
                SKKAD: {{ item.skkad || 'Biasa' }}
              </span>
            </div>
          </td>

          <!-- 9. Unit Pencipta -->
          <td class="text-nowrap">{{ item.penciptaArsip }}</td>

          <!-- 10. Aksi (Fixed / Sticky di Kanan) -->
          <td class="text-center text-nowrap sticky-action-td">
            <div class="d-inline-flex gap-1 justify-content-center flex-nowrap">
              <button
                type="button"
                class="btn btn-sm btn-light-info"
                title="Lihat Detail Lengkap (20 Field & Dokumen)"
                @click="emit('detail', item)"
              >
                <i class="ti ti-eye"></i>
              </button>
              <button
                v-if="canEdit"
                type="button"
                class="btn btn-sm btn-light-primary"
                title="Edit Arsip"
                @click="emit('edit', item)"
              >
                <i class="ti ti-pencil"></i>
              </button>
              <button
                v-if="canDelete"
                type="button"
                class="btn btn-sm btn-light-danger"
                title="Hapus Arsip"
                @click="emit('hapus', item)"
              >
                <i class="ti ti-trash"></i>
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { Arsip } from '~/types/arsip'

interface Props {
  list: Arsip[]
  canEdit?: boolean
  canDelete?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  list: () => [],
  canEdit: true,
  canDelete: true,
})

const emit = defineEmits<{
  detail: [item: Arsip]
  edit: [item: Arsip]
  hapus: [item: Arsip]
}>()

let dtInstance: any = null

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
  if (win.$ && win.$.fn && win.$.fn.DataTable && win.$.fn.DataTable.isDataTable('#tabel-pengelolaan-arsip')) {
    try {
      win.$('#tabel-pengelolaan-arsip').DataTable().destroy()
    } catch (_) {}
  }
  dtInstance = null
}

async function initOrRefreshDataTable() {
  await ensureDataTablesLoaded()
  await nextTick()

  const win = window as any
  if (!win.$ || !win.$.fn || !win.$.fn.DataTable) return

  destroyDataTable()
  await nextTick()

  dtInstance = win.$('#tabel-pengelolaan-arsip').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    columnDefs: [{ orderable: false, targets: [9] }],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ arsip',
      infoEmpty: 'Menampilkan 0 arsip',
      infoFiltered: '(disaring dari _MAX_ total arsip)',
      zeroRecords: 'Tidak ada data arsip yang sesuai filter',
      emptyTable: 'Tidak ada data arsip yang ditemukan',
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

// Ketika filter berubah, hanya re-draw DataTable tanpa me-reload komponen halaman
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

function isPermanen(val?: string) {
  return String(val || '').toLowerCase().includes('permanen')
}

function skkadBadge(val?: string) {
  const v = String(val || 'Biasa').toUpperCase()
  if (v === 'RAHASIA') return 'bg-light-danger text-danger border border-danger'
  if (v === 'TERBATAS') return 'bg-light-warning text-warning border border-warning'
  return 'bg-light-success text-success border border-success'
}

function kategoriBadge(val?: string) {
  const v = String(val || 'Inaktif').toUpperCase()
  if (v === 'VITAL') return 'bg-danger'
  if (v === 'TERJAGA') return 'bg-warning text-dark'
  if (v === 'USUL MUSNAH') return 'bg-dark'
  if (v === 'USUL SERAH') return 'bg-info'
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
