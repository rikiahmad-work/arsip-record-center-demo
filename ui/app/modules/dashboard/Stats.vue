<template>
  <div class="row">
    <UiAppStatCard
      :value="stats.totalMasuk"
      label="Total Laporan Masuk"
      icon="ti ti-file-analytics"
      bg-color="primary"
    />
    <UiAppStatCard
      :value="stats.selesai"
      label="Selesai Ditindaklanjuti"
      icon="ti ti-circle-check"
      bg-color="success"
    />
    <UiAppStatCard
      :value="stats.belumDitindaklanjuti"
      label="Belum Ditindaklanjuti"
      icon="ti ti-clock"
      bg-color="warning"
    />
  </div>
</template>

<script setup lang="ts">
interface Stats {
  totalMasuk: number
  selesai: number
  belumDitindaklanjuti: number
}

interface Props {
  stats: Stats
}

withDefaults(defineProps<Props>(), {
  stats: () => ({
    totalMasuk: 0,
    selesai: 0,
    belumDitindaklanjuti: 0,
  }),
})
</script>
