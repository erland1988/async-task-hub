package router

import (
	"async-task-hub/global"
	"github.com/gin-gonic/gin"
)

type RouterWeb struct {
}

func (router *RouterWeb) LoadRouter(r *gin.Engine) *gin.Engine {
	basePath := global.CONFIG.BASE_PATH
	webGroup := r.Group(basePath)
	{
		webGroup.Static("/static", "./static")
		webGroup.StaticFile("/favicon.ico", "./static/favicon.ico")
		webGroup.StaticFile("/robots.txt", "./static/robots.txt")

		webGroup.GET("/", func(c *gin.Context) {
			c.File("./public/backend/index.html")
		})

		webGroup.Static("/assets", "./public/backend/assets")
	}

	return r
}

func init() {
	global.RegisterRouter(&RouterWeb{})
}
