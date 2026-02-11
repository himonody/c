package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"habit/config"
	"habit/internal/router"
	"habit/pkg/database"
	"habit/pkg/logger"
	"habit/pkg/middleware"
	"habit/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(&cfg.Logger); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application...")

	// Initialize JWT
	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg.JWT.RefreshThresholdHours)
	logger.Info("JWT initialized",
		zap.Int("expire_hours", cfg.JWT.ExpireHours),
		zap.Int("refresh_threshold_hours", cfg.JWT.RefreshThresholdHours),
	)

	// Initialize MySQL
	if err := database.InitMySQL(&cfg.Database, logger.Logger); err != nil {
		logger.Fatal("Failed to initialize MySQL", zap.Error(err))
	}
	defer database.CloseMySQL()

	// Initialize Redis
	if err := database.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer database.CloseRedis()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:         "Habit Tracker",
		ErrorHandler:    customErrorHandler,
		ReadBufferSize:  20 << 20,
		WriteBufferSize: 20 << 20,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	// Routes
	router.SetupRoutes(app)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server starting", zap.String("address", addr))

	// Graceful shutdown
	go func() {
		if err := app.Listen(addr); err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("Request error",
		zap.String("path", c.Path()),
		zap.Int("status", code),
		zap.Error(err),
	)

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
