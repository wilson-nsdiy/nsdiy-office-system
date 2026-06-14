import ArticleVersionsClient from './ArticleVersionsClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function ArticleVersionsPage() {
  return <ArticleVersionsClient />
}