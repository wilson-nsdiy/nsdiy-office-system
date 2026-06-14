'use client'

import { AuthProvider } from '@/contexts/auth'
import { ChunkErrorRecovery } from '@/components/ChunkErrorRecovery'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ChunkErrorRecovery>
      <AuthProvider>{children}</AuthProvider>
    </ChunkErrorRecovery>
  )
}
