<template>
  <div class="row g-3">
    <div class="col-md-3">
      <FormSelect
        :model-value="form.lokasiSimpan"
        label="Lokasi Simpan"
        :options="lokasiOptions"
        placeholder="-- Pilih Lokasi --"
        required
        @update:model-value="emit('update:lokasiSimpan', $event)"
      />
    </div>
    <div class="col-md-3">
      <FormInput
        :model-value="form.noRak"
        label="No. Rak"
        placeholder="Contoh: R-01"
        required
        @update:model-value="emit('update:noRak', $event)"
      />
    </div>
    <div class="col-md-3">
      <FormInput
        :model-value="form.noBoks"
        label="No. Boks"
        placeholder="Contoh: B-001"
        required
        @update:model-value="emit('update:noBoks', $event)"
      />
    </div>
    <div class="col-md-3">
      <FormInput
        :model-value="form.noFolder"
        label="No. Folder"
        placeholder="Contoh: F-01"
        required
        @update:model-value="emit('update:noFolder', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SelectOption } from '~/types/arsip'

interface FormData {
  lokasiSimpan: string
  noBoks: string
  noRak: string
  noFolder: string
}

interface Props {
  form: FormData
  lokasiOptions: SelectOption[]
}

defineProps<Props>()

const emit = defineEmits<{
  'update:lokasiSimpan': [value: string]
  'update:noBoks': [value: string]
  'update:noRak': [value: string]
  'update:noFolder': [value: string]
}>()
</script>
