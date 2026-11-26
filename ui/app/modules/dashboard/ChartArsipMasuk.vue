<template>
  <div class="card">
    <div class="card-header d-flex align-items-center justify-content-between">
      <h5 class="mb-0">Grafik Arsip Masuk per Bulan (2026)</h5>
      <span class="badge bg-light text-dark">Jan–Jul</span>
    </div>
    <div class="card-body">
      <div ref="chartRef" style="min-height: 320px;"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChartDataPoint } from '~/types/arsip'

interface Props {
  data: ChartDataPoint[]
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
})

const chartRef = ref<HTMLElement | null>(null)
let chartInstance: any = null

async function loadApexCharts(): Promise<any> {
  const win = window as any

  for (let i = 0; i < 50; i++) {
    if (win.ApexCharts) return win.ApexCharts
    await new Promise(r => setTimeout(r, 100))
  }

  if (!win.ApexCharts) {
    await new Promise<void>((resolve) => {
      const existing = document.querySelector('script[src="/admin/js/plugins/apexcharts.min.js"]')
      if (existing) {
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

  return win.ApexCharts ?? null
}

function destroyChart() {
  if (chartInstance) {
    try { chartInstance.destroy() } catch (_) {}
    chartInstance = null
  }
}

async function renderChart() {
  await nextTick()
  if (!chartRef.value) return

  const ApexCharts = await loadApexCharts()
  if (!ApexCharts) return

  destroyChart()
  if (chartRef.value) chartRef.value.innerHTML = ''
  await nextTick()

  const categories = props.data.map(d => d.bulan)
  const values = props.data.map(d => d.jumlah)

  chartInstance = new ApexCharts(chartRef.value, {
    chart: {
      type: 'bar',
      height: 320,
      toolbar: { show: false },
    },
    colors: ['#4680ff'],
    plotOptions: {
      bar: {
        borderRadius: 6,
        columnWidth: '45%',
        distributed: false,
      },
    },
    dataLabels: { enabled: false },
    series: [
      {
        name: 'Arsip Masuk',
        data: values,
      },
    ],
    xaxis: {
      categories,
      axisBorder: { show: false },
      axisTicks: { show: false },
      labels: {
        style: { colors: '#6b7280', fontSize: '12px' },
      },
    },
    yaxis: {
      labels: {
        style: { colors: '#6b7280', fontSize: '12px' },
      },
    },
    grid: {
      borderColor: '#f3f4f6',
      strokeDashArray: 4,
    },
    tooltip: {
      theme: 'light',
      y: {
        formatter: (val: number) => `${val} arsip`,
      },
    },
    fill: {
      type: 'solid',
      opacity: 1,
    },
  })

  await chartInstance.render()
}

onMounted(() => renderChart())

watch(() => props.data, () => renderChart(), { deep: true })

onBeforeUnmount(() => destroyChart())
</script>
