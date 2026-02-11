package handler

import (
	"habit/internal/app/wallet/dto"
	"habit/internal/app/wallet/service"
	"habit/pkg/response"
	"habit/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetWalletInfo 获取钱包信息
func (h *WalletHandler) GetWalletInfo(c *fiber.Ctx) error {
	// 获取用户ID
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	// 调用 service
	walletInfo, err := h.walletService.GetWalletInfo(userID)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, walletInfo)
}

func (h *WalletHandler) SetWalletAddress(c *fiber.Ctx) error {
	var req dto.SetWalletAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	address := strings.TrimSpace(req.Address)
	if address == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}
	if !utils.IsValidTRC20Address(address) {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	if err := h.walletService.SetWalletAddress(userID, address); err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}
