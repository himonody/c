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
	log, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.CloseDB()

	// 初始化Redis
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer database.CloseRedis()

	// 初始化repositories
	challengeRepo := repo.NewChallengeRepository(database.DB)
	userChallengeRepo := repo.NewUserChallengeRepository(database.DB)
	settlementRepo := repo.NewUserChallengeSettlementRepository(database.DB)
	poolDailyRepo := repo.NewChallengePoolRepository(database.DB)
	dailyStatRepo := repo.NewChallengeStatRepository(database.DB)
	totalStatRepo := repo.NewChallengeTotalStatRepository(database.DB)

	// 初始化结算队列服务
	settlementQueueService := service.NewSettlementQueueService(
		challengeRepo,
		userChallengeRepo,
		settlementRepo,
		poolDailyRepo,
		dailyStatRepo,
		totalStatRepo,
		log,
	)

	// 创建context用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动消费者
	go func() {
		log.Info("Starting settlement queue consumer...")
		settlementQueueService.StartConsumer(ctx)
	}()

	// 等待信号
	<-sigCh
	log.Info("Received shutdown signal, stopping consumer...")

	// 优雅关闭
	cancel()
	
	// 等待一段时间让消费者完成当前任务
	time.Sleep(5 * time.Second)
	log.Info("Settlement consumer stopped")
}
