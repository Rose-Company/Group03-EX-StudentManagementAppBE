// internal/models/user.go
package models

import (
	"Group03-EX-StudentManagementAppBE/common"
)

type User struct {
	ID                 string  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email              string  `json:"email" gorm:"uniqueIndex"`
	FirstName          *string `json:"first_name" gorm:"column:first_name"`
	LastName           *string `json:"last_name" gorm:"column:last_name"`
	Password           string  `json:"-" gorm:"password"`
	RoleID             string  `json:"role_id" gorm:"column:role"`
	IsBanned           bool    `json:"is_banned" gorm:"is_banned"`
	EmailNotifications bool    `json:"email_notifications" gorm:"email_notifications"`
	Provider           string  `gorm:"nullable"`
	Avatar             string  `json:"avatar" gorm:"avatar"`
}

func (User) TableName() string {
	return common.POSTGRES_TABLE_NAME_USERS
}
