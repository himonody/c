package dto

// WithdrawListRequest 提款记录列表请求
type WithdrawListRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type WithdrawInfo struct {
	BizID        string `json:"bizId"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
	ActualAmount string `json:"actualAmount"`
	Address      string `json:"address"`
	Status       int    `json:"status"`
	RejectReason string `json:"rejectReason"` //拒绝原因
	TxHash       string `json:"txHash"`       //交易 hash
	ReviewedAt   string `json:"reviewedAt"`   //审核时间
	CreatedAt    string `json:"createdAt"`    //申请时间
}

type WithdrawListResponse struct {
	List     []*WithdrawInfo `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

type WithdrawApplyRequest struct {
	Amount      string `json:"amount"`
	Address     string `json:"address"`
	PayPassword string `json:"payPassword"`
}

type WithdrawApplyResponse struct {
	BizID string `json:"bizId"`
}
