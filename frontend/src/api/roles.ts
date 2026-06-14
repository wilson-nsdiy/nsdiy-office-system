import apiClient from './client'
import { Role, RoleCreateInput, RoleUpdateInput, Permission, PaginatedResponse } from '@/types'

export const rolesApi = {
  list: async (): Promise<Role[]> => {
    const response = await apiClient.get('/api/roles')
    return response.data
  },

  getById: async (id: number): Promise<Role> => {
    const response = await apiClient.get(`/api/roles/${id}`)
    return response.data
  },

  create: async (data: RoleCreateInput): Promise<Role> => {
    const response = await apiClient.post('/api/roles', data)
    return response.data
  },

  update: async (id: number, data: RoleUpdateInput): Promise<Role> => {
    const response = await apiClient.put(`/api/roles/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/roles/${id}`)
  },

  getPermissions: async (id: number): Promise<Permission[]> => {
    const response = await apiClient.get(`/api/roles/${id}/permissions`)
    return response.data.permissions
  },

  updatePermissions: async (id: number, permissionIds: number[]): Promise<void> => {
    await apiClient.put(`/api/roles/${id}/permissions`, { permissionIds })
  },
}
