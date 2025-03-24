package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiTaskLog struct {
	ControllerApiBase
}

func (c *ControllerApiTaskLog) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")
	start := ctx.DefaultQuery("start", "")
	end := ctx.DefaultQuery("end", "")
	requestId := ctx.DefaultQuery("request_id", "")

	query := global.DB.Model(&model.TaskLog{})
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN?", adminInfo.AppIDs)
	}

	if start != "" {
		query = query.Where("created_at >=?", start)
	}
	if end != "" {
		query = query.Where("created_at <=?", end)
	}
	if requestId != "" {
		query = query.Where("request_id =?", requestId)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取任务日志列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志列表失败", nil)
		return
	}

	var taskLogLists []model.TaskLog
	if err := query.Order("id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&taskLogLists).Error; err != nil {
		global.Logger.Warn("获取任务日志列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志列表失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取任务日志列表成功", gin.H{
		"list":  taskLogLists,
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

	query := global.DB
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN?", adminInfo.AppIDs)
	}
	var taskLog model.TaskLog
	if err := query.First(&taskLog, id).Error; err != nil {
		global.Logger.Warn("获取任务日志详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务日志详情失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取任务日志详情成功", taskLog)
}
