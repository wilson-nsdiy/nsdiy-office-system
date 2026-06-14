import EditMediaContentClient from './EditMediaContentClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function EditMediaContentPage() {
  return <EditMediaContentClient />
}