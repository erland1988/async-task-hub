package model

import (
	"asynctaskhub/src/types"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type AdminRole string

const (
	GlobalAdmin AdminRole = "global_admin" // 全局管理员
	AppAdmin    AdminRole = "app_admin"    // 应用管理员
)

var Rolenames = map[AdminRole]string{
	GlobalAdmin: "全局管理员",
	AppAdmin:    "应用管理员",
}

var RolePermissions = map[AdminRole][]string{
	GlobalAdmin: {"0", "1", "11", "2", "21", "3", "31", "4", "41", "5", "51", "6", "61"},
	AppAdmin:    {"0", "2", "21", "3", "31", "4", "41", "5", "51"},
}

type Admin struct {
	ID        int              `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`    // 主键
	Username  string           `gorm:"type:varchar(32);not null;uniqueIndex:username" json:"username"` // 用户名
	Password  string           `gorm:"type:varchar(64);not null" json:"-"`                             // 密码
	Truename  string           `gorm:"type:varchar(32);default:null" json:"truename"`                  // 真实姓名
	Phone     string           `gorm:"type:varchar(32);default:null" json:"phone"`                     // 手机号
	Email     string           `gorm:"type:varchar(64);default:null" json:"email"`                     // 邮箱
	Role      AdminRole        `gorm:"type:varchar(32);not null;index:role" json:"role"`               // 角色
	ExpiresAt types.Customtime `gorm:"type:datetime;not null;index:expires_at" json:"expires_at"`      // 过期时间
	CreatedAt types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index:created_at" json:"created_at"`
	UpdatedAt types.Customtime `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:updated_at" json:"updated_at"`
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

func (m *Admin) MarshalJSON() ([]byte, error) {
	type Alias Admin
	return json.Marshal(&struct {
		*Alias
		RoleName string `json:"rolename"`
	}{
		Alias:    (*Alias)(m),
		RoleName: Rolenames[m.Role],
	})
}

func (m *Admin) GetAdminRoles() []string {
	return []string{string(GlobalAdmin), string(AppAdmin)}
}

func (m *Admin) TableName() string { return "admins" }
