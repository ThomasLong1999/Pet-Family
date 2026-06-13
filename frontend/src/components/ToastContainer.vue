<script setup lang="ts">
import { useToast } from '../composables/useToast'

const { toasts, dismiss } = useToast()
</script>

<template>
  <TransitionGroup name="toast" tag="div" class="toast-container" aria-live="polite" aria-atomic="false">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="toast"
      :class="`toast-${toast.type}`"
      :role="toast.type === 'error' ? 'alert' : 'status'"
    >
      <span class="toast-icon">{{ toast.type === 'error' ? '⚠' : '✓' }}</span>
      <span class="toast-message">{{ toast.message }}</span>
      <button class="toast-close" :aria-label="'close'" @click="dismiss(toast.id)">×</button>
    </div>
  </TransitionGroup>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 1rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  pointer-events: none;
  max-width: calc(100vw - 2rem);
}

.toast {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: 0.75rem;
  background: var(--bg-elevated);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  font-size: var(--text-sm);
  font-weight: 500;
  pointer-events: auto;
  max-width: 28rem;
}

.toast-error {
  border-left: 3px solid var(--danger);
  color: var(--danger);
}

.toast-success {
  border-left: 3px solid var(--success);
  color: var(--success);
}

.toast-icon {
  flex-shrink: 0;
  font-size: 1rem;
}

.toast-message {
  flex: 1;
  word-break: break-word;
  color: var(--fg-primary);
}

.toast-close {
  flex-shrink: 0;
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  background: none;
  color: var(--fg-tertiary);
  cursor: pointer;
  border-radius: 0.375rem;
  font-size: 1.125rem;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.toast-close:hover {
  background: var(--bg-subtle);
  color: var(--fg-primary);
}

/* Transition */
.toast-enter-active {
  transition: all var(--duration-normal) var(--ease-spring);
}

.toast-leave-active {
  transition: all var(--duration-fast) var(--ease-smooth);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-12px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}
</style>
