<template>
  <div class="d-flex align-items-center justify-content-center vh-100">
    <div class="text-center">
      <h1 class="display-1 text-muted mb-3">{{ error.statusCode }}</h1>
      <p class="fs-5 text-muted mb-4">{{ error.message || 'Halaman tidak ditemukan' }}</p>
      <button class="btn btn-primary" @click="handleError">Kembali ke Beranda</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const handleError = () => clearError({ redirect: '/' })
</script>
