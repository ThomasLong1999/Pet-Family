import { request, upload } from './client'
import type { Photo } from '../types'

export const photosApi = {
  list: (petId: string) => request<Photo[]>(`/pets/${petId}/photos`),

  upload: (petId: string, file: File, ageGroup: string, caption?: string) => {
    const form = new FormData()
    form.append('photo', file)
    form.append('age_group', ageGroup)
    if (caption) form.append('caption', caption)
    return upload<Photo>(`/pets/${petId}/photos`, form)
  },

  delete: (petId: string, photoId: string) =>
    request<void>(`/pets/${petId}/photos/${photoId}`, { method: 'DELETE' }),
}
