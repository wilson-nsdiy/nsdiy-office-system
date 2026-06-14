'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import { mediaContentsApi, mediaAccountsApi } from '@/api'
import { MediaContent, MediaAccount } from '@/types'

export default function MediaContentsPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()
  const [contents, setContents] = useState<MediaContent[]>([])
  const [accounts, setAccounts] = useState<MediaAccount[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [platform, setPlatform] = useState('')
  const [status, setStatus] = useState('')
  const [accountId, setAccountId] = useState<string>('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated) {
      loadAccounts()
      loadContents()
    }
  }, [isAuthenticated, page])

  const loadAccounts = async () => {
    try {
      const result = await mediaAccountsApi.list({ pageSize: 100 })
      setAccounts(result.items || [])
    } catch (error) {
      console.error('Failed to load accounts:', error)
    }
  }

  const loadContents = async () => {
    setLoading(true)
    try {
      const params: any = { keyword, platform, status, page, pageSize }
      if (accountId) params.accountId = Number(accountId)
      const result = await mediaContentsApi.list(params)
      setContents(result.items || [])
      setTotal(result.total)
    } catch (error) {
      console.error('Failed to load contents:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setPage(1)
    loadContents()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条内容吗？')) return
    try {
      await mediaContentsApi.delete(id)
      loadContents()
    } catch (error) {
      alert('删除失败')
    }
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-800',
      published: 'bg-green-100 text-green-800',
      archived: 'bg-yellow-100 text-yellow-800',
    }
    const labels: Record<string, string> = {
      draft: '草稿',
      published: '已发布',
      archived: '已归档',
    }
    return (
      <span className={`px-2 py-1 text-xs font-medium rounded-full ${colors[status] || ''}`}>
        {labels[status] || status}
      </span>
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
              <span className="ml-4 text-gray-500">/ 媒体内容</span>
            </div>
            <div className="flex items-center">
              <Link href="/media/contents/create" className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700">
                新建内容
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="mb-6 flex gap-4 flex-wrap">
          <input
            type="text"
            placeholder="搜索内容..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="flex-1 min-w-[200px] px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
          <select
            value={platform}
            onChange={(e) => setPlatform(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">全部平台</option>
            <option value="wechat">微信</option>
            <option value="weibo">微博</option>
            <option value="douyin">抖音</option>
            <option value="kuaishou">快手</option>
            <option value="bilibili">B站</option>
            <option value="xiaohongshu">小红书</option>
          </select>
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">全部状态</option>
            <option value="draft">草稿</option>
            <option value="published">已发布</option>
            <option value="archived">已归档</option>
          </select>
          <select
            value={accountId}
            onChange={(e) => setAccountId(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">全部账号</option>
            {accounts.map((a) => (
              <option key={a.id} value={a.id}>{a.name}</option>
            ))}
          </select>
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
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">平台</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">账号</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">浏览量</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center text-gray-500">加载中...</td>
                </tr>
              ) : contents.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
              ) : (
                contents.map((content) => (
                  <tr key={content.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <Link href={`/media/contents/${content.id}/edit`} className="text-indigo-600 hover:text-indigo-900">
                        {content.title}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{content.platform}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{content.accountName || '-'}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(content.status)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{content.views || 0}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(content.createdAt).toLocaleDateString('zh-CN')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <Link href={`/media/contents/${content.id}/edit`} className="text-indigo-600 hover:text-indigo-900 mr-4">编辑</Link>
                      <button onClick={() => handleDelete(content.id)} className="text-red-600 hover:text-red-900">删除</button>
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
              disabled={contents.length < pageSize}
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
