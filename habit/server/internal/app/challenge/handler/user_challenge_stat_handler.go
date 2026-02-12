package handler

import (
	"habit/internal/app/challenge/dto"
	"habit/internal/app/challenge/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserChallengeStatHandler struct {
	svc *service.UserChallengeStatService
}

func NewUserChallengeStatHandler(svc *service.UserChallengeStatService) *UserChallengeStatHandler {
	return &UserChallengeStatHandler{svc: svc}
}

// GetUserChallengeStats 获取用户挑战统计信息
// POST /api/challenge/user-stats
func (h *UserChallengeStatHandler) GetUserChallengeStats(c *fiber.Ctx) error {
	var req dto.UserChallengeStatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "invalid request body")
	}

	resp, err := h.svc.GetUserChallengeStats(&req)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal error")
	}

	return response.Success(c, resp)
}
