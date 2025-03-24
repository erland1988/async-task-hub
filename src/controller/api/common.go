package api

import (
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type ControllerApiCommon struct {
	ControllerApiBase
}

type reponsePie struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

func (c *ControllerApiCommon) Timeline(ctx *gin.Context) {
	var results []struct {
		AppName           string
		TaskName          string
		Status            string
		ExecutionDuration int64
		ExecutionCount    int
		UpdatedAt         time.Time
	}
	cacheKey := global.CacheKey("timeline")
	// 尝试从缓存中获取数据
	if err := global.GetFromCache(cacheKey, &results); err != nil {
		global.Logger.Warn("获取缓存失败", zap.Error(err))
	}
	if len(results) == 0 {
		query := global.DB.Model(&model.TaskQueue{}).
			Joins("LEFT JOIN tasks ON task_queues.task_id = tasks.id").
			Joins("LEFT JOIN applications ON task_queues.app_id = applications.id").
			Select("tasks.name AS task_name", "applications.name AS app_name", "task_queues.execution_status AS status", "task_queues.execution_duration", "task_queues.execution_count", "task_queues.updated_at").
			//Where("execution_time > ?", time.Now().Add(-72*time.Hour).Unix()).
			Where("execution_status = ?", model.TaskQueueCompleted).
			Order("updated_at DESC").
			Limit(5)

		if err := query.Find(&results).Error; err != nil {
			global.Logger.Warn("获取数据失败", zap.Error(err))
			c.JSONResponse(ctx, false, "获取数据失败", nil)
			return
		}

		if err := global.SetToCache(cacheKey, results, 30*time.Second); err != nil {
			global.Logger.Warn("设置缓存失败", zap.Error(err))
		}
	}

	var timeline []map[string]interface{}
	colors := []string{"#00bcd4", "#1ABC9C", "#3f51b5", "#f44336", "#009688"}
	for i, result := range results {
		desc := fmt.Sprintf("%s 已完成，耗时 %d 毫秒，执行次数 %d 次",
			result.TaskName, result.ExecutionDuration, result.ExecutionCount)

		timeline = append(timeline, map[string]interface{}{
			"content":     result.AppName,
			"description": desc,
			"timestamp":   result.UpdatedAt.Format("2006-01-02 15:04:05"),
			"color":       colors[i%len(colors)],
		})
	}
	c.JSONResponse(ctx, true, "获取数据成功", timeline)
}

func (c *ControllerApiCommon) Pie(ctx *gin.Context) {
	var pies []reponsePie

	cacheKey := global.CacheKey("pie")
	// 尝试从缓存中获取数据
	if err := global.GetFromCache(cacheKey, &pies); err != nil {
		global.Logger.Warn("获取缓存失败", zap.Error(err))
	}
	if len(pies) == 0 {
		//遍历执行状态
		for _, status := range []model.TaskQueueExecutionStatus{model.TaskQueuePending, model.TaskQueueCompleted, model.TaskQueueFailed} {
			var total int64
			if err := global.DB.Model(&model.TaskQueue{}).
				Where("execution_time > ?", time.Now().Add(-720*time.Hour).Unix()).
				Where("execution_status =?", status).
				Count(&total).Error; err != nil {
				global.Logger.Warn("Pie", zap.Error(err))
				c.JSONResponse(ctx, false, "获取数据失败", nil)
				return
			}
			pies = append(pies, reponsePie{
				Value: int(total),
				Name:  string(status),
			})
		}
		if err := global.SetToCache(cacheKey, pies, 2*time.Hour); err != nil {
			global.Logger.Warn("设置缓存失败", zap.Error(err))
		}
	}
	c.JSONResponse(ctx, true, "获取数据成功", pies)
}

func (c *ControllerApiCommon) Line(ctx *gin.Context) {
	lines := make(map[string][]interface{})

	cacheKey := global.CacheKey("line")
	// 尝试从缓存中获取数据
	if err := global.GetFromCache(cacheKey, &lines); err != nil {
		global.Logger.Warn("获取缓存失败", zap.Error(err))
	}
	if len(lines) == 0 {
		now := time.Now()
		for i := 1; i < 8; i++ {
			day := now.AddDate(0, 0, -i).Format("2006-01-02")
			parsedDay, _ := time.ParseInLocation("2006-01-02", day, time.Local)

			lines["day"] = append(lines["day"], day)

			executionStart := parsedDay.Unix()
			executionEnd := parsedDay.Add(24*time.Hour).Unix() - 1
			var totalExecution int64
			if err := global.DB.Model(&model.TaskQueue{}).Where("execution_status = ?", model.TaskQueueCompleted).Where("execution_time >= ?", executionStart).Where("execution_time <= ?", executionEnd).Count(&totalExecution).Error; err != nil {
				global.Logger.Warn("lines", zap.Error(err))
				c.JSONResponse(ctx, false, "获取数据失败", nil)
				return
			}
			lines["execution"] = append(lines["execution"], totalExecution)

			createdStart := parsedDay
			createdEnd := parsedDay.Add(24*time.Hour - time.Second)
			var totalCreated int64
			if err := global.DB.Model(&model.TaskQueue{}).Where("created_at >= ?", createdStart).Where("created_at <= ?", createdEnd).Count(&totalCreated).Error; err != nil {
				global.Logger.Warn("lines", zap.Error(err))
				c.JSONResponse(ctx, false, "获取数据失败", nil)
				return
			}
			lines["created"] = append(lines["created"], totalCreated)
		}
		if err := global.SetToCache(cacheKey, lines, 8*time.Hour); err != nil {
			global.Logger.Warn("设置缓存失败", zap.Error(err))
		}
	}
	c.JSONResponse(ctx, true, "获取数据成功", lines)
}

type responseHomes struct {
	Totals     map[string]int64         `json:"totals"`
	AppQueues  []map[string]interface{} `json:"appqueues"`
	TaskQueues []map[string]interface{} `json:"taskqueues"`
}

func (c *ControllerApiCommon) Home(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	out := responseHomes{}

	cacheKey := global.CacheKey("home")
	// 尝试从缓存中获取数据
	if err := global.GetFromCache(cacheKey, &out); err != nil {
		global.Logger.Warn("获取缓存失败", zap.Error(err))
	}
	// 缓存为空，重新计算数据
	if out.Totals == nil {
		totals, err := c.totals()
		if err != nil {
			c.JSONResponse(ctx, false, "获取数据失败", nil)
			return
		}
		out.Totals = totals

		appqueues, err := c.appDeQueues()
		if err != nil {
			c.JSONResponse(ctx, false, "获取数据失败", nil)
			return
		}
		out.AppQueues = appqueues

		taskqueues, err := c.taskDeQueues()
		if err != nil {
			c.JSONResponse(ctx, false, "获取数据失败", nil)
			return
		}
		out.TaskQueues = taskqueues
		global.SetToCache(cacheKey, out, 5*time.Minute)
	}

	c.JSONResponse(ctx, true, "获取数据成功", out)
}

func (c *ControllerApiCommon) totals() (map[string]int64, error) {
	totals := make(map[string]int64)
	var totalAdmin int64
	if err := global.DB.Model(&model.Admin{}).Count(&totalAdmin).Error; err != nil {
		global.Logger.Warn("totals", zap.Error(err))
		return nil, err
	}
	totals["admin"] = totalAdmin

	var totalApp int64
	if err := global.DB.Model(&model.Application{}).Count(&totalApp).Error; err != nil {
		global.Logger.Warn("totals", zap.Error(err))
		return nil, err
	}
	totals["app"] = totalApp

	var totalTask int64
	if err := global.DB.Model(&model.Task{}).Count(&totalTask).Error; err != nil {
		global.Logger.Warn("totals", zap.Error(err))
		return nil, err
	}
	totals["task"] = totalTask

	var totalQueue int64
	if err := global.DB.Model(&model.TaskQueue{}).Count(&totalQueue).Error; err != nil {
		global.Logger.Warn("totals", zap.Error(err))
		return nil, err
	}
	totals["queue"] = totalQueue
	return totals, nil
}

func (c *ControllerApiCommon) appDeQueues() ([]map[string]interface{}, error) {
	var results []struct {
		AppID int
		Name  string
		Total int64
	}

	query := global.DB.Model(&model.TaskQueue{}).
		Joins("LEFT JOIN applications ON task_queues.app_id = applications.id").
		Select("task_queues.app_id, applications.name, count(task_queues.id) as total").
		Where("execution_status = ?", model.TaskQueueCompleted).
		Group("task_queues.app_id").
		Order("total DESC").
		Limit(5)

	if err := query.Find(&results).Error; err != nil {
		global.Logger.Warn("appDeQueues", zap.Error(err))
		return nil, err
	}

	colors := []string{"#f25e43", "#00bcd4", "#64d572", "#e9a745", "#009688"}
	var ranks []map[string]interface{}
	var maxTotal int64 = 0
	if len(results) > 0 {
		maxTotal = results[0].Total
	}

	for i, result := range results {
		if result.Name == "" {
			continue
		}
		ranks = append(ranks, map[string]interface{}{
			"title":   result.Name,
			"value":   result.Total,
			"percent": (result.Total * 100) / maxTotal,
			"color":   colors[i%len(colors)],
		})
	}

	return ranks, nil
}

func (c *ControllerApiCommon) taskDeQueues() ([]map[string]interface{}, error) {
	var results []struct {
		TaskID int
		Name   string
		Total  int64
	}

	query := global.DB.Model(&model.TaskQueue{}).
		Joins("LEFT JOIN tasks ON task_queues.task_id = tasks.id").
		Select("task_queues.task_id, tasks.name, count(task_queues.id) as total").
		Where("execution_status = ?", model.TaskQueueCompleted).
		Group("task_queues.task_id").
		Order("total DESC").
		Limit(5)

	if err := query.Find(&results).Error; err != nil {
		global.Logger.Warn("taskDeQueues", zap.Error(err))
		return nil, err
	}

	colors := []string{"#f25e43", "#00bcd4", "#64d572", "#e9a745", "#009688"} // 颜色映射
	var ranks []map[string]interface{}
	var maxTotal int64 = 0
	if len(results) > 0 {
		maxTotal = results[0].Total // 最大值用于计算百分比
	}

	for i, result := range results {
		if result.Name == "" {
			continue
		}
		ranks = append(ranks, map[string]interface{}{
			"title":   result.Name,
			"value":   result.Total,
			"percent": (result.Total * 100) / maxTotal, // 百分比计算
			"color":   colors[i%len(colors)],
		})
	}

	return ranks, nil
}
