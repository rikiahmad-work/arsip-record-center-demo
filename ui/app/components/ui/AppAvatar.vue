<template>
  <div
    class="google-avatar-circle"
    :class="customClass"
    :style="{
      width: `${size}px`,
      height: `${size}px`,
      minWidth: `${size}px`,
      minHeight: `${size}px`,
      backgroundColor: bgColor,
      fontSize: computedFontSize,
    }"
    :title="name || 'Pengguna'"
    aria-hidden="true"
  >
    {{ initial }}
  </div>
</template>

<script setup lang="ts">
interface Props {
  name?: string
  size?: number | string
  length?: number
  fontSize?: string
  customClass?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  name: '',
  size: 28,
  length: 1,
  fontSize: '',
  customClass: '',
  color: '',
})

// Palet warna Google Material Avatar (kontras tinggi dengan teks putih)
const GOOGLE_AVATAR_COLORS = [
  '#1a73e8', // Google Blue
  '#d93025', // Google Red
  '#188038', // Google Green
  '#f29900', // Google Amber / Orange
  '#8e24aa', // Google Purple
  '#00897b', // Google Teal
  '#e91e63', // Google Pink
  '#3949ab', // Google Indigo
  '#0288d1', // Google Light Blue
  '#e64a19', // Google Deep Orange
  '#5e35b1', // Google Deep Purple
  '#00796b', // Google Dark Teal
  '#c2185b', // Google Rose
  '#455a64', // Google Blue Grey
  '#7b1fa2', // Google Violet
  '#388e3c', // Google Dark Green
]

function hashString(str: string): number {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return Math.abs(hash)
}

const bgColor = computed(() => {
  if (props.color) return props.color
  const n = (props.name || '').trim()
  if (!n) return '#1a73e8'
  const idx = hashString(n) % GOOGLE_AVATAR_COLORS.length
  return GOOGLE_AVATAR_COLORS[idx]
})

const initial = computed(() => {
  const val = (props.name || '').trim()
  if (!val) return '?'
  const words = val.split(/\s+/).filter(Boolean)
  if (props.length === 2 && words.length >= 2) {
    return (words[0].charAt(0) + words[words.length - 1].charAt(0)).toUpperCase()
  }
  return words[0].charAt(0).toUpperCase()
})

const computedFontSize = computed(() => {
  if (props.fontSize) return props.fontSize
  const numSize = Number(props.size) || 28
  // Ukuran font ~48% dari diameter lingkaran
  return `${Math.max(10, Math.round(numSize * 0.48))}px`
})
</script>

<style scoped>
.google-avatar-circle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  line-height: 1;
  user-select: none;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12);
  flex-shrink: 0;
  vertical-align: middle;
}
</style>
