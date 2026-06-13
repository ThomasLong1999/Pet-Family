import { ref } from 'vue'

export interface Toast {
  id: number
  message: string
  type: 'error' | 'success'
}

const toasts = ref<Toast[]>([])
let nextId = 0

function push(message: string, type: Toast['type'] = 'error', duration = 4000) {
  const id = ++nextId
  toasts.value.push({ id, message, type })
  setTimeout(() => dismiss(id), duration)
}

function dismiss(id: number) {
  toasts.value = toasts.value.filter((t) => t.id !== id)
}

/** Global toast API. Call `error(...)` / `success(...)` anywhere; render `<ToastContainer />` once. */
export function useToast() {
  return {
    toasts,
    error: (message: string) => push(message, 'error'),
    success: (message: string) => push(message, 'success', 2500),
    dismiss,
  }
}
