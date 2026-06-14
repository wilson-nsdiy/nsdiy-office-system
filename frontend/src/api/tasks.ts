import apiClient from './client'
import { Task, TaskCreateInput, TaskUpdateInput, PaginatedResponse } from '@/types'

export const tasksApi = {
  list: async (params: { projectId?: number; status?: string; priority?: string; assigneeId?: number; page?: number; pageSize?: number }): Promise<PaginatedResponse<Task>> => {
    const searchParams = new URLSearchParams()
    if (params.projectId) searchParams.append('projectId', params.projectId.toString())
    if (params.status) searchParams.append('status', params.status)
    if (params.priority) searchParams.append('priority', params.priority)
    if (params.assigneeId) searchParams.append('assigneeId', params.assigneeId.toString())
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/tasks?${searchParams.toString()}`)
    return response.data
  },

  getById: async (id: number): Promise<Task> => {
    const response = await apiClient.get(`/api/tasks/${id}`)
    return response.data
  },

  create: async (projectId: number, data: TaskCreateInput): Promise<Task> => {
    const response = await apiClient.post(`/api/tasks?projectId=${projectId}`, data)
    return response.data
  },

  update: async (id: number, data: TaskUpdateInput): Promise<Task> => {
    const response = await apiClient.put(`/api/tasks/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/tasks/${id}`)
  },
}
