<template>
  <div class="p-3 bg-light rounded border mb-3">
    <div class="row g-2 align-items-center">
      <!-- Search -->
      <div class="col-12 col-md-4 col-xl-3">
        <div class="input-group input-group-sm">
          <span class="input-group-text bg-white"><i class="ti ti-search"></i></span>
          <input
            :value="search"
            type="text"
            class="form-control"
            placeholder="Cari No. Item, Uraian, Kode, Boks/Folder..."
            @input="emit('update:search', ($event.target as HTMLInputElement).value)"
          >
        </div>
      </div>

      <!-- Filter Kode Klasifikasi -->
      <div class="col-6 col-md-3 col-xl-2">
        <select
          :value="kode"
          class="form-select form-select-sm"
          @change="emit('update:kode', ($event.target as HTMLSelectElement).value)"
        >
          <option value="">Semua Kode (435)</option>
          <option v-for="opt in kodeOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <!-- Filter Kategori Arsip -->
      <div class="col-6 col-md-2 col-xl-2">
        <select
          :value="kategori"
          class="form-select form-select-sm"
          @change="emit('update:kategori', ($event.target as HTMLSelectElement).value)"
        >
          <option value="">Semua Kategori</option>
          <option v-for="opt in kategoriOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <!-- Filter SKKAD -->
      <div class="col-6 col-md-2 col-xl-2">
        <select
          :value="skkad"
          class="form-select form-select-sm"
          @change="emit('update:skkad', ($event.target as HTMLSelectElement).value)"
        >
          <option value="">Semua SKKAD</option>
          <option v-for="opt in skkadOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <!-- Tombol Aksi Toolbar -->
      <div class="col-6 col-md-12 col-xl-3 d-flex justify-content-xl-end gap-2 flex-wrap">
        <button
          v-if="hasActiveFilter"
          type="button"
          class="btn btn-sm btn-outline-secondary"
          @click="emit('reset')"
        >
          <i class="ti ti-refresh me-1"></i>Reset
        </button>
        <button
          v-if="canWrite"
          type="button"
          class="btn btn-sm btn-light-primary"
          @click="emit('import')"
        >
          <i class="ti ti-upload me-1"></i>Import Excel/CSV
        </button>
        <button
          type="button"
          class="btn btn-sm btn-light-success"
          @click="emit('export')"
        >
          <i class="ti ti-download me-1"></i>Export DAI
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SelectOption } from '~/types/arsip'
import { KATEGORI_ARSIP_OPTIONS, SKKAD_OPTIONS } from '~/constants/arsip'

interface Props {
  search?: string
  kode?: string
  kategori?: string
  skkad?: string
  kodeOptions?: SelectOption[]
  canWrite?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  search: '',
  kode: '',
  kategori: '',
  skkad: '',
  kodeOptions: () => [],
  canWrite: true,
})

const kategoriOptions = KATEGORI_ARSIP_OPTIONS
const skkadOptions = SKKAD_OPTIONS

const hasActiveFilter = computed(() =>
  Boolean(props.search || props.kode || props.kategori || props.skkad)
)

const emit = defineEmits<{
  'update:search': [value: string]
  'update:kode': [value: string]
  'update:kategori': [value: string]
  'update:skkad': [value: string]
  reset: []
  tambah: []
  import: []
  export: []
}>()
</script>
