import { describe, it, expect } from 'vitest'
import { smartDate } from './smartDate'

describe('smartDate', () => {
  it('passes through already-formatted dates', () => {
    expect(smartDate('2026-06-14')).toBe('2026-06-14')
  })

  it('parses YYYYMMDD (8 digits)', () => {
    expect(smartDate('20260614')).toBe('2026-06-14')
  })

  it('parses YYYYMDD (7 digits) by padding month', () => {
    expect(smartDate('2026614')).toBe('2026-06-14')
  })

  it('parses YYMMDD (6 digits) as 20YYMMDD', () => {
    expect(smartDate('260614')).toBe('2026-06-14')
  })

  it('normalizes slash separators', () => {
    expect(smartDate('2026/6/14')).toBe('2026-06-14')
  })

  it('normalizes dot separators', () => {
    expect(smartDate('2026.6.14')).toBe('2026-06-14')
  })

  it('normalizes single-digit month/day with dashes', () => {
    expect(smartDate('2026-6-14')).toBe('2026-06-14')
  })

  it('returns input unchanged when unparseable', () => {
    expect(smartDate('hello')).toBe('hello')
  })
})
