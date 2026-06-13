import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Photo } from '../types'
import { photosApi } from '../api/photos'

export const usePhotosStore = defineStore('photos', () => {
  const photos = ref<Photo[]>([])
  const loading = ref(false)

  async function fetchPhotos(petId: string) {
    loading.value = true
    try {
      photos.value = await photosApi.list(petId)
    } finally {
      loading.value = false
    }
  }

  async function uploadPhoto(petId: string, file: File, caption?: string) {
    const photo = await photosApi.upload(petId, file, caption)
    photos.value.push(photo)
    return photo
  }

  async function deletePhoto(petId: string, photoId: string) {
    await photosApi.delete(petId, photoId)
    photos.value = photos.value.filter(p => p.id !== photoId)
  }

  // Group photos by age_group
  function groupedByAge() {
    const groups: Record<string, Photo[]> = {}
    for (const photo of photos.value) {
      if (!groups[photo.age_group]) groups[photo.age_group] = []
      groups[photo.age_group].push(photo)
    }
    return groups
  }

  return { photos, loading, fetchPhotos, uploadPhoto, deletePhoto, groupedByAge }
})
