package model

import "async-task-hub/src/types"

type Task struct {
	ID          int              `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`      // 任务ID
	AppID       int              `gorm:"type:bigint(20) unsigned;not null;index:app_id" json:"app_id"`     // 所属应用ID
	Name        string           `gorm:"type:varchar(32);not null" json:"name"`                            // 任务名称
	TaskCode    string           `gorm:"type:varchar(32);not null;uniqueIndex:task_code" json:"task_code"` // 任务标识
	ExecutorURL string           `gorm:"type:varchar(64);not null" json:"executor_url"`                    // 执行器URL
	CreatedAt   types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt   types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`
}

func (m *Task) TableName() string {
	return "tasks"
}
