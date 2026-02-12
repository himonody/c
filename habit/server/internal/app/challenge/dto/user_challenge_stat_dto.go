package dto

// UserChallengeStatRequest 用户挑战统计请求
type UserChallengeStatRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// UserChallengeStatInfo 用户挑战统计信息
type UserChallengeStatInfo struct {
	UserID           int64  `json:"userId"`           // 用户 ID
	UserName         string `json:"userName"`         // 用户名称
	LatestCheckinAt  string `json:"latestCheckinAt"`  // 最新打卡时间
	TotalCheckinDays int    `json:"totalCheckinDays"` // 累计打卡天数
	TotalProfit      string `json:"totalProfit"`      // 累计收益
}

// UserChallengeStatResponse 用户挑战统计响应
type UserChallengeStatResponse struct {
	List     []*UserChallengeStatInfo `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}
