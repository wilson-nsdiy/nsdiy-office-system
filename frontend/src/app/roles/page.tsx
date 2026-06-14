'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import { rolesApi, permissionsApi } from '@/api'
import { Role, Permission } from '@/types'

export default function RolesPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [showPermModal, setShowPermModal] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null)
  const [selectedPermIds, setSelectedPermIds] = useState<number[]>([])
  const [form, setForm] = useState({ name: '', code: '', description: '', roleType: '' })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated) {
      loadRoles()
      loadPermissions()
    }
  }, [isAuthenticated])

  const loadRoles = async () => {
    setLoading(true)
    try {
      const data = await rolesApi.list()
      setRoles(data || [])
    } catch (error) {
      console.error('Failed to load roles:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadPermissions = async () => {
    try {
      const data = await permissionsApi.listAll()
      setPermissions(data || [])
    } catch (error) {
      console.error('Failed to load permissions:', error)
    }
  }

  const handleCreate = () => {
    setEditingRole(null)
    setForm({ name: '', code: '', description: '', roleType: '' })
    setShowModal(true)
  }

  const handleEdit = (role: Role) => {
    setEditingRole(role)
    setForm({
      name: role.name,
      code: role.code,
      description: role.description || '',
      roleType: role.roleType || '',
    })
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个角色吗？')) return
    try {
      await rolesApi.delete(id)
      loadRoles()
    } catch (error: any) {
      alert(error.response?.data?.message || '删除失败')
    }
  }

  const handleManagePermissions = async (role: Role) => {
    setSelectedRoleId(role.id)
    try {
      const rolePerms = await rolesApi.getPermissions(role.id)
      setSelectedPermIds(rolePerms.map(p => p.id))
      setShowPermModal(true)
    } catch (error) {
      console.error('Failed to load role permissions:', error)
    }
  }

  const handleSavePermissions = async () => {
    if (!selectedRoleId) return
    setSaving(true)
    try {
      await rolesApi.updatePermissions(selectedRoleId, selectedPermIds)
      setShowPermModal(false)
    } catch (error: any) {
      alert(error.response?.data?.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)

    try {
      if (editingRole) {
        await rolesApi.update(editingRole.id, form)
      } else {
        await rolesApi.create(form)
      }
      setShowModal(false)
      loadRoles()
    } catch (error: any) {
      alert(error.response?.data?.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const togglePerm = (permId: number) => {
    setSelectedPermIds(prev =>
      prev.includes(permId) ? prev.filter(id => id !== permId) : [...prev, permId]
    )
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
              <span className="ml-4 text-gray-500">/ 角色管理</span>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/permissions" className="text-gray-600 hover:text-gray-900 text-sm">权限管理</Link>
              <button
                onClick={handleCreate}
                className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700"
              >
                新建角色
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">编码</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">描述</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-gray-500">加载中...</td>
                </tr>
              ) : roles.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-gray-500">暂无角色</td>
                </tr>
              ) : (
                roles.map((role) => (
                  <tr key={role.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{role.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{role.code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{role.description || '-'}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${role.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                        {role.isActive ? '启用' : '禁用'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {role.isBuiltin ? (
                        <span className="px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800">内置</span>
                      ) : role.isDefault ? (
                        <span className="px-2 py-1 text-xs font-medium rounded-full bg-yellow-100 text-yellow-800">默认</span>
                      ) : (
                        <span className="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800">自定义</span>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button onClick={() => handleManagePermissions(role)} className="text-indigo-600 hover:text-indigo-900 mr-3">权限</button>
                      <button onClick={() => handleEdit(role)} className="text-indigo-600 hover:text-indigo-900 mr-3">编辑</button>
                      {!role.isBuiltin && (
                        <button onClick={() => handleDelete(role.id)} className="text-red-600 hover:text-red-900">删除</button>
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
              <h2 className="text-lg font-bold mb-4">{editingRole ? '编辑角色' : '新建角色'}</h2>
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
                  <label className="block text-sm font-medium text-gray-700 mb-1">编码 *</label>
                  <input
                    type="text"
                    required
                    value={form.code}
                    onChange={(e) => setForm({ ...form, code: e.target.value })}
                    disabled={!!editingRole}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100"
                  />
                </div>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    rows={3}
                    value={form.description}
                    onChange={(e) => setForm({ ...form, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-1">类型</label>
                  <input
                    type="text"
                    value={form.roleType}
                    onChange={(e) => setForm({ ...form, roleType: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
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

      {showPermModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg w-full max-w-2xl mx-4 max-h-[80vh] overflow-hidden flex flex-col">
            <div className="p-6 border-b">
              <h2 className="text-lg font-bold">管理权限</h2>
            </div>
            <div className="p-6 overflow-y-auto flex-1">
              {permissions.length === 0 ? (
                <p className="text-gray-500 text-center">暂无权限</p>
              ) : (
                <div className="space-y-2">
                  {permissions.map((perm) => (
                    <label key={perm.id} className="flex items-center gap-3 p-2 hover:bg-gray-50 rounded">
                      <input
                        type="checkbox"
                        checked={selectedPermIds.includes(perm.id)}
                        onChange={() => togglePerm(perm.id)}
                        className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                      />
                      <span className="text-sm">{perm.name}</span>
                      <span className="text-xs text-gray-400">{perm.resourceType} - {perm.resourcePath}</span>
                    </label>
                  ))}
                </div>
              )}
            </div>
            <div className="p-6 border-t flex justify-end gap-4">
              <button
                onClick={() => setShowPermModal(false)}
                className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button
                onClick={handleSavePermissions}
                disabled={saving}
                className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50"
              >
                {saving ? '保存中...' : '保存'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
