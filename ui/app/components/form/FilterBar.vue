<template>
  <div class="row mb-4">
    <div class="col-md-12 d-flex justify-content-end">
      <div class="d-flex gap-2 shadow-sm p-2 rounded bg-white border align-items-center">

        <!-- Filter Bulan -->
        <div class="input-group" style="width: 200px;">
          <span class="input-group-text bg-light text-secondary">
            <i class="ti ti-calendar-event"></i>
          </span>
          <select v-model="selectedBulan" class="form-select" id="filterBulan">
            <option value="">Semua Bulan</option>
            <option v-for="b in bulanOptions" :key="b.value" :value="b.value">{{ b.label }}</option>
          </select>
        </div>

        <!-- Filter Tahun -->
        <div class="input-group" style="width: 160px;">
          <span class="input-group-text bg-light text-secondary">
            <i class="ti ti-calendar-time"></i>
          </span>
          <select v-model="selectedTahun" class="form-select" id="filterTahun">
            <option value="">Semua Tahun</option>
            <option v-for="t in tahunOptions" :key="t" :value="String(t)">{{ t }}</option>
          </select>
        </div>

        <!-- Tombol Cari -->
        <button class="btn btn-primary d-flex align-items-center gap-1" type="button" @click="onCari">
          <i class="ti ti-search"></i> Cari
        </button>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  startYear?: number
  endYear?: number
}

const props = withDefaults(defineProps<Props>(), {
  startYear: 2024,
  endYear: 2026,
})

const emit = defineEmits<{
  filter: [payload: { bulan: string; tahun: string }]
  reset: []
}>()

const selectedBulan = ref('')
const selectedTahun = ref('')

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

const tahunOptions = computed(() => {
  const years: number[] = []
  for (let y = props.endYear; y >= props.startYear; y--) years.push(y)
  return years
})

function onCari() {
  emit('filter', { bulan: selectedBulan.value, tahun: selectedTahun.value })
}

function onReset() {
  selectedBulan.value = ''
  selectedTahun.value = ''
  emit('reset')
}
</script>
