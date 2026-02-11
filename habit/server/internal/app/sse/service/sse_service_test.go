package service

import (
	"testing"
	"time"

	"habit/internal/app/sse/dto"
	"habit/pkg/database"

	"go.uber.org/zap/zaptest"
)

// MockLogger 用于测试
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...interface{}) {}
func (m *MockLogger) Error(msg string, fields ...interface{}) {}

func TestSSEService_SendToUser(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	// 模拟添加两个客户端
	client1 := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client2 := &SSEClient{
		ID:      "client2",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client3 := &SSEClient{
		ID:      "client3",
		UserID:  1002,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}

	svc.AddClient(1001, client1)
	svc.AddClient(1001, client2)
	svc.AddClient(1002, client3)

	// 单播给用户 1001
	svc.SendToUser(1001, "test", "hello user 1001")

	// 验证客户端 1001 收到消息
	select {
	case msg := <-client1.Channel:
		if msg.Event != "test" || msg.Data != "hello user 1001" {
			t.Errorf("client1 got wrong message: %+v", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("client1 should receive message")
	}

	select {
	case msg := <-client2.Channel:
		if msg.Event != "test" || msg.Data != "hello user 1001" {
			t.Errorf("client2 got wrong message: %+v", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("client2 should receive message")
	}

	// 验证客户端 1002 没收到消息
	select {
	case <-client3.Channel:
		t.Fatal("client3 should not receive message")
	case <-time.After(50 * time.Millisecond):
		// 正常，没收到
	}
}

func TestSSEService_Broadcast(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	// 模拟添加三个客户端
	client1 := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client2 := &SSEClient{
		ID:      "client2",
		UserID:  1002,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client3 := &SSEClient{
		ID:      "client3",
		UserID:  1003,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}

	svc.AddClient(1001, client1)
	svc.AddClient(1002, client2)
	svc.AddClient(1003, client3)

	// 广播
	svc.Broadcast("broadcast", "hello everyone")

	// 验证所有客户端都收到
	for i, client := range []*SSEClient{client1, client2, client3} {
		select {
		case msg := <-client.Channel:
			if msg.Event != "broadcast" || msg.Data != "hello everyone" {
				t.Errorf("client%d got wrong message: %+v", i+1, msg)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("client%d should receive broadcast message", i+1)
		}
	}
}

func TestSSEService_RemoveClient(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	client := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}

	svc.AddClient(1001, client)
	// 验证用户在线
	if !svc.IsUserOnline(1001) {
		t.Fatal("User should be online after adding client")
	}

	// 验证客户端存在
	svc.SendToUser(1001, "test", "before remove")
	select {
	case <-client.Channel:
		// 正常收到
	case <-time.After(100 * time.Millisecond):
		t.Fatal("client should receive message before remove")
	}

	// 移除客户端
	svc.RemoveClient(1001, client.ID)

	// 验证用户离线
	if svc.IsUserOnline(1001) {
		t.Fatal("User should be offline after removing all clients")
	}

	// 再次发送，应该收不到
	svc.SendToUser(1001, "test", "after remove")
	select {
	case <-client.Channel:
		t.Fatal("client should not receive message after remove")
	case <-time.After(50 * time.Millisecond):
		// 正常，没收到
	}
}

func TestSSEService_IsUserOnline(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	// 初始状态离线
	if svc.IsUserOnline(1001) {
		t.Fatal("User should be offline initially")
	}

	// 添加客户端后在线
	client := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	svc.AddClient(1001, client)

	if !svc.IsUserOnline(1001) {
		t.Fatal("User should be online after adding client")
	}

	// 移除客户端后离线
	svc.RemoveClient(1001, client.ID)

	if svc.IsUserOnline(1001) {
		t.Fatal("User should be offline after removing client")
	}
}

func TestSSEService_GetUsersOnlineStatus(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	// 测试空数组
	emptyResult := svc.GetUsersOnlineStatus([]int64{})
	if len(emptyResult) != 0 {
		t.Fatal("Empty input should return empty result")
	}

	// 添加三个用户，其中两个在线
	client1 := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client2 := &SSEClient{
		ID:      "client2",
		UserID:  1002,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	// 1003 不添加客户端，应该离线

	svc.AddClient(1001, client1)
	svc.AddClient(1002, client2)

	// 批量查询状态
	userIDs := []int64{1001, 1002, 1003, 1004} // 1004 也不存在
	statusMap := svc.GetUsersOnlineStatus(userIDs)

	// 验证结果
	if len(statusMap) != 4 {
		t.Fatalf("Expected 4 results, got %d", len(statusMap))
	}

	if !statusMap[1001] {
		t.Error("User 1001 should be online")
	}
	if !statusMap[1002] {
		t.Error("User 1002 should be online")
	}
	if statusMap[1003] {
		t.Error("User 1003 should be offline")
	}
	if statusMap[1004] {
		t.Error("User 1004 should be offline")
	}

	// 移除一个用户后再次查询
	svc.RemoveClient(1001, client1.ID)
	statusMap = svc.GetUsersOnlineStatus([]int64{1001, 1002})

	if statusMap[1001] {
		t.Error("User 1001 should be offline after removal")
	}
	if !statusMap[1002] {
		t.Error("User 1002 should still be online")
	}
}

func TestSSEService_GetOnlineUsers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	svc := NewSSEService(logger, nil) // 测试时 userRepo 可以为 nil

	// 检查 Redis 是否可用
	if database.RedisClient == nil {
		t.Skip("Redis not available, skipping test")
	}

	// 初始无在线用户
	online := svc.GetOnlineUsers()
	if len(online) != 0 {
		t.Fatalf("Expected 0 online users, got %d", len(online))
	}

	// 添加两个在线用户
	client1 := &SSEClient{
		ID:      "client1",
		UserID:  1001,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}
	client2 := &SSEClient{
		ID:      "client2",
		UserID:  1002,
		Channel: make(chan *dto.SSEMessage, 10),
		Done:    make(chan struct{}),
	}

	svc.AddClient(1001, client1)
	svc.AddClient(1002, client2)

	online = svc.GetOnlineUsers()
	if len(online) != 2 {
		t.Fatalf("Expected 2 online users, got %d", len(online))
	}

	// 移除一个用户
	svc.RemoveClient(1001, client1.ID)

	online = svc.GetOnlineUsers()
	if len(online) != 1 {
		t.Fatalf("Expected 1 online user, got %d", len(online))
	}
	if online[0] != 1002 {
		t.Fatalf("Expected user 1002 online, got %d", online[0])
	}
}
