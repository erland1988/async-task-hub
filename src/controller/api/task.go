package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"asynctaskhub/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerApiTask struct {
	ControllerApiBase
}

type requestCreateTask struct {
	Name        string `json:"name" bingding:"required"`         // 任务名称
	AppID       int    `json:"app_id" bingding:"required"`       // 应用ID
	TaskCode    string `json:"task_code" bingding:"required"`    // 任务编码
	ExecutorURL string `json:"executor_url" bingding:"required"` // 执行器地址
}

type requestUpdateTask struct {
	ID          int    `json:"id" bingding:"required"`           // 任务ID
	Name        string `json:"name" bingding:"required"`         // 任务名称
	TaskCode    string `json:"task_code" bingding:"required"`    // 任务编码
	ExecutorURL string `json:"executor_url" bingding:"required"` // 执行器地址
}

func (c *ControllerApiTask) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	page, pageSize := c.GetPaginationParams(ctx, "page", "pageSize")
	keywords := ctx.DefaultQuery("keywords", "")

	query := global.DB.Model(&model.Task{})
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN ?", adminInfo.AppIDs)
	}
	if keywords != "" {
		query = query.Where("(name LIKE ? or task_code LIKE ?)", "%"+keywords+"%", "%"+keywords+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取任务列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务列表失败", nil)
		return
	}

	var taskLists []model.Task
	if err := query.Order("id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&taskLists).Error; err != nil {
		global.Logger.Warn("获取任务列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务列表失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取任务列表成功", gin.H{
		"list":  taskLists,
		"total": total,
	})
}

func (c *ControllerApiTask) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	query := global.DB
	if adminInfo.Role != model.GlobalAdmin {
		query = query.Where("app_id IN (?)", adminInfo.AppIDs)
	}
	var task model.Task
	if err := query.First(&task, id).Error; err != nil {
		global.Logger.Warn("获取任务详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取任务详情失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取任务详情成功", &task)
}

func (c *ControllerApiTask) Create(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestCreateTask
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}
	if adminInfo.Role != model.GlobalAdmin {
		if common.InArray(input.AppID, adminInfo.AppIDs) == false {
			c.JSONResponse(ctx, false, "未授权访问", nil)
			return
		}
	}

	// 检查任务编码是否重复
	var count int64
	if err := global.DB.Model(&model.Task{}).Where("task_code =?", input.TaskCode).Count(&count).Error; err != nil {
		global.Logger.Warn("创建任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建任务失败", nil)
		return
	}
	if count > 0 {
		c.JSONResponse(ctx, false, "任务编码已存在", nil)
		return
	}

	task := model.Task{
		AppID:       input.AppID,
		TaskCode:    input.TaskCode,
		ExecutorURL: input.ExecutorURL,
		Name:        input.Name,
	}
	if err := global.DB.Create(&task).Error; err != nil {
		global.Logger.Warn("创建任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建任务失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "创建任务成功", nil)
}

func (c *ControllerApiTask) Update(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestUpdateTask
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	//检查任务编码是否重复
	var count int64
	if err := global.DB.Model(&model.Task{}).Where("task_code =? and id !=?", input.TaskCode, input.ID).Count(&count).Error; err != nil {
		global.Logger.Warn("更新任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新任务失败", nil)
		return
	}
	if count > 0 {
		c.JSONResponse(ctx, false, "任务编码已存在", nil)
		return
	}

	var task model.Task
	if err := global.DB.Where("id = ?", input.ID).First(&task).Error; err != nil {
		global.Logger.Warn("更新任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新任务失败", nil)
		return
	}

	if adminInfo.Role != model.GlobalAdmin {
		if false == common.InArray(task.AppID, adminInfo.AppIDs) {
			c.JSONResponse(ctx, false, "未授权访问", nil)
			return
		}
	}

	task.Name = input.Name
	task.TaskCode = input.TaskCode
	task.ExecutorURL = input.ExecutorURL
	if err := global.DB.Omit("app_id", "created_at").Save(&task).Error; err != nil {
		global.Logger.Warn("更新任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新任务失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "更新任务成功", nil)
}

func (c *ControllerApiTask) Delete(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	id := common.Str2Int(ctx.PostForm("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数错误", nil)
		return
	}

	var task model.Task
	if err := global.DB.First(&task, id).Error; err != nil {
		global.Logger.Warn("删除任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "删除任务失败", nil)
		return
	}

	if adminInfo.Role != model.GlobalAdmin {
		if false == common.InArray(task.AppID, adminInfo.AppIDs) {
			c.JSONResponse(ctx, false, "未授权访问", nil)
			return
		}
	}

	if err := service.NewTaskService().DeleteTask(id); err != nil {
		global.Logger.Warn("删除任务失败", zap.Error(err))
		c.JSONResponse(ctx, false, "删除任务失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "删除任务成功", nil)
}
