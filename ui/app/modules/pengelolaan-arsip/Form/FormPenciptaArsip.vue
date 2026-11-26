<template>
  <div class="row g-3">
        <div class="col-md-6">
          <FormSelect
            :model-value="form.unorId"
            label="Unit Kerja / Pencipta Arsip"
            :options="unorOptions"
            placeholder="-- Pilih Unit Kerja --"
            required
            @update:model-value="emit('update:unorId', $event)"
          />
        </div>
        <div class="col-md-6">
          <FormInput
            :model-value="form.petugasInput"
            label="Petugas Input"
            placeholder="Otomatis terisi"
            readonly
          />
        </div>
  </div>
</template>

<script setup lang="ts">
import type { SelectOption } from '~/types/arsip'

interface FormData {
  unorId: string
  petugasInput: string
}

interface Props {
  form: FormData
  unorOptions: SelectOption[]
}

defineProps<Props>()

const emit = defineEmits<{
  'update:unorId': [value: string]
  'update:petugasInput': [value: string]
}>()
</script>
