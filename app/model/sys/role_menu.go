package sys

import (
	"github.com/jinzhu/gorm"
	"review-order/app/model/db"
	"time"
)

// 角色-菜单
type RoleMenu struct {
	db.Model
	RoleID uint64 `gorm:"column:role_id;unique_index:uk_role_menu_role_id;not null;"` // 角色ID
	MenuID uint64 `gorm:"column:menu_id;unique_index:uk_role_menu_role_id;not null;"` // 菜单ID
}

// 表名
func (RoleMenu) TableName() string {
	return TableName("role_menu")
}

// 添加前
func (m *RoleMenu) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *RoleMenu) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}
