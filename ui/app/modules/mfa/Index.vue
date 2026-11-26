<template>
  <div>
    <TemplatePageTitle
      title="Pengaturan Autentikasi Dua Faktor (MFA)"
      :breadcrumbs="[{ label: 'Pengaturan Akun' }, { label: 'MFA' }]"
    />

    <div class="row">
      <div class="col-12 col-lg-8">
        <!-- Card Status MFA -->
        <div class="card mb-4">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0"><i class="ti ti-shield-lock me-2"></i>Status Keamanan Akun</h5>
            <span
              :class="[
                'badge fs-6 px-3 py-2',
                status.enabled ? 'bg-light-success text-success border border-success-subtle' : 'bg-light-warning text-warning border border-warning-subtle'
              ]"
            >
              {{ status.enabled ? 'MFA Aktif' : 'MFA Belum Aktif' }}
            </span>
          </div>
          <div class="card-body">
            <div v-if="loadingStatus" class="text-center py-4">
              <span class="spinner-border text-primary" role="status"></span>
              <p class="text-muted small mt-2">Memeriksa status MFA...</p>
            </div>

            <div v-else>
              <div class="d-flex align-items-center gap-3 mb-3">
                <div
                  :class="[
                    'avatar rounded-circle d-flex align-items-center justify-content-center',
                    status.enabled ? 'bg-light-success text-success' : 'bg-light-warning text-warning'
                  ]"
                  style="width: 56px; height: 56px;"
                >
                  <i :class="status.enabled ? 'ti ti-shield-check fs-2' : 'ti ti-shield-x fs-2'"></i>
                </div>
                <div>
                  <h6 class="mb-1">
                    {{ status.enabled ? 'Autentikasi dua faktor aktif melindungi akun Anda' : 'Akun Anda belum dilindungi autentikasi dua faktor' }}
                  </h6>
                  <p class="text-muted small mb-0">
                    <span v-if="status.enabled && status.enabledAt">
                      Diaktifkan pada: {{ formatDate(status.enabledAt) }}
                    </span>
                    <span v-else>
                      MFA diwajibkan oleh kebijakan keamanan sistem SiARSIP KemenPPPA.
                    </span>
                  </p>
                </div>
              </div>

              <!-- Info Recovery Codes -->
              <div v-if="status.enabled" class="p-3 bg-light rounded border mb-3">
                <div class="d-flex justify-content-between align-items-center">
                  <div>
                    <span class="fw-semibold text-dark d-block">Sisa Recovery Code Cadangan:</span>
                    <small class="text-muted">Kode darurat untuk login jika kehilangan aplikasi authenticator.</small>
                  </div>
                  <span class="badge bg-primary fs-6">{{ status.recoveryCodesRemaining }} / 10 Tersedia</span>
                </div>
              </div>

              <!-- Alert error / success -->
              <div v-if="alertMsg" :class="['alert py-2 px-3 mb-3 d-flex align-items-center gap-2', isAlertError ? 'alert-danger' : 'alert-success']">
                <i :class="isAlertError ? 'ti ti-alert-circle' : 'ti ti-circle-check'"></i>
                <span>{{ alertMsg }}</span>
              </div>

              <!-- Form Verifikasi Faktor Aktif (Step-Up) sebelum Reset -->
              <div v-if="status.enabled && !setupChallenge" class="mt-4 pt-3 border-top">
                <h6 class="mb-2">Kelola Kunci & Recovery Code</h6>
                <p class="text-muted small mb-3">
                  Untuk alasan keamanan, masukkan 6 digit kode OTP dari authenticator atau 1 buah recovery code Anda saat ini sebelum melakukan perubahan.
                </p>

                <div class="form-group mb-3">
                  <label class="form-label">Kode OTP / Recovery Code Saat Ini</label>
                  <input
                    v-model="verificationCode"
                    type="password"
                    class="form-control"
                    placeholder="Masukkan 6 digit OTP atau kode pemulihan"
                    autocomplete="one-time-code"
                  >
                </div>

                <div class="d-flex flex-wrap gap-2">
                  <button
                    type="button"
                    class="btn btn-outline-primary"
                    :disabled="actionLoading || !verificationCode"
                    @click="handleResetRecoveryCodes"
                  >
                    <span v-if="actionLoading" class="spinner-border spinner-border-sm me-1"></span>
                    <i v-else class="ti ti-refresh me-1"></i> Buat Ulang 10 Recovery Code
                  </button>
                  <button
                    type="button"
                    class="btn btn-outline-danger"
                    :disabled="actionLoading || !verificationCode"
                    @click="handleResetSelf"
                  >
                    <span v-if="actionLoading" class="spinner-border spinner-border-sm me-1"></span>
                    <i v-else class="ti ti-lock-off me-1"></i> Reset & Re-Setup MFA
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Panel Tampilan Recovery Codes Baru -->
        <div v-if="newRecoveryCodes.length > 0" class="card mb-4 border-primary">
          <div class="card-header bg-light-primary border-primary">
            <h5 class="text-primary mb-0"><i class="ti ti-key me-2"></i>10 Recovery Codes Baru</h5>
          </div>
          <div class="card-body">
            <p class="text-muted small mb-3">
              Kode-kode pemulihan berikut <strong>hanya ditampilkan satu kali</strong> saat ini. Silakan unduh atau salin dan simpan di tempat aman.
            </p>
            <div class="p-3 bg-light rounded border mb-3">
              <div class="row g-2 font-monospace small">
                <div v-for="(code, idx) in newRecoveryCodes" :key="idx" class="col-6 text-center py-1">
                  <span class="badge bg-white text-dark border px-2 py-1 w-100">{{ code }}</span>
                </div>
              </div>
            </div>
            <button type="button" class="btn btn-primary" @click="downloadNewRecoveryCodes">
              <i class="ti ti-download me-1"></i>Unduh Kode (.txt)
            </button>
          </div>
        </div>

        <!-- Panel Re-Setup MFA -->
        <div v-if="setupChallenge" class="card mb-4 border-warning">
          <div class="card-header bg-light-warning border-warning">
            <h5 class="text-warning-emphasis mb-0"><i class="ti ti-qrcode me-2"></i>Konfigurasi Ulang Aplikasi Authenticator</h5>
          </div>
          <div class="card-body text-center">
            <p class="text-muted small mb-3">Pindai kode QR baru berikut di aplikasi authenticator Anda:</p>
            <CAuthMfaQr :value="setupOtpauthUri" :size="200" />

            <div class="my-3">
              <label class="form-label small text-muted mb-1">Kunci Rahasia Manual:</label>
              <div class="p-2 bg-light rounded font-monospace small user-select-all border mx-auto" style="max-width: 400px;">
                {{ setupSecret }}
              </div>
            </div>

            <div class="mx-auto" style="max-width: 350px;">
              <div class="form-group mb-3">
                <label class="form-label">Masukkan 6 Digit OTP Baru</label>
                <input
                  v-model="newTotpCode"
                  type="text"
                  maxlength="6"
                  class="form-control text-center font-monospace fs-4"
                  placeholder="000000"
                  inputmode="numeric"
                >
              </div>
              <div class="d-grid gap-2">
                <button
                  type="button"
                  class="btn btn-warning"
                  :disabled="actionLoading || !newTotpCode"
                  @click="handleVerifyNewSetup"
                >
                  <span v-if="actionLoading" class="spinner-border spinner-border-sm me-1"></span>
                  Konfirmasi Setup Baru
                </button>
                <button type="button" class="btn btn-link btn-sm text-secondary" @click="setupChallenge = ''">
                  Batal
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Panduan MFA -->
      <div class="col-12 col-lg-4">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0"><i class="ti ti-info-circle me-2"></i>Panduan MFA</h5>
          </div>
          <div class="card-body small text-muted">
            <p><strong>Apa itu MFA?</strong><br>Autentikasi dua faktor menambahkan lapisan keamanan ekstra selain kata sandi Anda.</p>
            <p><strong>Aplikasi yang Didukung:</strong></p>
            <ul class="ps-3 mb-3">
              <li>Google Authenticator (Android / iOS)</li>
              <li>Microsoft Authenticator</li>
              <li>2FAS / Aegis Authenticator</li>
              <li>1Password / Bitwarden</li>
            </ul>
            <p><strong>Catatan Recovery Code:</strong><br>Setiap kode pemulihan hanya dapat dipakai 1 kali. Jika sisa kode tinggal sedikit, gunakan tombol <em>Buat Ulang 10 Recovery Code</em>.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { API_BASE_URL } from '~/constants/auth'

const { apiFetch, verifyMfa, getUser } = useSession()

const loadingStatus = ref(true)
const actionLoading = ref(false)
const alertMsg = ref('')
const isAlertError = ref(false)

const status = ref<{
  enabled: boolean
  enabledAt?: string | null
  recoveryCodesRemaining: number
}>({
  enabled: false,
  recoveryCodesRemaining: 0,
})

const verificationCode = ref('')
const newRecoveryCodes = ref<string[]>([])

// State Re-Setup
const setupChallenge = ref('')
const setupSecret = ref('')
const setupOtpauthUri = ref('')
const newTotpCode = ref('')

function showAlert(msg: string, isError = false) {
  alertMsg.value = msg
  isAlertError.value = isError
}

async function fetchStatus() {
  loadingStatus.value = true
  try {
    const res = await apiFetch('/auth/mfa/status')
    const json = await res.json()
    if (json.success && json.data) {
      status.value = {
        enabled: Boolean(json.data.enabled),
        enabledAt: json.data.enabled_at,
        recoveryCodesRemaining: json.data.recovery_codes_remaining || 0,
      }
    }
  } catch (err: any) {
    showAlert(err.message || 'Gagal memuat status MFA', true)
  } finally {
    loadingStatus.value = false
  }
}

async function handleResetRecoveryCodes() {
  if (!verificationCode.value) return
  actionLoading.value = true
  alertMsg.value = ''
  try {
    const res = await apiFetch('/auth/mfa/recovery-codes', {
      method: 'POST',
      body: JSON.stringify({ code: verificationCode.value.trim() }),
    })
    const json = await res.json()
    if (!res.ok || !json.success) {
      throw new Error(json.message || 'Gagal membuat ulang recovery codes')
    }
    newRecoveryCodes.value = json.data.recovery_codes || []
    verificationCode.value = ''
    showAlert('10 Recovery codes baru berhasil dibuat. Pastikan Anda mengunduh dan menyimpannya.')
    await fetchStatus()
  } catch (err: any) {
    showAlert(err.message, true)
  } finally {
    actionLoading.value = false
  }
}

async function handleResetSelf() {
  if (!verificationCode.value) return
  actionLoading.value = true
  alertMsg.value = ''
  try {
    const res = await apiFetch('/auth/mfa/reset', {
      method: 'POST',
      body: JSON.stringify({ code: verificationCode.value.trim() }),
    })
    const json = await res.json()
    if (!res.ok || !json.success) {
      throw new Error(json.message || 'Gagal mereset MFA')
    }
    setupChallenge.value = json.data.challenge
    setupSecret.value = json.data.secret
    setupOtpauthUri.value = json.data.otpauth_uri
    verificationCode.value = ''
    showAlert('MFA berhasil direset. Silakan lakukan setup ulang menggunakan QR Code di bawah.')
  } catch (err: any) {
    showAlert(err.message, true)
  } finally {
    actionLoading.value = false
  }
}

async function handleVerifyNewSetup() {
  if (!newTotpCode.value || !setupChallenge.value) return
  actionLoading.value = true
  alertMsg.value = ''
  try {
    const result = await verifyMfa(setupChallenge.value, newTotpCode.value.trim())
    if (!result.success) {
      throw new Error(result.message)
    }
    setupChallenge.value = ''
    newTotpCode.value = ''
    if (result.recoveryCodes && result.recoveryCodes.length > 0) {
      newRecoveryCodes.value = result.recoveryCodes
    }
    showAlert('Setup MFA baru berhasil diaktifkan!')
    await fetchStatus()
  } catch (err: any) {
    showAlert(err.message, true)
  } finally {
    actionLoading.value = false
  }
}

function downloadNewRecoveryCodes() {
  if (!newRecoveryCodes.value.length) return
  const user = getUser()
  const content = [
    '=======================================================',
    ' SiARSIP KemenPPPA — Multi-Factor Recovery Codes',
    ` Dibuat Pada: ${formatDateTimeIndo(new Date())}`,
    ` Akun: ${user?.username || user?.nip || 'akun'}`,
    '=======================================================',
    '',
    'PERINGATAN:',
    '- Setiap kode pemulihan di bawah ini HANYA DAPAT DIGUNAKAN 1 KALI.',
    '- Simpan file ini di tempat yang aman dan rahasia.',
    '- Jika Anda kehilangan ponsel / authenticator, gunakan salah satu kode ini.',
    '',
    ...newRecoveryCodes.value.map((c, i) => `[${(i + 1).toString().padStart(2, '0')}] ${c}`),
    '',
    '=======================================================',
  ].join('\n')

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `siarsip-recovery-codes.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

function formatDate(iso: string) {
  return formatDateTimeIndo(iso)
}

onMounted(() => {
  fetchStatus()
})
</script>
