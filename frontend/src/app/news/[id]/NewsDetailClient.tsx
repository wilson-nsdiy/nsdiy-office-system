'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth'
import { newsApi, newsGroupsApi } from '@/api'
import { NewsDetail, NewsGroup } from '@/types'

export default function NewsDetailClient() {
  const params = useParams()
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [news, setNews] = useState<NewsDetail | null>(null)
  const [groups, setGroups] = useState<NewsGroup[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login/')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated && params.id) {
      loadData()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, params.id])

  const loadData = async () => {
    setLoading(true)
    try {
      const [newsData, groupsData] = await Promise.all([
        newsApi.getById(Number(params.id)),
        newsGroupsApi.list()
      ])
      setNews(newsData)
      setGroups(groupsData || [])
    } catch (error) {
      console.error('Failed to load news:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('确定要删除这条新闻吗？')) return
    try {
      await newsApi.delete(Number(params.id))
      router.push('/news/')
    } catch (error) {
      alert('删除失败')
    }
  }

  if (isLoading || loading) {
    return <div className="min-h-screen flex items-center justify-center">加载中...</div>
  }

  if (!news) {
    return <div className="min-h-screen flex items-center justify-center">新闻不存在</div>
  }

  const group = groups.find(g => g.id === news.groupId)

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center h-16">
            <Link href="/dashboard/" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
            <span className="ml-4 text-gray-500">/ 新闻详情</span>
          </div>
        </div>
      </nav>

      <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg p-6">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{news.title}</h1>
              <p className="mt-1 text-sm text-gray-500">
                分类: {group?.name || '-'}
              </p>
            </div>
            <div className="flex gap-2">
              <Link
                href={`/news/${news.id}/edit/`}
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

          <div className="prose max-w-none">
            {news.content}
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200 text-sm text-gray-500">
            <p>创建者: {news.creatorNickname || news.creatorId}</p>
            <p>创建时间: {new Date(news.createdAt).toLocaleString('zh-CN')}</p>
            <p>更新时间: {new Date(news.updatedAt).toLocaleString('zh-CN')}</p>
          </div>
        </div>
      </main>
    </div>
  )
}