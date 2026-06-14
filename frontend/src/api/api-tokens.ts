import apiClient from './client'
import { ApiToken, ApiTokenCreateInput, PaginatedResponse } from '@/types'

export const apiTokensApi = {
  list: async (params: { keyword?: string; status?: string; page?: number; pageSize?: number }): Promise<PaginatedResponse<ApiToken>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.status) searchParams.append('status', params.status)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/api-tokens?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<ApiToken> => {
    const response = await apiClient.get(`/api/api-tokens/${id}`)
    return response.data
  },

  create: async (data: ApiTokenCreateInput): Promise<{ token: ApiToken; rawToken: string }> => {
    const response = await apiClient.post('/api/api-tokens', data)
    return response.data
  },

  update: async (id: number, data: Partial<ApiToken>): Promise<ApiToken> => {
    const response = await apiClient.put(`/api/api-tokens/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/api-tokens/${id}`)
  },
}
