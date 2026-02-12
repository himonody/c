import request from '@/utils/request'

// 配置管理相关接口
export interface ConfigListRequest {
  page: number
  pageSize: number
  configName?: string
  configKey?: string
}

export interface ConfigInfo {
  id: number
  configName: string
  configKey: string
  configValue: string
  configType: string
  isFrontend: string
  remark: string
  createBy: number
  updateBy: number
  createdAt: string
  updatedAt: string
}

export interface ConfigListResponse {
  list: ConfigInfo[]
  total: number
  page: number
  pageSize: number
}

export interface CreateConfigRequest {
  configName: string
  configKey: string
  configValue?: string
  configType?: string
  isFrontend?: string
  remark?: string
}

export interface UpdateConfigRequest {
  id: number
  configName: string
  configKey: string
  configValue?: string
  configType?: string
  isFrontend?: string
  remark?: string
}

export interface GetConfigRequest {
  id: number
}

// 获取配置列表
export const getConfigList = (params: ConfigListRequest) => {
  return request.post<ConfigListResponse>('/admin/config/list', params)
}

// 获取配置详情
export const getConfig = (params: GetConfigRequest) => {
  return request.post<ConfigInfo>('/admin/config/get', params)
}

// 创建配置
export const createConfig = (data: CreateConfigRequest) => {
  return request.post('/admin/config/create', data)
}

// 更新配置
export const updateConfig = (data: UpdateConfigRequest) => {
  return request.post('/admin/config/update', data)
}

// 删除配置
export const deleteConfig = (id: number) => {
  return request.post('/admin/config/delete', { id })
}
