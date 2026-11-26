import { ADMIN_MENU_ITEMS } from '~/constants/menu'
import { ADMIN_TOKEN_KEY, ADMIN_USER_KEY, API_BASE_URL, normalizeRoleKey, ROLE_LABELS } from '~/constants/auth'
import type { AdminRoleKey, AdminUserSession } from '~/types/index'

function resolveUserPermissions(user: AdminUserSession | null): string[] {
  if (!user) return []
  const norm = normalizeRoleKey(user.roleKey)
  if (norm === 'administrator' || norm === 'super-admin') {
    return [
      'dashboard.view', 'archives.view', 'archives.create', 'archives.edit', 'archives.delete', 'archives.import', 'archives.export',
      'transfers.view', 'transfers.create', 'transfers.verify', 'storage.view', 'storage.manage',
      'loans.view', 'loans.create', 'loans.approve', 'loans.return', 'disposal.view', 'disposal.propose', 'disposal.approve', 'disposal.execute',
      'reports.view', 'reports.export', 'users.view', 'users.create', 'users.edit', 'users.delete', 'roles.manage', 'settings.view', 'settings.manage',
    ]
  }
  // Gunakan permissions aktual dari user session
  if (Array.isArray(user.permissions)) {
    return user.permissions
  }
  return []
}

export function useSession() {
  const sessionUser = useState<AdminUserSession | null>('admin_session_user', () => null)

  function getToken(): string {
    if (import.meta.server) return ''
    return localStorage.getItem(ADMIN_TOKEN_KEY) || ''
  }

  function getUser(): AdminUserSession | null {
    if (sessionUser.value) {
      return sessionUser.value
    }
    if (import.meta.server) return null
    const raw = localStorage.getItem(ADMIN_USER_KEY)
    if (!raw) return null

    try {
      const parsed = JSON.parse(raw) as Partial<AdminUserSession>
      if (!parsed.name || !parsed.roleKey) return null
      const normRole = normalizeRoleKey(parsed.roleKey)
      const u: AdminUserSession = {
        id: parsed.id,
        username: parsed.username || '',
        name: parsed.name,
        nip: parsed.nip || '',
        email: parsed.email || '',
        roleKey: normRole,
        roleLabel: parsed.roleLabel || ROLE_LABELS[normRole] || normRole,
        unit: parsed.unit || '-',
        jabatan: parsed.jabatan || '',
        isSso: parsed.isSso || false,
        permissions: Array.isArray(parsed.permissions) ? parsed.permissions : [],
      }
      sessionUser.value = u
      return u
    } catch {
      return null
    }
  }

  function isLoggedIn(): boolean {
    return Boolean(getToken())
  }

  function getRoleKey(): AdminRoleKey | null {
    const user = getUser()
    return user ? normalizeRoleKey(user.roleKey) : null
  }

  function hasRole(allowedRoles: AdminRoleKey[]): boolean {
    const role = getRoleKey()
    if (!role) return false
    const normalizedAllowed = allowedRoles.map((r) => normalizeRoleKey(r))
    return normalizedAllowed.includes(role)
  }

  function canWrite(): boolean {
    const role = getRoleKey()
    return role === 'administrator' || role === 'operator'
  }

  function canApprove(): boolean {
    const role = getRoleKey()
    return role === 'administrator' || role === 'approver'
  }

  function isAdmin(): boolean {
    const role = getRoleKey()
    return role === 'administrator' || role === 'super-admin'
  }

  function hasPermission(permission: string): boolean {
    const user = getUser()
    if (!user) return false
    const norm = normalizeRoleKey(user.roleKey)
    // Administrator & Super-Admin memiliki hak akses penuh
    if (norm === 'administrator' || norm === 'super-admin') {
      return true
    }
    const perms = resolveUserPermissions(user)
    return perms.includes(permission)
  }

  function getMenuByRole(roleKey: AdminRoleKey | null) {
    if (!roleKey) return []
    const norm = normalizeRoleKey(roleKey)
    const user = getUser()
    const isAdm = norm === 'administrator' || norm === 'super-admin'
    const perms = resolveUserPermissions(user)

    return ADMIN_MENU_ITEMS.filter((item) => {
      const roleAllowed = item.roles.map((r) => normalizeRoleKey(r)).includes(norm)
      if (!roleAllowed) return false

      if (isAdm) return true

      if (item.permission) {
        return perms.includes(item.permission)
      }

      return true
    })
  }

  function getDefaultRoute(): string {
    const user = getUser()
    if (!user) return '/login'
    if (hasPermission('dashboard.view')) {
      return '/dashboard'
    }
    const menus = getMenuByRole(user.roleKey)
    if (menus.length > 0 && menus[0]?.to) {
      return menus[0].to
    }
    return '/dashboard'
  }

  function setSession(token: string, user: AdminUserSession) {
    if (import.meta.server) return
    localStorage.setItem(ADMIN_TOKEN_KEY, token)
    localStorage.setItem(ADMIN_USER_KEY, JSON.stringify(user))
    sessionUser.value = user
  }

  function clearSession() {
    if (import.meta.server) return
    localStorage.removeItem(ADMIN_TOKEN_KEY)
    localStorage.removeItem(ADMIN_USER_KEY)
    sessionUser.value = null
  }

  /** Login ke API backend (mendukung NIP maupun Username, serta MFA challenge) */
  async function login(
    identifier: string,
    password: string,
    rememberMe: boolean
  ): Promise<{
    success: boolean
    message: string
    mfaRequired?: boolean
    challenge?: string
    mode?: string
    secret?: string
    otpauthUri?: string
  }> {
    try {
      const trimmed = identifier.trim()
      const res = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nip: trimmed,
          username: trimmed,
          password,
          remember_me: rememberMe,
        }),
      })

      const json = await res.json()

      if (!res.ok || !json.success) {
        return { success: false, message: json.message || 'NIP/Username atau password salah' }
      }

      const d = json.data
      // Jika butuh MFA (setup atau verify)
      if (d.mfa_required) {
        return {
          success: true,
          message: json.message || 'Verifikasi MFA diperlukan',
          mfaRequired: true,
          challenge: d.challenge,
          mode: d.mode,
          secret: d.secret,
          otpauthUri: d.otpauth_uri,
        }
      }

      const normRole = normalizeRoleKey(d.user.role)
      const user: AdminUserSession = {
        id: d.user.id,
        username: d.user.username || d.user.nip,
        name: d.user.name,
        nip: d.user.nip,
        email: d.user.email,
        roleKey: normRole,
        roleLabel: d.user.role_label || ROLE_LABELS[normRole] || normRole,
        unit: d.user.unit_kerja || '-',
        jabatan: d.user.jabatan || '',
        isSso: d.user.is_sso || false,
        permissions: Array.isArray(d.user.permissions) ? d.user.permissions : [],
      }

      setSession(d.token, user)

      return { success: true, message: 'Login berhasil' }
    } catch {
      return { success: false, message: 'Tidak dapat terhubung ke server API' }
    }
  }

  /** Verifikasi kode MFA (TOTP atau Recovery Code) */
  async function verifyMfa(
    challenge: string,
    code: string
  ): Promise<{
    success: boolean
    message: string
    recoveryCodes?: string[]
  }> {
    try {
      const res = await fetch(`${API_BASE_URL}/auth/mfa/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          challenge,
          code,
        }),
      })

      const json = await res.json()
      if (!res.ok || !json.success) {
        return { success: false, message: json.message || 'Kode verifikasi MFA salah' }
      }

      const d = json.data
      const normRole = normalizeRoleKey(d.user.role)
      const user: AdminUserSession = {
        id: d.user.id,
        username: d.user.username || d.user.nip,
        name: d.user.name,
        nip: d.user.nip,
        email: d.user.email,
        roleKey: normRole,
        roleLabel: d.user.role_label || ROLE_LABELS[normRole] || normRole,
        unit: d.user.unit_kerja || '-',
        jabatan: d.user.jabatan || '',
        isSso: d.user.is_sso || false,
        permissions: Array.isArray(d.user.permissions) ? d.user.permissions : [],
      }

      setSession(d.token, user)

      return {
        success: true,
        message: 'Verifikasi MFA berhasil',
        recoveryCodes: d.recovery_codes || [],
      }
    } catch {
      return { success: false, message: 'Gagal menghubungi server verifikasi MFA' }
    }
  }

  /** Logout: hapus session & blacklist JTI di backend */
  async function logout(): Promise<void> {
    const token = getToken()
    if (token) {
      try {
        await fetch(`${API_BASE_URL}/auth/logout`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` },
        })
      } catch {
        // ignore network error on logout
      }
    }
    clearSession()
  }

  /**
   * apiFetch: wrapper fetch dengan Authorization header + 401 auto-redirect.
   * Mendukung baik JSON request maupun FormData (multipart).
   */
  async function apiFetch(path: string, options: RequestInit = {}): Promise<Response> {
    const token = getToken()
    const isFormData = typeof FormData !== 'undefined' && options.body instanceof FormData
    const headers: Record<string, string> = {
      ...(options.headers as Record<string, string>),
      Authorization: `Bearer ${token}`,
    }
    if (!isFormData && !headers['Content-Type']) {
      headers['Content-Type'] = 'application/json'
    }

    const res = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      headers,
    })

    if (res.status === 401 || res.status === 403) {
      const clone = res.clone()
      try {
        const body = await clone.json()
        if (body.message && body.message.toLowerCase().includes('dinonaktifkan')) {
          clearSession()
          if (import.meta.client) {
            alert('Akun Anda telah dinonaktifkan. Silakan hubungi administrator.')
            await navigateTo('/login')
          }
          throw new Error(body.message)
        }
      } catch (_) {}

      if (res.status === 401) {
        clearSession()
        if (import.meta.client) {
          await navigateTo('/login')
        }
        throw new Error('Sesi telah berakhir')
      }
    }

    return res
  }

  /**
   * refreshProfile: Mengambil profil terbaru dari backend dan menyinkronkan
   * hak akses (permissions) ke session tanpa mengharuskan pengguna login ulang.
   */
  async function refreshProfile(): Promise<AdminUserSession | null> {
    const token = getToken()
    if (!token) return null
    try {
      const res = await fetch(`${API_BASE_URL}/auth/profile`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) return null
      const json = await res.json()
      if (!json.success || !json.data) return null
      const u = json.data
      const normRole = normalizeRoleKey(u.role)
      const user: AdminUserSession = {
        id: u.id,
        username: u.username || u.nip,
        name: u.name,
        nip: u.nip,
        email: u.email,
        roleKey: normRole,
        roleLabel: u.role_label || ROLE_LABELS[normRole] || normRole,
        unit: u.unit_kerja || '-',
        jabatan: u.jabatan || '',
        isSso: u.is_sso || false,
        permissions: Array.isArray(u.permissions) ? u.permissions : [],
      }
      setSession(token, user)
      return user
    } catch {
      return null
    }
  }

  return {
    getToken,
    getUser,
    isLoggedIn,
    getRoleKey,
    hasRole,
    canWrite,
    canApprove,
    isAdmin,
    hasPermission,
    getMenuByRole,
    getDefaultRoute,
    refreshProfile,
    login,
    verifyMfa,
    logout,
    clearSession,
    apiFetch,
  }
}
