package handler

import (
	"habit/internal/app/leaderboard/dto"
	"habit/internal/app/leaderboard/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type LeaderboardHandler struct {
	svc *service.LeaderboardService
}

func NewLeaderboardHandler(svc *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{svc: svc}
}

// List 排行榜查询
func (h *LeaderboardHandler) List(c *fiber.Ctx) error {
	var req dto.LeaderboardListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	resp, err := h.svc.List(&req)
	if err != nil {
		if err.Error() == "bad_request" {
			return response.Error(c, response.CodeBadRequest, "bad_request")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}


// Total 总排行榜查询（自动取最新一期 rankDate）
func (h *LeaderboardHandler) Total(c *fiber.Ctx) error {
	var req dto.LeaderboardTotalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	resp, err := h.svc.Total(&req)
	if err != nil {
		if err.Error() == "bad_request" {
			return response.Error(c, response.CodeBadRequest, "bad_request")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}
