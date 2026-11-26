/**
 * Dinonaktifkan sengaja.
 *
 * Sebelumnya plugin ini inject/remove CSS saat runtime berdasarkan route.
 * Di SPA mode hal ini memicu FOUC (halaman sempat tampil tanpa CSS beberapa ms).
 *
 * Sekarang CSS dimuat langsung via `useHead` pada layout/page terkait:
 * - admin/login: layout admin & admin-login
 * - pelapor: halaman landing
 */
export default defineNuxtPlugin(() => {})
