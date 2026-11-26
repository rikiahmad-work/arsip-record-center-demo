<template>
  <TemplatePageTitle title="Tambah Arsip Inaktif" :breadcrumbs="[
    { label: 'Pengelolaan Arsip', to: '/pengelolaan-arsip' },
    { label: 'Tambah Arsip' },
  ]">
    <template #action>
      <button class="btn btn-outline-secondary" @click="$router.push('/pengelolaan-arsip')">
        <i class="ti ti-arrow-left me-1"></i>Kembali
      </button>
    </template>
  </TemplatePageTitle>

  <form @submit.prevent="simpanFinal">
    <div class="card">
      <div class="card-header">
        <h5 class="mb-0"><i class="ti ti-file-plus me-2"></i>Formulir Registrasi Arsip Record Center</h5>
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

        <!-- Section C: Lokasi (+ No. Folder) -->
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

        <!-- Section D: Retensi (Auto-fill JRA & SKKAD) -->
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

        <!-- Section F: Dokumen Pendukung (Digitalisasi — Nota Dinas #8) -->
        <h6 class="text-primary mb-3"><i class="ti ti-paperclip me-1"></i>F. Dokumen Pendukung (Digitalisasi Arsip)</h6>
        <ModulesPengelolaanArsipFormFormDokumen
          :file="formF.file"
          :keterangan-dokumen="formF.keterangan"
          @update:file="formF.file = $event"
          @update:keteranganDokumen="formF.keterangan = $event"
        />
      </div>

      <!-- Action Bar -->
      <div class="card-footer d-flex gap-2 justify-content-end">
        <button type="button" class="btn btn-outline-secondary" @click="$router.push('/pengelolaan-arsip')">
          Batal
        </button>
        <button type="button" class="btn btn-warning" :disabled="simpanLoading" @click="simpanDraf">
          <i class="ti ti-file me-1"></i>Simpan Draf
        </button>
        <button type="submit" class="btn btn-primary" :disabled="simpanLoading">
          <span v-if="simpanLoading" class="spinner-border spinner-border-sm me-1"></span>
          <i v-else class="ti ti-device-floppy me-1"></i>Simpan Data
        </button>
      </div>
    </div>
  </form>
</template>

<script setup lang="ts">
import type { KodeKlasifikasi, LokasiSimpan, UnitKerja } from '~/types/arsip'
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

useHead({ title: 'Tambah Arsip | Arsip Record Center KemenPPPA' })

const router = useRouter()
const { getUser, isAdmin } = useSession()
const { createArsip, uploadDokumenArsip, getKodeKlasifikasi, getLokasiSimpan, getUnitKerja } = useArsip()
const { showToast } = useToast()
const simpanLoading = ref(false)

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
  noItemArsip: '1',
  kodeKlasifikasi: '',
  jenisArsip: 'Konvensional',
  tingkatPerkembangan: 'Asli',
  uraianInformasiArsip: '',
})

const formB = reactive({
  kurunWaktu: String(new Date().getFullYear()),
  jumlah: '1 Lembar',
  keterangan: 'KONDISI BAIK',
  nomorDefinitif: '',
})

const formC = reactive({
  lokasiSimpan: 'Record Center',
  noBoks: '',
  noRak: '',
  noFolder: '',
})

const formD = reactive({
  jangkaSimpanInaktif: '3 Tahun',
  jangkaSimpanNasibAkhir: 'Musnah',
  skkad: 'Biasa',
  kategoriArsip: 'Inaktif',
})

const formE = reactive({
  unorId: '',
  penciptaArsip: '',
  petugasInput: '',
})

const formF = reactive<{ file: File | null; keterangan: string }>({
  file: null,
  keterangan: '',
})

onMounted(async () => {
  const user = getUser()
  if (user) {
    formE.petugasInput = user.name
  }

  const [kodes, lokasis, units] = await Promise.all([
    getKodeKlasifikasi(),
    getLokasiSimpan(),
    getUnitKerja(true),
  ])
  masterKodes.value = kodes
  masterLokasis.value = lokasis
  masterUnits.value = units

  if (units.length > 0 && !formE.unorId) {
    formE.unorId = units[0]!.id
    formE.penciptaArsip = units[0]!.namaUnit
  }
})

// ── Auto-fill JRA & SKKAD saat Kode Klasifikasi dipilih ─────────
watch(
  () => formA.kodeKlasifikasi,
  (selectedKode) => {
    const found = masterKodes.value.find((k) => k.kode === selectedKode)
    if (found) {
      formD.jangkaSimpanInaktif = found.jangkaSimpanInaktif || '3 Tahun'
      formD.jangkaSimpanNasibAkhir = found.nasibAkhir || 'Musnah'
      formD.skkad = found.skkad || 'Biasa'
    }
  }
)

// ── Auto-fill Pencipta saat Unit Kerja dipilih ──────────────────
watch(
  () => formE.unorId,
  (val) => {
    const found = masterUnits.value.find((u) => u.id === val)
    if (found) {
      formE.penciptaArsip = found.namaUnit
    }
  }
)

function buildPayload(status: string) {
  const selectedLokasi = masterLokasis.value.find((l) => l.namaLokasi === formC.lokasiSimpan)
  const selectedUnit = masterUnits.value.find((u) => u.id === formE.unorId)

  return {
    no_item_arsip: undefined,
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
    status,
  }
}

async function submitArsip(status: 'draft' | 'final') {
  simpanLoading.value = true
  try {
    const created = await createArsip(buildPayload(status))
    if (formF.file && created.id) {
      await uploadDokumenArsip(created.id, formF.file, formF.keterangan)
    }
    showToast(
      status === 'draft' ? 'Arsip berhasil disimpan sebagai draf' : 'Arsip berhasil disimpan',
      'success'
    )
    router.push(`/pengelolaan-arsip/${created.id}`)
  } catch (err: any) {
    showToast(err?.message || 'Gagal menyimpan arsip', 'danger')
  } finally {
    simpanLoading.value = false
  }
}

function simpanDraf() {
  submitArsip('draft')
}

function simpanFinal() {
  submitArsip('final')
}
</script>
