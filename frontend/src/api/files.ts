import apiClient from './client'
import { UploadFile, PaginatedResponse } from '@/types'

export const filesApi = {
  list: async (params: { keyword?: string; fileType?: string; uploaderId?: number; page?: number; pageSize?: number }): Promise<PaginatedResponse<UploadFile>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.fileType) searchParams.append('fileType', params.fileType)
    if (params.uploaderId) searchParams.append('uploaderId', params.uploaderId.toString())
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/files?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<UploadFile> => {
    const response = await apiClient.get(`/api/files/${id}`)
    return response.data
  },

  create: async (data: Partial<UploadFile>): Promise<UploadFile> => {
    const response = await apiClient.post('/api/files', data)
    return response.data
  },

  update: async (id: number, data: Partial<UploadFile>): Promise<UploadFile> => {
    const response = await apiClient.put(`/api/files/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/files/${id}`)
  },
}
