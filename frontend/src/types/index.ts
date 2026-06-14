export interface ApiResponse<T> {
  data?: T
  error?: string
  code?: string
  message?: string
  reason?: string
  metadata?: Record<string, string>
}

export interface ApiError {
  status: number
  code?: string
  reason?: string
  message: string
  metadata?: Record<string, string>
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface LoginRequest {
  username?: string
  email?: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
  tokenType: string
  expiresIn: number
  refreshExpiresIn: number
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface RefreshTokenResponse extends LoginResponse {}

export interface ResetPasswordRequest {
  email: string
}

export interface ResetPasswordConfirm {
  email: string
  verificationCode: string
  newPassword: string
  confirmPassword: string
}

export interface PasswordChangeRequest {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

export interface SafeUser {
  id: number
  username: string
  email: string
  nickname?: string
  roleId?: number
  userType: string
  isActive: boolean
  tokenVersion: number
  createdAt: string
  updatedAt: string
}

export interface UserResponse extends SafeUser {
  role?: {
    id: number
    name: string
    code: string
  }
}

export interface Role {
  id: number
  name: string
  code: string
  description?: string
  isActive: boolean
  isDefault: boolean
  isBuiltin: boolean
  roleType?: string
  createdAt: string
  updatedAt: string
}

export interface RoleCreateInput {
  name: string
  code: string
  description?: string
  roleType?: string
  isActive?: boolean
}

export interface RoleUpdateInput {
  name?: string
  description?: string
  roleType?: string
  isActive?: boolean
}

export interface Permission {
  id: number
  pid?: number
  name: string
  resourceType: string
  resourcePath: string
  httpMethod?: string
  description?: string
  isActive: boolean
  isBuiltin: boolean
  createdAt: string
  updatedAt: string
}

export interface PermissionCreateInput {
  pid?: number
  name: string
  resourceType: string
  resourcePath: string
  httpMethod?: string
  description?: string
  isActive?: boolean
}

export interface PermissionUpdateInput extends Partial<PermissionCreateInput> {}

export interface RolePermissionUpdate {
  permissionIds: number[]
}

export interface RolePermissionResponse {
  roleId: number
  roleName: string
  permissions: Permission[]
}

export interface NewsGroup {
  id: number
  name: string
  description?: string
  sortOrder?: number
  createdAt: string
  updatedAt: string
}

export interface NewsGroupCreateInput {
  name: string
  description?: string
  sortOrder?: number
}

export interface NewsGroupUpdateInput {
  name?: string
  description?: string
  sortOrder?: number
}

export interface News {
  id: number
  groupId: number
  title: string
  content?: string
  creatorId: number
  createdAt: string
  updatedAt: string
}

export interface NewsCreateInput {
  groupId: number
  title: string
  content?: string
}

export interface NewsUpdateInput {
  groupId?: number
  title?: string
  content?: string
}

export interface NewsListItem extends News {
  groupName?: string
  creatorNickname?: string
}

export interface NewsDetail extends News {
  groupName?: string
  creatorNickname?: string
}

export interface NewsListQuery {
  groupId?: number
  keyword?: string
  page?: number
  pageSize?: number
}

export type ArticleStatus = 'DRAFT' | 'PUBLISHED' | 'ARCHIVED'

export interface Article {
  id: number
  articleNo: string
  title: string
  content?: string
  summary?: string
  status: ArticleStatus
  authorId: number
  coverDescription?: string
  coverUrl?: string
  firstPublishedAt?: string
  createdAt: string
  updatedAt: string
}

export interface ArticleCreateInput {
  title: string
  content?: string
  summary?: string
  status?: ArticleStatus
  coverDescription?: string
  coverUrl?: string
}

export interface ArticleUpdateInput {
  title?: string
  content?: string
  summary?: string
  status?: ArticleStatus
  coverDescription?: string
  coverUrl?: string
}

export interface ArticleListItem extends Article {
  authorName?: string
  authorNickname?: string
}

export interface ArticleDetail extends Article {
  authorName?: string
  authorNickname?: string
}

export interface ArticleVersion {
  id: number
  articleId: number
  versionNo: number
  title: string
  content?: string
  coverDescription?: string
  summary?: string
  status: string
  editorId?: number
  editorName?: string
  editorNickname?: string
  editReason?: string
  createdAt: string
}

export interface ArticleListQuery {
  keyword?: string
  status?: ArticleStatus
  page?: number
  pageSize?: number
}

export type ProjectStatus = 'TODO' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED'
export type ProjectPriority = 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT'
export type ProjectRole = 'OWNER' | 'MANAGER' | 'MEMBER'

export interface Project {
  id: number
  name: string
  projectNo: string
  description?: string
  status: ProjectStatus
  priority: ProjectPriority
  expectedStartDate?: string
  expectedEndDate?: string
  startDate?: string
  endDate?: string
  ownerId: number
  ownerNickname?: string
  createdAt: string
  updatedAt: string
}

export interface ProjectCreateInput {
  name: string
  description?: string
  priority?: ProjectPriority
  expectedStartDate?: string
  expectedEndDate?: string
}

export interface ProjectUpdateInput {
  name?: string
  description?: string
  status?: ProjectStatus
  priority?: ProjectPriority
  expectedStartDate?: string
  expectedEndDate?: string
  startDate?: string
  endDate?: string
}

export interface ProjectMember {
  id: number
  projectId: number
  userId: number
  username: string
  nickname?: string
  email: string
  isActive: boolean
  role: ProjectRole
  joinedAt: string
}

export interface ProjectMemberCreateInput {
  userId: number
  role?: ProjectRole
}

export interface ProjectMemberUpdateInput {
  role: ProjectRole
}

export interface ProjectListQuery {
  keyword?: string
  page?: number
  pageSize?: number
}

export type TaskStatus = 'TODO' | 'IN_PROGRESS' | 'REVIEW' | 'DONE'
export type TaskPriority = 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT'

export interface Task {
  id: number
  projectId: number
  parentId?: number
  title: string
  description?: string
  status: TaskStatus
  priority: TaskPriority
  assigneeId?: number
  creatorId: number
  assigneeName?: string
  assigneeNickname?: string
  creatorName?: string
  creatorNickname?: string
  plannedStartDate?: string
  plannedEndDate?: string
  actualStartTime?: string
  actualEndTime?: string
  estimatedHours?: number
  createdAt: string
  updatedAt: string
}

export interface TaskCreateInput {
  title: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  assigneeId?: number
  parentId?: number
  plannedStartDate?: string
  plannedEndDate?: string
  estimatedHours?: number
}

export interface TaskUpdateInput {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  assigneeId?: number
  parentId?: number
  plannedStartDate?: string
  plannedEndDate?: string
  estimatedHours?: number
}

export interface TaskListQuery {
  projectId?: number
  status?: TaskStatus
  priority?: TaskPriority
  assigneeId?: number
  page?: number
  pageSize?: number
}

export interface MediaAccount {
  id: number
  name: string
  platform: string
  accountId: string
  avatar?: string
  status: string
  accessToken?: string
  refreshToken?: string
  tokenExpiresAt?: string
  createdAt: string
  updatedAt: string
}

export interface MediaAccountCreateInput {
  name: string
  platform: string
  accountId: string
  avatar?: string
}

export interface MediaAccountUpdateInput {
  name?: string
  status?: string
  accessToken?: string
  refreshToken?: string
  tokenExpiresAt?: string
}

export interface MediaContent {
  id: number
  title: string
  content?: string
  coverImage?: string
  platform: string
  accountId?: number
  accountName?: string
  status: string
  views?: number
  likes?: number
  comments?: number
  shares?: number
  publishTime?: string
  createdAt: string
  updatedAt: string
}

export interface MediaContentCreateInput {
  title: string
  content?: string
  coverImage?: string
  platform: string
  accountId?: number
  status?: string
  publishTime?: string
}

export interface MediaContentUpdateInput {
  title?: string
  content?: string
  coverImage?: string
  status?: string
  views?: number
  likes?: number
  comments?: number
  shares?: number
  publishTime?: string
}

export interface MediaContentVersion {
  id: number
  contentId: number
  versionNo: number
  title: string
  content?: string
  coverImage?: string
  status: string
  editorId?: number
  editorName?: string
  editorNickname?: string
  editReason?: string
  createdAt: string
}

export interface MediaAccountListQuery {
  keyword?: string
  platform?: string
  status?: string
  page?: number
  pageSize?: number
}

export interface MediaContentListQuery {
  keyword?: string
  platform?: string
  status?: string
  accountId?: number
  page?: number
  pageSize?: number
}

export interface UploadFile {
  id: number
  filename: string
  originalFilename: string
  filePath: string
  fileSize: number
  mimeType: string
  fileType: string
  extension: string
  uploaderId: number
  uploaderName?: string
  uploaderNickname?: string
  purpose?: string
  md5?: string
  referenceCount: number
  createdAt: string
}

export interface UploadFileListQuery {
  keyword?: string
  fileType?: string
  uploaderId?: number
  page?: number
  pageSize?: number
}

export interface ApiToken {
  id: number
  userId: number
  name: string
  tokenHash: string
  tokenPrefix: string
  status: string
  expiresAt?: string
  lastUsedAt?: string
  usageCount: number
  createdAt: string
  updatedAt: string
}

export interface ApiTokenCreateInput {
  name: string
  expiresAt?: string
}

export interface ApiTokenListQuery {
  keyword?: string
  status?: string
  page?: number
  pageSize?: number
}
