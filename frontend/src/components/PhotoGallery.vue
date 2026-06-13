<script setup lang="ts">
import { computed, ref, onBeforeUnmount } from 'vue'
import type { Photo } from '../types'
import { photosApi } from '../api/photos'
import { t } from '../composables/useI18n'
import { useToast } from '../composables/useToast'

const props = defineProps<{
  photos: Photo[]
  petId: string
  birthday: string
}>()

const emit = defineEmits<{ refresh: [] }>()

// Sorted groups: 1m, 2m, ..., 11m, 12m, 1y, 2y, ..., 20y
const sortedGroups = computed(() => {
  const map: Record<string, Photo[]> = {}
  for (const photo of props.photos) {
    if (!map[photo.age_group]) map[photo.age_group] = []
    map[photo.age_group].push(photo)
  }
  // Sort each group's photos by created_at
  for (const key in map) {
    map[key].sort((a, b) => a.created_at.localeCompare(b.created_at))
  }
  // Sort keys: months first then years
  const keys = Object.keys(map)
  keys.sort((a, b) => {
    const aM = a.endsWith('m')
    const bM = b.endsWith('m')
    if (aM !== bM) return aM ? -1 : 1
    return parseInt(a) - parseInt(b)
  })
  return keys.map(k => ({ key: k, photos: map[k] }))
})

const ageGroupLabel = (group: string) => {
  if (group.endsWith('m')) return t('photos.ageGroup.month', { n: parseInt(group) })
  if (group.endsWith('y')) return t('photos.ageGroup.year', { n: parseInt(group) })
  return group
}

const ageGroupOptions = computed(() => {
  const options: { value: string; label: string }[] = []
  for (let m = 1; m <= 12; m++) {
    options.push({ value: `${m}m`, label: t('photos.ageGroup.month', { n: m }) })
  }
  for (let y = 1; y <= 20; y++) {
    options.push({ value: `${y}y`, label: t('photos.ageGroup.year', { n: y }) })
  }
  return options
})

const suggestedAgeGroup = computed(() => {
  if (!props.birthday) return '1y'
  const bd = new Date(props.birthday)
  const now = new Date()
  const months = (now.getFullYear() - bd.getFullYear()) * 12 + (now.getMonth() - bd.getMonth())
  if (months < 1) return '1m'
  if (months <= 12) return `${months}m`
  return `${Math.floor(months / 12)}y`
})

const showUploadDialog = ref(false)
const uploadFile = ref<File | null>(null)
const uploadPreview = ref<string | null>(null)
const selectedAgeGroup = ref('')
const uploadCaption = ref('')
const uploading = ref(false)

function openUploadDialog() {
  uploadFile.value = null
  uploadPreview.value = null
  selectedAgeGroup.value = suggestedAgeGroup.value
  uploadCaption.value = ''
  showUploadDialog.value = true
}

function onFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  // Revoke previous preview if re-selecting
  if (uploadPreview.value) URL.revokeObjectURL(uploadPreview.value)
  uploadFile.value = input.files[0]
  uploadPreview.value = URL.createObjectURL(uploadFile.value)
}

async function confirmUpload() {
  if (!uploadFile.value || !selectedAgeGroup.value) return
  uploading.value = true
  try {
    await photosApi.upload(props.petId, uploadFile.value, selectedAgeGroup.value, uploadCaption.value || undefined)
    showUploadDialog.value = false
    emit('refresh')
  } catch (err) {
    toastError(t('error.uploadFailed') + (err as Error).message)
  } finally {
    uploading.value = false
  }
}

async function deletePhoto(photo: Photo) {
  if (!confirm(t('confirm.deletePhoto'))) return
  try {
    await photosApi.delete(props.petId, photo.id)
    emit('refresh')
  } catch (err) {
    toastError(t('error.deleteFailed') + (err as Error).message)
  }
}

const { error: toastError } = useToast()
const viewingPhoto = ref<Photo | null>(null)

// Release object URL when the gallery is unmounted
onBeforeUnmount(() => {
  if (uploadPreview.value) URL.revokeObjectURL(uploadPreview.value)
})
</script>

<template>
  <div class="photo-gallery">
    <div class="tab-header">
      <h3>{{ t('photos.title') }}</h3>
      <button class="btn btn-primary btn-sm" @click="openUploadDialog">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <line x1="7" y1="1" x2="7" y2="13" /><line x1="1" y1="7" x2="13" y2="7" />
        </svg>
        {{ t('photos.upload') }}
      </button>
    </div>

    <div v-if="sortedGroups.length === 0" class="empty-state">
      <p>{{ t('photos.empty') }}</p>
    </div>

    <div v-for="group in sortedGroups" :key="group.key" class="age-group">
      <h4 class="age-group-title">{{ ageGroupLabel(group.key) }}</h4>
      <div class="photo-row">
        <div v-for="photo in group.photos" :key="photo.id" class="photo-thumb" role="button" tabindex="0" :aria-label="photo.caption || ''" @click="viewingPhoto = photo" @keydown.enter.prevent="viewingPhoto = photo" @keydown.space.prevent="viewingPhoto = photo">
          <img :src="photo.url" :alt="photo.caption || ''" loading="lazy" />
          <button class="photo-delete" :aria-label="t('confirm.deletePhoto')" @click.stop="deletePhoto(photo)">×</button>
        </div>
      </div>
    </div>

    <!-- Upload dialog -->
    <Transition name="overlay">
      <div v-if="showUploadDialog" class="upload-dialog-overlay" role="dialog" aria-modal="true" :aria-label="t('photos.dialog.title')" @click.self="showUploadDialog = false" @keydown.esc="showUploadDialog = false">
        <div class="upload-dialog">
          <h3>{{ t('photos.dialog.title') }}</h3>
          <div
            class="upload-file-area"
            role="button"
            tabindex="0"
            :aria-label="t('photos.dialog.select')"
            @click="($refs.fileInput as HTMLInputElement)?.click()"
            @keydown.enter.prevent="($refs.fileInput as HTMLInputElement)?.click()"
            @keydown.space.prevent="($refs.fileInput as HTMLInputElement)?.click()"
          >
            <img v-if="uploadPreview" :src="uploadPreview" class="upload-preview" />
            <div v-else class="upload-placeholder">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
                <path d="M12 5v14M5 12h14" />
              </svg>
              <span>{{ t('photos.dialog.select') }}</span>
            </div>
            <input ref="fileInput" type="file" accept="image/*" class="hidden-input" @change="onFileSelected" />
          </div>
          <div class="upload-field">
            <label class="form-label">{{ t('photos.ageGroup') }} *</label>
            <select v-model="selectedAgeGroup" class="form-input">
              <option v-for="opt in ageGroupOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>
          </div>
          <div class="upload-field">
            <label class="form-label">{{ t('photos.caption') }}</label>
            <input v-model="uploadCaption" class="form-input" :placeholder="t('photos.caption.ph')" />
          </div>
          <div class="upload-actions">
            <button class="btn btn-ghost" @click="showUploadDialog = false">{{ t('form.cancel') }}</button>
            <button class="btn btn-primary" :disabled="!uploadFile || !selectedAgeGroup || uploading" @click="confirmUpload">
              {{ uploading ? t('photos.uploading') : t('photos.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Full view -->
    <Transition name="overlay">
      <div v-if="viewingPhoto" class="photo-viewer" role="dialog" aria-modal="true" :aria-label="viewingPhoto.caption || ''" @click="viewingPhoto = null" @keydown.esc="viewingPhoto = null">
        <img :src="viewingPhoto.url" class="full-photo" :alt="viewingPhoto.caption || ''" loading="lazy" />
        <p v-if="viewingPhoto.caption" class="photo-caption">{{ viewingPhoto.caption }}</p>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.tab-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
.tab-header h3 { font-size: var(--text-lg); font-weight: 700; }
.btn-sm { padding: 0.5rem 0.875rem; font-size: var(--text-xs); }
.empty-state { padding: 2rem; text-align: center; color: var(--fg-tertiary); font-size: var(--text-sm); }
.age-group { margin-bottom: 1.25rem; }
.age-group-title { font-size: var(--text-sm); font-weight: 700; color: var(--fg-secondary); margin-bottom: 0.5rem; }
.photo-row { display: flex; gap: 0.5rem; overflow-x: auto; padding-bottom: 0.5rem; scroll-snap-type: x mandatory; }
.photo-thumb {
  flex-shrink: 0; width: 5rem; height: 5rem; border-radius: 0.625rem; overflow: hidden; cursor: pointer;
  position: relative; scroll-snap-align: start; transition: transform var(--duration-fast) var(--ease-spring);
}
.photo-thumb:hover { transform: scale(1.05); }
.photo-thumb img { width: 100%; height: 100%; object-fit: cover; }
.photo-delete {
  position: absolute; top: 0.25rem; right: 0.25rem; width: 1.25rem; height: 1.25rem; border: none;
  background: rgba(0,0,0,0.5); color: white; border-radius: 50%; font-size: 0.75rem; cursor: pointer;
  display: none; align-items: center; justify-content: center;
}
.photo-thumb:hover .photo-delete { display: flex; }

.upload-dialog-overlay {
  position: fixed; inset: 0; z-index: 300; background: rgba(0,0,0,0.3); backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
}
.upload-dialog {
  background: var(--bg-elevated); border-radius: 1.25rem; padding: 1.5rem; width: 100%; max-width: 360px;
  animation: scaleIn var(--duration-normal) var(--ease-spring);
}
.upload-dialog h3 { font-size: var(--text-lg); font-weight: 700; margin-bottom: 1rem; }
.upload-file-area {
  width: 100%; height: 10rem; border: 2px dashed var(--bg-muted); border-radius: 0.75rem;
  overflow: hidden; cursor: pointer; margin-bottom: 1rem; transition: border-color var(--duration-fast) var(--ease-spring);
}
.upload-file-area:hover { border-color: var(--fg-tertiary); }
.upload-preview { width: 100%; height: 100%; object-fit: cover; }
.upload-placeholder {
  width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center;
  justify-content: center; gap: 0.5rem; color: var(--fg-tertiary); font-size: var(--text-sm);
}
.upload-field { margin-bottom: 0.875rem; }
.upload-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1rem; }
.hidden-input { display: none; }

.photo-viewer {
  position: fixed; inset: 0; z-index: 300; background: rgba(0,0,0,0.9);
  display: flex; flex-direction: column; align-items: center; justify-content: center; cursor: pointer;
}
.full-photo { max-width: 90%; max-height: 80vh; border-radius: 0.75rem; object-fit: contain; }
.photo-caption { color: white; margin-top: 1rem; font-size: var(--text-base); }
</style>
