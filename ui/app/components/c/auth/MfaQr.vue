<template>
  <div class="mfa-qr-container d-flex flex-column align-items-center">
    <canvas ref="canvasRef" class="border rounded p-2 bg-white shadow-sm"></canvas>
  </div>
</template>

<script setup lang="ts">
import QRCode from 'qrcode'

interface Props {
  value: string
  size?: number
}

const props = withDefaults(defineProps<Props>(), {
  value: '',
  size: 200,
})

const canvasRef = ref<HTMLCanvasElement | null>(null)

async function renderQR() {
  if (!props.value || !canvasRef.value) return
  try {
    await QRCode.toCanvas(canvasRef.value, props.value, {
      width: props.size,
      margin: 1,
      color: {
        dark: '#17212b',
        light: '#ffffff',
      },
    })
  } catch (err) {
    console.error('Gagal membuat QR Code MFA:', err)
  }
}

watch(() => props.value, () => {
  renderQR()
})

onMounted(() => {
  renderQR()
})
</script>

<style scoped>
.mfa-qr-container {
  display: inline-block;
}
</style>
