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

type UserRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: cache.NewCache(database.RedisClient),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *model.AppUser) error {
	return r.db.Create(user).Error
}

// FindByUsername finds a user by username (with cache)
func (r *UserRepository) FindByUsername(username string) (*model.AppUser, error) {
	ctx := context.Background()
	cacheKey := cache.UserByUsernameCacheKey(username)

	// Try to get from cache first
	var user model.AppUser
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// Cache miss or error, query from database
	if !errors.Is(err, redis.Nil) {
		// Log cache error but continue
	}

	err = r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = r.cache.Set(ctx, cacheKey, &user, cache.UserCacheExpiration)

	return &user, nil
}

// FindByID finds a user by ID (with cache)
func (r *UserRepository) FindByID(id int64) (*model.AppUser, error) {
	ctx := context.Background()
	cacheKey := cache.UserCacheKey(id)

	// Try to get from cache first
	var user model.AppUser
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// Cache miss or error, query from database
	if !errors.Is(err, redis.Nil) {
		// Log cache error but continue
	}

	err = r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = r.cache.Set(ctx, cacheKey, &user, cache.UserCacheExpiration)

	return &user, nil
}

// FindByRefCode finds a user by reference code
func (r *UserRepository) FindByRefCode(refCode string) (*model.AppUser, error) {
	var user model.AppUser
	err := r.db.Where("ref_code = ?", refCode).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin updates user's last login time and IP
func (r *UserRepository) UpdateLastLogin(userID int64, ip string) error {
	err := r.db.Model(&model.AppUser{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error

	if err == nil {
		// Invalidate cache
		r.invalidateUserCache(userID)
	}

	return err
}

// UpdateOnlineStatus updates user's online status
func (r *UserRepository) UpdateOnlineStatus(userID int64, status string) error {
	err := r.db.Model(&model.AppUser{}).
		Where("id = ?", userID).
		Update("online_status", status).Error

	if err == nil {
		// Invalidate cache
		r.invalidateUserCache(userID)
	}

	return err
}

// UsernameExists checks if username already exists
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.AppUser{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(userID int64, newPasswordHash string) error {
	err := r.db.Model(&model.AppUser{}).
		Where("id = ?", userID).
		Update("pwd", newPasswordHash).Error

	if err == nil {
		// Invalidate cache
		r.invalidateUserCache(userID)
	}

	return err
}

// UpdateProfile 更新用户资料
func (r *UserRepository) UpdateProfile(userID int64, updates map[string]interface{}) error {
	err := r.db.Model(&model.AppUser{}).
		Where("id = ?", userID).
		Updates(updates).Error

	if err == nil {
		// Invalidate cache
		r.invalidateUserCache(userID)
	}

	return err
}

// invalidateUserCache 清除用户缓存
func (r *UserRepository) invalidateUserCache(userID int64) {
	ctx := context.Background()
	cacheKey := cache.UserCacheKey(userID)
	_ = r.cache.Delete(ctx, cacheKey)
}

// FindByIDs 根据ID列表查找用户
func (r *UserRepository) FindByIDs(userIDs []int64) ([]*model.AppUser, error) {
	if len(userIDs) == 0 {
		return []*model.AppUser{}, nil
	}
	
	var users []*model.AppUser
	err := r.db.Where("id IN ?", userIDs).Find(&users).Error
	return users, err
}

// ListWithPagination 分页查询所有用户
func (r *UserRepository) ListWithPagination(page, pageSize int) ([]*model.AppUser, int64, error) {
	var users []*model.AppUser
	var total int64
	
	offset := (page - 1) * pageSize
	
	if err := r.db.Model(&model.AppUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}
