package model

import "async-task-hub/src/types"

type Login struct {
	ID        int              `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`      // 主键ID
	AdminID   int              `gorm:"type:bigint(20) unsigned;not null;index:admin_id" json:"admin_id"` // 管理员ID
	Token     string           `gorm:"type:varchar(64);not null;uniqueIndex:token" json:"token"`         // 登录Token
	ExpiresAt types.Customtime `gorm:"type:datetime;not null;index:expires_at" json:"expires_at"`        // 过期时间
	CreatedAt types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`

	Admin Admin `json:"admin" gorm:"foreignKey:AdminID"`
}

func (m *Login) TableName() string {
	return "logins"
}
