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
        placeholder="Kosongkan jika belum ada (Contoh: R-01)"
        @update:model-value="emit('update:noRak', $event)"
      />
    </div>
    <div class="col-md-3">
      <FormInput
        :model-value="form.noBoks"
        label="No. Boks"
        placeholder="Kosongkan jika belum ada (Contoh: B-001)"
        @update:model-value="emit('update:noBoks', $event)"
      />
    </div>
    <div class="col-md-3">
      <FormInput
        :model-value="form.noFolder"
        label="No. Folder"
        placeholder="Kosongkan jika belum ada (Contoh: F-01)"
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
