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
	if err = logger.InitLogger(&cfg.Logger); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// 初始化数据库
	if err = database.InitMySQL(&cfg.Database, logger.Logger); err != nil {
		logger.Logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.CloseMySQL()

	// 初始化Redis
	if err := database.InitRedis(&cfg.Redis); err != nil {
		logger.Logger.Fatal("Failed to initialize Redis", zap.Error(err))
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

	// 创建context用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动定时队列消费者
	startScheduledQueueConsumer(ctx, logger.Logger, challengeRepo, userChallengeRepo, checkinRepo, settlementRepo, poolDailyRepo, dailyStatRepo, totalStatRepo)

	// 等待信号
	<-sigCh
	logger.Logger.Info("Received shutdown signal, stopping services...")

	// 优雅关闭
	cancel()

	// 等待一段时间让服务完成当前任务
	time.Sleep(5 * time.Second)
	logger.Logger.Info("Settlement services stopped")
}

// startScheduledQueueConsumer 启动定时队列消费者
func startScheduledQueueConsumer(ctx context.Context, log *zap.Logger, challengeRepo *repo.ChallengeRepository, userChallengeRepo *repo.UserChallengeRepository, checkinRepo *repo.UserChallengeCheckinRepository, settlementRepo *repo.UserChallengeSettlementRepository, poolDailyRepo *repo.ChallengePoolRepository, dailyStatRepo *repo.ChallengeStatRepository, totalStatRepo *repo.ChallengeTotalStatRepository) {
	log.Info("Starting scheduled queue consumer...")

	// 创建定时结算服务
	dailySettlementService := service.NewDailySettlementService(
		challengeRepo,
		userChallengeRepo,
		checkinRepo,
		settlementRepo,
		poolDailyRepo,
		dailyStatRepo,
		totalStatRepo,
		log,
	)

	// 创建队列消费者服务
	settlementQueueService := service.NewSettlementQueueService(
		challengeRepo,
		userChallengeRepo,
		settlementRepo,
		poolDailyRepo,
		dailyStatRepo,
		totalStatRepo,
		log,
	)

	// 创建调度器
	scheduler := service.NewSettlementScheduler(dailySettlementService, log)
	
	// 启动定时器和队列消费者
	go scheduler.Start(ctx)
	go settlementQueueService.StartConsumer(ctx)
}
