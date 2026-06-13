import { watch, nextTick, onBeforeUnmount, type Ref } from 'vue'

const FOCUSABLE = [
  'a[href]',
  'button:not([disabled])',
  'textarea',
  'input',
  'select',
  '[tabindex]:not([tabindex="-1"])',
].join(',')

/**
 * Trap keyboard focus inside `containerRef` while `activeRef` is true.
 * - On activation: focus the first focusable element.
 * - Tab / Shift+Tab cycle within the container.
 * - Escape invokes `onEscape` (typically closes the dialog).
 *
 * Usage:
 *   const { onKeydown } = useFocusTrap(panelRef, isOpen, () => emit('close'))
 *   // then bind @keydown="onKeydown" on the dialog root
 */
export function useFocusTrap(
  containerRef: Ref<HTMLElement | null>,
  activeRef: Ref<boolean>,
  onEscape: () => void,
) {
  function getFocusable(): HTMLElement[] {
    if (!containerRef.value) return []
    return Array.from(containerRef.value.querySelectorAll<HTMLElement>(FOCUSABLE)).filter(
      (el) => el.offsetParent !== null || el === document.activeElement,
    )
  }

  function focusFirst() {
    nextTick(() => {
      const els = getFocusable()
      els[0]?.focus()
    })
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault()
      onEscape()
      return
    }
    if (e.key !== 'Tab') return
    const els = getFocusable()
    if (els.length === 0) return
    const first = els[0]
    const last = els[els.length - 1]
    const active = document.activeElement as HTMLElement

    if (e.shiftKey) {
      if (active === first || !containerRef.value?.contains(active)) {
        e.preventDefault()
        last.focus()
      }
    } else {
      if (active === last) {
        e.preventDefault()
        first.focus()
      }
    }
  }

  const stop = watch(activeRef, (active) => {
    if (active) focusFirst()
  })

  onBeforeUnmount(stop)

  return { onKeydown, focusFirst }
}
