package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type AppUserBalanceLog struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BizID        string          `gorm:"column:biz_id;size:64;not null;default:'';index:idx_biz_id" json:"bizId"`
	UserID       int64           `gorm:"column:user_id;not null;default:0;index:idx_user_type_created" json:"userId"`
	Type         int             `gorm:"column:type;not null;default:0;index:idx_user_type_created" json:"type"`
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(32,8);not null;default:0.00000000" json:"amount"`
	BeforeBalance decimal.Decimal `gorm:"column:before_balance;type:decimal(32,8);not null;default:0.00000000" json:"beforeBalance"`
	AfterBalance  decimal.Decimal `gorm:"column:after_balance;type:decimal(32,8);not null;default:0.00000000" json:"afterBalance"`
	BeforeFrozen  decimal.Decimal `gorm:"column:before_frozen;type:decimal(32,8);not null;default:0.00000000" json:"beforeFrozen"`
	AfterFrozen   decimal.Decimal `gorm:"column:after_frozen;type:decimal(32,8);not null;default:0.00000000" json:"afterFrozen"`
	Remark       string          `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	OperatorID   int64           `gorm:"column:operator_id;not null;default:0" json:"operatorId"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (AppUserBalanceLog) TableName() string {
	return "app_user_balance_log"
}
