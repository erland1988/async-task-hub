package router

import (
	"async-task-hub/global"
	"async-task-hub/src/controller/api"
	"async-task-hub/src/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

type RouterApi struct {
}

func (router *RouterApi) LoadRouter(r *gin.Engine) *gin.Engine {
	noAuthApiGroup := r.Group("/api")
	{
		controllerApiAdmin := &api.ControllerApiAdmin{}
		noAuthApiGroup.POST("admin/registry", middleware.RateLimiterMiddleware(middleware.NewRateLimiter(3, time.Minute*10)), controllerApiAdmin.Registry)
		noAuthApiGroup.POST("admin/login", middleware.RateLimiterMiddleware(middleware.NewRateLimiter(3, time.Minute)), controllerApiAdmin.Login)
	}
	apiGroup := r.Group("/api", middleware.LoginMiddleware())
	{

		controllerApiCommon := &api.ControllerApiCommon{}
		apiGroup.GET("common/home", controllerApiCommon.Home)
		apiGroup.GET("common/line", controllerApiCommon.Line)
		apiGroup.GET("common/pie", controllerApiCommon.Pie)
		apiGroup.GET("common/timeline", controllerApiCommon.Timeline)

		controllerApiAdmin := &api.ControllerApiAdmin{}
		apiGroup.POST("admin/loginout", controllerApiAdmin.Loginout)
		apiGroup.GET("admin/getList", controllerApiAdmin.GetList)
		apiGroup.GET("admin/getDetail", controllerApiAdmin.GetDetail)
		apiGroup.POST("admin/create", controllerApiAdmin.Create)
		apiGroup.POST("admin/update", controllerApiAdmin.Update)
		apiGroup.POST("admin/delete", controllerApiAdmin.Delete)

		apiGroup.POST("admin/resetPassword", controllerApiAdmin.ResetPassword)
		apiGroup.POST("admin/updateProfile", controllerApiAdmin.UpdateProfile)

		controllerApiLog := &api.ControllerApiLog{}
		apiGroup.GET("log/getList", controllerApiLog.GetList)
		apiGroup.GET("log/getDetail", controllerApiLog.GetDetail)

		controllerApiApplication := &api.ControllerApiApplication{}
		apiGroup.GET("app/getList", controllerApiApplication.GetList)
		apiGroup.GET("app/getDetail", controllerApiApplication.GetDetail)
		apiGroup.POST("app/create", controllerApiApplication.Create)
		apiGroup.POST("app/update", controllerApiApplication.Update)
		apiGroup.POST("app/delete", controllerApiApplication.Delete)

		controllerApiTask := &api.ControllerApiTask{}
		apiGroup.GET("task/getList", controllerApiTask.GetList)
		apiGroup.GET("task/getDetail", controllerApiTask.GetDetail)
		apiGroup.POST("task/create", controllerApiTask.Create)
		apiGroup.POST("task/update", controllerApiTask.Update)
		apiGroup.POST("task/delete", controllerApiTask.Delete)

		controllerApiTaskQueue := &api.ControllerApiTaskQueue{}
		apiGroup.GET("taskqueue/getList", controllerApiTaskQueue.GetList)
		apiGroup.GET("taskqueue/getDetail", controllerApiTaskQueue.GetDetail)
		apiGroup.POST("taskqueue/create", controllerApiTaskQueue.Create)

		controllerApiTaskLog := &api.ControllerApiTaskLog{}
		apiGroup.GET("tasklog/getList", controllerApiTaskLog.GetList)
		apiGroup.GET("tasklog/getDetail", controllerApiTaskLog.GetDetail)

		controllerApiConfig := &api.ControllerApiConfig{}
		apiGroup.GET("config/getConfigs", controllerApiConfig.GetConfigs)
		apiGroup.GET("config/getCustomerConfigs", controllerApiConfig.GetCustomerConfigs)
		apiGroup.POST("config/updateConfigs", controllerApiConfig.UpdateConfigs)
	}
	return r
}

func init() {
	global.RegisterRouter(&RouterApi{})
}
