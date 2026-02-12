import request from '@/utils/request'

// 用户管理相关接口
export interface UserListRequest {
  page: number
  pageSize: number
  username?: string
  status?: number
}

export interface UserInfo {
  id: number
  username: string
  email?: string
  phone?: string
  avatar?: string
  status: number
  role: number
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
}

export interface UserListResponse {
  list: UserInfo[]
  total: number
  page: number
  pageSize: number
}

export interface CreateUserRequest {
  username: string
  password: string
  email?: string
  phone?: string
  role: number
  status: number
}

export interface UpdateUserRequest {
  id: number
  username?: string
  email?: string
  phone?: string
  role?: number
  status?: number
}

// 获取用户列表
export const getUserList = (params: UserListRequest) => {
  return request.post<UserListResponse>('/admin/user/list', params)
}

// 获取用户详情
export const getUser = (id: number) => {
  return request.post<UserInfo>('/admin/user/get', { id })
}

// 创建用户
export const createUser = (data: CreateUserRequest) => {
  return request.post('/admin/user/create', data)
}

// 更新用户
export const updateUser = (data: UpdateUserRequest) => {
  return request.post('/admin/user/update', data)
}

// 删除用户
export const deleteUser = (id: number) => {
  return request.post('/admin/user/delete', { id })
}

// 重置用户密码
export const resetUserPassword = (id: number, newPassword: string) => {
  return request.post('/admin/user/reset-password', { id, newPassword })
}
