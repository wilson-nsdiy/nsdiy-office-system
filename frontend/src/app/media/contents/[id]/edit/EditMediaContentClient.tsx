'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth'
import { mediaContentsApi, mediaAccountsApi } from '@/api'
import { MediaAccount, MediaContent } from '@/types'

export default function EditMediaContentClient() {
  const params = useParams()
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [accounts, setAccounts] = useState<MediaAccount[]>([])
  const [form, setForm] = useState({
    title: '',
    content: '',
    coverImage: '',
    platform: 'wechat',
    accountId: '',
    status: 'draft',
    views: 0,
    likes: 0,
    comments: 0,
    shares: 0,
    publishTime: '',
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login/')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated && params.id) {
      loadData()
    }
  }, [isAuthenticated, params.id])

  const loadData = async () => {
    setLoading(true)
    try {
      const [contentData, accountsResult] = await Promise.all([
        mediaContentsApi.getById(Number(params.id)),
        mediaAccountsApi.list({ pageSize: 100 })
      ])
      setAccounts(accountsResult.items || [])
      setForm({
        title: contentData.title,
        content: contentData.content || '',
        coverImage: contentData.coverImage || '',
        platform: contentData.platform,
        accountId: contentData.accountId?.toString() || '',
        status: contentData.status,
        views: contentData.views || 0,
        likes: contentData.likes || 0,
        comments: contentData.comments || 0,
        shares: contentData.shares || 0,
        publishTime: contentData.publishTime ? contentData.publishTime.slice(0, 16) : '',
      })
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError('')

    try {
      const data: any = {
        title: form.title,
        content: form.content,
        coverImage: form.coverImage,
        platform: form.platform,
        status: form.status,
        views: form.views,
        likes: form.likes,
        comments: form.comments,
        shares: form.shares,
      }
      if (form.accountId) data.accountId = Number(form.accountId)
      if (form.publishTime) data.publishTime = form.publishTime

      await mediaContentsApi.update(Number(params.id), data)
      router.push('/media/contents/')
    } catch (err: any) {
      setError(err.response?.data?.message || '更新失败')
    } finally {
      setSaving(false)
    }
  }

  if (isLoading || loading) {
    return <div className="min-h-screen flex items-center justify-center">加载中...</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center h-16">
            <Link href="/dashboard/" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
            <span className="ml-4 text-gray-500">/ 编辑媒体内容</span>
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
            <label className="block text-sm font-medium text-gray-700 mb-1">标题 *</label>
            <input
              type="text"
              required
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">平台 *</label>
              <select
                required
                value={form.platform}
                onChange={(e) => setForm({ ...form, platform: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="wechat">微信</option>
                <option value="weibo">微博</option>
                <option value="douyin">抖音</option>
                <option value="kuaishou">快手</option>
                <option value="bilibili">B站</option>
                <option value="xiaohongshu">小红书</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">账号</label>
              <select
                value={form.accountId}
                onChange={(e) => setForm({ ...form, accountId: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="">选择账号</option>
                {accounts.map((a) => (
                  <option key={a.id} value={a.id}>{a.name}</option>
                ))}
              </select>
            </div>
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">封面图片URL</label>
            <input
              type="text"
              value={form.coverImage}
              onChange={(e) => setForm({ ...form, coverImage: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">内容</label>
            <textarea
              rows={10}
              value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">状态</label>
              <select
                value={form.status}
                onChange={(e) => setForm({ ...form, status: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="draft">草稿</option>
                <option value="published">已发布</option>
                <option value="archived">已归档</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">发布时间</label>
              <input
                type="datetime-local"
                value={form.publishTime}
                onChange={(e) => setForm({ ...form, publishTime: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>

          <div className="grid grid-cols-4 gap-4 mb-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">浏览量</label>
              <input
                type="number"
                value={form.views}
                onChange={(e) => setForm({ ...form, views: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">点赞数</label>
              <input
                type="number"
                value={form.likes}
                onChange={(e) => setForm({ ...form, likes: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">评论数</label>
              <input
                type="number"
                value={form.comments}
                onChange={(e) => setForm({ ...form, comments: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">分享数</label>
              <input
                type="number"
                value={form.shares}
                onChange={(e) => setForm({ ...form, shares: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>

          <div className="flex justify-end gap-4">
            <Link href="/media/contents/" className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
              取消
            </Link>
            <button
              type="submit"
              disabled={saving}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50"
            >
              {saving ? '保存中...' : '保存'}
            </button>
          </div>
        </form>
      </main>
    </div>
  )
}