import type { Metadata } from 'next'
import '@fontsource-variable/inter'
import './globals.css'
import { Providers } from './providers'

export const metadata: Metadata = {
  title: 'OA-NSDIY - 工作室OA管理系统',
  description: '工作室OA管理系统',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
