package model

import (
	"asynctaskhub/src/types"
	"errors"
	"gorm.io/gorm"
)

type TaskQueueExecutionStatus string

const (
	TaskQueuePending    TaskQueueExecutionStatus = "pending"    // 待执行
	TaskQueueProcessing TaskQueueExecutionStatus = "processing" // 执行中
	TaskQueueCompleted  TaskQueueExecutionStatus = "completed"  // 已完成
	TaskQueueFailed     TaskQueueExecutionStatus = "failed"     // 已失败
)

var TaskQueueExecutionStatusString = map[TaskQueueExecutionStatus]string{
	TaskQueuePending:    "待执行",
	TaskQueueProcessing: "执行中",
	TaskQueueCompleted:  "已完成",
	TaskQueueFailed:     "已失败",
}

func (t TaskQueueExecutionStatus) String() string {
	return TaskQueueExecutionStatusString[t]
}

type TaskQueue struct {
	ID                 int                      `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`                                               // 任务队列ID
	AppID              int                      `gorm:"type:bigint(20) unsigned;not null;index:app_id" json:"app_id"`                                              // 所属应用ID
	TaskID             int                      `gorm:"type:bigint(20) unsigned;not null;index:task_id" json:"task_id"`                                            // 所属任务ID
	Parameters         string                   `gorm:"type:text" json:"parameters"`                                                                               // 任务参数
	RelativeDelayTime  int                      `gorm:"type:int(11) unsigned;default:NULL" json:"relative_delay_time"`                                             // 相对延迟时间（以秒为单位）
	DelayExecutionTime types.Timestamp          `gorm:"type:int(11) unsigned;default:NULL" json:"delay_execution_time"`                                            // 绝对时间戳
	ExecutionTime      types.Timestamp          `gorm:"type:int(11) unsigned;default:NULL;index:execution_time,order:1" json:"execution_time"`                     // 确定的具体执行时间
	ExecutionStatus    TaskQueueExecutionStatus `gorm:"type:varchar(32);default:NULL;index:execution_status;index:execution_time,order:2" json:"execution_status"` // 执行状态
	ExecutionStart     *types.Customtime        `gorm:"type:datetime" json:"execution_start"`                                                                      // 任务开始时间
	ExecutionEnd       *types.Customtime        `gorm:"type:datetime" json:"execution_end"`                                                                        // 任务结束时间
	ExecutionDuration  int64                    `gorm:"type:int(11) unsigned;default:NULL" json:"execution_duration"`                                              // 任务执行时长(毫秒)
	ExecutionCount     int                      `gorm:"type:int(11) unsigned;default:0" json:"execution_count"`                                                    // 执行次数
	CreatedAt          types.Customtime         `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt          types.Customtime         `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`

	Task Task `gorm:"foreignKey:TaskID" json:"task"`
}

var validTaskQueueExecutionStatus = map[TaskQueueExecutionStatus]struct{}{
	TaskQueuePending:    {},
	TaskQueueProcessing: {},
	TaskQueueCompleted:  {},
	TaskQueueFailed:     {},
}

func (m *TaskQueue) BeforeSave(tx *gorm.DB) error {
	if _, exists := validTaskQueueExecutionStatus[m.ExecutionStatus]; !exists {
		return errors.New("invalid execution_status: " + string(m.ExecutionStatus))
	}
	return nil
}

func (m *TaskQueue) TableName() string {
	return "task_queues"
}
