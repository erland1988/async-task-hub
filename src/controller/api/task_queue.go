package api

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"async-task-hub/src/service/queue"
	"async-task-hub/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type ControllerApiTaskQueue struct {
	ControllerApiBase
}

type requestCreateTaskQueue struct {
	TaskCode           string          `json:"task_code" binding:"required"`
	Parameters         string          `json:"parameters"`
	RelativeDelayTime  int             `json:"relative_delay_time"`
	DelayExecutionTime types.Timestamp `json:"delay_execution_time"`
}

type responseTaskQueue struct {
	ID                    int                            `json:"id"`
	AppID                 int                            `json:"app_id"`
	TaskID                int                            `json:"task_id"`
	Parameters            string                         `json:"parameters"`
	RelativeDelayTime     int                            `json:"relative_delay_time"`
	DelayExecutionTime    types.Timestamp                `json:"delay_execution_time"`
	ExecutionTime         types.Timestamp                `json:"execution_time"`
	ExecutionStatus       model.TaskQueueExecutionStatus `json:"execution_status"`
	ExecutionStatusString string                         `json:"execution_status_string"`
	ExecutionStart        *types.Customtime              `json:"execution_start"`
	ExecutionEnd          *types.Customtime              `json:"execution_end"`
	ExecutionDuration     int64                          `json:"execution_duration"`
	ExecutionCount        int                            `json:"execution_count"`
	CreatedAt             types.Customtime               `json:"created_at"`
	UpdatedAt             types.Customtime               `json:"updated_at"`
	Taskname              string                         `json:"taskname"`
	ExecutorURL           string                         `json:"executor_url"`
	Appname               string                         `json:"appname"`
}

func (c *ControllerApiTaskQueue) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")

	start := ctx.DefaultQuery("start", "")
	end := ctx.DefaultQuery("end", "")

	query := global.DB.Model(&model.TaskQueue{})
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("task_queues.app_id IN ?", adminInfo.AppIDs)
	}

	if start != "" {
		query = query.Where("task_queues.created_at >=?", start)
	}

	if end != "" {
		query = query.Where("task_queues.created_at <=?", end)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取队列列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取队列列表失败", nil)
		return
	}

	var responseTaskQueues []responseTaskQueue
	if err := query.Omit("task_queues.parameters").Order("task_queues.id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&responseTaskQueues).Error; err != nil {
		global.Logger.Warn("获取队列列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取队列列表失败", nil)
		return
	}

	for i, taskQueue := range responseTaskQueues {
		responseTaskQueues[i].ExecutionStatusString = taskQueue.ExecutionStatus.String()
	}

	c.JSONResponse(ctx, true, "获取队列列表成功", gin.H{
		"list":  responseTaskQueues,
		"total": total,
	})
}

func (c *ControllerApiTaskQueue) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}

	query := global.DB.Model(&model.TaskQueue{}).
		Joins("LEFT JOIN tasks ON tasks.id = task_queues.task_id").
		Joins("LEFT JOIN applications ON applications.id = task_queues.app_id").
		Select("task_queues.*, tasks.name as taskname, tasks.executor_url, applications.name as appname")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("task_queues.app_id IN?", adminInfo.AppIDs)
	}
	var responseTaskQueue responseTaskQueue
	if err := query.First(&responseTaskQueue, id).Error; err != nil {
		global.Logger.Warn("获取任务队列详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务队列详情失败", nil)
		return
	}
	responseTaskQueue.ExecutionStatusString = responseTaskQueue.ExecutionStatus.String()
	c.JSONResponse(ctx, true, "获取任务队列详情成功", &responseTaskQueue)
}

func (c *ControllerApiTaskQueue) Create(ctx *gin.Context) {
	var input requestCreateTaskQueue

	appInfo := c.CheckApp(ctx)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	var task model.Task
	if err := global.DB.Where("app_id = ?", appInfo.ID).Where("task_code =?", input.TaskCode).First(&task).Error; err != nil {
		global.Logger.Warn("task_code错误", zap.Error(err))
		c.JSONResponse(ctx, false, "task_code错误", nil)
		return
	}

	var taskQueue model.TaskQueue
	taskQueue.AppID = appInfo.ID
	taskQueue.TaskID = task.ID
	taskQueue.Parameters = input.Parameters
	taskQueue.RelativeDelayTime = input.RelativeDelayTime
	taskQueue.DelayExecutionTime = input.DelayExecutionTime
	taskQueue.ExecutionStatus = model.TaskQueuePending
	taskQueue.ExecutionCount = 0

	if input.RelativeDelayTime != 0 {
		taskQueue.ExecutionTime = types.Timestamp(int(time.Now().Unix()) + input.RelativeDelayTime)
	} else if input.DelayExecutionTime != 0 {
		taskQueue.ExecutionTime = input.DelayExecutionTime
	} else {
		taskQueue.ExecutionTime = types.Timestamp(time.Now().Unix())
	}

	if err := global.DB.Omit("Task").Create(&taskQueue).Error; err != nil {
		global.Logger.Warn("创建任务队列失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建任务队列失败", nil)
		return
	}

	if err := queue.NewQueueService().PushTaskToQueue(&taskQueue); err != nil {
		global.Logger.Warn("创建任务队列失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建任务队列失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "创建任务队列成功", nil)
}
