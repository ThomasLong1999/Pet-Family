const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1'

// Default request timeout (ms). Uploads are exempt since large files may take longer.
const DEFAULT_TIMEOUT = 15000

function withTimeout(signal?: AbortSignal | null, timeout = DEFAULT_TIMEOUT): AbortSignal {
  if (signal) return signal
  const controller = new AbortController()
  setTimeout(() => controller.abort(), timeout)
  return controller.signal
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
    signal: withTimeout(options?.signal),
  })

  if (res.status === 204) return undefined as T
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || `HTTP ${res.status}`)
  }

  return res.json()
}

async function upload<T>(path: string, formData: FormData, signal?: AbortSignal | null): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    body: formData,
    signal,
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || `HTTP ${res.status}`)
  }

  return res.json()
}

export { request, upload }
export { API_BASE }
