'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth'
import { articlesApi } from '@/api'
import { ArticleDetail } from '@/types'

export default function ArticleDetailClient() {
  const params = useParams()
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [article, setArticle] = useState<ArticleDetail | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login/')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated && params.id) {
      loadArticle()
    }
  }, [isAuthenticated, params.id])

  const loadArticle = async () => {
    setLoading(true)
    try {
      const data = await articlesApi.getById(Number(params.id))
      setArticle(data)
    } catch (error) {
      console.error('Failed to load article:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('确定要删除这篇文章吗？')) return
    try {
      await articlesApi.delete(Number(params.id))
      router.push('/articles/')
    } catch (error) {
      alert('删除失败')
    }
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      DRAFT: 'bg-gray-100 text-gray-800',
      PUBLISHED: 'bg-green-100 text-green-800',
      ARCHIVED: 'bg-yellow-100 text-yellow-800',
    }
    const labels: Record<string, string> = {
      DRAFT: '草稿',
      PUBLISHED: '已发布',
      ARCHIVED: '已归档',
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

  if (!article) {
    return <div className="min-h-screen flex items-center justify-center">文章不存在</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center h-16">
            <Link href="/dashboard/" className="text-xl font-bold text-indigo-600">OA-NSDIY</Link>
            <span className="ml-4 text-gray-500">/ 文章详情</span>
          </div>
        </div>
      </nav>

      <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg p-6">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{article.title}</h1>
              <p className="mt-1 text-sm text-gray-500">
                {article.articleNo} · {getStatusBadge(article.status)}
              </p>
            </div>
            <div className="flex gap-2">
              <Link
                href={`/articles/${article.id}/edit/`}
                className="px-3 py-1 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
              >
                编辑
              </Link>
              <Link
                href={`/articles/${article.id}/versions/`}
                className="px-3 py-1 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
              >
                版本历史
              </Link>
              <button
                onClick={handleDelete}
                className="px-3 py-1 border border-red-300 rounded-md text-sm text-red-600 hover:bg-red-50"
              >
                删除
              </button>
            </div>
          </div>

          {article.coverUrl && (
            <div className="mb-6">
              <img src={article.coverUrl} alt={article.title} className="w-full h-64 object-cover rounded-lg" />
            </div>
          )}

          {article.summary && (
            <div className="mb-6 p-4 bg-gray-50 rounded-lg">
              <p className="text-gray-700">{article.summary}</p>
            </div>
          )}

          <div className="prose max-w-none">
            {article.content}
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200 text-sm text-gray-500">
            <p>作者: {article.authorName || article.authorNickname || '-'}</p>
            <p>创建时间: {new Date(article.createdAt).toLocaleString('zh-CN')}</p>
            <p>更新时间: {new Date(article.updatedAt).toLocaleString('zh-CN')}</p>
          </div>
        </div>
      </main>
    </div>
  )
}