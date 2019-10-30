package sys

import (
	"github.com/jinzhu/gorm"
	"review-order/app/model/db"
	"time"
)

// 用户-角色
type AdminRole struct {
	db.Model
	AdminsID uint64 `gorm:"column:admins_id;unique_index:uk_admins_role_admins_id;not null;"` // 管理员ID
	RoleID   uint64 `gorm:"column:role_id;unique_index:uk_admins_role_admins_id;not null;"`   // 角色ID
}

// 表名
func (AdminRole) TableName() string {
	return TableName("admins_role")
}

// 添加前
func (m *AdminRole) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *AdminRole) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
