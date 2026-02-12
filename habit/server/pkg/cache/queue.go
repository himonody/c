package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// QueueNames 队列名称常量
const (
	CheckinSettlementQueue = "challenge:checkin:settlement" // 打卡结算队列
)

// CheckinSettlementMessage 打卡结算消息
type CheckinSettlementMessage struct {
	UserChallengeID int64  `json:"userChallengeId"` // 用户挑战ID
	UserID          int64  `json:"userId"`          // 用户ID
	CheckinID       int64  `json:"checkinId"`       // 打卡ID
	CheckinDate     string `json:"checkinDate"`     // 打卡日期 YYYY-MM-DD
	Timestamp       int64  `json:"timestamp"`       // 时间戳
}

// Queue Redis队列操作
type Queue struct {
	client *redis.Client
	logger *zap.Logger
}

// NewQueue 创建队列实例
func NewQueue(client *redis.Client, logger *zap.Logger) *Queue {
	return &Queue{
		client: client,
		logger: logger,
	}
}

// EnqueueCheckinSettlement 将打卡结算任务加入队列
func (q *Queue) EnqueueCheckinSettlement(ctx context.Context, message *CheckinSettlementMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		q.logger.Error("Failed to marshal settlement message", zap.Error(err))
		return err
	}

	// 使用LPUSH将消息加入队列左侧
	err = q.client.LPush(ctx, CheckinSettlementQueue, data).Err()
	if err != nil {
		q.logger.Error("Failed to enqueue settlement message", zap.Error(err))
		return err
	}

	q.logger.Info("Settlement message enqueued successfully", 
		zap.Int64("userChallengeId", message.UserChallengeID),
		zap.Int64("userId", message.UserID),
		zap.String("checkinDate", message.CheckinDate))

	return nil
}

// DequeueCheckinSettlement 从队列中取出打卡结算任务
func (q *Queue) DequeueCheckinSettlement(ctx context.Context, timeout time.Duration) (*CheckinSettlementMessage, error) {
	// 使用BRPOP阻塞式获取消息，超时时间为timeout
	result, err := q.client.BRPop(ctx, timeout, CheckinSettlementQueue).Result()
	if err != nil {
		if err == redis.Nil {
			// 队列为空，超时
			return nil, nil
		}
		q.logger.Error("Failed to dequeue settlement message", zap.Error(err))
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid redis result format")
	}

	var message CheckinSettlementMessage
	err = json.Unmarshal([]byte(result[1]), &message)
	if err != nil {
		q.logger.Error("Failed to unmarshal settlement message", zap.Error(err))
		return nil, err
	}

	return &message, nil
}

// GetQueueLength 获取队列长度
func (q *Queue) GetQueueLength(ctx context.Context, queueName string) (int64, error) {
	length, err := q.client.LLen(ctx, queueName).Result()
	if err != nil {
		q.logger.Error("Failed to get queue length", zap.String("queue", queueName), zap.Error(err))
		return 0, err
	}
	return length, nil
}

// Enqueue 通用入队方法
func (q *Queue) Enqueue(ctx context.Context, queueName string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return q.client.LPush(ctx, queueName, data).Err()
}

// Dequeue 通用出队方法
func (q *Queue) Dequeue(ctx context.Context, queueName string, timeout time.Duration, dest interface{}) error {
	result, err := q.client.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // 队列为空，超时
		}
		return err
	}

	if len(result) < 2 {
		return fmt.Errorf("invalid redis result format")
	}

	return json.Unmarshal([]byte(result[1]), dest)
}
