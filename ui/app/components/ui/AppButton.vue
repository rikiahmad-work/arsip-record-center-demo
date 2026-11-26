<template>
  <button
    :type="type"
    :class="[
      'btn',
      variantClass,
      sizeClass,
      'd-inline-flex align-items-center gap-2',
      { disabled: disabled || loading },
    ]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
    <i v-else-if="icon" :class="icon"></i>
    <slot />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'success' | 'danger' | 'warning' | 'secondary' | 'light-primary' | 'light-success' | 'light-danger' | 'outline-primary' | 'outline-secondary'
  size?: 'sm' | 'md' | 'lg'
  icon?: string
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  type: 'button',
})

const variantClass = computed(() => {
  const map: Record<string, string> = {
    primary: 'btn-primary',
    success: 'btn-success',
    danger: 'btn-danger',
    warning: 'btn-warning',
    secondary: 'btn-secondary',
    'light-primary': 'btn-light-primary',
    'light-success': 'btn-light-success',
    'light-danger': 'btn-light-danger',
    'outline-primary': 'btn-outline-primary',
    'outline-secondary': 'btn-outline-secondary',
  }
  return map[props.variant] ?? 'btn-primary'
})

const sizeClass = computed(() => {
  const map: Record<string, string> = {
    sm: 'btn-sm',
    md: '',
    lg: 'btn-lg',
  }
  return map[props.size] ?? ''
})
</script>
