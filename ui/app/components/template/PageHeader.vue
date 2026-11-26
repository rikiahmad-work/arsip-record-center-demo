<template>
  <div class="page-header">
    <div class="page-block">
      <div class="row align-items-center">
        <!-- Breadcrumb -->
        <div class="col-md-12">
          <ul class="breadcrumb">
            <li class="breadcrumb-item">
              <NuxtLink to="/dashboard">Dashboard</NuxtLink>
            </li>
            <li
              v-for="(crumb, idx) in breadcrumbs"
              :key="idx"
              class="breadcrumb-item"
              :class="{ active: idx === breadcrumbs.length - 1 }"
              :aria-current="idx === breadcrumbs.length - 1 ? 'page' : undefined"
            >
              <NuxtLink v-if="crumb.to && idx < breadcrumbs.length - 1" :to="crumb.to">
                {{ crumb.label }}
              </NuxtLink>
              <span v-else>{{ crumb.label }}</span>
            </li>
          </ul>
        </div>

        <!-- Title + Action slot -->
        <div class="col-md-12">
          <div class="page-header-title d-flex flex-wrap justify-content-between align-items-center gap-2">
            <h2 class="mb-0">{{ title }}</h2>
            <div class="page-header-actions d-flex flex-wrap align-items-center gap-2">
              <slot name="action" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Breadcrumb {
  label: string
  to?: string
}

interface Props {
  title: string
  breadcrumbs?: Breadcrumb[]
}

withDefaults(defineProps<Props>(), {
  breadcrumbs: () => [],
})
</script>
