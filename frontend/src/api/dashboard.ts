import { request } from './client'
import type { Dashboard } from '../types'

export const dashboardApi = {
  get: () => request<Dashboard>('/dashboard'),
}
