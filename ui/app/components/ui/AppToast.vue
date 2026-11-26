<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast-slide">
        <div
          v-for="t in toasts"
          :key="t.id"
          :class="['toast-item', `toast-${t.type}`]"
          @click="dismissToast(t.id)"
        >
          <i :class="iconClass(t.type)"></i>
          <span>{{ t.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from '~/composables/useToast'
import type { ToastItem } from '~/composables/useToast'

const { toasts, dismissToast } = useToast()

function iconClass(type: ToastItem['type']) {
  const map: Record<string, string> = {
    success: 'ti ti-circle-check',
    danger:  'ti ti-circle-x',
    warning: 'ti ti-alert-triangle',
    info:    'ti ti-info-circle',
  }
  return map[type] || map.info
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 99999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 420px;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 20px;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  color: #fff;
}

.toast-success { background: #16a34a; }
.toast-danger  { background: #dc2626; }
.toast-warning { background: #d97706; }
.toast-info    { background: #1d4ed8; }

.toast-item i { font-size: 1.2rem; }

/* Animasi */
.toast-slide-enter-active { transition: all 0.3s ease-out; }
.toast-slide-leave-active { transition: all 0.2s ease-in; }
.toast-slide-enter-from { opacity: 0; transform: translateX(100px); }
.toast-slide-leave-to { opacity: 0; transform: translateX(100px); }
</style>
