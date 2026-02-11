package dto

// ConfigListRequest 配置列表请求
type ConfigListRequest struct {
	Page       int    `json:"page" form:"page"`             // 页码
	PageSize   int    `json:"pageSize" form:"pageSize"`   // 每页数量
	ConfigName string `json:"configName" form:"configName"` // 配置名称
	ConfigKey  string `json:"configKey" form:"configKey"`   // 配置键
}

// ConfigListResponse 配置列表响应
type ConfigListResponse struct {
	List     []*ConfigInfo `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// ConfigInfo 配置信息
type ConfigInfo struct {
	ID          int64  `json:"id"`
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	IsFrontend  string `json:"isFrontend"`
	Remark      string `json:"remark"`
	CreateBy    int64  `json:"createBy"`
	UpdateBy    int64  `json:"updateBy"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	IsFrontend  string `json:"isFrontend"` // Y/N
	Remark      string `json:"remark"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	ID          int64  `json:"id" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	IsFrontend  string `json:"isFrontend"` // Y/N
	Remark      string `json:"remark"`
}

// GetConfigRequest 获取配置请求
type GetConfigRequest struct {
	ID int64 `json:"id" form:"id"`
}
