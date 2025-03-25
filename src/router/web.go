package router

import (
	"async-task-hub/global"
	"github.com/gin-gonic/gin"
)

type RouterWeb struct {
}

func (router *RouterWeb) LoadRouter(r *gin.Engine) *gin.Engine {
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/robots.txt", "./static/robots.txt")

	r.GET("/", func(c *gin.Context) {
		c.File("./public/backend/index.html")
	})

	r.Static("/assets", "./public/backend/assets")

	return r
}

func init() {
	global.RegisterRouter(&RouterWeb{})
}
