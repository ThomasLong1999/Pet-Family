<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue'
import type { CreatePetRequest } from '../types'
import { petsApi } from '../api/pets'
import { onDateBlur } from '../composables/smartDate'
import { t } from '../composables/useI18n'
import { useToast } from '../composables/useToast'

const emit = defineEmits<{ close: []; saved: [] }>()

const form = ref<CreatePetRequest & { adopted_at: string }>({
  species: 'cat', name: '', breed: '', gender: 'male', birthday: '', color: '', adopted_at: '', note: '',
})
const saving = ref(false)
const avatarFile = ref<File | null>(null)
const { error: toastError } = useToast()
const avatarPreview = ref<string | null>(null)

const speciesOptions = [
  { value: 'cat', key: 'species.cat' },
  { value: 'dog', key: 'species.dog' },
  { value: 'hamster', key: 'species.hamster' },
  { value: 'rabbit', key: 'species.rabbit' },
]

function onAvatarSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  // Revoke previous preview if re-selecting
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
  avatarFile.value = input.files[0]
  avatarPreview.value = URL.createObjectURL(avatarFile.value)
}

async function submit() {
  if (!form.value.name || !form.value.birthday) return
  saving.value = true
  try {
    const payload: CreatePetRequest = { ...form.value, adopted_at: form.value.adopted_at || null }
    const pet = await petsApi.create(payload)
    if (avatarFile.value) {
      try {
        await petsApi.uploadAvatar(pet.id, avatarFile.value)
      } catch (err) {
        toastError(t('error.uploadFailed') + (err as Error).message)
      }
    }
    emit('saved')
  } catch (err) {
    toastError(t('error.saveFailed') + (err as Error).message)
  } finally {
    saving.value = false
  }
}

// Release object URL when the form is closed/unmounted
onBeforeUnmount(() => {
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
})
</script>

<template>
  <div class="form-panel">
    <div class="form-header">
      <h2>{{ t('form.addPet') }}</h2>
      <button class="btn btn-ghost" @click="emit('close')">{{ t('form.cancel') }}</button>
    </div>
    <form class="pet-form" @submit.prevent="submit">
      <div class="avatar-section">
        <label class="form-label">{{ t('form.avatar') }}</label>
        <div class="avatar-upload-area" @click="($refs.avatarInput as HTMLInputElement)?.click()">
          <img v-if="avatarPreview" :src="avatarPreview" class="avatar-preview" />
          <div v-else class="avatar-placeholder">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M12 5v14M5 12h14" /></svg>
            <span>{{ t('form.addPhoto') }}</span>
          </div>
          <input ref="avatarInput" type="file" accept="image/*" class="hidden-input" @change="onAvatarSelected" />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('form.type') }}</label>
          <select v-model="form.species" class="form-input">
            <option v-for="opt in speciesOptions" :key="opt.value" :value="opt.value">{{ t(opt.key) }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('form.name') }} *</label>
          <input v-model="form.name" class="form-input" :placeholder="t('form.name.ph')" required />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('form.breed') }}</label>
          <input v-model="form.breed" class="form-input" :placeholder="t('form.breed.ph')" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('form.gender') }}</label>
          <select v-model="form.gender" class="form-input">
            <option value="male">{{ t('form.gender.male') }}</option>
            <option value="female">{{ t('form.gender.female') }}</option>
          </select>
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('form.birthday') }} *</label>
          <input v-model="form.birthday" type="text" class="form-input" :placeholder="t('form.birthday.ph')" required @blur="onDateBlur" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('form.adoptedAt') }}</label>
          <input v-model="form.adopted_at" type="text" class="form-input" :placeholder="t('form.adoptedAt.ph')" @blur="onDateBlur" />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">{{ t('form.color') }}</label>
        <input v-model="form.color" class="form-input" :placeholder="t('form.color.ph')" />
      </div>
      <div class="form-group">
        <label class="form-label">{{ t('form.note') }}</label>
        <textarea v-model="form.note" class="form-input" rows="2" :placeholder="t('form.note.ph')"></textarea>
      </div>
      <button type="submit" class="btn btn-primary submit-btn" :disabled="saving">
        {{ saving ? t('form.saving') : t('form.save') }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.form-panel { background: var(--bg-elevated); border-radius: 1.5rem; width: 100%; max-width: 480px; padding: 1.5rem; animation: scaleIn var(--duration-normal) var(--ease-spring); max-height: 90vh; overflow-y: auto; }
.form-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; }
.form-header h2 { font-size: var(--text-xl); font-weight: 800; }
.pet-form { display: flex; flex-direction: column; gap: 1rem; }
.avatar-section { display: flex; flex-direction: column; }
.avatar-upload-area { width: 6rem; height: 6rem; border-radius: 1.25rem; border: 2px dashed var(--bg-muted); overflow: hidden; cursor: pointer; transition: border-color var(--duration-fast) var(--ease-spring); align-self: center; }
.avatar-upload-area:hover { border-color: var(--fg-tertiary); }
.avatar-preview { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.25rem; color: var(--fg-tertiary); font-size: var(--text-xs); font-weight: 500; }
.hidden-input { display: none; }
.form-row { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.75rem; }
.form-group { display: flex; flex-direction: column; }
textarea.form-input { resize: vertical; min-height: 3rem; }
.submit-btn { margin-top: 0.5rem; width: 100%; padding: 0.875rem; font-size: var(--text-base); }
@media (max-width: 480px) { .form-row { grid-template-columns: 1fr; } }
</style>
