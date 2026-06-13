<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'

const props = withDefaults(
  defineProps<{
    title?: string
    label?: string // aria-label fallback when no title
    closeOnOverlay?: boolean
  }>(),
  { title: '', label: '', closeOnOverlay: true },
)

const emit = defineEmits<{ close: [] }>()

const panelRef = ref<HTMLElement | null>(null)

const FOCUSABLE = [
  'a[href]',
  'button:not([disabled])',
  'textarea',
  'input',
  'select',
  '[tabindex]:not([tabindex="-1"])',
].join(',')

function getFocusable(): HTMLElement[] {
  if (!panelRef.value) return []
  return Array.from(panelRef.value.querySelectorAll<HTMLElement>(FOCUSABLE)).filter(
    (el) => el.offsetParent !== null || el === document.activeElement,
  )
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.preventDefault()
    emit('close')
    return
  }
  if (e.key !== 'Tab') return
  const els = getFocusable()
  if (els.length === 0) return
  const first = els[0]
  const last = els[els.length - 1]
  const active = document.activeElement as HTMLElement
  if (e.shiftKey) {
    if (active === first || !panelRef.value?.contains(active)) {
      e.preventDefault()
      last.focus()
    }
  } else {
    if (active === last) {
      e.preventDefault()
      first.focus()
    }
  }
}

function onOverlayClick() {
  if (props.closeOnOverlay) emit('close')
}

// Focus first element on mount
nextTick(() => getFocusable()[0]?.focus())

// Restore focus to trigger element on close is handled by parent via unmount.
onBeforeUnmount(() => {
  // nothing — parent manages lifecycle
})

// Prevent body scroll while open
watch(
  () => true,
  () => {
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    onBeforeUnmount(() => {
      document.body.style.overflow = prev
    })
  },
  { immediate: true },
)
</script>

<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="onOverlayClick">
      <div
        ref="panelRef"
        class="modal-panel"
        role="dialog"
        aria-modal="true"
        :aria-label="label || title"
        @keydown="onKeydown"
      >
        <slot />
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.overlay-enter-active {
  transition: opacity var(--duration-normal) var(--ease-spring);
}

.overlay-leave-active {
  transition: opacity var(--duration-fast) var(--ease-smooth);
}

.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .overlay-enter-active,
  .overlay-leave-active {
    transition: none;
  }
}
</style>
