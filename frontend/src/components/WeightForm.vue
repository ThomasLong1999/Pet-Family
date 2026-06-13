<script setup lang="ts">
import { ref } from 'vue'
import type { CreateWeightRequest } from '../types'
import { weightsApi } from '../api/weights'
import { onDateBlur } from '../composables/smartDate'
import { t } from '../composables/useI18n'
import { useToast } from '../composables/useToast'

const props = defineProps<{
  petId: string
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const today = new Date().toISOString().slice(0, 10)
const form = ref<CreateWeightRequest>({
  weight: 0,
  recorded_at: today,
  note: '',
})

const { error: toastError } = useToast()
const saving = ref(false)

async function submit() {
  if (form.value.weight <= 0) return
  saving.value = true
  try {
    await weightsApi.create(props.petId, form.value)
    emit('saved')
  } catch (err) {
    toastError(t('error.saveFailed') + (err as Error).message)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="form-panel">
    <div class="form-header">
      <h2>{{ t('weight.record') }}</h2>
      <button class="btn btn-ghost" @click="emit('close')">{{ t('form.cancel') }}</button>
    </div>

    <form class="weight-form" @submit.prevent="submit">
      <div class="form-group">
        <label class="form-label">{{ t('weight.weightLabel') }} *</label>
        <input v-model.number="form.weight" type="number" step="0.01" min="0" class="form-input" placeholder="4.50" required />
      </div>

      <div class="form-group">
        <label class="form-label">{{ t('weight.date') }}</label>
        <input v-model="form.recorded_at" type="text" class="form-input" :placeholder="t('weight.date.ph')" @blur="onDateBlur" />
      </div>

      <div class="form-group">
        <label class="form-label">{{ t('weight.note') }}</label>
        <input v-model="form.note" class="form-input" :placeholder="t('weight.note.ph')" />
      </div>

      <button type="submit" class="btn btn-primary submit-btn" :disabled="saving">
        {{ saving ? t('form.saving') : t('form.saveChanges') }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.form-panel {
  background: var(--bg-elevated);
  border-radius: 1rem;
  width: 100%;
  max-width: 360px;
  padding: 1.5rem;
  animation: scaleIn var(--duration-normal) var(--ease-spring);
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.form-header h2 { font-size: var(--text-lg); font-weight: 700; }

.weight-form { display: flex; flex-direction: column; gap: 0.875rem; }
.form-group { display: flex; flex-direction: column; }
.submit-btn { margin-top: 0.5rem; width: 100%; }
</style>
