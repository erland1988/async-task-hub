package router

import (
	"async-task-hub/global"
	"async-task-hub/src/controller/web"
	"github.com/gin-gonic/gin"
	"strings"
)

type RouterWeb struct {
}

func (router *RouterWeb) LoadRouter(r *gin.Engine) *gin.Engine {
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/robots.txt", "./static/robots.txt")

	r.GET("/", web.ControllerWebHome{}.Index)

	r.Static("/backend/assets", "./public/backend/assets")
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/backend" {
			c.Redirect(301, "/backend/")
			return
		}
		if strings.HasPrefix(path, "/backend") {
			c.File("./public/backend/index.html")
		}
	})
	return r
}

func init() {
	global.RegisterRouter(&RouterWeb{})
}
