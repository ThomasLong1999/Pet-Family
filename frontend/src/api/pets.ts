import { request, upload } from './client'
import type { Pet, CreatePetRequest, UpdatePetRequest } from '../types'

export const petsApi = {
  list: () => request<Pet[]>('/pets'),

  get: (id: string) => request<Pet>(`/pets/${id}`),

  create: (data: CreatePetRequest) =>
    request<Pet>('/pets', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: string, data: UpdatePetRequest) =>
    request<Pet>(`/pets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: string) =>
    request<void>(`/pets/${id}`, { method: 'DELETE' }),

  uploadAvatar: (id: string, file: File) => {
    const form = new FormData()
    form.append('avatar', file)
    return upload<{ avatar_url: string; dominant_color: string }>(`/pets/${id}/avatar`, form)
  },
}
