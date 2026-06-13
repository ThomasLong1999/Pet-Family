import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import type { WeightRecord } from '../types'

// Mock vue-echarts to avoid canvas rendering in happy-dom (no real 2d context).
vi.mock('vue-echarts', () => ({
  default: {
    name: 'VChart',
    props: { option: Object },
    template: '<div class="vchart-mock" />',
  },
}))

// Import after mock is registered
const WeightChart = (await import('./WeightChart.vue')).default

function makeRecord(weight: number, date: string): WeightRecord {
  return {
    id: `${date}-${weight}`,
    pet_id: 'pet-1',
    weight,
    recorded_at: date,
    note: null,
    created_at: date,
  }
}

describe('WeightChart', () => {
  it('shows empty state when no weights', () => {
    const wrapper = mount(WeightChart, { props: { weights: [] } })
    expect(wrapper.text()).toContain('暂无体重记录')
    expect(wrapper.find('.empty-chart').exists()).toBe(true)
    expect(wrapper.find('.vchart-mock').exists()).toBe(false)
  })

  it('renders chart when weights exist', () => {
    const weights = [
      makeRecord(4.0, '2024-01-01'),
      makeRecord(4.5, '2024-02-01'),
    ]
    const wrapper = mount(WeightChart, { props: { weights } })
    expect(wrapper.find('.vchart-mock').exists()).toBe(true)
    expect(wrapper.find('.empty-chart').exists()).toBe(false)
  })
})
