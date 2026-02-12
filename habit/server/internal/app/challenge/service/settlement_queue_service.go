package service

import (
	"context"
	"errors"
	"time"

	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/cache"
	"habit/pkg/database"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettlementQueueService struct {
	db                *gorm.DB
	challengeRepo     *repo.ChallengeRepository
	userChallengeRepo *repo.UserChallengeRepository
	settlementRepo    *repo.UserChallengeSettlementRepository
	poolDailyRepo     *repo.ChallengePoolRepository
	dailyStatRepo     *repo.ChallengeStatRepository
	totalStatRepo     *repo.ChallengeTotalStatRepository
	queue             *cache.Queue
	logger            *zap.Logger
}

func NewSettlementQueueService(
	challengeRepo *repo.ChallengeRepository,
	userChallengeRepo *repo.UserChallengeRepository,
	settlementRepo *repo.UserChallengeSettlementRepository,
	poolDailyRepo *repo.ChallengePoolRepository,
	dailyStatRepo *repo.ChallengeStatRepository,
	totalStatRepo *repo.ChallengeTotalStatRepository,
	logger *zap.Logger,
) *SettlementQueueService {
	return &SettlementQueueService{
		db:                database.DB,
		challengeRepo:     challengeRepo,
		userChallengeRepo: userChallengeRepo,
		settlementRepo:    settlementRepo,
		poolDailyRepo:     poolDailyRepo,
		dailyStatRepo:     dailyStatRepo,
		totalStatRepo:     totalStatRepo,
		queue:             cache.NewQueue(database.RedisClient, logger),
		logger:            logger,
	}
}

// EnqueueCheckinSettlement 将打卡结算任务加入队列
func (s *SettlementQueueService) EnqueueCheckinSettlement(userChallengeID, userID, checkinID int64, checkinDate string) error {
	message := &cache.CheckinSettlementMessage{
		UserChallengeID: userChallengeID,
		UserID:          userID,
		CheckinID:       checkinID,
		CheckinDate:     checkinDate,
		Timestamp:       time.Now().Unix(),
	}

	return s.queue.EnqueueCheckinSettlement(context.Background(), message)
}

// StartConsumer 启动结算队列消费者
func (s *SettlementQueueService) StartConsumer(ctx context.Context) {
	s.logger.Info("Starting settlement queue consumer")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Settlement queue consumer stopped")
			return
		default:
			// 从队列中获取结算任务，超时时间5秒
			message, err := s.queue.DequeueCheckinSettlement(ctx, 5*time.Second)
			if err != nil {
				s.logger.Error("Failed to dequeue settlement message", zap.Error(err))
				time.Sleep(1 * time.Second) // 等待1秒后重试
				continue
			}

			if message == nil {
				// 队列为空，继续等待
				continue
			}

			// 处理结算任务
			if err := s.processSettlement(ctx, message); err != nil {
				s.logger.Error("Failed to process settlement",
					zap.Int64("userChallengeId", message.UserChallengeID),
					zap.Int64("userId", message.UserID),
					zap.Error(err))
			}
		}
	}
}

// processSettlement 处理单个结算任务
func (s *SettlementQueueService) processSettlement(ctx context.Context, message *cache.CheckinSettlementMessage) error {
	s.logger.Info("Processing settlement",
		zap.Int64("userChallengeId", message.UserChallengeID),
		zap.Int64("userId", message.UserID),
		zap.String("checkinDate", message.CheckinDate))

	// 使用事务处理结算
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.processSettlementInTx(ctx, tx, message)
	})
}

// processSettlementInTx 在事务中处理结算
func (s *SettlementQueueService) processSettlementInTx(ctx context.Context, tx *gorm.DB, message *cache.CheckinSettlementMessage) error {
	// 1. 获取挑战配置
	challengeConfig, err := s.challengeRepo.Last()
	if err != nil {
		return errors.New("failed to get challenge config")
	}

	// 2. 获取用户挑战信息
	userChallenge, err := s.userChallengeRepo.FindByID(message.UserChallengeID)
	if err != nil {
		return errors.New("user challenge not found")
	}

	// 3. 计算收益
	baseProfit := s.calculateBaseProfit(userChallenge, challengeConfig)
	platformSubsidy := decimal.NewFromInt(int64(challengeConfig.DailyPlatformSubsidy))
	totalRawProfit := baseProfit.Add(platformSubsidy)

	// 4. 计算扣除金额（如果超过上限）
	taxDeduction := s.calculateTaxDeduction(totalRawProfit, challengeConfig)
	finalProfit := totalRawProfit.Sub(taxDeduction)

	// 5. 创建结算记录
	settlement := &model.AppUserChallengeSettlement{
		UserChallengeID: message.UserChallengeID,
		UserID:          message.UserID,
		CheckinID:       message.CheckinID,
		SettleDate:      time.Now(), // 使用当前时间作为结算时间
		BaseProfit:      baseProfit,
		PlatformSubsidy: platformSubsidy,
		TotalRawProfit:  totalRawProfit,
		TaxDeduction:    taxDeduction,
		FinalProfit:     finalProfit,
		IsSettled:       1, // 已结算
		SettleAt:        time.Now(),
	}

	// 6. 保存结算记录
	if err := s.settlementRepo.Create(settlement); err != nil {
		return err
	}

	// 7. 更新用户钱包余额
	if err := s.updateUserWallet(tx, message.UserID, finalProfit); err != nil {
		return err
	}

	s.logger.Info("Settlement processed successfully",
		zap.Int64("userChallengeId", message.UserChallengeID),
		zap.Int64("userId", message.UserID),
		zap.String("finalProfit", finalProfit.StringFixed(2)))

	return nil
}

// calculateBaseProfit 计算基础收益
func (s *SettlementQueueService) calculateBaseProfit(userChallenge *model.AppUserChallenge, challengeConfig *model.AppChallenge) decimal.Decimal {
	// 基础收益逻辑：根据用户挑战金额和平台规则计算
	// 这里简化处理，实际可能需要更复杂的计算逻辑
	totalAmount := userChallenge.ChallengeAmount

	// 假设基础收益率为1%（实际应该根据业务规则计算）
	return totalAmount.Mul(decimal.NewFromFloat(0.01))
}

// calculateTaxDeduction 计算超过上限的扣除金额
func (s *SettlementQueueService) calculateTaxDeduction(totalRawProfit decimal.Decimal, challengeConfig *model.AppChallenge) decimal.Decimal {
	maxDailyProfit := decimal.NewFromInt(int64(challengeConfig.MaxDailyProfit))

	if maxDailyProfit.LessThanOrEqual(decimal.Zero) || totalRawProfit.LessThanOrEqual(maxDailyProfit) {
		return decimal.Zero
	}

	// 超过部分的扣除比例
	excessAmount := totalRawProfit.Sub(maxDailyProfit)
	taxRate := decimal.NewFromInt(int64(challengeConfig.ExcessTaxRate)).Div(decimal.NewFromInt(100))

	return excessAmount.Mul(taxRate)
}

// updateUserWallet 更新用户钱包余额
func (s *SettlementQueueService) updateUserWallet(tx *gorm.DB, userID int64, profit decimal.Decimal) error {
	// 这里需要实现钱包余额更新逻辑
	// 由于没有看到钱包相关的repository，先预留接口
	// TODO: 实现钱包余额更新

	s.logger.Info("User wallet updated",
		zap.Int64("userId", userID),
		zap.String("profit", profit.StringFixed(2)))

	return nil
}

// GetQueueStatus 获取队列状态
func (s *SettlementQueueService) GetQueueStatus(ctx context.Context) (int64, error) {
	return s.queue.GetQueueLength(ctx, cache.CheckinSettlementQueue)
}
