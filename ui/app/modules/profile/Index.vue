<template>
  <div>
    <TemplatePageTitle
      title="Profil Saya"
      :breadcrumbs="[{ label: 'Pengaturan Akun' }, { label: 'Profil Saya' }]"
    />

    <div class="row g-4">
      <!-- Kolom Kiri: Data Identitas Pegawai (Readonly) -->
      <div class="col-12 col-lg-7">
        <div class="card h-100">
          <div class="card-header d-flex align-items-center justify-content-between">
            <h5 class="mb-0"><i class="ti ti-id me-2 text-primary"></i>Informasi Akun & Data Pegawai</h5>
            <span class="badge bg-light-primary text-primary border border-primary-subtle px-2 py-1">
              {{ profile.roleLabel || profile.roleKey }}
            </span>
          </div>

          <div class="card-body">
            <div class="alert alert-light border small text-muted d-flex align-items-center gap-2 mb-4">
              <i class="ti ti-info-circle text-primary fs-5 flex-shrink-0"></i>
              <div>
                Data identitas kedinasan (NIP, Unit Kerja, Jabatan, dan Peran) dikelola secara terpusat oleh <strong>Administrator Sistem</strong>. Hubungi administrator apabila terdapat pembaruan data kedinasan.
              </div>
            </div>

            <div class="row g-3">
              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">Nama Lengkap</label>
                <input
                  type="text"
                  class="form-control bg-light"
                  :value="profile.name || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">NIP (Nomor Induk Pegawai)</label>
                <input
                  type="text"
                  class="form-control bg-light font-monospace"
                  :value="profile.nip || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">Username Akun</label>
                <input
                  type="text"
                  class="form-control bg-light font-monospace text-primary"
                  :value="profile.username || profile.nip || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">Email Kedinasan</label>
                <input
                  type="text"
                  class="form-control bg-light"
                  :value="profile.email || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">Unit Kerja / Pencipta Arsip</label>
                <input
                  type="text"
                  class="form-control bg-light"
                  :value="profile.unit || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold">Jabatan</label>
                <input
                  type="text"
                  class="form-control bg-light"
                  :value="profile.jabatan || '-'"
                  readonly
                >
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold d-block">Status Akun</label>
                <span class="badge bg-light-success text-success border border-success-subtle px-3 py-2">
                  <i class="ti ti-check me-1"></i>Akun Aktif
                </span>
              </div>

              <div class="col-md-6">
                <label class="form-label small text-muted fw-semibold d-block">Tipe Autentikasi</label>
                <span v-if="profile.isSso" class="badge bg-light-info text-info border border-info-subtle px-3 py-2">
                  <i class="ti ti-link me-1"></i>Single Sign-On (SSO)
                </span>
                <span v-else class="badge bg-light-secondary text-secondary border border-secondary-subtle px-3 py-2">
                  <i class="ti ti-key me-1"></i>Autentikasi Lokal
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Kolom Kanan: Ubah Kata Sandi & Pintasan MFA -->
      <div class="col-12 col-lg-5 d-flex flex-column gap-4">
        <!-- Card Ubah Kata Sandi -->
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0"><i class="ti ti-lock me-2 text-warning"></i>Ubah Kata Sandi</h5>
          </div>

          <div class="card-body">
            <div v-if="profile.isSso" class="alert alert-info small mb-0">
              <i class="ti ti-info-circle me-1"></i>
              Akun Anda terhubung dengan <strong>Single Sign-On (SSO)</strong> Kemen PPPA. Perubahan kata sandi dilakukan secara terpusat melalui portal SSO instansi.
            </div>

            <form v-else @submit.prevent="submitChangePassword">
              <div v-if="pwdErrors.length" class="alert alert-danger py-2 small mb-3">
                <ul class="mb-0 ps-3">
                  <li v-for="(err, i) in pwdErrors" :key="i">{{ err }}</li>
                </ul>
              </div>

              <div class="mb-3">
                <label class="form-label small fw-semibold">
                  Kata Sandi Saat Ini <span class="text-danger">*</span>
                </label>
                <input
                  v-model="pwdForm.oldPassword"
                  type="password"
                  autocomplete="current-password"
                  class="form-control"
                  placeholder="Masukkan kata sandi saat ini"
                  required
                >
              </div>

              <div class="mb-3">
                <label class="form-label small fw-semibold">
                  Kata Sandi Baru <span class="text-danger">*</span>
                </label>
                <input
                  v-model="pwdForm.newPassword"
                  type="password"
                  autocomplete="new-password"
                  class="form-control"
                  placeholder="Minimal 8 karakter"
                  required
                >
              </div>

              <div class="mb-3">
                <label class="form-label small fw-semibold">
                  Konfirmasi Kata Sandi Baru <span class="text-danger">*</span>
                </label>
                <input
                  v-model="pwdForm.confirmPassword"
                  type="password"
                  autocomplete="new-password"
                  class="form-control"
                  placeholder="Ulangi kata sandi baru"
                  required
                >
              </div>

              <div class="d-flex justify-content-end">
                <button
                  type="submit"
                  class="btn btn-primary d-inline-flex align-items-center gap-2"
                  :disabled="pwdSaving"
                >
                  <span v-if="pwdSaving" class="spinner-border spinner-border-sm"></span>
                  <i v-else class="ti ti-device-floppy"></i>
                  <span>Perbarui Kata Sandi</span>
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Card Pintasan Pengaturan Keamanan MFA -->
        <div class="card">
          <div class="card-header d-flex align-items-center justify-content-between">
            <h5 class="mb-0"><i class="ti ti-shield-lock me-2 text-success"></i>Autentikasi Dua Faktor (MFA)</h5>
            <span
              :class="[
                'badge px-2 py-1',
                mfaStatusEnabled ? 'bg-light-success text-success border border-success-subtle' : 'bg-light-warning text-warning border border-warning-subtle'
              ]"
            >
              {{ mfaStatusEnabled ? 'MFA Aktif' : 'MFA Belum Aktif' }}
            </span>
          </div>

          <div class="card-body">
            <p class="text-muted small mb-3">
              Lindungi akun SiARSIP Anda dengan autentikasi dua faktor (TOTP Authenticator App) dan simpan kode cadangan pemulihan darurat.
            </p>

            <NuxtLink to="/setting/mfa" class="btn btn-outline-primary btn-sm w-100 d-inline-flex align-items-center justify-content-center gap-1">
              <i class="ti ti-settings"></i>
              <span>Buka Pengaturan MFA & Recovery Codes</span>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '~/composables/useToast'
import { useSession } from '~/composables/useSession'

const { getUser, apiFetch } = useSession()
const { showToast } = useToast()

const profile = ref({
  id: 0,
  username: '',
  name: '',
  nip: '',
  email: '',
  roleKey: '',
  roleLabel: '',
  unit: '',
  jabatan: '',
  isSso: false,
})

const mfaStatusEnabled = ref(false)

const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const pwdSaving = ref(false)
const pwdErrors = ref<string[]>([])

async function loadProfile() {
  const localUser = getUser()
  if (localUser) {
    profile.value = {
      id: localUser.id,
      username: localUser.username || localUser.nip,
      name: localUser.name,
      nip: localUser.nip,
      email: localUser.email,
      roleKey: localUser.roleKey,
      roleLabel: localUser.roleLabel,
      unit: localUser.unit,
      jabatan: localUser.jabatan,
      isSso: localUser.isSso,
    }
  }

  try {
    const res = await apiFetch('/auth/me')
    const json = await res.json()
    if (json.success && json.data) {
      const u = json.data
      profile.value = {
        id: u.id,
        username: u.username || u.nip,
        name: u.name,
        nip: u.nip,
        email: u.email,
        roleKey: u.role,
        roleLabel: u.role_label || u.role,
        unit: u.unit_kerja || '-',
        jabatan: u.jabatan || '-',
        isSso: Boolean(u.is_sso),
      }
      mfaStatusEnabled.value = Boolean(u.mfa_enabled)
    }
  } catch (_) {}
}

async function submitChangePassword() {
  pwdErrors.value = []

  if (!pwdForm.oldPassword) {
    pwdErrors.value.push('Kata sandi saat ini wajib diisi.')
    return
  }
  if (!pwdForm.newPassword || pwdForm.newPassword.length < 8) {
    pwdErrors.value.push('Kata sandi baru minimal 8 karakter.')
    return
  }
  if (pwdForm.newPassword !== pwdForm.confirmPassword) {
    pwdErrors.value.push('Konfirmasi kata sandi baru tidak sesuai.')
    return
  }

  pwdSaving.value = true
  try {
    const res = await apiFetch('/auth/change-password', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        old_password: pwdForm.oldPassword,
        new_password: pwdForm.newPassword,
        new_password_confirmation: pwdForm.confirmPassword,
      }),
    })
    const json = await res.json()

    if (json.success) {
      showToast('Kata sandi Anda berhasil diperbarui', 'success')
      pwdForm.oldPassword = ''
      pwdForm.newPassword = ''
      pwdForm.confirmPassword = ''
    } else {
      if (json.errors && typeof json.errors === 'object') {
        pwdErrors.value = Object.values(json.errors).flat() as string[]
      } else {
        pwdErrors.value = [json.message || 'Gagal mengubah kata sandi']
      }
      showToast(json.message || 'Gagal mengubah kata sandi', 'danger')
    }
  } catch (err) {
    if (err instanceof Error) {
      pwdErrors.value = [err.message]
      showToast(err.message, 'danger')
    }
  } finally {
    pwdSaving.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>
