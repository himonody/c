package handler

import (
	"habit/internal/app/challenge/dto"
	"habit/internal/app/challenge/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ChallengeStatHandler struct {
	svc       *service.ChallengeStatService
	actionSvc *service.ChallengeActionService
}

func NewChallengeStatHandler(svc *service.ChallengeStatService, actionSvc *service.ChallengeActionService) *ChallengeStatHandler {
	return &ChallengeStatHandler{
		svc:       svc,
		actionSvc: actionSvc,
	}
}

// TotalStat 获取平台累计统计
func (h *ChallengeStatHandler) TotalStat(c *fiber.Ctx) error {
	resp, err := h.svc.GetTotalStat()
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}
	return response.Success(c, resp)
}

// Start 开始挑战
func (h *ChallengeStatHandler) Start(c *fiber.Ctx) error {
	// 从上下文获取用户ID（假设有认证中间件）
	userID := c.Locals("userID").(int64)

	var req dto.ChallengeStartRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "invalid request body")
	}

	resp, err := h.actionSvc.Start(userID, &req)
	if err != nil {
		return response.Error(c, response.CodeBadRequest, err.Error())
	}

	return response.Success(c, resp)
}

// Money 增加挑战金
func (h *ChallengeStatHandler) Money(c *fiber.Ctx) error {
	// 从上下文获取用户ID（假设有认证中间件）
	userID := c.Locals("userID").(int64)

	var req dto.ChallengeMoneyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "invalid request body")
	}

	resp, err := h.actionSvc.Money(userID, &req)
	if err != nil {
		return response.Error(c, response.CodeBadRequest, err.Error())
	}

	return response.Success(c, resp)
}

// Query 查询用户今天和明天的挑战记录
func (h *ChallengeStatHandler) Query(c *fiber.Ctx) error {
	// 从上下文获取用户ID（假设有认证中间件）
	userID := c.Locals("userID").(int64)

	resp, err := h.actionSvc.Query(userID)
	if err != nil {
		return response.Error(c, response.CodeInternalError, err.Error())
	}

	return response.Success(c, resp)
}

// Checkin 打卡
func (h *ChallengeStatHandler) Checkin(c *fiber.Ctx) error {
	// 从上下文获取用户ID（假设有认证中间件）
	userID := c.Locals("userID").(int64)
	
	var req dto.CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "invalid request body")
	}

	resp, err := h.actionSvc.Checkin(userID, &req)
	if err != nil {
		return response.Error(c, response.CodeBadRequest, err.Error())
	}

	return response.Success(c, resp)
}
