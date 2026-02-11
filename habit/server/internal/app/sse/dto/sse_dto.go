package dto

// SSEQueueMessage Redis 队列消息格式
type SSEQueueMessage struct {
	UserID *int64 `json:"userId,omitempty"` // nil 表示广播
	Event  string `json:"event"`            // 事件名
	Data   string `json:"data"`             // 数据内容
}

// SSESendRequest HTTP POST 请求格式（本地调用）
type SSESendRequest struct {
	UserID *int64 `json:"userId,omitempty"` // nil 表示广播
	Event  string `json:"event"`            // 事件名
	Data   string `json:"data"`             // 数据内容
}

type SSEMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
