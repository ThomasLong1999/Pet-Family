import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Pet, CreatePetRequest, UpdatePetRequest } from '../types'
import { petsApi } from '../api/pets'

export const usePetsStore = defineStore('pets', () => {
  const pets = ref<Pet[]>([])
  const loading = ref(false)
  const currentPet = ref<Pet | null>(null)

  async function fetchPets() {
    loading.value = true
    try {
      pets.value = await petsApi.list()
    } finally {
      loading.value = false
    }
  }

  async function fetchPet(id: string) {
    loading.value = true
    try {
      currentPet.value = await petsApi.get(id)
    } finally {
      loading.value = false
    }
  }

  async function createPet(data: CreatePetRequest) {
    const pet = await petsApi.create(data)
    pets.value.unshift(pet)
    return pet
  }

  async function updatePet(id: string, data: UpdatePetRequest) {
    const pet = await petsApi.update(id, data)
    const idx = pets.value.findIndex(p => p.id === id)
    if (idx >= 0) pets.value[idx] = pet
    if (currentPet.value?.id === id) currentPet.value = pet
    return pet
  }

  async function deletePet(id: string) {
    await petsApi.delete(id)
    pets.value = pets.value.filter(p => p.id !== id)
    if (currentPet.value?.id === id) currentPet.value = null
  }

  async function uploadAvatar(id: string, file: File) {
    const result = await petsApi.uploadAvatar(id, file)
    const idx = pets.value.findIndex(p => p.id === id)
    if (idx >= 0) {
      pets.value[idx] = { ...pets.value[idx], avatar_url: result.avatar_url, dominant_color: result.dominant_color }
    }
    if (currentPet.value?.id === id) {
      currentPet.value = { ...currentPet.value, avatar_url: result.avatar_url, dominant_color: result.dominant_color }
    }
    return result
  }

  return { pets, loading, currentPet, fetchPets, fetchPet, createPet, updatePet, deletePet, uploadAvatar }
})
