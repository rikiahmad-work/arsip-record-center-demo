<template>
  <div>
    <!-- TAHAP 1: FORM KREDENSIAL (USERNAME & PASSWORD) -->
    <form v-if="mfaStage === 'credentials'" @submit.prevent="emit('submit')">
      <div class="d-flex justify-content-between align-items-end mb-4">
        <h3 class="mb-0"><b>Login SSO</b></h3>
      </div>

      <div v-if="errorMsg" class="alert alert-danger py-2 px-3 mb-3 d-flex align-items-center gap-2" role="alert">
        <i class="ti ti-alert-circle"></i>
        <span>{{ errorMsg }}</span>
      </div>

      <div class="form-group mb-3">
        <label class="form-label">NIP / Username</label>
        <input
          :value="username"
          type="text"
          class="form-control"
          placeholder="Masukkan NIP atau Username"
          autocomplete="username"
          @input="emit('update:username', ($event.target as HTMLInputElement).value)"
        >
      </div>

      <div class="form-group mb-3">
        <label class="form-label">Password</label>
        <div class="input-group">
          <input
            :value="password"
            :type="showPassword ? 'text' : 'password'"
            class="form-control"
            placeholder="Password"
            autocomplete="current-password"
            @input="emit('update:password', ($event.target as HTMLInputElement).value)"
          >
          <button
            type="button"
            class="btn btn-outline-secondary"
            tabindex="-1"
            @click="showPassword = !showPassword"
          >
            <i :class="showPassword ? 'ti ti-eye-off' : 'ti ti-eye'"></i>
          </button>
        </div>
      </div>

      <div class="d-flex mt-1 justify-content-between align-items-center">
        <div class="form-check">
          <input
            id="customCheckc1"
            class="form-check-input input-primary"
            type="checkbox"
            :checked="remember"
            @change="emit('update:remember', ($event.target as HTMLInputElement).checked)"
          >
          <label class="form-check-label text-muted" for="customCheckc1">Simpan Password</label>
        </div>
        <h5 class="text-secondary f-w-400 mb-0" style="font-size: 0.875rem;">Lupa Password?</h5>
      </div>

      <div class="d-grid mt-4">
        <button type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          Login
        </button>
      </div>
    </form>

    <!-- TAHAP 2A: SETUP MFA PERTAMA KALI -->
    <div v-else-if="mfaStage === 'setup'">
      <div class="text-center mb-3">
        <h4 class="mb-1"><b>Setup Autentikasi 2-Faktor</b></h4>
        <p class="text-muted small">Pindai kode QR menggunakan Google Authenticator atau aplikasi TOTP lainnya.</p>
      </div>

      <div v-if="errorMsg" class="alert alert-danger py-2 px-3 mb-3 d-flex align-items-center gap-2" role="alert">
        <i class="ti ti-alert-circle"></i>
        <span>{{ errorMsg }}</span>
      </div>

      <div class="text-center mb-3">
        <CAuthMfaQr :value="mfaOtpauthUri" :size="180" />
      </div>

      <div class="mb-3 text-center">
        <label class="form-label small text-muted mb-1">Atau masukkan kode rahasia secara manual:</label>
        <div class="p-2 bg-light rounded font-monospace small user-select-all border text-break">
          {{ mfaSecret }}
        </div>
      </div>

      <form @submit.prevent="emit('verifyMfa')">
        <div class="form-group mb-3">
          <label class="form-label text-center d-block">Masukkan 6 Digit Kode OTP</label>
          <input
            :value="mfaCode"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            maxlength="6"
            class="form-control text-center font-monospace fs-4 letter-spacing-2"
            placeholder="000000"
            autocomplete="one-time-code"
            autofocus
            @input="emit('update:mfaCode', ($event.target as HTMLInputElement).value)"
          >
        </div>

        <div class="d-grid mb-2">
          <button type="submit" class="btn btn-primary" :disabled="loading || !mfaCode">
            <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            Verifikasi & Aktifkan MFA
          </button>
        </div>

        <div class="text-center">
          <button type="button" class="btn btn-link btn-sm text-secondary" @click="emit('cancelMfa')">
            <i class="ti ti-arrow-left me-1"></i>Kembali ke Login
          </button>
        </div>
      </form>
    </div>

    <!-- TAHAP 2B: VERIFIKASI MFA (LOGIN RUTIN) -->
    <div v-else-if="mfaStage === 'verify'">
      <div class="text-center mb-3">
        <div class="avatar avatar-md bg-light-primary text-primary mx-auto mb-2 rounded-circle p-2 d-inline-flex align-items-center justify-content-center" style="width: 48px; height: 48px;">
          <i class="ti ti-shield-lock fs-3"></i>
        </div>
        <h4 class="mb-1"><b>Verifikasi Keamanan</b></h4>
        <p class="text-muted small">Buka aplikasi Authenticator Anda dan masukkan 6 digit kode yang tampil.</p>
      </div>

      <div v-if="errorMsg" class="alert alert-danger py-2 px-3 mb-3 d-flex align-items-center gap-2" role="alert">
        <i class="ti ti-alert-circle"></i>
        <span>{{ errorMsg }}</span>
      </div>

      <form @submit.prevent="emit('verifyMfa')">
        <div class="form-group mb-3">
          <label class="form-label text-center d-block">Kode Authenticator (TOTP)</label>
          <input
            :value="mfaCode"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            maxlength="6"
            class="form-control text-center font-monospace fs-4 letter-spacing-2"
            placeholder="000000"
            autocomplete="one-time-code"
            autofocus
            @input="emit('update:mfaCode', ($event.target as HTMLInputElement).value)"
          >
        </div>

        <div class="d-grid mb-3">
          <button type="submit" class="btn btn-primary" :disabled="loading || !mfaCode">
            <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            Verifikasi
          </button>
        </div>

        <div class="d-flex justify-content-between align-items-center">
          <button type="button" class="btn btn-link btn-sm text-secondary px-0" @click="emit('cancelMfa')">
            <i class="ti ti-arrow-left me-1"></i>Ganti Akun
          </button>
          <button type="button" class="btn btn-link btn-sm text-primary px-0" @click="emit('switchToRecovery')">
            Gunakan Recovery Code
          </button>
        </div>
      </form>
    </div>

    <!-- TAHAP 2C: MASUK MENGGUNAKAN RECOVERY CODE -->
    <div v-else-if="mfaStage === 'recovery'">
      <div class="text-center mb-3">
        <div class="avatar avatar-md bg-light-warning text-warning mx-auto mb-2 rounded-circle p-2 d-inline-flex align-items-center justify-content-center" style="width: 48px; height: 48px;">
          <i class="ti ti-key fs-3"></i>
        </div>
        <h4 class="mb-1"><b>Masuk dengan Recovery Code</b></h4>
        <p class="text-muted small">Gunakan salah satu dari 10 kode pemulihan cadangan yang pernah Anda simpan.</p>
      </div>

      <div v-if="errorMsg" class="alert alert-danger py-2 px-3 mb-3 d-flex align-items-center gap-2" role="alert">
        <i class="ti ti-alert-circle"></i>
        <span>{{ errorMsg }}</span>
      </div>

      <form @submit.prevent="emit('verifyMfa')">
        <div class="form-group mb-3">
          <label class="form-label">Recovery Code (Contoh: XXXXX-XXXXX-XXXXX-XXXXX)</label>
          <input
            :value="mfaCode"
            type="text"
            class="form-control font-monospace text-center text-uppercase"
            placeholder="XXXXX-XXXXX-XXXXX-XXXXX"
            autocomplete="off"
            autofocus
            @input="emit('update:mfaCode', ($event.target as HTMLInputElement).value)"
          >
        </div>

        <div class="d-grid mb-3">
          <button type="submit" class="btn btn-warning" :disabled="loading || !mfaCode">
            <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            Gunakan Kode Pemulihan
          </button>
        </div>

        <div class="d-flex justify-content-between align-items-center">
          <button type="button" class="btn btn-link btn-sm text-secondary px-0" @click="emit('cancelMfa')">
            <i class="ti ti-arrow-left me-1"></i>Batal
          </button>
          <button type="button" class="btn btn-link btn-sm text-primary px-0" @click="emit('switchToVerify')">
            Gunakan Kode Authenticator
          </button>
        </div>
      </form>
    </div>

    <!-- TAHAP 3: TAMPILKAN RECOVERY CODES PASCA SETUP SUKSES -->
    <div v-else-if="mfaStage === 'recovery_display'">
      <div class="text-center mb-3">
        <div class="avatar avatar-md bg-light-success text-success mx-auto mb-2 rounded-circle p-2 d-inline-flex align-items-center justify-content-center" style="width: 48px; height: 48px;">
          <i class="ti ti-check fs-3"></i>
        </div>
        <h4 class="mb-1"><b>MFA Berhasil Diaktifkan!</b></h4>
        <p class="text-muted small">
          Simpan 10 recovery codes ini di tempat aman. Setiap kode hanya dapat digunakan 1 kali apabila Anda kehilangan perangkat authenticator.
        </p>
      </div>

      <div class="p-3 bg-light rounded border mb-3">
        <div class="row g-2 font-monospace small">
          <div v-for="(code, idx) in recoveryCodes" :key="idx" class="col-6 text-center py-1">
            <span class="badge bg-white text-dark border px-2 py-1 w-100">{{ code }}</span>
          </div>
        </div>
      </div>

      <div class="d-grid gap-2 mb-2">
        <button type="button" class="btn btn-outline-primary" @click="emit('downloadRecovery')">
          <i class="ti ti-download me-1"></i>Unduh Kode (.txt)
        </button>
        <button type="button" class="btn btn-primary" @click="emit('finishSetup')">
          Lanjutkan ke Dashboard <i class="ti ti-arrow-right ms-1"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  username: string
  password: string
  remember: boolean
  loading: boolean
  errorMsg: string
  mfaStage: 'credentials' | 'setup' | 'verify' | 'recovery' | 'recovery_display'
  mfaCode: string
  mfaSecret: string
  mfaOtpauthUri: string
  recoveryCodes: string[]
}

defineProps<Props>()

const showPassword = ref(false)

const emit = defineEmits<{
  'update:username': [value: string]
  'update:password': [value: string]
  'update:remember': [value: boolean]
  'update:mfaCode': [value: string]
  submit: []
  verifyMfa: []
  cancelMfa: []
  switchToRecovery: []
  switchToVerify: []
  downloadRecovery: []
  finishSetup: []
}>()
</script>

<style scoped>
.letter-spacing-2 {
  letter-spacing: 0.25rem;
}
</style>
