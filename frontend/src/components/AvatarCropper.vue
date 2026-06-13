<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{ src: string }>()
const emit = defineEmits<{ confirm: [blob: Blob]; cancel: [] }>()

const canvas = ref<HTMLCanvasElement | null>(null)
const img = ref<HTMLImageElement | null>(null)
const scale = ref(1)
const offsetX = ref(0)
const offsetY = ref(0)
const dragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

onMounted(() => {
  const i = new Image()
  i.crossOrigin = 'anonymous'
  i.onload = () => {
    img.value = i
    // Fit image to canvas
    const size = 280
    scale.value = Math.max(size / i.width, size / i.height)
    offsetX.value = (size - i.width * scale.value) / 2
    offsetY.value = (size - i.height * scale.value) / 2
    draw()
  }
  i.src = props.src
})

function draw() {
  if (!canvas.value || !img.value) return
  const ctx = canvas.value.getContext('2d')!
  const size = 280
  ctx.clearRect(0, 0, size, size)
  ctx.drawImage(img.value, offsetX.value, offsetY.value, img.value.width * scale.value, img.value.height * scale.value)
}

function onPointerDown(e: PointerEvent) {
  dragging.value = true
  dragStart.value = { x: e.clientX - offsetX.value, y: e.clientY - offsetY.value }
}

function onPointerMove(e: PointerEvent) {
  if (!dragging.value) return
  offsetX.value = e.clientX - dragStart.value.x
  offsetY.value = e.clientY - dragStart.value.y
  draw()
}

function onPointerUp() { dragging.value = false }

function onWheel(e: WheelEvent) {
  e.preventDefault()
  const delta = e.deltaY > 0 ? -0.05 : 0.05
  scale.value = Math.max(0.1, scale.value + delta)
  draw()
}

function onSliderChange(e: Event) {
  scale.value = parseFloat((e.target as HTMLInputElement).value)
  if (img.value) {
    const size = 280
    offsetX.value = (size - img.value.width * scale.value) / 2
    offsetY.value = (size - img.value.height * scale.value) / 2
  }
  draw()
}

function confirm() {
  if (!canvas.value || !img.value) return
  // Export at 256x256
  const out = document.createElement('canvas')
  out.width = 256
  out.height = 256
  const ctx = out.getContext('2d')!
  const ratio = 256 / 280
  ctx.drawImage(
    img.value,
    offsetX.value * ratio,
    offsetY.value * ratio,
    img.value.width * scale.value * ratio,
    img.value.height * scale.value * ratio,
  )
  out.toBlob((blob) => { if (blob) emit('confirm', blob) }, 'image/jpeg', 0.9)
}
</script>

<template>
  <div class="cropper-overlay" @click.self="emit('cancel')">
    <div class="cropper-panel">
      <h3>调整头像</h3>
      <div class="cropper-area"
        @pointerdown="onPointerDown" @pointermove="onPointerMove" @pointerup="onPointerUp" @pointerleave="onPointerUp"
        @wheel="onWheel">
        <canvas ref="canvas" width="280" height="280" class="cropper-canvas" />
      </div>
      <div class="cropper-controls">
        <input type="range" min="0.1" max="5" step="0.01" :value="scale" @input="onSliderChange" class="zoom-slider" />
      </div>
      <div class="cropper-actions">
        <button class="btn btn-ghost" @click="emit('cancel')">取消</button>
        <button class="btn btn-primary" @click="confirm">确认</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cropper-overlay {
  position: fixed; inset: 0; z-index: 400;
  background: rgba(0,0,0,0.4); backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
}
.cropper-panel {
  background: var(--bg-elevated); border-radius: 1.25rem;
  padding: 1.5rem; width: 340px;
  animation: scaleIn var(--duration-normal) var(--ease-spring);
}
.cropper-panel h3 { font-size: var(--text-lg); font-weight: 700; margin-bottom: 1rem; text-align: center; }
.cropper-area {
  width: 280px; height: 280px; margin: 0 auto; border-radius: 1rem; overflow: hidden;
  cursor: grab; touch-action: none; background: var(--bg-subtle);
}
.cropper-area:active { cursor: grabbing; }
.cropper-canvas { display: block; }
.cropper-controls { margin-top: 0.75rem; text-align: center; }
.zoom-slider { width: 80%; accent-color: var(--fg-primary); }
.cropper-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1rem; }
</style>
