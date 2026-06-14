'use client'

import React, { createContext, useContext, useState, useEffect, useCallback } from 'react'
import apiClient, {
  scheduleTokenRefresh,
  startProactiveRefresh,
  stopProactiveRefresh,
} from '@/api/client'
import { SafeUser, LoginRequest, LoginResponse, ApiError } from '@/types'

interface AuthContextType {
  user: SafeUser | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (data: LoginRequest) => Promise<void>
  logout: () => void
  refreshToken: () => Promise<void>
  getAuthError: (error: unknown) => ApiError
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<SafeUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const getAuthError = useCallback((error: unknown): ApiError => {
    if (error && typeof error === 'object' && 'status' in error) {
      return error as ApiError
    }
    if (error instanceof Error) {
      return { status: 0, message: error.message }
    }
    return { status: 0, message: 'An unknown error occurred' }
  }, [])

  const fetchUser = useCallback(async () => {
    try {
      const response = await apiClient.get('/api/auth/me')
      setUser(response.data as SafeUser)
    } catch (error) {
      setUser(null)
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('tokenExpiresAt')
      stopProactiveRefresh()
    }
  }, [])

  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('accessToken')
      if (token) {
        await fetchUser()
        startProactiveRefresh()
      }
      setIsLoading(false)
    }
    initAuth()

    return () => {
      stopProactiveRefresh()
    }
  }, [fetchUser])

  const login = async (data: LoginRequest) => {
    const response = await apiClient.post('/api/auth/login', data)
    const loginData = response.data as LoginResponse
    localStorage.setItem('accessToken', loginData.accessToken)
    localStorage.setItem('refreshToken', loginData.refreshToken)

    if (loginData.expiresIn) {
      scheduleTokenRefresh(loginData.expiresIn)
    }

    await fetchUser()
    startProactiveRefresh()
  }

  const logout = () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('tokenExpiresAt')
    stopProactiveRefresh()
    setUser(null)
    window.location.href = '/login'
  }

  const refreshToken = async () => {
    const refreshTokenValue = localStorage.getItem('refreshToken')
    if (!refreshTokenValue) {
      throw new Error('No refresh token')
    }

    const response = await apiClient.post('/api/auth/refresh', {
      refreshToken: refreshTokenValue,
    })
    const refreshData = response.data as LoginResponse
    localStorage.setItem('accessToken', refreshData.accessToken)
    localStorage.setItem('refreshToken', refreshData.refreshToken)

    if (refreshData.expiresIn) {
      scheduleTokenRefresh(refreshData.expiresIn)
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        logout,
        refreshToken,
        getAuthError,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
