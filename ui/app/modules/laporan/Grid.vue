<template>
  <div class="dt-responsive table-responsive">
    <table
      id="tabel-laporan-dai"
      class="table table-bordered table-hover align-middle mb-0"
    >
      <thead class="table-light">
        <tr>
          <th class="text-center" style="width: 45px;">1. NO</th>
          <th class="text-nowrap">2. NO. ITEM</th>
          <th class="text-center text-nowrap">3. KODE</th>
          <th style="min-width: 260px;">4. URAIAN INFORMASI ARSIP</th>
          <th class="text-nowrap">5. JENIS</th>
          <th class="text-center text-nowrap">6. KURUN WAKTU</th>
          <th class="text-nowrap">7. TINGKAT</th>
          <th class="text-nowrap">8. JUMLAH</th>
          <th class="text-nowrap">9. KETERANGAN</th>
          <th class="text-nowrap">10. NO. DEFINITIF</th>
          <th class="text-nowrap">11. LOKASI SIMPAN</th>
          <th class="text-center text-nowrap">12. BOKS</th>
          <th class="text-center text-nowrap">13. RAK</th>
          <th class="text-nowrap">14. RETENSI INAKTIF</th>
          <th class="text-center text-nowrap">15. NASIB AKHIR</th>
          <th class="text-center text-nowrap">16. SKKAD</th>
          <th class="text-center text-nowrap">17. KATEGORI</th>
          <th class="text-nowrap">18. PENCIPTA ARSIP</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, idx) in rows" :key="`${row.no_item_arsip}-${idx}`">
          <td class="text-center">{{ idx + 1 }}</td>
          <td class="text-nowrap fw-semibold">{{ row.no_item_arsip }}</td>
          <td class="text-center text-nowrap">
            <span class="badge bg-light-primary text-primary fw-bold">{{ row.kode_klasifikasi }}</span>
          </td>
          <td style="min-width: 260px; max-width: 340px; white-space: normal;">
            {{ row.uraian_informasi }}
          </td>
          <td class="text-nowrap">{{ row.jenis_arsip }}</td>
          <td class="text-center text-nowrap">{{ row.kurun_waktu }}</td>
          <td class="text-nowrap">{{ row.tingkat_perkembangan }}</td>
          <td class="text-nowrap">{{ row.jumlah }}</td>
          <td class="text-nowrap">{{ row.keterangan || '-' }}</td>
          <td class="text-nowrap">{{ row.nomor_definitif || '-' }}</td>
          <td class="text-nowrap">{{ row.lokasi_simpan }}</td>
          <td class="text-center text-nowrap">{{ row.no_boks || '-' }}</td>
          <td class="text-center text-nowrap">{{ row.no_rak || '-' }}</td>
          <td class="text-nowrap">{{ row.jangka_simpan_inaktif || '-' }}</td>
          <td class="text-center text-nowrap">
            <span :class="['badge', (row.nasib_akhir || '').toLowerCase().includes('permanen') ? 'bg-primary' : 'bg-secondary']">
              {{ row.nasib_akhir || 'Musnah' }}
            </span>
          </td>
          <td class="text-center text-nowrap">
            <span :class="['badge', skkadBadge(row.skkad)]">{{ row.skkad }}</span>
          </td>
          <td class="text-center text-nowrap">
            <span :class="['badge', kategoriBadge(row.kategori_arsip)]">{{ row.kategori_arsip }}</span>
          </td>
          <td class="text-nowrap">{{ row.pencipta_arsip }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { DAIRow } from '~/types/arsip'

interface Props {
  rows: DAIRow[]
}

const props = withDefaults(defineProps<Props>(), {
  rows: () => [],
})

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
  if (win.$ && win.$.fn?.DataTable && win.$.fn.DataTable.isDataTable('#tabel-laporan-dai')) {
    try {
      win.$('#tabel-laporan-dai').DataTable().destroy()
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

  win.$('#tabel-laporan-dai').DataTable({
    pageLength: 10,
    lengthMenu: [10, 25, 50, 100],
    order: [[0, 'asc']],
    language: {
      search: 'Cari:',
      lengthMenu: 'Tampilkan _MENU_ data',
      info: 'Menampilkan _START_ - _END_ dari _TOTAL_ berkas',
      infoEmpty: 'Menampilkan 0 berkas',
      infoFiltered: '(disaring dari _MAX_ total berkas)',
      zeroRecords: 'Tidak ada data laporan arsip yang sesuai filter',
      emptyTable: 'Tidak ada data arsip untuk ditampilkan',
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
  () => props.rows,
  async () => {
    destroyDataTable()
    await nextTick()
    await initDataTable()
  },
  { flush: 'pre' }
)

onBeforeUnmount(() => destroyDataTable())

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
