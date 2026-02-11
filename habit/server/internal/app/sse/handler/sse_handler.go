package handler

import (
	"habit/internal/app/sse/dto"
	"habit/internal/app/sse/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type SSEHandler struct {
	svc *service.SSEService
}

func NewSSEHandler(svc *service.SSEService) *SSEHandler {
	return &SSEHandler{svc: svc}
}

// Connect 建立 SSE 连接（GET）
func (h *SSEHandler) Connect(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return response.Error(c, response.CodeUnauthorized, "unauthorized")
	}

	return h.svc.Stream(c, userID)
}

// Send 推送消息（POST）- 本地调用发送
func (h *SSEHandler) Send(c *fiber.Ctx) error {
	var req dto.SSESendRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if req.Event == "" || req.Data == "" {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	// 单播或广播
	if req.UserID != nil {
		h.svc.SendToUser(*req.UserID, req.Event, req.Data)
	} else {
		h.svc.Broadcast(req.Event, req.Data)
	}

	return response.SuccessWithMessage(c, "message sent", nil)
}
