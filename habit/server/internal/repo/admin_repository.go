package repo

import (
	"context"
	"errors"

	"habit/internal/model"
	"habit/pkg/cache"
	"habit/pkg/database"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		db:    db,
		cache: cache.NewCache(database.RedisClient),
	}
}

// FindByUsername finds an admin user by username (with cache)
func (r *AdminRepository) FindByUsername(username string) (*model.AdminUser, error) {
	ctx := context.Background()
	cacheKey := cache.AdminByUsernameCacheKey(username)

	// Try to get from cache first
	var admin model.AdminUser
	err := r.cache.Get(ctx, cacheKey, &admin)
	if err == nil {
		return &admin, nil
	}

	// Cache miss or error, query from database
	if !errors.Is(err, redis.Nil) {
		// Log cache error but continue
	}

	err = r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = r.cache.Set(ctx, cacheKey, &admin, cache.AdminCacheExpiration)

	return &admin, nil
}

// FindByID finds an admin user by ID (with cache)
func (r *AdminRepository) FindByID(id int64) (*model.AdminUser, error) {
	ctx := context.Background()
	cacheKey := cache.AdminCacheKey(id)

	// Try to get from cache first
	var admin model.AdminUser
	err := r.cache.Get(ctx, cacheKey, &admin)
	if err == nil {
		return &admin, nil
	}

	// Cache miss or error, query from database
	if !errors.Is(err, redis.Nil) {
		// Log cache error but continue
	}

	err = r.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = r.cache.Set(ctx, cacheKey, &admin, cache.AdminCacheExpiration)

	return &admin, nil
}

// Create creates a new admin user
func (r *AdminRepository) Create(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}

// Update updates admin user
func (r *AdminRepository) Update(admin *model.AdminUser) error {
	err := r.db.Save(admin).Error
	if err == nil {
		// Invalidate cache
		r.invalidateAdminCache(admin.ID)
	}
	return err
}

// invalidateAdminCache 清除管理员缓存
func (r *AdminRepository) invalidateAdminCache(adminID int64) {
	ctx := context.Background()
	cacheKey := cache.AdminCacheKey(adminID)
	_ = r.cache.Delete(ctx, cacheKey)
}
