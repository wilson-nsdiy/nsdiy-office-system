'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth'
import { projectsApi, tasksApi } from '@/api'
import { Project, ProjectMember, Task } from '@/types'

export default function ProjectDetailClient() {
  const params = useParams()
  const router = useRouter()
  const { isAuthenticated, isLoading, user } = useAuth()
  const [project, setProject] = useState<Project | null>(null)
  const [members, setMembers] = useState<ProjectMember[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [activeTab, setActiveTab] = useState<'members' | 'tasks'>('members')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login/')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated && params.projectNo) {
      loadProject()
    }
  }, [isAuthenticated, params.projectNo])

  const loadProject = async () => {
    setLoading(true)
    try {
      const projectData = await projectsApi.getByProjectNo(params.projectNo as string)
      setProject(projectData)

      const membersData = await projectsApi.getMembers(params.projectNo as string)
      setMembers(membersData || [])

      const tasksResult = await tasksApi.list({ projectId: projectData.id })
      setTasks(tasksResult.items || [])
    } catch (error) {
      console.error('Failed to load project:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('确定要删除这个项目吗？')) return
    try {
      await projectsApi.delete(params.projectNo as string)
      router.push('/projects/')
    } catch (error) {
      alert('删除失败')
    }
  }

  const handleRemoveMember = async (userId: number) => {
    if (!confirm('确定要移除这个成员吗？')) return
    try {
      await projectsApi.removeMember(params.projectNo as string, userId)
      loadProject()
    } catch (error) {
      alert('移除失败')
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

  if (isLoading || loading) {
    return <div className="min-h-screen flex items-center justify-center">加载中...</div>
  }

  if (!project) {
    return <div className="min-h-screen flex items-center justify-center">项目不存在</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center h-16">
            <Link href="/dashboard/" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
            <span className="ml-4 text-gray-500">/ 项目详情</span>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{project.name}</h1>
              <p className="mt-1 text-sm text-gray-500">
                {project.projectNo} · {getStatusBadge(project.status)}
              </p>
            </div>
            <div className="flex gap-2">
              <Link
                href={`/projects/${project.projectNo}/edit/`}
                className="px-3 py-1 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
              >
                编辑
              </Link>
              <button
                onClick={handleDelete}
                className="px-3 py-1 border border-red-300 rounded-md text-sm text-red-600 hover:bg-red-50"
              >
                删除
              </button>
            </div>
          </div>

          {project.description && (
            <p className="mt-4 text-gray-600">{project.description}</p>
          )}

          <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-500">优先级:</span>
              <span className="ml-2 text-gray-900">{project.priority}</span>
            </div>
            <div>
              <span className="text-gray-500">负责人:</span>
              <span className="ml-2 text-gray-900">{project.ownerNickname || '-'}</span>
            </div>
            <div>
              <span className="text-gray-500">创建时间:</span>
              <span className="ml-2 text-gray-900">{new Date(project.createdAt).toLocaleDateString('zh-CN')}</span>
            </div>
            <div>
              <span className="text-gray-500">更新时间:</span>
              <span className="ml-2 text-gray-900">{new Date(project.updatedAt).toLocaleDateString('zh-CN')}</span>
            </div>
          </div>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <div className="border-b border-gray-200">
            <nav className="flex">
              <button
                onClick={() => setActiveTab('members')}
                className={`px-6 py-3 text-sm font-medium ${activeTab === 'members' ? 'border-b-2 border-indigo-500 text-indigo-600' : 'text-gray-500 hover:text-gray-700'}`}
              >
                成员 ({members.length})
              </button>
              <button
                onClick={() => setActiveTab('tasks')}
                className={`px-6 py-3 text-sm font-medium ${activeTab === 'tasks' ? 'border-b-2 border-indigo-500 text-indigo-600' : 'text-gray-500 hover:text-gray-700'}`}
              >
                任务 ({tasks.length})
              </button>
            </nav>
          </div>

          <div className="p-6">
            {activeTab === 'members' ? (
              <div>
                <div className="mb-4 flex justify-end">
                  <Link href={`/projects/${project.projectNo}/members/add/`} className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700">
                    添加成员
                  </Link>
                </div>
                {members.length === 0 ? (
                  <p className="text-gray-500 text-center py-4">暂无成员</p>
                ) : (
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">用户名</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">昵称</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">邮箱</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">角色</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">加入时间</th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {members.map((member) => (
                        <tr key={member.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{member.username}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{member.nickname || '-'}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{member.email}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{member.role}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(member.joinedAt).toLocaleDateString('zh-CN')}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                            {member.role !== 'OWNER' && (
                              <button
                                onClick={() => handleRemoveMember(member.userId)}
                                className="text-red-600 hover:text-red-900"
                              >
                                移除
                              </button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            ) : (
              <div>
                <div className="mb-4 flex justify-end">
                  <Link href={`/tasks/create/?projectId=${project.id}`} className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700">
                    创建任务
                  </Link>
                </div>
                {tasks.length === 0 ? (
                  <p className="text-gray-500 text-center py-4">暂无任务</p>
                ) : (
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">标题</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">优先级</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">负责人</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {tasks.map((task) => (
                        <tr key={task.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{task.title}</td>
                          <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(task.status)}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{task.priority}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {task.assigneeName || task.assigneeNickname || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(task.createdAt).toLocaleDateString('zh-CN')}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}