package service

import (
	"errors"

	"habit/internal/app/challenge/dto"
	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/database"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserChallengeStatService struct {
	db             *gorm.DB
	userRepo       *repo.UserRepository
	challengeRepo  *repo.ChallengeRepository
	checkinRepo    *repo.UserChallengeCheckinRepository
	settlementRepo *repo.UserChallengeSettlementRepository
	logger         *zap.Logger
}

func NewUserChallengeStatService(
	userRepo *repo.UserRepository,
	challengeRepo *repo.ChallengeRepository,
	checkinRepo *repo.UserChallengeCheckinRepository,
	settlementRepo *repo.UserChallengeSettlementRepository,
	logger *zap.Logger,
) *UserChallengeStatService {
	return &UserChallengeStatService{
		db:             database.DB,
		userRepo:       userRepo,
		challengeRepo:  challengeRepo,
		checkinRepo:    checkinRepo,
		settlementRepo: settlementRepo,
		logger:         logger,
	}
}

// GetUserChallengeStats 获取用户挑战统计信息（分页）
func (s *UserChallengeStatService) GetUserChallengeStats(req *dto.UserChallengeStatRequest) (*dto.UserChallengeStatResponse, error) {
	// 参数验证和默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 优先从打卡表获取用户ID（按最新打卡时间排序）
	userIDs, total, err := s.checkinRepo.GetUserIDsByLatestCheckin(req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to get user IDs by latest checkin", zap.Error(err))
		return nil, errors.New("failed to get user IDs")
	}

	// 根据用户ID获取用户信息
	users, err := s.userRepo.FindByIDs(userIDs)
	if err != nil {
		s.logger.Error("Failed to get users by IDs", zap.Error(err))
		return nil, errors.New("failed to get users")
	}

	// 创建用户映射
	userMap := make(map[int64]*model.AppUser)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 获取每个用户的统计信息
	result := make([]*dto.UserChallengeStatInfo, 0, len(userIDs))

	for _, userID := range userIDs {
		stat, err := s.getUserStat(userID, userMap[userID])
		if err != nil {
			s.logger.Warn("Failed to get user stat", zap.Int64("user_id", userID), zap.Error(err))
			continue
		}
		result = append(result, stat)
	}

	return &dto.UserChallengeStatResponse{
		List:     result,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// getUserStat 获取单个用户的统计信息
func (s *UserChallengeStatService) getUserStat(userID int64, user *model.AppUser) (*dto.UserChallengeStatInfo, error) {
	// 获取最新打卡时间
	latestCheckin, err := s.checkinRepo.GetLatestByUserID(userID)
	var latestCheckinAt string
	if err == nil && latestCheckin != nil && latestCheckin.CheckinTime != nil {
		latestCheckinAt = latestCheckin.CheckinTime.Format("2006-01-02 15:04:05")
	}

	// 获取累计打卡天数
	totalCheckinDays, err := s.checkinRepo.CountDistinctCheckinDatesByUserID(userID)
	if err != nil {
		totalCheckinDays = 0
	}

	// 获取累计收益
	totalProfit, err := s.settlementRepo.SumFinalProfitByUserID(userID)
	if err != nil {
		totalProfit = decimal.Zero
	}

	// 获取用户名称
	userName := ""
	if user != nil {
		userName = user.Nickname
		if userName == "" {
			userName = user.Username
		}
	}

	return &dto.UserChallengeStatInfo{
		UserID:           userID,
		UserName:         userName,
		LatestCheckinAt:  latestCheckinAt,
		TotalCheckinDays: totalCheckinDays,
		TotalProfit:      totalProfit.StringFixed(2),
	}, nil
}
