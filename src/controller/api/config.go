package api

import (
	"asynctaskhub/global"
	"asynctaskhub/src/model"
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

	var configs []model.Config
	if err := global.DB.Find(&configs).Error; err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取配置失败", nil)
		return
	}
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	c.JSONResponse(ctx, true, "获取配置成功", configMap)
}

func (c *ControllerApiConfig) GetCustomerConfigs(ctx *gin.Context) {
	_ = c.CheckAdmin(ctx)

	keys := []string{
		"notice",
	}
	var configs []model.Config
	if err := global.DB.Where("`key` IN (?)", keys).Find(&configs).Error; err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取配置失败", nil)
		return
	}
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	c.JSONResponse(ctx, true, "获取配置成功", configMap)
}

func (c *ControllerApiConfig) UpdateConfigs(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	var configs []model.Config
	if err := ctx.ShouldBindJSON(&configs); err != nil {
		c.JSONResponse(ctx, false, "参数解析失败", nil)
		return
	}

	for _, config := range configs {
		if err := global.DB.Where("`key` = ?", config.Key).Updates(&config).Error; err != nil {
			global.Logger.Warn("修改配置失败", zap.Error(err))
			c.JSONResponse(ctx, false, "修改配置失败", nil)
			return
		}
	}

	c.JSONResponse(ctx, true, "修改配置成功", nil)
}
