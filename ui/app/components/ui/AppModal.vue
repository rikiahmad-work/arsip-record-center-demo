<template>
  <div
    :id="modalId"
    class="modal fade show"
    tabindex="-1"
    :aria-labelledby="`${modalId}Label`"
    aria-modal="true"
    role="dialog"
    style="display: block; background: rgba(0, 0, 0, 0.5);"
    @click.self="emit('close')"
  >
    <div :class="['modal-dialog modal-dialog-centered', sizeClass, scrollable ? 'modal-dialog-scrollable' : '']">
      <div class="modal-content shadow">
        <div :class="['modal-header', headerClass]">
          <h5 :id="`${modalId}Label`" class="modal-title">{{ title }}</h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
            @click="emit('close')"
          ></button>
        </div>
        <div class="modal-body">
          <slot />
          <slot name="body" />
        </div>
        <div v-if="$slots.footer" class="modal-footer">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  id?: string
  title: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  scrollable?: boolean
  headerClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  id: undefined,
  size: 'md',
  scrollable: false,
  headerClass: '',
})

const emit = defineEmits<{
  close: []
}>()

const modalId = computed(() => props.id || `modal-${Math.random().toString(36).slice(2, 8)}`)

const sizeClass = computed(() => {
  const map: Record<string, string> = {
    sm: 'modal-sm',
    md: '',
    lg: 'modal-lg',
    xl: 'modal-xl',
  }
  return map[props.size] ?? ''
})
</script>
