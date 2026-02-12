// 状态枚举
export enum Status {
  DISABLED = 0,
  ENABLED = 1,
}

// 用户状态枚举
export enum UserStatus {
  INACTIVE = 0,
  ACTIVE = 1,
  BANNED = 2,
}

// 挑战状态枚举
export enum ChallengeStatus {
  PENDING = 0,
  ACTIVE = 1,
  COMPLETED = 2,
  CANCELLED = 3,
}

// 打卡状态枚举
export enum CheckinStatus {
  FAILED = 0,
  SUCCESS = 1,
  PENDING = 2,
}

// 性别枚举
export enum Gender {
  MALE = 1,
  FEMALE = 2,
  UNKNOWN = 3,
}

// 操作类型枚举
export enum ActionType {
  CREATE = 'create',
  UPDATE = 'update',
  DELETE = 'delete',
  VIEW = 'view',
}

// 日志级别枚举
export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
}

// 主题枚举
export enum Theme {
  LIGHT = 'light',
  DARK = 'dark',
}

// 语言枚举
export enum Language {
  ZH_CN = 'zh-CN',
  EN_US = 'en-US',
}
