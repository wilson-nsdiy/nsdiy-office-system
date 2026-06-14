import ArticleDetailClient from './ArticleDetailClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function ArticleDetailPage() {
  return <ArticleDetailClient />
}