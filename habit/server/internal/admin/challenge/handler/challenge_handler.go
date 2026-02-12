package handler

import (
	"habit/internal/admin/challenge/dto"
	"habit/internal/admin/challenge/service"
	"habit/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ChallengeHandler struct {
	svc *service.ChallengeService
}

func NewChallengeHandler(svc *service.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{svc: svc}
}

// List 挑战列表
func (h *ChallengeHandler) List(c *fiber.Ctx) error {
	var req dto.ChallengeListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	resp, err := h.svc.List(&req)
	if err != nil {
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.Success(c, resp)
}

// Create 新增挑战
func (h *ChallengeHandler) Create(c *fiber.Ctx) error {
	var req dto.ChallengeUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if req.CycleDays <= 0 {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if err := h.svc.Create(&req); err != nil {
		if err.Error() == "bad_request" {
			return response.Error(c, response.CodeBadRequest, "bad_request")
		}
		if err.Error() == "duplicate request" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "请勿重复提交")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}

// Update 编辑挑战
func (h *ChallengeHandler) Update(c *fiber.Ctx) error {
	var req dto.ChallengeUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if req.ID <= 0 || req.CycleDays <= 0 {
		return response.Error(c, response.CodeBadRequest, "bad_request")
	}

	if err := h.svc.Update(&req); err != nil {
		if err.Error() == "bad_request" {
			return response.Error(c, response.CodeBadRequest, "bad_request")
		}
		if err.Error() == "duplicate request" {
			return response.ErrorWithMessage(c, response.CodeBadRequest, "请勿重复提交")
		}
		return response.Error(c, response.CodeInternalError, "internal_error")
	}

	return response.SuccessWithMessage(c, "success", nil)
}
