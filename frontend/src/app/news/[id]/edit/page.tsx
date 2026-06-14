import EditNewsClient from './EditNewsClient'

export function generateStaticParams() {
  return [{ id: '_' }]
}

export default function EditNewsPage() {
  return <EditNewsClient />
}