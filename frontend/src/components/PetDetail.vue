<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import type { Pet, WeightRecord, HealthRecord, Photo } from '../types'
import { petsApi } from '../api/pets'
import { weightsApi } from '../api/weights'
import { healthApi } from '../api/health'
import { photosApi } from '../api/photos'
import { smartDate } from '../composables/smartDate'
import { t } from '../composables/useI18n'
import AvatarCropper from './AvatarCropper.vue'
import WeightChart from './WeightChart.vue'
import HealthTimeline from './HealthTimeline.vue'
import PhotoGallery from './PhotoGallery.vue'
import WeightForm from './WeightForm.vue'
import HealthForm from './HealthForm.vue'
import PetEditForm from './PetEditForm.vue'

const props = defineProps<{
  pet: Pet
}>()

const emit = defineEmits<{
  close: []
  updated: []
}>()

const activeTab = ref<'info' | 'weight' | 'health' | 'photos'>('info')
const pet = ref<Pet>({ ...props.pet })
const weights = ref<WeightRecord[]>([])
const healthRecords = ref<HealthRecord[]>([])
const photos = ref<Photo[]>([])
const showWeightForm = ref(false)
const showHealthForm = ref(false)
const showEditForm = ref(false)

watch(() => props.pet, (p) => { pet.value = { ...p } })

onMounted(() => {
  loadData()
})

async function loadData() {
  try {
    const [w, h, p] = await Promise.all([
      weightsApi.list(pet.value.id),
      healthApi.list(pet.value.id),
      photosApi.list(pet.value.id),
    ])
    weights.value = w
    healthRecords.value = h
    photos.value = p
  } catch (err) {
    console.error('Failed to load pet data:', err)
  }
}

async function refreshPet() {
  try {
    pet.value = await petsApi.get(pet.value.id)
  } catch {}
}

async function onWeightAdded() {
  showWeightForm.value = false
  weights.value = await weightsApi.list(pet.value.id)
}

async function onHealthAdded() {
  showHealthForm.value = false
  healthRecords.value = await healthApi.list(pet.value.id)
}

const cropSrc = ref<string | null>(null)
const avatarInput = ref<HTMLInputElement | null>(null)

function onAvatarClick() {
  avatarInput.value?.click()
}

function onAvatarFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  cropSrc.value = URL.createObjectURL(input.files[0])
  input.value = ''
}

async function onCropConfirm(blob: Blob) {
  cropSrc.value = null
  try {
    const file = new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
    const result = await petsApi.uploadAvatar(pet.value.id, file)
    pet.value = { ...pet.value, avatar_url: result.avatar_url, dominant_color: result.dominant_color }
    emit('updated')
  } catch (err) {
    console.error('Avatar upload failed:', err)
  }
}

function onCropCancel() {
  if (cropSrc.value) URL.revokeObjectURL(cropSrc.value)
  cropSrc.value = null
}

async function deletePet() {
  if (!confirm(t('passed.delete.confirm', { name: pet.value.name }))) return
  try {
    await petsApi.delete(pet.value.id)
    emit('updated')
  } catch (err) {
    console.error('Delete failed:', err)
  }
}

async function markPassed() {
  const today = new Date().toISOString().slice(0, 10)
  const date = prompt(t('passed.prompt'), today)
  if (!date) return
  try {
    const formatted = smartDate(date)
    await petsApi.update(pet.value.id, { passed_at: formatted })
    await refreshPet()
    emit('updated')
  } catch (err) {
    alert('Error: ' + (err as Error).message)
  }
}

async function clearPassed() {
  if (!confirm(t('passed.clear.confirm', { name: pet.value.name }))) return
  try {
    await petsApi.update(pet.value.id, { passed_at: '' })
    await refreshPet()
    emit('updated')
  } catch (err) {
    alert('保存失败：' + (err as Error).message)
  }
}

async function onPetSaved() {
  showEditForm.value = false
  await refreshPet()
}

const speciesIcon = computed(() => {
  const map: Record<string, string> = { cat: '🐱', dog: '🐶', hamster: '🐹', rabbit: '🐰' }
  return map[pet.value.species] || '🐾'
})

const ageText = computed(() => {
  if (!pet.value.birthday) return ''
  const bd = new Date(pet.value.birthday)
  const now = new Date()
  const months = (now.getFullYear() - bd.getFullYear()) * 12 + (now.getMonth() - bd.getMonth())
  if (months < 1) return '< 1个月'
  if (months < 12) return `${months}个月`
  const years = Math.floor(months / 12)
  const remainMonths = months % 12
  if (remainMonths === 0) return `${years}岁`
  return `${years}岁${remainMonths}个月`
})
</script>

<template>
  <div class="detail-panel" @click.stop>
    <!-- Close button -->
    <button class="close-btn" @click="emit('close')">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <line x1="4" y1="4" x2="16" y2="16" /><line x1="16" y1="4" x2="4" y2="16" />
      </svg>
    </button>

    <!-- Header area -->
    <div class="detail-header">
      <div class="avatar-upload-area" @click="onAvatarClick">
        <img v-if="pet.avatar_url" :src="pet.avatar_url" class="detail-avatar" />
        <div v-else class="detail-avatar-placeholder">{{ speciesIcon }}</div>
        <div class="avatar-overlay">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="white" stroke-width="1.5">
            <path d="M15 13l-2-2a4 4 0 10-1 1l2 2" />
          </svg>
        </div>
        <input ref="avatarInput" type="file" accept="image/*" class="avatar-input" @change="onAvatarFileSelected" />
      </div>

      <div class="detail-header-info">
        <div class="detail-name-row">
          <h2 class="detail-name">{{ pet.name }}</h2>
          <span class="species-badge-sm">{{ speciesIcon }}</span>
        </div>
        <p class="detail-breed">{{ pet.breed }} · {{ pet.gender === 'male' ? '公' : pet.gender === 'female' ? '母' : '' }} · {{ ageText }}</p>
        <p v-if="pet.color" class="detail-color">{{ pet.color }}</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="detail-tabs">
      <button
        v-for="tab in (['info', 'weight', 'health', 'photos'] as const)"
        :key="tab"
        class="tab-btn"
        :class="{ active: activeTab === tab }"
        @click="activeTab = tab"
      >
        {{ { info: t('tab.info'), weight: t('tab.weight'), health: t('tab.health'), photos: t('tab.photos') }[tab] }}
      </button>
    </div>

    <!-- Tab content -->
    <div class="detail-content">
      <!-- Info tab -->
      <div v-if="activeTab === 'info'" class="tab-content info-tab">
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('info.breed') }}</span>
            <span class="info-value">{{ pet.breed || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('info.gender') }}</span>
            <span class="info-value">{{ pet.gender === 'male' ? '公' : pet.gender === 'female' ? '母' : '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('info.birthday') }}</span>
            <span class="info-value">{{ pet.birthday }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('info.color') }}</span>
            <span class="info-value">{{ pet.color || '-' }}</span>
          </div>
          <div v-if="pet.adopted_at" class="info-item">
            <span class="info-label">{{ t('info.adoptedAt') }}</span>
            <span class="info-value">{{ pet.adopted_at }}</span>
          </div>
          <div v-if="pet.passed_at" class="info-item">
            <span class="info-label">{{ t('detail.passedDate') }}</span>
            <span class="info-value">{{ pet.passed_at }}</span>
          </div>
          <div v-if="pet.note" class="info-item full-width">
            <span class="info-label">{{ t('info.note') }}</span>
            <span class="info-value">{{ pet.note }}</span>
          </div>
        </div>

        <div class="info-actions">
          <button class="btn btn-secondary" @click="showEditForm = true">{{ t('detail.editInfo') }}</button>
          <button v-if="!pet.passed_at" class="btn btn-ghost" @click="markPassed">🌈 {{ pet.species === 'cat' ? t('detail.markPassed') : t('detail.markPassed.other') }}</button>
          <button v-else class="btn btn-ghost" @click="clearPassed">{{ t('detail.clearPassed') }}</button>
          <button class="btn btn-danger" @click="deletePet">{{ t('form.delete') }}</button>
        </div>
      </div>

      <!-- Weight tab -->
      <div v-if="activeTab === 'weight'" class="tab-content">
        <div class="tab-header">
          <h3>{{ t('weight.title') }}</h3>
          <button class="btn btn-primary btn-sm" @click="showWeightForm = true">{{ t('weight.record') }}</button>
        </div>
        <WeightChart :weights="weights" />
        <div class="weight-list">
          <div v-for="w in weights" :key="w.id" class="weight-item">
            <span class="weight-date">{{ w.recorded_at.slice(0, 10) }}</span>
            <span class="weight-value">{{ w.weight.toFixed(2) }} kg</span>
            <button class="btn-icon" @click="async () => { await weightsApi.delete(pet.id, w.id); weights = await weightsApi.list(pet.id) }">×</button>
          </div>
        </div>
      </div>

      <!-- Health tab -->
      <div v-if="activeTab === 'health'" class="tab-content">
        <div class="tab-header">
          <h3>{{ t('health.title') }}</h3>
          <button class="btn btn-primary btn-sm" @click="showHealthForm = true">添加记录</button>
        </div>
        <HealthTimeline :records="healthRecords" @refresh="async () => { healthRecords = await healthApi.list(pet.id) }" />
      </div>

      <!-- Photos tab -->
      <div v-if="activeTab === 'photos'" class="tab-content">
        <PhotoGallery :photos="photos" :pet-id="pet.id" :birthday="pet.birthday" @refresh="async () => { photos = await photosApi.list(pet.id) }" />
      </div>
    </div>

    <!-- Edit form modal -->
    <Transition name="overlay">
      <div v-if="showEditForm" class="mini-overlay" @click.self="showEditForm = false">
        <PetEditForm :pet="pet" @close="showEditForm = false" @saved="onPetSaved" />
      </div>
    </Transition>

    <!-- Weight form modal -->
    <Transition name="overlay">
      <div v-if="showWeightForm" class="mini-overlay" @click.self="showWeightForm = false">
        <WeightForm :pet-id="pet.id" @close="showWeightForm = false" @saved="onWeightAdded" />
      </div>
    </Transition>

    <!-- Health form modal -->
    <Transition name="overlay">
      <div v-if="showHealthForm" class="mini-overlay" @click.self="showHealthForm = false">
        <HealthForm :pet-id="pet.id" @close="showHealthForm = false" @saved="onHealthAdded" />
      </div>
    </Transition>

    <!-- Avatar cropper -->
    <AvatarCropper v-if="cropSrc" :src="cropSrc" @confirm="onCropConfirm" @cancel="onCropCancel" />
  </div>
</template>

<style scoped>
.detail-panel {
  background: var(--bg-elevated);
  border-radius: 1.5rem;
  width: 100%;
  max-width: 640px;
  height: 80vh;
  display: flex;
  flex-direction: column;
  position: relative;
  animation: scaleIn var(--duration-normal) var(--ease-spring);
  box-shadow: 0 8px 60px rgba(0,0,0,0.15);
}

.close-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  z-index: 10;
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.625rem;
  border: none;
  background: var(--bg-subtle);
  color: var(--fg-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--duration-fast) var(--ease-spring);
}
.close-btn:hover { background: var(--bg-muted); transform: scale(0.95); }

/* Header */
.detail-header {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.5rem;
  border-bottom: 1px solid var(--bg-subtle);
}

.avatar-upload-area {
  position: relative;
  width: 4.5rem;
  height: 4.5rem;
  border-radius: 1.25rem;
  overflow: hidden;
  flex-shrink: 0;
  cursor: pointer;
}

.detail-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detail-avatar-placeholder {
  width: 100%;
  height: 100%;
  background: var(--bg-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-spring);
}

.avatar-upload-area:hover .avatar-overlay { opacity: 1; }

.avatar-input { display: none; }

.detail-name-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.detail-name {
  font-size: var(--text-xl);
  font-weight: 800;
  letter-spacing: -0.02em;
}

.species-badge-sm { font-size: 1.25rem; }

.detail-breed, .detail-color {
  font-size: var(--text-sm);
  color: var(--fg-secondary);
  margin-top: 0.25rem;
}

/* Tabs */
.detail-tabs {
  display: flex;
  border-bottom: 1px solid var(--bg-subtle);
  padding: 0 1.5rem;
}

.tab-btn {
  flex: 1;
  padding: 0.875rem 0;
  font-family: var(--font-family);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--fg-tertiary);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-spring);
}

.tab-btn.active {
  color: var(--fg-primary);
  border-bottom-color: var(--fg-primary);
}

.tab-btn:hover:not(.active) { color: var(--fg-secondary); }

/* Content */
.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}

.tab-content { animation: fadeIn var(--duration-fast) var(--ease-spring); }

.tab-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.tab-header h3 {
  font-size: var(--text-lg);
  font-weight: 700;
}

.btn-sm { padding: 0.5rem 0.875rem; font-size: var(--text-xs); }

/* Info tab */
.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.info-item { display: flex; flex-direction: column; gap: 0.25rem; }
.info-item.full-width { grid-column: 1 / -1; }

.info-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--fg-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: var(--text-base);
  color: var(--fg-primary);
  font-weight: 500;
}

.info-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid var(--bg-subtle);
}

/* Weight list */
.weight-list { margin-top: 1rem; }

.weight-item {
  display: flex;
  align-items: center;
  padding: 0.625rem 0;
  border-bottom: 1px solid var(--bg-subtle);
  font-size: var(--text-sm);
}

.weight-date { flex: 1; color: var(--fg-secondary); }
.weight-value { font-weight: 600; margin-right: 0.5rem; }

.btn-icon {
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  background: none;
  color: var(--fg-tertiary);
  cursor: pointer;
  border-radius: 0.375rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all var(--duration-fast);
}
.btn-icon:hover { background: var(--danger-soft); color: var(--danger); }

/* Mini overlay for forms */
.mini-overlay {
  position: absolute;
  inset: 0;
  z-index: 20;
  background: rgba(120, 113, 108, 0.35);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 1.5rem;
}

/* Mobile: bottom sheet */
@media (max-width: 768px) {
  .detail-panel {
    max-width: 100%;
    height: 85vh;
    border-radius: 1.25rem 1.25rem 0 0;
    align-self: flex-end;
    margin-bottom: 0;
  }
}
</style>
