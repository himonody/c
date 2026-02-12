package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"habit/config"
	"habit/internal/app/challenge/service"
	"habit/internal/repo"
	"habit/pkg/database"
	"habit/pkg/logger"

	"go.uber.org/zap"
)

var (
	configFile = flag.String("config", "config/config.yaml", "config file path")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Logger); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// 初始化数据库
	if err := database.InitMySQL(&cfg.Database, logger.Logger); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.CloseMySQL()

	// 初始化Redis
	if err := database.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer database.CloseRedis()

	// 初始化repositories
	challengeRepo := repo.NewChallengeRepository(database.DB)
	userChallengeRepo := repo.NewUserChallengeRepository(database.DB)
	checkinRepo := repo.NewUserChallengeCheckinRepository(database.DB)
	settlementRepo := repo.NewUserChallengeSettlementRepository(database.DB)
	poolDailyRepo := repo.NewChallengePoolRepository(database.DB)
	dailyStatRepo := repo.NewChallengeStatRepository(database.DB)
	totalStatRepo := repo.NewChallengeTotalStatRepository(database.DB)

	// 初始化服务
	dailySettlementService := service.NewDailySettlementService(
		challengeRepo,
		userChallengeRepo,
		checkinRepo,
		settlementRepo,
		poolDailyRepo,
		dailyStatRepo,
		totalStatRepo,
		logger.Logger,
	)

	// 初始化调度器
	scheduler := service.NewSettlementScheduler(dailySettlementService, logger.Logger)

	// 创建context用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动调度器
	go func() {
		logger.Info("Starting settlement scheduler...")
		scheduler.Start(ctx)
	}()

	// 等待信号
	<-sigCh
	logger.Info("Received shutdown signal, stopping scheduler...")

	// 优雅关闭
	cancel()
	
	// 等待一段时间让调度器完成当前任务
	time.Sleep(5 * time.Second)
	logger.Info("Settlement scheduler stopped")
}
