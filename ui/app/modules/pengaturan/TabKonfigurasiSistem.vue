<template>
  <div>
    <h6 class="mb-3">Konfigurasi Umum Sistem</h6>

    <form @submit.prevent="simpan">
      <div class="row g-3">
        <div class="col-md-6">
          <label class="form-label">Nama Aplikasi</label>
          <input v-model="form.namaAplikasi" class="form-control" placeholder="Contoh: SiARSIP Record Center">
        </div>

        <div class="col-md-6">
          <label class="form-label">Default Jangka Simpan Inaktif</label>
          <input v-model="form.defaultJangkaSimpanInaktif" class="form-control" placeholder="Contoh: 5 Tahun">
        </div>

        <div class="col-md-6">
          <label class="form-label">Default Nasib Akhir</label>
          <select v-model="form.defaultNasibAkhir" class="form-select">
            <option value="Musnah">Musnah</option>
            <option value="Permanen">Permanen</option>
          </select>
        </div>

        <div class="col-md-6">
          <label class="form-label">Notifikasi Jatuh Tempo (hari)</label>
          <input v-model.number="form.notifikasiJatuhTempoHari" type="number" min="1" class="form-control">
        </div>

        <div class="col-12">
          <label class="form-label">Logo Aplikasi</label>
          <div class="d-flex align-items-center gap-3">
            <img
              v-if="form.logoAplikasi"
              :src="form.logoAplikasi"
              alt="Logo"
              style="height: 48px; width: auto; object-fit: contain"
            >
            <input type="file" accept="image/*" class="form-control" @change="onLogoChange">
          </div>
          <small class="text-muted">Format gambar (PNG/JPG). Maksimal 1 MB.</small>
        </div>
      </div>

      <div class="mt-4 d-flex justify-content-end">
        <button type="submit" class="btn btn-primary" :disabled="saving">
          <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
          <i v-else class="ti ti-device-floppy me-1"></i>
          Simpan Konfigurasi
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { KonfigurasiSistem } from '~/types/arsip'

interface Props {
  initial?: KonfigurasiSistem | null
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  initial: null,
  saving: false,
})

const emit = defineEmits<{
  save: [payload: KonfigurasiSistem]
}>()

const form = reactive<KonfigurasiSistem>({
  namaAplikasi: 'SiARSIP Record Center',
  defaultJangkaSimpanInaktif: '5 Tahun',
  defaultNasibAkhir: 'Musnah',
  notifikasiJatuhTempoHari: 30,
  logoAplikasi: null,
})

watch(
  () => props.initial,
  (value) => {
    if (!value) return
    form.namaAplikasi = value.namaAplikasi ?? 'SiARSIP Record Center'
    form.defaultJangkaSimpanInaktif = value.defaultJangkaSimpanInaktif ?? '5 Tahun'
    form.defaultNasibAkhir = value.defaultNasibAkhir ?? 'Musnah'
    form.notifikasiJatuhTempoHari = Number(value.notifikasiJatuhTempoHari ?? 30)
    form.logoAplikasi = value.logoAplikasi ?? null
  },
  { immediate: true },
)

function onLogoChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  if (file.size > 1024 * 1024) {
    alert('Ukuran logo maksimal 1 MB.')
    target.value = ''
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    form.logoAplikasi = typeof reader.result === 'string' ? reader.result : null
  }
  reader.readAsDataURL(file)
}

function simpan() {
  emit('save', { ...form })
}
</script>
