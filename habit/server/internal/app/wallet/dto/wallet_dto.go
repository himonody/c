package dto

// WalletInfo 钱包信息
type WalletInfo struct {
	UserId    int64  `json:"userId"`
	PayStatus string `json:"payStatus"` // 1-启用 2-禁用
	Balance   string `json:"balance"`   // 余额（字符串，避免科学计数）
	Frozen    string `json:"frozen"`    // 冻结金额
	TotalR    string `json:"totalR"`    // 总充值
	TotalW    string `json:"totalW"`    // 总提现
	TotalRe   string `json:"totalRe"`   // 打卡总收益
	TotalI    string `json:"totalI"`    // 邀请总收益
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SetWalletAddressRequest struct {
	Address string `json:"address"`
}
