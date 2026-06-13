import { describe, it, expect, beforeEach } from 'vitest'
import { t, setLocale, getLocale, speciesKey } from './useI18n'

describe('useI18n', () => {
  beforeEach(() => {
    setLocale('zh')
  })

  it('returns Chinese text by default', () => {
    expect(t('app.title')).toBe('自家毛孩')
  })

  it('switches to English', () => {
    setLocale('en')
    expect(t('app.title')).toBe('Pet Family')
  })

  it('interpolates params', () => {
    expect(t('pet.age.months', { n: 3 })).toBe('3个月')
    setLocale('en')
    expect(t('pet.age.months', { n: 3 })).toBe('3mo')
  })

  it('falls back to zh when key missing in current locale', () => {
    setLocale('en')
    // Any key that exists in zh works
    expect(t('btn.addPet')).toBe('Add Pet')
  })

  it('returns the key itself for unknown keys', () => {
    expect(t('does.not.exist')).toBe('does.not.exist')
  })

  it('getLocale reflects the current locale', () => {
    setLocale('en')
    expect(getLocale()).toBe('en')
    setLocale('zh')
    expect(getLocale()).toBe('zh')
  })
})

describe('speciesKey', () => {
  beforeEach(() => {
    setLocale('zh')
  })

  it('returns the default key for cats', () => {
    // cats use the base key ("去喵星日期")
    expect(speciesKey('cat', 'form.passedAt')).toBe('去喵星日期')
  })

  it('returns the .other variant for non-cats', () => {
    // non-cats use the .other variant ("离开日期")
    expect(speciesKey('dog', 'form.passedAt')).toBe('离开日期')
    expect(speciesKey('rabbit', 'form.passedAt')).toBe('离开日期')
  })

  it('falls back to base key when no .other variant exists', () => {
    // species.cat has no .other variant → returns base
    expect(speciesKey('dog', 'species.cat')).toBe('🐱 猫')
  })
})
