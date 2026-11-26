<template>
  <span :class="['badge', colorClass]">{{ label }}</span>
</template>

<script setup lang="ts">
type BadgeStatus =
  | 'menunggu-verifikasi'
  | 'dalam-proses'
  | 'selesai'
  | 'ditolak'
  | 'aktif'
  | 'nonaktif'
  | string

interface Props {
  status: BadgeStatus
  label?: string
}

const props = defineProps<Props>()

const STATUS_MAP: Record<string, { color: string; text: string }> = {
  'menunggu-verifikasi': { color: 'bg-secondary', text: 'Menunggu Verifikasi' },
  'dalam-proses':        { color: 'bg-warning',   text: 'Dalam Proses' },
  selesai:               { color: 'bg-success',   text: 'Selesai' },
  ditolak:               { color: 'bg-danger',    text: 'Ditolak' },
  aktif:                 { color: 'bg-success',   text: 'Aktif' },
  nonaktif:              { color: 'bg-danger',    text: 'Nonaktif' },
}

const colorClass = computed(() => {
  return STATUS_MAP[props.status]?.color ?? 'bg-secondary'
})

const label = computed(() => {
  return props.label ?? STATUS_MAP[props.status]?.text ?? props.status
})
</script>
