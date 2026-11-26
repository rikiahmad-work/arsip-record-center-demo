<template>
  <UiAppModal title="Catat Peminjaman Arsip Baru" size="lg" @close="emit('close')">
    <form @submit.prevent="simpan">
      <div class="row g-3">
        <div class="col-md-6">
          <label class="form-label">Nama Peminjam <span class="text-danger">*</span></label>
          <input
            v-model="form.peminjam"
            type="text"
            class="form-control"
            placeholder="Masukkan nama pegawai peminjam"
            required
          >
        </div>

        <div class="col-md-6">
          <FormSelect
            v-model="form.unitKerjaId"
            label="Unit Kerja Peminjam"
            :options="unitOptions"
            placeholder="-- Pilih Unit Kerja --"
            search-placeholder="Cari unit kerja..."
            searchable
            no-margin
          />
        </div>

        <div class="col-12">
          <FormSelect
            v-model="form.arsipId"
            label="Pilih Berkas Arsip"
            :options="arsipSelectOptions"
            placeholder="-- Cari & Pilih Berkas Arsip --"
            search-placeholder="Ketik No. Item, Kode Klasifikasi, atau Uraian Arsip..."
            required
            searchable
            no-margin
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
                <span class="text-muted d-block">Lokasi Simpan</span>
                <strong>{{ selectedArsip.lokasiSimpan }} (Rak {{ selectedArsip.noRak || '-' }} / Boks {{ selectedArsip.noBoks || '-' }})</strong>
              </div>
              <div class="col-md-4">
                <span class="text-muted d-block">Klasifikasi Keamanan</span>
                <span class="badge bg-secondary">{{ selectedArsip.skkad }}</span>
              </div>
              <div class="col-12 mt-1">
                <span class="text-muted d-block">Uraian Informasi Arsip</span>
                <span class="fw-medium text-dark">{{ selectedArsip.uraianInformasi || selectedArsip.uraianInformasiArsip }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="col-md-6">
          <label class="form-label">Tanggal Pinjam <span class="text-danger">*</span></label>
          <input v-model="form.tanggalPinjam" type="date" class="form-control" required>
        </div>

        <div class="col-md-6">
          <label class="form-label">Rencana Tanggal Kembali <span class="text-danger">*</span></label>
          <input v-model="form.tanggalKembaliRencana" type="date" class="form-control" required>
        </div>

        <div class="col-12">
          <label class="form-label">Keperluan / Catatan Peminjaman</label>
          <textarea
            v-model="form.catatan"
            class="form-control"
            rows="2"
            placeholder="Contoh: Diperlukan untuk pemeriksaan audit internal / referensi penyusunan kebijakan"
          ></textarea>
        </div>
      </div>
    </form>

    <template #footer>
      <button type="button" class="btn btn-outline-secondary" @click="emit('close')">Batal</button>
      <button type="button" class="btn btn-primary" :disabled="submitting" @click="simpan">
        <span v-if="submitting" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="ti ti-device-floppy me-1"></i>Simpan Peminjaman
      </button>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { Arsip, UnitKerja } from '~/types/arsip'

interface Props {
  arsipList: Arsip[]
  unitList: UnitKerja[]
  submitting?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  arsipList: () => [],
  unitList: () => [],
  submitting: false,
})

const emit = defineEmits<{
  close: []
  save: [payload: Record<string, any>]
}>()

const today = new Date().toISOString().slice(0, 10)
const nextWeek = new Date(Date.now() + 7 * 86400000).toISOString().slice(0, 10)

const form = reactive({
  arsipId: '',
  peminjam: '',
  unitKerjaId: '',
  tanggalPinjam: today,
  tanggalKembaliRencana: nextWeek,
  catatan: '',
})

const errorMsg = ref('')

const arsipSelectOptions = computed(() =>
  props.arsipList.map((a) => ({
    value: a.id,
    label: `${a.noItemArsip} [${a.kodeKlasifikasi}] — ${(a.uraianInformasi || a.uraianInformasiArsip).slice(0, 80)}`,
  }))
)

const unitOptions = computed(() =>
  props.unitList.map((u) => ({
    value: u.id,
    label: `${u.kodeUnit} — ${u.namaUnit}`,
  }))
)

const selectedArsip = computed(() =>
  props.arsipList.find((a) => String(a.id) === String(form.arsipId))
)

function simpan() {
  errorMsg.value = ''
  if (!form.arsipId) {
    errorMsg.value = 'Silakan pilih berkas arsip yang akan dipinjam'
    return
  }
  if (!form.peminjam.trim()) {
    errorMsg.value = 'Nama peminjam wajib diisi'
    return
  }

  emit('save', {
    arsip_id: Number(form.arsipId),
    peminjam: form.peminjam.trim(),
    unit_kerja_id: form.unitKerjaId ? Number(form.unitKerjaId) : null,
    tanggal_pinjam: form.tanggalPinjam,
    tanggal_kembali_rencana: form.tanggalKembaliRencana,
    catatan: form.catatan || null,
  })
}
</script>
