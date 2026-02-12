package repo

import (
	"errors"
	"habit/internal/model"

	"gorm.io/gorm"
)

type ChallengeTotalStatRepository struct {
	db *gorm.DB
}

func NewChallengeTotalStatRepository(db *gorm.DB) *ChallengeTotalStatRepository {
	return &ChallengeTotalStatRepository{db: db}
}

// GetOrCreate 获取或创建累计统计记录
func (r *ChallengeTotalStatRepository) GetOrCreate() (*model.AppChallengeTotalStat, error) {
	var stat model.AppChallengeTotalStat

	// 尝试获取记录
	err := r.db.Where("id = ?", 1).First(&stat).Error
	if err == nil {
		return &stat, nil
	}
	
	// 如果记录不存在，创建默认记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stat = model.AppChallengeTotalStat{
			ID: 1,
		}
		if err := r.db.Create(&stat).Error; err != nil {
			return nil, err
		}
		return &stat, nil
	}

	return nil, err
}

// Update 更新累计统计
func (r *ChallengeTotalStatRepository) Update(stat *model.AppChallengeTotalStat) error {
	return r.db.Save(stat).Error
}

// IncrementJoinCount 增加参与人次
func (r *ChallengeTotalStatRepository) IncrementJoinCount() error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_join_cnt", gorm.Expr("total_join_cnt + ?", 1)).Error
}

// IncrementSuccessCount 增加成功人次
func (r *ChallengeTotalStatRepository) IncrementSuccessCount() error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_success_cnt", gorm.Expr("total_success_cnt + ?", 1)).Error
}

// IncrementFailCount 增加失败人次
func (r *ChallengeTotalStatRepository) IncrementFailCount() error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_fail_cnt", gorm.Expr("total_fail_cnt + ?", 1)).Error
}

// AddJoinAmount 增加参与金额
func (r *ChallengeTotalStatRepository) AddJoinAmount(amount string) error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_join_amount", gorm.Expr("total_join_amount + ?", amount)).Error
}

// AddSuccessAmount 增加成功金额
func (r *ChallengeTotalStatRepository) AddSuccessAmount(amount string) error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_success_amount", gorm.Expr("total_success_amount + ?", amount)).Error
}

// AddFailAmount 增加失败金额
func (r *ChallengeTotalStatRepository) AddFailAmount(amount string) error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_fail_amount", gorm.Expr("total_fail_amount + ?", amount)).Error
}

// AddPlatformBonus 增加平台补贴
func (r *ChallengeTotalStatRepository) AddPlatformBonus(amount string) error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_platform_bonus", gorm.Expr("total_platform_bonus + ?", amount)).Error
}

// AddPoolAmount 增加奖池金额
func (r *ChallengeTotalStatRepository) AddPoolAmount(amount string) error {
	return r.db.Model(&model.AppChallengeTotalStat{}).
		Where("id = ?", 1).
		UpdateColumn("total_pool_amount", gorm.Expr("total_pool_amount + ?", amount)).Error
}
