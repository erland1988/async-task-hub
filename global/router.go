package global

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Router interface {
	LoadRouter(r *gin.Engine) *gin.Engine
}

var (
	routers []Router
	once    sync.Once
)

func RegisterRouter(router Router) {
	routers = append(routers, router)
}

func InitRouter(r *gin.Engine) *gin.Engine {
	once.Do(func() {
		for _, router := range routers {
			r = router.LoadRouter(r)
		}
	})
	return r
}
