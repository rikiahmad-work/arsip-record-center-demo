<script setup lang="ts">
// Root redirect: jika sudah login → dashboard, jika belum → login
definePageMeta({ layout: false })

const router = useRouter()

onMounted(async () => {
  const { isLoggedIn, getDefaultRoute } = useSession()
  if (isLoggedIn()) {
    const target = getDefaultRoute()
    const dest = target === '/login' ? '/dashboard' : target
    await navigateTo(dest)
  } else {
    await navigateTo('/login')
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <p class="text-gray-500">Mengalihkan...</p>
  </div>
</template>
