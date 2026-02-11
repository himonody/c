package repo

import (
	"context"
	"habit/internal/model"
	"habit/pkg/cache"
	"habit/pkg/database"

	"gorm.io/gorm"
)

type ChallengeStatRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewChallengeStatRepository(db *gorm.DB) *ChallengeStatRepository {
	return &ChallengeStatRepository{
		db:    db,
		cache: cache.NewCache(database.RedisClient),
	}
}

func (r *ChallengeStatRepository) GetTotalStat() (*model.AppChallengeTotalStat, error) {
	ctx := context.Background()
	cacheKey := cache.ChallengeTotalStatKey()

	// 尝试从缓存读取
	var cached model.AppChallengeTotalStat
	err := r.cache.Get(ctx, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	var stat model.AppChallengeTotalStat
	if err := r.db.Where("id = ?", 1).First(&stat).Error; err != nil {
		return nil, err
	}

	// 回写缓存
	_ = r.cache.Set(ctx, cacheKey, &stat, cache.ChallengeStatExpiration)
	return &stat, nil
}
