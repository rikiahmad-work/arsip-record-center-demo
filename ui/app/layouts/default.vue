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

/* ============================================================
   RESPONSIVE MOBILE ENHANCEMENTS (Desktop 100% Tidak Terpengaruh)
   ============================================================ */

/* Mobile Scrollable Tabs */
@media (max-width: 767.98px) {
  .nav-tabs {
    display: flex !important;
    flex-wrap: nowrap !important;
    overflow-x: auto !important;
    overflow-y: hidden !important;
    -webkit-overflow-scrolling: touch;
    white-space: nowrap !important;
    border-bottom: 1px solid #e2e5e8;
    padding-bottom: 4px;
    margin-bottom: 1rem !important;
    scrollbar-width: thin;
  }
  .nav-tabs::-webkit-scrollbar {
    height: 4px;
  }
  .nav-tabs::-webkit-scrollbar-thumb {
    background: #ced4da;
    border-radius: 4px;
  }
  .nav-tabs .nav-item {
    flex: 0 0 auto !important;
  }
  .nav-tabs .nav-link {
    padding: 0.5rem 0.85rem !important;
    font-size: 0.875rem !important;
  }

  /* Responsive DataTables Controls */
  .dataTables_wrapper .dataTables_length,
  .dataTables_wrapper .dataTables_filter {
    text-align: left !important;
    float: none !important;
    width: 100% !important;
    margin-bottom: 0.75rem !important;
  }
  .dataTables_wrapper .dataTables_filter label,
  .dataTables_wrapper .dataTables_length label {
    display: flex !important;
    align-items: center;
    width: 100%;
    gap: 0.5rem;
  }
  .dataTables_wrapper .dataTables_filter input {
    width: 100% !important;
    margin-left: 0 !important;
  }
  .dataTables_wrapper .dataTables_info {
    text-align: center !important;
    float: none !important;
    margin-bottom: 0.5rem !important;
    width: 100%;
  }
  .dataTables_wrapper .dataTables_paginate {
    text-align: center !important;
    float: none !important;
    display: flex !important;
    justify-content: center !important;
    margin-top: 0.5rem !important;
    width: 100%;
  }
  .dataTables_wrapper .dataTables_paginate .pagination {
    flex-wrap: wrap !important;
    justify-content: center !important;
    gap: 3px;
  }
}

/* Responsive Table Smooth Scrolling */
.table-responsive,
.user-table-wrapper,
.dt-responsive {
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
}

/* Extra Small Screens (Smartphones < 576px) */
@media (max-width: 575.98px) {
  .page-header {
    margin-bottom: 1rem !important;
  }
  .page-header-title h2 {
    font-size: 1.35rem !important;
  }
  .pc-container .pc-content {
    padding: 1rem 0.75rem !important;
  }
  .card-header {
    padding: 0.75rem 1rem !important;
  }
  .card-body {
    padding: 1rem 0.75rem !important;
  }
  .modal-dialog {
    margin: 0.5rem !important;
  }
  .modal-content {
    border-radius: 8px !important;
  }
}
</style>
