import { useSelector } from 'react-redux'
import { RootState } from '@/redux'
import { AdminUserInfo } from '@/api/auth'

// 权限枚举
export enum Permission {
  // 挑战管理
  CHALLENGE_VIEW = 'challenge:view',
  CHALLENGE_CREATE = 'challenge:create',
  CHALLENGE_UPDATE = 'challenge:update',
  CHALLENGE_DELETE = 'challenge:delete',
  
  // 用户管理
  USER_VIEW = 'user:view',
  USER_CREATE = 'user:create',
  USER_UPDATE = 'user:update',
  USER_DELETE = 'user:delete',
  
  // 系统配置
  CONFIG_VIEW = 'config:view',
  CONFIG_CREATE = 'config:create',
  CONFIG_UPDATE = 'config:update',
  CONFIG_DELETE = 'config:delete',
  
  // 系统设置
  SETTINGS_VIEW = 'settings:view',
  SETTINGS_UPDATE = 'settings:update',
  
  // 仪表盘
  DASHBOARD_VIEW = 'dashboard:view',
}

// 角色权限映射
const ROLE_PERMISSIONS: Record<number, Permission[]> = {
  1: [ // 超级管理员 - 所有权限
    Permission.CHALLENGE_VIEW,
    Permission.CHALLENGE_CREATE,
    Permission.CHALLENGE_UPDATE,
    Permission.CHALLENGE_DELETE,
    Permission.USER_VIEW,
    Permission.USER_CREATE,
    Permission.USER_UPDATE,
    Permission.USER_DELETE,
    Permission.CONFIG_VIEW,
    Permission.CONFIG_CREATE,
    Permission.CONFIG_UPDATE,
    Permission.CONFIG_DELETE,
    Permission.SETTINGS_VIEW,
    Permission.SETTINGS_UPDATE,
    Permission.DASHBOARD_VIEW,
  ],
  2: [ // 管理员 - 部分权限
    Permission.CHALLENGE_VIEW,
    Permission.CHALLENGE_CREATE,
    Permission.CHALLENGE_UPDATE,
    Permission.USER_VIEW,
    Permission.USER_UPDATE,
    Permission.CONFIG_VIEW,
    Permission.CONFIG_UPDATE,
    Permission.DASHBOARD_VIEW,
  ],
  3: [ // 普通用户 - 基本权限
    Permission.CHALLENGE_VIEW,
    Permission.DASHBOARD_VIEW,
  ],
}

// 检查用户是否有指定权限
export const hasPermission = (
  user: AdminUserInfo | null,
  permission: Permission
): boolean => {
  if (!user) return false
  
  const permissions = ROLE_PERMISSIONS[user.role] || []
  return permissions.includes(permission)
}

// 检查用户是否有任一权限
export const hasAnyPermission = (
  user: AdminUserInfo | null,
  permissions: Permission[]
): boolean => {
  if (!user) return false
  
  return permissions.some(permission => hasPermission(user, permission))
}

// 检查用户是否有所有权限
export const hasAllPermissions = (
  user: AdminUserInfo | null,
  permissions: Permission[]
): boolean => {
  if (!user) return false
  
  return permissions.every(permission => hasPermission(user, permission))
}

// 权限检查 Hook
export const usePermission = (permission: Permission) => {
  const { auth } = useSelector((state: RootState) => state)
  return hasPermission(auth.userInfo, permission)
}

// 多权限检查 Hook
export const usePermissions = (permissions: Permission[], requireAll = false) => {
  const { auth } = useSelector((state: RootState) => state)
  return requireAll 
    ? hasAllPermissions(auth.userInfo, permissions)
    : hasAnyPermission(auth.userInfo, permissions)
}
