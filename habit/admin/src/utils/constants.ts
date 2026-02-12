// API 端点常量
export const API_ENDPOINTS = {
  // 认证相关
  AUTH: {
    LOGIN: '/admin/auth/login',
    LOGOUT: '/admin/auth/logout',
    ME: '/admin/auth/me',
  },
  
  // 挑战管理
  CHALLENGE: {
    LIST: '/admin/challenge/list',
    CREATE: '/admin/challenge/create',
    UPDATE: '/admin/challenge/update',
    DELETE: '/admin/challenge/delete',
    GET: '/admin/challenge/get',
  },
  
  // 配置管理
  CONFIG: {
    LIST: '/admin/config/list',
    CREATE: '/admin/config/create',
    UPDATE: '/admin/config/update',
    DELETE: '/admin/config/delete',
    GET: '/admin/config/get',
  },
  
  // 用户管理
  USER: {
    LIST: '/admin/user/list',
    CREATE: '/admin/user/create',
    UPDATE: '/admin/user/update',
    DELETE: '/admin/user/delete',
    GET: '/admin/user/get',
    RESET_PASSWORD: '/admin/user/reset-password',
  },
} as const

// HTTP 状态码
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
} as const

// 响应代码
export const RESPONSE_CODE = {
  SUCCESS: 0,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_ERROR: 500,
} as const

// 分页默认配置
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
} as const

// 文件上传配置
export const UPLOAD = {
  MAX_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_TYPES: ['image/jpeg', 'image/png', 'image/gif'],
  MAX_COUNT: 1,
} as const

// 缓存键名
export const CACHE_KEYS = {
  USER_INFO: 'user_info',
  PERMISSIONS: 'permissions',
  CONFIG: 'config',
} as const

// 时间格式
export const DATE_FORMATS = {
  DATE: 'YYYY-MM-DD',
  TIME: 'HH:mm:ss',
  DATETIME: 'YYYY-MM-DD HH:mm:ss',
  MONTH: 'YYYY-MM',
  YEAR: 'YYYY',
} as const

// 正则表达式
export const REGEX = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE: /^1[3-9]\d{9}$/,
  PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$/,
  USERNAME: /^[a-zA-Z0-9_]{3,20}$/,
  CONFIG_KEY: /^[a-zA-Z][a-zA-Z0-9_]*$/,
} as const

// 消息提示
export const MESSAGES = {
  SUCCESS: '操作成功',
  ERROR: '操作失败',
  NETWORK_ERROR: '网络错误',
  UNAUTHORIZED: '未授权访问',
  FORBIDDEN: '权限不足',
  NOT_FOUND: '资源不存在',
  SERVER_ERROR: '服务器错误',
  CONFIRM_DELETE: '确定要删除吗？',
  SAVE_SUCCESS: '保存成功',
  SAVE_ERROR: '保存失败',
  LOAD_ERROR: '加载失败',
} as const
