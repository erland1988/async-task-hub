package model

import (
	"asynctaskhub/common"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AdminRole string

const (
	GlobalAdmin AdminRole = "global_admin" // 全局管理员
	AppAdmin    AdminRole = "app_admin"    // 应用管理员
)

var RoleNames = map[AdminRole]string{
	GlobalAdmin: "全局管理员",
	AppAdmin:    "应用管理员",
}

var RolePermissions = map[AdminRole][]string{
	GlobalAdmin: {"0", "1", "11", "2", "21", "3", "31", "4", "41", "5", "51", "6", "61", "62"},
	AppAdmin:    {"0", "2", "21", "3", "31", "4", "41", "5", "51"},
}

type Admin struct {
	ID        int       `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`    // 主键
	Username  string    `gorm:"type:varchar(32);not null;uniqueIndex:username" json:"username"` // 用户名
	Password  string    `gorm:"type:varchar(64);not null" json:"password"`                      // 密码
	Truename  string    `gorm:"type:varchar(32);default:null" json:"truename"`                  // 真实姓名
	Phone     string    `gorm:"type:varchar(32);default:null" json:"phone"`                     // 手机号
	Email     string    `gorm:"type:varchar(64);default:null" json:"email"`                     // 邮箱
	Role      AdminRole `gorm:"type:varchar(32);not null;index:role" json:"role"`               // 角色
	ExpiresAt time.Time `gorm:"type:datetime;not null;index:expires_at" json:"expires_at"`      // 过期时间
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`
}

var validAdminRoles = map[AdminRole]struct{}{
	GlobalAdmin: {},
	AppAdmin:    {},
}

func (m *Admin) BeforeSave(tx *gorm.DB) error {
	if _, exists := validAdminRoles[m.Role]; !exists {
		return errors.New("invalid role: " + string(m.Role))
	}
	return nil
}

func (m *Admin) UnmarshalJSON(data []byte) error {
	type Alias Admin
	aux := &struct {
		ExpiresAt string `json:"expires_at"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.ExpiresAt != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", aux.ExpiresAt, time.Local)
		if err != nil {
			return errors.New("expires_at 时间格式错误，请使用 'YYYY-MM-DD HH:mm:ss'")
		}
		m.ExpiresAt = t
	}
	return nil
}

func (m *Admin) MarshalJSON() ([]byte, error) {
	type Alias Admin
	return json.Marshal(&struct {
		*Alias
		RoleName  string `json:"role_name"`
		ExpiresAt string `json:"expires_at"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(m),
		RoleName:  RoleNames[m.Role],
		ExpiresAt: common.FormatTime(&m.ExpiresAt),
		CreatedAt: common.FormatTime(&m.CreatedAt),
		UpdatedAt: common.FormatTime(&m.UpdatedAt),
	})
}
func (m *Admin) GetAdminRoles() []string {
	return []string{string(GlobalAdmin), string(AppAdmin)}
}

func (m *Admin) TableName() string { return "admins" }
