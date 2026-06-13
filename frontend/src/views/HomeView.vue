<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { dashboardApi } from '../api/dashboard'
import { petsApi } from '../api/pets'
import { t, setLocale, isZh } from '../composables/useI18n'
import PetCard from '../components/PetCard.vue'
import AddPetCard from '../components/AddPetCard.vue'
import PetDetail from '../components/PetDetail.vue'
import PetForm from '../components/PetForm.vue'
import type { Pet, PetSummary, Dashboard } from '../types'

const dashboard = ref<Dashboard | null>(null)
const pets = ref<Pet[]>([])
const loading = ref(false)

const selectedPet = ref<Pet | null>(null)
const showAddForm = ref(false)
const showDetail = ref(false)

function toggleLocale() {
  setLocale(isZh.value ? 'en' : 'zh')
}

onMounted(async () => {
  loading.value = true
  try {
    const [dash, petList] = await Promise.all([dashboardApi.get(), petsApi.list()])
    dashboard.value = dash
    pets.value = petList
  } finally {
    loading.value = false
  }
})

function getPetSummary(petId: string): PetSummary | null {
  return dashboard.value?.pets.find(p => p.id === petId) ?? null
}

function openDetail(pet: Pet) { selectedPet.value = pet; showDetail.value = true }
function closeDetail() { showDetail.value = false; selectedPet.value = null }
function openAddForm() { showAddForm.value = true }
function closeAddForm() { showAddForm.value = false }

async function handlePetAdded() { showAddForm.value = false; await refreshData() }
async function handlePetUpdated() { showDetail.value = false; selectedPet.value = null; await refreshData() }

async function refreshData() {
  const [dash, petList] = await Promise.all([dashboardApi.get(), petsApi.list()])
  dashboard.value = dash
  pets.value = petList
}
</script>

<template>
  <div class="home">
    <header class="home-header">
      <div class="header-inner">
        <div class="header-brand">
          <span class="header-icons">🐱🐶🐹🐰</span>
          <h1 class="header-title">{{ t('app.title') }}</h1>
        </div>
        <div class="header-actions">
          <button
            class="btn btn-primary"
            :aria-label="t('btn.addPet')"
            @click="openAddForm"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true">
              <line x1="8" y1="2" x2="8" y2="14" /><line x1="2" y1="8" x2="14" y2="8" />
            </svg>
            {{ t('btn.addPet') }}
          </button>
          <button class="btn btn-ghost lang-btn" :aria-label="isZh ? 'English' : '中文'" @click="toggleLocale">{{ t('lang.toggle') }}</button>
        </div>
      </div>
    </header>

    <main class="pet-grid-container">
      <div class="pet-grid">
        <PetCard
v-for="(pet, index) in pets" :key="pet.id"
          :pet="pet" :summary="getPetSummary(pet.id)"
          :reminders="(dashboard?.reminders ?? []).filter(r => r.pet_id === pet.id)"
          :style="{ animationDelay: `${index * 60}ms` }" class="pet-card-animate"
          @click="openDetail(pet)" @avatar-changed="refreshData" />
        <AddPetCard @click="openAddForm" />
      </div>
      <div v-if="pets.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">🐾</div>
        <h2>{{ t('empty.title') }}</h2>
        <p>{{ t('empty.desc') }}</p>
      </div>
    </main>

    <Transition name="overlay">
      <div
        v-if="showDetail && selectedPet"
        class="overlay"
        role="dialog"
        aria-modal="true"
        :aria-label="selectedPet.name"
        @click.self="closeDetail"
        @keydown.esc="closeDetail"
      >
        <PetDetail :pet="selectedPet" @close="closeDetail" @updated="handlePetUpdated" />
      </div>
    </Transition>

    <Transition name="overlay">
      <div
        v-if="showAddForm"
        class="overlay"
        role="dialog"
        aria-modal="true"
        :aria-label="t('form.addPet')"
        @click.self="closeAddForm"
        @keydown.esc="closeAddForm"
      >
        <PetForm @close="closeAddForm" @saved="handlePetAdded" />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.home { min-height: 100dvh; background: var(--bg-base); }
.home-header {
  position: sticky; top: 0; z-index: 100;
  background: rgba(250,250,249,0.85); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0,0,0,0.04);
}
.header-inner {
  max-width: 1200px; margin: 0 auto; padding: 1rem 2rem;
  display: flex; align-items: center; justify-content: space-between;
}
.header-brand { display: flex; align-items: center; gap: 0.75rem; }
.header-icons { font-size: 1.25rem; letter-spacing: 0.125rem; }
.header-title { font-size: var(--text-xl); font-weight: 700; color: var(--fg-primary); letter-spacing: -0.02em; }
.header-actions { display: flex; align-items: center; gap: 0.5rem; }
.lang-btn { font-weight: 700; font-size: var(--text-xs); min-width: 3rem; }
.pet-grid-container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
.pet-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
.pet-card-animate { animation: fadeInUp var(--duration-slow) var(--ease-spring) both; }
.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 6rem 2rem; text-align: center; animation: fadeInUp var(--duration-slow) var(--ease-spring); }
.empty-icon { font-size: 4rem; margin-bottom: 1.5rem; opacity: 0.6; }
.empty-state h2 { font-size: var(--text-xl); font-weight: 700; margin-bottom: 0.5rem; color: var(--fg-primary); }
.empty-state p { color: var(--fg-tertiary); font-size: var(--text-base); }
.overlay { position: fixed; inset: 0; z-index: 200; background: rgba(0,0,0,0.3); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; padding: 1rem; }
.overlay-enter-active { transition: opacity var(--duration-normal) var(--ease-spring); }
.overlay-leave-active { transition: opacity var(--duration-fast) var(--ease-smooth); }
.overlay-enter-from, .overlay-leave-to { opacity: 0; }
@media (max-width: 768px) {
  .header-inner { padding: 0.75rem 1rem; }
  .pet-grid-container { padding: 1rem; }
  .pet-grid { grid-template-columns: repeat(2, 1fr); gap: 0.75rem; }
}
@media (max-width: 480px) { .pet-grid { grid-template-columns: 1fr; } }
</style>
