import NewsDetailClient from './NewsDetailClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function NewsDetailPage() {
  return <NewsDetailClient />
}