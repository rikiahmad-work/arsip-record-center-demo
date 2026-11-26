<template>
  <div class="card">
    <div v-if="title" class="card-header border-bottom-0 pb-0">
      <h5 class="mb-0">{{ title }}</h5>
      <slot name="header-action" />
    </div>
    <div class="card-body">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
}

defineProps<Props>()
</script>
