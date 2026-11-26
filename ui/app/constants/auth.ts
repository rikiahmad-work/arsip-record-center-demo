import type { AdminRoleKey } from '~/types/index'

export const ADMIN_TOKEN_KEY = 'admin_token'
export const ADMIN_USER_KEY = 'admin_user'

// API Backend URL — dari .env (NUXT_PUBLIC_API_BASE_URL) atau fallback default port 8000
declare const __NUXT_PUBLIC_API_BASE_URL__: string | undefined

export function normalizeApiUrl(raw?: string): string {
  let url = (raw || '').trim()
  if (!url) {
    url = 'http://localhost:8000/api'
  }
  return url.endsWith('/api') ? url : `${url.replace(/\/+$/, '')}/api`
}

function resolveRawApiUrl(): string {
  if (typeof __NUXT_PUBLIC_API_BASE_URL__ !== 'undefined' && __NUXT_PUBLIC_API_BASE_URL__) {
    return __NUXT_PUBLIC_API_BASE_URL__
  }
  if (import.meta.env?.NUXT_PUBLIC_API_BASE_URL) {
    return import.meta.env.NUXT_PUBLIC_API_BASE_URL as string
  }
  if (typeof window !== 'undefined' && (window as any).__NUXT__?.config?.public?.apiBase) {
    return (window as any).__NUXT__.config.public.apiBase as string
  }
  if (typeof process !== 'undefined' && process.env?.NUXT_PUBLIC_API_BASE_URL) {
    return process.env.NUXT_PUBLIC_API_BASE_URL
  }
  return ''
}

export const API_BASE_URL = normalizeApiUrl(resolveRawApiUrl())

// Pemetaan role lama ke 4 role kanonik
export const LEGACY_ROLE_MAP: Record<string, AdminRoleKey> = {
  'super-admin': 'administrator',
  'admin-arsip': 'operator',
}

export function normalizeRoleKey(raw?: string | null): AdminRoleKey {
  if (!raw) return 'viewer'
  const lower = raw.toLowerCase().trim()
  if (LEGACY_ROLE_MAP[lower]) {
    return LEGACY_ROLE_MAP[lower]
  }
  if (['administrator', 'operator', 'viewer', 'approver'].includes(lower)) {
    return lower as AdminRoleKey
  }
  return 'viewer'
}

export const ROLE_LABELS: Record<string, string> = {
  administrator: 'Administrator',
  operator: 'Operator',
  viewer: 'Viewer',
  approver: 'Approver',
  'super-admin': 'Administrator',
  'admin-arsip': 'Operator',
}
