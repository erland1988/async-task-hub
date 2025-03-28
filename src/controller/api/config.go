package api

import (
	"async-task-hub/global"
	"async-task-hub/src/model"
	"async-task-hub/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiConfig struct {
	ControllerApiBase
}

func (c *ControllerApiConfig) GetConfigs(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	configMap, err := service.NewConfigService().GetConfigMap()
	if err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取配置失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取配置成功", configMap)
}

func (c *ControllerApiConfig) GetCustomerConfigs(ctx *gin.Context) {
	_ = c.CheckAdmin(ctx)

	configMap := make(map[string]string)
	keys := []string{
		"notice",
	}
	for _, key := range keys {
		value, err := service.NewConfigService().GetConfig(key)
		if err != nil {
			value = ""
		}
		configMap[key] = value
	}
	c.JSONResponse(ctx, true, "获取配置成功", configMap)
}

func (c *ControllerApiConfig) UpdateConfigs(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	var configs map[string]string
	if err := ctx.ShouldBindJSON(&configs); err != nil {
		c.JSONResponse(ctx, false, "参数解析失败", nil)
		return
	}

	for key, value := range configs {
		if err := service.NewConfigService().UpdateConfig(key, value); err != nil {
			global.Logger.Warn("更新配置失败", zap.Error(err))
			c.JSONResponse(ctx, false, "更新配置失败", nil)
			return
		}
	}

	service.NewLogService().CreateLog(adminInfo.ID, "修改配置", configs)
	c.JSONResponse(ctx, true, "修改配置成功", nil)
}
