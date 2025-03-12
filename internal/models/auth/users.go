// internal/models/user.go
package models

import (
	"Group03-EX-StudentManagementAppBE/common"
)

type User struct {
	ID       string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"-" gorm:"password"`
	RoleID   string `json:"role_id" gorm:"column:role"`
}

func (User) TableName() string {
	return common.POSTGRES_TABLE_NAME_USERS
}
