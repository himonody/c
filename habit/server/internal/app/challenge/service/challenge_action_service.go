package service

import (
	"errors"
	"time"

	"habit/internal/app/challenge/dto"
	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/database"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChallengeActionService struct {
	db                       *gorm.DB
	userRepo                 *repo.UserRepository
	challengeRepo            *repo.ChallengeRepository
	userChallengeRepo        *repo.UserChallengeRepository
	userChallengeCheckinRepo *repo.UserChallengeCheckinRepository
	settlementQueueService   *SettlementQueueService
	logger                   *zap.Logger
}

func NewChallengeActionService(
	userRepo *repo.UserRepository,
	challengeRepo *repo.ChallengeRepository,
	userChallengeRepo *repo.UserChallengeRepository,
	userChallengeCheckinRepo *repo.UserChallengeCheckinRepository,
	settlementQueueService *SettlementQueueService,
	logger *zap.Logger,
) *ChallengeActionService {
	return &ChallengeActionService{
		db:                       database.DB,
		userRepo:                 userRepo,
		challengeRepo:            challengeRepo,
		userChallengeRepo:        userChallengeRepo,
		userChallengeCheckinRepo: userChallengeCheckinRepo,
		settlementQueueService:   settlementQueueService,
		logger:                   logger,
	}
}

// Start 开始挑战
func (s *ChallengeActionService) Start(userID int64, req *dto.ChallengeStartRequest) (*dto.ChallengeStartResponse, error) {
	// 验证预充值金额
	preRecharge, err := decimal.NewFromString(req.PreRecharge)
	if err != nil {
		return nil, errors.New("invalid pre recharge amount")
	}
	if preRecharge.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("pre recharge amount must be greater than 0")
	}

	// 检查用户是否已有活跃挑战
	existingChallenge, err := s.userChallengeRepo.GetActiveChallengeByUserID(userID)
	if err == nil && existingChallenge != nil {
		return nil, errors.New("user already has an active challenge")
	}

	// 获取当前有效的挑战配置
	challengeConfig, err := s.challengeRepo.Last()
	if err != nil {
		s.logger.Error("Failed to get challenge config", zap.Error(err))
		return nil, errors.New("failed to get challenge config")
	}

	// 检查预充值金额是否超过上限
	if challengeConfig.MaxDepositAmount > 0 {
		maxAmount := decimal.NewFromInt(int64(challengeConfig.MaxDepositAmount))
		if preRecharge.GreaterThan(maxAmount) {
			return nil, errors.New("pre recharge amount exceeds maximum limit")
		}
	}

	var userChallenge *model.AppUserChallenge

	// 使用事务处理
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 计算挑战开始和结束日期
		now := time.Now()

		// 根据打卡时间计算开始日期
		var startDate time.Time
		currentTime := now.Format("15:04:05")
		checkinStartTime := challengeConfig.StartTime

		// 如果当前时间在打卡开始时间之前，今天开始；否则明天开始
		if currentTime < checkinStartTime {
			// 6:00之前创建，今天开始
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		} else {
			// 6:00之后创建，明天开始
			startDate = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		}

		endDate := startDate.AddDate(0, 0, challengeConfig.CycleDays)

		// 创建用户挑战记录
		userChallenge = &model.AppUserChallenge{
			UserID:          userID,
			ChallengeID:     int64(challengeConfig.ID),
			PoolID:          0,            // 暂时设为0，后续可能需要奖池逻辑
			ChallengeAmount: decimal.Zero, // 预充值作为主要挑战金额
			PreRecharge:     preRecharge,
			StartDate:       startDate,
			EndDate:         endDate,
			Status:          1, // 进行中
			FailReason:      2, // 默认失败原因：未打卡
			CreatedAt:       now,
		}

		if err := s.userChallengeRepo.Create(userChallenge); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to start challenge", zap.Error(err))
		return nil, errors.New("failed to start challenge")
	}

	return &dto.ChallengeStartResponse{
		UserChallengeID: userChallenge.ID,
		ChallengeID:     userChallenge.ChallengeID,
		ChallengeAmount: userChallenge.ChallengeAmount.StringFixed(2),
		PreRecharge:     userChallenge.PreRecharge.StringFixed(2),
		StartDate:       userChallenge.StartDate.Format("2006-01-02"),
		EndDate:         userChallenge.EndDate.Format("2006-01-02"),
		Status:          userChallenge.Status,
	}, nil
}

// Money 增加挑战金
func (s *ChallengeActionService) Money(userID int64, req *dto.ChallengeMoneyRequest) (*dto.ChallengeMoneyResponse, error) {
	// 验证增加金额
	addAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errors.New("invalid amount")
	}
	if addAmount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than 0")
	}
	// 获取更新后的记录用于响应
	userChallenge, err := s.userChallengeRepo.FindByID(req.UserChallengeID)
	if err != nil {
		s.logger.Error("Failed to get updated challenge", zap.Error(err))
		return nil, errors.New("failed to get updated challenge")
	}
	// 获取当前有效的挑战配置
	challengeConfig, err := s.challengeRepo.Last()
	if err != nil {
		s.logger.Error("Failed to get challenge config", zap.Error(err))
		return nil, errors.New("failed to get challenge config")
	}

	// 检查预充值金额是否超过上限
	if challengeConfig.MaxDepositAmount > 0 {
		maxAmount := decimal.NewFromInt(int64(challengeConfig.MaxDepositAmount))
		if addAmount.Add(userChallenge.ChallengeAmount).Add(userChallenge.PreRecharge).GreaterThan(maxAmount) {
			return nil, errors.New("pre recharge amount exceeds maximum limit")
		}
	}
	// 直接更新预充值金额，使用WHERE条件确保权限和状态
	err = s.db.Model(&model.AppUserChallenge{}).
		Where("id = ? AND user_id = ? AND status = ?", req.UserChallengeID, userID, 1).
		Update("pre_recharge", gorm.Expr("pre_recharge + ?", addAmount)).Error

	if err != nil {
		s.logger.Error("Failed to update pre recharge amount", zap.Error(err))
		return nil, errors.New("failed to update pre recharge amount")
	}

	// 计算更新前后的金额（这里简化处理，实际可能需要从日志或缓存获取）
	return &dto.ChallengeMoneyResponse{
		UserChallengeID: userChallenge.ID,
		ChallengeID:     userChallenge.ChallengeID,
		ChallengeAmount: userChallenge.ChallengeAmount.StringFixed(2),
		PreRecharge:     userChallenge.PreRecharge.Add(addAmount).StringFixed(2),
		StartDate:       userChallenge.StartDate.Format("2006-01-02"),
		EndDate:         userChallenge.EndDate.Format("2006-01-02"),
		Status:          userChallenge.Status,
	}, nil
}

// Query 查询用户今天和明天的挑战记录
func (s *ChallengeActionService) Query(userID int64) (*dto.ChallengeQueryResponse, error) {
	// 一次性获取今天和明天的挑战记录
	challenges, err := s.userChallengeRepo.GetTodayAndTomorrowChallenges(userID)
	if err != nil {
		s.logger.Error("Failed to get challenges", zap.Error(err))
		return nil, errors.New("failed to get challenges")
	}

	// 分离今天和明天的挑战记录
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	var todayChallenge, tomorrowChallenge *model.AppUserChallenge

	for _, challenge := range challenges {
		challengeDate := challenge.StartDate.Format("2006-01-02")
		if challengeDate == today {
			todayChallenge = challenge
		} else if challengeDate == tomorrow {
			tomorrowChallenge = challenge
		}
	}

	// 转换为响应格式
	response := &dto.ChallengeQueryResponse{
		TodayChallenge:    s.convertToChallengeInfo(todayChallenge),
		TomorrowChallenge: s.convertToChallengeInfo(tomorrowChallenge),
	}

	return response, nil
}

// convertToChallengeInfo 将模型转换为DTO
func (s *ChallengeActionService) convertToChallengeInfo(challenge *model.AppUserChallenge) *dto.UserChallengeInfo {
	if challenge == nil {
		return nil
	}

	info := &dto.UserChallengeInfo{
		ID:              challenge.ID,
		UserID:          challenge.UserID,
		UserChallengeID: challenge.ID,
		ChallengeID:     challenge.ChallengeID,
		ChallengeAmount: challenge.ChallengeAmount.StringFixed(2),
		PreRecharge:     challenge.PreRecharge.StringFixed(2),
		StartDate:       challenge.StartDate.Format("2006-01-02"),
		EndDate:         challenge.EndDate.Format("2006-01-02"),
		Status:          challenge.Status,
		FailReason:      challenge.FailReason,
		CreatedAt:       challenge.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 设置状态描述
	switch challenge.Status {
	case 1:
		info.StatusText = "进行中"
	case 2:
		info.StatusText = "成功"
	case 3:
		info.StatusText = "失败"
	default:
		info.StatusText = "未知"
	}

	// 设置完成时间
	if challenge.FinishedAt != nil {
		finishedAtStr := challenge.FinishedAt.Format("2006-01-02 15:04:05")
		info.FinishedAt = &finishedAtStr
	}

	return info
}

// Checkin 打卡
func (s *ChallengeActionService) Checkin(userID int64, req *dto.CheckinRequest) (*dto.CheckinResponse, error) {
	// 获取用户挑战记录
	userChallenge, err := s.userChallengeRepo.FindByID(req.UserChallengeID)
	if err != nil {
		s.logger.Error("Failed to get user challenge", zap.Error(err))
		return nil, errors.New("user challenge not found")
	}

	// 验证挑战是否属于当前用户
	if userChallenge.UserID != userID {
		return nil, errors.New("unauthorized operation")
	}

	// 验证挑战状态是否为进行中
	if userChallenge.Status != 1 {
		return nil, errors.New("challenge is not active")
	}

	// 获取挑战配置
	challengeConfig, err := s.challengeRepo.Last()
	if err != nil {
		s.logger.Error("Failed to get challenge config", zap.Error(err))
		return nil, errors.New("failed to get challenge config")
	}

	// 检查是否可以打卡
	canCheckin, err := s.userChallengeCheckinRepo.CheckCanCheckin(
		req.UserChallengeID,
		challengeConfig.StartTime,
		challengeConfig.EndTime,
	)
	if err != nil {
		s.logger.Error("Failed to check can checkin", zap.Error(err))
		return nil, errors.New("failed to check checkin availability")
	}

	if !canCheckin {
		return &dto.CheckinResponse{
			UserChallengeID: req.UserChallengeID,
			CheckinTime:     time.Now().Format("2006-01-02 15:04:05"),
			CheckinDate:     time.Now().Format("2006-01-02"),
			Status:          2, // 失败
			Message:         "不在打卡时间范围内或今日已打卡",
			ChallengeInfo:   s.convertToChallengeInfo(userChallenge),
		}, nil
	}

	// 创建打卡记录
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	checkin := &model.AppUserChallengeCheckin{
		UserChallengeID: req.UserChallengeID,
		UserID:          userID,
		CheckinDate:     today,
		CheckinTime:     &now,
		MoodCode:        req.MoodCode,
		MoodText:        req.MoodText,
		ContentType:     1, // 默认图片类型
		Status:          1, // 打卡成功
		CreatedAt:       now,
	}

	err = s.userChallengeCheckinRepo.Create(checkin)
	if err != nil {
		s.logger.Error("Failed to create checkin record", zap.Error(err))
		return nil, errors.New("failed to create checkin record")
	}

	// 将结算任务加入Redis队列
	if s.settlementQueueService != nil {
		if err := s.settlementQueueService.EnqueueCheckinSettlement(
			req.UserChallengeID,
			userID,
			checkin.ID,
			today.Format("2006-01-02"),
		); err != nil {
			// 队列加入失败不影响打卡成功，只记录日志
			s.logger.Error("Failed to enqueue settlement task", zap.Error(err))
		}
	}

	return &dto.CheckinResponse{
		UserChallengeID: req.UserChallengeID,
		CheckinTime:     now.Format("2006-01-02 15:04:05"),
		CheckinDate:     today.Format("2006-01-02"),
		Status:          1, // 成功
		Message:         "打卡成功",
		ChallengeInfo:   s.convertToChallengeInfo(userChallenge),
	}, nil
}
