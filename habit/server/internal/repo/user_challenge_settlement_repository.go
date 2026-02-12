package repo

import (
	"habit/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserChallengeSettlementRepository struct {
	db *gorm.DB
}

func NewUserChallengeSettlementRepository(db *gorm.DB) *UserChallengeSettlementRepository {
	return &UserChallengeSettlementRepository{db: db}
}

// SumFinalProfitByUserID 统计用户累计收益
func (r *UserChallengeSettlementRepository) SumFinalProfitByUserID(userID int64) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.Model(&model.AppUserChallengeSettlement{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(final_profit), 0)").
		Scan(&total).Error
	return total, err
}

// Create 创建结算记录
func (r *UserChallengeSettlementRepository) Create(settlement *model.AppUserChallengeSettlement) error {
	return r.db.Create(settlement).Error
}
