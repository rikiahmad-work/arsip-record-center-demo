<template>
  <TemplatePageHeader :title="title" :breadcrumbs="breadcrumbs">
    <template #action>
      <slot name="action" />
    </template>
  </TemplatePageHeader>
</template>

<script setup lang="ts">
import type { Breadcrumb } from '~/components/template/PageHeader.vue'

interface Props {
  title: string
  breadcrumbs?: Breadcrumb[]
}

withDefaults(defineProps<Props>(), {
  breadcrumbs: () => [],
})
</script>
