package handler

import (
	"habit/internal/app/withdraw/dto"
	"habit/internal/app/withdraw/service"
	"habit/pkg/response"
	"habit/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type WithdrawHandler struct {
	withdrawService *service.WithdrawService
}

func NewWithdrawHandler(withdrawService *service.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{withdrawService: withdrawService}
}

func (h *WithdrawHandler) List(c *fiber.Ctx) error {
	var req dto.WithdrawListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}
	resp, err := h.withdrawService.GetWithdrawList(userID, &req)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}

func (h *WithdrawHandler) Apply(c *fiber.Ctx) error {
	var req dto.WithdrawApplyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if strings.TrimSpace(req.PayPassword) == "" {
		return response.ErrorWithMessage(c, response.CodeBadRequest, "请设置支付密码")
	}

	address := strings.TrimSpace(req.Address)
	if address == "" {
		return response.ErrorWithMessage(c, response.CodeBadRequest, "请设置提款地址")
	}
	if !utils.IsValidTRC20Address(address) {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	amountStr := strings.TrimSpace(req.Amount)
	amount, err := decimal.NewFromString(amountStr)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	bizID, err := h.withdrawService.ApplyWithdraw(userID, amount, address, c.IP(), req.PayPassword)
	if err != nil {
		if err.Error() == "duplicate request" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "请勿重复提交")
		}
		if err.Error() == "insufficient balance" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "余额不足")
		}
		if err.Error() == "pay password not set" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "请设置支付密码")
		}
		if err.Error() == "invalid pay password" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "支付密码错误")
		}
		if err.Error() == "amount too small" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "金额过小")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, &dto.WithdrawApplyResponse{BizID: bizID})
}
