<template>
  <TemplatePageTitle title="Master Pengaturan" :breadcrumbs="[{ label: 'Master Pengaturan' }]">
    <template #action>
      <button class="btn btn-outline-secondary" type="button" :disabled="loading" @click="loadAll(true)">
        <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="ti ti-refresh me-1"></i>
        Muat Ulang
      </button>
    </template>
  </TemplatePageTitle>

  <div class="card">
    <div class="card-body">
      <!-- Nav Tabs -->
      <ul class="nav nav-tabs mb-4" id="pengaturanTabs" role="tablist">
        <li class="nav-item" role="presentation">
          <button class="nav-link" :class="{ active: activeTab === 'kode' }" type="button" @click="activeTab = 'kode'">
            <i class="ti ti-tag me-1"></i>Kode Klasifikasi
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button class="nav-link" :class="{ active: activeTab === 'lokasi' }" type="button" @click="activeTab = 'lokasi'">
            <i class="ti ti-map-pin me-1"></i>Lokasi Simpan
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button class="nav-link" :class="{ active: activeTab === 'unit' }" type="button" @click="activeTab = 'unit'">
            <i class="ti ti-building me-1"></i>Unit Kerja
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button class="nav-link" :class="{ active: activeTab === 'konfigurasi' }" type="button" @click="activeTab = 'konfigurasi'">
            <i class="ti ti-settings me-1"></i>Konfigurasi Sistem
          </button>
        </li>
      </ul>

      <!-- Tab Content -->
      <div class="tab-content">
        <div v-show="activeTab === 'kode'" class="tab-pane fade show active">
          <ModulesPengaturanTabKodeKlasifikasi
            :list="kodeKlasifikasiList"
            :saving="saving"
            @save="saveKodeKlasifikasi"
            @hapus="(item) => openConfirm('kode', item)"
          />
        </div>
        <div v-show="activeTab === 'lokasi'" class="tab-pane fade show active">
          <ModulesPengaturanTabLokasiSimpan
            :list="lokasiSimpanList"
            :saving="saving"
            @save="saveLokasiSimpan"
            @hapus="(item) => openConfirm('lokasi', item)"
          />
        </div>
        <div v-show="activeTab === 'unit'" class="tab-pane fade show active">
          <ModulesPengaturanTabUnitKerja
            :list="unitKerjaList"
            :saving="saving"
            @save="saveUnitKerja"
            @hapus="(item) => openConfirm('unit', item)"
          />
        </div>
        <div v-show="activeTab === 'konfigurasi'" class="tab-pane fade show active">
          <ModulesPengaturanTabKonfigurasiSistem
            :initial="konfigurasi"
            :saving="saving"
            @save="saveKonfigurasi"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Konfirmasi Hapus -->
  <UiAppModal v-if="confirmTarget" title="Hapus Data Master" size="sm" @close="confirmTarget = null">
    <p class="mb-0">
      Hapus <strong>{{ confirmLabel }}</strong>? Tindakan ini tidak dapat dibatalkan.
    </p>
    <template #footer>
      <button type="button" class="btn btn-outline-secondary" :disabled="saving" @click="confirmTarget = null">
        Batal
      </button>
      <button type="button" class="btn btn-danger" :disabled="saving" @click="confirmDelete">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        Hapus
      </button>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { KodeKlasifikasi, KonfigurasiSistem, LokasiSimpan, UnitKerja } from '~/types/arsip'
import { useArsip } from '~/composables/useArsip'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Master Pengaturan | Arsip Record Center KemenPPPA' })

const {
  getKodeKlasifikasi,
  createKodeKlasifikasi,
  updateKodeKlasifikasi,
  deleteKodeKlasifikasi,
  getLokasiSimpan,
  createLokasiSimpan,
  updateLokasiSimpan,
  deleteLokasiSimpan,
  getUnitKerja,
  createUnitKerja,
  updateUnitKerja,
  deleteUnitKerja,
  getKonfigurasi,
  updateKonfigurasi,
} = useArsip()

const { showToast } = useToast()

const activeTab = ref('kode')
const loading = ref(false)
const saving = ref(false)

const kodeKlasifikasiList = ref<KodeKlasifikasi[]>([])
const lokasiSimpanList = ref<LokasiSimpan[]>([])
const unitKerjaList = ref<UnitKerja[]>([])
const konfigurasi = ref<KonfigurasiSistem | null>(null)

type ConfirmType = 'kode' | 'lokasi' | 'unit'
const confirmTarget = ref<{ type: ConfirmType; item: KodeKlasifikasi | LokasiSimpan | UnitKerja } | null>(null)

const confirmLabel = computed(() => {
  if (!confirmTarget.value) return ''
  const { type, item } = confirmTarget.value
  if (type === 'kode') return (item as KodeKlasifikasi).nama
  if (type === 'lokasi') return (item as LokasiSimpan).namaLokasi
  return (item as UnitKerja).namaUnit
})

async function loadAll(forceRefresh = false) {
  loading.value = true
  try {
    const [kode, lokasi, unit, konf] = await Promise.all([
      getKodeKlasifikasi(forceRefresh),
      getLokasiSimpan(forceRefresh),
      getUnitKerja(forceRefresh),
      getKonfigurasi(),
    ])
    kodeKlasifikasiList.value = kode
    lokasiSimpanList.value = lokasi
    unitKerjaList.value = unit
    konfigurasi.value = konf
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal memuat data master', 'danger')
  } finally {
    loading.value = false
  }
}

async function saveKodeKlasifikasi(payload: Record<string, any>) {
  saving.value = true
  try {
    const { id, ...body } = payload
    if (id) {
      await updateKodeKlasifikasi(id, body)
      showToast('Kode klasifikasi berhasil diperbarui', 'success')
    } else {
      await createKodeKlasifikasi(body)
      showToast('Kode klasifikasi berhasil ditambahkan', 'success')
    }
    kodeKlasifikasiList.value = await getKodeKlasifikasi(true)
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal menyimpan kode klasifikasi', 'danger')
  } finally {
    saving.value = false
  }
}

async function saveLokasiSimpan(payload: Record<string, any>) {
  saving.value = true
  try {
    const { id, ...body } = payload
    if (id) {
      await updateLokasiSimpan(id, body)
      showToast('Lokasi simpan berhasil diperbarui', 'success')
    } else {
      await createLokasiSimpan(body)
      showToast('Lokasi simpan berhasil ditambahkan', 'success')
    }
    lokasiSimpanList.value = await getLokasiSimpan(true)
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal menyimpan lokasi simpan', 'danger')
  } finally {
    saving.value = false
  }
}

async function saveUnitKerja(payload: Record<string, any>) {
  saving.value = true
  try {
    const { id, ...body } = payload
    if (id) {
      await updateUnitKerja(id, body)
      showToast('Unit kerja berhasil diperbarui', 'success')
    } else {
      await createUnitKerja(body)
      showToast('Unit kerja berhasil ditambahkan', 'success')
    }
    unitKerjaList.value = await getUnitKerja(true)
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal menyimpan unit kerja', 'danger')
  } finally {
    saving.value = false
  }
}

async function saveKonfigurasi(payload: KonfigurasiSistem) {
  saving.value = true
  try {
    konfigurasi.value = await updateKonfigurasi(payload)
    showToast('Konfigurasi sistem berhasil disimpan', 'success')
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal menyimpan konfigurasi', 'danger')
  } finally {
    saving.value = false
  }
}

function openConfirm(type: ConfirmType, item: KodeKlasifikasi | LokasiSimpan | UnitKerja) {
  confirmTarget.value = { type, item }
}

async function confirmDelete() {
  if (!confirmTarget.value) return
  const { type, item } = confirmTarget.value
  saving.value = true
  try {
    if (type === 'kode') {
      await deleteKodeKlasifikasi(item.id)
      kodeKlasifikasiList.value = await getKodeKlasifikasi(true)
      showToast('Kode klasifikasi berhasil dihapus', 'success')
    } else if (type === 'lokasi') {
      await deleteLokasiSimpan(item.id)
      lokasiSimpanList.value = await getLokasiSimpan(true)
      showToast('Lokasi simpan berhasil dihapus', 'success')
    } else {
      await deleteUnitKerja(item.id)
      unitKerjaList.value = await getUnitKerja(true)
      showToast('Unit kerja berhasil dihapus', 'success')
    }
    confirmTarget.value = null
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Gagal menghapus data', 'danger')
  } finally {
    saving.value = false
  }
}

onMounted(() => loadAll())
</script>
