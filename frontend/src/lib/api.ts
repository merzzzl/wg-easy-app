import { useRawInitData } from '@tma.js/sdk-react'

export type UserStatus = 'pending' | string

export interface User {
  id: number
  telegram_id: number
  username: string
  language_code: string
  chat_id: number
  status: UserStatus
  created_at: string
  updated_at: string
}

export interface Tunnel {
  id: number
  user_id: number
  wg_client_name: string
  wg_client_id: string
  created_at: string
}

export interface MeResponse {
  user: User
  max_tunnels: number
  used_tunnels: number
}

class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, initDataRaw: string | undefined, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      'tg-token': initDataRaw ?? '',
      ...(init?.headers ?? {}),
    },
  })

  if (!response.ok) {
    let message = 'Server error'

    try {
      const payload = (await response.json()) as { error?: string }
      message = payload.error ?? message
    } catch {
      message = response.statusText || message
    }

    throw new ApiError(response.status, message)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export function useApi() {
  const initDataRaw = useRawInitData()

  return {
    getMe: () => request<MeResponse>('/api/v1/me', initDataRaw),
    listTunnels: () => request<Tunnel[]>('/api/v1/tunnels', initDataRaw),
    createTunnel: () => request<Tunnel>('/api/v1/tunnels', initDataRaw, { method: 'POST' }),
    deleteTunnel: (tunnelId: number) => request<{ ok: boolean }>(`/api/v1/tunnels/${tunnelId}`, initDataRaw, { method: 'DELETE' }),
    getTunnelQR: (tunnelId: number) => request<{ svg: string }>(`/api/v1/tunnels/${tunnelId}/qr`, initDataRaw),
    sendTunnelConfig: (tunnelId: number) => request<{ ok: boolean }>(`/api/v1/tunnels/${tunnelId}/config`, initDataRaw),
  }
}

export { ApiError }
