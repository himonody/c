package model

import "time"

// AdminSysConfig 系统配置表
type AdminSysConfig struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName  string    `gorm:"column:config_name;size:128" json:"config_name"`
	ConfigKey   string    `gorm:"column:config_key;size:128;index" json:"config_key"`
	ConfigValue string    `gorm:"column:config_value;size:255" json:"config_value"`
	ConfigType  string    `gorm:"column:config_type;size:64" json:"config_type"`
	IsFrontend  string    `gorm:"column:is_frontend;size:1" json:"is_frontend"` // Y/N
	Remark      string    `gorm:"column:remark;size:255" json:"remark"`
	CreateBy    int64     `gorm:"column:create_by" json:"create_by"`
	UpdateBy    int64     `gorm:"column:update_by" json:"update_by"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (AdminSysConfig) TableName() string {
	return "admin_sys_config"
}
