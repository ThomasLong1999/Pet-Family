export interface Pet {
  id: string
  species: 'cat' | 'dog' | 'hamster' | 'rabbit'
  name: string
  breed: string
  gender: 'male' | 'female'
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
  species: string
  name: string
  breed: string
  gender: string
  birthday: string
  color: string
  adopted_at?: string | null
  note?: string | null
}

export interface UpdatePetRequest {
  species?: string
  name?: string
  breed?: string
  gender?: string
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
  type: 'vaccine' | 'deworming' | 'checkup'
  name: string
  date: string
  next_date: string | null
  note: string | null
  report_url: string
  created_at: string
  updated_at: string
}

export interface CreateHealthRequest {
  type: string
  name: string
  date: string
  next_date?: string | null
  note?: string | null
}

export interface UpdateHealthRequest {
  type?: string
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
  species: string
  avatar_url: string
  dominant_color: string
  latest_weight: number | null
  weight_trend: number | null
  passed_at: string | null
}

export interface HealthReminder {
  pet_id: string
  pet_name: string
  record_type: string
  name: string
  next_date: string
  days_left: number
}

export interface Dashboard {
  pets: PetSummary[]
  reminders: HealthReminder[]
}
