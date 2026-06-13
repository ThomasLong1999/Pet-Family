<script setup lang="ts">
import { ref } from 'vue'
import type { Pet, UpdatePetRequest } from '../types'
import { petsApi } from '../api/pets'
import { onDateBlur } from '../composables/smartDate'

const props = defineProps<{
  pet: Pet
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const form = ref<UpdatePetRequest & {
  species: string
  name: string
  breed: string
  gender: string
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

const saving = ref(false)

const speciesOptions = [
  { value: 'cat', label: '🐱 猫' },
  { value: 'dog', label: '🐶 狗' },
  { value: 'hamster', label: '🐹 仓鼠' },
  { value: 'rabbit', label: '🐰 兔子' },
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
    alert('保存失败：' + (err as Error).message)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="form-panel">
    <div class="form-header">
      <h2>编辑宠物</h2>
      <button class="btn btn-ghost" @click="emit('close')">取消</button>
    </div>

    <form @submit.prevent="submit" class="pet-form">
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">类型</label>
          <select v-model="form.species" class="form-input">
            <option v-for="opt in speciesOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">名字 *</label>
          <input v-model="form.name" class="form-input" placeholder="宠物名字" required />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">品种</label>
          <input v-model="form.breed" class="form-input" placeholder="如：英短蓝猫" />
        </div>
        <div class="form-group">
          <label class="form-label">性别</label>
          <select v-model="form.gender" class="form-input">
            <option value="male">公 ♂</option>
            <option value="female">母 ♀</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">生日 *</label>
          <input v-model="form.birthday" type="text" class="form-input" placeholder="如 20240301" required @blur="onDateBlur" />
        </div>
        <div class="form-group">
          <label class="form-label">毛色</label>
          <input v-model="form.color" class="form-input" placeholder="如：蓝灰色" />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">领养日期</label>
          <input v-model="form.adopted_at" type="text" class="form-input" placeholder="如 20240415" @blur="onDateBlur" />
        </div>
        <div class="form-group">
          <label class="form-label">去喵星日期</label>
          <input v-model="form.passed_at" type="text" class="form-input" placeholder="如 20260614" @blur="onDateBlur" />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">备注</label>
        <textarea v-model="form.note" class="form-input" rows="2" placeholder="备注信息..."></textarea>
      </div>

      <button type="submit" class="btn btn-primary submit-btn" :disabled="saving">
        {{ saving ? '保存中...' : '保存' }}
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
