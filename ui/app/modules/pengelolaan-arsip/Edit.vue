<template>
  <TemplatePageTitle :title="`Edit Arsip — ${arsip?.noItemArsip || ''}`" :breadcrumbs="[
    { label: 'Pengelolaan Arsip', to: '/pengelolaan-arsip' },
    { label: 'Edit Arsip' },
  ]">
    <template #action>
      <button class="btn btn-outline-secondary" @click="$router.push(`/pengelolaan-arsip/${id}`)">
        <i class="ti ti-arrow-left me-1"></i>Kembali
      </button>
    </template>
  </TemplatePageTitle>

  <div v-if="loading" class="text-center py-5">
    <div class="spinner-border text-primary" role="status">
      <span class="visually-hidden">Memuat...</span>
    </div>
    <p class="mt-2 text-muted">Memuat data arsip...</p>
  </div>

  <div v-else-if="!arsip" class="text-center py-5 text-muted">
    <i class="ti ti-inbox" style="font-size: 48px;"></i>
    <p class="mt-3">Arsip tidak ditemukan</p>
  </div>

  <form v-else @submit.prevent="simpan">
    <div class="card">
      <div class="card-header">
        <h5 class="mb-0"><i class="ti ti-file-edit me-2"></i>Edit Data Arsip Record Center</h5>
      </div>
      <div class="card-body">
        <!-- Section A: Identitas -->
        <h6 class="text-primary mb-3"><i class="ti ti-file-description me-1"></i>A. Identitas Arsip</h6>
        <ModulesPengelolaanArsipFormFormIdentitasArsip
          :form="formA"
          :kode-options="kodeOptions"
          @update:noItemArsip="formA.noItemArsip = $event"
          @update:kodeKlasifikasi="formA.kodeKlasifikasi = $event"
          @update:jenisArsip="formA.jenisArsip = $event"
          @update:tingkatPerkembangan="formA.tingkatPerkembangan = $event"
          @update:uraianInformasiArsip="formA.uraianInformasiArsip = $event"
        />

        <hr class="my-4">

        <!-- Section B: Detail Fisik -->
        <h6 class="text-primary mb-3"><i class="ti ti-box me-1"></i>B. Detail Fisik &amp; Kondisi</h6>
        <ModulesPengelolaanArsipFormFormDetailFisik
          :form="formB"
          @update:kurunWaktu="formB.kurunWaktu = $event"
          @update:jumlah="formB.jumlah = $event"
          @update:keterangan="formB.keterangan = $event"
          @update:nomorDefinitif="formB.nomorDefinitif = $event"
        />

        <hr class="my-4">

        <!-- Section C: Lokasi -->
        <h6 class="text-primary mb-3"><i class="ti ti-map-pin me-1"></i>C. Lokasi Penyimpanan (Rak, Boks, Folder)</h6>
        <ModulesPengelolaanArsipFormFormLokasiSimpan
          :form="formC"
          :lokasi-options="lokasiOptions"
          @update:lokasiSimpan="formC.lokasiSimpan = $event"
          @update:noBoks="formC.noBoks = $event"
          @update:noRak="formC.noRak = $event"
          @update:noFolder="formC.noFolder = $event"
        />

        <hr class="my-4">

        <!-- Section D: Retensi -->
        <h6 class="text-primary mb-3"><i class="ti ti-clock me-1"></i>D. Jangka Simpan &amp; Retensi (JRA &amp; SKKAD)</h6>
        <ModulesPengelolaanArsipFormFormRetensi
          :form="formD"
          :can-override-skkad="isAdmin()"
          @update:jangkaSimpanInaktif="formD.jangkaSimpanInaktif = $event"
          @update:jangkaSimpanNasibAkhir="formD.jangkaSimpanNasibAkhir = $event"
          @update:skkad="formD.skkad = $event"
          @update:kategoriArsip="formD.kategoriArsip = $event"
        />

        <hr class="my-4">

        <!-- Section E: Pencipta -->
        <h6 class="text-primary mb-3"><i class="ti ti-building me-1"></i>E. Pencipta Arsip</h6>
        <ModulesPengelolaanArsipFormFormPenciptaArsip
          :form="formE"
          :unor-options="unitOptions"
          @update:unorId="formE.unorId = $event"
          @update:petugasInput="formE.petugasInput = $event"
        />

        <hr class="my-4">

        <!-- Section F: Tambah Dokumen Pendukung Baru -->
        <h6 class="text-primary mb-3"><i class="ti ti-paperclip me-1"></i>F. Unggah Dokumen Pendukung Tambahan</h6>
        <ModulesPengelolaanArsipFormFormDokumen
          :file="formF.file"
          :keterangan-dokumen="formF.keterangan"
          @update:file="formF.file = $event"
          @update:keteranganDokumen="formF.keterangan = $event"
        />
      </div>

      <div class="card-footer d-flex gap-2 justify-content-end">
        <button type="button" class="btn btn-outline-secondary" @click="$router.push(`/pengelolaan-arsip/${id}`)">
          Batal
        </button>
        <button type="submit" class="btn btn-primary" :disabled="simpanLoading">
          <span v-if="simpanLoading" class="spinner-border spinner-border-sm me-1"></span>
          <i v-else class="ti ti-device-floppy me-1"></i>Simpan Perubahan
        </button>
      </div>
    </div>
  </form>
</template>

<script setup lang="ts">
import type { Arsip, KodeKlasifikasi, LokasiSimpan, UnitKerja } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

import ModulesPengelolaanArsipFormFormIdentitasArsip from '~/modules/pengelolaan-arsip/Form/FormIdentitasArsip.vue'
import ModulesPengelolaanArsipFormFormDetailFisik from '~/modules/pengelolaan-arsip/Form/FormDetailFisik.vue'
import ModulesPengelolaanArsipFormFormLokasiSimpan from '~/modules/pengelolaan-arsip/Form/FormLokasiSimpan.vue'
import ModulesPengelolaanArsipFormFormRetensi from '~/modules/pengelolaan-arsip/Form/FormRetensi.vue'
import ModulesPengelolaanArsipFormFormPenciptaArsip from '~/modules/pengelolaan-arsip/Form/FormPenciptaArsip.vue'
import ModulesPengelolaanArsipFormFormDokumen from '~/modules/pengelolaan-arsip/Form/FormDokumen.vue'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const route = useRoute()
const router = useRouter()
const { isAdmin } = useSession()
const { getArsipById, updateArsip, uploadDokumenArsip, getKodeKlasifikasi, getLokasiSimpan, getUnitKerja } = useArsip()
const { showToast } = useToast()

const id = route.params.id as string
const arsip = ref<Arsip | null>(null)
const loading = ref(true)
const simpanLoading = ref(false)
const initialLoadDone = ref(false)

useHead({
  title: computed(() => `Edit Arsip — ${arsip.value?.noItemArsip || ''} | Arsip Record Center KemenPPPA`),
})

const masterKodes = ref<KodeKlasifikasi[]>([])
const masterLokasis = ref<LokasiSimpan[]>([])
const masterUnits = ref<UnitKerja[]>([])

const kodeOptions = computed(() =>
  masterKodes.value.map((k) => ({ value: k.kode, label: `${k.kode} — ${k.nama}` }))
)
const lokasiOptions = computed(() =>
  masterLokasis.value.map((l) => ({ value: l.namaLokasi, label: l.namaLokasi }))
)
const unitOptions = computed(() =>
  masterUnits.value.map((u) => ({ value: u.id, label: `${u.kodeUnit} — ${u.namaUnit}` }))
)

const formA = reactive({
  noItemArsip: '',
  kodeKlasifikasi: '',
  jenisArsip: 'Konvensional',
  tingkatPerkembangan: 'Asli',
  uraianInformasiArsip: '',
})

const formB = reactive({
  kurunWaktu: '',
  jumlah: '',
  keterangan: '',
  nomorDefinitif: '',
})

const formC = reactive({
  lokasiSimpan: '',
  noBoks: '',
  noRak: '',
  noFolder: '',
})

const formD = reactive({
  jangkaSimpanInaktif: '',
  jangkaSimpanNasibAkhir: '',
  skkad: '',
  kategoriArsip: '',
})

const formE = reactive({
  unorId: '',
  penciptaArsip: '',
  petugasInput: 'Petugas Arsip',
})

const formF = reactive<{ file: File | null; keterangan: string }>({
  file: null,
  keterangan: '',
})

onMounted(async () => {
  try {
    const [data, kodes, lokasis, units] = await Promise.all([
      getArsipById(id),
      getKodeKlasifikasi(),
      getLokasiSimpan(),
      getUnitKerja(),
    ])

    arsip.value = data
    masterKodes.value = kodes
    masterLokasis.value = lokasis
    masterUnits.value = units

    if (data) {
      formA.noItemArsip = data.noItemArsip
      formA.kodeKlasifikasi = data.kodeKlasifikasi
      formA.jenisArsip = data.jenisArsip
      formA.tingkatPerkembangan = data.tingkatPerkembangan
      formA.uraianInformasiArsip = data.uraianInformasi || data.uraianInformasiArsip

      formB.kurunWaktu = data.kurunWaktu
      formB.jumlah = data.jumlah
      formB.keterangan = data.keterangan
      formB.nomorDefinitif = data.nomorDefinitif || ''

      formC.lokasiSimpan = data.lokasiSimpan
      formC.noBoks = data.noBoks
      formC.noRak = data.noRak
      formC.noFolder = data.noFolder || ''

      formD.jangkaSimpanInaktif = data.jangkaSimpanInaktif
      formD.jangkaSimpanNasibAkhir = (data.nasibAkhir || data.jangkaSimpanNasibAkhir || 'Musnah') as string
      formD.skkad = data.skkad
      formD.kategoriArsip = data.kategoriArsip

      const matchedUnit = units.find((u) => u.namaUnit === data.penciptaArsip || Number(u.id) === data.unitKerjaId)
      formE.unorId = matchedUnit ? matchedUnit.id : (units[0]?.id || '')
      formE.penciptaArsip = data.penciptaArsip
      formE.petugasInput = data.petugasInput || 'Petugas Arsip'
    }
  } finally {
    loading.value = false
    nextTick(() => {
      initialLoadDone.value = true
    })
  }
})

watch(
  () => formA.kodeKlasifikasi,
  (selectedKode) => {
    if (!initialLoadDone.value) return
    const found = masterKodes.value.find((k) => k.kode === selectedKode)
    if (found) {
      formD.jangkaSimpanInaktif = found.jangkaSimpanInaktif || '3 Tahun'
      formD.jangkaSimpanNasibAkhir = found.nasibAkhir || 'Musnah'
      formD.skkad = found.skkad || 'Biasa'
    }
  }
)

watch(
  () => formE.unorId,
  (val) => {
    const found = masterUnits.value.find((u) => u.id === val)
    if (found) formE.penciptaArsip = found.namaUnit
  }
)

function buildPayload() {
  const selectedLokasi = masterLokasis.value.find((l) => l.namaLokasi === formC.lokasiSimpan)
  const selectedUnit = masterUnits.value.find((u) => u.id === formE.unorId)

  return {
    no_item_arsip: formA.noItemArsip,
    no_folder: formC.noFolder || null,
    kode_klasifikasi: formA.kodeKlasifikasi,
    uraian_informasi: formA.uraianInformasiArsip,
    uraian_informasi_arsip: formA.uraianInformasiArsip,
    jenis_arsip: formA.jenisArsip,
    kurun_waktu: formB.kurunWaktu,
    tingkat_perkembangan: formA.tingkatPerkembangan,
    jumlah: formB.jumlah,
    keterangan: formB.keterangan || undefined,
    nomor_definitif: formB.nomorDefinitif || null,
    lokasi_simpan_id: selectedLokasi ? Number(selectedLokasi.id) : undefined,
    lokasi_simpan: formC.lokasiSimpan,
    no_boks: formC.noBoks,
    no_rak: formC.noRak,
    jangka_simpan_inaktif: formD.jangkaSimpanInaktif,
    nasib_akhir: formD.jangkaSimpanNasibAkhir || 'Musnah',
    skkad: formD.skkad,
    kategori_arsip: formD.kategoriArsip,
    unit_kerja_id: selectedUnit ? Number(selectedUnit.id) : undefined,
    pencipta_arsip: formE.penciptaArsip,
    status: 'final',
  }
}

async function simpan() {
  simpanLoading.value = true
  try {
    await updateArsip(id, buildPayload())
    if (formF.file) {
      await uploadDokumenArsip(id, formF.file, formF.keterangan)
    }
    showToast('Arsip berhasil diperbarui', 'success')
    router.push(`/pengelolaan-arsip/${id}`)
  } catch (err: any) {
    showToast(err?.message || 'Gagal memperbarui arsip', 'danger')
  } finally {
    simpanLoading.value = false
  }
}
</script>
