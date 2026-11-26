<template>
  <div class="row">
    <!-- 1. Total Arsip -->
    <div class="col-md-6 col-xl-3">
      <div
        class="card stat-card bg-primary text-white cursor-pointer"
        @click="goToArsip()"
      >
        <div class="card-body">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <h4 class="mb-1 text-white f-w-700">{{ stats.totalArsip }}</h4>
              <h6 class="mb-0 text-white-50 f-w-500">Total Arsip</h6>
            </div>
            <div class="stat-icon">
              <i class="ti ti-archive"></i>
            </div>
          </div>
          <div class="mt-3 d-flex align-items-center justify-content-between">
            <span class="badge bg-white text-primary">
              <i class="ti ti-database me-1"></i>Record Center KemenPPPA
            </span>
            <NuxtLink
              to="/pengelolaan-arsip"
              class="stat-arrow-btn text-white"
              title="Lihat Semua Arsip"
              @click.stop
            >
              <i class="ti ti-arrow-right"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 2. Arsip Inaktif -->
    <div class="col-md-6 col-xl-3">
      <div
        class="card stat-card bg-success text-white cursor-pointer"
        @click="goToArsip({ kategori: 'Inaktif' })"
      >
        <div class="card-body">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <h4 class="mb-1 text-white f-w-700">{{ stats.arsipInaktif }}</h4>
              <h6 class="mb-0 text-white-50 f-w-500">Arsip Inaktif</h6>
            </div>
            <div class="stat-icon">
              <i class="ti ti-file-text"></i>
            </div>
          </div>
          <div class="mt-3 d-flex align-items-center justify-content-between">
            <span class="badge bg-white text-success">
              <i class="ti ti-check me-1"></i>Tersimpan di Boks/Rak
            </span>
            <NuxtLink
              :to="{ path: '/pengelolaan-arsip', query: { kategori: 'Inaktif' } }"
              class="stat-arrow-btn text-white"
              title="Lihat Arsip Inaktif"
              @click.stop
            >
              <i class="ti ti-arrow-right"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. Arsip Vital & Terjaga -->
    <div class="col-md-6 col-xl-3">
      <div
        class="card stat-card bg-warning text-white cursor-pointer"
        @click="goToArsip({ kategori: 'Vital' })"
      >
        <div class="card-body">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <h4 class="mb-1 text-white f-w-700">{{ stats.arsipVital ?? stats.arsipAktif }}</h4>
              <h6 class="mb-0 text-white-50 f-w-500">Arsip Vital & Terjaga</h6>
            </div>
            <div class="stat-icon">
              <i class="ti ti-shield-lock"></i>
            </div>
          </div>
          <div class="mt-3 d-flex align-items-center justify-content-between">
            <span class="badge bg-white text-warning">
              <i class="ti ti-lock me-1"></i>Terjaga: {{ stats.arsipTerjaga ?? 0 }} berkas
            </span>
            <NuxtLink
              :to="{ path: '/pengelolaan-arsip', query: { kategori: 'Vital' } }"
              class="stat-arrow-btn text-white"
              title="Lihat Arsip Vital & Terjaga"
              @click.stop
            >
              <i class="ti ti-arrow-right"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 4. Jatuh Tempo / Usul Susut -->
    <div class="col-md-6 col-xl-3">
      <div
        class="card stat-card bg-danger text-white cursor-pointer"
        @click="goToArsip({ kategori: 'Usul Musnah' })"
      >
        <div class="card-body">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <h4 class="mb-1 text-white f-w-700">{{ stats.jatuhTempoRetensi }}</h4>
              <h6 class="mb-0 text-white-50 f-w-500">Jatuh Tempo / Usul Susut</h6>
            </div>
            <div class="stat-icon">
              <i class="ti ti-clock"></i>
            </div>
          </div>
          <div class="mt-3 d-flex align-items-center justify-content-between">
            <span class="badge bg-white text-danger">
              <i class="ti ti-alert-triangle me-1"></i>Dipinjam: {{ stats.arsipDipinjam ?? 0 }} berkas
            </span>
            <NuxtLink
              :to="{ path: '/pengelolaan-arsip', query: { kategori: 'Usul Musnah' } }"
              class="stat-arrow-btn text-white"
              title="Lihat Arsip Jatuh Tempo / Usul Susut"
              @click.stop
            >
              <i class="ti ti-arrow-right"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardStats } from '~/types/arsip'

interface Props {
  stats: DashboardStats
}

withDefaults(defineProps<Props>(), {
  stats: () => ({
    totalArsip: 0,
    arsipVital: 0,
    arsipAktif: 0,
    arsipInaktif: 0,
    arsipTerjaga: 0,
    jatuhTempoRetensi: 0,
    arsipDipinjam: 0,
  }),
})

const router = useRouter()

function goToArsip(query?: Record<string, string>) {
  router.push({
    path: '/pengelolaan-arsip',
    query: query || {},
  })
}
</script>

<style scoped>
.stat-card {
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.14);
}

.stat-card .card-body {
  padding: 1.25rem;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.stat-arrow-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.22);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  transition: background-color 0.15s ease, transform 0.15s ease;
}

.stat-card:hover .stat-arrow-btn {
  background: rgba(255, 255, 255, 0.35);
  transform: translateX(2px);
}

.badge {
  font-weight: 500;
  font-size: 0.75rem;
  padding: 0.25rem 0.65rem;
  border-radius: 6px;
}
</style>
