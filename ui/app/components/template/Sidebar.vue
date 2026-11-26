<template>
  <nav class="pc-sidebar">
    <div class="navbar-wrapper">
      <!-- Logo sesuai docs/template-ui/UI-admin/dashboard/index.html -->
      <div class="m-header">
        <NuxtLink :to="logoTarget" class="b-brand text-primary" style="text-decoration: none; display: flex; align-items: center; gap: 10px;" @click="closeMobileSidebar">
          <img src="/admin/images/logokpppa.png" alt="Logo KemenPPPA" style="width: 40px;">
          <div style="display: flex; flex-direction: column; line-height: 1.15;">
            <span style="font-weight: bold; font-size: 1.2rem; color: #1890ff;">ARSIP</span>
            <span style="font-weight: 600; font-size: 0.75rem; color: #262626; letter-spacing: 0.5px;">RECORD CENTER</span>
          </div>
        </NuxtLink>
      </div>

      <!-- Nav Items -->
      <div class="navbar-content">
        <ul class="pc-navbar">
          <template v-for="(items, groupName) in groupedMenu" :key="groupName">
            <li v-if="groupName !== 'Utama' && items.length > 0" class="pc-item pc-caption">
              <label>{{ groupName }}</label>
              <i class="ti ti-dashboard"></i>
            </li>

            <li v-for="item in items" :key="item.key" class="pc-item" :class="{ active: isActive(item.to) }">
              <NuxtLink :to="item.to" class="pc-link" @click="closeMobileSidebar">
                <span class="pc-micon"><i :class="item.icon"></i></span>
                <span class="pc-mtext">{{ item.label }}</span>
              </NuxtLink>
            </li>
          </template>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import type { AdminMenuItem } from '~/types/index'

const route = useRoute()
const { getUser, getMenuByRole, getDefaultRoute } = useSession()

const logoTarget = computed(() => getDefaultRoute())

const menuItems = computed<AdminMenuItem[]>(() => {
  // Mengikuti reaktivitas rute dan session pengguna saat ini
  const _ = route.path
  const user = getUser()
  return getMenuByRole(user ? user.roleKey : null)
})

const groupedMenu = computed<Record<string, AdminMenuItem[]>>(() => {
  return menuItems.value.reduce((acc, item) => {
    const group = item.group || 'Menu'
    if (!acc[group]) acc[group] = []
    acc[group].push(item)
    return acc
  }, {} as Record<string, AdminMenuItem[]>)
})

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

function closeMobileSidebar() {
  if (process.client) {
    document.querySelector('.pc-sidebar')?.classList.remove('mob-sidebar-active')
    document.querySelector('.pc-menu-overlay')?.remove()
  }
}

watch(() => route.path, () => {
  closeMobileSidebar()
})
</script>
