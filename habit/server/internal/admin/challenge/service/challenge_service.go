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

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChallengeService struct {
	db       *gorm.DB
	repo     *repo.ChallengeRepository
	poolRepo *repo.ChallengePoolRepository
	logger   *zap.Logger
}

func NewChallengeService(challengeRepo *repo.ChallengeRepository, poolRepo *repo.ChallengePoolRepository, logger *zap.Logger) *ChallengeService {
	return &ChallengeService{
		db:       database.DB,
		repo:     challengeRepo,
		poolRepo: poolRepo,
		logger:   logger,
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

	rows, total, err := s.repo.ListWithLatestPool(req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list challenge", zap.Error(err))
		return nil, errors.New("failed to list challenge")
	}

	list := make([]*dto.ChallengeInfo, 0, len(rows))
	for _, row := range rows {
		c := row.AppChallenge
		info := &dto.ChallengeInfo{
			ID:            c.ID,
			DayCount:      c.DayCount,
			Amount:        c.Amount.String(),
			CheckinStart:  c.CheckinStart,
			CheckinEnd:    c.CheckinEnd,
			PlatformBonus: c.PlatformBonus.String(),
			Status:        c.Status,
			Sort:          c.Sort,
			CreatedAt:     c.CreatedAt.Format(time.DateTime),
		}
		info.Pool = &dto.ChallengePoolInfo{
			TotalAmount: row.PoolTotalAmount.String(),
			Settled:     row.PoolSettled,
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
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return errors.New("bad_request")
	}
	platformBonus, err := decimal.NewFromString(req.PlatformBonus)
	if err != nil {
		return errors.New("bad_request")
	}

	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:admin:challenge:create:%d:%s", req.DayCount, amount.String())
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

	var startDatePtr *time.Time
	if req.StartDate != "" {
		t, err := time.Parse(time.DateTime, req.StartDate)
		if err != nil {
			return errors.New("bad_request")
		}
		startDatePtr = &t
	}
	var endDatePtr *time.Time
	if req.EndDate != "" {
		t, err := time.Parse(time.DateTime, req.EndDate)
		if err != nil {
			return errors.New("bad_request")
		}
		endDatePtr = &t
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		challenge := &model.AppChallenge{
			DayCount:      req.DayCount,
			Amount:        amount,
			CheckinStart:  req.CheckinStart,
			CheckinEnd:    req.CheckinEnd,
			PlatformBonus: platformBonus,
			Status:        req.Status,
			Sort:          req.Sort,
		}
		if err := s.repo.Create(tx, challenge); err != nil {
			return err
		}

		pool := &model.AppChallengePool{
			ChallengeID: challenge.ID,
			StartDate:   startDatePtr,
			EndDate:     endDatePtr,
			TotalAmount: platformBonus,
		}
		if err := s.poolRepo.Create(tx, pool); err != nil {
			return err
		}
		return nil
	})
}

func (s *ChallengeService) Update(req *dto.ChallengeUpsertRequest) error {
	if req.ID <= 0 {
		return errors.New("bad_request")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return errors.New("bad_request")
	}
	platformBonus, err := decimal.NewFromString(req.PlatformBonus)
	if err != nil {
		return errors.New("bad_request")
	}

	var startDatePtr *time.Time
	if req.StartDate != "" {
		t, err := time.Parse(time.DateTime, req.StartDate)
		if err != nil {
			return errors.New("bad_request")
		}
		startDatePtr = &t
	}
	var endDatePtr *time.Time
	if req.EndDate != "" {
		t, err := time.Parse(time.DateTime, req.EndDate)
		if err != nil {
			return errors.New("bad_request")
		}
		endDatePtr = &t
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		challenge, err := s.repo.FindByID(req.ID)
		if err != nil {
			return err
		}
		challenge.DayCount = req.DayCount
		challenge.Amount = amount
		challenge.CheckinStart = req.CheckinStart
		challenge.CheckinEnd = req.CheckinEnd
		challenge.PlatformBonus = platformBonus
		challenge.Status = req.Status
		challenge.Sort = req.Sort
		if err := s.repo.Update(tx, challenge); err != nil {
			return err
		}

		pool := &model.AppChallengePool{
			ID:          req.PoolID,
			ChallengeID: req.ID,
			StartDate:   startDatePtr,
			EndDate:     endDatePtr,
			TotalAmount: challenge.PlatformBonus.Sub(platformBonus),
			Settled:     0,
		}
		if req.PoolID > 0 {
			if err := s.poolRepo.Update(tx, pool); err != nil {
				return err
			}
			return nil
		}
		if err := s.poolRepo.Create(tx, pool); err != nil {
			return err
		}
		return nil
	})
}
