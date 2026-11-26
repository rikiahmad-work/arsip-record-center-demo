<template>
  <!-- ===== Modal Tambah Akun ===== -->
  <UiAppModal id="modalTambahAkun" title="Tambah Akun Baru" size="lg">
    <template #body>
      <div class="row g-2">
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.nama"
            label="Nama Lengkap"
            placeholder="Masukkan nama lengkap"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.nip"
            label="NIP"
            placeholder="Masukkan NIP"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.jabatan"
            label="Jabatan"
            placeholder="Masukkan jabatan"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.unor"
            label="Unit Organisasi (UNOR)"
            placeholder="Masukkan unit organisasi"
            required
          />
        </div>
        <div class="col-md-6">
          <FormSelect
            v-model="tambahForm.role"
            label="Role"
            :options="roleOptions"
            placeholder="Pilih role"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.email"
            type="email"
            label="Email"
            placeholder="Masukkan email"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.password"
            type="password"
            label="Password"
            placeholder="Masukkan password"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="tambahForm.konfirmasiPassword"
            type="password"
            label="Konfirmasi Password"
            placeholder="Ulangi password"
            required
          />
        </div>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Batal</button>
      <UiAppButton variant="primary" icon="ti ti-plus" @click="onTambah">
        Tambah Akun
      </UiAppButton>
    </template>
  </UiAppModal>

  <!-- ===== Modal Edit Akun ===== -->
  <UiAppModal id="modalEditAkun" title="Edit Akun" size="lg">
    <template #body>
      <div class="row g-2">
        <div class="col-md-6">
          <FormInput
            v-model="editForm.nama"
            label="Nama Lengkap"
            placeholder="Masukkan nama lengkap"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="editForm.nip"
            label="NIP"
            placeholder="Masukkan NIP"
            readonly
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="editForm.jabatan"
            label="Jabatan"
            placeholder="Masukkan jabatan"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="editForm.unor"
            label="Unit Organisasi (UNOR)"
            placeholder="Masukkan unit organisasi"
            required
          />
        </div>
        <div class="col-md-6">
          <FormSelect
            v-model="editForm.role"
            label="Role"
            :options="roleOptions"
            placeholder="Pilih role"
            required
          />
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="editForm.email"
            type="email"
            label="Email"
            placeholder="Masukkan email"
          />
        </div>
        <div class="col-12">
          <p class="text-muted small mb-0">
            <i class="ti ti-info-circle me-1"></i>
            Kosongkan password jika tidak ingin mengubah password.
          </p>
        </div>
        <div class="col-md-6">
          <FormInput
            v-model="editForm.passwordBaru"
            type="password"
            label="Password Baru (opsional)"
            placeholder="Biarkan kosong jika tidak berubah"
          />
        </div>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Batal</button>
      <UiAppButton variant="primary" icon="ti ti-device-floppy" @click="onEdit">
        Simpan Perubahan
      </UiAppButton>
    </template>
  </UiAppModal>

  <!-- ===== Modal Nonaktifkan ===== -->
  <UiAppModal id="modalNonaktif" title="Nonaktifkan Akun" size="sm">
    <template #body>
      <div class="text-center py-2">
        <i class="ti ti-user-off text-danger" style="font-size: 2.5rem;"></i>
        <p class="mt-3">
          Apakah Anda yakin ingin <strong>menonaktifkan</strong> akun
          <strong>{{ selectedRow?.nama }}</strong>?
        </p>
        <p class="text-muted small">Akun tidak akan bisa login setelah dinonaktifkan.</p>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Batal</button>
      <UiAppButton variant="danger" icon="ti ti-user-off" @click="onNonaktif">
        Nonaktifkan
      </UiAppButton>
    </template>
  </UiAppModal>

  <!-- ===== Modal Aktifkan ===== -->
  <UiAppModal id="modalAktif" title="Aktifkan Akun" size="sm">
    <template #body>
      <div class="text-center py-2">
        <i class="ti ti-user-check text-success" style="font-size: 2.5rem;"></i>
        <p class="mt-3">
          Apakah Anda yakin ingin <strong>mengaktifkan kembali</strong> akun
          <strong>{{ selectedRow?.nama }}</strong>?
        </p>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Batal</button>
      <UiAppButton variant="success" icon="ti ti-user-check" @click="onAktif">
        Aktifkan
      </UiAppButton>
    </template>
  </UiAppModal>

  <!-- ===== Modal Hapus ===== -->
  <UiAppModal id="modalHapusAkun" title="Hapus Akun" size="sm">
    <template #body>
      <div class="text-center py-2">
        <i class="ti ti-trash text-danger" style="font-size: 2.5rem;"></i>
        <p class="mt-3">
          Apakah Anda yakin ingin <strong>menghapus</strong> akun
          <strong>{{ selectedRow?.nama }}</strong>?
        </p>
        <p class="text-danger small">Tindakan ini tidak dapat dibatalkan.</p>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Batal</button>
      <UiAppButton variant="danger" icon="ti ti-trash" @click="onHapus">
        Hapus
      </UiAppButton>
    </template>
  </UiAppModal>

  <!-- ===== Modal Sukses ===== -->
  <UiAppModal id="modalSuksesAkun" title="Berhasil" size="sm">
    <template #body>
      <div class="text-center py-3">
        <i class="ti ti-circle-check text-success" style="font-size: 3rem;"></i>
        <p class="mt-3 mb-0">{{ suksesMessage }}</p>
      </div>
    </template>
    <template #footer>
      <button type="button" class="btn btn-success w-100" data-bs-dismiss="modal">OK</button>
    </template>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { UserRow } from '~/types/index'

interface Props {
  selectedRow?: UserRow | null
}

withDefaults(defineProps<Props>(), {
  selectedRow: null,
})

const emit = defineEmits<{
  tambah:   [data: typeof tambahForm]
  edit:     [data: typeof editForm]
  nonaktif: [row: UserRow]
  aktif:    [row: UserRow]
  hapus:    [row: UserRow]
}>()

// --- Form Tambah ---
const tambahForm = reactive({
  nama: '', nip: '', jabatan: '', unor: '',
  role: '', email: '', password: '', konfirmasiPassword: '',
})

function onTambah() {
  if (!tambahForm.nama || !tambahForm.nip || !tambahForm.password) {
    alert('Nama, NIP, dan Password wajib diisi.')
    return
  }
  if (tambahForm.password !== tambahForm.konfirmasiPassword) {
    alert('Password dan konfirmasi password tidak cocok.')
    return
  }
  emit('tambah', { ...tambahForm })
  Object.assign(tambahForm, { nama: '', nip: '', jabatan: '', unor: '', role: '', email: '', password: '', konfirmasiPassword: '' })
  closeModal('modalTambahAkun')
}

// --- Form Edit ---
const editForm = reactive({
  nama: '', nip: '', jabatan: '', unor: '',
  role: '', email: '', passwordBaru: '',
})

// Isi form edit dari selectedRow
const props = getCurrentInstance()?.props as any
watch(
  () => props?.selectedRow,
  (row: UserRow | null) => {
    if (row) {
      Object.assign(editForm, {
        nama: row.nama, nip: row.nip, jabatan: row.jabatan,
        unor: row.unitKerja, role: row.role, email: row.email ?? '',
        passwordBaru: '',
      })
    }
  },
  { immediate: true },
)

function onEdit() {
  emit('edit', { ...editForm })
  closeModal('modalEditAkun')
}

// --- Nonaktif ---
function onNonaktif() {
  const row = (getCurrentInstance()?.props as any)?.selectedRow
  if (row) emit('nonaktif', row)
  closeModal('modalNonaktif')
}

// --- Aktif ---
function onAktif() {
  const row = (getCurrentInstance()?.props as any)?.selectedRow
  if (row) emit('aktif', row)
  closeModal('modalAktif')
}

// --- Hapus ---
function onHapus() {
  const row = (getCurrentInstance()?.props as any)?.selectedRow
  if (row) emit('hapus', row)
  closeModal('modalHapusAkun')
}

// --- Options ---
const roleOptions = [
  { value: 'Super Admin', label: 'Super Admin' },
  { value: 'Admin Inspektorat', label: 'Admin Inspektorat' },
]

// --- Pesan Sukses ---
const suksesMessage = ref('Operasi berhasil dilakukan.')

// --- Helper Bootstrap Modal ---
function closeModal(id: string) {
  const win = window as any
  if (win.bootstrap) {
    const el = document.getElementById(id)
    if (el) win.bootstrap.Modal.getInstance(el)?.hide()
  }
}

function openModal(id: string) {
  const win = window as any
  if (win.bootstrap) {
    const el = document.getElementById(id)
    if (el) win.bootstrap.Modal.getOrCreateInstance(el).show()
  }
}

defineExpose({ openModal, suksesMessage, editForm })
</script>
