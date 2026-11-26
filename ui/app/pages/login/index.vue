<template>
  <!-- Mengikuti template docs/template-ui/UI-admin/index.html -->
  <div class="auth-main">
    <div class="auth-wrapper v3">
      <div class="auth-form">
        <div class="logo-container text-center">
          <img
            src="/admin/images/logokpppa.png"
            alt="Logo Kemen PPPA"
            style="width: 200px; height: 200px; object-fit: contain;"
          >
        </div>
        <h2 style="color: white; text-align: center; font-weight: 700;">ARSIP RECORD CENTER KEMENPPPA</h2>

        <div class="card my-5">
          <div class="card-body">
            <CAuthLoginForm
              :username="username"
              :password="password"
              :remember="remember"
              :loading="loading"
              :error-msg="errorMsg"
              :mfa-stage="mfaStage"
              :mfa-code="mfaCode"
              :mfa-secret="mfaSecret"
              :mfa-otpauth-uri="mfaOtpauthUri"
              :recovery-codes="recoveryCodes"
              @update:username="username = $event"
              @update:password="password = $event"
              @update:remember="remember = $event"
              @update:mfa-code="mfaCode = $event"
              @submit="handleLogin"
              @verify-mfa="handleVerifyMfa"
              @cancel-mfa="handleCancelMfa"
              @switch-to-recovery="mfaStage = 'recovery'; mfaCode = ''; errorMsg = ''"
              @switch-to-verify="mfaStage = 'verify'; mfaCode = ''; errorMsg = ''"
              @download-recovery="downloadRecoveryCodesTxt"
              @finish-setup="proceedToDashboard"
            />
          </div>
        </div>

        <div class="auth-footer row">
          <div class="col my-1 text-center">
            <p class="m-0" style="color: white;">
              Copyright &copy; {{ currentYear }} All Right Reserved |
              <a href="https://www.kemenpppa.go.id" target="_blank" rel="noopener" style="color: #fff; text-decoration: underline;">
                Kemen PPPA
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'login' })

useHead({ title: 'Login | Arsip Record Center KemenPPPA' })

const router = useRouter()
const currentYear = new Date().getFullYear()

const username = ref('')
const password = ref('')
const remember = ref(true)
const loading  = ref(false)
const errorMsg = ref('')

// State Multi-Tahap MFA
const mfaStage = ref<'credentials' | 'setup' | 'verify' | 'recovery' | 'recovery_display'>('credentials')
const mfaChallenge = ref('')
const mfaSecret = ref('')
const mfaOtpauthUri = ref('')
const mfaCode = ref('')
const recoveryCodes = ref<string[]>([])

const { login, verifyMfa, isLoggedIn, getDefaultRoute } = useSession()

async function handleLogin() {
  errorMsg.value = ''
  loading.value = true

  if (!username.value || !password.value) {
    errorMsg.value = 'NIP/Username dan password harus diisi'
    loading.value = false
    return
  }

  const result = await login(username.value, password.value, remember.value)

  if (result.success) {
    if (remember.value) {
      localStorage.setItem('remember_username', username.value)
    } else {
      localStorage.removeItem('remember_username')
    }

    if (result.mfaRequired) {
      loading.value = false
      mfaChallenge.value = result.challenge || ''
      mfaSecret.value = result.secret || ''
      mfaOtpauthUri.value = result.otpauthUri || ''
      mfaCode.value = ''

      if (result.mode === 'setup') {
        mfaStage.value = 'setup'
      } else {
        mfaStage.value = 'verify'
      }
      return
    }

    loading.value = false
    await proceedToDashboard()
  } else {
    errorMsg.value = result.message
    loading.value = false
  }
}

async function handleVerifyMfa() {
  errorMsg.value = ''
  loading.value = true

  if (!mfaCode.value) {
    errorMsg.value = 'Kode verifikasi harus diisi'
    loading.value = false
    return
  }

  const result = await verifyMfa(mfaChallenge.value, mfaCode.value.trim())

  if (result.success) {
    loading.value = false
    // Jika ada recovery codes yang dikembalikan (mode setup selesai)
    if (result.recoveryCodes && result.recoveryCodes.length > 0) {
      recoveryCodes.value = result.recoveryCodes
      mfaStage.value = 'recovery_display'
    } else {
      await proceedToDashboard()
    }
  } else {
    errorMsg.value = result.message
    loading.value = false
  }
}

function handleCancelMfa() {
  mfaStage.value = 'credentials'
  mfaChallenge.value = ''
  mfaSecret.value = ''
  mfaOtpauthUri.value = ''
  mfaCode.value = ''
  errorMsg.value = ''
}

function downloadRecoveryCodesTxt() {
  if (!recoveryCodes.value.length) return
  const content = [
    '=======================================================',
    ' SiARSIP KemenPPPA — Multi-Factor Recovery Codes',
    ` Dibuat Pada: ${formatDateTimeIndo(new Date())}`,
    ` Akun: ${username.value}`,
    '=======================================================',
    '',
    'PERINGATAN:',
    '- Setiap kode pemulihan di bawah ini HANYA DAPAT DIGUNAKAN 1 KALI.',
    '- Simpan file ini di tempat yang aman dan rahasia.',
    '- Jika Anda kehilangan ponsel / authenticator, gunakan salah satu kode ini.',
    '',
    ...recoveryCodes.value.map((c, i) => `[${(i + 1).toString().padStart(2, '0')}] ${c}`),
    '',
    '=======================================================',
  ].join('\n')

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `siarsip-recovery-codes-${username.value || 'akun'}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

async function proceedToDashboard() {
  const target = getDefaultRoute()
  const dest = target === '/login' ? '/dashboard' : target
  await navigateTo(dest)
}

onMounted(async () => {
  if (isLoggedIn()) {
    const target = getDefaultRoute()
    const dest = target === '/login' ? '/dashboard' : target
    await navigateTo(dest)
  }

  const remembered = localStorage.getItem('remember_username')
  if (remembered) {
    username.value = remembered
  }
})
</script>
