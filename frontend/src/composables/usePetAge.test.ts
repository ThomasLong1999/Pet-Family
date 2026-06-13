import { describe, it, expect, beforeEach } from 'vitest'
import { ageLabel } from './usePetAge'
import { setLocale } from './useI18n'

function monthsAgo(months: number): string {
  const d = new Date()
  d.setMonth(d.getMonth() - months)
  return d.toISOString().slice(0, 10)
}

describe('ageLabel', () => {
  beforeEach(() => setLocale('zh'))

  it('returns empty for missing birthday', () => {
    expect(ageLabel(null)).toBe('')
    expect(ageLabel(undefined)).toBe('')
    expect(ageLabel('')).toBe('')
  })

  it('returns "less than 1 month" for newborns', () => {
    const future = '2999-01-01'
    expect(ageLabel(future)).toBe('< 1个月')
  })

  it('formats months under a year', () => {
    setLocale('zh')
    expect(ageLabel(monthsAgo(3))).toBe('3个月')
    setLocale('en')
    expect(ageLabel(monthsAgo(3))).toBe('3mo')
  })

  it('formats whole years without remaining months', () => {
    setLocale('zh')
    expect(ageLabel(monthsAgo(24))).toBe('2岁')
    setLocale('en')
    expect(ageLabel(monthsAgo(24))).toBe('2y')
  })

  it('formats years with remaining months', () => {
    setLocale('zh')
    expect(ageLabel(monthsAgo(30))).toBe('2岁6个月')
    setLocale('en')
    expect(ageLabel(monthsAgo(30))).toBe('2y6mo')
  })

  it('handles invalid date gracefully', () => {
    // Invalid date → monthsSince returns 0 → treated as newborn
    expect(ageLabel('not-a-date')).toBe('< 1个月')
  })
})
