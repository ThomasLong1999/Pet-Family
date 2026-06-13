/**
 * Smart date parsing — accepts flexible input and returns YYYY-MM-DD.
 * Supported formats:
 *   20260614   → 2026-06-14
 *   2026614    → 2026-06-14
 *   2026/6/14  → 2026-06-14
 *   2026-6-14  → 2026-06-14
 *   2026-06-14 → 2026-06-14
 */
export function smartDate(input: string): string {
  const s = input.trim().replace(/[/.]/g, '-')
  // Already valid format
  if (/^\d{4}-\d{2}-\d{2}$/.test(s)) return s

  // Try to extract year, month, day from digits
  const digits = s.replace(/\D/g, '')
  if (digits.length === 8) {
    // YYYYMMDD
    return `${digits.slice(0, 4)}-${digits.slice(4, 6)}-${digits.slice(6, 8)}`
  }
  if (digits.length === 7) {
    // YYYYMDD or YYYYMMD — assume YYYY M DD (single digit month)
    return `${digits.slice(0, 4)}-0${digits.slice(4, 5)}-${digits.slice(5, 7)}`
  }
  if (digits.length === 6) {
    // Could be YYMMDD — treat as 20YYMMDD
    return `20${digits.slice(0, 2)}-${digits.slice(2, 4)}-${digits.slice(4, 6)}`
  }

  // Fallback: try to parse with dashes
  const parts = s.split('-')
  if (parts.length === 3) {
    const y = parts[0].padStart(4, '20')
    const m = parts[1].padStart(2, '0')
    const d = parts[2].padStart(2, '0')
    return `${y}-${m}-${d}`
  }

  return input
}

/**
 * Vue directive-like helper: formats a date input value on blur.
 */
export function onDateBlur(e: Event): void {
  const input = e.target as HTMLInputElement
  if (!input.value) return
  const formatted = smartDate(input.value)
  if (formatted !== input.value) {
    input.value = formatted
    input.dispatchEvent(new Event('input', { bubbles: true }))
  }
}
