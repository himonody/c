package repo

import (
	"habit/internal/model"

	"gorm.io/gorm"
)

type InviteRepository struct {
	db *gorm.DB
}

func NewInviteRepository(db *gorm.DB) *InviteRepository {
	return &InviteRepository{db: db}
}

func (r *InviteRepository) FindByUserID(userID int64) (*model.AppUser, error) {
	var u model.AppUser
	if err := r.db.Where("id = ?", userID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *InviteRepository) ListInvitedUsers(refID int, page, pageSize int) ([]*model.AppUser, int64, error) {
	var list []*model.AppUser
	var total int64

	query := r.db.Model(&model.AppUser{}).Where("ref_id = ?", refID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
