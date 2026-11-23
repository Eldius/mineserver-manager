const TOKEN_KEY = 'mineserver.token'

export interface Instance {
  id: string
  name: string
  path: string
  install_date: string
  running: boolean
  manageable: boolean
}

export interface LoginResponse {
  token: string
  expires_at: string
}

export interface InstallRequest {
  name: string
  version: string
  memory_limit: string
  motd?: string
  level_name?: string
  seed?: string
  rcon_enabled?: boolean
  rcon_password?: string
}

export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers)
  const token = getToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  if (init?.body) {
    headers.set('Content-Type', 'application/json')
  }

  const res = await fetch(path, { ...init, headers })

  if (res.status === 401) {
    clearToken()
    window.location.hash = '#/login'
    throw new ApiError(401, 'session expired')
  }

  const body: unknown = await res.json().catch(() => ({}))
  if (!res.ok) {
    const message =
      typeof body === 'object' && body !== null && 'error' in body
        ? String((body as { error: unknown }).error)
        : `request failed with status ${res.status}`
    throw new ApiError(res.status, message)
  }
  return body as T
}

export async function login(password: string): Promise<void> {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  })
  const body: unknown = await res.json().catch(() => ({}))
  if (!res.ok) {
    const message =
      typeof body === 'object' && body !== null && 'error' in body
        ? String((body as { error: unknown }).error)
        : `login failed with status ${res.status}`
    throw new ApiError(res.status, message)
  }
  setToken((body as LoginResponse).token)
}

export function listInstances(): Promise<{ instances: Instance[] }> {
  return apiFetch('/api/instances')
}

export function listVersions(): Promise<{ versions: string[] }> {
  return apiFetch('/api/versions')
}

export function installInstance(req: InstallRequest): Promise<{ name: string; path: string }> {
  return apiFetch('/api/instances', { method: 'POST', body: JSON.stringify(req) })
}

export function instanceAction(id: string, action: 'start' | 'stop' | 'backup'): Promise<unknown> {
  return apiFetch(`/api/instances/${id}/${action}`, { method: 'POST' })
}
