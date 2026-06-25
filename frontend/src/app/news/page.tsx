'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/contexts/auth'
import { newsApi, newsGroupsApi } from '@/api'
import { News, NewsGroup } from '@/types'

export default function NewsPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const router = useRouter()
  const [newsList, setNewsList] = useState<News[]>([])
  const [groups, setGroups] = useState<NewsGroup[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [groupId, setGroupId] = useState<string>('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated) {
      loadGroups()
      loadNews()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, page, groupId])

  const loadGroups = async () => {
    try {
      const data = await newsGroupsApi.list()
      setGroups(data || [])
    } catch (error) {
      console.error('Failed to load groups:', error)
    }
  }

  const loadNews = async () => {
    setLoading(true)
    try {
      const params: any = { page, pageSize }
      if (keyword) params.keyword = keyword
      if (groupId) params.groupId = Number(groupId)
      const result = await newsApi.list(params)
      setNewsList(result.items || [])
      setTotal(result.total)
    } catch (error) {
      console.error('Failed to load news:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setPage(1)
    loadNews()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条新闻吗？')) return
    try {
      await newsApi.delete(id)
      loadNews()
    } catch (error) {
      alert('删除失败')
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
              <span className="ml-4 text-gray-500">/ 新闻管理</span>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/news-groups" className="text-gray-600 hover:text-gray-900 text-sm">分类管理</Link>
              <Link href="/news/create" className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-indigo-700">
                新建新闻
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="mb-6 flex gap-4">
          <input
            type="text"
            placeholder="搜索新闻..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
          <select
            value={groupId}
            onChange={(e) => setGroupId(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">全部分类</option>
            {groups.map((g) => (
              <option key={g.id} value={g.id}>{g.name}</option>
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
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">分类</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建者</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-6 py-4 text-center text-gray-500">加载中...</td>
                </tr>
              ) : newsList.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
              ) : (
                newsList.map((news) => (
                  <tr key={news.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <Link href={`/news/${news.id}`} className="text-indigo-600 hover:text-indigo-900">
                        {news.title}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {groups.find(g => g.id === news.groupId)?.name || '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{news.creatorId}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(news.createdAt).toLocaleDateString('zh-CN')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <Link href={`/news/${news.id}/edit`} className="text-indigo-600 hover:text-indigo-900 mr-4">编辑</Link>
                      <button onClick={() => handleDelete(news.id)} className="text-red-600 hover:text-red-900">删除</button>
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
              disabled={newsList.length < pageSize}
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
