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
	ID                   int    `json:"id"`
	IsAutoSettle         bool   `json:"isAutoSettle"`
	SettleTime           string `json:"settleTime"`
	CycleDays            int    `json:"cycleDays"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	MaxDepositAmount     int    `json:"maxDepositAmount"`
	MinWithdrawAmount    int    `json:"minWithdrawAmount"`
	MaxDailyProfit       int    `json:"maxDailyProfit"`
	ExcessTaxRate        int    `json:"excessTaxRate"`
	MinDailyProfit       int    `json:"minDailyProfit"`
	DailyPlatformSubsidy int    `json:"dailyPlatformSubsidy"`
	UncheckDeductRate    int    `json:"uncheckDeductRate"`
	MinUncheckUsers      int    `json:"minUncheckUsers"`
	CommissionFollow     int    `json:"commissionFollow"`
	CommissionJoin       int    `json:"commissionJoin"`
	CommissionL1         int    `json:"commissionL1"`
	CommissionL2         int    `json:"commissionL2"`
	CommissionL3         int    `json:"commissionL3"`
	UpdatedAt            string `json:"updatedAt"`
}

type ChallengeListResponse struct {
	List     []*ChallengeInfo `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

type ChallengeUpsertRequest struct {
	ID                   int64  `json:"id"`
	IsAutoSettle         bool   `json:"isAutoSettle"`
	SettleTime           string `json:"settleTime"`
	CycleDays            int    `json:"cycleDays"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	MaxDepositAmount     int    `json:"maxDepositAmount"`
	MinWithdrawAmount    int    `json:"minWithdrawAmount"`
	MaxDailyProfit       int    `json:"maxDailyProfit"`
	ExcessTaxRate        int    `json:"excessTaxRate"`
	MinDailyProfit       int    `json:"minDailyProfit"`
	DailyPlatformSubsidy int    `json:"dailyPlatformSubsidy"`
	UncheckDeductRate    int    `json:"uncheckDeductRate"`
	MinUncheckUsers      int    `json:"minUncheckUsers"`
	CommissionFollow     int    `json:"commissionFollow"`
	CommissionJoin       int    `json:"commissionJoin"`
	CommissionL1         int    `json:"commissionL1"`
	CommissionL2         int    `json:"commissionL2"`
	CommissionL3         int    `json:"commissionL3"`
}
