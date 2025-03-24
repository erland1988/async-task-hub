package model

import (
	"asynctaskhub/common"
	"encoding/json"
	"time"
)

type Application struct {
	ID        int       `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`      // 应用ID
	Name      string    `gorm:"type:varchar(32);not null" json:"name"`                            // 应用名称
	AppKey    string    `gorm:"type:varchar(32);not null;uniqueIndex:app_key" json:"app_key"`     // 应用标识
	AppSecret string    `gorm:"type:varchar(64);not null" json:"app_secret"`                      // 应用秘钥
	AdminID   int       `gorm:"type:bigint(20) unsigned;not null;index:admin_id" json:"admin_id"` // 管理员ID
	Remark    string    `gorm:"type:text;column:remark" json:"remark"`                            // 备注
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`
}

func (m *Application) MarshalJSON() ([]byte, error) {
	type Alias Application
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(m),
		CreatedAt: common.FormatTime(&m.CreatedAt),
		UpdatedAt: common.FormatTime(&m.UpdatedAt),
	})
}

func (m *Application) TableName() string {
	return "applications"
}
