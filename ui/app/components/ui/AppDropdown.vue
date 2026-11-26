<template>
  <div class="dropdown">
    <button
      :class="['btn', `btn-${variant}`, size ? `btn-${size}` : '', 'dropdown-toggle']"
      type="button"
      data-bs-toggle="dropdown"
      aria-expanded="false"
    >
      <i v-if="icon" :class="[icon, 'me-1']"></i>
      {{ label }}
    </button>
    <ul class="dropdown-menu">
      <slot name="items" />
    </ul>
  </div>
</template>

<script setup lang="ts">
interface Props {
  label?: string
  variant?: string
  size?: 'sm' | 'lg' | ''
  icon?: string
}

withDefaults(defineProps<Props>(), {
  label: 'Aksi',
  variant: 'light-primary',
  size: 'sm',
})
</script>
