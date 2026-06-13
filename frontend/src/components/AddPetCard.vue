<script setup lang="ts">
import { t } from '../composables/useI18n'
const emit = defineEmits<{ click: [] }>()
</script>

<template>
  <button type="button" class="add-card" :aria-label="t('btn.addPet')" @click="emit('click')">
    <svg width="28" height="28" viewBox="0 0 28 28" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" aria-hidden="true">
      <line x1="14" y1="6" x2="14" y2="22" /><line x1="6" y1="14" x2="22" y2="14" />
    </svg>
    <span>{{ t('btn.addPet') }}</span>
  </button>
</template>

<style scoped>
.add-card {
  cursor: pointer; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.5rem;
  aspect-ratio: 4 / 3; border: 2px dashed var(--bg-muted); border-radius: var(--card-radius);
  background: transparent;
  color: var(--fg-tertiary); font-size: var(--text-sm); font-weight: 600;
  font-family: var(--font-family);
  transition: all var(--duration-fast) var(--ease-spring);
}
.add-card:hover { border-color: var(--fg-tertiary); color: var(--fg-secondary); transform: translateY(-2px); }
.add-card:focus-visible { outline: 2px solid var(--info); outline-offset: 2px; }
</style>
