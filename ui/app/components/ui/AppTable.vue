<template>
  <div :class="['dt-responsive table-responsive', wrapperClass]">
    <table
      :id="tableId"
      :class="['table', striped ? '' : '', bordered ? 'table-bordered' : '', 'nowrap']"
    >
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key" :class="col.thClass">
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <slot name="rows">
          <tr v-if="!rows || rows.length === 0">
            <td :colspan="columns.length" class="text-center text-muted py-4">
              Tidak ada data
            </td>
          </tr>
          <tr v-for="(row, idx) in rows" :key="idx">
            <td v-for="col in columns" :key="col.key" :class="col.tdClass">
              <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
                {{ row[col.key] }}
              </slot>
            </td>
          </tr>
        </slot>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
export interface TableColumn {
  key: string
  label: string
  thClass?: string
  tdClass?: string
}

interface Props {
  tableId: string
  columns: TableColumn[]
  rows?: Record<string, any>[]
  striped?: boolean
  bordered?: boolean
  /** Aktifkan integrasi DataTables.js */
  datatable?: boolean
  datatableOptions?: Record<string, any>
  /** Class tambahan untuk wrapper tabel */
  wrapperClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  striped: true,
  bordered: true,
  datatable: true,
  rows: () => [],
  wrapperClass: '',
})

onMounted(async () => {
  if (!props.datatable) return

  // Load jQuery terlebih dahulu, lalu DataTables (urutan wajib)
  const scripts = [
    '/admin/js/plugins/jquery.min.js',
    '/admin/js/plugins/jquery.dataTables.min.js',
    '/admin/js/plugins/dataTables.bootstrap5.min.js',
  ]

  const loadScript = (src: string) =>
    new Promise<void>((resolve) => {
      // Jika sudah ada di DOM dan window.$ sudah tersedia (khusus jQuery), langsung resolve
      if (document.querySelector(`script[src="${src}"]`)) { resolve(); return }
      const s = document.createElement('script')
      s.src = src
      s.onload = () => resolve()
      s.onerror = () => resolve()
      document.body.appendChild(s)
    })

  const loadLink = (href: string) => {
    if (document.querySelector(`link[href="${href}"]`)) return
    const l = document.createElement('link')
    l.rel = 'stylesheet'
    l.href = href
    document.head.appendChild(l)
  }

  loadLink('/admin/css/plugins/dataTables.bootstrap5.min.css')
  await scripts.reduce((p, s) => p.then(() => loadScript(s)), Promise.resolve())

  await nextTick()

  const win = window as any
  if (win.$ && win.$(`#${props.tableId}`).DataTable) {
    win.$(`#${props.tableId}`).DataTable({
      language: {
        search:           'Cari:',
        lengthMenu:       'Tampilkan _MENU_ data',
        info:             'Menampilkan _START_ - _END_ dari _TOTAL_ data',
        infoEmpty:        'Tidak ada data',
        infoFiltered:     '(difilter dari _MAX_ total data)',
        zeroRecords:      'Data tidak ditemukan',
        paginate: {
          first:    'Pertama',
          last:     'Terakhir',
          next:     'Selanjutnya',
          previous: 'Sebelumnya',
        },
      },
      responsive: true,
      ...props.datatableOptions,
    })
  }
})

onBeforeUnmount(() => {
  const win = window as any
  if (win.$ && win.$(`#${props.tableId}`).DataTable) {
    try {
      win.$(`#${props.tableId}`).DataTable().destroy()
    } catch (_) {}
  }
})
</script>
