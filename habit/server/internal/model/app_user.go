package model

import (
	"time"
)

type AppUser struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	LevelID      int        `gorm:"column:level_id;not null;default:1" json:"level_id"`
	Username     string     `gorm:"column:username;not null;uniqueIndex:uk_username" json:"username"`
	Nickname     string     `gorm:"column:nickname;not null;default:''" json:"nickname"`
	Avatar       string     `gorm:"column:avatar" json:"avatar"`
	Pwd          string     `gorm:"column:pwd;not null" json:"-"`
	RefCode      string     `gorm:"column:ref_code;uniqueIndex:uk_ref_code" json:"ref_code"`
	RefID        int        `gorm:"column:ref_id;not null;default:0;index:idx_ref_id" json:"ref_id"`
	FriendCode   string     `gorm:"column:friend_code" json:"friend_code"`
	FriendID     string     `gorm:"column:friend_id" json:"friend_id"`
	Status       string     `gorm:"column:status;not null;default:'1';index:idx_status" json:"status"`
	OnlineStatus string     `gorm:"column:online_status;not null;default:'1';index:idx_online_status" json:"online_status"`
	RegisterAt   time.Time  `gorm:"column:register_at;not null;default:CURRENT_TIMESTAMP;index:idx_register_at" json:"register_at"`
	RegisterIP   string     `gorm:"column:register_ip" json:"register_ip"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at;index:idx_last_login_at" json:"last_login_at"`
	LastLoginIP  string     `gorm:"column:last_login_ip" json:"last_login_ip"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AppUser) TableName() string {
	return "app_user"
}
