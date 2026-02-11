package sse_test

import (
	"habit/internal/app/sse/service"
	"habit/pkg/logger"

	"go.uber.org/zap"
)

// ExampleBusinessService 演示如何在业务代码中注入 SSE Service 并调用发送
type ExampleBusinessService struct {
	sseSvc *service.SSEService
	logger *zap.Logger
}

func NewExampleBusinessService(sseSvc *service.SSEService) *ExampleBusinessService {
	return &ExampleBusinessService{
		sseSvc: sseSvc,
		logger: logger.Logger,
	}
}

// OnUserLogin 用户登录时推送通知
func (b *ExampleBusinessService) OnUserLogin(userID int64, username string) {
	// 本地方法调用发送（单播）
	b.sseSvc.SendToUser(userID, "login", `{"message":"欢迎回来","username":"`+username+`"}`)
	b.logger.Info("Sent login notification", zap.Int64("userID", userID))
}

// OnSystemAnnouncement 系统公告推送（广播）
func (b *ExampleBusinessService) OnSystemAnnouncement(title, content string) {
	// 本地方法调用发送（广播）
	b.sseSvc.Broadcast("announcement", `{"title":"`+title+`","content":"`+content+`"}`)
	b.logger.Info("Broadcast system announcement", zap.String("title", title))
}

// OnChallengeCompleted 挑战完成推送（单播）
func (b *ExampleBusinessService) OnChallengeCompleted(userID int64, challengeID int64) {
	data := `{"challengeId":` + string(rune(challengeID)) + `,"message":"挑战完成"}`
	b.sseSvc.SendToUser(userID, "challenge", data)
	b.logger.Info("Sent challenge completion", zap.Int64("userID", userID), zap.Int64("challengeID", challengeID))
}

// DemoBusinessUsage 在业务代码中使用 SSE Service 发送消息
func DemoBusinessUsage() {
	// 初始化 SSE Service（通常在 router 层初始化，这里仅作演示）
	sseSvc := service.NewSSEService(logger.Logger, nil) // 示例中 userRepo 为 nil

	// 初始化业务服务，注入 SSE Service
	bizSvc := NewExampleBusinessService(sseSvc)

	// 业务场景调用
	bizSvc.OnUserLogin(1001, "alice")
	bizSvc.OnSystemAnnouncement("系统维护", "今晚 22:00 进行系统维护")
	bizSvc.OnChallengeCompleted(1001, 123)
}
