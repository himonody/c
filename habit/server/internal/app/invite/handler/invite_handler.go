package handler

import (
	"habit/internal/app/invite/dto"
	"habit/internal/app/invite/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type InviteHandler struct {
	svc *service.InviteService
}

func NewInviteHandler(svc *service.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

// Info 获取邀请信息（friend_code 和 invite_url）
func (h *InviteHandler) Info(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	resp, err := h.svc.GetInviteInfo(userID)
	if err != nil {
		if err.Error() == "user_not_found" {
			return response.Error(c, response.CodeNotFound, "user_not_found")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}

// Friends 获取邀请的好友列表（avatar、nickname、username）
func (h *InviteHandler) Friends(c *fiber.Ctx) error {
	var req dto.InviteFriendsRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	resp, err := h.svc.ListInvitedFriends(int(userID), &req)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}
