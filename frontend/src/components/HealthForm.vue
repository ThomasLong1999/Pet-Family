<script setup lang="ts">
import { ref, computed } from 'vue'
import type { CreateHealthRequest, HealthType } from '../types'
import { healthApi } from '../api/health'
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
const form = ref<CreateHealthRequest>({
  type: 'vaccine',
  name: '',
  date: today,
  next_date: '',
  note: '',
})

const saving = ref(false)
const reportFile = ref<File | null>(null)
const reportName = ref('')
const { error: toastError } = useToast()

const typeOptions: { value: HealthType; key: string }[] = [
  { value: 'vaccine', key: 'health.type.vaccine' },
  { value: 'deworming', key: 'health.type.deworming' },
  { value: 'checkup', key: 'health.type.checkup' },
]

const quickOptions: Record<string, string[]> = {
  vaccine: ['猫三联', '狂犬'],
  deworming: ['大宠爱', '海乐妙', '拜达尔', '普安特', '福来恩', '爱沃克', '喜倍安'],
  checkup: [],
}

const currentQuickOptions = computed(() => quickOptions[form.value.type] || [])

function selectQuickOption(name: string) {
  form.value.name = name
}

function onReportSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  reportFile.value = input.files[0]
  reportName.value = reportFile.value.name
}

async function submit() {
  if (!form.value.name || !form.value.date) return
  saving.value = true
  try {
    const payload = {
      ...form.value,
      next_date: form.value.next_date || null,
      note: form.value.note || null,
    }

    if (form.value.type === 'checkup' && reportFile.value) {
      await healthApi.createWithReport(props.petId, payload, reportFile.value)
    } else {
      await healthApi.create(props.petId, payload)
    }
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
      <h2>{{ t('health.add') }}</h2>
      <button class="btn btn-ghost" @click="emit('close')">{{ t('form.cancel') }}</button>
    </div>

    <form class="health-form" @submit.prevent="submit">
      <!-- Type selector -->
      <div class="form-group">
        <label class="form-label">{{ t('health.type') }} *</label>
        <div class="type-selector">
          <button
            v-for="opt in typeOptions"
            :key="opt.value"
            type="button"
            class="type-btn"
            :class="{ active: form.type === opt.value }"
            @click="form.type = opt.value; form.name = ''"
          >
            {{ t(opt.key) }}
          </button>
        </div>
      </div>

      <!-- Quick options -->
      <div v-if="currentQuickOptions.length > 0" class="form-group">
        <label class="form-label">{{ t('health.quickSelect') }}</label>
        <div class="quick-options">
          <button
            v-for="opt in currentQuickOptions"
            :key="opt"
            type="button"
            class="quick-btn"
            :class="{ active: form.name === opt }"
            @click="selectQuickOption(opt)"
          >
            {{ opt }}
          </button>
        </div>
      </div>

      <!-- Name (auto-filled by quick option or manual input) -->
      <div class="form-group">
        <label class="form-label">{{ t('health.name') }} *</label>
        <input v-model="form.name" class="form-input" :placeholder="t('health.name.ph')" required />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('health.date') }} *</label>
          <input v-model="form.date" type="text" class="form-input" :placeholder="t('weight.date.ph')" required @blur="onDateBlur" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('health.nextDate') }}</label>
          <input v-model="form.next_date" type="text" class="form-input" :placeholder="t('weight.date.ph')" @blur="onDateBlur" />
        </div>
      </div>

      <!-- Report upload (checkup only) -->
      <div v-if="form.type === 'checkup'" class="form-group">
        <label class="form-label">{{ t('health.report') }}</label>
        <label class="report-upload">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
            <path d="M8 2v8M4 6l4-4 4 4" />
            <path d="M2 11v2a1 1 0 001 1h10a1 1 0 001-1v-2" />
          </svg>
          <span v-if="reportName">{{ reportName }}</span>
          <span v-else class="upload-hint">{{ t('health.report.ph') }}</span>
          <input type="file" accept=".pdf,.jpg,.jpeg,.png,.doc,.docx" class="hidden-input" @change="onReportSelected" />
        </label>
      </div>

      <div class="form-group">
        <label class="form-label">{{ t('health.note') }}</label>
        <input v-model="form.note" class="form-input" :placeholder="t('health.note.ph')" />
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
  max-width: 400px;
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

.health-form { display: flex; flex-direction: column; gap: 0.875rem; }
.form-group { display: flex; flex-direction: column; }
.form-row { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.75rem; }

/* Type selector */
.type-selector {
  display: flex;
  gap: 0.5rem;
}

.type-btn {
  flex: 1;
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  border: 1.5px solid var(--bg-muted);
  background: var(--bg-subtle);
  font-family: var(--font-family);
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--fg-secondary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-spring);
}

.type-btn:hover { border-color: var(--fg-tertiary); }

.type-btn.active {
  border-color: var(--fg-primary);
  background: var(--fg-primary);
  color: var(--bg-elevated);
}

/* Quick options */
.quick-options {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.quick-btn {
  padding: 0.375rem 0.75rem;
  border-radius: 2rem;
  border: 1.5px solid var(--bg-muted);
  background: var(--bg-elevated);
  font-family: var(--font-family);
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--fg-secondary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-spring);
}

.quick-btn:hover { border-color: var(--fg-tertiary); }

.quick-btn.active {
  border-color: var(--fg-primary);
  background: var(--fg-primary);
  color: var(--bg-elevated);
}

/* Report upload */
.report-upload {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 0.875rem;
  border-radius: 0.625rem;
  border: 1.5px dashed var(--bg-muted);
  cursor: pointer;
  font-size: var(--text-sm);
  color: var(--fg-secondary);
  transition: border-color var(--duration-fast) var(--ease-spring);
}

.report-upload:hover { border-color: var(--fg-tertiary); }

.upload-hint { color: var(--fg-tertiary); }

.hidden-input { display: none; }

.submit-btn { margin-top: 0.5rem; width: 100%; }

@media (max-width: 400px) {
  .form-row { grid-template-columns: 1fr; }
  .type-selector { flex-wrap: wrap; }
}
</style>
