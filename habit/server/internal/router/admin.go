package router

import (
	adminAuthHandler "habit/internal/admin/auth/handler"
	adminAuthService "habit/internal/admin/auth/service"
	challengeHandler "habit/internal/admin/challenge/handler"
	challengeService "habit/internal/admin/challenge/service"
	configHandler "habit/internal/admin/config/handler"
	configService "habit/internal/admin/config/service"
	"habit/internal/repo"
	"habit/pkg/database"
	"habit/pkg/logger"
	"habit/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes 设置 Admin 端路由
func SetupAdminRoutes(api fiber.Router, adminHandler *adminAuthHandler.AdminAuthHandler, adminAuthService *adminAuthService.AdminAuthService) {
	// Initialize config dependencies
	configRepo := repo.NewConfigRepository(database.DB)
	configSvc := configService.NewConfigService(configRepo, logger.Logger)
	configHdl := configHandler.NewConfigHandler(configSvc)

	// Initialize challenge dependencies
	challengeRepo := repo.NewChallengeRepository(database.DB)
	challengePoolRepo := repo.NewChallengePoolRepository(database.DB)
	challengeSvc := challengeService.NewChallengeService(challengeRepo, challengePoolRepo, logger.Logger)
	challengeHdl := challengeHandler.NewChallengeHandler(challengeSvc)

	// Admin 路由分组
	admin := api.Group("/admin")

	// 公开路由
	admin.Post("/auth/login", adminHandler.Login)

	// 需要认证的路由
	adminProtected := admin.Group("")
	adminProtected.Use(middleware.AdminAuthMiddleware(adminAuthService))
	adminProtected.Post("/auth/logout", adminHandler.Logout)
	adminProtected.Post("/auth/me", adminHandler.GetAdminInfo)

	// 系统配置路由
	adminProtected.Post("/config/list", configHdl.List)
	adminProtected.Post("/config/get", configHdl.Get)
	adminProtected.Post("/config/create", configHdl.Create)
	adminProtected.Post("/config/update", configHdl.Update)

	// 挑战配置路由
	adminProtected.Post("/challenge/list", challengeHdl.List)
	adminProtected.Post("/challenge/create", challengeHdl.Create)
	adminProtected.Post("/challenge/update", challengeHdl.Update)
}
