'use client'

import { useEffect, ReactNode } from 'react'

interface ChunkErrorRecoveryProps {
  children: ReactNode
}

export function ChunkErrorRecovery({ children }: ChunkErrorRecoveryProps) {
  useEffect(() => {
    const handleChunkError = (error: ErrorEvent) => {
      const isChunkLoadError =
        error.message?.includes('Failed to fetch dynamically imported module') ||
        error.message?.includes('Loading chunk') ||
        error.message?.includes('Loading CSS chunk') ||
        error.message?.includes('ChunkLoadError')

      if (isChunkLoadError) {
        const reloadKey = 'chunk_reload_attempted'
        const lastReload = sessionStorage.getItem(reloadKey)
        const now = Date.now()

        if (!lastReload || now - parseInt(lastReload) > 10000) {
          sessionStorage.setItem(reloadKey, now.toString())
          console.warn('Chunk load error detected, reloading page...')
          window.location.reload()
        } else {
          console.error('Chunk load error persists. Please clear browser cache.')
        }
      }
    }

    window.addEventListener('error', handleChunkError)
    return () => window.removeEventListener('error', handleChunkError)
  }, [])

  return <>{children}</>
}
