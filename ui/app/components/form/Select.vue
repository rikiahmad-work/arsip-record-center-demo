<template>
  <div ref="containerRef" :class="[noMargin ? '' : 'form-group mb-3', 'position-relative']">
    <label v-if="label" :for="selectId" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <!-- Mode Searchable Dropdown (Select2-style) -->
    <div v-if="searchable" class="s2-container" :class="{ 'is-disabled': disabled, 'is-invalid': !!error }">
      <button
        :id="selectId"
        type="button"
        class="form-select text-start d-flex align-items-center justify-content-between s2-trigger"
        :class="[size === 'sm' ? 'form-select-sm' : '', error ? 'is-invalid' : '']"
        :disabled="disabled"
        @click="toggleDropdown"
      >
        <span class="text-truncate" :class="selectedOption ? 'text-dark' : 'text-muted'">
          {{ selectedOption ? selectedOption.label : (placeholder || '-- Pilih --') }}
        </span>
        <span
          v-if="clearable && modelValue !== '' && modelValue !== undefined && !disabled"
          class="s2-clear ms-1 me-2 text-muted"
          title="Kosongkan pilihan"
          @click.stop="clearSelection"
        >
          &times;
        </span>
      </button>

      <!-- Dropdown Panel dengan Input Pencarian -->
      <div v-if="isOpen" class="s2-dropdown shadow-sm border rounded bg-white">
        <div class="p-2 border-bottom bg-light">
          <div class="input-group input-group-sm">
            <span class="input-group-text bg-white border-end-0">
              <i class="ti ti-search text-muted"></i>
            </span>
            <input
              ref="searchInputRef"
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              :placeholder="searchPlaceholder"
              @keydown.esc.stop="closeDropdown"
              @keydown.down.prevent="moveHighlight(1)"
              @keydown.up.prevent="moveHighlight(-1)"
              @keydown.enter.prevent="selectHighlighted"
            >
          </div>
        </div>

        <ul ref="listRef" class="list-unstyled mb-0 s2-options-list">
          <li
            v-if="placeholder"
            class="s2-option text-muted"
            :class="{ 'is-selected': modelValue === '' || modelValue === undefined }"
            @click="selectValue('')"
          >
            {{ placeholder }}
          </li>
          <li
            v-for="(opt, idx) in filteredOptions"
            :key="opt.value"
            class="s2-option"
            :class="{
              'is-selected': String(opt.value) === String(modelValue),
              'is-highlighted': idx === highlightedIndex,
              'is-disabled': opt.disabled,
            }"
            @mouseenter="highlightedIndex = idx"
            @click="!opt.disabled && selectValue(String(opt.value))"
          >
            {{ opt.label }}
          </li>
          <li v-if="filteredOptions.length === 0" class="px-3 py-2 text-muted small text-center">
            Tidak ditemukan hasil untuk "{{ searchQuery }}"
          </li>
        </ul>
      </div>
    </div>

    <!-- Mode Native Select -->
    <select
      v-else
      :id="selectId"
      :class="['form-select', size === 'sm' ? 'form-select-sm' : '', error ? 'is-invalid' : '']"
      :value="modelValue"
      :disabled="disabled"
      @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option
        v-for="opt in options"
        :key="opt.value"
        :value="opt.value"
        :disabled="opt.disabled"
      >
        {{ opt.label }}
      </option>
    </select>

    <div v-if="error" class="invalid-feedback d-block">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
export interface SelectOption {
  value: string | number
  label: string
  disabled?: boolean
}

interface Props {
  modelValue?: string | number
  label?: string
  options: SelectOption[]
  placeholder?: string
  searchPlaceholder?: string
  error?: string
  required?: boolean
  disabled?: boolean
  searchable?: boolean
  clearable?: boolean
  noMargin?: boolean
  size?: 'sm' | 'md'
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  searchable: true,
  clearable: false,
  noMargin: false,
  size: 'md',
  searchPlaceholder: 'Ketik untuk mencari...',
})

const emit = defineEmits<{ 'update:modelValue': [val: string] }>()

const selectId = computed(() => props.id ?? `select-${Math.random().toString(36).slice(2, 7)}`)

const containerRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)
const listRef = ref<HTMLElement | null>(null)

const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(0)

const selectedOption = computed(() =>
  props.options.find((o) => String(o.value) === String(props.modelValue) && String(props.modelValue) !== '')
)

const filteredOptions = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return props.options
  return props.options.filter(
    (o) =>
      String(o.label).toLowerCase().includes(q) ||
      String(o.value).toLowerCase().includes(q)
  )
})

watch(searchQuery, () => {
  highlightedIndex.value = 0
})

async function toggleDropdown() {
  if (props.disabled) return
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    searchQuery.value = ''
    highlightedIndex.value = 0
    await nextTick()
    searchInputRef.value?.focus()
  }
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
}

function selectValue(val: string) {
  emit('update:modelValue', val)
  closeDropdown()
}

function clearSelection() {
  emit('update:modelValue', '')
  closeDropdown()
}

function moveHighlight(step: number) {
  if (filteredOptions.value.length === 0) return
  const max = filteredOptions.value.length - 1
  let next = highlightedIndex.value + step
  if (next < 0) next = max
  if (next > max) next = 0
  highlightedIndex.value = next
}

function selectHighlighted() {
  const opt = filteredOptions.value[highlightedIndex.value]
  if (opt && !opt.disabled) {
    selectValue(String(opt.value))
  }
}

function handleClickOutside(e: MouseEvent) {
  if (!isOpen.value) return
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})
</script>

<style scoped>
.s2-container {
  position: relative;
  width: 100%;
}

.s2-trigger {
  background-color: #fff;
  cursor: pointer;
  padding-right: 2rem;
}

.s2-clear {
  font-size: 1.15rem;
  line-height: 1;
  font-weight: 700;
  cursor: pointer;
}

.s2-clear:hover {
  color: #dc3545 !important;
}

.s2-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 1060;
}

.s2-options-list {
  max-height: 240px;
  overflow-y: auto;
}

.s2-option {
  padding: 0.45rem 0.75rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.1s ease;
  border-bottom: 1px solid #f5f6f8;
}

.s2-option:last-child {
  border-bottom: none;
}

.s2-option:hover,
.s2-option.is-highlighted {
  background-color: #e6f4ff;
  color: #0958d9;
}

.s2-option.is-selected {
  background-color: #1890ff;
  color: #ffffff;
  font-weight: 600;
}

.s2-option.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
