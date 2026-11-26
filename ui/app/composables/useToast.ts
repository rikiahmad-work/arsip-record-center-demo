import { ref } from 'vue'

export interface ToastItem {
  id: number
  message: string
  type: 'success' | 'danger' | 'warning' | 'info'
}

const toasts = ref<ToastItem[]>([])
let nextId = 0

export function useToast() {
  function showToast(message: string, type: ToastItem['type'] = 'info', duration = 4000) {
    const id = ++nextId
    toasts.value.push({ id, message, type })
    if (duration > 0) {
      setTimeout(() => dismissToast(id), duration)
    }
  }

  function dismissToast(id: number) {
    const idx = toasts.value.findIndex((t) => t.id === id)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }

  return { toasts, showToast, dismissToast }
}
