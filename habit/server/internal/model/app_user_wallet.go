package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// AppUserWallet 用户钱包表
type AppUserWallet struct {
	UserID    int64           `gorm:"column:user_id;uniqueIndex:uk_user_id"`
	PayPwd    string          `gorm:"column:pay_pwd;not null;default:''"`
	PayStatus string          `gorm:"column:pay_status;size:1;not null;default:'1'"` // 1-启用 2-禁用
	Balance   decimal.Decimal `gorm:"column:balance;type:decimal(30,2);not null;default:0.00"`
	Frozen    decimal.Decimal `gorm:"column:frozen;type:decimal(30,2);not null;default:0.00"`
	TotalR    decimal.Decimal `gorm:"column:total_r;type:decimal(30,2);not null;default:0.00"`  // 总充值
	TotalW    decimal.Decimal `gorm:"column:total_w;type:decimal(30,2);not null;default:0.00"`  // 总提现
	TotalRe   decimal.Decimal `gorm:"column:total_re;type:decimal(30,2);not null;default:0.00"` // 打卡总收益
	TotalI    decimal.Decimal `gorm:"column:total_i;type:decimal(30,2);not null;default:0.00"`  // 邀请总收益
	Address   string          `gorm:"column:address;not null;default:''" json:"address"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (AppUserWallet) TableName() string {
	return "app_user_wallet"
}
