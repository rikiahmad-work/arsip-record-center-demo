<template>
  <div ref="containerRef" :class="[noMargin ? '' : 'form-group mb-3', 'position-relative']">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <div class="input-group">
      <span class="input-group-text bg-light text-secondary cursor-pointer" title="Buka bantuan pemilih" @click="togglePopover">
        <i class="ti ti-calendar-time"></i>
      </span>
      <input
        :id="inputId"
        ref="inputRef"
        type="text"
        class="form-control"
        :class="[error ? 'is-invalid' : '', disabled ? 'bg-light text-muted' : '']"
        :value="modelValue"
        :placeholder="placeholder || 'Contoh: 2021 - 2022, 2026, atau 28 Juli 2021'"
        :disabled="disabled"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @keydown.esc="closePopover"
      />
      <button
        v-if="!disabled"
        type="button"
        class="btn btn-outline-secondary"
        title="Buka bantuan pemilih kurun waktu"
        @click="togglePopover"
      >
        <i class="ti ti-chevron-down" style="font-size: 14px;"></i>
      </button>
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </div>
    <div class="form-text text-muted small mt-1">
      <i class="ti ti-info-circle me-1"></i>Bisa diketik bebas berupa rentang tahun (2021 - 2022), tanggal (28 Juli 2021), atau tahun tunggal (2026).
    </div>

    <!-- Popover Helper Bantuan Kurun Waktu -->
    <div
      v-if="isOpen && !disabled"
      class="kurun-picker-popover shadow-lg border rounded bg-white p-3 position-absolute"
      style="z-index: 1050; top: 100%; left: 0; margin-top: 4px;"
    >
      <!-- Navigasi Tab Mode -->
      <ul class="nav nav-pills nav-fill mb-3 bg-light p-1 rounded small">
        <li class="nav-item">
          <button
            type="button"
            class="nav-link py-1 px-2"
            :class="{ active: activeTab === 'range' }"
            @click.stop="activeTab = 'range'"
          >
            <i class="ti ti-arrows-left-right me-1"></i>Rentang Tahun
          </button>
        </li>
        <li class="nav-item">
          <button
            type="button"
            class="nav-link py-1 px-2"
            :class="{ active: activeTab === 'year' }"
            @click.stop="activeTab = 'year'"
          >
            <i class="ti ti-calendar me-1"></i>Tahun Tunggal
          </button>
        </li>
        <li class="nav-item">
          <button
            type="button"
            class="nav-link py-1 px-2"
            :class="{ active: activeTab === 'date' }"
            @click.stop="activeTab = 'date'"
          >
            <i class="ti ti-calendar-event me-1"></i>Tanggal
          </button>
        </li>
      </ul>

      <!-- Tab 1: Rentang Tahun -->
      <div v-if="activeTab === 'range'" class="tab-pane-content">
        <div class="row g-2 align-items-center mb-2">
          <div class="col-5">
            <label class="form-label small text-muted mb-1">Tahun Awal</label>
            <input
              v-model="rangeStart"
              type="number"
              class="form-control form-control-sm text-center"
              min="1950"
              max="2100"
              placeholder="YYYY"
            />
          </div>
          <div class="col-2 text-center pt-3 text-muted">
            <span>s/d</span>
          </div>
          <div class="col-5">
            <label class="form-label small text-muted mb-1">Tahun Akhir</label>
            <input
              v-model="rangeEnd"
              type="number"
              class="form-control form-control-sm text-center"
              min="1950"
              max="2100"
              placeholder="YYYY"
            />
          </div>
        </div>

        <!-- Tombol Preset Rentang Cepat -->
        <div class="mb-3">
          <div class="small text-muted mb-1">Pilihan Cepat:</div>
          <div class="d-flex flex-wrap gap-1">
            <button
              v-for="preset in rangePresets"
              :key="preset"
              type="button"
              class="btn btn-xs btn-outline-secondary py-0 px-2"
              style="font-size: 11px;"
              @click.stop="applyPreset(preset)"
            >
              {{ preset }}
            </button>
          </div>
        </div>

        <button
          type="button"
          class="btn btn-sm btn-primary w-100"
          :disabled="!rangeStart || !rangeEnd"
          @click.stop="applyRange"
        >
          <i class="ti ti-check me-1"></i>Gunakan Rentang ({{ rangeStart }} - {{ rangeEnd }})
        </button>
      </div>

      <!-- Tab 2: Tahun Tunggal -->
      <div v-else-if="activeTab === 'year'" class="tab-pane-content">
        <div class="d-flex align-items-center justify-content-between mb-2 pb-1 border-bottom">
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary py-0 px-2"
            @click.stop="decadeStart -= 12"
          >
            <i class="ti ti-chevron-left"></i>
          </button>
          <span class="fw-semibold text-dark small">
            {{ decadeStart }} &ndash; {{ decadeStart + 11 }}
          </span>
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary py-0 px-2"
            @click.stop="decadeStart += 12"
          >
            <i class="ti ti-chevron-right"></i>
          </button>
        </div>

        <div class="row g-1 mb-2">
          <div v-for="yr in visibleYears" :key="yr" class="col-4">
            <button
              type="button"
              class="btn btn-sm w-100 py-1 text-center"
              :class="[
                String(modelValue).trim() === String(yr)
                  ? 'btn-primary text-white fw-bold'
                  : yr === currentYear
                  ? 'btn-outline-primary'
                  : 'btn-light text-dark'
              ]"
              @click.stop="selectSingleYear(yr)"
            >
              {{ yr }}
            </button>
          </div>
        </div>

        <button
          type="button"
          class="btn btn-xs btn-link text-primary w-100 text-center p-0"
          @click.stop="selectSingleYear(currentYear)"
        >
          Pilih Tahun Ini ({{ currentYear }})
        </button>
      </div>

      <!-- Tab 3: Tanggal Spesifik -->
      <div v-else-if="activeTab === 'date'" class="tab-pane-content">
        <label class="form-label small text-muted mb-1">Pilih Tanggal Dokumen:</label>
        <input
          v-model="tempDate"
          type="date"
          class="form-control form-control-sm mb-2"
          @change="applyDate"
        />
        <div v-if="tempDate" class="text-center small text-primary fw-medium mb-2">
          Hasil: {{ formatDateIndo(tempDate) }}
        </div>
        <button
          type="button"
          class="btn btn-sm btn-primary w-100"
          :disabled="!tempDate"
          @click.stop="applyDate"
        >
          <i class="ti ti-check me-1"></i>Gunakan Format Tanggal
        </button>
      </div>

      <!-- Footer -->
      <div class="d-flex justify-content-between pt-2 mt-2 border-top">
        <button
          type="button"
          class="btn btn-xs btn-link text-danger p-0 text-decoration-none"
          @click.stop="clearValue"
        >
          Hapus
        </button>
        <button
          type="button"
          class="btn btn-xs btn-link text-secondary p-0 text-decoration-none"
          @click.stop="closePopover"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

interface Props {
  modelValue?: string
  label?: string
  placeholder?: string
  error?: string
  required?: boolean
  disabled?: boolean
  noMargin?: boolean
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '',
  noMargin: false,
})

const emit = defineEmits<{
  'update:modelValue': [val: string]
}>()

const currentYear = new Date().getFullYear()
const inputId = computed(() => props.id ?? `kurunpicker-${Math.random().toString(36).slice(2, 7)}`)

const containerRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const isOpen = ref(false)
const activeTab = ref<'range' | 'year' | 'date'>('range')

// State untuk rentang tahun
const rangeStart = ref(currentYear - 1)
const rangeEnd = ref(currentYear)

const rangePresets = computed(() => [
  `${currentYear - 1} - ${currentYear}`,
  `${currentYear - 2} - ${currentYear}`,
  `${currentYear - 3} - ${currentYear}`,
  '2021 - 2022',
  '2020 - 2024',
])

// State untuk tahun tunggal
const decadeStart = ref(Math.floor(currentYear / 12) * 12)
const visibleYears = computed(() => {
  const years: number[] = []
  for (let i = 0; i < 12; i++) {
    years.push(decadeStart.value + i)
  }
  return years
})

// State untuk tanggal
const tempDate = ref('')

// Inisialisasi nilai saat dibuka
watch(
  () => isOpen.value,
  (open) => {
    if (open) {
      parseExistingValue(props.modelValue)
    }
  }
)

function parseExistingValue(val: string) {
  if (!val) return
  const trimmed = val.trim()

  // Cek apakah format rentang: "YYYY - YYYY"
  const rangeMatch = trimmed.match(/^(\d{4})\s*-\s*(\d{4})$/)
  if (rangeMatch) {
    activeTab.value = 'range'
    rangeStart.value = Number(rangeMatch[1])
    rangeEnd.value = Number(rangeMatch[2])
    return
  }

  // Cek apakah 4 digit tahun: "YYYY"
  const yearMatch = trimmed.match(/^(\d{4})$/)
  if (yearMatch) {
    activeTab.value = 'year'
    const yr = Number(yearMatch[1])
    decadeStart.value = Math.floor(yr / 12) * 12
    return
  }

  // Cek jika mengandung nama bulan Indonesia atau format tanggal
  const bulanIndo = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']
  if (bulanIndo.some((b) => trimmed.toLowerCase().includes(b.toLowerCase()))) {
    activeTab.value = 'date'
  }
}

function applyRange() {
  if (rangeStart.value && rangeEnd.value) {
    emit('update:modelValue', `${rangeStart.value} - ${rangeEnd.value}`)
    closePopover()
  }
}

function applyPreset(preset: string) {
  emit('update:modelValue', preset)
  closePopover()
}

function selectSingleYear(yr: number) {
  emit('update:modelValue', String(yr))
  closePopover()
}

const BULAN_NAMES = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

function formatDateIndo(dateStr: string): string {
  if (!dateStr) return ''
  const [y, m, d] = dateStr.split('-')
  if (!y || !m || !d) return dateStr
  const monthIdx = parseInt(m, 10) - 1
  const monthName = BULAN_NAMES[monthIdx] || m
  return `${parseInt(d, 10)} ${monthName} ${y}`
}

function applyDate() {
  if (tempDate.value) {
    emit('update:modelValue', formatDateIndo(tempDate.value))
    closePopover()
  }
}

function clearValue() {
  emit('update:modelValue', '')
  closePopover()
}

function openPopover() {
  if (props.disabled) return
  parseExistingValue(props.modelValue)
  isOpen.value = true
}

function closePopover() {
  isOpen.value = false
}

function togglePopover() {
  if (props.disabled) return
  if (isOpen.value) {
    closePopover()
  } else {
    openPopover()
  }
}

function handleClickOutside(event: MouseEvent) {
  if (containerRef.value && !containerRef.value.contains(event.target as Node)) {
    closePopover()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
.kurun-picker-popover {
  animation: fadeIn 0.15s ease-in-out;
  min-width: 280px;
  max-width: 380px;
  width: calc(100vw - 2rem);
}
@media (min-width: 480px) {
  .kurun-picker-popover {
    width: 380px;
  }
}
.btn-xs {
  font-size: 11px;
  line-height: 1.4;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
