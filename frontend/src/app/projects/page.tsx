'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import { projectsApi } from '@/api'
import { Project } from '@/types'

export default function ProjectsPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()
  const [projects, setProjects] = useState<Project[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated) {
      loadProjects()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, page])

  const loadProjects = async () => {
    setLoading(true)
    try {
      const result = await projectsApi.list({ keyword, page, pageSize })
      setProjects(result.items || [])
      setTotal(result.total)
    } catch (error) {
      console.error('Failed to load projects:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setPage(1)
    loadProjects()
  }

  const handleDelete = async (projectNo: string) => {
    if (!confirm('确定要删除这个项目吗？')) return
    try {
      await projectsApi.delete(projectNo)
      loadProjects()
    } catch (error) {
      alert('删除失败')
    }
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      TODO: 'bg-gray-100 text-gray-800',
      IN_PROGRESS: 'bg-blue-100 text-blue-800',
      COMPLETED: 'bg-green-100 text-green-800',
      CANCELLED: 'bg-red-100 text-red-800',
    }
    const labels: Record<string, string> = {
      TODO: '待开始',
      IN_PROGRESS: '进行中',
      COMPLETED: '已完成',
      CANCELLED: '已取消',
    }
    return (
      <span className={`px-2 py-1 text-xs font-medium rounded-full ${colors[status] || ''}`}>
        {labels[status] || status}
      </span>
    )
  }

  const getPriorityBadge = (priority: string) => {
    const colors: Record<string, string> = {
      LOW: 'text-gray-500',
      MEDIUM: 'text-yellow-600',
      HIGH: 'text-orange-600',
      URGENT: 'text-red-600',
    }
    const labels: Record<string, string> = {
      LOW: '低',
      MEDIUM: '中',
      HIGH: '高',
      URGENT: '紧急',
    }
    return <span className={`text-xs font-medium ${colors[priority] || ''}`}>{labels[priority] || priority}</span>
  }

  if (isLoading) {
    return <div className="min-h-screen flex items-center justify-center">加载中...</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link href="/dashboard" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
              <span className="ml-4 text-gray-500">/ 项目管理</span>
            </div>
            <div className="flex items-center">
              <Link href="/projects/create" className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700">
                新建项目
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="mb-6 flex gap-4">
          <input
            type="text"
            placeholder="搜索项目..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
          <button
            onClick={handleSearch}
            className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700"
          >
            搜索
          </button>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">项目编号</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">优先级</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">负责人</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center text-gray-500">加载中...</td>
                </tr>
              ) : projects.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
              ) : (
                projects.map((project) => (
                  <tr key={project.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{project.projectNo}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <Link href={`/projects/${project.projectNo}`} className="text-indigo-600 hover:text-indigo-900">
                        {project.name}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(project.status)}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{getPriorityBadge(project.priority)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{project.ownerNickname || '-'}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(project.createdAt).toLocaleDateString('zh-CN')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <Link href={`/projects/${project.projectNo}/edit`} className="text-indigo-600 hover:text-indigo-900 mr-4">编辑</Link>
                      <button onClick={() => handleDelete(project.projectNo)} className="text-red-600 hover:text-red-900">删除</button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div className="mt-4 flex justify-between items-center">
          <span className="text-sm text-gray-500">共 {total} 条</span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
            >
              上一页
            </button>
            <span className="px-3 py-1 text-sm">第 {page} 页</span>
            <button
              onClick={() => setPage(p => p + 1)}
              disabled={projects.length < pageSize}
              className="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
            >
              下一页
            </button>
          </div>
        </div>
      </main>
    </div>
  )
}
