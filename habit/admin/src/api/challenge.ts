import request from '@/utils/request'

// 挑战管理相关接口
export interface ChallengeListRequest {
  page: number
  pageSize: number
}

export interface ChallengeInfo {
  id: number
  isAutoSettle: boolean
  settleTime: string
  cycleDays: number
  startTime: string
  endTime: string
  maxDepositAmount: number
  minWithdrawAmount: number
  maxDailyProfit: number
  excessTaxRate: number
  minDailyProfit: number
  dailyPlatformSubsidy: number
  uncheckDeductRate: number
  minUncheckUsers: number
  commissionFollow: number
  commissionJoin: number
  commissionL1: number
  commissionL2: number
  commissionL3: number
  updatedAt: string
}

export interface ChallengeListResponse {
  list: ChallengeInfo[]
  total: number
  page: number
  pageSize: number
}

export interface ChallengeUpsertRequest {
  id?: number
  isAutoSettle: boolean
  settleTime: string
  cycleDays: number
  startTime: string
  endTime: string
  maxDepositAmount: number
  minWithdrawAmount: number
  maxDailyProfit: number
  excessTaxRate: number
  minDailyProfit: number
  dailyPlatformSubsidy: number
  uncheckDeductRate: number
  minUncheckUsers: number
  commissionFollow: number
  commissionJoin: number
  commissionL1: number
  commissionL2: number
  commissionL3: number
}

// 获取挑战列表
export const getChallengeList = (params: ChallengeListRequest) => {
  return request.post<ChallengeListResponse>('/admin/challenge/list', params)
}

// 创建挑战
export const createChallenge = (data: ChallengeUpsertRequest) => {
  return request.post('/admin/challenge/create', data)
}

// 更新挑战
export const updateChallenge = (data: ChallengeUpsertRequest) => {
  return request.post('/admin/challenge/update', data)
}

// 删除挑战
export const deleteChallenge = (id: number) => {
  return request.post('/admin/challenge/delete', { id })
}
