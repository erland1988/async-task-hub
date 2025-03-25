package model

import "async-task-hub/src/types"

type Application struct {
	ID        int              `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`      // 应用ID
	Name      string           `gorm:"type:varchar(32);not null" json:"name"`                            // 应用名称
	AppKey    string           `gorm:"type:varchar(32);not null;uniqueIndex:app_key" json:"app_key"`     // 应用标识
	AppSecret string           `gorm:"type:varchar(64);not null" json:"app_secret"`                      // 应用秘钥
	AdminID   int              `gorm:"type:bigint(20) unsigned;not null;index:admin_id" json:"admin_id"` // 管理员ID
	Remark    string           `gorm:"type:text;column:remark" json:"remark"`                            // 备注
	CreatedAt types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`
}

func (m *Application) TableName() string {
	return "applications"
}
