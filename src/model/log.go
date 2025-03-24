package model

import (
	"asynctaskhub/common"
	"encoding/json"
	"time"
)

type Log struct {
	Id        int       `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`      // 日志ID
	Operation string    `gorm:"type:varchar(50);not null;index:operation" json:"operation"`       // 操作类型
	Details   string    `gorm:"type:text" json:"details"`                                         // 操作详情
	AdminID   int       `gorm:"type:bigint(20) unsigned;not null;index:admin_id" json:"admin_id"` // 管理员ID
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`

	Admin Admin `gorm:"foreignKey:AdminID" json:"admin"`
}

func (m *Log) MarshalJSON() ([]byte, error) {
	type Alias Log
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     (*Alias)(m),
		CreatedAt: common.FormatTime(&m.CreatedAt),
	})
}

func (m *Log) TableName() string {
	return "logs"
}
