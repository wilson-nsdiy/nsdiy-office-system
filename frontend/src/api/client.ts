import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse, ApiError } from '@/types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ==================== Token Refresh State ====================

let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

function subscribeTokenRefresh(callback: (token: string) => void): void {
  refreshSubscribers.push(callback)
}

function onTokenRefreshed(token: string): void {
  refreshSubscribers.forEach((callback) => callback(token))
  refreshSubscribers = []
}

// ==================== Proactive Token Refresh ====================

const TOKEN_REFRESH_BUFFER = 120 * 1000 // 120 seconds before expiry
let tokenRefreshTimeoutId: ReturnType<typeof setTimeout> | null = null

export function scheduleTokenRefresh(expiresInSeconds: number): void {
  if (tokenRefreshTimeoutId) {
    clearTimeout(tokenRefreshTimeoutId)
    tokenRefreshTimeoutId = null
  }

  const expiresAtMs = Date.now() + expiresInSeconds * 1000
  localStorage.setItem('tokenExpiresAt', String(expiresAtMs))

  const refreshInMs = Math.max(0, expiresAtMs - Date.now() - TOKEN_REFRESH_BUFFER)

  if (refreshInMs <= 0) {
    performTokenRefresh()
    return
  }

  tokenRefreshTimeoutId = setTimeout(() => {
    performTokenRefresh()
  }, refreshInMs)
}

export function scheduleTokenRefreshAt(expiresAtMs: number): void {
  if (tokenRefreshTimeoutId) {
    clearTimeout(tokenRefreshTimeoutId)
    tokenRefreshTimeoutId = null
  }

  const refreshInMs = Math.max(0, expiresAtMs - Date.now() - TOKEN_REFRESH_BUFFER)

  if (refreshInMs <= 0) {
    performTokenRefresh()
    return
  }

  tokenRefreshTimeoutId = setTimeout(() => {
    performTokenRefresh()
  }, refreshInMs)
}

async function performTokenRefresh(): Promise<void> {
  const refreshToken = localStorage.getItem('refreshToken')
  if (!refreshToken) return

  try {
    const response = await axios.post(
      `${API_BASE_URL}/api/auth/refresh`,
      { refreshToken },
      { headers: { 'Content-Type': 'application/json' } }
    )

    const apiResponse = response.data as ApiResponse<{
      accessToken: string
      refreshToken: string
      expiresIn: number
    }>

    if (apiResponse.data) {
      const { accessToken, refreshToken: newRefreshToken, expiresIn } = apiResponse.data
      localStorage.setItem('accessToken', accessToken)
      localStorage.setItem('refreshToken', newRefreshToken)
      scheduleTokenRefresh(expiresIn)
    }
  } catch (error) {
    console.error('Proactive token refresh failed:', error)
  }
}

export function startProactiveRefresh(): void {
  const expiresAt = localStorage.getItem('tokenExpiresAt')
  if (expiresAt) {
    const expiresAtMs = parseInt(expiresAt, 10)
    if (!isNaN(expiresAtMs)) {
      scheduleTokenRefreshAt(expiresAtMs)
    }
  }
}

export function stopProactiveRefresh(): void {
  if (tokenRefreshTimeoutId) {
    clearTimeout(tokenRefreshTimeoutId)
    tokenRefreshTimeoutId = null
  }
}

// ==================== Helpers ====================

function getUserTimezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
  } catch {
    return 'UTC'
  }
}

function getLocale(): string {
  if (typeof navigator !== 'undefined') {
    return navigator.language || 'zh-CN'
  }
  return 'zh-CN'
}

// ==================== Request Interceptor ====================

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('accessToken')
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`
      }

      if (config.headers) {
        config.headers['Accept-Language'] = getLocale()
      }

      if (config.method === 'get') {
        if (!config.params) {
          config.params = {}
        }
        config.params.timezone = getUserTimezone()
      }
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// ==================== Response Interceptor ====================

apiClient.interceptors.response.use(
  (response) => {
    const apiResponse = response.data as ApiResponse<unknown>
    if (apiResponse && typeof apiResponse === 'object' && 'data' in apiResponse) {
      response.data = apiResponse.data
    }
    return response
  },
  async (error: AxiosError<ApiResponse<unknown>>) => {
    if (error.code === 'ERR_CANCELED' || axios.isCancel(error)) {
      return Promise.reject(error)
    }

    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response) {
      const { status, data } = error.response
      const apiData = (typeof data === 'object' && data !== null ? data : {}) as Record<string, any>

      if (status === 401 && !originalRequest._retry) {
        const refreshToken = localStorage.getItem('refreshToken')
        const isAuthEndpoint =
          originalRequest.url?.includes('/auth/login') ||
          originalRequest.url?.includes('/auth/register') ||
          originalRequest.url?.includes('/auth/refresh')

        if (refreshToken && !isAuthEndpoint) {
          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              subscribeTokenRefresh((newToken: string) => {
                if (newToken) {
                  originalRequest._retry = true
                  if (originalRequest.headers) {
                    originalRequest.headers.Authorization = `Bearer ${newToken}`
                  }
                  resolve(apiClient(originalRequest))
                } else {
                  reject({
                    status,
                    code: apiData.code,
                    message: apiData.message || 'Token refresh failed',
                  } as ApiError)
                }
              })
            })
          }

          originalRequest._retry = true
          isRefreshing = true

          try {
            const refreshResponse = await axios.post(
              `${API_BASE_URL}/api/auth/refresh`,
              { refreshToken },
              { headers: { 'Content-Type': 'application/json' } }
            )

            const refreshData = refreshResponse.data as ApiResponse<{
              accessToken: string
              refreshToken: string
              expiresIn: number
            }>

            if (refreshData.data) {
              const { accessToken, refreshToken: newRefreshToken, expiresIn } = refreshData.data
              localStorage.setItem('accessToken', accessToken)
              localStorage.setItem('refreshToken', newRefreshToken)
              scheduleTokenRefresh(expiresIn)

              onTokenRefreshed(accessToken)
              isRefreshing = false

              if (originalRequest.headers) {
                originalRequest.headers.Authorization = `Bearer ${accessToken}`
              }
              return apiClient(originalRequest)
            }

            throw new Error('Token refresh failed')
          } catch (refreshError) {
            onTokenRefreshed('')
            isRefreshing = false

            localStorage.removeItem('accessToken')
            localStorage.removeItem('refreshToken')
            localStorage.removeItem('tokenExpiresAt')
            stopProactiveRefresh()

            if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
              window.location.href = '/login'
            }

            return Promise.reject({
              status: 401,
              code: 'TOKEN_REFRESH_FAILED',
              message: 'Session expired. Please log in again.',
            } as ApiError)
          }
        }

        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('tokenExpiresAt')
        stopProactiveRefresh()

        if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
          window.location.href = '/login'
        }
      }

      return Promise.reject({
        status,
        code: apiData.code,
        reason: apiData.reason,
        message: apiData.message || apiData.error || error.message,
        metadata: apiData.metadata,
      } as ApiError)
    }

    return Promise.reject({
      status: 0,
      message: 'Network error. Please check your connection.',
    } as ApiError)
  }
)

export default apiClient
