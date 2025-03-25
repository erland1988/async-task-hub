package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerWebHome struct {
	ControllerWebBase
}

func (controller ControllerWebHome) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "web/index.html", gin.H{})
}
