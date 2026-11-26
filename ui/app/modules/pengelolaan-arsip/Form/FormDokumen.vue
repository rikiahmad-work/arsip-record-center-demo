<template>
  <div class="row g-3">
    <div class="col-md-6">
      <label class="form-label fw-semibold">File Dokumen Digitalisasi (PDF / JPG / PNG, Maks 10 MB)</label>
      <input
        type="file"
        class="form-control"
        accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png"
        @change="onFileChange"
      >
      <small v-if="selectedFileName" class="text-success d-block mt-1">
        <i class="ti ti-paperclip me-1"></i>File dipilih: <strong>{{ selectedFileName }}</strong>
      </small>
      <small v-else class="text-muted d-block mt-1">
        Opsional saat tambah arsip; dapat juga diunggah melalui halaman Detail Arsip.
      </small>
    </div>
    <div class="col-md-6">
      <FormInput
        :model-value="keteranganDokumen"
        label="Keterangan Dokumen Pendukung"
        placeholder="Contoh: Hasil scan surat keputusan asli"
        @update:model-value="emit('update:keteranganDokumen', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  file?: File | null
  keteranganDokumen?: string
}

const props = withDefaults(defineProps<Props>(), {
  file: null,
  keteranganDokumen: '',
})

const emit = defineEmits<{
  'update:file': [value: File | null]
  'update:keteranganDokumen': [value: string]
}>()

const selectedFileName = computed(() => props.file?.name || '')

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files?.[0] || null
  emit('update:file', f)
}
</script>
