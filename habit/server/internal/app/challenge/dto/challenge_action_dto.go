package dto

// ChallengeStartRequest 开始挑战请求
type ChallengeStartRequest struct {
	PreRecharge string `json:"preRecharge"` // 预充值金额（字符串格式）
}

// ChallengeStartResponse 开始挑战响应
type ChallengeStartResponse struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	ChallengeID     int64  `json:"challengeId"`     // 挑战配置ID
	ChallengeAmount string `json:"challengeAmount"` // 挑战金额
	PreRecharge     string `json:"preRecharge"`     // 预充值金额
	StartDate       string `json:"startDate"`       // 开始日期
	EndDate         string `json:"endDate"`         // 结束日期
	Status          int8   `json:"status"`          // 状态：1进行中
}

// ChallengeMoneyRequest 增加挑战金请求
type ChallengeMoneyRequest struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	Amount          string `json:"amount"`          // 增加金额（字符串格式）
}

// ChallengeMoneyResponse 增加挑战金响应
type ChallengeMoneyResponse struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	ChallengeID     int64  `json:"challengeId"`     // 挑战配置ID
	ChallengeAmount string `json:"challengeAmount"` // 挑战金额
	PreRecharge     string `json:"preRecharge"`     // 预充值金额
	StartDate       string `json:"startDate"`       // 开始日期
	EndDate         string `json:"endDate"`         // 结束日期
	Status          int8   `json:"status"`          // 状态：1进行中
}

// ChallengeQueryRequest 查询挑战记录请求
type ChallengeQueryRequest struct {
	// 可以扩展其他查询条件，目前只需要用户ID（从认证中获取）
}

// ChallengeQueryResponse 查询挑战记录响应
type ChallengeQueryResponse struct {
	TodayChallenge *UserChallengeInfo `json:"todayChallenge"` // 今天的挑战记录
	TomorrowChallenge *UserChallengeInfo `json:"tomorrowChallenge"` // 明天的挑战记录
}

// UserChallengeInfo 用户挑战信息
type UserChallengeInfo struct {
	ID              int64  `json:"id"`              // 记录ID
	UserID          int64  `json:"userId"`          // 用户ID
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	ChallengeID     int64  `json:"challengeId"`     // 挑战配置ID
	ChallengeAmount string `json:"challengeAmount"` // 挑战金额
	PreRecharge     string `json:"preRecharge"`     // 预充值金额
	StartDate       string `json:"startDate"`       // 开始日期
	EndDate         string `json:"endDate"`         // 结束日期
	Status          int8   `json:"status"`          // 状态
	StatusText      string `json:"statusText"`      // 状态描述
	FailReason      int8   `json:"failReason"`      // 失败原因
	CreatedAt       string `json:"createdAt"`       // 创建时间
	FinishedAt      *string `json:"finishedAt"`     // 完成时间
}

// CheckinRequest 打卡请求
type CheckinRequest struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	MoodCode        int8   `json:"moodCode"`         // 心情枚举：1开心 2平静 3一般 4疲惫 5低落 6爆棚
	MoodText        string `json:"moodText"`         // 心情文字描述（最多200字）
}

// CheckinResponse 打卡响应
type CheckinResponse struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	CheckinTime     string `json:"checkinTime"`     // 打卡时间
	CheckinDate     string `json:"checkinDate"`     // 打卡日期
	Status          int8   `json:"status"`          // 打卡状态：1成功，2失败
	Message         string `json:"message"`         // 打卡消息
	ChallengeInfo   *UserChallengeInfo `json:"challengeInfo"` // 挑战信息
}
