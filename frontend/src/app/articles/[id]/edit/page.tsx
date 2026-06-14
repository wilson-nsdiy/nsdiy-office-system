import EditArticleClient from './EditArticleClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function EditArticlePage() {
  return <EditArticleClient />
}