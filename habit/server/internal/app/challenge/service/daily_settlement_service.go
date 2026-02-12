package service

import (

	"errors"
	"time"

	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/database"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DailySettlementService struct {
	db                   *gorm.DB
	challengeRepo        *repo.ChallengeRepository
	userChallengeRepo    *repo.UserChallengeRepository
	checkinRepo          *repo.UserChallengeCheckinRepository
	settlementRepo       *repo.UserChallengeSettlementRepository
	poolDailyRepo        *repo.ChallengePoolRepository
	dailyStatRepo        *repo.ChallengeStatRepository
	totalStatRepo        *repo.ChallengeTotalStatRepository
	logger               *zap.Logger
}

func NewDailySettlementService(
	challengeRepo *repo.ChallengeRepository,
	userChallengeRepo *repo.UserChallengeRepository,
	checkinRepo *repo.UserChallengeCheckinRepository,
	settlementRepo *repo.UserChallengeSettlementRepository,
	poolDailyRepo *repo.ChallengePoolRepository,
	dailyStatRepo *repo.ChallengeStatRepository,
	totalStatRepo *repo.ChallengeTotalStatRepository,
	logger *zap.Logger,
) *DailySettlementService {
	return &DailySettlementService{
		db:                   database.DB,
		challengeRepo:        challengeRepo,
		userChallengeRepo:    userChallengeRepo,
		checkinRepo:          checkinRepo,
		settlementRepo:       settlementRepo,
		poolDailyRepo:        poolDailyRepo,
		dailyStatRepo:        dailyStatRepo,
		totalStatRepo:        totalStatRepo,
		logger:               logger,
	}
}

// ExecuteDailySettlement 执行每日结算
func (s *DailySettlementService) ExecuteDailySettlement(settleDate time.Time) error {
	s.logger.Info("Starting daily settlement", zap.Time("settleDate", settleDate))

	// 获取挑战配置
	challengeConfig, err := s.challengeRepo.Last()
	if err != nil {
		return errors.New("failed to get challenge config")
	}

	// 使用事务处理整个结算过程
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.executeDailySettlementInTx(tx, settleDate, challengeConfig)
	})
}

// executeDailySettlementInTx 在事务中执行每日结算
func (s *DailySettlementService) executeDailySettlementInTx(tx *gorm.DB, settleDate time.Time, challengeConfig *model.AppChallenge) error {
	// 1. 获取当天所有成功打卡的用户
	successfulCheckins, err := s.getSuccessfulCheckins(tx, settleDate)
	if err != nil {
		return err
	}

	// 2. 计算奖池金额
	poolAmount, err := s.calculatePoolAmount(tx, settleDate, challengeConfig, successfulCheckins)
	if err != nil {
		return err
	}

	// 3. 计算平均收益
	avgProfit := decimal.Zero
	if len(successfulCheckins) > 0 {
		avgProfit = poolAmount.Div(decimal.NewFromInt(int64(len(successfulCheckins))))
	}

	// 4. 为每个成功打卡的用户创建结算记录
	for _, checkin := range successfulCheckins {
		if err := s.createUserSettlement(tx, checkin, challengeConfig, avgProfit, settleDate); err != nil {
			s.logger.Error("Failed to create user settlement", 
				zap.Int64("checkinId", checkin.ID),
				zap.Int64("userId", checkin.UserID),
				zap.Error(err))
			continue
		}
	}

	// 5. 创建奖池日结算记录
	if err := s.createPoolDailyRecord(tx, settleDate, poolAmount, successfulCheckins, challengeConfig); err != nil {
		return err
	}

	// 6. 创建/更新每日统计记录
	if err := s.createOrUpdateDailyStat(tx, settleDate, successfulCheckins, poolAmount); err != nil {
		return err
	}

	// 7. 更新累计统计
	if err := s.updateTotalStat(successfulCheckins, poolAmount); err != nil {
		return err
	}

	s.logger.Info("Daily settlement completed successfully",
		zap.Time("settleDate", settleDate),
		zap.Int("successCount", len(successfulCheckins)),
		zap.String("poolAmount", poolAmount.StringFixed(2)),
		zap.String("avgProfit", avgProfit.StringFixed(2)))

	return nil
}

// getSuccessfulCheckins 获取当天所有成功打卡的用户
func (s *DailySettlementService) getSuccessfulCheckins(tx *gorm.DB, settleDate time.Time) ([]*model.AppUserChallengeCheckin, error) {
	var checkins []*model.AppUserChallengeCheckin
	
	err := tx.Where("DATE(checkin_date) = ? AND status = ?", settleDate.Format("2006-01-02"), 1).
		Find(&checkins).Error
	
	if err != nil {
		return nil, err
	}

	return checkins, nil
}

// calculatePoolAmount 计算奖池金额
func (s *DailySettlementService) calculatePoolAmount(tx *gorm.DB, settleDate time.Time, challengeConfig *model.AppChallenge, successfulCheckins []*model.AppUserChallengeCheckin) (decimal.Decimal, error) {
	// 1. 计算未打卡用户扣除金额
	failDeductPool := decimal.Zero
	
	// 获取当天有挑战但未打卡的用户
	var failedChallenges []*model.AppUserChallenge
	err := tx.Where("start_date <= ? AND end_date >= ? AND status = ?", 
		settleDate, settleDate, 1).
		Find(&failedChallenges).Error
	
	if err != nil {
		return decimal.Zero, err
	}

	// 计算未打卡扣除
	successfulUserIDs := make(map[int64]bool)
	for _, checkin := range successfulCheckins {
		successfulUserIDs[checkin.UserID] = true
	}

	for _, challenge := range failedChallenges {
		if !successfulUserIDs[challenge.UserID] {
			// 该用户有挑战但未打卡，需要扣除
			totalAmount := challenge.ChallengeAmount.Add(challenge.PreRecharge)
			deductRate := decimal.NewFromInt(int64(challengeConfig.UncheckDeductRate)).Div(decimal.NewFromInt(100))
			deductAmount := totalAmount.Mul(deductRate)
			failDeductPool = failDeductPool.Add(deductAmount)
		}
	}

	// 2. 平台补贴
	platformSubsidyPool := decimal.NewFromInt(int64(challengeConfig.DailyPlatformSubsidy)).Mul(decimal.NewFromInt(int64(len(successfulCheckins))))

	// 3. 总可分配金额
	totalDistributable := failDeductPool.Add(platformSubsidyPool)

	return totalDistributable, nil
}

// createUserSettlement 为用户创建结算记录
func (s *DailySettlementService) createUserSettlement(tx *gorm.DB, checkin *model.AppUserChallengeCheckin, challengeConfig *model.AppChallenge, avgProfit decimal.Decimal, settleDate time.Time) error {
	// 获取用户挑战信息
	userChallenge, err := s.userChallengeRepo.FindByID(checkin.UserChallengeID)
	if err != nil {
		return err
	}

	// 计算基础收益
	baseProfit := s.calculateBaseProfit(userChallenge, challengeConfig)
	
	// 平台补贴
	platformSubsidy := decimal.NewFromInt(int64(challengeConfig.DailyPlatformSubsidy))
	
	// 总原始收益
	totalRawProfit := baseProfit.Add(platformSubsidy)
	
	// 计算超额扣除
	taxDeduction := s.calculateTaxDeduction(avgProfit, challengeConfig)
	
	// 最终收益
	finalProfit := avgProfit.Sub(taxDeduction)

	// 创建结算记录
	settlement := &model.AppUserChallengeSettlement{
		UserChallengeID: checkin.UserChallengeID,
		UserID:          checkin.UserID,
		CheckinID:       checkin.ID,
		SettleDate:      settleDate,
		BaseProfit:      baseProfit,
		PlatformSubsidy: platformSubsidy,
		TotalRawProfit:  totalRawProfit,
		TaxDeduction:    taxDeduction,
		FinalProfit:     finalProfit,
		IsSettled:       1, // 已结算
		SettleAt:        time.Now(),
	}

	return s.settlementRepo.Create(settlement)
}

// createPoolDailyRecord 创建奖池日结算记录
func (s *DailySettlementService) createPoolDailyRecord(tx *gorm.DB, settleDate time.Time, poolAmount decimal.Decimal, successfulCheckins []*model.AppUserChallengeCheckin, challengeConfig *model.AppChallenge) error {
	// 计算失败用户数
	failUsers := 0
	
	// 获取当天活跃挑战总数
	var totalActiveChallenges int64
	tx.Model(&model.AppUserChallenge{}).
		Where("start_date <= ? AND end_date >= ? AND status = ?", settleDate, settleDate, 1).
		Count(&totalActiveChallenges)
	
	failUsers = int(totalActiveChallenges) - len(successfulCheckins)

	// 计算系统回收金额
	systemTaxIncome := decimal.Zero
	
	poolDaily := &model.AppChallengePoolDaily{
		PoolDate:              settleDate,
		TotalUsers:            int(totalActiveChallenges),
		SuccessUsers:          len(successfulCheckins),
		FailUsers:             failUsers,
		FailDeductPool:        s.calculateFailDeductPool(failUsers, challengeConfig),
		PlatformSubsidyPool:   decimal.NewFromInt(int64(challengeConfig.DailyPlatformSubsidy)).Mul(decimal.NewFromInt(int64(len(successfulCheckins)))),
		TotalDistributable:    poolAmount,
		AvgProfit:             poolAmount.Div(decimal.NewFromInt(int64(len(successfulCheckins)))),
		SystemTaxIncome:       systemTaxIncome,
		CreatedAt:             time.Now(),
	}

	return s.poolDailyRepo.Create(tx, poolDaily)
}

// createOrUpdateDailyStat 创建或更新每日统计
func (s *DailySettlementService) createOrUpdateDailyStat(tx *gorm.DB, settleDate time.Time, successfulCheckins []*model.AppUserChallengeCheckin, poolAmount decimal.Decimal) error {
	// 计算各项统计数据
	joinUserCnt := len(successfulCheckins)
	successUserCnt := len(successfulCheckins)
	failUserCnt := 0
	
	// 计算金额统计
	joinAmount := decimal.Zero
	successAmount := decimal.Zero
	failAmount := decimal.Zero
	
	for _, checkin := range successfulCheckins {
		userChallenge, err := s.userChallengeRepo.FindByID(checkin.UserChallengeID)
		if err != nil {
			continue
		}
		
		totalAmount := userChallenge.ChallengeAmount.Add(userChallenge.PreRecharge)
		joinAmount = joinAmount.Add(totalAmount)
		successAmount = successAmount.Add(totalAmount)
	}

	platformBonus := decimal.NewFromInt(int64(0)) // 根据实际业务逻辑计算

	// 查找是否已存在统计记录
	var existingStat model.AppChallengeDailyStat
	err := tx.Where("stat_date = ?", settleDate).First(&existingStat).Error
	
	if err == nil {
		// 更新现有记录
		existingStat.JoinUserCnt = joinUserCnt
		existingStat.SuccessUserCnt = successUserCnt
		existingStat.FailUserCnt = failUserCnt
		existingStat.JoinAmount = joinAmount
		existingStat.SuccessAmount = successAmount
		existingStat.FailAmount = failAmount
		existingStat.PlatformBonus = platformBonus
		existingStat.PoolAmount = poolAmount
		existingStat.UpdatedAt = time.Now()
		
		return tx.Save(&existingStat).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		dailyStat := &model.AppChallengeDailyStat{
			StatDate:      settleDate,
			JoinUserCnt:   joinUserCnt,
			SuccessUserCnt: successUserCnt,
			FailUserCnt:    failUserCnt,
			JoinAmount:    joinAmount,
			SuccessAmount: successAmount,
			FailAmount:    failAmount,
			PlatformBonus: platformBonus,
			PoolAmount:    poolAmount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		
		return tx.Create(dailyStat).Error
	}
	
	return err
}

// updateTotalStat 更新累计统计
func (s *DailySettlementService) updateTotalStat(successfulCheckins []*model.AppUserChallengeCheckin, poolAmount decimal.Decimal) error {
	totalStat, err := s.totalStatRepo.GetOrCreate()
	if err != nil {
		return err
	}

	// 更新累计数据
	totalStat.TotalJoinCnt += len(successfulCheckins)
	totalStat.TotalSuccessCnt += len(successfulCheckins)
	totalStat.TotalPoolAmount = totalStat.TotalPoolAmount.Add(poolAmount)

	return s.totalStatRepo.Update(totalStat)
}

// calculateBaseProfit 计算基础收益
func (s *DailySettlementService) calculateBaseProfit(userChallenge *model.AppUserChallenge, challengeConfig *model.AppChallenge) decimal.Decimal {
	totalAmount := userChallenge.ChallengeAmount.Add(userChallenge.PreRecharge)
	return totalAmount.Mul(decimal.NewFromFloat(0.01)) // 1%基础收益率
}

// calculateTaxDeduction 计算超额扣除
func (s *DailySettlementService) calculateTaxDeduction(avgProfit decimal.Decimal, challengeConfig *model.AppChallenge) decimal.Decimal {
	maxDailyProfit := decimal.NewFromInt(int64(challengeConfig.MaxDailyProfit))
	
	if maxDailyProfit.LessThanOrEqual(decimal.Zero) || avgProfit.LessThanOrEqual(maxDailyProfit) {
		return decimal.Zero
	}

	excessAmount := avgProfit.Sub(maxDailyProfit)
	taxRate := decimal.NewFromInt(int64(challengeConfig.ExcessTaxRate)).Div(decimal.NewFromInt(100))
	
	return excessAmount.Mul(taxRate)
}

// calculateFailDeductPool 计算未打卡扣除总金额
func (s *DailySettlementService) calculateFailDeductPool(failUsers int, challengeConfig *model.AppChallenge) decimal.Decimal {
	// 这里简化处理，实际应该根据每个未打卡用户的挑战金额计算
	// 可以返回一个估算值或0，具体根据业务需求
	return decimal.Zero
}
