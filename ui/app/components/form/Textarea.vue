<template>
  <div class="form-group mb-3">
    <label v-if="label" :for="textareaId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>
    <textarea
      :id="textareaId"
      :class="['form-control', error ? 'is-invalid' : '']"
      :placeholder="placeholder"
      :rows="rows"
      :disabled="disabled"
      :value="modelValue"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    ></textarea>
    <div v-if="error" class="invalid-feedback">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue?: string
  label?: string
  placeholder?: string
  rows?: number
  error?: string
  required?: boolean
  disabled?: boolean
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  rows: 3,
})

const emit = defineEmits<{ 'update:modelValue': [val: string] }>()

const textareaId = computed(() => props.id ?? `textarea-${Math.random().toString(36).slice(2, 7)}`)
</script>
