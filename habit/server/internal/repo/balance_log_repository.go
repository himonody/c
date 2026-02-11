package repo

import (
	"habit/internal/model"

	"gorm.io/gorm"
)

type BalanceLogRepository struct {
	db *gorm.DB
}

func NewBalanceLogRepository(db *gorm.DB) *BalanceLogRepository {
	return &BalanceLogRepository{db: db}
}

func (r *BalanceLogRepository) Create(tx *gorm.DB, log *model.AppUserBalanceLog) error {
	if tx != nil {
		return tx.Create(log).Error
	}
	return r.db.Create(log).Error
}
