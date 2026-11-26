import type { AdminRoleKey } from '~/types/index'
import { normalizeRoleKey } from '~/constants/auth'

export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return

  const { getRoleKey, getDefaultRoute } = useSession()
  const role = getRoleKey()

  if (!role) {
    return navigateTo('/login')
  }

  const rawAllowed = (to.meta.allowedRoles as AdminRoleKey[] | undefined) || []
  if (rawAllowed.length > 0) {
    const normalizedAllowed = rawAllowed.map((r) => normalizeRoleKey(r))
    if (!normalizedAllowed.includes(normalizeRoleKey(role))) {
      return navigateTo(getDefaultRoute())
    }
  }
})
