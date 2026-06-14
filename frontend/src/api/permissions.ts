import apiClient from './client'
import { Permission, PermissionCreateInput, PermissionUpdateInput } from '@/types'

export const permissionsApi = {
  list: async (resourceType?: string, keyword?: string): Promise<Permission[]> => {
    const params = new URLSearchParams()
    if (resourceType) params.append('resourceType', resourceType)
    if (keyword) params.append('keyword', keyword)
    const response = await apiClient.get(`/api/permissions?${params.toString()}`)
    return response.data
  },

  listAll: async (): Promise<Permission[]> => {
    const response = await apiClient.get('/api/permissions/all')
    return response.data
  },

  getById: async (id: number): Promise<Permission> => {
    const response = await apiClient.get(`/api/permissions/${id}`)
    return response.data
  },

  create: async (data: PermissionCreateInput): Promise<Permission> => {
    const response = await apiClient.post('/api/permissions', data)
    return response.data
  },

  update: async (id: number, data: PermissionUpdateInput): Promise<Permission> => {
    const response = await apiClient.put(`/api/permissions/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/permissions/${id}`)
  },
}
