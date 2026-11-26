<template>
  <header class="pc-header">
    <div class="header-wrapper">

      <!-- Tombol toggle sidebar (kiri) -->
      <div class="me-auto pc-mob-drp">
        <ul class="list-unstyled">
          <li class="pc-h-item pc-sidebar-collapse">
            <a href="#" class="pc-head-link ms-0" id="sidebar-hide">
              <i class="ti ti-menu-2"></i>
            </a>
          </li>
          <li class="pc-h-item pc-sidebar-popup">
            <a href="#" class="pc-head-link ms-0" id="mobile-collapse">
              <i class="ti ti-menu-2"></i>
            </a>
          </li>
        </ul>
      </div>

      <!-- Actions (kanan) -->
      <div class="ms-auto">
        <ul class="list-unstyled">

          <!-- Notifikasi -->
          <li class="dropdown pc-h-item">
            <a
              class="pc-head-link dropdown-toggle arrow-none me-0"
              data-bs-toggle="dropdown"
              href="#"
              role="button"
              aria-haspopup="false"
              aria-expanded="false"
            >
              <i class="ti ti-bell"></i>
              <span v-if="notifCount > 0" class="badge bg-success pc-h-badge">{{ notifCount }}</span>
            </a>
            <div class="dropdown-menu dropdown-notification dropdown-menu-end pc-h-dropdown">
              <div class="dropdown-header d-flex align-items-center justify-content-between">
                <h5 class="m-0">Notifikasi</h5>
                <a
                  href="#!"
                  class="pc-head-link bg-transparent"
                  title="Tandai semua telah dibaca"
                  @click.prevent="clearNotifications"
                >
                  <i class="ti ti-circle-check text-success"></i>
                </a>
              </div>
              <div class="dropdown-divider"></div>
              <div
                class="dropdown-header px-0 text-wrap header-notification-scroll position-relative"
                style="max-height: calc(100vh - 215px)"
              >
                <div v-if="notifications.length" class="list-group list-group-flush w-100">
                  <div
                    v-for="notif in notifications"
                    :key="notif.id"
                    class="list-group-item list-group-item-action"
                  >
                    <div class="d-flex">
                      <div class="flex-shrink-0">
                        <div class="user-avtar bg-light-primary">
                          <i :class="notif.icon"></i>
                        </div>
                      </div>
                      <div class="flex-grow-1 ms-2">
                        <p class="text-body mb-0 small">{{ notif.message }}</p>
                        <span class="text-muted" style="font-size: 0.75rem;">{{ notif.time }}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <div v-else class="text-center py-4 text-muted small">
                  <i class="ti ti-bell-off fs-4 d-block mb-1 text-secondary opacity-50"></i>
                  Tidak ada notifikasi baru
                </div>
              </div>
              <div class="dropdown-divider"></div>
              <div class="text-center py-2">
                <a href="#!" class="link-primary">Lihat semua</a>
              </div>
            </div>
          </li>

          <!-- Profil User -->
          <li class="dropdown pc-h-item header-user-profile">
            <a
              class="pc-head-link dropdown-toggle arrow-none me-0 d-inline-flex align-items-center"
              data-bs-toggle="dropdown"
              href="#"
              role="button"
              aria-haspopup="false"
              data-bs-auto-close="outside"
              aria-expanded="false"
            >
              <UiAppAvatar :name="userName" :size="28" custom-class="user-avtar" />
              <span>{{ userName }}</span>
            </a>
            <div class="dropdown-menu dropdown-user-profile dropdown-menu-end pc-h-dropdown">
              <div class="dropdown-header">
                <div class="d-flex mb-1 align-items-center">
                  <div class="flex-shrink-0">
                    <UiAppAvatar :name="userName" :size="38" custom-class="user-avtar" />
                  </div>
                  <div class="flex-grow-1 ms-3">
                    <h6 class="mb-1">{{ userName }}</h6>
                    <span>{{ userRole }}</span>
                  </div>
                </div>
              </div>
              <div class="tab-content">
                <div class="tab-pane fade show active" role="tabpanel">
                  <NuxtLink to="/setting/mfa" class="dropdown-item">
                    <i class="ti ti-shield-lock"></i>
                    <span>Pengaturan MFA</span>
                  </NuxtLink>
                  <NuxtLink to="/setting/profile" class="dropdown-item">
                    <i class="ti ti-user"></i>
                    <span>Profil Saya</span>
                  </NuxtLink>
                  <a class="dropdown-item" href="#" @click.prevent="logout">
                    <i class="ti ti-power"></i>
                    <span>Keluar</span>
                  </a>
                </div>
              </div>
            </div>
          </li>

        </ul>
      </div>

    </div>
  </header>
</template>

<script setup lang="ts">
const router = useRouter()
const { getUser, logout: sessionLogout } = useSession()

// Data user dari localStorage (akan diganti dengan store/auth saat backend tersedia)
const userName = ref('Admin Record Center')
const userRole = ref('Administrator')

onMounted(() => {
  const user = getUser()
  if (user) {
    userName.value = user.name
    userRole.value = user.roleLabel
  }
})

interface NotificationItem {
  id: number
  icon: string
  message: string
  time: string
}

const notifications = ref<NotificationItem[]>([])

const notifCount = computed(() => notifications.value.length)

function clearNotifications() {
  notifications.value = []
}

function handleSidebarHideClick(event: Event) {
  event.preventDefault()
  const sidebar = document.querySelector('.pc-sidebar')
  if (!sidebar) return

  sidebar.classList.toggle('pc-sidebar-hide')
}

function handleMobileCollapseClick(event: Event) {
  event.preventDefault()
  const sidebar = document.querySelector('.pc-sidebar')
  if (!sidebar) return

  const isActive = sidebar.classList.contains('mob-sidebar-active')
  const existingOverlay = sidebar.querySelector('.pc-menu-overlay')

  if (isActive) {
    sidebar.classList.remove('mob-sidebar-active')
    if (existingOverlay) existingOverlay.remove()
    return
  }

  sidebar.classList.add('mob-sidebar-active')
  if (!existingOverlay) {
    sidebar.insertAdjacentHTML('beforeend', '<div class="pc-menu-overlay"></div>')
  }

  const overlay = sidebar.querySelector('.pc-menu-overlay')
  overlay?.addEventListener('click', () => {
    sidebar.classList.remove('mob-sidebar-active')
    overlay.remove()
  }, { once: true })
}

onMounted(() => {
  const btnHide = document.querySelector('#sidebar-hide')
  const btnMobile = document.querySelector('#mobile-collapse')

  btnHide?.addEventListener('click', handleSidebarHideClick)
  btnMobile?.addEventListener('click', handleMobileCollapseClick)
})

onBeforeUnmount(() => {
  const btnHide = document.querySelector('#sidebar-hide')
  const btnMobile = document.querySelector('#mobile-collapse')

  btnHide?.removeEventListener('click', handleSidebarHideClick)
  btnMobile?.removeEventListener('click', handleMobileCollapseClick)
})

async function logout() {
  await sessionLogout()
  router.push('/login')
}
</script>

<style scoped>
/* Pastikan tombol profil dan avatar tidak tertutup oleh pseudo-element ::before bawaan template Mantis */
:deep(.header-user-profile .pc-head-link) {
  position: relative;
  z-index: 2;
  overflow: visible;
  transition: background-color 0.15s ease;
}
:deep(.header-user-profile .pc-head-link::before) {
  display: none !important;
}
:deep(.header-user-profile .pc-head-link:hover) {
  background-color: rgba(0, 0, 0, 0.05) !important;
  border-radius: 8px;
}
:deep(.header-user-profile .pc-head-link .user-avtar),
:deep(.header-user-profile .pc-head-link span) {
  position: relative;
  z-index: 3;
}
</style>
