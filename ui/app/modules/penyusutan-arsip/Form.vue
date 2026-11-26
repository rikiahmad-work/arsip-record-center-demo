<template>
  <UiAppModal title="Ajukan Usulan Penyusutan Arsip (JRA)" size="lg" @close="emit('close')">
    <form @submit.prevent="simpan">
      <div class="row g-3">
        <div class="col-12">
          <FormSelect
            v-model="form.arsipId"
            label="Pilih Berkas Arsip yang Akan Disusutkan"
            :options="arsipSelectOptions"
            placeholder="-- Cari & Pilih Berkas Arsip --"
            search-placeholder="Ketik No. Item, Kode Klasifikasi, atau Uraian Arsip..."
            required
            searchable
            no-margin
            @update:model-value="onArsipChange"
          />
          <div v-if="errorMsg" class="text-danger small mt-1">{{ errorMsg }}</div>
        </div>

        <div v-if="selectedArsip" class="col-12">
          <div class="p-3 bg-light rounded border">
            <div class="row g-2 small">
              <div class="col-md-4">
                <span class="text-muted d-block">Kode Klasifikasi</span>
                <span class="badge bg-light-primary text-primary fw-bold">{{ selectedArsip.kodeKlasifikasi }}</span>
              </div>
              <div class="col-md-4">
                <span class="text-muted d-block">Jangka Simpan Inaktif (JRA)</span>
                <strong>{{ selectedArsip.jangkaSimpanInaktif || '-' }}</strong>
              </div>
              <div class="col-md-4">
                <span class="text-muted d-block">Nasib Akhir (JRA)</span>
                <span :class="['badge', (selectedArsip.nasibAkhir || '').toLowerCase().includes('permanen') ? 'bg-primary' : 'bg-secondary']">
                  {{ selectedArsip.nasibAkhir || selectedArsip.jangkaSimpanNasibAkhir || 'Musnah' }}
                </span>
              </div>
              <div class="col-12 mt-1">
                <span class="text-muted d-block">Uraian Informasi Arsip</span>
                <span class="fw-medium text-dark">{{ selectedArsip.uraianInformasi || selectedArsip.uraianInformasiArsip }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="col-md-6">
          <FormSelect
            v-model="form.jenisUsul"
            label="Jenis Usulan Penyusutan"
            :options="jenisUsulOptions"
            required
            :searchable="false"
            no-margin
          />
          <small class="text-muted">Otomatis disesuaikan dengan Nasib Akhir JRA (Musnah / Serah Permanen)</small>
        </div>

        <div class="col-md-6">
          <label class="form-label">Tanggal Jatuh Tempo Retensi <span class="text-danger">*</span></label>
          <input v-model="form.tanggalRetensi" type="date" class="form-control" required>
        </div>

        <div class="col-12">
          <label class="form-label">Alasan / Dasar Pengajuan Penyusutan</label>
          <textarea
            v-model="form.alasan"
            class="form-control"
            rows="3"
            placeholder="Contoh: Masa retensi inaktif telah habis sesuai JRA KemenPPPA dan tidak memiliki nilai guna primer/sekunder lagi"
          ></textarea>
        </div>
      </div>
    </form>

    <template #footer>
      <button type="button" class="btn btn-outline-secondary" @click="emit('close')">Batal</button>
      <button type="button" class="btn btn-primary" :disabled="submitting" @click="simpan">
        <span v-if="submitting" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="ti ti-send me-1"></i>Ajukan Penyusutan
      </button>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { Arsip } from '~/types/arsip'

interface Props {
  arsipList: Arsip[]
  submitting?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  arsipList: () => [],
  submitting: false,
})

const emit = defineEmits<{
  close: []
  save: [payload: Record<string, any>]
}>()

const today = new Date().toISOString().slice(0, 10)

const form = reactive({
  arsipId: '',
  jenisUsul: 'Musnah',
  tanggalRetensi: today,
  alasan: '',
})

const errorMsg = ref('')

const jenisUsulOptions = [
  { value: 'Musnah', label: 'Usul Musnah (Pemusnahan Arsip)' },
  { value: 'Serah', label: 'Usul Serah (Penyerahan Arsip Statis/Permanen ke ANRI)' },
]

const arsipSelectOptions = computed(() =>
  props.arsipList.map((a) => ({
    value: a.id,
    label: `${a.noItemArsip} [${a.kodeKlasifikasi}] — ${(a.uraianInformasi || a.uraianInformasiArsip).slice(0, 75)} (${a.nasibAkhir || 'Musnah'})`,
  }))
)

const selectedArsip = computed(() =>
  props.arsipList.find((a) => String(a.id) === String(form.arsipId))
)

function onArsipChange(val: string) {
  const arsip = props.arsipList.find((a) => String(a.id) === String(val))
  if (!arsip) return
  const nasib = (arsip.nasibAkhir || arsip.jangkaSimpanNasibAkhir || '').toLowerCase()
  form.jenisUsul = nasib.includes('permanen') ? 'Serah' : 'Musnah'
}

function simpan() {
  errorMsg.value = ''
  if (!form.arsipId) {
    errorMsg.value = 'Silakan pilih berkas arsip yang akan diajukan'
    return
  }

  emit('save', {
    arsip_id: Number(form.arsipId),
    jenis_usul: form.jenisUsul,
    tanggal_retensi: form.tanggalRetensi,
    alasan: form.alasan || null,
  })
}
</script>
