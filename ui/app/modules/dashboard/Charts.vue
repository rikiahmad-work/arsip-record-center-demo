<template>
  <!-- Row 1: Bar Charts -->
  <div class="row">
    <div class="col-md-6">
      <div class="card">
        <div class="card-header">
          <h5>Jenis Laporan Masuk</h5>
        </div>
        <div class="card-body">
          <div ref="refJenisLaporan" style="min-height: 300px;"></div>
        </div>
      </div>
    </div>
    <div class="col-md-6">
      <div class="card">
        <div class="card-header">
          <h5>Laporan Berdasarkan Unit Kerja</h5>
        </div>
        <div class="card-body">
          <div ref="refUnitKerja" style="min-height: 300px;"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Row 2: Pie / Donut Charts -->
  <div class="row">
    <div class="col-md-6">
      <div class="card">
        <div class="card-header">
          <h5>Presentase Tindak Lanjut Laporan</h5>
        </div>
        <div class="card-body d-flex justify-content-center align-items-center" style="min-height: 320px;">
          <div ref="refKelayakan"></div>
        </div>
      </div>
    </div>
    <div class="col-md-6">
      <div class="card">
        <div class="card-header">
          <h5>Presentase Penyelesaian Laporan</h5>
        </div>
        <div class="card-body d-flex justify-content-center align-items-center" style="min-height: 320px;">
          <div ref="refPenyelesaian"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  filterBulan?: string
  filterTahun?: string
}

const props = withDefaults(defineProps<Props>(), {
  filterBulan: '',
  filterTahun: '',
})

// Template refs
const refJenisLaporan = ref<HTMLElement | null>(null)
const refUnitKerja    = ref<HTMLElement | null>(null)
const refKelayakan    = ref<HTMLElement | null>(null)
const refPenyelesaian = ref<HTMLElement | null>(null)

let apexInstances: any[] = []

// -------------------------------------------------------
// Tunggu ApexCharts tersedia di window (di-load oleh admin layout)
// Fallback: inject script sendiri jika belum ada setelah 2 detik
// -------------------------------------------------------
async function loadApexCharts(): Promise<any> {
  const win = window as any

  // Polling sampai window.ApexCharts tersedia (maks 5 detik)
  for (let i = 0; i < 50; i++) {
    if (win.ApexCharts) return win.ApexCharts
    await new Promise(r => setTimeout(r, 100))
  }

  // Fallback: inject script sendiri
  if (!win.ApexCharts) {
    await new Promise<void>((resolve) => {
      // Jika script tag sudah ada tapi belum selesai, tunggu saja via polling
      const existing = document.querySelector('script[src="/admin/js/plugins/apexcharts.min.js"]')
      if (existing) {
        // Tunggu lagi sebentar
        setTimeout(resolve, 1000)
        return
      }
      const s = document.createElement('script')
      s.src = '/admin/js/plugins/apexcharts.min.js'
      s.onload  = () => resolve()
      s.onerror = () => resolve()
      document.body.appendChild(s)
    })
  }

  if (!win.ApexCharts) {
    console.error('[DashboardCharts] ApexCharts tidak tersedia')
  }
  return win.ApexCharts ?? null
}

// -------------------------------------------------------
// Hancurkan semua instance lama dan bersihkan innerHTML
// -------------------------------------------------------
function destroyCharts() {
  apexInstances.forEach((c) => { try { c.destroy() } catch (_) {} })
  apexInstances = []

  // Kosongkan elemen agar ApexCharts bisa render ulang
  ;[refJenisLaporan, refUnitKerja, refKelayakan, refPenyelesaian].forEach((r) => {
    if (r.value) r.value.innerHTML = ''
  })
}

// -------------------------------------------------------
// Render semua chart
// -------------------------------------------------------
async function renderCharts() {
  // Pastikan DOM sudah siap
  await nextTick()

  const ApexCharts = await loadApexCharts()
  if (!ApexCharts) return

  // Hancurkan instance lama SETELAH ApexCharts confirmed tersedia
  destroyCharts()
  await nextTick()

  // ── 1. Bar Chart – Jenis Laporan Masuk ──────────────
  if (refJenisLaporan.value) {
    const chart = new ApexCharts(refJenisLaporan.value, {
      chart: { type: 'bar', height: 300, toolbar: { show: false } },
      colors: ['#1890ff'],
      plotOptions: { bar: { borderRadius: 4, horizontal: false, columnWidth: '50%' } },
      dataLabels: { enabled: false },
      series: [{ name: 'Total Laporan', data: [45, 30, 25] }],
      xaxis: { categories: ['Korupsi', 'Kolusi', 'Nepotisme'] },
    })
    await chart.render()
    apexInstances.push(chart)
  }

  // ── 2. Bar Chart – Unit Kerja (horizontal) ──────────
  if (refUnitKerja.value) {
    const chart = new ApexCharts(refUnitKerja.value, {
      chart: { type: 'bar', height: 300, toolbar: { show: false } },
      colors: ['#2ca87f'],
      plotOptions: { bar: { borderRadius: 4, horizontal: true } },
      dataLabels: { enabled: false },
      series: [{ name: 'Jumlah Laporan', data: [35, 28, 22, 15] }],
      xaxis: { categories: ['Inspektorat', 'Biro Datin', 'Biro Sumber Daya dan Organisasi', 'Biro Humas & Umum'] },
    })
    await chart.render()
    apexInstances.push(chart)
  }

  // ── 3. Pie Chart – Tindak Lanjut Laporan ────────────
  if (refKelayakan.value) {
    const chart = new ApexCharts(refKelayakan.value, {
      chart: { type: 'pie', height: 300 },
      labels: ['Ditindaklanjuti', 'Tidak Dapat Ditindaklanjuti (Ditolak)'],
      series: [85, 39],
      colors: ['#1890ff', '#e58a00'],
      legend: { position: 'bottom' },
    })
    await chart.render()
    apexInstances.push(chart)
  }

  // ── 4. Donut Chart – Penyelesaian Laporan ───────────
  if (refPenyelesaian.value) {
    const chart = new ApexCharts(refPenyelesaian.value, {
      chart: { type: 'donut', height: 300 },
      labels: ['Selesai', 'Sedang Diproses', 'Belum Diproses'],
      series: [89, 20, 15],
      colors: ['#2ca87f', '#1890ff', '#e58a00'],
      legend: { position: 'bottom' },
    })
    await chart.render()
    apexInstances.push(chart)
  }
}

onMounted(() => renderCharts())

watch([() => props.filterBulan, () => props.filterTahun], () => {
  renderCharts()
})

onBeforeUnmount(() => destroyCharts())
</script>
