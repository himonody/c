package router

import (
	appAuthHandler "habit/internal/app/auth/handler"
	appAuthService "habit/internal/app/auth/service"
	challengeHandler "habit/internal/app/challenge/handler"
	challengeService "habit/internal/app/challenge/service"
	inviteHandler "habit/internal/app/invite/handler"
	inviteService "habit/internal/app/invite/service"
	leaderboardHandler "habit/internal/app/leaderboard/handler"
	leaderboardService "habit/internal/app/leaderboard/service"
	sseHandler "habit/internal/app/sse/handler"
	sseService "habit/internal/app/sse/service"
	walletHandler "habit/internal/app/wallet/handler"
	walletService "habit/internal/app/wallet/service"
	withdrawHandler "habit/internal/app/withdraw/handler"
	withdrawService "habit/internal/app/withdraw/service"
	"habit/internal/repo"
	"habit/pkg/database"
	"habit/pkg/logger"
	"habit/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAppRoutes 设置 App 端路由
func SetupAppRoutes(api fiber.Router, appAuthHdl *appAuthHandler.AuthHandler, authService *appAuthService.AuthService) {
	// Initialize wallet dependencies
	walletRepo := repo.NewWalletRepository(database.DB)
	walletSvc := walletService.NewWalletService(walletRepo, logger.Logger)
	walletHdl := walletHandler.NewWalletHandler(walletSvc)

	// Initialize withdraw dependencies
	withdrawRepo := repo.NewWithdrawRepository(database.DB)
	balanceLogRepo := repo.NewBalanceLogRepository(database.DB)
	withdrawSvc := withdrawService.NewWithdrawService(withdrawRepo, walletRepo, balanceLogRepo, logger.Logger)
	withdrawHdl := withdrawHandler.NewWithdrawHandler(withdrawSvc)

	// Initialize challenge stat dependencies
	challengeRepo := repo.NewChallengeStatRepository(database.DB)
	challengeSvc := challengeService.NewChallengeStatService(challengeRepo, logger.Logger)
	challengeHdl := challengeHandler.NewChallengeStatHandler(challengeSvc)

	// Initialize leaderboard dependencies
	leaderboardRepo := repo.NewLeaderboardRepository(database.DB)
	leaderboardSvc := leaderboardService.NewLeaderboardService(leaderboardRepo, logger.Logger)
	leaderboardHdl := leaderboardHandler.NewLeaderboardHandler(leaderboardSvc)

	// Initialize invite dependencies
	inviteRepo := repo.NewInviteRepository(database.DB)
	inviteSvc := inviteService.NewInviteService(inviteRepo, logger.Logger)
	inviteHdl := inviteHandler.NewInviteHandler(inviteSvc)

	// Initialize SSE dependencies
	userRepo := repo.NewUserRepository(database.DB)
	sseSvc := sseService.NewSSEService(logger.Logger, userRepo)
	sseHdl := sseHandler.NewSSEHandler(sseSvc)

	// App 路由分组
	app := api.Group("/app")

	// 公开路由
	app.Post("/auth/register", appAuthHdl.Register)
	app.Post("/auth/login", appAuthHdl.Login)

	// 需要认证的路由
	appProtected := app.Group("")
	appProtected.Use(middleware.AuthMiddleware(authService))
	appProtected.Post("/auth/logout", appAuthHdl.Logout)
	appProtected.Post("/auth/me", appAuthHdl.GetUserInfo)
	appProtected.Post("/auth/change-password", appAuthHdl.ChangePassword)
	appProtected.Post("/auth/set-pay-password", appAuthHdl.SetPayPassword)
	appProtected.Post("/auth/update-profile", appAuthHdl.UpdateProfile)

	// 钱包路由
	appProtected.Post("/wallet/info", walletHdl.GetWalletInfo)
	appProtected.Post("/wallet/address", walletHdl.SetWalletAddress)

	// 提款路由
	appProtected.Post("/withdraw/list", withdrawHdl.List)
	appProtected.Post("/withdraw/apply", withdrawHdl.Apply)

	// 挑战路由
	appProtected.Post("/challenge/total-stat", challengeHdl.TotalStat)
	appProtected.Post("/challenge/start", challengeHdl.Start)
	appProtected.Post("/challenge/money", challengeHdl.Start)

	// 排行榜路由
	appProtected.Post("/leaderboard/list", leaderboardHdl.List)
	appProtected.Post("/leaderboard/total", leaderboardHdl.Total)

	// 邀请路由
	appProtected.Post("/invite/info", inviteHdl.Info)
	appProtected.Post("/invite/friends", inviteHdl.Friends)

	// SSE 路由
	appProtected.Get("/sse/connect", sseHdl.Connect)
	appProtected.Post("/sse/send", sseHdl.Send)
}
