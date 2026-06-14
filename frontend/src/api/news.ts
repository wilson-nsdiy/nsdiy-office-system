import apiClient from './client'
import { NewsGroup, NewsGroupCreateInput, NewsGroupUpdateInput, News, NewsCreateInput, NewsUpdateInput, NewsDetail, PaginatedResponse } from '@/types'

export const newsGroupsApi = {
  list: async (): Promise<NewsGroup[]> => {
    const response = await apiClient.get('/api/news-groups')
    return response.data
  },

  getById: async (id: number): Promise<NewsGroup> => {
    const response = await apiClient.get(`/api/news-groups/${id}`)
    return response.data
  },

  create: async (data: NewsGroupCreateInput): Promise<NewsGroup> => {
    const response = await apiClient.post('/api/news-groups', data)
    return response.data
  },

  update: async (id: number, data: NewsGroupUpdateInput): Promise<NewsGroup> => {
    const response = await apiClient.put(`/api/news-groups/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/news-groups/${id}`)
  },
}

export const newsApi = {
  list: async (params: { groupId?: number; keyword?: string; page?: number; pageSize?: number }): Promise<PaginatedResponse<News>> => {
    const searchParams = new URLSearchParams()
    if (params.groupId) searchParams.append('groupId', params.groupId.toString())
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/news?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<NewsDetail> => {
    const response = await apiClient.get(`/api/news/${id}`)
    return response.data
  },

  create: async (data: NewsCreateInput): Promise<News> => {
    const response = await apiClient.post('/api/news', data)
    return response.data
  },

  update: async (id: number, data: NewsUpdateInput): Promise<News> => {
    const response = await apiClient.put(`/api/news/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/news/${id}`)
  },
}
