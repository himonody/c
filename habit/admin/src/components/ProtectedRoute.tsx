import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { RootState } from '@/redux'
import { Permission, usePermission } from '@/utils/permission'

interface ProtectedRouteProps {
  children: React.ReactNode
  requiredPermission?: Permission
  fallback?: React.ReactNode
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiredPermission,
  fallback = <Navigate to="/404" replace />
}) => {
  const location = useLocation()
  const { auth } = useSelector((state: RootState) => state)
  
  // 检查是否已登录
  if (!auth.isLoggedIn) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }
  
  // 如果需要特定权限，检查权限
  if (requiredPermission) {
    const hasPermission = usePermission(requiredPermission)
    if (!hasPermission) {
      return fallback
    }
  }
  
  return <>{children}</>
}

export default ProtectedRoute
