<template>
  <div class="form-group mb-3">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>
    <div class="input-group">
      <span class="input-group-text bg-light text-secondary">
        <i class="ti ti-calendar-event"></i>
      </span>
      <input
        :id="inputId"
        type="date"
        :class="['form-control', error ? 'is-invalid' : '']"
        :value="modelValue"
        :disabled="disabled"
        :min="min"
        :max="max"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      />
      <div v-if="error" class="invalid-feedback">{{ error }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: string
  label?: string
  error?: string
  required?: boolean
  disabled?: boolean
  min?: string
  max?: string
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
})

const emit = defineEmits<{ 'update:modelValue': [val: string] }>()

const inputId = computed(() => props.id ?? `datepicker-${Math.random().toString(36).slice(2, 7)}`)
</script>
