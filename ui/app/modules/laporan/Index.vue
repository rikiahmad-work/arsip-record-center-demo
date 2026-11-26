<template>
  <div>
    <TemplatePageTitle title="Laporan &amp; Rekapitulasi Arsip" :breadcrumbs="[{ label: 'Laporan Arsip' }]">
      <template #action>
        <button
          type="button"
          class="btn btn-outline-primary d-inline-flex align-items-center gap-1"
          @click="cetakHalaman"
        >
          <i class="ti ti-printer"></i>
          <span>Cetak Ringkasan</span>
        </button>
      </template>
    </TemplatePageTitle>

    <!-- 4 Stat Cards Rekapitulasi -->
    <div class="row mb-4">
      <div v-for="s in statCards" :key="s.label" class="col-sm-6 col-xl-3">
        <div class="card mb-0" :class="s.bg">
          <div class="card-body py-3">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h4 class="mb-0" :class="s.textClass">{{ s.value }}</h4>
                <small :class="s.textClass + '-50'">{{ s.label }}</small>
              </div>
              <i :class="[s.icon, s.textClass]" style="font-size: 1.75rem;"></i>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 2 Grafik Rekapitulasi (Per Kode Klasifikasi & Per Unit Kerja) -->
    <div class="row mb-4">
      <div class="col-xl-6">
        <div class="card h-100 mb-0">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Distribusi per Kode Klasifikasi (Top 10)</h5>
            <span class="badge bg-light-primary text-primary">{{ perKode.length }} Kode Aktif</span>
          </div>
          <div class="card-body">
            <div ref="refChartKode" style="min-height: 290px;"></div>
          </div>
        </div>
      </div>

      <div class="col-xl-6">
        <div class="card h-100 mb-0">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Distribusi per Unit Kerja &amp; Keamanan</h5>
            <span class="badge bg-light-success text-success">{{ perUnitKerja.length }} Unit Kerja</span>
          </div>
          <div class="card-body">
            <div ref="refChartUnit" style="min-height: 290px;"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Card Filter & Unduh / Export Laporan DAI 18 Kolom -->
    <div class="row">
      <div class="col-sm-12">
        <div class="card">
          <div class="card-body">
            <div class="mb-3">
              <h5 class="mb-1"><i class="ti ti-file-download me-1 text-primary"></i>Unduh Laporan Arsip Inaktif (DAI 18 Kolom)</h5>
              <p class="text-muted small mb-0">
                Pilih kriteria filter di bawah ini untuk mengunduh laporan Daftar Arsip Inaktif berformat resmi 18 kolom ANRI (Excel / CSV).
              </p>
            </div>

            <!-- Bar Filter Dedicated & Tombol Export -->
            <div class="p-3 bg-light rounded border">
              <div class="row g-2 align-items-end">
                <div class="col-12 col-md-3">
                  <label class="form-label small text-muted fw-semibold mb-1">
                    <i class="ti ti-tag me-1"></i>Kode Klasifikasi
                  </label>
                  <FormSelect
                    v-model="filter.kode"
                    :options="kodeSelectOptions"
                    placeholder="Semua Kode Klasifikasi"
                    search-placeholder="Cari kode atau nama..."
                    size="sm"
                    searchable
                    clearable
                    no-margin
                  />
                </div>

                <div class="col-6 col-md-3">
                  <label class="form-label small text-muted fw-semibold mb-1">
                    <i class="ti ti-building me-1"></i>Unit Kerja Pencipta
                  </label>
                  <FormSelect
                    v-model="filter.unitKerja"
                    :options="unitSelectOptions"
                    placeholder="Semua Unit Kerja"
                    search-placeholder="Cari unit kerja..."
                    size="sm"
                    searchable
                    clearable
                    no-margin
                  />
                </div>

                <div class="col-6 col-md-2">
                  <label class="form-label small text-muted fw-semibold mb-1">
                    <i class="ti ti-category me-1"></i>Kategori Arsip
                  </label>
                  <FormSelect
                    v-model="filter.kategori"
                    :options="kategoriOptions"
                    placeholder="Semua Kategori"
                    search-placeholder="Cari kategori..."
                    size="sm"
                    searchable
                    clearable
                    no-margin
                  />
                </div>

                <div class="col-6 col-md-2">
                  <label class="form-label small text-muted fw-semibold mb-1">
                    <i class="ti ti-shield-lock me-1"></i>Keamanan (SKKAD)
                  </label>
                  <FormSelect
                    v-model="filter.skkad"
                    :options="skkadOptions"
                    placeholder="Semua SKKAD"
                    search-placeholder="Cari SKKAD..."
                    size="sm"
                    searchable
                    clearable
                    no-margin
                  />
                </div>

                <!-- Tombol Aksi: Export & Reset Filter Berdampingan -->
                <div class="col-12 col-md-2 d-flex gap-1">
                  <button
                    type="button"
                    class="btn btn-sm btn-success flex-grow-1 d-inline-flex align-items-center justify-content-center gap-1 text-nowrap"
                    :disabled="exporting"
                    @click="exportDAICSV"
                  >
                    <span v-if="exporting" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                    <i v-else class="ti ti-download"></i>
                    <span>{{ exporting ? 'Mengunduh...' : 'Unduh Laporan' }}</span>
                  </button>
                  <button
                    v-if="hasActiveFilter"
                    type="button"
                    class="btn btn-sm btn-outline-danger d-inline-flex align-items-center justify-content-center"
                    title="Reset Filter"
                    @click="resetFilter"
                  >
                    <i class="ti ti-filter-off"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- Loader Progress saat Proses Export Berlangsung -->
            <div v-if="exporting" class="text-center py-4 mt-3 bg-light rounded border border-success-subtle">
              <div class="spinner-border text-success" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
              <h6 class="mt-2 mb-1 text-success fw-semibold">Sedang Menyiapkan Berkas Laporan...</h6>
              <p class="text-muted small mb-0">Sistem sedang mengompilasi data inventaris 18 kolom DAI sesuai filter yang dipilih.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DAIRow, KodeKlasifikasi, UnitKerja } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Laporan Arsip | Arsip Record Center KemenPPPA' })

const { apiFetch } = useSession()
const { getLaporanRekap, getLaporanDAI, getKodeKlasifikasi, getUnitKerja } = useArsip()
const { showToast } = useToast()

const allRows = ref<DAIRow[]>([])
const kodeList = ref<KodeKlasifikasi[]>([])
const unitList = ref<UnitKerja[]>([])
const initialLoading = ref(true)
const exporting = ref(false)

const rekapSummary = reactive({
  totalArsip: 0,
  totalPeminjaman: 0,
  totalPenyusutan: 0,
})

const perKode = ref<{ label: string; jumlah: number }[]>([])
const perUnitKerja = ref<{ label: string; jumlah: number }[]>([])

const filter = reactive({
  kode: '',
  unitKerja: '',
  kategori: '',
  skkad: '',
})

const hasActiveFilter = computed(() =>
  Boolean(filter.kode || filter.unitKerja || filter.kategori || filter.skkad)
)

const kodeSelectOptions = computed(() =>
  kodeList.value.map((k) => ({
    value: k.kode,
    label: `${k.kode} — ${k.nama}`,
  }))
)

const unitSelectOptions = computed(() =>
  unitList.value.map((u) => ({
    value: u.namaUnit,
    label: `${u.kodeUnit} — ${u.namaUnit}`,
  }))
)

const kategoriOptions = [
  { value: 'Inaktif', label: 'Inaktif' },
  { value: 'Vital', label: 'Vital' },
  { value: 'Terjaga', label: 'Terjaga' },
  { value: 'Usul Musnah', label: 'Usul Musnah' },
  { value: 'Usul Serah', label: 'Usul Serah' },
]

const skkadOptions = [
  { value: 'Biasa', label: 'Biasa' },
  { value: 'Terbatas', label: 'Terbatas' },
  { value: 'Rahasia', label: 'Rahasia' },
]

const filteredRows = computed(() =>
  allRows.value.filter((r) => {
    if (filter.kode && r.kode_klasifikasi !== filter.kode) return false
    if (filter.unitKerja && r.pencipta_arsip !== filter.unitKerja) return false
    if (filter.kategori && (r.kategori_arsip || '').toLowerCase() !== filter.kategori.toLowerCase()) return false
    if (filter.skkad && (r.skkad || '').toLowerCase() !== filter.skkad.toLowerCase()) return false
    return true
  })
)

const statCards = computed(() => [
  {
    label: 'Total Arsip Terdaftar',
    value: rekapSummary.totalArsip || allRows.value.length,
    bg: 'bg-primary text-white',
    textClass: 'text-white',
    icon: 'ti ti-archive',
  },
  {
    label: 'Arsip Sesuai Filter',
    value: filteredRows.value.length,
    bg: 'bg-success text-white',
    textClass: 'text-white',
    icon: 'ti ti-filter',
  },
  {
    label: 'Total Transaksi Peminjaman',
    value: rekapSummary.totalPeminjaman,
    bg: 'bg-info text-white',
    textClass: 'text-white',
    icon: 'ti ti-book',
  },
  {
    label: 'Total Usulan Penyusutan',
    value: rekapSummary.totalPenyusutan,
    bg: 'bg-warning text-white',
    textClass: 'text-white',
    icon: 'ti ti-trash',
  },
])

const refChartKode = ref<HTMLElement | null>(null)
const refChartUnit = ref<HTMLElement | null>(null)
let chartKodeInstance: any = null
let chartUnitInstance: any = null

async function loadApexCharts(): Promise<any> {
  const win = window as any
  if (win.ApexCharts) return win.ApexCharts

  return new Promise((resolve) => {
    const existing = document.querySelector('script[src*="apexcharts"]')
    if (existing) {
      existing.addEventListener('load', () => resolve((window as any).ApexCharts))
      if ((window as any).ApexCharts) resolve((window as any).ApexCharts)
      return
    }
    const script = document.createElement('script')
    script.src = '/admin/js/plugins/apexcharts.min.js'
    script.onload = () => resolve((window as any).ApexCharts)
    document.body.appendChild(script)
  })
}

function destroyCharts() {
  if (chartKodeInstance) {
    try {
      chartKodeInstance.destroy()
    } catch (_) {}
    chartKodeInstance = null
  }
  if (chartUnitInstance) {
    try {
      chartUnitInstance.destroy()
    } catch (_) {}
    chartUnitInstance = null
  }
}

async function renderCharts() {
  const ApexCharts = await loadApexCharts()
  if (!ApexCharts) return

  destroyCharts()
  await nextTick()

  if (refChartKode.value && perKode.value.length) {
    const items = perKode.value.slice(0, 10)
    chartKodeInstance = new ApexCharts(refChartKode.value, {
      chart: { type: 'bar', height: 290, toolbar: { show: false } },
      series: [{ name: 'Jumlah Berkas', data: items.map((i) => i.jumlah) }],
      colors: ['#4680ff'],
      plotOptions: { bar: { horizontal: true, borderRadius: 4, barHeight: '60%' } },
      dataLabels: { enabled: true, style: { fontSize: '11px' } },
      xaxis: { categories: items.map((i) => i.label) },
    })
    chartKodeInstance.render()
  }

  if (refChartUnit.value && perUnitKerja.value.length) {
    const items = perUnitKerja.value.slice(0, 8)
    chartUnitInstance = new ApexCharts(refChartUnit.value, {
      chart: { type: 'bar', height: 290, toolbar: { show: false } },
      series: [{ name: 'Jumlah Berkas', data: items.map((i) => i.jumlah) }],
      colors: ['#2ca87f'],
      plotOptions: { bar: { horizontal: false, borderRadius: 4, columnWidth: '55%' } },
      dataLabels: { enabled: false },
      xaxis: { categories: items.map((i) => i.label) },
    })
    chartUnitInstance.render()
  }
}

async function loadData() {
  initialLoading.value = true
  try {
    const [rekap, dai, kodes, units] = await Promise.all([
      getLaporanRekap(),
      getLaporanDAI(),
      getKodeKlasifikasi(),
      getUnitKerja(),
    ])

    allRows.value = dai
    kodeList.value = kodes
    unitList.value = units

    if (rekap) {
      rekapSummary.totalArsip = Number(rekap.total_arsip || dai.length)
      rekapSummary.totalPeminjaman = Number(rekap.total_peminjaman || 0)
      rekapSummary.totalPenyusutan = Number(rekap.total_penyusutan || 0)
      perKode.value = Array.isArray(rekap.per_kode) ? rekap.per_kode : []
      perUnitKerja.value = Array.isArray(rekap.per_unit_kerja) ? rekap.per_unit_kerja : []
    }
  } catch (err: any) {
    showToast(err?.message || 'Gagal memuat laporan arsip', 'danger')
  } finally {
    initialLoading.value = false
    await renderCharts()
  }
}

onMounted(() => loadData())
onBeforeUnmount(() => destroyCharts())

function resetFilter() {
  filter.kode = ''
  filter.unitKerja = ''
  filter.kategori = ''
  filter.skkad = ''
}

function cetakHalaman() {
  window.print()
}

async function exportDAICSV() {
  exporting.value = true
  try {
    const qs = new URLSearchParams()
    if (filter.kode) qs.set('kode', filter.kode)
    if (filter.unitKerja) qs.set('unit_kerja', filter.unitKerja)
    if (filter.kategori) qs.set('kategori', filter.kategori)
    if (filter.skkad) qs.set('skkad', filter.skkad)

    const res = await apiFetch(`/laporan/export${qs.toString() ? `?${qs.toString()}` : ''}`)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `rekap-laporan-arsip-18-kolom-${new Date().toISOString().slice(0, 10)}.csv`
    link.click()
    URL.revokeObjectURL(url)
    showToast('Laporan 18 kolom berhasil diekspor ke CSV', 'success')
  } catch {
    showToast('Gagal mengekspor laporan', 'danger')
  } finally {
    exporting.value = false
  }
}
</script>
