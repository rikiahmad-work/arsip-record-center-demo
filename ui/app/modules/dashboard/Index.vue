<template>
  <!-- Page Header -->
  <TemplatePageTitle title="Dashboard Record Center" :breadcrumbs="[{ label: 'Dashboard' }]" />

  <!-- Filter Bar Bulan & Tahun sesuai docs/template-ui/UI-admin/dashboard/index.html -->
  <div class="row mb-4">
    <div class="col-12 d-flex justify-content-md-end justify-content-start">
      <div class="d-flex flex-wrap gap-2 shadow-sm p-2 rounded bg-white border align-items-center w-100 w-md-auto">
        <div class="input-group dashboard-filter-bulan flex-grow-1 flex-md-grow-0" style="min-width: 140px;">
          <span class="input-group-text bg-light text-secondary">
            <i class="ti ti-calendar-event"></i>
          </span>
          <select v-model="filterBulan" class="form-select">
            <option value="">Semua Bulan</option>
            <option v-for="b in bulanOptions" :key="b.value" :value="b.value">{{ b.label }}</option>
          </select>
        </div>

        <div class="input-group dashboard-filter-tahun flex-grow-1 flex-md-grow-0" style="min-width: 120px;">
          <span class="input-group-text bg-light text-secondary">
            <i class="ti ti-calendar-time"></i>
          </span>
          <select v-model="filterTahun" class="form-select">
            <option value="2026">2026</option>
            <option value="2025">2025</option>
            <option value="2024">2024</option>
            <option value="2023">2023</option>
          </select>
        </div>

        <button class="btn btn-primary d-flex align-items-center justify-content-center gap-1 flex-grow-1 flex-md-grow-0" type="button" :disabled="loading" @click="loadDashboardData">
          <span v-if="loading" class="spinner-border spinner-border-sm"></span>
          <i v-else class="ti ti-search"></i>
          <span>Cari</span>
        </button>
      </div>
    </div>
  </div>

  <!-- 4 Stat Cards Utama -->
  <ModulesDashboardStatCards :stats="dashboardStats" />

  <!-- Row 1: Tren Arsip Masuk Bulanan (Bar) + Distribusi per Unit Kerja (Horizontal Bar) -->
  <div class="row mt-2">
    <div class="col-xl-7">
      <ModulesDashboardChartArsipMasuk :data="chartData" />
    </div>
    <div class="col-xl-5">
      <div class="card h-100">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Arsip Berdasarkan Unit Kerja Pencipta</h5>
          <span class="badge bg-light-primary text-primary">{{ perUnitKerja.length }} Unit</span>
        </div>
        <div class="card-body">
          <div ref="refChartUnit" style="min-height: 300px;"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Row 2: Pie Chart Kategori Arsip + Donut Chart Keamanan SKKAD + Ringkasan Cepat -->
  <div class="row mt-4">
    <div class="col-xl-4">
      <div class="card h-100">
        <div class="card-header">
          <h5 class="mb-0">Distribusi Kategori Arsip</h5>
        </div>
        <div class="card-body d-flex justify-content-center align-items-center" style="min-height: 290px;">
          <div ref="refChartKategori" class="w-100"></div>
        </div>
      </div>
    </div>

    <div class="col-xl-4">
      <div class="card h-100">
        <div class="card-header">
          <h5 class="mb-0">Klasifikasi Keamanan (SKKAD)</h5>
        </div>
        <div class="card-body d-flex justify-content-center align-items-center" style="min-height: 290px;">
          <div ref="refChartSKKAD" class="w-100"></div>
        </div>
      </div>
    </div>

    <div class="col-xl-4">
      <div class="card h-100">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Ringkasan Master & Operasional</h5>
          <NuxtLink to="/pengelolaan-arsip" class="btn btn-sm btn-light-primary">Lihat Arsip</NuxtLink>
        </div>
        <div class="card-body">
          <ul class="list-group list-group-flush">
            <li class="list-group-item d-flex justify-content-between align-items-center px-0">
              <span class="text-muted"><i class="ti ti-tag me-2 text-primary"></i>Total Kode Klasifikasi (JRA)</span>
              <span class="badge bg-primary rounded-pill">{{ ringkasan.totalKodeKlasifikasi }}</span>
            </li>
            <li class="list-group-item d-flex justify-content-between align-items-center px-0">
              <span class="text-muted"><i class="ti ti-map-pin me-2 text-success"></i>Lokasi Simpan Aktif</span>
              <span class="badge bg-success rounded-pill">{{ ringkasan.lokasiSimpanAktif }}</span>
            </li>
            <li class="list-group-item d-flex justify-content-between align-items-center px-0">
              <span class="text-muted"><i class="ti ti-building me-2 text-info"></i>Unit Kerja Pencipta</span>
              <span class="badge bg-info rounded-pill">{{ ringkasan.totalUnitKerja }}</span>
            </li>
            <li class="list-group-item d-flex justify-content-between align-items-center px-0">
              <span class="text-muted"><i class="ti ti-archive me-2 text-warning"></i>Total Transaksi Peminjaman</span>
              <span class="badge bg-warning rounded-pill">{{ ringkasan.totalPeminjaman }}</span>
            </li>
            <li class="list-group-item d-flex justify-content-between align-items-center px-0 border-bottom-0">
              <span class="text-muted"><i class="ti ti-trash me-2 text-danger"></i>Pengajuan Penyusutan (JRA)</span>
              <span class="badge bg-danger rounded-pill">{{ ringkasan.totalPenyusutan }}</span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { DUMMY_DASHBOARD_STATS, DUMMY_CHART_DATA } from '~/constants/arsip'
import type { DashboardStats, ChartDataPoint } from '~/types/arsip'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Dashboard | Arsip Record Center KemenPPPA' })

const { getDashboardStats, getChartData, getLaporanRekap, getKodeKlasifikasi, getLokasiSimpan, getUnitKerja } = useArsip()

const filterBulan = ref('')
const filterTahun = ref('2026')
const loading = ref(false)

const bulanOptions = [
  { value: '01', label: 'Januari' },
  { value: '02', label: 'Februari' },
  { value: '03', label: 'Maret' },
  { value: '04', label: 'April' },
  { value: '05', label: 'Mei' },
  { value: '06', label: 'Juni' },
  { value: '07', label: 'Juli' },
  { value: '08', label: 'Agustus' },
  { value: '09', label: 'September' },
  { value: '10', label: 'Oktober' },
  { value: '11', label: 'November' },
  { value: '12', label: 'Desember' },
]

const dashboardStats = ref<DashboardStats>({
  ...DUMMY_DASHBOARD_STATS,
  arsipVital: DUMMY_DASHBOARD_STATS.arsipAktif,
})
const chartData = ref<ChartDataPoint[]>([...DUMMY_CHART_DATA])

const perUnitKerja = ref<{ label: string; jumlah: number }[]>([])
const perKategori = ref<{ label: string; jumlah: number }[]>([])
const perSKKAD = ref<{ label: string; jumlah: number }[]>([])

const ringkasan = reactive({
  totalKodeKlasifikasi: 435,
  lokasiSimpanAktif: 2,
  totalUnitKerja: 10,
  totalPeminjaman: 0,
  totalPenyusutan: 0,
})

const refChartUnit = ref<HTMLElement | null>(null)
const refChartKategori = ref<HTMLElement | null>(null)
const refChartSKKAD = ref<HTMLElement | null>(null)
let extraCharts: any[] = []

async function waitApexCharts(): Promise<any> {
  const win = window as any
  for (let i = 0; i < 50; i++) {
    if (win.ApexCharts) return win.ApexCharts
    await new Promise((r) => setTimeout(r, 100))
  }
  return win.ApexCharts ?? null
}

function destroyExtraCharts() {
  extraCharts.forEach((c) => {
    try {
      c.destroy()
    } catch {}
  })
  extraCharts = []
  ;[refChartUnit, refChartKategori, refChartSKKAD].forEach((r) => {
    if (r.value) r.value.innerHTML = ''
  })
}

async function renderExtraCharts() {
  await nextTick()
  const ApexCharts = await waitApexCharts()
  if (!ApexCharts) return

  destroyExtraCharts()
  await nextTick()

  // 1. Horizontal Bar Chart — Unit Kerja
  if (refChartUnit.value) {
    const unitItems = perUnitKerja.value.length > 0
      ? perUnitKerja.value
      : [
          { label: 'Biro Data dan Informasi', jumlah: 3 },
          { label: 'Biro Umum & SDM', jumlah: 2 },
        ]
    const cUnit = new ApexCharts(refChartUnit.value, {
      chart: { type: 'bar', height: 300, toolbar: { show: false } },
      colors: ['#2ca87f'],
      plotOptions: { bar: { borderRadius: 4, horizontal: true, barHeight: '55%' } },
      dataLabels: { enabled: true },
      series: [{ name: 'Jumlah Arsip', data: unitItems.map((u) => Number(u.jumlah || 0)) }],
      xaxis: { categories: unitItems.map((u) => u.label || 'Lainnya') },
    })
    await cUnit.render()
    extraCharts.push(cUnit)
  }

  // 2. Pie Chart — Kategori Arsip
  if (refChartKategori.value) {
    const katItems = perKategori.value.length > 0
      ? perKategori.value
      : [
          { label: 'Inaktif', jumlah: dashboardStats.value.arsipInaktif || 1 },
          { label: 'Vital', jumlah: dashboardStats.value.arsipVital || 1 },
          { label: 'Terjaga', jumlah: dashboardStats.value.arsipTerjaga || 1 },
        ]
    const cKat = new ApexCharts(refChartKategori.value, {
      chart: { type: 'pie', height: 280 },
      labels: katItems.map((k) => k.label || 'Inaktif'),
      series: katItems.map((k) => Number(k.jumlah || 0)),
      colors: ['#1890ff', '#e58a00', '#2ca87f', '#dc2626', '#6b7280'],
      legend: { position: 'bottom' },
    })
    await cKat.render()
    extraCharts.push(cKat)
  }

  // 3. Donut Chart — SKKAD
  if (refChartSKKAD.value) {
    const skkadItems = perSKKAD.value.length > 0
      ? perSKKAD.value
      : [
          { label: 'Biasa', jumlah: 1 },
          { label: 'Terbatas', jumlah: 1 },
          { label: 'Rahasia', jumlah: 1 },
        ]
    const cSkkad = new ApexCharts(refChartSKKAD.value, {
      chart: { type: 'donut', height: 280 },
      labels: skkadItems.map((s) => s.label || 'Biasa'),
      series: skkadItems.map((s) => Number(s.jumlah || 0)),
      colors: ['#2ca87f', '#e58a00', '#dc2626'],
      legend: { position: 'bottom' },
    })
    await cSkkad.render()
    extraCharts.push(cSkkad)
  }
}

async function loadDashboardData() {
  loading.value = true
  try {
    const [stats, chart, rekap, kodes, lokasis, units] = await Promise.all([
      getDashboardStats(),
      getChartData(Number(filterTahun.value) || 2026),
      getLaporanRekap({ tahun: filterTahun.value }),
      getKodeKlasifikasi(),
      getLokasiSimpan(),
      getUnitKerja(),
    ])

    dashboardStats.value = stats
    chartData.value = chart

    if (kodes.length > 0) ringkasan.totalKodeKlasifikasi = kodes.length
    if (lokasis.length > 0) ringkasan.lokasiSimpanAktif = lokasis.filter((l) => l.statusKey === 'aktif').length
    if (units.length > 0) ringkasan.totalUnitKerja = units.length

    if (rekap) {
      ringkasan.totalPeminjaman = Number(rekap.total_peminjaman || 0)
      ringkasan.totalPenyusutan = Number(rekap.total_penyusutan || 0)
      perUnitKerja.value = Array.isArray(rekap.per_unit_kerja) ? rekap.per_unit_kerja : []
      perKategori.value = Array.isArray(rekap.per_kategori) ? rekap.per_kategori : []
      perSKKAD.value = Array.isArray(rekap.per_skkad) ? rekap.per_skkad : []
    }

    await renderExtraCharts()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})

onBeforeUnmount(() => {
  destroyExtraCharts()
})
</script>

<style scoped>
.list-group-item {
  border-left: none;
  border-right: none;
  font-size: 0.875rem;
}
@media (min-width: 768px) {
  .dashboard-filter-bulan {
    width: 200px !important;
  }
  .dashboard-filter-tahun {
    width: 160px !important;
  }
}
@media (max-width: 767.98px) {
  .dashboard-filter-bulan,
  .dashboard-filter-tahun {
    width: 100% !important;
  }
}
</style>
