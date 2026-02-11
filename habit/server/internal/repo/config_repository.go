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

type ConfigRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{
		db:    db,
		cache: cache.NewCache(database.RedisClient),
	}
}

// FindByID 根据ID查询配置（带缓存）
func (r *ConfigRepository) FindByID(id int64) (*model.AdminSysConfig, error) {
	ctx := context.Background()
	cacheKey := cache.SystemConfigIDKey(id)

	// 尝试从缓存获取
	var config model.AdminSysConfig
	err := r.cache.Get(ctx, cacheKey, &config)
	if err == nil {
		return &config, nil
	}

	// 缓存未命中，查询数据库
	if !errors.Is(err, redis.Nil) {
		// 记录缓存错误但继续
	}

	err = r.db.Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}

	// 存入缓存
	_ = r.cache.Set(ctx, cacheKey, &config, cache.SystemConfigExpiration)

	return &config, nil
}

// FindByKey 根据配置键查询配置（带缓存）
func (r *ConfigRepository) FindByKey(configKey string) (*model.AdminSysConfig, error) {
	ctx := context.Background()
	cacheKey := cache.SystemConfigKey(configKey)

	// 尝试从缓存获取
	var config model.AdminSysConfig
	err := r.cache.Get(ctx, cacheKey, &config)
	if err == nil {
		return &config, nil
	}

	// 缓存未命中，查询数据库
	if !errors.Is(err, redis.Nil) {
		// 记录缓存错误但继续
	}

	err = r.db.Where("config_key = ?", configKey).First(&config).Error
	if err != nil {
		return nil, err
	}

	// 存入缓存
	_ = r.cache.Set(ctx, cacheKey, &config, cache.SystemConfigExpiration)

	return &config, nil
}

// List 查询配置列表（分页）
func (r *ConfigRepository) List(page, pageSize int, configName, configKey string) ([]*model.AdminSysConfig, int64, error) {
	var configs []*model.AdminSysConfig
	var total int64

	query := r.db.Model(&model.AdminSysConfig{})

	// 条件查询
	if configName != "" {
		query = query.Where("config_name LIKE ?", "%"+configName+"%")
	}
	if configKey != "" {
		query = query.Where("config_key LIKE ?", "%"+configKey+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// Create 创建配置
func (r *ConfigRepository) Create(config *model.AdminSysConfig) error {
	return r.db.Create(config).Error
}

// Update 更新配置
func (r *ConfigRepository) Update(config *model.AdminSysConfig) error {
	err := r.db.Save(config).Error
	if err == nil {
		// 清除缓存
		r.invalidateConfigCache(config.ID, config.ConfigKey)
	}
	return err
}

// Delete 删除配置
func (r *ConfigRepository) Delete(id int64) error {
	var config model.AdminSysConfig
	if err := r.db.Where("id = ?", id).First(&config).Error; err != nil {
		return err
	}

	err := r.db.Delete(&config).Error
	if err == nil {
		// 清除缓存
		r.invalidateConfigCache(config.ID, config.ConfigKey)
	}
	return err
}

// ConfigKeyExists 检查配置键是否存在
func (r *ConfigRepository) ConfigKeyExists(configKey string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&model.AdminSysConfig{}).Where("config_key = ?", configKey)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// invalidateConfigCache 清除配置缓存
func (r *ConfigRepository) invalidateConfigCache(id int64, configKey string) {
	ctx := context.Background()
	_ = r.cache.Delete(ctx, cache.SystemConfigIDKey(id))
	_ = r.cache.Delete(ctx, cache.SystemConfigKey(configKey))
}
