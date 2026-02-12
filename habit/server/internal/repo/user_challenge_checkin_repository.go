package repo

import (
	"habit/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserChallengeCheckinRepository struct {
	db *gorm.DB
}

func NewUserChallengeCheckinRepository(db *gorm.DB) *UserChallengeCheckinRepository {
	return &UserChallengeCheckinRepository{db: db}
}

// GetLatestByUserID 获取用户最新打卡记录
func (r *UserChallengeCheckinRepository) GetLatestByUserID(userID int64) (*model.AppUserChallengeCheckin, error) {
	var checkin model.AppUserChallengeCheckin
	err := r.db.Where("user_id = ?", userID).
		Order("checkin_time DESC, created_at DESC").
		First(&checkin).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

// CountDistinctCheckinDatesByUserID 统计用户累计打卡天数（去重日期）
func (r *UserChallengeCheckinRepository) CountDistinctCheckinDatesByUserID(userID int64) (int, error) {
	var count int64
	err := r.db.Model(&model.AppUserChallengeCheckin{}).
		Where("user_id = ? AND status = ?", userID, 1). // 只统计成功打卡
		Count(&count).Error
	return int(count), err
}

// GetUserIDsByLatestCheckin 分页获取按最新打卡时间排序的用户ID列表
func (r *UserChallengeCheckinRepository) GetUserIDsByLatestCheckin(page, pageSize int) ([]int64, int64, error) {
	var userIDs []int64
	var total int64

	// 获取有打卡记录的用户总数
	subQuery := r.db.Model(&model.AppUserChallengeCheckin{}).
		Select("user_id, MAX(checkin_time) as latest_checkin_time").
		Where("status = ?", 1). // 只统计成功打卡
		Group("user_id")

	if err := r.db.Table("(?) as latest_checkins", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页获取用户ID（按最新打卡时间排序）
	offset := (page - 1) * pageSize
	query := r.db.Model(&model.AppUserChallengeCheckin{}).
		Select("user_id").
		Where("id IN (?)",
			r.db.Table("(?) as latest_checkins", subQuery).
				Select("MIN(id)").
				Order("latest_checkin_time DESC, user_id ASC").
				Offset(offset).
				Limit(pageSize),
		).
		Order("checkin_time DESC, user_id ASC")

	if err := query.Find(&userIDs).Error; err != nil {
		return nil, 0, err
	}

	return userIDs, total, nil
}

// Create 创建打卡记录
func (r *UserChallengeCheckinRepository) Create(checkin *model.AppUserChallengeCheckin) error {
	return r.db.Create(checkin).Error
}

// GetTodayCheckin 获取用户今天的打卡记录
func (r *UserChallengeCheckinRepository) GetTodayCheckin(userID, userChallengeID int64) (*model.AppUserChallengeCheckin, error) {
	var checkin model.AppUserChallengeCheckin
	today := time.Now().Format("2006-01-02")

	err := r.db.Where("user_id = ? AND user_challenge_id = ? AND checkin_date = ?",
		userID, userChallengeID, today).First(&checkin).Error

	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

// UpdateCheckinSuccess 更新打卡为成功状态
func (r *UserChallengeCheckinRepository) UpdateCheckinSuccess(id int64) error {
	now := time.Now()
	return r.db.Model(&model.AppUserChallengeCheckin{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       1, // 打卡成功
			"checkin_time": now,
		}).Error
}

// CheckCanCheckin 检查是否可以打卡（在打卡时间范围内）
func (r *UserChallengeCheckinRepository) CheckCanCheckin(userChallengeID int64, startTime, endTime string) (bool, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 检查今天是否已经打卡成功
	var count int64
	err := r.db.Model(&model.AppUserChallengeCheckin{}).
		Where("user_challenge_id = ? AND checkin_date = ? AND status = 1", userChallengeID, today).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// 如果已经打卡成功，返回false
	if count > 0 {
		return false, nil
	}

	// 检查当前时间是否在打卡时间范围内
	currentTime := now.Format("15:04:05")

	return currentTime >= startTime && currentTime <= endTime, nil
}
