package repo

import (
	"habit/internal/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserChallengeRepository struct {
	db *gorm.DB
}

func NewUserChallengeRepository(db *gorm.DB) *UserChallengeRepository {
	return &UserChallengeRepository{db: db}
}

// Create 创建用户挑战记录
func (r *UserChallengeRepository) Create(userChallenge *model.AppUserChallenge) error {
	return r.db.Create(userChallenge).Error
}

// FindByID 根据ID查找用户挑战
func (r *UserChallengeRepository) FindByID(id int64) (*model.AppUserChallenge, error) {
	var userChallenge model.AppUserChallenge
	err := r.db.Where("id = ?", id).First(&userChallenge).Error
	if err != nil {
		return nil, err
	}
	return &userChallenge, nil
}

// FindByUserIDAndStatus 根据用户ID和状态查找用户挑战
func (r *UserChallengeRepository) FindByUserIDAndStatus(userID int64, status int8) (*model.AppUserChallenge, error) {
	var userChallenge model.AppUserChallenge
	err := r.db.Where("user_id = ? AND status = ?", userID, status).First(&userChallenge).Error
	if err != nil {
		return nil, err
	}
	return &userChallenge, nil
}

// UpdateAmount 更新预充值金额
func (r *UserChallengeRepository) UpdateAmount(id int64, newAmount decimal.Decimal) error {
	return r.db.Model(&model.AppUserChallenge{}).
		Where("id = ?", id).
		Update("pre_recharge", newAmount).Error
}

// UpdatePreRecharge 更新预充值金额
func (r *UserChallengeRepository) UpdatePreRecharge(id int64, newPreRecharge decimal.Decimal) error {
	return r.db.Model(&model.AppUserChallenge{}).
		Where("id = ?", id).
		Update("pre_recharge", newPreRecharge).Error
}

// UpdateBothAmounts 同时更新挑战金额和预充值金额
func (r *UserChallengeRepository) UpdateBothAmounts(id int64, preRecharge decimal.Decimal) error {
	updates := map[string]interface{}{
		"pre_recharge":     preRecharge,
	}
	return r.db.Model(&model.AppUserChallenge{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateStatus 更新挑战状态
func (r *UserChallengeRepository) UpdateStatus(id int64, status int8, failReason int8) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == 3 { // 失败状态
		updates["fail_reason"] = failReason
		updates["finished_at"] = time.Now()
	} else if status == 2 { // 成功状态
		updates["finished_at"] = time.Now()
	}
	
	return r.db.Model(&model.AppUserChallenge{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetTodayAndTomorrowChallenges 一次性获取用户今天和明天的挑战记录
func (r *UserChallengeRepository) GetTodayAndTomorrowChallenges(userID int64) ([]*model.AppUserChallenge, error) {
	var challenges []*model.AppUserChallenge
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	
	err := r.db.Where("user_id = ? AND DATE(start_date) IN (?, ?)", userID, today, tomorrow).
		Order("start_date ASC").
		Find(&challenges).Error
	
	if err != nil {
		return nil, err
	}
	return challenges, nil
}

// GetTodayChallenge 获取用户今天的挑战记录
func (r *UserChallengeRepository) GetTodayChallenge(userID int64) (*model.AppUserChallenge, error) {
	var userChallenge model.AppUserChallenge
	today := time.Now().Format("2006-01-02")
	
	err := r.db.Where("user_id = ? AND DATE(start_date) = ?", userID, today).First(&userChallenge).Error
	if err != nil {
		return nil, err
	}
	return &userChallenge, nil
}

// GetTomorrowChallenge 获取用户明天的挑战记录
func (r *UserChallengeRepository) GetTomorrowChallenge(userID int64) (*model.AppUserChallenge, error) {
	var userChallenge model.AppUserChallenge
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	
	err := r.db.Where("user_id = ? AND DATE(start_date) = ?", userID, tomorrow).First(&userChallenge).Error
	if err != nil {
		return nil, err
	}
	return &userChallenge, nil
}

// GetActiveChallengeByUserID 获取用户的活跃挑战
func (r *UserChallengeRepository) GetActiveChallengeByUserID(userID int64) (*model.AppUserChallenge, error) {
	var userChallenge model.AppUserChallenge
	err := r.db.Where("user_id = ? AND status = ?", userID, 1).First(&userChallenge).Error
	if err != nil {
		return nil, err
	}
	return &userChallenge, nil
}
