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
	basePath := global.CONFIG.BASE_PATH
	apiGroup := r.Group(basePath)
	{
		noAuthGroup := apiGroup.Group("/api")
		{
			controllerApiAdmin := &api.ControllerApiAdmin{}
			noAuthGroup.POST("admin/register", middleware.RateLimiterMiddleware(middleware.NewRateLimiter(3, time.Minute*10)), controllerApiAdmin.Register)
			noAuthGroup.POST("admin/login", middleware.RateLimiterMiddleware(middleware.NewRateLimiter(3, time.Minute)), controllerApiAdmin.Login)
		}
		authGroup := apiGroup.Group("/api", middleware.LoginMiddleware())
		{

			controllerApiCommon := &api.ControllerApiCommon{}
			authGroup.GET("common/home", controllerApiCommon.Home)
			authGroup.GET("common/line", controllerApiCommon.Line)
			authGroup.GET("common/pie", controllerApiCommon.Pie)
			authGroup.GET("common/timeline", controllerApiCommon.Timeline)

			controllerApiAdmin := &api.ControllerApiAdmin{}
			authGroup.POST("admin/loginout", controllerApiAdmin.Loginout)
			authGroup.GET("admin/getList", controllerApiAdmin.GetList)
			authGroup.GET("admin/getDetail", controllerApiAdmin.GetDetail)
			authGroup.POST("admin/create", controllerApiAdmin.Create)
			authGroup.POST("admin/update", controllerApiAdmin.Update)
			authGroup.POST("admin/delete", controllerApiAdmin.Delete)

			authGroup.POST("admin/resetPassword", controllerApiAdmin.ResetPassword)
			authGroup.POST("admin/updateProfile", controllerApiAdmin.UpdateProfile)

			controllerApiLog := &api.ControllerApiLog{}
			authGroup.GET("log/getList", controllerApiLog.GetList)
			authGroup.GET("log/getDetail", controllerApiLog.GetDetail)

			controllerApiApplication := &api.ControllerApiApplication{}
			authGroup.GET("app/getList", controllerApiApplication.GetList)
			authGroup.GET("app/getDetail", controllerApiApplication.GetDetail)
			authGroup.POST("app/create", controllerApiApplication.Create)
			authGroup.POST("app/update", controllerApiApplication.Update)
			authGroup.POST("app/delete", controllerApiApplication.Delete)

			controllerApiTask := &api.ControllerApiTask{}
			authGroup.GET("task/getList", controllerApiTask.GetList)
			authGroup.GET("task/getDetail", controllerApiTask.GetDetail)
			authGroup.POST("task/create", controllerApiTask.Create)
			authGroup.POST("task/update", controllerApiTask.Update)
			authGroup.POST("task/delete", controllerApiTask.Delete)

			controllerApiTaskQueue := &api.ControllerApiTaskQueue{}
			authGroup.GET("taskqueue/getList", controllerApiTaskQueue.GetList)
			authGroup.GET("taskqueue/getDetail", controllerApiTaskQueue.GetDetail)
			authGroup.POST("taskqueue/create", controllerApiTaskQueue.Create)

			controllerApiTaskLog := &api.ControllerApiTaskLog{}
			authGroup.GET("tasklog/getList", controllerApiTaskLog.GetList)
			authGroup.GET("tasklog/getDetail", controllerApiTaskLog.GetDetail)

			controllerApiConfig := &api.ControllerApiConfig{}
			authGroup.GET("config/getConfigs", controllerApiConfig.GetConfigs)
			authGroup.GET("config/getCustomerConfigs", controllerApiConfig.GetCustomerConfigs)
			authGroup.POST("config/updateConfigs", controllerApiConfig.UpdateConfigs)
		}
	}
	return r
}

func init() {
	global.RegisterRouter(&RouterApi{})
}
