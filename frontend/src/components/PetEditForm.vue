<script setup lang="ts">
import { ref } from 'vue'
import type { Pet, UpdatePetRequest } from '../types'
import { petsApi } from '../api/pets'
import { onDateBlur } from '../composables/smartDate'
import { t, speciesKey } from '../composables/useI18n'
import { useToast } from '../composables/useToast'

const props = defineProps<{
  pet: Pet
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const form = ref<UpdatePetRequest & {
  species: Pet['species']
  name: string
  breed: string
  gender: Pet['gender']
  birthday: string
  color: string
  adopted_at: string | null
  passed_at: string | null
  note: string | null
}>({
  species: props.pet.species,
  name: props.pet.name,
  breed: props.pet.breed,
  gender: props.pet.gender,
  birthday: props.pet.birthday,
  color: props.pet.color,
  adopted_at: props.pet.adopted_at ?? '',
  passed_at: props.pet.passed_at ?? '',
  note: props.pet.note ?? '',
})

const { error: toastError } = useToast()
const saving = ref(false)

const speciesOptions = [
  { value: 'cat', key: 'species.cat' },
  { value: 'dog', key: 'species.dog' },
  { value: 'hamster', key: 'species.hamster' },
  { value: 'rabbit', key: 'species.rabbit' },
]

async function submit() {
  if (!form.value.name || !form.value.birthday) return
  saving.value = true
  try {
    await petsApi.update(props.pet.id, {
      species: form.value.species,
      name: form.value.name,
      breed: form.value.breed,
      gender: form.value.gender,
      birthday: form.value.birthday,
      color: form.value.color,
      adopted_at: form.value.adopted_at || null,
      passed_at: form.value.passed_at || null,
      note: form.value.note || null,
    })
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
      <h2>{{ t('form.editPet') }}</h2>
      <button class="btn btn-ghost" @click="emit('close')">{{ t('form.cancel') }}</button>
    </div>

    <form class="pet-form" @submit.prevent="submit">
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('form.type') }}</label>
          <select v-model="form.species" class="form-input">
            <option v-for="opt in speciesOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.key) }}
            </option>
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
          <label class="form-label">{{ t('form.color') }}</label>
          <input v-model="form.color" class="form-input" :placeholder="t('form.color.ph')" />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('form.adoptedAt') }}</label>
          <input v-model="form.adopted_at" type="text" class="form-input" :placeholder="t('form.adoptedAt.ph')" @blur="onDateBlur" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ speciesKey(props.pet.species, 'form.passedAt') }}</label>
          <input v-model="form.passed_at" type="text" class="form-input" :placeholder="t('form.passedAt.ph')" @blur="onDateBlur" />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">{{ t('form.note') }}</label>
        <textarea v-model="form.note" class="form-input" rows="2" :placeholder="t('form.note.ph')"></textarea>
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
  border-radius: 1.5rem;
  width: 100%;
  max-width: 480px;
  padding: 1.5rem;
  animation: scaleIn var(--duration-normal) var(--ease-spring);
  max-height: 90vh;
  overflow-y: auto;
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.form-header h2 {
  font-size: var(--text-xl);
  font-weight: 800;
}

.pet-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

textarea.form-input { resize: vertical; min-height: 3rem; }

.submit-btn {
  margin-top: 0.5rem;
  width: 100%;
  padding: 0.875rem;
  font-size: var(--text-base);
}

@media (max-width: 480px) {
  .form-row { grid-template-columns: 1fr; }
}
</style>
