// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  devServer: {
    host: '0.0.0.0',
    port: 3000,
  },

  /**
   * @nuxtjs/tailwindcss diaktifkan dengan cssPath: false agar TIDAK
   * meng-inject Tailwind CSS secara global.
   *
   * CSS template (pelapor & admin) dimuat di app.head.link agar
   * tersedia di HTML statis sejak awal → tidak FOUC di SPA mode.
   * Layout hanya override favicon saat navigasi ke halaman admin.
   */
  modules: [
    ['@nuxtjs/tailwindcss', { cssPath: false }],
  ],

  css: [],

  app: {
    head: {
      title: 'Arsip Record Center | KemenPPPA',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0, user-scalable=0, minimal-ui' },
        { name: 'description', content: 'Arsip Record Center — Kementerian Pemberdayaan Perempuan dan Perlindungan Anak Republik Indonesia' },
        { name: 'keywords', content: 'Arsip Record Center, Arsip, Record Center, KemenPPPA, Kearsipan' },
      ],
      link: [
        // Favicon
        { rel: 'icon', type: 'image/png', href: '/admin/images/logokpppa.png' },
        // Google Fonts admin
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Public+Sans:wght@300;400;500;600;700&display=swap' },
        // Admin CSS
        { rel: 'stylesheet', href: '/admin/fonts/tabler-icons.min.css' },
        { rel: 'stylesheet', href: '/admin/fonts/feather.css' },
        { rel: 'stylesheet', href: '/admin/fonts/fontawesome.css' },
        { rel: 'stylesheet', href: '/admin/fonts/material.css' },
        { rel: 'stylesheet', href: '/admin/css/style.css' },
        { rel: 'stylesheet', href: '/admin/css/style-preset.css' },
        { rel: 'stylesheet', href: '/admin/css/plugins/dataTables.bootstrap5.min.css' },
      ],
      script: [
        { src: 'https://cdn.jsdelivr.net/npm/chart.js', defer: true },
      ],
    },
  },

  $production: {
    devtools: { enabled: false },
  },

  vite: {
    build: {
      cssCodeSplit: true,
      rollupOptions: {
        external: (id: string) => id.startsWith('/admin/'),
      },
    },
    assetsInclude: ['**/*.png', '**/*.jpg', '**/*.jpeg', '**/*.gif', '**/*.svg', '**/*.ico', '**/*.webp'],
  },

  ssr: false,

  components: [
    { path: '~/components' },
    { path: '~/modules', prefix: 'Modules' },
  ],
})
