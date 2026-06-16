'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import apiClient from '@/api/client'

export default function Home() {
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [checkingSetup, setCheckingSetup] = useState(true)

  useEffect(() => {
    const checkSetup = async () => {
      try {
        const response = await apiClient.get('/api/setup/status')
        const data = response.data as { needsSetup: boolean }
        if (data.needsSetup) {
          router.replace('/setup')
          return
        }
      } catch {
        // If setup endpoint fails, continue to normal flow
      }
      setCheckingSetup(false)
    }

    checkSetup()
  }, [router])

  useEffect(() => {
    if (checkingSetup || isLoading) return

    if (isAuthenticated) {
      router.replace('/dashboard/')
    } else {
      router.replace('/login/')
    }
  }, [checkingSetup, isLoading, isAuthenticated, router])

  return (
    <div className="min-h-screen flex items-center justify-center">
      加载中...
    </div>
  )
}