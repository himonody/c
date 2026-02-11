package dto

type ChallengeListRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type ChallengePoolInfo struct {
	TotalAmount string `json:"totalAmount"`
	Settled     int8   `json:"settled"`
}

type ChallengeInfo struct {
	ID            int64              `json:"id"`
	DayCount      int                `json:"dayCount"`
	Amount        string             `json:"amount"`
	CheckinStart  uint16             `json:"checkinStart"`
	CheckinEnd    uint16             `json:"checkinEnd"`
	PlatformBonus string             `json:"platformBonus"`
	Status        int8               `json:"status"`
	Sort          int8               `json:"sort"`
	CreatedAt     string             `json:"createdAt"`
	Pool          *ChallengePoolInfo `json:"pool"`
}

type ChallengeListResponse struct {
	List     []*ChallengeInfo `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

type ChallengeUpsertRequest struct {
	ID            int64  `json:"id"`
	DayCount      int    `json:"dayCount"`
	Amount        string `json:"amount"`
	CheckinStart  uint16 `json:"checkinStart"`
	CheckinEnd    uint16 `json:"checkinEnd"`
	PlatformBonus string `json:"platformBonus"`
	Status        int8   `json:"status"`
	Sort          int8   `json:"sort"`

	PoolID    int64  `json:"poolId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
