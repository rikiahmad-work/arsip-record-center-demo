<template>
  <div class="row g-3">
        <div class="col-md-4">
          <FormInput
            :model-value="form.noItemArsip"
            label="No. Item Arsip"
            placeholder="Otomatis"
            :readonly="!!form.noItemArsip"
            @update:model-value="emit('update:noItemArsip', $event)"
          />
        </div>
        <div class="col-md-4">
          <FormSelect
            :model-value="form.kodeKlasifikasi"
            label="Kode Klasifikasi"
            :options="kodeOptions"
            placeholder="-- Pilih Kode --"
            required
            @update:model-value="emit('update:kodeKlasifikasi', $event)"
          />
        </div>
        <div class="col-md-4">
          <FormSelect
            :model-value="form.jenisArsip"
            label="Jenis Arsip"
            :options="jenisOptions"
            placeholder="-- Pilih Jenis --"
            required
            @update:model-value="emit('update:jenisArsip', $event)"
          />
        </div>
        <div class="col-md-4">
          <FormSelect
            :model-value="form.tingkatPerkembangan"
            label="Tingkat Perkembangan"
            :options="tingkatOptions"
            placeholder="-- Pilih Tingkat --"
            required
            @update:model-value="emit('update:tingkatPerkembangan', $event)"
          />
        </div>
        <div class="col-md-8">
          <FormTextarea
            :model-value="form.uraianInformasiArsip"
            label="Uraian Informasi Arsip"
            placeholder="Deskripsi/uraian isi arsip"
            required
            @update:model-value="emit('update:uraianInformasiArsip', $event)"
          />
        </div>
  </div>
</template>

<script setup lang="ts">
import type { SelectOption } from '~/types/arsip'
import { JENIS_ARSIP_OPTIONS, TINGKAT_PERKEMBANGAN_OPTIONS } from '~/constants/arsip'

interface FormData {
  noItemArsip: string
  kodeKlasifikasi: string
  jenisArsip: string
  tingkatPerkembangan: string
  uraianInformasiArsip: string
}

interface Props {
  form: FormData
  kodeOptions: SelectOption[]
}

defineProps<Props>()

const jenisOptions = JENIS_ARSIP_OPTIONS
const tingkatOptions = TINGKAT_PERKEMBANGAN_OPTIONS

const emit = defineEmits<{
  'update:noItemArsip': [value: string]
  'update:kodeKlasifikasi': [value: string]
  'update:jenisArsip': [value: string]
  'update:tingkatPerkembangan': [value: string]
  'update:uraianInformasiArsip': [value: string]
}>()
</script>
