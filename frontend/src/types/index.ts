// Shared domain unions — used across Pet DTOs and components
export type Species = 'cat' | 'dog' | 'hamster' | 'rabbit'
export type Gender = 'male' | 'female'
export type HealthType = 'vaccine' | 'deworming' | 'checkup'

export interface Pet {
  id: string
  species: Species
  name: string
  breed: string
  gender: Gender
  birthday: string
  color: string
  avatar_url: string
  dominant_color: string
  adopted_at: string | null
  passed_at: string | null
  note: string | null
  created_at: string
  updated_at: string
}

export interface CreatePetRequest {
  species: Species
  name: string
  breed: string
  gender: Gender
  birthday: string
  color: string
  adopted_at?: string | null
  note?: string | null
}

export interface UpdatePetRequest {
  species?: Species
  name?: string
  breed?: string
  gender?: Gender
  birthday?: string
  color?: string
  adopted_at?: string | null
  passed_at?: string | null
  note?: string | null
}

export interface WeightRecord {
  id: string
  pet_id: string
  weight: number
  recorded_at: string
  note: string | null
  created_at: string
}

export interface CreateWeightRequest {
  weight: number
  recorded_at: string
  note?: string | null
}

export interface HealthRecord {
  id: string
  pet_id: string
  type: HealthType
  name: string
  date: string
  next_date: string | null
  note: string | null
  report_url: string
  created_at: string
  updated_at: string
}

export interface CreateHealthRequest {
  type: HealthType
  name: string
  date: string
  next_date?: string | null
  note?: string | null
}

export interface UpdateHealthRequest {
  type?: HealthType
  name?: string
  date?: string
  next_date?: string | null
  note?: string | null
}

export interface Photo {
  id: string
  pet_id: string
  url: string
  age_group: string
  caption: string | null
  created_at: string
}

export interface PetSummary {
  id: string
  name: string
  species: Species
  avatar_url: string
  dominant_color: string
  latest_weight: number | null
  weight_trend: number | null
  passed_at: string | null
}

export interface HealthReminder {
  pet_id: string
  pet_name: string
  record_type: HealthType
  name: string
  next_date: string
  days_left: number
}

export interface Dashboard {
  pets: PetSummary[]
  reminders: HealthReminder[]
}
