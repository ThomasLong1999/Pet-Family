import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { WeightRecord, CreateWeightRequest } from '../types'
import { weightsApi } from '../api/weights'

export const useWeightsStore = defineStore('weights', () => {
  const records = ref<WeightRecord[]>([])
  const loading = ref(false)

  async function fetchWeights(petId: string) {
    loading.value = true
    try {
      records.value = await weightsApi.list(petId)
    } finally {
      loading.value = false
    }
  }

  async function addWeight(petId: string, data: CreateWeightRequest) {
    const record = await weightsApi.create(petId, data)
    records.value.unshift(record)
    return record
  }

  async function deleteWeight(petId: string, recordId: string) {
    await weightsApi.delete(petId, recordId)
    records.value = records.value.filter(r => r.id !== recordId)
  }

  return { records, loading, fetchWeights, addWeight, deleteWeight }
})
