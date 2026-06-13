import { request, upload, API_BASE } from './client'
import type { HealthRecord, CreateHealthRequest, UpdateHealthRequest } from '../types'

export const healthApi = {
  list: (petId: string, type?: string) => {
    const params = type ? `?type=${type}` : ''
    return request<HealthRecord[]>(`/pets/${petId}/health${params}`)
  },

  create: (petId: string, data: CreateHealthRequest) =>
    request<HealthRecord>(`/pets/${petId}/health`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  createWithReport: (petId: string, data: CreateHealthRequest, reportFile?: File) => {
    const form = new FormData()
    form.append('type', data.type)
    form.append('name', data.name)
    form.append('date', data.date)
    if (data.next_date) form.append('next_date', data.next_date)
    if (data.note) form.append('note', data.note)
    if (reportFile) form.append('report', reportFile)
    return upload<HealthRecord>(`/pets/${petId}/health`, form)
  },

  update: (petId: string, recordId: string, data: UpdateHealthRequest) =>
    request<HealthRecord>(`/pets/${petId}/health/${recordId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (petId: string, recordId: string) =>
    request<void>(`/pets/${petId}/health/${recordId}`, { method: 'DELETE' }),
}
