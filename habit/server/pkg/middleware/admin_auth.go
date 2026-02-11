package middleware

import (
	"strings"

	"habit/internal/admin/auth/service"
	"habit/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AdminAuthMiddleware validates admin JWT token and handles auto-renewal
func AdminAuthMiddleware(adminAuthService *service.AdminAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token and get claims
		adminID, claims, err := adminAuthService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Check if token should be refreshed
		if utils.ShouldRefreshToken(claims) {
			newToken, refreshed, err := adminAuthService.RefreshToken(adminID, claims.Username, token)
			if err == nil && refreshed {
				// Set new token in response header
				c.Set("X-New-Token", newToken)
			}
		}

		// Store admin ID and username in context
		c.Locals("userID", adminID)
		c.Locals("username", claims.Username)
		c.Locals("token", token)

		return c.Next()
	}
}
