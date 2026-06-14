import EditProjectClient from './EditProjectClient'

export function generateStaticParams() {
  return [{ projectNo: '_' }]
}

export default function EditProjectPage() {
  return <EditProjectClient />
}