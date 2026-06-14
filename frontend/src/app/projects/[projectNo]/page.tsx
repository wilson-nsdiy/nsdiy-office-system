import ProjectDetailClient from './ProjectDetailClient'

export function generateStaticParams() {
  return [{ projectNo: '_' }]
}

export default function ProjectDetailPage() {
  return <ProjectDetailClient />
}