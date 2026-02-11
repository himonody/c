package repo

import (
	"time"

	"habit/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LeaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

func (r *LeaderboardRepository) ListDaily(rankType int8, rankDate time.Time, page, pageSize int) ([]*model.AppChallengeRankDaily, int64, error) {
	var list []*model.AppChallengeRankDaily
	var total int64

	query := r.db.Model(&model.AppChallengeRankDaily{}).
		Where("rank_type = ?", rankType).
		Where("rank_date = ?", rankDate.Format("2006-01-02"))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("rank_no ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

type LeaderboardTotalRow struct {
	UserID int64           `gorm:"column:user_id"`
	Value  decimal.Decimal `gorm:"column:value"`
}

func (r *LeaderboardRepository) ListTotalSum(rankType int8, page, pageSize int) ([]*LeaderboardTotalRow, int64, error) {
	var list []*LeaderboardTotalRow
	var total int64

	base := r.db.Model(&model.AppChallengeRankDaily{}).Where("rank_type = ?", rankType)
	if err := base.Distinct("user_id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.AppChallengeRankDaily{}).
		Select("user_id, SUM(value) AS value").
		Where("rank_type = ?", rankType).
		Group("user_id").
		Order("value DESC")

	if err := query.Offset(offset).Limit(pageSize).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *LeaderboardRepository) ListLatest(rankType int8, page, pageSize int) ([]*model.AppChallengeRankDaily, int64, error) {
	var list []*model.AppChallengeRankDaily
	var total int64

	maxDateSub := r.db.Model(&model.AppChallengeRankDaily{}).
		Select("MAX(rank_date)").
		Where("rank_type = ?", rankType)

	query := r.db.Model(&model.AppChallengeRankDaily{}).
		Where("rank_type = ?", rankType).
		Where("rank_date = (?)", maxDateSub)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("rank_no ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
