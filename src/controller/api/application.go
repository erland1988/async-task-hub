package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"asynctaskhub/src/service"
	"asynctaskhub/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiApplication struct {
	ControllerApiBase
}
type requestCreateApplication struct {
	Name      string `json:"name" binding:"required"`       // 应用名称
	AppKey    string `json:"app_key" binding:"required"`    // 应用标识
	AppSecret string `json:"app_secret" binding:"required"` // 应用秘钥
	Remark    string `json:"remark"`                        // 备注
}

type requestUpdateApplication struct {
	ID        int    `json:"id" binding:"required"`         // 应用ID
	Name      string `json:"name" binding:"required"`       // 应用名称
	AppKey    string `json:"app_key" binding:"required"`    // 应用标识
	AppSecret string `json:"app_secret" binding:"required"` // 应用秘钥
	Remark    string `json:"remark"`                        // 备注
}

type responseApplication struct {
	ID        int              `json:"id"`         // 应用ID
	Name      string           `json:"name"`       // 应用名称
	AppKey    string           `json:"app_key"`    // 应用标识
	AppSecret string           `json:"app_secret"` // 应用秘钥
	Remark    string           `json:"remark"`     // 备注
	CreatedAt types.Customtime `json:"created_at"` // 创建时间
	UpdatedAt types.Customtime `json:"updated_at"` // 更新时间
	Username  string           `json:"username"`   // 更新时间
}

func (c *ControllerApiApplication) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")
	keywords := ctx.DefaultQuery("keywords", "")

	query := global.DB.Model(&model.Application{}).
		Joins("LEFT JOIN admins ON applications.admin_id = admins.id").
		Select("applications.*, admins.username")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("applications.admin_id = ?", adminInfo.ID)
	}

	if keywords != "" {
		query = query.Where("(applications.name like ? or applications.app_key like ?)", "%"+keywords+"%", "%"+keywords+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取应用列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取应用列表失败", nil)
		return
	}

	var responseApplications []responseApplication
	if err := query.Order("applications.id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&responseApplications).Error; err != nil {
		global.Logger.Warn("获取应用列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取应用列表失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取应用列表成功", gin.H{
		"list":  responseApplications,
		"total": total,
	})
}

func (c *ControllerApiApplication) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}

	query := global.DB.Model(&model.Application{}).
		Joins("LEFT JOIN admins ON applications.admin_id = admins.id").
		Select("applications.*, admins.username")
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("applications.admin_id =?", adminInfo.ID)
	}
	var responseApplication responseApplication
	if err := query.First(&responseApplication, id).Error; err != nil {
		global.Logger.Warn("获取应用详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取应用详情失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取应用详情成功", &responseApplication)
}

func (c *ControllerApiApplication) Create(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestCreateApplication
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数错误", zap.Error(err))
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}
	// 校验应用标识是否重复
	var count int64
	if err := global.DB.Model(&model.Application{}).Where("app_key =?", input.AppKey).Count(&count).Error; err != nil {
		global.Logger.Warn("创建应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建应用失败", nil)
		return
	}
	if count > 0 {
		c.JSONResponse(ctx, false, "应用标识已存在", nil)
		return
	}

	application := model.Application{
		Name:      input.Name,
		AppKey:    input.AppKey,
		AppSecret: input.AppSecret,
		Remark:    input.Remark,
		AdminID:   adminInfo.ID,
	}
	if err := global.DB.Create(&application).Error; err != nil {
		global.Logger.Warn("创建应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建应用失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "创建任务成功", nil)
}

func (c *ControllerApiApplication) Update(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestUpdateApplication
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数错误", zap.Error(err))
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}
	// 校验应用标识是否重复
	var count int64
	if err := global.DB.Model(&model.Application{}).Where("app_key =? and id !=?", input.AppKey, input.ID).Count(&count).Error; err != nil {
		global.Logger.Warn("更新应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新应用失败", nil)
		return
	}
	if count > 0 {
		c.JSONResponse(ctx, false, "应用标识已存在", nil)
		return
	}

	query := global.DB.Model(&model.Application{})
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("admin_id =?", adminInfo.ID)
	}

	var app model.Application
	if err := query.Where("id = ?", input.ID).First(&app).Error; err != nil {
		global.Logger.Warn("更新应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新应用失败", nil)
		return
	}

	app.Name = input.Name
	app.AppKey = input.AppKey
	app.AppSecret = input.AppSecret
	app.Remark = input.Remark
	if err := query.Omit("created_at", "admin_id").Save(&app).Error; err != nil {
		global.Logger.Warn("更新应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新应用失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "更新应用成功", nil)
}

func (c *ControllerApiApplication) Delete(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.PostForm("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}

	if adminInfo.Role != model.GlobalAdmin {
		if false == common.InArray(id, adminInfo.AppIDs) {
			c.JSONResponse(ctx, false, "未授权访问", nil)
			return
		}
	}

	if err := service.NewTaskService().DeleteApp(id); err != nil {
		global.Logger.Warn("删除应用失败", zap.Error(err))
		c.JSONResponse(ctx, false, "删除应用失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "删除应用成功", nil)
}
