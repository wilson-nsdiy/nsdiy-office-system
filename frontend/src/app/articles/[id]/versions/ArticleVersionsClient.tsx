'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth'
import { articlesApi } from '@/api'
import { ArticleVersion } from '@/types'

export default function ArticleVersionsClient() {
  const params = useParams()
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()
  const [versions, setVersions] = useState<ArticleVersion[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [loading, setLoading] = useState(true)
  const [selectedVersion, setSelectedVersion] = useState<ArticleVersion | null>(null)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login/')
    }
  }, [isLoading, isAuthenticated, router])

  useEffect(() => {
    if (isAuthenticated && params.id) {
      loadVersions()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, params.id, page])

  const loadVersions = async () => {
    setLoading(true)
    try {
      const result = await articlesApi.getVersions(Number(params.id), page, pageSize)
      setVersions(result.items || [])
      setTotal(result.total)
    } catch (error) {
      console.error('Failed to load versions:', error)
    } finally {
      setLoading(false)
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
            <span className="ml-4 text-gray-500">/ 版本历史</span>
          </div>
        </div>
      </nav>

      <main className="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="mb-4">
          <Link href={`/articles/${params.id}/`} className="text-indigo-600 hover:text-indigo-900">
            &larr; 返回文章详情
          </Link>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">版本</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">编辑者</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">修改原因</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">时间</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {versions.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-4 text-center text-gray-500">暂无版本记录</td>
                </tr>
              ) : (
                versions.map((version) => (
                  <tr key={version.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">v{version.versionNo}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{version.title}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{version.status}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {version.editorName || version.editorNickname || '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {version.editReason || '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(version.createdAt).toLocaleString('zh-CN')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                      <button
                        onClick={() => setSelectedVersion(version)}
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        查看
                      </button>
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
              disabled={versions.length < pageSize}
              className="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
            >
              下一页
            </button>
          </div>
        </div>
      </main>

      {selectedVersion && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex justify-between items-start mb-4">
                <h2 className="text-xl font-bold">版本 v{selectedVersion.versionNo}</h2>
                <button onClick={() => setSelectedVersion(null)} className="text-gray-500 hover:text-gray-700">&times;</button>
              </div>
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium text-gray-500">标题</label>
                  <p className="text-gray-900">{selectedVersion.title}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">状态</label>
                  <p className="text-gray-900">{selectedVersion.status}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">内容</label>
                  <div className="mt-1 p-3 bg-gray-50 rounded text-sm whitespace-pre-wrap">{selectedVersion.content || '-'}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">摘要</label>
                  <p className="text-gray-900">{selectedVersion.summary || '-'}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">修改原因</label>
                  <p className="text-gray-900">{selectedVersion.editReason || '-'}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">编辑时间</label>
                  <p className="text-gray-900">{new Date(selectedVersion.createdAt).toLocaleString('zh-CN')}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}