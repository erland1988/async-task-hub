package middleware

import (
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type AdminInfo struct {
	ID     int             `json:"id"`
	Role   model.AdminRole `json:"role"`
	AppIDs []int           `json:"app_ids,omitempty"` // 只有 app_admin 才有此字段
}

type AppInfo struct {
	ID int `json:"id"`
}

type AuthInfo struct {
	IsApp   bool `json:"is_app"`
	IsAdmin bool `json:"is_admin"`
	AdminInfo
	AppInfo
}

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authInfo AuthInfo
		token := c.GetHeader("Authorization")
		if token == "" {
			appKey := c.GetHeader("X-App-Key")
			appSecret := c.GetHeader("X-App-Secret")
			if appKey == "" || appSecret == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的app凭证"})
				c.Abort()
				return
			}

			var app model.Application
			if err := global.DB.Where("app_key = ? AND app_secret = ?", appKey, appSecret).First(&app).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的app凭证"})
				c.Abort()
				return
			}

			var admin model.Admin
			if err := global.DB.Where("id =?", app.AdminID).First(&admin).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的app凭证"})
				c.Abort()
				return
			}

			if time.Time(admin.ExpiresAt).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的app凭证"})
				c.Abort()
				return
			}

			authInfo.IsApp = true
			authInfo.AppInfo.ID = app.ID
			c.Set("auth_info", authInfo)
			c.Next()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		var login model.Login
		if err := global.DB.Where("token = ?", token).Order("expires_at desc").First(&login).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或者过期的token"})
			c.Abort()
			return
		}

		if time.Time(login.ExpiresAt).Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或者过期的token"})
			c.Abort()
			return
		}

		var admin model.Admin
		if err := global.DB.Where("id = ?", login.AdminID).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或者过期的token"})
			c.Abort()
			return
		}

		if time.Time(admin.ExpiresAt).Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或者过期的token"})
			c.Abort()
			return
		}

		var appIDs []int
		if admin.Role == model.AppAdmin {
			var apps []model.Application
			if err := global.DB.Where("admin_id =?", admin.ID).Find(&apps).Error; err != nil {
				global.Logger.Warn(err.Error(), zap.Error(err))
			}
			for _, app := range apps {
				appIDs = append(appIDs, app.ID)
			}
		}

		authInfo.IsAdmin = true
		authInfo.AdminInfo.ID = admin.ID
		authInfo.AdminInfo.Role = admin.Role
		authInfo.AdminInfo.AppIDs = appIDs
		c.Set("auth_info", authInfo)
		c.Next()
	}
}
