<template>
  <!-- Loader -->
  <div class="loader-bg">
    <div class="loader-track">
      <div class="loader-fill"></div>
    </div>
  </div>

  <TemplateSidebar />
  <TemplateHeader />

  <div class="pc-container">
    <div class="pc-content">
      <slot />
    </div>
  </div>

  <TemplateFooter />

  <!-- Toast Notification Global -->
  <UiAppToast />
</template>

<script setup lang="ts">
/**
 * CSS admin dimuat langsung via useHead agar tidak ada FOUC saat halaman pertama kali dirender.
 * JS Mantis tetap dimuat berurutan di onMounted.
 */
useHead({
  link: [
    { rel: 'icon', type: 'image/png', href: '/admin/images/logokpppa.png' },
  ],
})

onMounted(async () => {
  // Pastikan body attributes Mantis terset
  document.body.setAttribute('data-pc-preset', 'preset-1')
  document.body.setAttribute('data-pc-direction', 'ltr')
  document.body.setAttribute('data-pc-theme', 'light')
  // Hapus class pelapor jika ada
  document.body.classList.remove('index-page')

  const scripts = [
    '/admin/js/plugins/popper.min.js',
    '/admin/js/plugins/simplebar.min.js',
    '/admin/js/plugins/bootstrap.min.js',
    '/admin/js/fonts/custom-font.js',
    '/admin/js/pcoded.js',
    '/admin/js/plugins/feather.min.js',
    '/admin/js/plugins/apexcharts.min.js',
  ]

  const loadScript = (src: string) =>
    new Promise<void>((resolve) => {
      if (document.querySelector(`script[src="${src}"]`)) { resolve(); return }
      const s = document.createElement('script')
      s.src = src
      s.onload = () => resolve()
      s.onerror = () => resolve()
      document.body.appendChild(s)
    })

  await scripts.reduce((p, src) => p.then(() => loadScript(src)), Promise.resolve())

  // font_change() di pcoded.js butuh elemen <link id="main-font-link"> di <head>.
  // Inject sekali jika belum ada agar setAttribute tidak throw null error.
  if (!document.querySelector('#main-font-link')) {
    const link = document.createElement('link')
    link.id   = 'main-font-link'
    link.rel  = 'stylesheet'
    link.href = 'https://fonts.googleapis.com/css2?family=Public+Sans:wght@300;400;500;600;700&display=swap'
    document.head.appendChild(link)
  }

  const win = window as any
  if (win.layout_change)        win.layout_change('light')
  if (win.change_box_container) win.change_box_container('false')
  if (win.layout_rtl_change)    win.layout_rtl_change('false')
  if (win.preset_change)        win.preset_change('preset-1')
  if (win.font_change)          win.font_change('Public-Sans')
})
</script>

<style>
/* input-line-bottom — sesuai template UI-Admin laporan-WBS */
.input-line-bottom {
  border: none !important;
  border-bottom: 1px solid #ced4da !important;
  border-radius: 0 !important;
  padding-left: 0 !important;
  background-color: transparent !important;
  box-shadow: none !important;
}
.input-line-bottom:focus {
  border-bottom: 2px solid #5cb85c !important;
}
.input-line-bottom::placeholder {
  color: #999 !important;
}
</style>
