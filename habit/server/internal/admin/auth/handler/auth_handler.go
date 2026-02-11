package handler

import (
	"habit/internal/admin/auth/dto"
	"habit/internal/admin/auth/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminAuthHandler struct {
	adminAuthService *service.AdminAuthService
}

func NewAdminAuthHandler(adminAuthService *service.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminAuthService: adminAuthService,
	}
}

// Login godoc
// @Summary Admin login
// @Description Admin login with username and password
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.AdminLoginRequest true "Login credentials"
// @Success 200 {object} dto.AdminLoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/auth/login [post]
func (h *AdminAuthHandler) Login(c *fiber.Ctx) error {
	var req dto.AdminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// Get client IP
	ip := c.IP()

	// Call service
	resp, err := h.adminAuthService.Login(&req, ip)
	if err != nil {
		if err.Error() == "invalid username or password" {
			return response.Error(c, response.CodeUnauthorized, "invalid_credentials")
		}
		if err.Error() == "account is disabled" {
			return response.Error(c, response.CodeForbidden, "account_disabled")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "login_success", resp)
}

// Logout godoc
// @Summary Admin logout
// @Description Logout current admin
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.AdminLogoutResponse
// @Failure 401 {object} map[string]interface{}
// @Router /admin/auth/logout [post]
func (h *AdminAuthHandler) Logout(c *fiber.Ctx) error {
	// Get admin ID and token from context (set by auth middleware)
	adminID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	token, ok := c.Locals("token").(string)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// Call service
	if err := h.adminAuthService.Logout(adminID, token); err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "logout_success", dto.AdminLogoutResponse{
		Message: "logout successful",
	})
}

// GetAdminInfo godoc
// @Summary Get current admin info
// @Description Get current logged in admin information
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.AdminUserInfo
// @Failure 401 {object} map[string]interface{}
// @Router /admin/auth/me [post]
func (h *AdminAuthHandler) GetAdminInfo(c *fiber.Ctx) error {
	// Get admin ID from context (set by auth middleware)
	adminID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	return response.Success(c, fiber.Map{
		"admin_id": adminID,
	})
}
