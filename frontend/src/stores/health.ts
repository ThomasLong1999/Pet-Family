import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HealthRecord, CreateHealthRequest, UpdateHealthRequest } from '../types'
import { healthApi } from '../api/health'

export const useHealthStore = defineStore('health', () => {
  const records = ref<HealthRecord[]>([])
  const loading = ref(false)

  async function fetchRecords(petId: string, type?: string) {
    loading.value = true
    try {
      records.value = await healthApi.list(petId, type)
    } finally {
      loading.value = false
    }
  }

  async function addRecord(petId: string, data: CreateHealthRequest) {
    const record = await healthApi.create(petId, data)
    records.value.unshift(record)
    return record
  }

  async function updateRecord(petId: string, recordId: string, data: UpdateHealthRequest) {
    const record = await healthApi.update(petId, recordId, data)
    const idx = records.value.findIndex(r => r.id === recordId)
    if (idx >= 0) records.value[idx] = record
    return record
  }

  async function deleteRecord(petId: string, recordId: string) {
    await healthApi.delete(petId, recordId)
    records.value = records.value.filter(r => r.id !== recordId)
  }

  return { records, loading, fetchRecords, addRecord, updateRecord, deleteRecord }
})
