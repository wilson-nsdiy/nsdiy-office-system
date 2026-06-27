'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import { permissionsApi } from '@/api'
import { Permission } from '@/types'

export default function PermissionsPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [editingPerm, setEditingPerm] = useState<Permission | null>(null)
  const [form, setForm] = useState({
    name: '',
    resourceType: 'api',
    resourcePath: '',
    httpMethod: 'GET',
    description: '',
    isActive: true,
  })
  const [filter, setFilter] = useState({ resourceType: '', keyword: '' })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated) {
      loadPermissions()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, filter])

  const loadPermissions = async () => {
    setLoading(true)
    try {
      const data = await permissionsApi.list(filter.resourceType, filter.keyword)
      setPermissions(data || [])
    } catch (error) {
      console.error('Failed to load permissions:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingPerm(null)
    setForm({
      name: '',
      resourceType: 'api',
      resourcePath: '',
      httpMethod: 'GET',
      description: '',
      isActive: true,
    })
    setShowModal(true)
  }

  const handleEdit = (perm: Permission) => {
    setEditingPerm(perm)
    setForm({
      name: perm.name,
      resourceType: perm.resourceType,
      resourcePath: perm.resourcePath,
      httpMethod: perm.httpMethod || 'GET',
      description: perm.description || '',
      isActive: perm.isActive,
    })
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个权限吗？')) return
    try {
      await permissionsApi.delete(id)
      loadPermissions()
    } catch (error: any) {
      alert(error.response?.data?.message || '删除失败')
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)

    try {
      if (editingPerm) {
        await permissionsApi.update(editingPerm.id, form)
      } else {
        await permissionsApi.create(form)
      }
      setShowModal(false)
      loadPermissions()
    } catch (error: any) {
      alert(error.response?.data?.message || '保存失败')
    } finally {
      setSaving(false)
    }
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
              <span className="ml-4 text-gray-500">/ 权限管理</span>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleCreate}
                className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700"
              >
                新建权限
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="mb-6 flex gap-4">
          <input
            type="text"
            placeholder="搜索权限..."
            value={filter.keyword}
            onChange={(e) => setFilter({ ...filter, keyword: e.target.value })}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
          <select
            value={filter.resourceType}
            onChange={(e) => setFilter({ ...filter, resourceType: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">全部类型</option>
            <option value="api">API</option>
            <option value="menu">菜单</option>
            <option value="button">按钮</option>
            <option value="cate">分类</option>
          </select>
          <button
            onClick={loadPermissions}
            className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700"
          >
            搜索
          </button>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源类型</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源路径</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">HTTP方法</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-gray-500">加载中...</td>
                </tr>
              ) : permissions.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-gray-500">暂无权限</td>
                </tr>
              ) : (
                permissions.map((perm) => (
                  <tr key={perm.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{perm.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{perm.resourceType}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{perm.resourcePath}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{perm.httpMethod || '-'}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${perm.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                        {perm.isActive ? '启用' : '禁用'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button onClick={() => handleEdit(perm)} className="text-indigo-600 hover:text-indigo-900 mr-3">编辑</button>
                      {!perm.isBuiltin && (
                        <button onClick={() => handleDelete(perm.id)} className="text-red-600 hover:text-red-900">删除</button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </main>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg w-full max-w-md mx-4">
            <div className="p-6">
              <h2 className="text-lg font-bold mb-4">{editingPerm ? '编辑权限' : '新建权限'}</h2>
              <form onSubmit={handleSubmit}>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">名称 *</label>
                  <input
                    type="text"
                    required
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">资源类型 *</label>
                  <select
                    required
                    value={form.resourceType}
                    onChange={(e) => setForm({ ...form, resourceType: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option value="api">API</option>
                    <option value="menu">菜单</option>
                    <option value="button">按钮</option>
                    <option value="cate">分类</option>
                  </select>
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">资源路径 *</label>
                  <input
                    type="text"
                    required
                    value={form.resourcePath}
                    onChange={(e) => setForm({ ...form, resourcePath: e.target.value })}
                    placeholder="/api/xxx"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">HTTP方法</label>
                  <select
                    value={form.httpMethod}
                    onChange={(e) => setForm({ ...form, httpMethod: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option value="GET">GET</option>
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="DELETE">DELETE</option>
                  </select>
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    rows={2}
                    value={form.description}
                    onChange={(e) => setForm({ ...form, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="mb-6">
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={form.isActive}
                      onChange={(e) => setForm({ ...form, isActive: e.target.checked })}
                      className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <span className="text-sm text-gray-700">启用</span>
                  </label>
                </div>
                <div className="flex justify-end gap-4">
                  <button
                    type="button"
                    onClick={() => setShowModal(false)}
                    className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    disabled={saving}
                    className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50"
                  >
                    {saving ? '保存中...' : '保存'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
