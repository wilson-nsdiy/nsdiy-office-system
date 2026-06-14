import apiClient from './client'
import { Project, ProjectCreateInput, ProjectUpdateInput, ProjectMember, ProjectMemberCreateInput, ProjectMemberUpdateInput, PaginatedResponse } from '@/types'

export const projectsApi = {
  list: async (params: { keyword?: string; page?: number; pageSize?: number }): Promise<PaginatedResponse<Project>> => {
    const searchParams = new URLSearchParams()
    if (params.keyword) searchParams.append('keyword', params.keyword)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.pageSize) searchParams.append('pageSize', params.pageSize.toString())
    const response = await apiClient.get(`/api/projects?${searchParams.toString()}`)
    return response.data
  },

  getByProjectNo: async (projectNo: string): Promise<Project> => {
    const response = await apiClient.get(`/api/projects/${projectNo}`)
    return response.data
  },

  create: async (data: ProjectCreateInput): Promise<Project> => {
    const response = await apiClient.post('/api/projects', data)
    return response.data
  },

  update: async (projectNo: string, data: ProjectUpdateInput): Promise<Project> => {
    const response = await apiClient.put(`/api/projects/${projectNo}`, data)
    return response.data
  },

  delete: async (projectNo: string): Promise<void> => {
    await apiClient.delete(`/api/projects/${projectNo}`)
  },

  getMembers: async (projectNo: string): Promise<ProjectMember[]> => {
    const response = await apiClient.get(`/api/projects/${projectNo}/members`)
    return response.data
  },

  addMember: async (projectNo: string, data: ProjectMemberCreateInput): Promise<void> => {
    await apiClient.post(`/api/projects/${projectNo}/members`, data)
  },

  updateMemberRole: async (projectNo: string, userId: number, data: ProjectMemberUpdateInput): Promise<void> => {
    await apiClient.put(`/api/projects/${projectNo}/members/${userId}`, data)
  },

  removeMember: async (projectNo: string, userId: number): Promise<void> => {
    await apiClient.delete(`/api/projects/${projectNo}/members/${userId}`)
  },
}
