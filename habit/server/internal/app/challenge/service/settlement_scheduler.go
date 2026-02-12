package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type SettlementScheduler struct {
	dailySettlementService *DailySettlementService
	logger                 *zap.Logger
}

func NewSettlementScheduler(dailySettlementService *DailySettlementService, logger *zap.Logger) *SettlementScheduler {
	return &SettlementScheduler{
		dailySettlementService: dailySettlementService,
		logger:                 logger,
	}
}

// Start 启动定时结算调度器
func (s *SettlementScheduler) Start(ctx context.Context) {
	s.logger.Info("Starting settlement scheduler")

	// 立即执行一次结算检查
	go s.runDailySettlement(ctx)

	// 每分钟检查一次是否到了结算时间
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Settlement scheduler stopped")
			return
		case <-ticker.C:
			go s.runDailySettlement(ctx)
		}
	}
}

// runDailySettlement 执行每日结算检查
func (s *SettlementScheduler) runDailySettlement(ctx context.Context) {
	now := time.Now()
	
	// 检查是否到了结算时间 (06:10:00)
	if now.Hour() == 6 && now.Minute() == 10 && now.Second() < 30 {
		s.logger.Info("Triggering daily settlement", zap.Time("now", now))
		
		// 执行昨天的结算
		yesterday := now.AddDate(0, 0, -1)
		if err := s.dailySettlementService.ExecuteDailySettlement(yesterday); err != nil {
			s.logger.Error("Failed to execute daily settlement", 
				zap.Time("settleDate", yesterday),
				zap.Error(err))
		} else {
			s.logger.Info("Daily settlement completed successfully", 
				zap.Time("settleDate", yesterday))
		}
	}
}

// StartWithCron 使用cron表达式启动调度器 (可选实现)
func (s *SettlementScheduler) StartWithCron(ctx context.Context, cronExpr string) {
	s.logger.Info("Starting settlement scheduler with cron", zap.String("cron", cronExpr))
	
	// 这里可以集成cron库，比如 github.com/robfig/cron/v3
	// 为了简化，暂时使用简单的定时检查
	s.Start(ctx)
}
