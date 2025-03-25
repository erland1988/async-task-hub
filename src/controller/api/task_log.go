package api

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"async-task-hub/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiTaskLog struct {
	ControllerApiBase
}

type responseTaskLog struct {
	ID             int                  `json:"id"`
	AppID          int                  `json:"app_id"`
	TaskID         int                  `json:"task_id"`
	TaskQueueID    int                  `json:"task_queue_id"`
	RequestID      string               `json:"request_id"`
	Action         model.TaskLogAction  `json:"action"`
	ActionString   string               `json:"action_string"`
	Message        string               `json:"message"`
	MilliTimestamp types.MilliTimestamp `json:"milli_timestamp"`
	CreatedAt      types.Customtime     `json:"created_at"`
	TaskName       string               `json:"task_name"`
	AppName        string               `json:"app_name"`
}

func (c *ControllerApiTaskLog) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")
	start := ctx.DefaultQuery("start", "")
	end := ctx.DefaultQuery("end", "")
	requestId := ctx.DefaultQuery("request_id", "")

	query := global.DB.Model(&model.TaskLog{}).
		Select("task_logs.*,tasks.name as task_name,applications.name as app_name").
		Joins("LEFT JOIN tasks ON tasks.id = task_logs.task_id").
		Joins("LEFT JOIN applications ON applications.id = task_logs.app_id")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("task_logs.app_id IN?", adminInfo.AppIDs)
	}

	if start != "" {
		query = query.Where("task_logs.created_at >=?", start)
	}
	if end != "" {
		query = query.Where("task_logs.created_at <=?", end)
	}
	if requestId != "" {
		query = query.Where("task_logs.request_id =?", requestId)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取任务日志列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志列表失败", nil)
		return
	}

	var responseTaskLogs []responseTaskLog
	if err := query.Order("task_logs.id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&responseTaskLogs).Error; err != nil {
		global.Logger.Warn("获取任务日志列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志列表失败", nil)
		return
	}

	for i, taskLog := range responseTaskLogs {
		responseTaskLogs[i].ActionString = taskLog.Action.String()
	}
	c.JSONResponse(ctx, true, "获取任务日志列表成功", gin.H{
		"list":  responseTaskLogs,
		"total": total,
	})
}

func (c *ControllerApiTaskLog) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "获取任务日志详情失败", nil)
		return
	}

	query := global.DB.Model(&model.TaskLog{}).
		Select("task_logs.*,tasks.name as task_name,applications.name as app_name").
		Joins("LEFT JOIN tasks ON tasks.id = task_logs.task_id").
		Joins("LEFT JOIN applications ON applications.id = task_logs.app_id")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("task_logs.app_id IN?", adminInfo.AppIDs)
	}
	var responseTaskLog responseTaskLog
	if err := query.First(&responseTaskLog, id).Error; err != nil {
		global.Logger.Warn("获取任务日志详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志详情失败", nil)
		return
	}
	responseTaskLog.ActionString = responseTaskLog.Action.String()

	global.Logger.Info("获取任务日志详情成功", zap.Any("responseTaskLog", &responseTaskLog))
	c.JSONResponse(ctx, true, "获取任务日志详情成功", &responseTaskLog)
}
