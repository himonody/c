import request from '@/utils/request'

// 登录接口
export interface AdminLoginRequest {
  username: string
  password: string
}

export interface AdminUserInfo {
  id: number
  username: string
  nickName: string
  role: number
  status: number
}

export interface AdminLoginResponse {
  token: string
  adminInfo: AdminUserInfo
}

export interface AdminLogoutResponse {
  message: string
}

export const login = (data: AdminLoginRequest) => {
  return request.post<AdminLoginResponse>('/admin/auth/login', data)
}

// 登出
export const logout = () => {
  return request.post<AdminLogoutResponse>('/admin/auth/logout')
}

// 获取管理员信息
export const getAdminInfo = () => {
  return request.post<{ admin_id: number }>('/admin/auth/me')
}
