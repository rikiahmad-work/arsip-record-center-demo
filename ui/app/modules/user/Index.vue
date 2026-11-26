<template>
  <TemplatePageTitle title="Manajemen User" :breadcrumbs="[{ label: 'Manajemen User' }]">
    <template #action>
      <button class="btn btn-primary" @click="openTambah">
        <i class="ti ti-plus me-1"></i>Tambah User
      </button>
    </template>
  </TemplatePageTitle>

  <!-- 4 Stat Cards -->
  <div class="row mb-3">
    <div v-for="card in statCards" :key="card.label" class="col-md-6 col-xl-3">
      <div class="card border-0 shadow-sm">
        <div class="card-body d-flex align-items-center justify-content-between">
          <div>
            <div class="text-muted small">{{ card.label }}</div>
            <h4 class="mb-0 mt-1">{{ card.value }}</h4>
          </div>
          <div class="stat-icon" :class="card.iconBg">
            <i :class="card.icon"></i>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="card">
    <div class="card-body">
      <!-- Nav Tabs -->
      <ul class="nav nav-tabs mb-4" role="tablist">
        <li class="nav-item" role="presentation">
          <button
            class="nav-link"
            :class="{ active: activeTab === 'users' }"
            type="button"
            @click="activeTab = 'users'"
          >
            <i class="ti ti-users me-1"></i>Daftar Pengguna
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button
            class="nav-link"
            :class="{ active: activeTab === 'rbac' }"
            type="button"
            @click="activeTab = 'rbac'"
          >
            <i class="ti ti-shield-lock me-1"></i>Pengaturan Hak Akses
          </button>
        </li>
      </ul>

      <div class="tab-content">
        <!-- Tab 1: Daftar Pengguna -->
        <div v-show="activeTab === 'users'" class="tab-pane fade show active">
          <!-- Bar Filter Dedicated di Dalam Card Tabel (Seragam dengan Pengelolaan Arsip) -->
          <div class="p-3 bg-light rounded border mb-3">
            <div class="row g-2 align-items-end">
              <div class="col-12 col-md-3">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-search me-1"></i>Pencarian
                </label>
                <input
                  v-model="filter.search"
                  type="text"
                  name="user-filter-search"
                  autocomplete="off"
                  class="form-control form-control-sm"
                  placeholder="Cari nama / username / NIP / email / UNOR"
                >
              </div>
              <div class="col-6 col-md-3">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-shield me-1"></i>Role
                </label>
                <select v-model="filter.role" class="form-select form-select-sm">
                  <option value="">Semua Role</option>
                  <option v-for="role in roleOptions" :key="role.key" :value="role.key">
                    {{ role.name }}
                  </option>
                </select>
              </div>
              <div class="col-6 col-md-2">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-toggle-left me-1"></i>Status
                </label>
                <select v-model="filter.status" class="form-select form-select-sm">
                  <option value="">Semua Status</option>
                  <option value="active">Aktif</option>
                  <option value="inactive">Nonaktif</option>
                </select>
              </div>
              <div class="col-6 col-md-2">
                <label class="form-label small text-muted fw-semibold mb-1">
                  <i class="ti ti-trash me-1"></i>Data Terhapus
                </label>
                <select v-model="filter.trashed" class="form-select form-select-sm" @change="fetchUsers(true)">
                  <option value="">Tidak</option>
                  <option value="with">Ikut Sertakan</option>
                  <option value="only">Hanya Terhapus</option>
                </select>
              </div>
              <div class="col-6 col-md-2 d-flex gap-1">
                <button
                  type="button"
                  class="btn btn-sm w-100 d-inline-flex align-items-center justify-content-center gap-1"
                  :class="hasActiveFilter ? 'btn-outline-danger' : 'btn-outline-secondary opacity-50'"
                  :disabled="!hasActiveFilter"
                  @click="resetFilter"
                >
                  <i class="ti ti-filter-off"></i>
                  <span>Reset Filter</span>
                </button>
              </div>
            </div>
          </div>

          <div v-if="initialLoading" class="text-center py-4">
            <div class="spinner-border text-primary" role="status"></div>
            <div class="text-muted small mt-2">Memuat data pengguna...</div>
          </div>
          <ModulesUserGrid
            v-else
            :rows="filteredUserList"
            :current-user-id="currentUser?.id"
            @edit="openEdit"
            @reset-password="openResetPassword"
            @reset-mfa="openResetMfa"
            @nonaktif="onNonaktif"
            @aktif="onAktif"
            @hapus="konfirmasiHapus"
            @restore="konfirmasiRestore"
          />
        </div>

        <!-- Tab 2: Pengaturan Hak Akses -->
        <div v-show="activeTab === 'rbac'" class="tab-pane fade show active">
          <div v-if="rolesLoading" class="text-center py-4">
            <div class="spinner-border text-primary" role="status"></div>
            <div class="text-muted small mt-2">Memuat data peran...</div>
          </div>
          <ModulesUserRbacMatrix
            v-else
            :roles="roles"
            :permission-groups="permissionGroups"
            :saving="savingRole"
            @save="savePermissions"
            @add-role="addRole"
            @delete-role="deleteRole"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Tambah/Edit User -->
  <UiAppModal
    v-if="showForm"
    :title="formMode === 'tambah' ? 'Tambah Akun Baru' : 'Edit Akun'"
    size="lg"
    @close="showForm = false"
  >
    <div v-if="formErrors.length" class="alert alert-danger py-2">
      <ul class="mb-0 ps-3">
        <li v-for="(err, i) in formErrors" :key="i">{{ err }}</li>
      </ul>
    </div>

    <div class="row g-3">
      <div class="col-md-6">
        <label class="form-label">Nama Lengkap <span class="text-danger">*</span></label>
        <input v-model="form.name" class="form-control" placeholder="Nama lengkap" required>
      </div>
      <div class="col-md-6">
        <label class="form-label">NIP <span class="text-danger">*</span></label>
        <input v-model="form.nip" class="form-control" placeholder="NIP" required>
      </div>
      <div class="col-md-6">
        <label class="form-label">Username</label>
        <input v-model="form.username" class="form-control" placeholder="Username (opsional)">
      </div>
      <div class="col-md-6">
        <label class="form-label">Email <span class="text-danger">*</span></label>
        <input v-model="form.email" type="email" class="form-control" placeholder="email@instansi.go.id" required>
      </div>
      <div class="col-md-6">
        <label class="form-label">Unit Kerja</label>
        <input v-model="form.unitKerja" class="form-control" placeholder="Unit kerja">
      </div>
      <div class="col-md-6">
        <label class="form-label">Jabatan</label>
        <input v-model="form.jabatan" class="form-control" placeholder="Jabatan">
      </div>
      <div class="col-md-6">
        <label class="form-label">Role <span class="text-danger">*</span></label>
        <select v-model="form.role" class="form-select" required>
          <option value="" disabled>— Pilih role —</option>
          <option v-for="role in roleOptions" :key="role.key" :value="role.key">
            {{ role.name }}
          </option>
        </select>
      </div>
      <div class="col-md-6">
        <label class="form-label">Status Akun</label>
        <div class="form-check form-switch mt-2">
          <input
            v-model="form.isActive"
            class="form-check-input"
            type="checkbox"
            id="user-active"
            :disabled="formMode === 'edit' && editingId === currentUser?.id"
          >
          <label class="form-check-label" for="user-active">
            {{ form.isActive ? 'Aktif' : 'Nonaktif' }}
            <span v-if="formMode === 'edit' && editingId === currentUser?.id" class="text-muted small ms-1">(Akun sendiri)</span>
          </label>
        </div>
      </div>
      <!-- Pemisah Keamanan & Kata Sandi di Bagian Bawah -->
      <div class="col-12">
        <hr class="my-2 text-muted opacity-25">
        <div class="d-flex align-items-center justify-content-between mb-2">
          <label class="form-label fw-bold text-dark mb-0">
            <i class="ti ti-lock me-1"></i>Keamanan & Kata Sandi
          </label>
          <div class="form-check form-switch mb-0">
            <input v-model="form.isSso" class="form-check-input" type="checkbox" id="user-sso">
            <label class="form-check-label small" for="user-sso">Akun SSO (tanpa kata sandi lokal)</label>
          </div>
        </div>
      </div>

      <template v-if="!form.isSso">
        <div class="col-md-6">
          <label class="form-label">
            {{ formMode === 'tambah' ? 'Kata Sandi' : 'Kata Sandi Baru (opsional)' }}
            <span v-if="formMode === 'tambah'" class="text-danger">*</span>
          </label>
          <input
            v-model="form.password"
            type="password"
            autocomplete="new-password"
            class="form-control"
            :placeholder="formMode === 'tambah' ? 'Minimal 8 karakter' : 'Kosongkan jika tidak diubah'"
          >
        </div>
        <div class="col-md-6">
          <label class="form-label">
            {{ formMode === 'tambah' ? 'Konfirmasi Kata Sandi' : 'Konfirmasi Kata Sandi Baru' }}
            <span v-if="formMode === 'tambah' || form.password" class="text-danger">*</span>
          </label>
          <input
            v-model="form.passwordConfirmation"
            type="password"
            autocomplete="new-password"
            class="form-control"
            :placeholder="formMode === 'tambah' ? 'Ulangi kata sandi' : 'Ulangi jika mengubah kata sandi'"
          >
        </div>
      </template>
      <div v-else class="col-12">
        <div class="alert alert-light border small text-muted mb-0 py-2">
          <i class="ti ti-info-circle me-1"></i>Akun SSO melakukan autentikasi melalui Identity Provider terintegrasi dan tidak memerlukan kata sandi lokal.
        </div>
      </div>
    </div>

    <div class="d-flex justify-content-end gap-2 mt-4">
      <button type="button" class="btn btn-secondary" @click="showForm = false">Batal</button>
      <button type="button" class="btn btn-primary" :disabled="saving" @click="simpanUser">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        Simpan
      </button>
    </div>
  </UiAppModal>

  <!-- Modal Reset Kata Sandi -->
  <UiAppModal v-if="resetTarget" title="Reset Kata Sandi" size="sm" @close="resetTarget = null">
    <p class="mb-2">
      Reset kata sandi untuk <strong>{{ resetTarget.nama }}</strong> ({{ resetTarget.nip }}).
    </p>
    <div class="mb-2">
      <label class="form-label">Kata Sandi Baru</label>
      <input
        v-model="resetPwdForm.password"
        type="password"
        autocomplete="new-password"
        class="form-control"
        placeholder="Minimal 8 karakter"
      >
    </div>
    <div class="mb-2">
      <label class="form-label">Konfirmasi Kata Sandi</label>
      <input
        v-model="resetPwdForm.passwordConfirmation"
        type="password"
        autocomplete="new-password"
        class="form-control"
      >
    </div>
    <div class="d-flex justify-content-end gap-2 mt-3">
      <button type="button" class="btn btn-secondary" @click="resetTarget = null">Batal</button>
      <button type="button" class="btn btn-warning" :disabled="saving" @click="submitResetPassword">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        Reset
      </button>
    </div>
  </UiAppModal>

  <!-- Modal Konfirmasi Hapus / Restore -->
  <UiAppModal v-if="confirmTarget" :title="confirmMode === 'hapus' ? 'Hapus Pengguna' : 'Pulihkan Pengguna'" size="sm" @close="confirmTarget = null">
    <p v-if="confirmMode === 'hapus'" class="mb-0">
      Hapus pengguna <strong>{{ confirmTarget.nama }}</strong>? Data akan dipindahkan ke tempat sampah (soft delete)
      dan dapat dipulihkan kembali.
    </p>
    <p v-else class="mb-0">
      Pulihkan pengguna <strong>{{ confirmTarget.nama }}</strong> dari tempat sampah?
    </p>
    <div class="d-flex justify-content-end gap-2 mt-3">
      <button type="button" class="btn btn-secondary" @click="confirmTarget = null">Batal</button>
      <button
        type="button"
        class="btn"
        :class="confirmMode === 'hapus' ? 'btn-danger' : 'btn-success'"
        :disabled="saving"
        @click="submitConfirm"
      >
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        {{ confirmMode === 'hapus' ? 'Hapus' : 'Pulihkan' }}
      </button>
    </div>
  </UiAppModal>

  <!-- Modal Reset MFA Pengguna (Admin Step-Up) -->
  <UiAppModal v-if="mfaResetTarget" title="Reset MFA Pengguna" size="sm" @close="mfaResetTarget = null">
    <p class="mb-2">
      Reset status MFA untuk <strong>{{ mfaResetTarget.nama }}</strong> ({{ mfaResetTarget.nip }})? Pengguna akan diminta melakukan konfigurasi ulang MFA pada login berikutnya.
    </p>
    <div class="alert alert-warning py-2 small mb-3">
      <i class="ti ti-shield-alert me-1"></i><strong>Step-Up Verification:</strong> Masukkan 6 digit kode OTP Administrator Anda untuk menyetujui reset ini.
    </div>
    <div class="mb-3">
      <label class="form-label">Kode OTP Administrator</label>
      <input
        v-model="adminMfaOtp"
        type="password"
        maxlength="6"
        autocomplete="one-time-code"
        class="form-control text-center font-monospace fs-5"
        placeholder="000000"
      >
    </div>
    <div class="d-flex justify-content-end gap-2 mt-3">
      <button type="button" class="btn btn-secondary" @click="mfaResetTarget = null">Batal</button>
      <button type="button" class="btn btn-danger" :disabled="saving || !adminMfaOtp" @click="submitResetMfa">
        <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
        Reset MFA
      </button>
    </div>
  </UiAppModal>
</template>

<script setup lang="ts">
import type { UserRow } from '~/types/index'
import type { RoleItem, PermissionGroup } from '~/modules/user/RbacMatrix.vue'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

useHead({ title: 'Manajemen User — Arsip Record Center KemenPPPA' })

const { apiFetch, refreshProfile, getUser } = useSession()
const { showToast } = useToast()
const currentUser = computed(() => getUser())

// ─── State ──────────────────────────────────────────────────────
const activeTab = ref<'users' | 'rbac'>('users')
const userList = ref<UserRow[]>([])
const initialLoading = ref(true)
const refreshing = ref(false)
const saving = ref(false)

const roles = ref<RoleItem[]>([])
const permissionGroups = ref<PermissionGroup[]>([])
const rolesLoading = ref(true)
const savingRole = ref(false)

const stats = reactive({
  total: 0,
  aktif: 0,
  nonaktif: 0,
  administrator: 0,
})

const filter = reactive({
  search: '',
  role: '',
  status: '',
  trashed: '',
})

const hasActiveFilter = computed(() =>
  Boolean(filter.search || filter.role || filter.status || filter.trashed)
)

// Filter instan di sisi client/DataTables tanpa berkedip atau me-reload seluruh card/halaman
const filteredUserList = computed(() => {
  const q = filter.search.trim().toLowerCase()
  return userList.value.filter((u) => {
    if (filter.role && u.role !== filter.role) return false
    if (filter.status === 'active' && !u.isActive) return false
    if (filter.status === 'inactive' && u.isActive) return false
    if (q) {
      const matchNama = (u.nama || '').toLowerCase().includes(q)
      const matchUsername = (u.username || '').toLowerCase().includes(q)
      const matchNip = (u.nip || '').toLowerCase().includes(q)
      const matchEmail = (u.email || '').toLowerCase().includes(q)
      const matchUnor = (u.unitKerja || '').toLowerCase().includes(q)
      const matchJabatan = (u.jabatan || '').toLowerCase().includes(q)
      if (!matchNama && !matchUsername && !matchNip && !matchEmail && !matchUnor && !matchJabatan) return false
    }
    return true
  })
})

const statCards = computed(() => [
  { label: 'Total Pengguna', value: stats.total, icon: 'ti ti-users', iconBg: 'bg-primary-subtle text-primary' },
  { label: 'Pengguna Aktif', value: stats.aktif, icon: 'ti ti-user-check', iconBg: 'bg-success-subtle text-success' },
  { label: 'Nonaktif', value: stats.nonaktif, icon: 'ti ti-user-off', iconBg: 'bg-warning-subtle text-warning' },
  { label: 'Administrator', value: stats.administrator, icon: 'ti ti-shield-check', iconBg: 'bg-danger-subtle text-danger' },
])

const roleOptions = computed(() => {
  if (roles.value.length > 0) {
    return roles.value.map((r) => ({ key: r.key, name: r.name }))
  }
  return [
    { key: 'administrator', name: 'Administrator' },
    { key: 'operator', name: 'Operator' },
    { key: 'approver', name: 'Approver' },
    { key: 'viewer', name: 'Viewer' },
  ]
})

// ─── Fetch ──────────────────────────────────────────────────────
async function fetchStats() {
  try {
    const res = await apiFetch('/users/stats')
    const json = await res.json()
    if (json.success && json.data) {
      stats.total = Number(json.data.total ?? 0)
      stats.aktif = Number(json.data.aktif ?? 0)
      stats.nonaktif = Number(json.data.nonaktif ?? 0)
      stats.administrator = Number(json.data.administrator ?? 0)
    }
  } catch (err) {
    console.error('[fetchStats]', err)
  }
}

async function fetchUsers(isBackground = false) {
  if (isBackground) {
    refreshing.value = true
  } else {
    initialLoading.value = true
  }
  try {
    const params = new URLSearchParams({ per_page: '100' })
    if (filter.trashed) params.set('trashed', filter.trashed)

    const res = await apiFetch(`/users?${params.toString()}`)
    const json = await res.json()

    if (json.success) {
      const rawData = Array.isArray(json.data) ? json.data : []
      userList.value = rawData.map((u: any) => ({
        id: Number(u.id),
        username: u.username,
        nama: u.name,
        nip: u.nip,
        email: u.email,
        unitKerja: u.unit_kerja || '-',
        jabatan: u.jabatan || '',
        role: u.role,
        roleLabel: u.role_label || u.role,
        isSso: Boolean(u.is_sso),
        isActive: Boolean(u.is_active),
        mfaEnabled: Boolean(u.mfa_enabled),
        lastLoginAt: u.last_login_at ? formatDateTimeIndo(u.last_login_at) : '-',
        createdAt: u.created_at ? formatDateTimeIndo(u.created_at) : '-',
        deletedAt: u.deleted_at || null,
      }))
    } else {
      showToast(json.message || 'Gagal memuat data pengguna', 'danger')
    }
  } catch (err) {
    if (err instanceof Error && err.message !== 'Sesi telah berakhir') {
      showToast(err.message, 'danger')
    }
  } finally {
    initialLoading.value = false
    refreshing.value = false
  }
}

async function fetchRoles() {
  rolesLoading.value = true
  try {
    const res = await apiFetch('/roles')
    const json = await res.json()
    if (json.success && json.data) {
      roles.value = Array.isArray(json.data.roles) ? json.data.roles : []
      permissionGroups.value = Array.isArray(json.data.permission_groups) ? json.data.permission_groups : []
    } else {
      showToast(json.message || 'Gagal memuat data peran', 'danger')
    }
  } catch (err) {
    if (err instanceof Error && err.message !== 'Sesi telah berakhir') {
      showToast(err.message, 'danger')
    }
  } finally {
    rolesLoading.value = false
  }
}

function resetFilter() {
  const hadTrashed = Boolean(filter.trashed)
  filter.search = ''
  filter.role = ''
  filter.status = ''
  filter.trashed = ''
  if (hadTrashed) {
    fetchUsers(true)
  }
}

// ─── Form Tambah/Edit ───────────────────────────────────────────
const showForm = ref(false)
const formMode = ref<'tambah' | 'edit'>('tambah')
const formErrors = ref<string[]>([])

const form = reactive({
  id: 0,
  username: '',
  name: '',
  nip: '',
  email: '',
  unitKerja: '',
  jabatan: '',
  role: '',
  password: '',
  passwordConfirmation: '',
  isSso: false,
  isActive: true,
})

function resetForm() {
  form.id = 0
  form.username = ''
  form.name = ''
  form.nip = ''
  form.email = ''
  form.unitKerja = ''
  form.jabatan = ''
  form.role = ''
  form.password = ''
  form.passwordConfirmation = ''
  form.isSso = false
  form.isActive = true
  formErrors.value = []
}

function openTambah() {
  resetForm()
  formMode.value = 'tambah'
  showForm.value = true
}

function openEdit(row: UserRow) {
  resetForm()
  formMode.value = 'edit'
  form.id = row.id
  form.username = row.username || ''
  form.name = row.nama
  form.nip = row.nip
  form.email = row.email || ''
  form.unitKerja = row.unitKerja === '-' ? '' : row.unitKerja
  form.jabatan = row.jabatan || ''
  form.role = row.role
  form.isSso = row.isSso
  form.isActive = row.isActive
  showForm.value = true
}

async function simpanUser() {
  formErrors.value = []

  if (!form.name || !form.nip || !form.email || !form.role) {
    formErrors.value.push('Nama, NIP, email, dan role wajib diisi.')
    return
  }
  if (!form.isSso && formMode.value === 'tambah' && form.password.length < 8) {
    formErrors.value.push('Kata sandi minimal 8 karakter untuk akun non-SSO.')
    return
  }
  if (form.password && form.password !== form.passwordConfirmation) {
    formErrors.value.push('Konfirmasi kata sandi tidak cocok.')
    return
  }

  saving.value = true
  try {
    const body: Record<string, any> = {
      username: form.username || form.nip,
      name: form.name,
      nip: form.nip,
      email: form.email,
      role: form.role,
      unit_kerja: form.unitKerja || null,
      jabatan: form.jabatan || null,
      is_sso: form.isSso,
      is_active: form.isActive,
    }
    if (form.password) {
      body.password = form.password
      body.password_confirmation = form.passwordConfirmation
    }

    const url = formMode.value === 'tambah' ? '/users' : `/users/${form.id}`
    const method = formMode.value === 'tambah' ? 'POST' : 'PUT'

    const res = await apiFetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const json = await res.json()

    if (json.success) {
      showToast(
        formMode.value === 'tambah' ? 'Pengguna berhasil ditambahkan' : 'Pengguna berhasil diperbarui',
        'success'
      )
      showForm.value = false
      await Promise.all([fetchUsers(true), fetchStats()])
    } else {
      const errs: string[] = []
      if (json.errors && typeof json.errors === 'object') {
        Object.values(json.errors).forEach((v: any) => {
          if (Array.isArray(v)) errs.push(...v)
          else if (typeof v === 'string') errs.push(v)
        })
      }
      formErrors.value = errs.length ? errs : [json.message || 'Gagal menyimpan pengguna']
    }
  } catch (err) {
    formErrors.value = [err instanceof Error ? err.message : 'Tidak dapat terhubung ke server']
  } finally {
    saving.value = false
  }
}

// ─── Reset Password ─────────────────────────────────────────────
const resetTarget = ref<UserRow | null>(null)
const resetPwdForm = reactive({ password: '', passwordConfirmation: '' })

function openResetPassword(row: UserRow) {
  resetTarget.value = row
  resetPwdForm.password = ''
  resetPwdForm.passwordConfirmation = ''
}

async function submitResetPassword() {
  if (!resetTarget.value) return
  if (resetPwdForm.password.length < 8) {
    showToast('Kata sandi baru minimal 8 karakter', 'warning')
    return
  }
  if (resetPwdForm.password !== resetPwdForm.passwordConfirmation) {
    showToast('Konfirmasi kata sandi tidak cocok', 'warning')
    return
  }

  saving.value = true
  try {
    const row = resetTarget.value
    const res = await apiFetch(`/users/${row.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: row.username || row.nip,
        name: row.nama,
        nip: row.nip,
        email: row.email,
        role: row.role,
        unit_kerja: row.unitKerja === '-' ? null : row.unitKerja,
        jabatan: row.jabatan || null,
        is_sso: row.isSso,
        is_active: row.isActive,
        password: resetPwdForm.password,
        password_confirmation: resetPwdForm.passwordConfirmation,
      }),
    })
    const json = await res.json()
    if (json.success) {
      showToast('Kata sandi berhasil direset', 'success')
      resetTarget.value = null
    } else {
      showToast(json.message || 'Gagal mereset kata sandi', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    saving.value = false
  }
}

// ─── Status / Hapus / Restore ───────────────────────────────────
const confirmTarget = ref<UserRow | null>(null)
const confirmMode = ref<'hapus' | 'restore'>('hapus')

// ─── Reset MFA (Admin Step-Up) ──────────────────────────────────
const mfaResetTarget = ref<UserRow | null>(null)
const adminMfaOtp = ref('')

function openResetMfa(row: UserRow) {
  mfaResetTarget.value = row
  adminMfaOtp.value = ''
}

async function submitResetMfa() {
  if (!mfaResetTarget.value || !adminMfaOtp.value) return
  saving.value = true
  try {
    const row = mfaResetTarget.value
    const res = await apiFetch(`/users/${row.id}/mfa/reset`, {
      method: 'POST',
      body: JSON.stringify({ code: adminMfaOtp.value.trim() }),
    })
    const json = await res.json()
    if (json.success) {
      showToast(`MFA untuk pengguna ${row.nama} berhasil direset`, 'success')
      mfaResetTarget.value = null
      await fetchUsers(true)
    } else {
      showToast(json.message || 'Gagal mereset MFA pengguna', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    saving.value = false
  }
}

function onNonaktif(row: UserRow) {
  toggleStatus(row, false)
}

function onAktif(row: UserRow) {
  toggleStatus(row, true)
}

async function toggleStatus(row: UserRow, isActive: boolean) {
  saving.value = true
  try {
    const res = await apiFetch(`/users/${row.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: row.username || row.nip,
        name: row.nama,
        nip: row.nip,
        email: row.email,
        role: row.role,
        unit_kerja: row.unitKerja === '-' ? null : row.unitKerja,
        jabatan: row.jabatan || null,
        is_sso: row.isSso,
        is_active: isActive,
      }),
    })
    const json = await res.json()
    if (json.success) {
      showToast(isActive ? 'Pengguna diaktifkan' : 'Pengguna dinonaktifkan', 'success')
      await Promise.all([fetchUsers(true), fetchStats()])
    } else {
      showToast(json.message || 'Gagal mengubah status pengguna', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    saving.value = false
  }
}

function konfirmasiHapus(row: UserRow) {
  confirmTarget.value = row
  confirmMode.value = 'hapus'
}

function konfirmasiRestore(row: UserRow) {
  confirmTarget.value = row
  confirmMode.value = 'restore'
}

async function submitConfirm() {
  if (!confirmTarget.value) return
  saving.value = true
  try {
    const row = confirmTarget.value
    const url = confirmMode.value === 'hapus' ? `/users/${row.id}` : `/users/${row.id}/restore`
    const method = confirmMode.value === 'hapus' ? 'DELETE' : 'PUT'

    const res = await apiFetch(url, { method })
    const json = await res.json()

    if (json.success) {
      showToast(
        confirmMode.value === 'hapus' ? 'Pengguna berhasil dihapus' : 'Pengguna berhasil dipulihkan',
        'success'
      )
      confirmTarget.value = null
      await Promise.all([fetchUsers(true), fetchStats()])
    } else {
      showToast(json.message || 'Gagal memproses permintaan', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    saving.value = false
  }
}

// ─── RBAC ───────────────────────────────────────────────────────
async function savePermissions(roleId: number, permissions: string[]) {
  savingRole.value = true
  try {
    const res = await apiFetch(`/roles/${roleId}/permissions`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ permissions }),
    })
    const json = await res.json()
    if (json.success) {
      showToast('Hak akses berhasil disimpan', 'success')
      await fetchRoles()
      await refreshProfile()
    } else {
      showToast(json.message || 'Gagal menyimpan hak akses', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    savingRole.value = false
  }
}

async function addRole(name: string) {
  savingRole.value = true
  try {
    const res = await apiFetch('/roles', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, permissions: ['dashboard.view'] }),
    })
    const json = await res.json()
    if (json.success) {
      showToast('Peran baru berhasil ditambahkan', 'success')
      await fetchRoles()
    } else {
      showToast(json.message || 'Gagal menambahkan peran', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    savingRole.value = false
  }
}

async function deleteRole(roleId: number) {
  savingRole.value = true
  try {
    const res = await apiFetch(`/roles/${roleId}`, { method: 'DELETE' })
    const json = await res.json()
    if (json.success) {
      showToast('Peran berhasil dihapus', 'success')
      await fetchRoles()
    } else {
      showToast(json.message || 'Gagal menghapus peran', 'danger')
    }
  } catch (err) {
    showToast(err instanceof Error ? err.message : 'Tidak dapat terhubung ke server', 'danger')
  } finally {
    savingRole.value = false
  }
}

onMounted(() => {
  fetchUsers()
  fetchStats()
  fetchRoles()
})
</script>

<style scoped>
.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.35rem;
}
</style>
