package service

import (
	"errors"
	"time"

	"habit/internal/admin/config/dto"
	"habit/internal/model"
	"habit/internal/repo"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConfigService struct {
	configRepo *repo.ConfigRepository
	logger     *zap.Logger
}

func NewConfigService(configRepo *repo.ConfigRepository, logger *zap.Logger) *ConfigService {
	return &ConfigService{
		configRepo: configRepo,
		logger:     logger,
	}
}

// GetConfigList 获取配置列表
func (s *ConfigService) GetConfigList(req *dto.ConfigListRequest) (*dto.ConfigListResponse, error) {
	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	configs, total, err := s.configRepo.List(req.Page, req.PageSize, req.ConfigName, req.ConfigKey)
	if err != nil {
		s.logger.Error("Failed to get config list", zap.Error(err))
		return nil, errors.New("failed to get config list")
	}

	// 转换为 DTO
	list := make([]*dto.ConfigInfo, 0, len(configs))
	for _, config := range configs {
		list = append(list, &dto.ConfigInfo{
			ID:          config.ID,
			ConfigName:  config.ConfigName,
			ConfigKey:   config.ConfigKey,
			ConfigValue: config.ConfigValue,
			ConfigType:  config.ConfigType,
			IsFrontend:  config.IsFrontend,
			Remark:      config.Remark,
			CreateBy:    config.CreateBy,
			UpdateBy:    config.UpdateBy,
			CreatedAt:   config.CreatedAt.Format(time.DateTime),
			UpdatedAt:   config.UpdatedAt.Format(time.DateTime),
		})
	}

	return &dto.ConfigListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetConfig 获取配置详情
func (s *ConfigService) GetConfig(id int64) (*dto.ConfigInfo, error) {
	config, err := s.configRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("config not found")
		}
		s.logger.Error("Failed to get config", zap.Error(err))
		return nil, errors.New("failed to get config")
	}

	return &dto.ConfigInfo{
		ID:          config.ID,
		ConfigName:  config.ConfigName,
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  config.ConfigType,
		IsFrontend:  config.IsFrontend,
		Remark:      config.Remark,
		CreateBy:    config.CreateBy,
		UpdateBy:    config.UpdateBy,
		CreatedAt:   config.CreatedAt.Format(time.DateTime),
		UpdatedAt:   config.UpdatedAt.Format(time.DateTime),
	}, nil
}

// CreateConfig 创建配置
func (s *ConfigService) CreateConfig(req *dto.CreateConfigRequest, adminID int64) error {
	// 检查配置键是否已存在
	exists, err := s.configRepo.ConfigKeyExists(req.ConfigKey, 0)
	if err != nil {
		s.logger.Error("Failed to check config key", zap.Error(err))
		return errors.New("internal server error")
	}
	if exists {
		return errors.New("config key already exists")
	}

	config := &model.AdminSysConfig{
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		IsFrontend:  req.IsFrontend,
		Remark:      req.Remark,
		CreateBy:    adminID,
		UpdateBy:    adminID,
	}

	if err := s.configRepo.Create(config); err != nil {
		s.logger.Error("Failed to create config", zap.Error(err))
		return errors.New("failed to create config")
	}

	return nil
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(req *dto.UpdateConfigRequest, adminID int64) error {
	// 检查配置是否存在
	config, err := s.configRepo.FindByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("config not found")
		}
		s.logger.Error("Failed to get config", zap.Error(err))
		return errors.New("internal server error")
	}

	// 检查配置键是否与其他配置重复
	exists, err := s.configRepo.ConfigKeyExists(req.ConfigKey, req.ID)
	if err != nil {
		s.logger.Error("Failed to check config key", zap.Error(err))
		return errors.New("internal server error")
	}
	if exists {
		return errors.New("config key already exists")
	}

	// 更新配置
	config.ConfigName = req.ConfigName
	config.ConfigKey = req.ConfigKey
	config.ConfigValue = req.ConfigValue
	config.ConfigType = req.ConfigType
	config.IsFrontend = req.IsFrontend
	config.Remark = req.Remark
	config.UpdateBy = adminID

	if err := s.configRepo.Update(config); err != nil {
		s.logger.Error("Failed to update config", zap.Error(err))
		return errors.New("failed to update config")
	}

	return nil
}
