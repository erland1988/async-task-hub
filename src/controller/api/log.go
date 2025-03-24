package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiLog struct {
	ControllerApiBase
}

func (c *ControllerApiLog) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")

	keywords := ctx.DefaultQuery("keywords", "")

	query := global.DB.Preload("Admin").Model(&model.Log{}).Joins("LEFT JOIN admins ON logs.admin_id = admins.id")
	if keywords != "" {
		query = query.Where("(logs.operation LIKE ?)", "%"+keywords+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn(err.Error(), zap.Error(err))
		c.JSONResponse(ctx, false, "获取操作日志失败", nil)
		return
	}

	var logs []model.Log
	if err := query.Omit("logs.details").Order("logs.id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&logs).Error; err != nil {
		global.Logger.Warn(err.Error(), zap.Error(err))
		c.JSONResponse(ctx, false, "获取操作日志失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取操作日志成功", gin.H{
		"list":  logs,
		"total": total,
	})
}

func (c *ControllerApiLog) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	query := global.DB.Preload("Admin")
	var log model.Log
	if err := query.First(&log, id).Error; err != nil {
		global.Logger.Warn("获取操作日志失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取操作日志失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取操作日志成功", &log)
}
