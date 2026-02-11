package repo

import (
	"habit/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}

func (r *ChallengeRepository) Create(tx *gorm.DB, challenge *model.AppChallenge) error {
	if tx != nil {
		return tx.Create(challenge).Error
	}
	return r.db.Create(challenge).Error
}

func (r *ChallengeRepository) Update(tx *gorm.DB, challenge *model.AppChallenge) error {
	if tx != nil {
		return tx.Save(challenge).Error
	}
	return r.db.Save(challenge).Error
}

func (r *ChallengeRepository) FindByID(id int64) (*model.AppChallenge, error) {
	var c model.AppChallenge
	if err := r.db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ChallengeRepository) List(page, pageSize int) ([]*model.AppChallenge, int64, error) {
	var list []*model.AppChallenge
	var total int64

	query := r.db.Model(&model.AppChallenge{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

type ChallengeWithLatestPool struct {
	model.AppChallenge
	PoolTotalAmount decimal.Decimal `gorm:"column:pool_total_amount"`
	PoolSettled     int8            `gorm:"column:pool_settled"`
}

func (r *ChallengeRepository) ListWithLatestPool(page, pageSize int) ([]*ChallengeWithLatestPool, int64, error) {
	var list []*ChallengeWithLatestPool
	var total int64

	base := r.db.Table("app_challenge AS c")
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := r.db.Table("app_challenge AS c").Select("c.*, p.total_amount AS pool_total_amount, p.settled AS pool_settled")
	query = query.Joins("LEFT JOIN app_challenge_pool p ON p.challenge_id = c.id)")
	if err := query.Offset(offset).Limit(pageSize).Order("c.created_at DESC").Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
