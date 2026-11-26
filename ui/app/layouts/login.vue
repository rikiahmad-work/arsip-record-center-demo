<template>
  <!-- Pre-loader sesuai docs/template-ui/UI-admin/index.html -->
  <div class="loader-bg">
    <div class="loader-track">
      <div class="loader-fill"></div>
    </div>
  </div>

  <slot />
</template>

<script setup lang="ts">
useHead({
  link: [
    { rel: 'icon', type: 'image/png', href: '/admin/images/logokpppa.png' },
  ],
})

onMounted(() => {
  document.body.setAttribute('data-pc-preset', 'preset-1')
  document.body.setAttribute('data-pc-direction', 'ltr')
  document.body.setAttribute('data-pc-theme', 'light')
  document.body.classList.remove('index-page')

  const win = window as any
  if (win.layout_change) win.layout_change('light')
  if (win.change_box_container) win.change_box_container('false')
  if (win.layout_rtl_change) win.layout_rtl_change('false')
  if (win.preset_change) win.preset_change('preset-1')
  if (win.font_change) win.font_change('Public-Sans')
})
</script>
