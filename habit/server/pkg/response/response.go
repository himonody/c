package response

import "github.com/gofiber/fiber/v2"

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	lang := GetLang(c)
	return c.JSON(Response{
		Code: 200,
		Msg:  GetMessage(lang, "success"),
		Data: data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *fiber.Ctx, msgKey string, data interface{}) error {
	lang := GetLang(c)
	return c.JSON(Response{
		Code: 200,
		Msg:  GetMessage(lang, msgKey),
		Data: data,
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code int, msgKey string) error {
	lang := GetLang(c)
	return c.JSON(Response{
		Code: code,
		Msg:  GetMessage(lang, msgKey),
		Data: nil,
	})
}

// ErrorWithMessage 错误响应（自定义消息）
func ErrorWithMessage(c *fiber.Ctx, code int, message string) error {
	return c.JSON(Response{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}

// Common response codes
const (
	CodeSuccess         = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeInternalError   = 500
	CodeValidationError = 422
	CodeTooManyRequests = 429
)
