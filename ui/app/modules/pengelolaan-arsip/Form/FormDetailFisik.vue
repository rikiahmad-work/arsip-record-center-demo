<template>
  <div class="row g-3">
    <div class="col-md-6">
      <FormKurunWaktuPicker
        :model-value="form.kurunWaktu"
        label="Kurun Waktu"
        placeholder="Contoh: 2021 - 2022, 2026, atau 28 Juli 2021"
        required
        @update:model-value="emit('update:kurunWaktu', $event)"
      />
    </div>
    <div class="col-md-6">
      <FormInput
        :model-value="form.jumlah"
        label="Jumlah"
        placeholder="Contoh: 1 Berkas / 1 File PDF"
        required
        @update:model-value="emit('update:jumlah', $event)"
      />
    </div>
        <div class="col-md-6">
          <FormInput
            :model-value="form.keterangan"
            label="Keterangan / Kondisi"
            placeholder="Kondisi fisik arsip"
            @update:model-value="emit('update:keterangan', $event)"
          />
        </div>
        <div class="col-md-6">
          <FormInput
            :model-value="form.nomorDefinitif"
            label="Nomor Definitif"
            placeholder="Kosongkan jika belum ada"
            @update:model-value="emit('update:nomorDefinitif', $event)"
          />
        </div>
  </div>
</template>

<script setup lang="ts">
interface FormData {
  kurunWaktu: string
  jumlah: string
  keterangan: string
  nomorDefinitif: string
}

interface Props {
  form: FormData
}

defineProps<Props>()

const emit = defineEmits<{
  'update:kurunWaktu': [value: string]
  'update:jumlah': [value: string]
  'update:keterangan': [value: string]
  'update:nomorDefinitif': [value: string]
}>()
</script>
