package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type AppUserWithdraw struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BizID        string          `gorm:"column:biz_id;size:64;not null;default:'';uniqueIndex:unq_biz_id" json:"bizId"`
	UserID       int64           `gorm:"column:user_id;not null;default:0;index:idx_user_updated" json:"userId"`
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(32,8);not null;default:0.00000000" json:"amount"`
	Fee          decimal.Decimal `gorm:"column:fee;type:decimal(32,8);not null;default:0.00000000" json:"fee"`
	ActualAmount decimal.Decimal `gorm:"column:actual_amount;type:decimal(32,8);not null;default:0.00000000" json:"actualAmount"`
	Address      string          `gorm:"column:address;size:128;not null;default:''" json:"address"`
	ApplyIP      string          `gorm:"column:apply_ip;size:45;not null;default:''" json:"applyIp"`
	Status       int             `gorm:"column:status;not null;default:1;index:idx_status_updated" json:"status"`
	RejectReason string          `gorm:"column:reject_reason;size:255;not null;default:''" json:"rejectReason"`
	TxHash       string          `gorm:"column:tx_hash;size:128;not null;default:''" json:"txHash"`
	ReviewID     int64           `gorm:"column:review_id;not null;default:0" json:"reviewId"`
	ReviewIP     string          `gorm:"column:review_ip;size:45;not null;default:''" json:"reviewIp"`
	ReviewedAt   *time.Time      `gorm:"column:reviewed_at" json:"reviewedAt"`
	Version      int             `gorm:"column:version;not null;default:0" json:"version"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AppUserWithdraw) TableName() string {
	return "app_user_withdraw"
}
