<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Pet, PetSummary, HealthReminder } from '../types'
import { petsApi } from '../api/pets'
import { t } from '../composables/useI18n'
import AvatarCropper from './AvatarCropper.vue'

const props = defineProps<{ pet: Pet; summary: PetSummary | null; reminders: HealthReminder[] }>()
const emit = defineEmits<{ click: []; avatarChanged: [] }>()

const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const cropSrc = ref<string | null>(null)

const speciesIcon = computed(() => {
  const map: Record<string, string> = { cat: '🐱', dog: '🐶', hamster: '🐹', rabbit: '🐰' }
  return map[props.pet.species] || '🐾'
})

const cardBgStyle = computed(() => {
  const hex = props.pet.dominant_color
  if (!hex || hex === '#888888') return {}
  return { background: `linear-gradient(135deg, ${hex}0C 0%, ${hex}18 100%)` }
})

const genderText = computed(() =>
  props.pet.gender === 'male' ? '♂' : props.pet.gender === 'female' ? '♀' : ''
)
const genderColor = computed(() =>
  props.pet.gender === 'male' ? '#3b82f6' : props.pet.gender === 'female' ? '#ec4899' : 'var(--fg-tertiary)'
)

const ageDisplay = computed(() => {
  if (props.pet.passed_at) return `${props.pet.birthday} — ${props.pet.passed_at}`
  if (!props.pet.birthday) return ''
  const bd = new Date(props.pet.birthday)
  const now = new Date()
  const months = (now.getFullYear() - bd.getFullYear()) * 12 + (now.getMonth() - bd.getMonth())
  if (months < 1) return t('pet.age.lessThan1')
  if (months < 12) return t('pet.age.months', { n: months })
  const years = Math.floor(months / 12)
  const remainMonths = months % 12
  return remainMonths === 0
    ? t('pet.age.yearsRound', { y: years })
    : t('pet.age.years', { y: years, m: remainMonths })
})

const latestWeight = computed(() => props.summary?.latest_weight ?? null)
const weightTrend = computed(() => props.summary?.weight_trend ?? null)
const weightTrendIcon = computed(() => {
  if (!weightTrend.value) return ''
  return weightTrend.value > 0 ? '↑' : weightTrend.value < 0 ? '↓' : '→'
})
const weightTrendColor = computed(() => {
  if (!weightTrend.value) return 'var(--fg-tertiary)'
  return weightTrend.value > 0.05 ? '#f59e0b' : weightTrend.value < -0.05 ? '#22c55e' : 'var(--fg-tertiary)'
})
const hasReminders = computed(() => props.reminders.length > 0)

function onBubbleClick(e: Event) { e.stopPropagation() }

function onUploadClick(e: Event) {
  e.stopPropagation()
  fileInput.value?.click()
}

function onFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  // Show cropper with the selected image
  cropSrc.value = URL.createObjectURL(input.files[0])
  // Reset input so same file can be re-selected
  input.value = ''
}

async function onCropConfirm(blob: Blob) {
  cropSrc.value = null
  uploading.value = true
  try {
    const file = new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
    await petsApi.uploadAvatar(props.pet.id, file)
    emit('avatarChanged')
  } catch (err) {
    console.error('Avatar upload failed:', err)
  } finally {
    uploading.value = false
  }
}

function onCropCancel() {
  if (cropSrc.value) {
    URL.revokeObjectURL(cropSrc.value)
  }
  cropSrc.value = null
}
</script>

<template>
  <div class="pet-card" :style="cardBgStyle" @click="emit('click')">
    <!-- Reminder bubble -->
    <div v-if="hasReminders" class="reminder-bubble" @click="onBubbleClick">
      <span class="bubble-count">{{ reminders.length }}</span>
      <div class="reminder-popover">
        <div v-for="r in reminders" :key="r.name + r.next_date" class="reminder-item"
          :class="{ urgent: r.days_left <= 3, warning: r.days_left > 3 }">
          <span class="reminder-name">{{ r.name }}</span>
          <span class="reminder-date">{{ t('reminder.daysLeft', { n: r.days_left }) }}</span>
        </div>
      </div>
    </div>

    <!-- Left: Photo -->
    <div class="card-photo-side">
      <img v-if="pet.avatar_url" :src="pet.avatar_url" :alt="pet.name" class="card-photo" />
      <div v-else class="card-photo-placeholder">{{ speciesIcon }}</div>
      <button class="avatar-upload-btn" @click="onUploadClick" :disabled="uploading">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <path d="M10.5 3.5L9 2a2.5 2.5 0 10-.6.6L10.5 3.5z" />
        </svg>
      </button>
      <input ref="fileInput" type="file" accept="image/*" class="hidden-file" @change="onFileSelected" />
    </div>

    <!-- Right: Info -->
    <div class="card-info-side">
      <div class="card-name-row">
        <span class="card-species-icon">{{ speciesIcon }}</span>
        <span class="card-name">{{ pet.name }}</span>
      </div>
      <div class="card-meta-row">
        <span class="card-breed">{{ pet.breed }}</span>
        <span v-if="genderText" class="card-gender" :style="{ color: genderColor }">{{ genderText }}</span>
      </div>
      <div class="card-meta-row">
        <span class="card-age">{{ ageDisplay }}</span>
      </div>
      <div v-if="latestWeight !== null" class="card-weight">
        <span class="weight-value">{{ latestWeight.toFixed(1) }} {{ t('pet.weight') }}</span>
        <span v-if="weightTrend && weightTrend !== 0" class="weight-trend" :style="{ color: weightTrendColor }">
          {{ weightTrendIcon }} {{ Math.abs(weightTrend).toFixed(1) }}
        </span>
      </div>
    </div>

    <!-- Avatar cropper overlay -->
    <AvatarCropper v-if="cropSrc" :src="cropSrc" @confirm="onCropConfirm" @cancel="onCropCancel" />
  </div>
</template>

<style scoped>
.pet-card {
  cursor: pointer; display: flex; flex-direction: row; position: relative;
  aspect-ratio: 4 / 3; background: var(--card-bg); border-radius: var(--card-radius);
  box-shadow: var(--card-shadow); overflow: hidden;
  transition: box-shadow var(--duration-normal) var(--ease-spring), transform var(--duration-normal) var(--ease-spring);
}
.pet-card:hover { box-shadow: var(--card-shadow-hover); transform: translateY(-2px); }
.pet-card:active { transform: scale(0.98); }

.reminder-bubble { position: absolute; top: 0.375rem; right: 0.375rem; z-index: 3; }
.bubble-count {
  display: flex; align-items: center; justify-content: center;
  min-width: 1.25rem; height: 1.25rem; padding: 0 0.25rem;
  border-radius: 0.625rem; background: var(--danger); color: white;
  font-size: 0.625rem; font-weight: 700; cursor: pointer;
  transition: transform var(--duration-fast) var(--ease-spring);
}
.bubble-count:hover { transform: scale(1.1); }
.reminder-popover {
  display: none; position: absolute; top: calc(100% + 0.375rem); right: 0;
  min-width: 180px; background: var(--bg-elevated); border-radius: 0.625rem;
  box-shadow: 0 4px 24px rgba(0,0,0,0.12); padding: 0.375rem; z-index: 10;
}
.reminder-bubble:hover .reminder-popover { display: block; animation: scaleIn var(--duration-fast) var(--ease-spring); }
.reminder-item {
  padding: 0.375rem 0.5rem; border-radius: 0.375rem;
  display: flex; justify-content: space-between; align-items: center;
  gap: 0.375rem; font-size: 0.6875rem;
}
.reminder-item.urgent { background: var(--danger-soft); color: var(--danger); }
.reminder-item.warning { background: var(--warning-soft); color: #92400e; }
.reminder-name { font-weight: 600; }
.reminder-date { white-space: nowrap; font-weight: 500; }

.card-photo-side { width: 45%; flex-shrink: 0; position: relative; overflow: hidden; background: var(--bg-subtle); }
.card-photo { width: 100%; height: 100%; object-fit: cover; transition: transform var(--duration-slow) var(--ease-spring); }
.pet-card:hover .card-photo { transform: scale(1.05); }
.card-photo-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; font-size: 2.5rem; opacity: 0.3; }
.avatar-upload-btn {
  position: absolute; bottom: 0.375rem; right: 0.375rem; z-index: 4;
  width: 1.5rem; height: 1.5rem; border-radius: 0.375rem; border: none;
  background: rgba(255,255,255,0.85); backdrop-filter: blur(8px);
  color: var(--fg-secondary); cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  opacity: 0; transition: all var(--duration-fast) var(--ease-spring);
}
.pet-card:hover .avatar-upload-btn { opacity: 1; }
.avatar-upload-btn:hover { background: #fff; transform: scale(1.1); }
.hidden-file { display: none; }

.card-info-side { flex: 1; padding: 0.75rem 0.875rem; display: flex; flex-direction: column; justify-content: center; gap: 0.25rem; min-width: 0; }
.card-name-row { display: flex; align-items: center; gap: 0.25rem; }
.card-species-icon { font-size: 1rem; flex-shrink: 0; }
.card-name { font-size: var(--text-lg); font-weight: 700; color: var(--fg-primary); letter-spacing: -0.01em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-meta-row { display: flex; align-items: center; gap: 0.375rem; font-size: var(--text-xs); color: var(--fg-tertiary); }
.card-breed { font-weight: 500; }
.card-gender { font-weight: 700; }
.card-age { font-weight: 500; }
.card-weight { display: flex; align-items: center; gap: 0.375rem; font-size: var(--text-xs); }
.weight-value { font-weight: 700; color: var(--fg-secondary); }
.weight-trend { font-weight: 600; font-size: 0.6875rem; }
</style>
