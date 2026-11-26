export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook('app:chunkError', () => {
    // Otomatis reload jika ada mismatch chunk setelah deploy baru
    window.location.reload()
  })
})
