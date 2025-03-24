package model

import (
	"asynctaskhub/src/types"
	"errors"
	"gorm.io/gorm"
)

type TaskLogAction string

const (
	TaskLogActionStart TaskLogAction = "start" // 开始
	TaskLogActionEnd   TaskLogAction = "end"   // 结束
)

var TaskLogActionString = map[TaskLogAction]string{
	TaskLogActionStart: "开始",
	TaskLogActionEnd:   "结束",
}

func (t TaskLogAction) String() string {
	return TaskLogActionString[t]
}

type TaskLog struct {
	ID             int                  `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`                                                                             // 日志ID
	AppID          int                  `gorm:"type:bigint(20) unsigned;not null;index:app_id,order:1" json:"app_id"`                                                                    // 所属应用ID
	TaskID         int                  `gorm:"type:bigint(20) unsigned;not null;index:task_id,order:1" json:"task_id"`                                                                  // 任务ID
	TaskQueueID    int                  `gorm:"type:bigint(20) unsigned;not null;index:task_queue_id,order:1" json:"task_queue_id"`                                                      // 任务队列ID
	RequestID      string               `gorm:"type:varchar(64);index:request_id,order:1" json:"request_id"`                                                                             // 任务上下文ID
	Action         TaskLogAction        `gorm:"type:varchar(32);not null;index:request_id,order:2;index:app_id,order:2;index:task_id,order:2;index:task_queue_id,order:2" json:"action"` // 任务动作
	Message        string               `gorm:"type:text" json:"message"`                                                                                                                // 日志信息
	MilliTimestamp types.MilliTimestamp `gorm:"type:bigint(20) unsigned;not null" json:"milli_timestamp"`                                                                                // 日志时间戳(毫秒)
	CreatedAt      types.Customtime     `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
}

var validTaskLogActions = map[TaskLogAction]struct{}{
	TaskLogActionStart: {},
	TaskLogActionEnd:   {},
}

func (m *TaskLog) BeforeSave(tx *gorm.DB) error {
	if _, exists := validTaskLogActions[m.Action]; !exists {
		return errors.New("invalid TaskLogAction: " + string(m.Action))
	}
	return nil
}

func (m *TaskLog) TableName() string {
	return "task_logs"
}
