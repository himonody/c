package repo

import (
	"habit/internal/model"

	"gorm.io/gorm"
)

type WithdrawRepository struct {
	db *gorm.DB
}

func NewWithdrawRepository(db *gorm.DB) *WithdrawRepository {
	return &WithdrawRepository{db: db}
}

func (r *WithdrawRepository) Create(tx *gorm.DB, withdraw *model.AppUserWithdraw) error {
	if tx != nil {
		return tx.Create(withdraw).Error
	}
	return r.db.Create(withdraw).Error
}

func (r *WithdrawRepository) ListByUserID(userID int64, page, pageSize int) ([]*model.AppUserWithdraw, int64, error) {
	var list []*model.AppUserWithdraw
	var total int64

	query := r.db.Model(&model.AppUserWithdraw{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
