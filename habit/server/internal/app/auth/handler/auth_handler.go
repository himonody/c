package handler

import (
	"habit/internal/app/auth/dto"
	"habit/internal/app/auth/service"
	"habit/pkg/response"
	"habit/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration info"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// Validate username format (6-12 characters, alphanumeric and special chars)
	if errMsg := utils.GetUsernameValidationError(req.Username); errMsg != "" {
		return response.Error(c, response.CodeValidationError, "username_invalid")
	}

	// Validate password format (6-12 characters, alphanumeric and special chars)
	if errMsg := utils.GetPasswordValidationError(req.Password); errMsg != "" {
		return response.Error(c, response.CodeValidationError, "password_invalid")
	}

	// Get client IP
	ip := c.IP()

	// Call service
	resp, err := h.authService.Register(&req, ip)
	if err != nil {
		if err.Error() == "username already exists" {
			return response.Error(c, response.CodeBadRequest, "username_exists")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "register_success", resp)
}

// Login godoc
// @Summary User login
// @Description Login with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
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
	resp, err := h.authService.Login(&req, ip)
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
// @Summary User logout
// @Description Logout current user
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.LogoutResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user ID and token from context (set by auth middleware)
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	token, ok := c.Locals("token").(string)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// Call service
	if err := h.authService.Logout(userID, token); err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "logout_success", dto.LogoutResponse{
		Message: "logout successful",
	})
}

// GetUserInfo godoc
// @Summary Get current user info
// @Description Get current logged in user information
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.UserInfo
// @Failure 401 {object} map[string]interface{}
// @Router /auth/me [get]
func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// 调用 service 获取用户信息（带缓存）
	userInfo, err := h.authService.GetUserInfo(userID)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, userInfo)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 基本验证
	if req.OldPassword == "" || req.NewPassword == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 获取用户ID
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// 调用 service
	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		if err.Error() == "old password is incorrect" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "旧密码错误")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}

// SetPayPassword 设置/修改支付密码
func (h *AuthHandler) SetPayPassword(c *fiber.Ctx) error {
	var req dto.SetPayPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 基本验证
	if req.PayPassword == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 获取用户ID
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// 调用 service
	err := h.authService.SetPayPassword(userID, &req)
	if err != nil {
		if err.Error() == "old pay password is required" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "需要提供旧支付密码")
		}
		if err.Error() == "old pay password is incorrect" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "旧支付密码错误")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}

// UpdateProfile 修改用户资料
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 获取用户ID
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// 调用 service
	err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		if err.Error() == "no fields to update" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "没有需要更新的字段")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}
