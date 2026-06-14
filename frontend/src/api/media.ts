import apiClient from './client'
import { MediaAccount, MediaAccountCreateInput, MediaAccountUpdateInput, MediaContent, MediaContentCreateInput, MediaContentUpdateInput, MediaContentVersion, PaginatedResponse } from '@/types'

export const mediaAccountsApi = {
  list: async (params: { keyword?: string; platform?: string; status?: string; page?: number; pageSize?: number }): Promise<PaginatedResponse<MediaAccount>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.platform) searchParams.append('platform', params.platform)
    if (params.status) searchParams.append('status', params.status)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/media/accounts?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<MediaAccount> => {
    const response = await apiClient.get(`/api/media/accounts/${id}`)
    return response.data
  },

  create: async (data: MediaAccountCreateInput): Promise<MediaAccount> => {
    const response = await apiClient.post('/api/media/accounts', data)
    return response.data
  },

  update: async (id: number, data: MediaAccountUpdateInput): Promise<MediaAccount> => {
    const response = await apiClient.put(`/api/media/accounts/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/media/accounts/${id}`)
  },
}

export const mediaContentsApi = {
  list: async (params: { keyword?: string; platform?: string; status?: string; accountId?: number; page?: number; pageSize?: number }): Promise<PaginatedResponse<MediaContent>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.platform) searchParams.append('platform', params.platform)
    if (params.status) searchParams.append('status', params.status)
    if (params.accountId) searchParams.append('accountId', params.accountId.toString())
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/media/contents?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<MediaContent> => {
    const response = await apiClient.get(`/api/media/contents/${id}`)
    return response.data
  },

  create: async (data: MediaContentCreateInput): Promise<MediaContent> => {
    const response = await apiClient.post('/api/media/contents', data)
    return response.data
  },

  update: async (id: number, data: MediaContentUpdateInput): Promise<MediaContent> => {
    const response = await apiClient.put(`/api/media/contents/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/media/contents/${id}`)
  },

  getVersions: async (id: number, page?: number, pageSize?: number): Promise<PaginatedResponse<MediaContentVersion>> => {
    const searchParams = new URLSearchParams()
    if (page) searchParams.append('page', page.toString())
    if (pageSize) searchParams.append('pageSize', pageSize.toString())
    const response = await apiClient.get(`/api/media/contents/${id}/versions?${searchParams.toString()}`)
    return response.data
  },

  getVersion: async (versionId: number): Promise<MediaContentVersion> => {
    const response = await apiClient.get(`/api/media/contents/versions/${versionId}`)
    return response.data
  },
}
