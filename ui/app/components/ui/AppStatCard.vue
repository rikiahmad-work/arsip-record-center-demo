<template>
  <div class="col-md-4">
    <div :class="['card', `bg-${bgColor}`, 'text-white']">
      <div class="card-body">
        <div class="row align-items-center">
          <div class="col-8">
            <h4 class="mb-1 text-white">{{ value }}</h4>
            <h6 class="mb-0 text-white">{{ label }}</h6>
          </div>
          <div class="col-4 text-end">
            <i :class="[icon, 'f-36 text-white-75']"></i>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  value: number | string
  label: string
  icon: string
  bgColor?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
}

withDefaults(defineProps<Props>(), {
  bgColor: 'primary',
})
</script>
