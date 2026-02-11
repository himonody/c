package repo

import (
	"habit/internal/model"

	"gorm.io/gorm"
)

type ChallengePoolRepository struct {
	db *gorm.DB
}

func NewChallengePoolRepository(db *gorm.DB) *ChallengePoolRepository {
	return &ChallengePoolRepository{db: db}
}

func (r *ChallengePoolRepository) Create(tx *gorm.DB, pool *model.AppChallengePool) error {
	if tx != nil {
		return tx.Create(pool).Error
	}
	return r.db.Create(pool).Error
}

func (r *ChallengePoolRepository) Update(tx *gorm.DB, pool *model.AppChallengePool) error {
	if tx != nil {
		return tx.Save(pool).Error
	}
	return r.db.Save(pool).Error
}

func (r *ChallengePoolRepository) FindLatestByChallengeIDs(challengeIDs []int64) ([]*model.AppChallengePool, error) {
	if len(challengeIDs) == 0 {
		return []*model.AppChallengePool{}, nil
	}
	var list []*model.AppChallengePool
	// NOTE: For admin list purpose, we pick latest pool per challenge by max(id).
	// Use a subquery to get max(id) for each challenge_id.
	sub := r.db.Model(&model.AppChallengePool{}).Select("MAX(id) as id").Where("challenge_id IN ?", challengeIDs).Group("challenge_id")
	if err := r.db.Model(&model.AppChallengePool{}).Where("id IN (?)", sub).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
