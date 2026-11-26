<template>
  <Transition name="fade">
    <div
      v-if="visible"
      :class="['alert', alertClass, dismissible ? 'alert-dismissible' : '', 'd-flex align-items-center gap-2']"
      role="alert"
    >
      <i :class="iconClass"></i>
      <span>{{ message }}</span>
      <button
        v-if="dismissible"
        type="button"
        class="btn-close ms-auto"
        @click="visible = false"
        aria-label="Close"
      ></button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
interface Props {
  type?: 'success' | 'danger' | 'warning' | 'info'
  message: string
  dismissible?: boolean
  modelValue?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'info',
  dismissible: true,
  modelValue: true,
})

const emit = defineEmits<{ 'update:modelValue': [val: boolean] }>()

const visible = ref(props.modelValue)

watch(() => props.modelValue, (v) => { visible.value = v })
watch(visible, (v) => emit('update:modelValue', v))

const alertClass = computed(() => `alert-${props.type}`)

const iconClass = computed(() => {
  const map: Record<string, string> = {
    success: 'ti ti-circle-check',
    danger:  'ti ti-circle-x',
    warning: 'ti ti-alert-triangle',
    info:    'ti ti-info-circle',
  }
  return map[props.type ?? 'info']
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
