package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"habit/internal/app/sse/dto"
	"habit/internal/repo"
	"habit/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SSEClient struct {
	ID      string
	UserID  int64
	Channel chan *dto.SSEMessage
	Done    chan struct{}
}

type SSEService struct {
	logger   *zap.Logger
	rdb      *redis.Client
	userRepo *repo.UserRepository

	// userID -> []*SSEClient
	clients sync.Map // map[int64][]*SSEClient

	// 在线状态配置
	onlineKey     string        // Redis ZSET key
	heartbeatTTL  time.Duration // 心跳间隔
	cleanupTicker time.Duration // 清理间隔
	offlineAfter  time.Duration // 超过此时间视为离线
}

func NewSSEService(logger *zap.Logger, userRepo *repo.UserRepository) *SSEService {
	svc := &SSEService{
		logger:        logger,
		rdb:           database.RedisClient,
		userRepo:      userRepo,
		onlineKey:     "sse:online",
		heartbeatTTL:  30 * time.Second, // 心跳间隔
		cleanupTicker: 5 * time.Minute,  // 清理间隔
		offlineAfter:  5 * time.Minute,  // 超过5分钟视为离线
	}

	// 启动清理协程
	go svc.startCleanup()

	return svc
}

// AddClient 添加客户端连接
func (s *SSEService) AddClient(userID int64, client *SSEClient) {
	if client == nil {
		s.logger.Error("Cannot add nil client", zap.Int64("userID", userID))
		return
	}

	if client.Channel == nil {
		s.logger.Error("Client channel is nil", zap.Int64("userID", userID), zap.String("clientID", client.ID))
		return
	}

	var list []*SSEClient
	if val, ok := s.clients.Load(userID); ok {
		if existingList, ok := val.([]*SSEClient); ok {
			list = existingList
		} else {
			s.logger.Error("Invalid client list type in sync.Map", zap.Int64("userID", userID))
			list = []*SSEClient{}
		}
	}

	// 检查是否已存在相同ID的客户端
	for _, existingClient := range list {
		if existingClient.ID == client.ID {
			s.logger.Warn("Client ID already exists, skipping", zap.Int64("userID", userID), zap.String("clientID", client.ID))
			return
		}
	}

	list = append(list, client)
	s.clients.Store(userID, list)

	// 更新用户在线状态（Redis ZSET）
	s.updateUserOnline(userID)
	s.logger.Info("User came online", zap.Int64("userID", userID), zap.String("clientID", client.ID))

	// 启动清理协程
	go s.monitorClient(userID, client)
}

// RemoveClient 移除客户端连接
func (s *SSEService) RemoveClient(userID int64, clientID string) {
	if clientID == "" {
		s.logger.Error("Cannot remove client with empty ID", zap.Int64("userID", userID))
		return
	}

	if val, ok := s.clients.Load(userID); ok {
		if existingList, ok := val.([]*SSEClient); ok {
			list := existingList
			var newList []*SSEClient
			found := false

			for _, c := range list {
				if c.ID == clientID {
					found = true
					// 安全关闭客户端通道
					if c.Channel != nil {
						select {
						case <-c.Done:
							// 已经关闭
						default:
							close(c.Done)
						}
					}
					continue
				}
				newList = append(newList, c)
			}

			if !found {
				s.logger.Warn("Client ID not found for removal", zap.Int64("userID", userID), zap.String("clientID", clientID))
				return
			}

			if len(newList) == 0 {
				s.clients.Delete(userID)
				// 如果用户没有连接了，从 Redis ZSET 中移除
				s.removeUserOnline(userID)
				s.logger.Info("User went offline", zap.Int64("userID", userID), zap.String("clientID", clientID))
			} else {
				s.clients.Store(userID, newList)
				s.logger.Info("Client removed, user still online", zap.Int64("userID", userID), zap.String("clientID", clientID), zap.Int("remainingClients", len(newList)))
			}
		} else {
			s.logger.Error("Invalid client list type in sync.Map during removal", zap.Int64("userID", userID))
		}
	} else {
		s.logger.Warn("No clients found for user during removal", zap.Int64("userID", userID), zap.String("clientID", clientID))
	}
}

// monitorClient 监控客户端连接，断开时清理
func (s *SSEService) monitorClient(userID int64, client *SSEClient) {
	if client == nil || client.Done == nil {
		s.logger.Error("Invalid client for monitoring", zap.Int64("userID", userID))
		return
	}

	select {
	case <-client.Done:
		s.RemoveClient(userID, client.ID)
		s.logger.Info("SSE client disconnected", zap.Int64("userID", userID), zap.String("clientID", client.ID))
	case <-time.After(30 * time.Minute):
		// 超时清理
		s.RemoveClient(userID, client.ID)
		s.logger.Info("SSE client timeout removed", zap.Int64("userID", userID), zap.String("clientID", client.ID))
	}
}

// SendToUser 单播给指定用户的所有连接
func (s *SSEService) SendToUser(userID int64, event, data string) {
	if userID <= 0 {
		s.logger.Error("Invalid user ID for SendToUser", zap.Int64("userID", userID))
		return
	}

	if event == "" {
		s.logger.Error("Empty event for SendToUser", zap.Int64("userID", userID))
		return
	}

	if val, ok := s.clients.Load(userID); ok {
		if existingList, ok := val.([]*SSEClient); ok {
			list := existingList
			if len(list) == 0 {
				s.logger.Debug("No active clients for user", zap.Int64("userID", userID))
				return
			}

			msg := &dto.SSEMessage{Event: event, Data: data}
			sentCount := 0
			
			for _, client := range list {
				if client == nil || client.Channel == nil {
					s.logger.Warn("Invalid client in list", zap.Int64("userID", userID))
					continue
				}

				select {
				case client.Channel <- msg:
					sentCount++
				default:
					// 通道满，跳过
					continue
				}
			}
			
			s.logger.Debug("Message sent to user", zap.Int64("userID", userID), zap.String("event", event), zap.Int("sentCount", sentCount), zap.Int("totalClients", len(list)))
		} else {
			s.logger.Error("Invalid client list type in sync.Map for SendToUser", zap.Int64("userID", userID))
		}
	} else {
		s.logger.Debug("User not connected", zap.Int64("userID", userID))
	}
}

// Broadcast 广播给所有连接
func (s *SSEService) Broadcast(event, data string) {
	if event == "" {
		s.logger.Error("Empty event for Broadcast")
		return
	}

	msg := &dto.SSEMessage{Event: event, Data: data}
	totalSent := 0
	totalClients := 0
	
	s.clients.Range(func(key, value interface{}) bool {
		if existingList, ok := value.([]*SSEClient); ok {
			list := existingList
			totalClients += len(list)
			
			for _, client := range list {
				if client == nil || client.Channel == nil {
					s.logger.Warn("Invalid client in broadcast list")
					continue
				}

				select {
				case client.Channel <- msg:
					totalSent++
				default:
					// 通道满，跳过
					continue
				}
			}
		} else {
			s.logger.Error("Invalid client list type in sync.Map during broadcast")
		}
		return true
	})
	
	s.logger.Debug("Broadcast message sent", zap.String("event", event), zap.Int("sentCount", totalSent), zap.Int("totalClients", totalClients))
}

// Stream SSE 流响应
func (s *SSEService) Stream(c *fiber.Ctx, userID int64) error {
	if c == nil {
		s.logger.Error("Nil fiber context in Stream")
		return errors.New("fiber context is nil")
	}

	if userID <= 0 {
		s.logger.Error("Invalid user ID in Stream", zap.Int64("userID", userID))
		return errors.New("invalid user ID")
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	clientID := generateClientID()
	if clientID == "" {
		s.logger.Error("Failed to generate client ID", zap.Int64("userID", userID))
		return errors.New("failed to generate client ID")
	}

	client := &SSEClient{
		ID:      clientID,
		UserID:  userID,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}

	s.AddClient(userID, client)
	defer func() {
		if client.Done != nil {
			select {
			case <-client.Done:
				// 已经关闭
			default:
				close(client.Done)
			}
		}
	}()

	// 发送初始连接事件
	initMsg := &dto.SSEMessage{Event: "connected", Data: `{"status":"connected"}`}
	if err := s.writeEvent(c, initMsg); err != nil {
		s.logger.Error("Failed to write initial event", zap.Error(err), zap.Int64("userID", userID), zap.String("clientID", clientID))
		return err
	}

	// 监听消息推送
	for {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				s.logger.Info("Client channel closed", zap.Int64("userID", userID), zap.String("clientID", clientID))
				return nil
			}
			if msg == nil {
				s.logger.Warn("Received nil message", zap.Int64("userID", userID), zap.String("clientID", clientID))
				continue
			}
			if err := s.writeEvent(c, msg); err != nil {
				s.logger.Error("Failed to write event", zap.Error(err), zap.Int64("userID", userID), zap.String("clientID", clientID))
				return err
			}
		case <-c.Context().Done():
			s.logger.Info("Client context done", zap.Int64("userID", userID), zap.String("clientID", clientID))
			return nil
		}
	}
}

func (s *SSEService) writeEvent(c *fiber.Ctx, msg *dto.SSEMessage) error {
	data, _ := json.Marshal(msg.Data)
	if _, err := c.Writef("event: %s\ndata: %s\n\n", msg.Event, string(data)); err != nil {
		return err
	}
	if err := c.Context().Err(); err != nil {
		return err
	}
	return c.Send([]byte{})
}

func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// updateUserOnline 更新用户在线状态（Redis ZSET + 数据库）
func (s *SSEService) updateUserOnline(userID int64) {
	ctx := context.Background()
	score := float64(time.Now().Unix())

	// 更新 Redis ZSET
	err := s.rdb.ZAdd(ctx, s.onlineKey, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
	if err != nil {
		s.logger.Error("Failed to update user online status in Redis", zap.Error(err), zap.Int64("userID", userID))
	}

	// 更新数据库在线状态
	if s.userRepo != nil {
		err = s.userRepo.UpdateOnlineStatus(userID, "2") // "2" 表示在线
		if err != nil {
			s.logger.Error("Failed to update user online status in database", zap.Error(err), zap.Int64("userID", userID))
		}
	}
}

// removeUserOnline 移除用户在线状态（Redis ZSET + 数据库）
func (s *SSEService) removeUserOnline(userID int64) {
	ctx := context.Background()

	// 从 Redis ZSET 中移除
	err := s.rdb.ZRem(ctx, s.onlineKey, userID).Err()
	if err != nil {
		s.logger.Error("Failed to remove user online status from Redis", zap.Error(err), zap.Int64("userID", userID))
	}

	// 更新数据库在线状态
	if s.userRepo != nil {
		err = s.userRepo.UpdateOnlineStatus(userID, "1") // "1" 表示离线
		if err != nil {
			s.logger.Error("Failed to update user offline status in database", zap.Error(err), zap.Int64("userID", userID))
		}
	}
}

// IsUserOnline 检查用户是否在线（Redis ZSET）
func (s *SSEService) IsUserOnline(userID int64) bool {
	ctx := context.Background()
	cutoff := float64(time.Now().Add(-s.offlineAfter).Unix())

	// 检查用户在 ZSET 中的分数是否大于截止时间
	result, err := s.rdb.ZScore(ctx, s.onlineKey, strconv.FormatInt(userID, 10)).Result()
	if errors.Is(err, redis.Nil) {
		return false
	}
	if err != nil {
		s.logger.Error("Failed to check user online status", zap.Error(err), zap.Int64("userID", userID))
		return false
	}

	return result > cutoff
}

// GetOnlineUsers 获取所有在线用户ID列表（Redis ZSET）
func (s *SSEService) GetOnlineUsers() []int64 {
	ctx := context.Background()
	cutoff := float64(time.Now().Add(-s.offlineAfter).Unix())

	// 获取所有分数大于截止时间的用户
	result, err := s.rdb.ZRangeByScore(ctx, s.onlineKey, &redis.ZRangeBy{
		Min: strconv.FormatFloat(cutoff, 'f', 0, 64),
		Max: "+inf",
	}).Result()
	if err != nil {
		s.logger.Error("Failed to get online users", zap.Error(err))
		return nil
	}

	var users []int64
	for _, idStr := range result {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			users = append(users, id)
		}
	}
	return users
}

// GetUsersOnlineStatus 批量查询用户在线状态
func (s *SSEService) GetUsersOnlineStatus(userIDs []int64) map[int64]bool {
	if len(userIDs) == 0 {
		return map[int64]bool{}
	}
	ctx := context.Background()
	cutoff := float64(time.Now().Add(-s.offlineAfter).Unix())

	result := make(map[int64]bool, len(userIDs))

	// 构建 Redis Pipeline 批量查询
	pipeline := s.rdb.Pipeline()
	cmds := make([]*redis.FloatCmd, len(userIDs))

	for i, userID := range userIDs {
		cmds[i] = pipeline.ZScore(ctx, s.onlineKey, strconv.FormatInt(userID, 10))
	}

	// 执行批量查询
	_, err := pipeline.Exec(ctx)
	if err != nil {
		s.logger.Error("Failed to batch query user online status", zap.Error(err))
		// 如果批量查询失败，逐个查询
		for _, userID := range userIDs {
			result[userID] = s.IsUserOnline(userID)
		}
		return result
	}

	// 处理结果
	for i, userID := range userIDs {
		score, err := cmds[i].Result()
		if errors.Is(err, redis.Nil) {
			result[userID] = false
		} else if err != nil {
			s.logger.Error("Failed to get user score", zap.Error(err), zap.Int64("userID", userID))
			result[userID] = false
		} else {
			result[userID] = score > cutoff
		}
	}

	return result
}

// startCleanup 启动清理协程
func (s *SSEService) startCleanup() {
	ticker := time.NewTicker(s.cleanupTicker)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredUsers()
	}
}

// cleanupExpiredUsers 清理过期的在线用户
func (s *SSEService) cleanupExpiredUsers() {
	ctx := context.Background()
	cutoff := float64(time.Now().Add(-s.offlineAfter).Unix())

	// 获取需要移除的用户ID列表
	expiredUsers, err := s.rdb.ZRangeByScore(ctx, s.onlineKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatFloat(cutoff, 'f', 0, 64),
	}).Result()
	if err != nil {
		s.logger.Error("Failed to get expired users", zap.Error(err))
		return
	}

	// 从 Redis ZSET 中移除过期用户
	removed, err := s.rdb.ZRemRangeByScore(ctx, s.onlineKey, "-inf", strconv.FormatFloat(cutoff, 'f', 0, 64)).Result()
	if err != nil {
		s.logger.Error("Failed to cleanup expired users", zap.Error(err))
		return
	}

	if removed > 0 {
		s.logger.Info("Cleaned up expired users", zap.Int64("count", removed))

		// 批量更新数据库中的在线状态为离线
		if s.userRepo != nil && len(expiredUsers) > 0 {
			s.batchUpdateOfflineStatus(expiredUsers)
		}
	}
}

// batchUpdateOfflineStatus 批量更新用户离线状态
func (s *SSEService) batchUpdateOfflineStatus(userIDs []string) {
	if len(userIDs) == 0 {
		return
	}

	// 转换字符串ID为int64
	userIDInts := make([]int64, 0, len(userIDs))
	for _, idStr := range userIDs {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			userIDInts = append(userIDInts, id)
		}
	}

	// 批量更新数据库
	for _, userID := range userIDInts {
		err := s.userRepo.UpdateOnlineStatus(userID, "1") // "1" 表示离线
		if err != nil {
			s.logger.Error("Failed to batch update user offline status", zap.Error(err), zap.Int64("userID", userID))
		}
	}

	s.logger.Info("Batch updated offline status", zap.Int("count", len(userIDInts)))
}

// SyncOnlineStatusToDB 同步 Redis 在线状态到数据库
func (s *SSEService) SyncOnlineStatusToDB() error {
	if s.userRepo == nil {
		return errors.New("user repository is nil")
	}

	// 获取所有在线用户
	onlineUsers := s.GetOnlineUsers()
	if len(onlineUsers) == 0 {
		return nil
	}

	// 批量设置在线用户状态
	for _, userID := range onlineUsers {
		err := s.userRepo.UpdateOnlineStatus(userID, "2") // "2" 表示在线
		if err != nil {
			s.logger.Error("Failed to sync online status", zap.Error(err), zap.Int64("userID", userID))
		}
	}

	s.logger.Info("Synced online status to database", zap.Int("count", len(onlineUsers)))
	return nil
}

// StartConsumer 启动 Redis 队列消费者
func (s *SSEService) StartConsumer(ctx context.Context, queueKey string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("SSE consumer stopped")
				return
			default:
				result, err := s.rdb.BRPop(ctx, 0, queueKey).Result()
				if err != nil {
					if !errors.Is(err, redis.Nil) {
						s.logger.Error("Failed to pop from SSE queue", zap.Error(err))
					}
					time.Sleep(1 * time.Second)
					continue
				}

				if len(result) < 2 {
					continue
				}

				var msg dto.SSEQueueMessage
				if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
					s.logger.Error("Failed to unmarshal SSE queue message", zap.Error(err), zap.String("data", result[1]))
					continue
				}

				// 分发消息
				if msg.UserID != nil {
					s.SendToUser(*msg.UserID, msg.Event, msg.Data)
				} else {
					s.Broadcast(msg.Event, msg.Data)
				}
			}
		}
	}()
}
