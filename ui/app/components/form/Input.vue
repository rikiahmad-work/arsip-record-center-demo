<template>
  <div class="form-group mb-3">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>
    <input
      :id="inputId"
      :type="type"
      :class="['form-control', error ? 'is-invalid' : '']"
      :placeholder="placeholder"
      :value="modelValue"
      :disabled="disabled"
      :readonly="readonly"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <div v-if="error" class="invalid-feedback">{{ error }}</div>
    <div v-if="hint" class="form-text text-muted">{{ hint }}</div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: string
  label?: string
  type?: string
  placeholder?: string
  error?: string
  hint?: string
  required?: boolean
  disabled?: boolean
  readonly?: boolean
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  modelValue: '',
})

const emit = defineEmits<{ 'update:modelValue': [val: string] }>()

const inputId = computed(() => props.id ?? `input-${Math.random().toString(36).slice(2, 7)}`)
</script>
