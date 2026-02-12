package repo

import (
	"context"
	"habit/internal/model"
	"habit/pkg/cache"
	"habit/pkg/database"
	"time"

	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db  *gorm.DB
	rdb *cache.Cache
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{
		db:  db,
		rdb: cache.NewCache(database.RedisClient),
	}
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
	if err := query.Offset(offset).Limit(pageSize).Order("updated_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ChallengeRepository) Last() (*model.AppChallenge, error) {
	var last *model.AppChallenge
	_ = r.rdb.Get(context.Background(), cache.AppChallengeKey(), &last)
	if last == nil {
		if err := r.db.Model(&model.AppChallenge{}).Last(&last).Error; err != nil {
			return nil, err
		}
		_ = r.rdb.Set(context.Background(), cache.AppChallengeKey(), last, time.Duration(24)*time.Hour)
	}
	return last, nil
}
