package api

import (
	"asynctaskhub/common"
	"asynctaskhub/src/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
type ControllerApiBase struct {
}

func (controller *ControllerApiBase) getAuthInfo(ctx *gin.Context) *middleware.AuthInfo {
	info, exists := ctx.Get("auth_info")
	if !exists {
		return nil
	}

	authInfo, ok := info.(middleware.AuthInfo)
	if !ok {
		return nil
	}

	return &authInfo
}

// CheckApp /**
//   - @Description: 检查应用是否登录
//   - @receiver controller
//   - @param ctx
//   - @return *middleware.AppInfo
//     */
func (controller *ControllerApiBase) CheckApp(ctx *gin.Context) *middleware.AppInfo {
	authInfo := controller.getAuthInfo(ctx)
	if authInfo == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	if !authInfo.IsApp {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	if authInfo.AppInfo.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	return &authInfo.AppInfo
}

// CheckAdmin /**
//   - @Description: 检查管理员是否登录
//   - @receiver controller
//   - @param ctx
//   - @return *middleware.AdminInfo
//     */
func (controller *ControllerApiBase) CheckAdmin(ctx *gin.Context) *middleware.AdminInfo {
	authInfo := controller.getAuthInfo(ctx)
	if authInfo == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	if !authInfo.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	if authInfo.AdminInfo.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		ctx.Abort()
		return nil
	}
	return &authInfo.AdminInfo
}

func (controller *ControllerApiBase) GetPaginationParams(ctx *gin.Context, pageKey, pageSizeKey string) (page, pageSize int) {
	page = common.Str2Int(ctx.DefaultQuery(pageKey, "1"))
	pageSize = common.Str2Int(ctx.DefaultQuery(pageSizeKey, "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

func (controller *ControllerApiBase) JSONResponse(ctx *gin.Context, success bool, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}
