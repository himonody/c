package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"habit/internal/admin/challenge/dto"
	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/cache"
	"habit/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChallengeService struct {
	db     *gorm.DB
	repo   *repo.ChallengeRepository
	logger *zap.Logger
	rdb    *cache.Cache
}

func NewChallengeService(challengeRepo *repo.ChallengeRepository, logger *zap.Logger) *ChallengeService {
	return &ChallengeService{
		db:     database.DB,
		repo:   challengeRepo,
		logger: logger,
		rdb:    cache.NewCache(database.RedisClient),
	}
}

func (s *ChallengeService) List(req *dto.ChallengeListRequest) (*dto.ChallengeListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	rows, total, err := s.repo.List(req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list challenge", zap.Error(err))
		return nil, errors.New("failed to list challenge")
	}

	if len(rows) > 0 {
		err = s.rdb.Set(context.Background(), cache.AppChallengeKey(), rows[0], time.Duration(24)*time.Hour)
		if err != nil {
			s.logger.Error("Failed to list challenge", zap.Error(err))
		}
	}

	list := make([]*dto.ChallengeInfo, 0, len(rows))
	for _, c := range rows {
		info := &dto.ChallengeInfo{
			ID:                   c.ID,
			IsAutoSettle:         c.IsAutoSettle,
			SettleTime:           c.SettleTime,
			CycleDays:            c.CycleDays,
			StartTime:            c.StartTime,
			EndTime:              c.EndTime,
			MaxDepositAmount:     c.MaxDepositAmount,
			MinWithdrawAmount:    c.MinWithdrawAmount,
			MaxDailyProfit:       c.MaxDailyProfit,
			ExcessTaxRate:        c.ExcessTaxRate,
			MinDailyProfit:       c.MinDailyProfit,
			DailyPlatformSubsidy: c.DailyPlatformSubsidy,
			UncheckDeductRate:    c.UncheckDeductRate,
			MinUncheckUsers:      c.MinUncheckUsers,
			CommissionFollow:     c.CommissionFollow,
			CommissionJoin:       c.CommissionJoin,
			CommissionL1:         c.CommissionL1,
			CommissionL2:         c.CommissionL2,
			CommissionL3:         c.CommissionL3,
		}
		if c.UpdatedAt != nil {
			info.UpdatedAt = c.UpdatedAt.Format(time.DateTime)
		}

		list = append(list, info)
	}

	return &dto.ChallengeListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *ChallengeService) Create(req *dto.ChallengeUpsertRequest) error {

	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:admin:challenge:create:%d", req.CycleDays)
	lock, locked, err := cache.AcquireLockWithRenew(ctx, database.RedisClient, lockKey, 5*time.Second)
	if err != nil {
		s.logger.Error("Failed to acquire challenge create lock", zap.Error(err))
		return errors.New("internal server error")
	}
	if !locked {
		return errors.New("duplicate request")
	}
	defer func() {
		_ = lock.Unlock(ctx)
	}()

	return s.db.Transaction(func(tx *gorm.DB) error {
		challenge := &model.AppChallenge{
			IsAutoSettle:         req.IsAutoSettle,
			SettleTime:           req.SettleTime,
			CycleDays:            req.CycleDays,
			StartTime:            req.StartTime,
			EndTime:              req.EndTime,
			MaxDepositAmount:     req.MaxDepositAmount,
			MinWithdrawAmount:    req.MinWithdrawAmount,
			MaxDailyProfit:       req.MaxDailyProfit,
			ExcessTaxRate:        req.ExcessTaxRate,
			MinDailyProfit:       req.MinDailyProfit,
			DailyPlatformSubsidy: req.DailyPlatformSubsidy,
			UncheckDeductRate:    req.UncheckDeductRate,
			MinUncheckUsers:      req.MinUncheckUsers,
			CommissionFollow:     req.CommissionFollow,
			CommissionJoin:       req.CommissionJoin,
			CommissionL1:         req.CommissionL1,
			CommissionL2:         req.CommissionL2,
			CommissionL3:         req.CommissionL3,
		}
		if err := s.repo.Create(tx, challenge); err != nil {
			return err
		}
		return nil
	})
}

func (s *ChallengeService) Update(req *dto.ChallengeUpsertRequest) error {
	if req.ID <= 0 {
		return errors.New("bad_request")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			_ = s.rdb.Delete(context.Background(), cache.AppChallengeKey())
		}()
		challenge, err := s.repo.FindByID(req.ID)
		if err != nil {
			return err
		}
		challenge.IsAutoSettle = req.IsAutoSettle
		challenge.SettleTime = req.SettleTime
		challenge.CycleDays = req.CycleDays
		challenge.StartTime = req.StartTime
		challenge.EndTime = req.EndTime
		challenge.MaxDepositAmount = req.MaxDepositAmount
		challenge.MinWithdrawAmount = req.MinWithdrawAmount
		challenge.MaxDailyProfit = req.MaxDailyProfit
		challenge.ExcessTaxRate = req.ExcessTaxRate
		challenge.MinDailyProfit = req.MinDailyProfit
		challenge.DailyPlatformSubsidy = req.DailyPlatformSubsidy
		challenge.UncheckDeductRate = req.UncheckDeductRate
		challenge.MinUncheckUsers = req.MinUncheckUsers
		challenge.CommissionFollow = req.CommissionFollow
		challenge.CommissionJoin = req.CommissionJoin
		challenge.CommissionL1 = req.CommissionL1
		challenge.CommissionL2 = req.CommissionL2
		challenge.CommissionL3 = req.CommissionL3
		if err := s.repo.Update(tx, challenge); err != nil {
			return err
		}
		return nil
	})
}
