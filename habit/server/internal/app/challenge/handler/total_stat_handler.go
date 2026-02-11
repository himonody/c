package handler

import (
	"habit/internal/app/challenge/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ChallengeStatHandler struct {
	svc *service.ChallengeStatService
}

func NewChallengeStatHandler(svc *service.ChallengeStatService) *ChallengeStatHandler {
	return &ChallengeStatHandler{svc: svc}
}

// TotalStat 获取平台累计统计
func (h *ChallengeStatHandler) TotalStat(c *fiber.Ctx) error {
	resp, err := h.svc.GetTotalStat()
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}
	return response.Success(c, resp)
}
func (h *ChallengeStatHandler) Start(c *fiber.Ctx) error {

}

func (h *ChallengeStatHandler) Money(c *fiber.Ctx) error {

}
