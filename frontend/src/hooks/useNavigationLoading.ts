'use client'

import { useState, useCallback, useRef } from 'react'

interface UseNavigationLoadingReturn {
  isLoading: boolean
  startNavigation: () => void
  endNavigation: () => void
}

export function useNavigationLoading(): UseNavigationLoadingReturn {
  const [isLoading, setIsLoading] = useState(false)
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const startNavigation = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    setIsLoading(true)
  }, [])

  const endNavigation = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    timeoutRef.current = setTimeout(() => {
      setIsLoading(false)
    }, 100)
  }, [])

  return { isLoading, startNavigation, endNavigation }
}
