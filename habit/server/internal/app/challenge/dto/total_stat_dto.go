package dto

type ChallengeTotalStatResponse struct {
	TotalUserCnt       int    `json:"totalUserCnt"`
	TotalJoinCnt       int    `json:"totalJoinCnt"`
	TotalSuccessCnt    int    `json:"totalSuccessCnt"`
	TotalFailCnt       int    `json:"totalFailCnt"`
	TotalJoinAmount    string `json:"totalJoinAmount"`
	TotalSuccessAmount string `json:"totalSuccessAmount"`
	TotalFailAmount    string `json:"totalFailAmount"`
	TotalPlatformBonus string `json:"totalPlatformBonus"`
	TotalPoolAmount    string `json:"totalPoolAmount"`
	UpdatedAt          string `json:"updatedAt"`
}
