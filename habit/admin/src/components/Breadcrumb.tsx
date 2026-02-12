import React from 'react'
import { Breadcrumb } from 'antd'
import { useLocation, Link } from 'react-router-dom'
import type { BreadcrumbItemType } from 'antd/es/breadcrumb/Breadcrumb'

const CustomBreadcrumb: React.FC = () => {
  const location = useLocation()
  
  const getBreadcrumbItems = (): BreadcrumbItemType[] => {
    const pathSegments = location.pathname.split('/').filter(Boolean)
    const items: BreadcrumbItemType[] = [
      {
        title: <Link to="/dashboard">首页</Link>,
      }
    ]

    const pathMap: Record<string, string> = {
      'dashboard': '仪表盘',
      'challenge': '挑战管理',
      'config': '系统配置',
      'user': '用户管理',
      'settings': '系统设置',
      'login': '登录',
      'edit': '编辑',
      'list': '列表',
      '404': '页面不存在'
    }

    let currentPath = ''
    pathSegments.forEach((segment, index) => {
      currentPath += `/${segment}`
      const title = pathMap[segment] || segment
      
      if (index === pathSegments.length - 1) {
        items.push({ title })
      } else {
        items.push({
          title: <Link to={currentPath}>{title}</Link>,
        })
      }
    })

    return items
  }

  // 不在登录页和404页显示面包屑
  if (location.pathname === '/login' || location.pathname === '/404') {
    return null
  }

  return (
    <Breadcrumb
      style={{ margin: '16px 0' }}
      items={getBreadcrumbItems()}
    />
  )
}

export default CustomBreadcrumb
