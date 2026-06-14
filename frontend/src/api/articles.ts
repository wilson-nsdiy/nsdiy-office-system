import apiClient from './client'
import { Article, ArticleCreateInput, ArticleUpdateInput, ArticleDetail, ArticleListItem, ArticleVersion, PaginatedResponse } from '@/types'

export const articlesApi = {
  list: async (params: { keyword?: string; status?: string; page?: number; pageSize?: number }): Promise<PaginatedResponse<ArticleListItem>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.status) searchParams.append('status', params.status)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/articles?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<ArticleDetail> => {
    const response = await apiClient.get(`/api/articles/${id}`)
    return response.data
  },

  create: async (data: ArticleCreateInput): Promise<Article> => {
    const response = await apiClient.post('/api/articles', data)
    return response.data
  },

  update: async (id: number, data: ArticleUpdateInput, editReason?: string): Promise<Article> => {
    const params = editReason ? `?editReason=${encodeURIComponent(editReason)}` : ''
    const response = await apiClient.put(`/api/articles/${id}${params}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/articles/${id}`)
  },

  getVersions: async (id: number, page?: number, pageSize?: number): Promise<PaginatedResponse<ArticleVersion>> => {
    const searchParams = new URLSearchParams()
    if (page) searchParams.append('page', page.toString())
    if (pageSize) searchParams.append('pageSize', pageSize.toString())
    const response = await apiClient.get(`/api/articles/${id}/versions?${searchParams.toString()}`)
    return response.data
  },

  getVersion: async (id: number, versionId: number): Promise<ArticleVersion> => {
    const response = await apiClient.get(`/api/articles/${id}/versions/${versionId}`)
    return response.data
  },
}
