package api

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"asynctaskhub/src/service"
	"asynctaskhub/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type ControllerApiAdmin struct {
	ControllerApiBase
}

type requestCreateAdmin struct {
	Username  string          `json:"username" binding:"required"`   // 用户名
	Password  string          `json:"password" binding:"required"`   // 密码
	Truename  string          `json:"truename" binding:"required"`   // 真实姓名
	Phone     string          `json:"phone"`                         // 手机号
	Email     string          `json:"email"`                         // 邮箱
	Role      model.AdminRole `json:"role"     binding:"required"`   // 角色
	ExpiresAt string          `json:"expires_at" binding:"required"` // 到期时间
}

type requestUpdateAdmin struct {
	ID        int             `json:"id" binding:"required"`       // ID
	Username  string          `json:"username" binding:"required"` // 用户名
	Password  string          `json:"password"`                    // 密码
	Truename  string          `json:"truename"`                    // 真实姓名
	Phone     string          `json:"phone"`                       // 手机号
	Email     string          `json:"email"`                       // 邮箱
	Role      model.AdminRole `json:"role"`                        // 角色
	ExpiresAt string          `json:"expires_at"`                  // 到期时间
}

type requestUpdateProfile struct {
	Truename string `json:"truename" binding:"required"` // 真实姓名
	Phone    string `json:"phone"`                       // 手机号
	Email    string `json:"email"`                       // 邮箱
}

type requestResetPassword struct {
	OldPassword     string `json:"old_password" binding:"required"`     // 旧密码
	NewPassword     string `json:"new_password" binding:"required"`     // 新密码
	ConfirmPassword string `json:"confirm_password" binding:"required"` // 确认密码
}

func (c *ControllerApiAdmin) GetList(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	page, pageSize := c.GetPaginationParams(ctx, "page", "page_size")
	keywords := ctx.DefaultQuery("keywords", "")

	query := global.DB.Model(&model.Admin{})
	if keywords != "" {
		query = query.Where("(truename LIKE ? or username LIKE ?)", "%"+keywords+"%", "%"+keywords+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Warn("获取管理员列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取管理员列表失败", nil)
		return
	}

	var adminLists []model.Admin
	if err := query.Omit("password").Order("id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&adminLists).Error; err != nil {
		global.Logger.Warn("获取管理员列表失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取管理员列表失败", nil)
		return
	}

	c.JSONResponse(ctx, true, "获取管理员列表成功", gin.H{
		"list":  adminLists,
		"total": total,
	})
}

func (c *ControllerApiAdmin) GetDetail(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	id := common.Str2Int(ctx.Query("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	query := global.DB
	var admin model.Admin
	if err := query.Omit("password").First(&admin, id).Error; err != nil {
		global.Logger.Warn("获取管理员详情失败", zap.Error(err))
		c.JSONResponse(ctx, false, "获取管理员详情失败", nil)
		return
	}
	c.JSONResponse(ctx, true, "获取管理员详情成功", &admin)
}

func (c *ControllerApiAdmin) Create(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	var input requestCreateAdmin
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	if input.Email != "" {
		if err := common.ValidateEmail(input.Email); err != nil {
			c.JSONResponse(ctx, false, "邮箱格式错误", nil)
			return
		}
	}
	admin := model.Admin{
		Username:  input.Username,
		Password:  common.HashMD5(input.Password),
		Truename:  input.Truename,
		Phone:     input.Phone,
		Email:     input.Email,
		Role:      input.Role,
		ExpiresAt: common.FormatDatetime(input.ExpiresAt),
	}
	if err := global.DB.Create(&admin).Error; err != nil {
		global.Logger.Warn("创建管理员失败", zap.Error(err))
		c.JSONResponse(ctx, false, "创建管理员失败", nil)
		return
	}
	service.NewLogService().CreateLog(adminInfo.ID, "创建管理员", admin)
	c.JSONResponse(ctx, true, "创建管理员成功", nil)
}

func (c *ControllerApiAdmin) UpdateProfile(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestUpdateProfile
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}
	var admin model.Admin
	if err := global.DB.First(&admin, adminInfo.ID).Error; err != nil {
		global.Logger.Warn("更新用户信息失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新用户信息失败", nil)
		return
	}

	if input.Email != "" {
		if err := common.ValidateEmail(input.Email); err != nil {
			c.JSONResponse(ctx, false, "邮箱格式错误", nil)
			return
		}
	}
	admin.Truename = input.Truename
	admin.Phone = input.Phone
	admin.Email = input.Email
	if err := global.DB.Select("truename", "phone", "email", "updated_at").Save(&admin).Error; err != nil {
		global.Logger.Warn("更新用户信息失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新用户信息失败", nil)
		return
	}
	service.NewLogService().CreateLog(adminInfo.ID, "更新用户信息", admin)
	c.JSONResponse(ctx, true, "更新用户信息成功", nil)
}

func (c *ControllerApiAdmin) ResetPassword(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var input requestResetPassword
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}
	if input.NewPassword != input.ConfirmPassword {
		c.JSONResponse(ctx, false, "两次密码不一致", nil)
		return
	}
	var admin model.Admin
	if err := global.DB.First(&admin, adminInfo.ID).Error; err != nil {
		global.Logger.Warn("更新密码失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新密码失败", nil)
		return
	}
	if admin.Password != common.HashMD5(input.OldPassword) {
		c.JSONResponse(ctx, false, "旧密码错误", nil)
		return
	}
	admin.Password = common.HashMD5(input.NewPassword)
	if err := global.DB.Select("password", "updated_at").Save(&admin).Error; err != nil {
		global.Logger.Warn("更新密码失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新密码失败", nil)
		return
	}
	service.NewLogService().CreateLog(adminInfo.ID, "更新密码", admin)
	c.JSONResponse(ctx, true, "更新密码成功", nil)
}

func (c *ControllerApiAdmin) Update(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	var input requestUpdateAdmin
	if err := ctx.ShouldBindJSON(&input); err != nil {
		global.Logger.Warn("参数异常", zap.Error(err))
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	var admin model.Admin
	if err := global.DB.Where("id = ?", input.ID).First(&admin).Error; err != nil {
		c.JSONResponse(ctx, false, "更新管理员失败", nil)
		return
	}
	if input.Email != "" {
		if err := common.ValidateEmail(input.Email); err != nil {
			c.JSONResponse(ctx, false, "邮箱格式错误", nil)
			return
		}
	}
	admin.Username = input.Username
	if input.Password != "" {
		admin.Password = common.HashMD5(input.Password)
	}
	admin.Truename = input.Truename
	admin.Phone = input.Phone
	admin.Email = input.Email
	if input.Role != "" {
		admin.Role = input.Role
	}
	if input.ExpiresAt != "" {
		expiresAt := common.FormatDatetime(input.ExpiresAt)
		admin.ExpiresAt = expiresAt
	}
	if err := global.DB.Omit("created_at").Save(&admin).Error; err != nil {
		global.Logger.Warn("更新管理员失败", zap.Error(err))
		c.JSONResponse(ctx, false, "更新管理员失败", nil)
		return
	}
	service.NewLogService().CreateLog(adminInfo.ID, "更新管理员", admin)
	c.JSONResponse(ctx, true, "更新管理员成功", nil)
}

func (c *ControllerApiAdmin) Delete(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)
	if adminInfo.Role != model.GlobalAdmin {
		c.JSONResponse(ctx, false, "未授权访问", nil)
		return
	}

	id := common.Str2Int(ctx.PostForm("id"))
	if id == 0 {
		c.JSONResponse(ctx, false, "参数异常", nil)
		return
	}

	if err := service.NewTaskService().DeleteAdmin(id); err != nil {
		global.Logger.Warn("删除管理员失败", zap.Error(err))
		c.JSONResponse(ctx, false, "删除管理员失败", nil)
		return
	}

	service.NewLogService().CreateLog(adminInfo.ID, "删除管理员", id)
	c.JSONResponse(ctx, true, "删除管理员成功", nil)
}

func (c *ControllerApiAdmin) Registry(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		c.JSONResponse(ctx, false, "用户名或密码不能为空", nil)
		return
	}
	password = common.HashMD5(password)

	var admin model.Admin
	admin.Username = username
	admin.Password = password
	admin.Role = model.AppAdmin

	expiresAt := time.Now().AddDate(0, 3, 0)
	admin.ExpiresAt = types.Customtime(expiresAt)
	if err := global.DB.Create(&admin).Error; err != nil {
		global.Logger.Warn("注册管理员失败", zap.Error(err))
		c.JSONResponse(ctx, false, "注册管理员失败", nil)
		return
	}
	service.NewLogService().CreateLog(admin.ID, "注册", admin)
	c.JSONResponse(ctx, true, "注册管理员成功", admin)
}

func (c *ControllerApiAdmin) Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username == "" || password == "" {
		c.JSONResponse(ctx, false, "用户名或密码不能为空", nil)
		return
	}
	password = common.HashMD5(password)
	var admin model.Admin
	if err := global.DB.Where("username = ? AND password = ?", username, password).First(&admin).Error; err != nil {
		global.Logger.Warn("用户名或密码错误", zap.Error(err))
		c.JSONResponse(ctx, false, "用户名或密码错误", nil)
		return
	}

	if time.Time(admin.ExpiresAt).Before(time.Now()) {
		c.JSONResponse(ctx, false, "账号已过期", nil)
		return
	}

	token, _ := common.HashToken()
	var login model.Login
	login.AdminID = admin.ID
	login.Token = token
	expiresAt := time.Now().Add(time.Hour * 2)
	login.ExpiresAt = types.Customtime(expiresAt)
	if err := global.DB.Create(&login).Error; err != nil {
		global.Logger.Warn("登录失败", zap.Error(err))
		c.JSONResponse(ctx, false, "登录失败", nil)
		return
	}
	service.NewLogService().CreateLog(admin.ID, "登录", login)
	c.JSONResponse(ctx, true, "登录成功", gin.H{
		"username":   admin.Username,
		"rolename":   model.Rolenames[admin.Role],
		"truename":   admin.Truename,
		"phone":      admin.Phone,
		"email":      admin.Email,
		"expires_at": common.FormatTime(&admin.ExpiresAt),
		"token":      token,
		"permiss":    model.RolePermissions[admin.Role],
	})
}

func (c *ControllerApiAdmin) Loginout(ctx *gin.Context) {
	adminInfo := c.CheckAdmin(ctx)

	var login model.Login
	if err := global.DB.Where("admin_id =?", adminInfo.ID).First(&login).Error; err != nil {
		global.Logger.Warn("退出登录失败", zap.Error(err))
		c.JSONResponse(ctx, false, "退出登录失败", nil)
		return
	}
	if err := global.DB.Delete(&login).Error; err != nil {
		global.Logger.Warn("退出登录失败", zap.Error(err))
		c.JSONResponse(ctx, false, "退出登录失败", nil)
		return
	}
	service.NewLogService().CreateLog(adminInfo.ID, "退出登录", login)
	c.JSONResponse(ctx, true, "退出登录成功", nil)
}
