package middleware

import (
	"time"

	"habit/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger middleware logs HTTP requests
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request
		logger.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("ip", c.IP()),
			zap.Duration("latency", time.Since(start)),
		)

		return err
	}
}
