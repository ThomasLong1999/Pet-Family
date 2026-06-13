import { request } from './client'
import type { WeightRecord, CreateWeightRequest } from '../types'

export const weightsApi = {
  list: (petId: string) => request<WeightRecord[]>(`/pets/${petId}/weights`),

  create: (petId: string, data: CreateWeightRequest) =>
    request<WeightRecord>(`/pets/${petId}/weights`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  delete: (petId: string, recordId: string) =>
    request<void>(`/pets/${petId}/weights/${recordId}`, { method: 'DELETE' }),
}
