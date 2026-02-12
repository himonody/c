// 应用配置
export const APP_CONFIG = {
  name: 'Habit Admin',
  version: '1.0.0',
  description: '习惯打卡管理后台系统',
}

// API 配置
export const API_CONFIG = {
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
}

// 分页配置
export const PAGINATION_CONFIG = {
  defaultPageSize: 10,
  pageSizeOptions: ['10', '20', '50', '100'],
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`,
}

// 表格配置
export const TABLE_CONFIG = {
  scroll: { x: 1200 },
  size: 'middle' as const,
  bordered: false,
}

// 主题配置
export const THEME_CONFIG = {
  primaryColor: '#1890ff',
  successColor: '#52c41a',
  warningColor: '#faad14',
  errorColor: '#f5222d',
  infoColor: '#1890ff',
}

// 本地存储键名
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USER_INFO: 'userInfo',
  LANGUAGE: 'language',
  THEME: 'theme',
  SIDEBAR_COLLAPSED: 'sidebarCollapsed',
}

// 路由配置
export const ROUTE_CONFIG = {
  home: '/dashboard',
  login: '/login',
  notFound: '/404',
}

// 文件上传配置
export const UPLOAD_CONFIG = {
  maxSize: 5 * 1024 * 1024, // 5MB
  acceptTypes: ['image/jpeg', 'image/png', 'image/gif'],
  maxCount: 1,
}
