<template>
  <div>
    <div class="alert alert-info py-2 px-3 mb-3 d-flex align-items-center gap-2 small">
      <i class="ti ti-info-circle fs-5"></i>
      <span>
        <strong>Auto-fill JRA &amp; SKKAD:</strong> Jangka Simpan Inaktif, Nasib Akhir, dan SKKAD terisi otomatis dari Master 435 Kode Klasifikasi saat kode dipilih.
      </span>
    </div>

    <div class="row g-3">
      <div class="col-md-6">
        <FormInput
          :model-value="form.jangkaSimpanInaktif"
          label="Jangka Simpan Inaktif (Auto dari JRA)"
          placeholder="Otomatis dari Kode Klasifikasi"
          readonly
          required
          @update:model-value="emit('update:jangkaSimpanInaktif', $event)"
        />
      </div>
      <div class="col-md-6">
        <FormInput
          :model-value="form.jangkaSimpanNasibAkhir"
          label="Nasib Akhir (Auto dari JRA)"
          placeholder="Otomatis dari Kode Klasifikasi"
          readonly
          required
          @update:model-value="emit('update:jangkaSimpanNasibAkhir', $event)"
        />
      </div>
      <div class="col-md-6">
        <FormSelect
          :model-value="form.skkad"
          label="SKKAD (Klasifikasi Keamanan & Akses)"
          :options="skkadOptions"
          placeholder="-- Otomatis dari Kode Klasifikasi --"
          :disabled="!canOverrideSkkad"
          required
          @update:model-value="emit('update:skkad', $event)"
        />
        <small class="text-muted">Nilai standar diambil dari master kode klasifikasi (override hanya Administrator).</small>
      </div>
      <div class="col-md-6">
        <FormSelect
          :model-value="form.kategoriArsip"
          label="Kategori Arsip Record Center"
          :options="kategoriOptions"
          placeholder="-- Pilih Kategori --"
          required
          @update:model-value="emit('update:kategoriArsip', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SKKAD_OPTIONS, KATEGORI_ARSIP_OPTIONS } from '~/constants/arsip'

interface FormData {
  jangkaSimpanInaktif: string
  jangkaSimpanNasibAkhir: string
  skkad: string
  kategoriArsip: string
}

interface Props {
  form: FormData
  canOverrideSkkad?: boolean
}

withDefaults(defineProps<Props>(), {
  canOverrideSkkad: false,
})

const skkadOptions = SKKAD_OPTIONS
const kategoriOptions = KATEGORI_ARSIP_OPTIONS

const emit = defineEmits<{
  'update:jangkaSimpanInaktif': [value: string]
  'update:jangkaSimpanNasibAkhir': [value: string]
  'update:skkad': [value: string]
  'update:kategoriArsip': [value: string]
}>()
</script>
