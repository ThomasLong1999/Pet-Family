import { computed } from 'vue'
import { t } from './useI18n'

/** Total whole months between a birthday and `now`. */
function monthsSince(birthday: string, now = new Date()): number {
  const bd = new Date(birthday)
  if (Number.isNaN(bd.getTime())) return 0
  return (now.getFullYear() - bd.getFullYear()) * 12 + (now.getMonth() - bd.getMonth())
}

/**
 * Human-readable, i18n-aware age label, e.g. "< 1个月", "3个月", "2岁".
 * Returns '' when birthday is missing.
 */
export function ageLabel(birthday: string | null | undefined): string {
  if (!birthday) return ''
  const months = monthsSince(birthday)
  if (months < 1) return t('pet.age.lessThan1')
  if (months < 12) return t('pet.age.months', { n: months })
  const years = Math.floor(months / 12)
  const remainMonths = months % 12
  return remainMonths === 0
    ? t('pet.age.yearsRound', { y: years })
    : t('pet.age.years', { y: years, m: remainMonths })
}

/**
 * Composable wrapper for reactive age labels.
 * `passRange` — if set, shows "birthday — passed_at" instead of an age.
 */
export function usePetAge() {
  const ageText = (birthday: string | null | undefined, passedAt?: string | null) =>
    computed(() => {
      if (passedAt) return `${birthday ?? ''} — ${passedAt}`
      return ageLabel(birthday)
    })
  return { ageText, ageLabel }
}
