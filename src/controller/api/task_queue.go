package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"asynctaskhub/src/service/queue"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ControllerApiTaskQueue struct {
	ControllerApiBase
}

func (c *ControllerApiTaskQueue) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")

	start := ctx.DefaultQuery("start", "")
	end := ctx.DefaultQuery("end", "")

	query := global.DB.Preload("Task", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "app_id", "name", "task_code")
	}).Model(&model.TaskQueue{})
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN ?", adminInfo.AppIDs)
	}

	if start != "" {
		query = query.Where("created_at >=?", start)
	}

	if end != "" {
		query = query.Where("created_at <=?", end)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取队列列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取队列列表失败", nil)
		return
	}

	var taskQueueLists []model.TaskQueue
	if err := query.Omit("parameters").Order("id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&taskQueueLists).Error; err != nil {
		global.Logger.Warn("获取队列列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取队列列表失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取队列列表成功", gin.H{
		"list":  taskQueueLists,
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

	query := global.DB.Preload("Task")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN?", adminInfo.AppIDs)
	}
	var taskQueue model.TaskQueue
	if err := query.First(&taskQueue, id).Error; err != nil {
		global.Logger.Warn("获取任务队列详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务队列详情失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取任务队列详情成功", &taskQueue)
}

type TaskQueueRequest struct {
	TaskCode           string `json:"task_code" binding:"required"`
	Parameters         string `json:"parameters"`
	RelativeDelayTime  int    `json:"relative_delay_time"`
	DelayExecutionTime int    `json:"delay_execution_time"`
}

func (c *ControllerApiTaskQueue) Create(ctx *gin.Context) {
	var input TaskQueueRequest

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
	taskQueue.CreatedAt = time.Now()
	taskQueue.UpdatedAt = time.Now()

	if input.RelativeDelayTime != 0 {
		taskQueue.ExecutionTime = int(time.Now().Unix()) + input.RelativeDelayTime
	} else if input.DelayExecutionTime != 0 {
		taskQueue.ExecutionTime = input.DelayExecutionTime
	} else {
		taskQueue.ExecutionTime = int(time.Now().Unix())
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
