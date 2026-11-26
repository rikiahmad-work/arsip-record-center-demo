<template>
  <TemplatePageTitle title="Cetak Label Arsip (Boks, Folder, & Arsip Inaktif)" :breadcrumbs="[{ label: 'Cetak Label' }]">
    <template #action>
      <button
        type="button"
        class="btn btn-primary d-inline-flex align-items-center gap-1"
        :disabled="printing"
        @click="handleCetakPDF"
      >
        <span v-if="printing" class="spinner-border spinner-border-sm"></span>
        <i v-else class="ti ti-printer"></i>
        <span>Generate &amp; Cetak PDF</span>
      </button>
    </template>
  </TemplatePageTitle>

  <div class="row g-4">
    <!-- Kolom Kiri: Pilihan Jenis Label & Parameter -->
    <div class="col-xl-6">
      <div class="card h-100 mb-0">
        <div class="card-header">
          <h5 class="mb-0"><i class="ti ti-adjustments-horizontal me-2"></i>Pengaturan &amp; Parameter Label</h5>
          <small class="text-muted">Pilih jenis label fisik atau ambil otomatis dari berkas arsip terdaftar</small>
        </div>

        <div class="card-body">
          <!-- Pilihan 3 Jenis Label -->
          <label class="form-label fw-semibold mb-2">1. Pilih Jenis Format Label</label>
          <div class="row g-2 mb-4">
            <div v-for="t in labelTypes" :key="t.key" class="col-md-4">
              <div
                class="p-3 rounded border cursor-pointer text-center label-type-card"
                :class="selectedType === t.key ? 'border-primary bg-light-primary' : 'bg-white'"
                @click="selectedType = t.key"
              >
                <i :class="[t.icon, selectedType === t.key ? 'text-primary' : 'text-secondary']" style="font-size: 1.75rem;"></i>
                <div class="fw-bold mt-1" :class="selectedType === t.key ? 'text-primary' : 'text-dark'">
                  {{ t.title }}
                </div>
                <small class="text-muted d-block" style="font-size: 0.75rem;">{{ t.subtitle }}</small>
              </div>
            </div>
          </div>

          <!-- Ambil Otomatis dari Data Arsip -->
          <div class="p-3 bg-light rounded border mb-3">
            <FormSelect
              v-model="selectedArsipId"
              label="2. Isi Otomatis dari Berkas Arsip Terdaftar (Opsional)"
              :options="arsipOptions"
              placeholder="-- Cari Berkas Arsip untuk Auto-Fill Label --"
              search-placeholder="Ketik No. Item, Kode Klasifikasi, No. Boks, atau Uraian..."
              searchable
              clearable
              no-margin
              @update:model-value="onSelectArsip"
            />
          </div>

          <!-- Form Parameter Label -->
          <div class="row g-3">
            <div class="col-md-8">
              <FormSelect
                v-model="form.unitKerja"
                label="Unit Kerja / Pencipta Arsip"
                :options="unitOptions"
                placeholder="-- Pilih atau Ketik Unit Kerja --"
                search-placeholder="Cari unit kerja..."
                searchable
                no-margin
              />
            </div>

            <div class="col-md-4">
              <label class="form-label">Tahun / Kurun Waktu</label>
              <input v-model="form.tahun" type="text" class="form-control" placeholder="Contoh: 2024 - 2026">
            </div>

            <div class="col-md-6">
              <FormSelect
                v-model="form.kode"
                label="Kode Klasifikasi"
                :options="kodeOptions"
                placeholder="-- Pilih Kode Klasifikasi --"
                search-placeholder="Cari dari 435 kode klasifikasi..."
                searchable
                clearable
                no-margin
              />
            </div>

            <div class="col-md-6">
              <label class="form-label">Rentang Nomor Item Arsip</label>
              <input v-model="form.rentangItem" type="text" class="form-control" placeholder="Contoh: 1 s/d 25">
            </div>

            <div class="col-md-4">
              <label class="form-label">Nomor Rak</label>
              <input v-model="form.noRak" type="text" class="form-control" placeholder="Contoh: R-01">
            </div>

            <div class="col-md-4">
              <label class="form-label">Nomor Boks</label>
              <input v-model="form.noBoks" type="text" class="form-control" placeholder="Contoh: B-001">
            </div>

            <div class="col-md-4">
              <label class="form-label">Nomor Folder</label>
              <input v-model="form.noFolder" type="text" class="form-control" placeholder="Contoh: F-01">
            </div>

            <div class="col-12">
              <label class="form-label">Uraian Ringkas Isi Boks / Folder</label>
              <textarea
                v-model="form.uraian"
                class="form-control"
                rows="2"
                placeholder="Masukkan uraian ringkas berkas arsip untuk dicetak pada label"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="card-footer d-flex justify-content-between align-items-center">
          <button type="button" class="btn btn-outline-secondary" @click="resetForm">
            <i class="ti ti-refresh me-1"></i>Reset Parameter
          </button>
          <div class="d-flex gap-2">
            <button type="button" class="btn btn-outline-primary" :disabled="printing" @click="handleDownloadPDF">
              <i class="ti ti-download me-1"></i>Unduh File PDF
            </button>
            <button type="button" class="btn btn-primary" :disabled="printing" @click="handleCetakPDF">
              <span v-if="printing" class="spinner-border spinner-border-sm me-1"></span>
              <i v-else class="ti ti-printer me-1"></i>Buka &amp; Cetak PDF
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Kolom Kanan: Live Preview Visual Label -->
    <div class="col-xl-6">
      <div class="card h-100 mb-0">
        <div class="card-header d-flex justify-content-between align-items-center">
          <div>
            <h5 class="mb-0"><i class="ti ti-eye me-2"></i>Pratinjau Desain Label (Live Preview)</h5>
            <small class="text-muted">Simulasi tampilan fisik sebelum dicetak ke PDF standar Record Center</small>
          </div>
          <span class="badge bg-light-primary text-primary">{{ activeTypeTitle }}</span>
        </div>

        <div class="card-body d-flex flex-column justify-content-center align-items-center bg-light p-4">
          <!-- Kartu Preview Label Standar KemenPPPA -->
          <div class="label-preview-sheet bg-white shadow-sm w-100" style="max-width: 520px;">
            <!-- Header Instansi -->
            <div class="label-header d-flex align-items-center gap-3 p-3 border-bottom border-dark border-2">
              <img src="/admin/images/logokpppa.png" alt="Logo KemenPPPA" style="width: 48px; height: 48px; object-fit: contain;">
              <div class="text-center flex-grow-1">
                <div class="fw-bold text-uppercase" style="font-size: 0.85rem; letter-spacing: 0.4px;">
                  Kementerian Pemberdayaan Perempuan dan Perlindungan Anak
                </div>
                <div class="fw-bold text-primary text-uppercase" style="font-size: 0.95rem;">
                  RECORD CENTER — {{ activeTypeTitle }}
                </div>
              </div>
            </div>

            <!-- Body Tabel Label -->
            <div class="p-3">
              <table class="table table-bordered border-dark mb-0 align-middle">
                <tbody>
                  <tr>
                    <td class="bg-light fw-bold small" style="width: 38%;">UNIT KERJA / PENCIPTA</td>
                    <td class="fw-semibold">{{ form.unitKerja || '-' }}</td>
                  </tr>
                  <tr>
                    <td class="bg-light fw-bold small">KODE KLASIFIKASI</td>
                    <td>
                      <span class="badge bg-primary fs-6">{{ form.kode || '-' }}</span>
                    </td>
                  </tr>
                  <tr>
                    <td class="bg-light fw-bold small">KURUN WAKTU / TAHUN</td>
                    <td class="fw-semibold">{{ form.tahun || '-' }}</td>
                  </tr>
                  <tr v-if="selectedType === 'box' || selectedType === 'inaktif'">
                    <td class="bg-light fw-bold small">RENTANG NO. ITEM</td>
                    <td class="fw-semibold">{{ form.rentangItem || '-' }}</td>
                  </tr>
                  <tr>
                    <td class="bg-light fw-bold small">URAIAN INFORMASI</td>
                    <td class="small">{{ form.uraian || '-' }}</td>
                  </tr>
                </tbody>
              </table>

              <!-- Kotak Nomor Rak / Boks / Folder Besar -->
              <div class="row g-2 mt-2 text-center">
                <div class="col-4">
                  <div class="border border-dark rounded p-2">
                    <div class="small text-muted fw-bold">NO. RAK</div>
                    <div class="fs-4 fw-bold text-dark">{{ form.noRak || '-' }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="border border-dark rounded p-2 bg-light-primary">
                    <div class="small text-primary fw-bold">NO. BOKS</div>
                    <div class="fs-4 fw-bold text-primary">{{ form.noBoks || '-' }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="border border-dark rounded p-2">
                    <div class="small text-muted fw-bold">NO. FOLDER</div>
                    <div class="fs-4 fw-bold text-dark">{{ form.noFolder || '-' }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <p class="text-muted small mt-3 mb-0 text-center">
            <i class="ti ti-info-circle me-1"></i>
            Klik tombol <strong>Buka &amp; Cetak PDF</strong> untuk menghasilkan dokumen PDF vektor resmi dari server (`/api/label/{{ selectedType }}`).
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Arsip, KodeKlasifikasi, UnitKerja } from '~/types/arsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Cetak Label Arsip | Arsip Record Center KemenPPPA' })

const { apiFetch } = useSession()
const { getArsipList, getKodeKlasifikasi, getUnitKerja } = useArsip()
const { showToast } = useToast()

type LabelTypeKey = 'box' | 'folder' | 'inaktif'

const labelTypes: { key: LabelTypeKey; title: string; subtitle: string; icon: string }[] = [
  {
    key: 'box',
    title: 'Label Boks Arsip',
    subtitle: 'Punggung Boks Record Center',
    icon: 'ti ti-box',
  },
  {
    key: 'folder',
    title: 'Label Folder',
    subtitle: 'Tab Map / Sampul Berkas',
    icon: 'ti ti-folder',
  },
  {
    key: 'inaktif',
    title: 'Label Arsip Inaktif',
    subtitle: 'Stiker Identitas Berkas DAI',
    icon: 'ti ti-tag',
  },
]

const selectedType = ref<LabelTypeKey>('box')
const selectedArsipId = ref('')
const printing = ref(false)

const arsipList = ref<Arsip[]>([])
const kodeList = ref<KodeKlasifikasi[]>([])
const unitList = ref<UnitKerja[]>([])

const form = reactive({
  unitKerja: 'Biro Data dan Informasi',
  kode: 'HK.00.00',
  tahun: '2026',
  noRak: 'R-01',
  noBoks: 'B-001',
  noFolder: 'F-01',
  rentangItem: '1 s/d 10',
  uraian: 'Berkas Arsip Dinamis Inaktif Record Center KemenPPPA',
})

const activeTypeTitle = computed(() =>
  labelTypes.find((t) => t.key === selectedType.value)?.title.toUpperCase() || 'LABEL ARSIP'
)

const arsipOptions = computed(() =>
  arsipList.value.map((a) => ({
    value: a.id,
    label: `${a.noItemArsip} [${a.kodeKlasifikasi}] — Boks ${a.noBoks || '-'} / Folder ${a.noFolder || '-'} — ${(a.uraianInformasi || a.uraianInformasiArsip).slice(0, 60)}`,
  }))
)

const kodeOptions = computed(() =>
  kodeList.value.map((k) => ({
    value: k.kode,
    label: `${k.kode} — ${k.nama}`,
  }))
)

const unitOptions = computed(() =>
  unitList.value.map((u) => ({
    value: u.namaUnit,
    label: `${u.kodeUnit} — ${u.namaUnit}`,
  }))
)

onMounted(async () => {
  try {
    const [arsipRes, kodes, units] = await Promise.all([
      getArsipList({ per_page: 500 }),
      getKodeKlasifikasi(),
      getUnitKerja(),
    ])
    arsipList.value = arsipRes.data
    kodeList.value = kodes
    unitList.value = units

    if (arsipRes.data.length > 0) {
      const first = arsipRes.data[0]
      if (first) applyArsipToForm(first)
    }
  } catch {
    // Fallback default values remain
  }
})

function applyArsipToForm(a: Arsip) {
  form.unitKerja = a.penciptaArsip || 'Biro Data dan Informasi'
  form.kode = a.kodeKlasifikasi || 'HK.00.00'
  form.tahun = a.kurunWaktu || '2026'
  form.noRak = a.noRak || 'R-01'
  form.noBoks = a.noBoks || 'B-001'
  form.noFolder = a.noFolder || 'F-01'
  form.rentangItem = a.noItemArsip || '1'
  form.uraian = a.uraianInformasi || a.uraianInformasiArsip || ''
}

function onSelectArsip(id: string) {
  if (!id) return
  const found = arsipList.value.find((a) => String(a.id) === String(id))
  if (found) {
    applyArsipToForm(found)
    showToast(`Parameter label diisi otomatis dari arsip ${found.noItemArsip}`, 'info')
  }
}

function resetForm() {
  selectedArsipId.value = ''
  form.unitKerja = 'Biro Data dan Informasi'
  form.kode = ''
  form.tahun = String(new Date().getFullYear())
  form.noRak = 'R-01'
  form.noBoks = 'B-001'
  form.noFolder = 'F-01'
  form.rentangItem = '1 s/d 10'
  form.uraian = ''
}

function buildQueryParams() {
  const qs = new URLSearchParams()
  if (form.unitKerja) qs.set('unit_kerja', form.unitKerja)
  if (form.kode) qs.set('kode', form.kode)
  if (form.tahun) qs.set('tahun', form.tahun)
  if (form.noRak) qs.set('no_rak', form.noRak)
  if (form.noBoks) qs.set('no_boks', form.noBoks)
  if (form.noFolder) qs.set('no_folder', form.noFolder)
  if (form.rentangItem) qs.set('rentang_item', form.rentangItem)
  if (form.uraian) qs.set('uraian', form.uraian)
  return qs.toString()
}

async function handleCetakPDF() {
  printing.value = true
  try {
    const res = await apiFetch(`/label/${selectedType.value}?${buildQueryParams()}`)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    window.open(url, '_blank')
    showToast(`PDF ${activeTypeTitle.value} berhasil dibuat`, 'success')
  } catch {
    showToast('Gagal membuat PDF label arsip', 'danger')
  } finally {
    printing.value = false
  }
}

async function handleDownloadPDF() {
  printing.value = true
  try {
    const res = await apiFetch(`/label/${selectedType.value}?${buildQueryParams()}`)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `label-${selectedType.value}-${form.noBoks || 'arsip'}.pdf`
    link.click()
    URL.revokeObjectURL(url)
    showToast(`File PDF ${activeTypeTitle.value} berhasil diunduh`, 'success')
  } catch {
    showToast('Gagal mengunduh PDF label arsip', 'danger')
  } finally {
    printing.value = false
  }
}
</script>

<style scoped>
.label-type-card {
  transition: all 0.15s ease;
}
.label-type-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}
.label-preview-sheet {
  border: 2px solid #262626;
  border-radius: 6px;
}
</style>
