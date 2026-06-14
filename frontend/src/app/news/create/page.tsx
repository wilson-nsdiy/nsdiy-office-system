'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { newsApi, newsGroupsApi } from '@/api'
import { NewsGroup } from '@/types'

export default function CreateNewsPage() {
  const params = useParams()
  const router = useRouter()
  const [groups, setGroups] = useState<NewsGroup[]>([])
  const [form, setForm] = useState({
    groupId: '',
    title: '',
    content: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const isEdit = !!params.id

  useEffect(() => {
    loadGroups()
    if (isEdit) {
      loadNews()
    }
  }, [params.id])

  const loadGroups = async () => {
    try {
      const data = await newsGroupsApi.list()
      setGroups(data || [])
    } catch (error) {
      console.error('Failed to load groups:', error)
    }
  }

  const loadNews = async () => {
    try {
      const data = await newsApi.getById(Number(params.id))
      setForm({
        groupId: data.groupId.toString(),
        title: data.title,
        content: data.content || '',
      })
    } catch (error) {
      console.error('Failed to load news:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      const data = {
        groupId: Number(form.groupId),
        title: form.title,
        content: form.content,
      }

      if (isEdit) {
        await newsApi.update(Number(params.id), data)
      } else {
        await newsApi.create(data)
      }
      router.push('/news')
    } catch (err: any) {
      setError(err.response?.data?.message || '保存失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center h-16">
            <Link href="/dashboard" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
            <span className="ml-4 text-gray-500">/ {isEdit ? '编辑新闻' : '新建新闻'}</span>
          </div>
        </div>
      </nav>

      <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <form onSubmit={handleSubmit} className="bg-white shadow rounded-lg p-6">
          {error && (
            <div className="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded">
              {error}
            </div>
          )}

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">分类 *</label>
            <select
              required
              value={form.groupId}
              onChange={(e) => setForm({ ...form, groupId: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">请选择分类</option>
              {groups.map((g) => (
                <option key={g.id} value={g.id}>{g.name}</option>
              ))}
            </select>
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">标题 *</label>
            <input
              type="text"
              required
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">内容</label>
            <textarea
              rows={10}
              value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="flex justify-end gap-4">
            <Link href="/news" className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
              取消
            </Link>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? '保存中...' : '保存'}
            </button>
          </div>
        </form>
      </main>
    </div>
  )
}
