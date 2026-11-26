<template>
  <div ref="containerRef" :class="[noMargin ? '' : 'form-group mb-3', 'position-relative']">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <div class="input-group">
      <span class="input-group-text bg-light text-secondary cursor-pointer" @click="togglePopover">
        <i class="ti ti-calendar"></i>
      </span>
      <input
        :id="inputId"
        ref="inputRef"
        type="text"
        class="form-control"
        :class="[error ? 'is-invalid' : '', disabled ? 'bg-light text-muted' : '']"
        :value="modelValue"
        :placeholder="placeholder || 'Pilih Tahun (contoh: 2026)'"
        :disabled="disabled"
        maxlength="4"
        inputmode="numeric"
        @focus="openPopover"
        @input="onInput"
        @keydown.esc="closePopover"
      />
      <button
        v-if="!disabled"
        type="button"
        class="btn btn-outline-secondary"
        title="Buka pemilih tahun"
        @click="togglePopover"
      >
        <i class="ti ti-chevron-down" style="font-size: 14px;"></i>
      </button>
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </div>

    <!-- Popover Grid Pemilih Tahun -->
    <div
      v-if="isOpen && !disabled"
      class="year-picker-popover shadow border rounded bg-white p-3 position-absolute"
      style="z-index: 1050; min-width: 280px; width: 100%; top: 100%; left: 0; margin-top: 4px;"
    >
      <!-- Header Navigasi Dekade -->
      <div class="d-flex align-items-center justify-content-between mb-2 pb-2 border-bottom">
        <button
          type="button"
          class="btn btn-sm btn-outline-secondary py-1 px-2"
          title="12 Tahun Sebelumnya"
          @click.stop="prevDecade"
        >
          <i class="ti ti-chevron-left"></i>
        </button>
        <span class="fw-semibold text-dark fs-6">
          {{ decadeStart }} &ndash; {{ decadeStart + 11 }}
        </span>
        <button
          type="button"
          class="btn btn-sm btn-outline-secondary py-1 px-2"
          title="12 Tahun Berikutnya"
          @click.stop="nextDecade"
        >
          <i class="ti ti-chevron-right"></i>
        </button>
      </div>

      <!-- Grid 12 Tahun (3 kolom x 4 baris) -->
      <div class="row g-2 mb-3">
        <div v-for="yr in visibleYears" :key="yr" class="col-4">
          <button
            type="button"
            class="btn btn-sm w-100 py-2 text-center"
            :class="[
              String(modelValue) === String(yr)
                ? 'btn-primary text-white fw-bold shadow-sm'
                : yr === currentYear
                ? 'btn-outline-primary fw-medium'
                : 'btn-light text-dark'
            ]"
            @click.stop="selectYear(yr)"
          >
            {{ yr }}
          </button>
        </div>
      </div>

      <!-- Footer Tombol Cepat -->
      <div class="d-flex justify-content-between pt-2 border-top">
        <button
          type="button"
          class="btn btn-sm btn-link text-decoration-none p-0 text-primary"
          @click.stop="selectYear(currentYear)"
        >
          <i class="ti ti-clock me-1"></i>Tahun Ini ({{ currentYear }})
        </button>
        <button
          type="button"
          class="btn btn-sm btn-link text-decoration-none p-0 text-secondary"
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
  modelValue?: string | number
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
const inputId = computed(() => props.id ?? `yearpicker-${Math.random().toString(36).slice(2, 7)}`)

const containerRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const isOpen = ref(false)

// Dekade awal tampilan (misal 2016-2027 atau 2020-2031)
const baseYear = computed(() => {
  const val = Number(props.modelValue)
  if (!isNaN(val) && val >= 1900 && val <= 2100) {
    return val
  }
  return currentYear
})

const decadeStart = ref(Math.floor(baseYear.value / 12) * 12)

watch(
  () => props.modelValue,
  (newVal) => {
    const val = Number(newVal)
    if (!isNaN(val) && val >= 1900 && val <= 2100) {
      decadeStart.value = Math.floor(val / 12) * 12
    }
  }
)

const visibleYears = computed(() => {
  const years: number[] = []
  for (let i = 0; i < 12; i++) {
    years.push(decadeStart.value + i)
  }
  return years
})

function prevDecade() {
  decadeStart.value -= 12
}

function nextDecade() {
  decadeStart.value += 12
}

function selectYear(yr: number) {
  emit('update:modelValue', String(yr))
  isOpen.value = false
}

function onInput(e: Event) {
  const target = e.target as HTMLInputElement
  // Hanya izinkan angka
  const cleaned = target.value.replace(/\D/g, '').slice(0, 4)
  target.value = cleaned
  emit('update:modelValue', cleaned)

  if (cleaned.length === 4) {
    const parsed = Number(cleaned)
    if (!isNaN(parsed) && parsed >= 1900 && parsed <= 2100) {
      decadeStart.value = Math.floor(parsed / 12) * 12
    }
  }
}

function openPopover() {
  if (props.disabled) return
  const val = Number(props.modelValue)
  if (!isNaN(val) && val >= 1900 && val <= 2100) {
    decadeStart.value = Math.floor(val / 12) * 12
  } else {
    decadeStart.value = Math.floor(currentYear / 12) * 12
  }
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

// Click outside handling
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
.year-picker-popover {
  animation: fadeIn 0.15s ease-in-out;
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
