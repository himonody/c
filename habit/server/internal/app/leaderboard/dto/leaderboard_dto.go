package dto

type LeaderboardListRequest struct {
	RankType int8   `json:"rankType"`
	RankDate string `json:"rankDate"` // YYYY-MM-DD，不传默认当天
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type LeaderboardTotalRequest struct {
	RankType int8 `json:"rankType"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
}

type LeaderboardItem struct {
	UserID int64  `json:"userId"`
	Value  string `json:"value"`
	RankNo int    `json:"rankNo"`
}

type LeaderboardListResponse struct {
	List     []*LeaderboardItem `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

type LeaderboardTotalResponse struct {
	List     []*LeaderboardItem `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}
