package model

import (
	"time"
)

type AdminUser struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	NickName  string    `gorm:"column:nick_name" json:"nick_name"`
	Role      int       `gorm:"column:role;default:1" json:"role"` // 1:superadmin 2:user
	Salt      string    `gorm:"column:salt" json:"-"`
	Remark    string    `gorm:"column:remark" json:"remark"`
	Status    int       `gorm:"column:status;default:1" json:"status"` // 1:启用 2:禁用
	CreateBy  int64     `gorm:"column:create_by" json:"create_by"`
	UpdateBy  int64     `gorm:"column:update_by" json:"update_by"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "admin_sys_user"
}
