export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return

  const { isLoggedIn, hasPermission, getDefaultRoute, isAdmin } = useSession()

  if (!isLoggedIn()) {
    return navigateTo('/login')
  }

  // Administrator & Super-Admin selalu memiliki akses penuh (standar industri)
  if (isAdmin()) {
    return
  }

  // Pemetaan rute ke permission yang dipersyaratkan sesuai Matriks RBAC
  const routePermissions: Record<string, string> = {
    '/dashboard': 'dashboard.view',
    '/pengelolaan-arsip': 'archives.view',
    '/peminjaman-arsip': 'loans.view',
    '/penyusutan-arsip': 'disposal.view',
    '/laporan': 'reports.view',
    '/label': 'reports.export',
    '/master/user': 'roles.manage',
    '/master/pengaturan': 'settings.view',
  }

  // Sub-halaman pengelolaan arsip
  if (to.path === '/pengelolaan-arsip/tambah') {
    if (!hasPermission('archives.create')) {
      return navigateTo('/pengelolaan-arsip')
    }
    return
  }
  if (to.path.endsWith('/edit') && to.path.startsWith('/pengelolaan-arsip/')) {
    if (!hasPermission('archives.edit')) {
      return navigateTo('/pengelolaan-arsip')
    }
    return
  }

  // Cek kecocokan rute utama
  for (const [routePath, requiredPerm] of Object.entries(routePermissions)) {
    if (to.path === routePath || to.path.startsWith(routePath + '/')) {
      if (!hasPermission(requiredPerm)) {
        const fallback = getDefaultRoute()
        if (fallback === to.path) {
          return navigateTo('/login')
        }
        return navigateTo(fallback)
      }
      break
    }
  }
})
