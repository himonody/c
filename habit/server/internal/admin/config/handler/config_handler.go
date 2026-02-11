package handler

import (
	"habit/internal/admin/config/dto"
	"habit/internal/admin/config/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler(configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// List 获取配置列表
func (h *ConfigHandler) List(c *fiber.Ctx) error {
	var req dto.ConfigListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	resp, err := h.configService.GetConfigList(&req)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}

// Get 获取配置详情
func (h *ConfigHandler) Get(c *fiber.Ctx) error {
	var req dto.GetConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if req.ID <= 0 {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	resp, err := h.configService.GetConfig(req.ID)
	if err != nil {
		if err.Error() == "config not found" {
			return response.Error(c, response.CodeNotFound, "not_found")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}

// Create 创建配置
func (h *ConfigHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 基本验证
	if req.ConfigName == "" || req.ConfigKey == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 获取管理员ID
	adminID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	err := h.configService.CreateConfig(&req, adminID)
	if err != nil {
		if err.Error() == "config key already exists" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "配置键已存在")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}

// Update 更新配置
func (h *ConfigHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 基本验证
	if req.ID <= 0 || req.ConfigName == "" || req.ConfigKey == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 获取管理员ID
	adminID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	err := h.configService.UpdateConfig(&req, adminID)
	if err != nil {
		if err.Error() == "config not found" {
			return response.Error(c, response.CodeNotFound, "not_found")
		}
		if err.Error() == "config key already exists" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "配置键已存在")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}
