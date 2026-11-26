<template>
  <TemplatePageTitle :title="`Detail Arsip — ${arsip?.noItemArsip || ''}`" :breadcrumbs="[
    { label: 'Pengelolaan Arsip', to: '/pengelolaan-arsip' },
    { label: 'Detail Arsip' },
  ]">
    <template #action>
      <div class="d-flex gap-2">
        <button class="btn btn-outline-secondary" @click="$router.push('/pengelolaan-arsip')">
          <i class="ti ti-arrow-left me-1"></i>Kembali
        </button>
        <button
          v-if="canWrite()"
          class="btn btn-primary"
          @click="$router.push(`/pengelolaan-arsip/${id}/edit`)"
        >
          <i class="ti ti-pencil me-1"></i>Edit Arsip
        </button>
      </div>
    </template>
  </TemplatePageTitle>

  <div v-if="loading" class="text-center py-5">
    <div class="spinner-border text-primary" role="status">
      <span class="visually-hidden">Memuat...</span>
    </div>
    <p class="mt-2 text-muted">Memuat detail arsip...</p>
  </div>

  <div v-else-if="!arsip" class="text-center py-5 text-muted">
    <i class="ti ti-inbox" style="font-size: 48px;"></i>
    <p class="mt-3">Arsip tidak ditemukan</p>
  </div>

  <!-- Satu Card Utama Utuh: 2 Kolom Kiri & Kanan Tanpa Dipisah Per Sub-Judul -->
  <div v-else class="card shadow-sm border">
    <div class="card-header bg-white py-3 d-flex justify-content-between align-items-center">
      <h5 class="mb-0 text-dark">
        <i class="ti ti-file-text me-2 text-primary"></i>Detail Data Berkas Arsip Record Center
      </h5>
      <div class="d-flex align-items-center gap-2">
        <span :class="['badge', arsip.status === 'final' ? 'bg-success' : 'bg-warning']">
          {{ arsip.status === 'final' ? 'Final' : 'Draft' }}
        </span>
        <span class="badge" :class="kategoriBadge(arsip.kategoriArsip)">
          {{ arsip.kategoriArsip || 'Inaktif' }}
        </span>
      </div>
    </div>

    <div class="card-body">
      <!-- Dua Kolom Kiri & Kanan Berdampingan -->
      <div class="row g-4">
        <!-- Kolom Kiri: Informasi Identitas & Fisik Arsip -->
        <div class="col-lg-6">
          <table class="table table-sm table-borderless align-middle mb-0">
            <tbody>
              <tr>
                <td class="text-muted" style="width: 175px; vertical-align: top;">No. Item Arsip</td>
                <td class="fw-bold font-monospace text-dark">{{ arsip.noItemArsip || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted" style="vertical-align: top;">Kode Klasifikasi</td>
                <td>
                  <span class="badge bg-light-primary text-primary fw-bold font-monospace px-2 py-1">
                    {{ arsip.kodeKlasifikasi || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted" style="vertical-align: top;">Uraian Informasi</td>
                <td class="text-dark" style="line-height: 1.5;">
                  {{ arsip.uraianInformasi || arsip.uraianInformasiArsip || '-' }}
                </td>
              </tr>
              <tr>
                <td class="text-muted">Jenis Arsip</td>
                <td>
                  <span class="badge bg-light-secondary text-dark border">
                    {{ arsip.jenisArsip || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Tingkat Perkembangan</td>
                <td>
                  <span class="badge bg-light-success text-success border border-success-subtle">
                    {{ arsip.tingkatPerkembangan || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Kurun Waktu</td>
                <td class="font-monospace text-dark fw-medium">{{ arsip.kurunWaktu || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Jumlah Berkas</td>
                <td class="text-dark fw-medium">{{ arsip.jumlah || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Kondisi / Keterangan</td>
                <td class="text-dark">{{ arsip.keterangan || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Nomor Definitif</td>
                <td class="font-monospace text-dark fw-medium">{{ arsip.nomorDefinitif || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Unit Pengolah / Pencipta</td>
                <td class="text-dark fw-medium">{{ arsip.penciptaArsip || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Kolom Kanan: Lokasi Simpan, Retensi JRA, & Keamanan -->
        <div class="col-lg-6">
          <table class="table table-sm table-borderless align-middle mb-0">
            <tbody>
              <tr>
                <td class="text-muted" style="width: 175px;">Lokasi Simpan</td>
                <td class="text-dark fw-medium">{{ arsip.lokasiSimpan || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Nomor Rak</td>
                <td>
                  <span class="badge bg-light text-dark border font-monospace px-2 py-1">
                    {{ arsip.noRak || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Nomor Boks</td>
                <td>
                  <span class="badge bg-light text-dark border font-monospace px-2 py-1">
                    {{ arsip.noBoks || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Nomor Folder</td>
                <td>
                  <span class="badge bg-light text-dark border font-monospace px-2 py-1">
                    {{ arsip.noFolder || '-' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Jangka Simpan Inaktif</td>
                <td class="text-dark fw-medium">{{ arsip.jangkaSimpanInaktif || '-' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Nasib Akhir (JRA)</td>
                <td>
                  <span class="badge bg-primary px-2 py-1">
                    {{ arsip.nasibAkhir || arsip.jangkaSimpanNasibAkhir || 'Musnah' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Batas Retensi (31 Des)</td>
                <td>
                  <span class="text-dark fw-medium me-2">{{ arsip.tanggalJatuhTempo || '-' }}</span>
                  <span
                    v-if="arsip.isJatuhTempo"
                    class="badge bg-light-danger text-danger border border-danger px-2 py-1"
                  >
                    <i class="ti ti-alert-triangle me-1"></i>{{ arsip.rekomendasiRetensi || 'Jatuh Tempo' }}
                  </span>
                </td>
              </tr>
              <tr v-if="arsip.sedangDipinjam || arsip.dalamPenyusutan">
                <td class="text-muted">Status Operasional</td>
                <td>
                  <div class="d-flex flex-wrap gap-1">
                    <span v-if="arsip.sedangDipinjam" class="badge bg-warning text-dark px-2 py-1">
                      <i class="ti ti-book-upload me-1"></i>Sedang Dipinjam ({{ arsip.statusPinjam }})
                    </span>
                    <span v-if="arsip.dalamPenyusutan" class="badge bg-dark text-white px-2 py-1">
                      <i class="ti ti-recycle me-1"></i>{{ arsip.statusPenyusutan || 'Proses Penyusutan' }}
                    </span>
                  </div>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Klasifikasi Keamanan (SKKAD)</td>
                <td>
                  <span class="badge px-2 py-1" :class="skkadBadge(arsip.skkad)">
                    {{ arsip.skkad || 'Biasa' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Kategori Berkas</td>
                <td>
                  <span class="badge px-2 py-1" :class="kategoriBadge(arsip.kategoriArsip)">
                    {{ arsip.kategoriArsip || 'Inaktif' }}
                  </span>
                </td>
              </tr>
              <tr>
                <td class="text-muted">Petugas Input</td>
                <td class="text-dark">{{ arsip.petugasInput || 'Administrator Sistem' }}</td>
              </tr>
              <tr>
                <td class="text-muted">Terdaftar Pada</td>
                <td class="text-dark small">{{ formatDateTimeIndo(arsip.createdAt) }}</td>
              </tr>
              <tr v-if="arsip.updatedAt">
                <td class="text-muted">Terakhir Diperbarui</td>
                <td class="text-dark small">{{ formatDateTimeIndo(arsip.updatedAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Garis Pemisah Halus ke Dokumen Terlampir -->
      <hr class="my-4">

      <!-- Dokumen Pendukung & Digitalisasi Arsip (Masih di dalam 1 Card) -->
      <div class="d-flex justify-content-between align-items-center mb-3">
        <div>
          <h6 class="mb-0 text-dark fw-bold">
            <i class="ti ti-paperclip me-1 text-primary"></i>Berkas Digital &amp; Dokumen Pendukung
          </h6>
          <small class="text-muted">Daftar file hasil alih media / scan dokumen fisik (PDF, JPG, PNG)</small>
        </div>
        <span class="badge bg-light-primary text-primary border border-primary-subtle">
          {{ dokumenList.length }} File Terlampir
        </span>
      </div>

      <div v-if="dokumenList.length === 0" class="text-center text-muted py-4 bg-light rounded border">
        <i class="ti ti-file-off fs-2 d-block mb-1 text-secondary opacity-50"></i>
        <p class="mb-0 small">Belum ada dokumen digitalisasi fisik yang dilampirkan untuk berkas arsip ini.</p>
      </div>

      <div v-else class="table-responsive border rounded">
        <table class="table table-sm table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th class="ps-3">Nama Berkas</th>
              <th>Format File</th>
              <th>Ukuran File</th>
              <th>Keterangan</th>
              <th class="text-center" style="width: 170px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="doc in dokumenList" :key="doc.id">
              <td class="ps-3 fw-medium text-dark">
                <i class="ti ti-file-text me-2 text-primary"></i>{{ doc.namaFile }}
              </td>
              <td><span class="badge bg-light text-dark border font-monospace">{{ doc.mimeType }}</span></td>
              <td class="font-monospace small">{{ formatBytes(doc.ukuranBytes) }}</td>
              <td>{{ doc.keterangan || '-' }}</td>
              <td class="text-center">
                <div class="d-inline-flex gap-1">
                  <button
                    type="button"
                    class="btn btn-sm btn-light-info"
                    title="Buka / Preview Berkas"
                    @click="handleDownload(doc, 'inline')"
                  >
                    <i class="ti ti-eye me-1"></i>Buka
                  </button>
                  <button
                    v-if="canWrite()"
                    type="button"
                    class="btn btn-sm btn-light-danger"
                    title="Hapus Dokumen"
                    @click="handleDeleteDokumen(doc)"
                  >
                    <i class="ti ti-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Arsip, DokumenArsipItem } from '~/types/arsip'
import { useToast } from '~/composables/useToast'
import { useSession } from '~/composables/useSession'
import { useArsip } from '~/composables/useArsip'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const route = useRoute()
const { canWrite, apiFetch } = useSession()
const { getArsipById, listDokumenArsip, deleteDokumenArsip } = useArsip()
const { showToast } = useToast()

const id = route.params.id as string
const arsip = ref<Arsip | null>(null)
const dokumenList = ref<DokumenArsipItem[]>([])
const loading = ref(true)

async function fetchDetail() {
  loading.value = true
  try {
    const [data, docs] = await Promise.all([
      getArsipById(id),
      listDokumenArsip(id),
    ])
    arsip.value = data
    dokumenList.value = docs
  } catch {
    arsip.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchDetail())

useHead({
  title: computed(() => `Detail Arsip — ${arsip.value?.noItemArsip || ''} | Arsip Record Center KemenPPPA`),
})

async function handleDownload(doc: DokumenArsipItem, disposition: 'inline' | 'attachment') {
  try {
    const res = await apiFetch(`/dokumen/${doc.id}/download?disposition=${disposition}`)
    if (!res.ok) {
      showToast('Anda tidak memiliki hak akses untuk mengunduh dokumen ini', 'danger')
      return
    }
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    window.open(url, '_blank')
  } catch {
    showToast('Gagal membuka dokumen', 'danger')
  }
}

async function handleDeleteDokumen(doc: DokumenArsipItem) {
  try {
    const ok = await deleteDokumenArsip(doc.id)
    if (ok) {
      showToast('Dokumen berhasil dihapus', 'success')
      dokumenList.value = await listDokumenArsip(id)
    }
  } catch (err: any) {
    showToast(err?.message || 'Gagal menghapus dokumen', 'danger')
  }
}

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

function skkadBadge(skkad?: string): string {
  const upper = (skkad || '').toUpperCase()
  if (upper === 'RAHASIA') return 'bg-danger text-white'
  if (upper === 'TERBATAS') return 'bg-warning text-dark'
  return 'bg-success text-white'
}

function kategoriBadge(kategori?: string): string {
  const upper = (kategori || '').toUpperCase()
  if (upper === 'VITAL') return 'bg-warning text-dark'
  if (upper === 'TERJAGA') return 'bg-info text-white'
  if (upper === 'USUL MUSNAH') return 'bg-danger text-white'
  if (upper === 'USUL SERAH') return 'bg-primary text-white'
  return 'bg-secondary text-white'
}
</script>
