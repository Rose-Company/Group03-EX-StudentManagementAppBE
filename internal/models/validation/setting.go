package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"time"
)

// ValidationSetting represents the validation_settings table
type ValidationSetting struct {
	ID            int       `gorm:"primaryKey;column:id"`
	ValidationKey string    `gorm:"column:validation_key"`
	IsEnabled     bool      `gorm:"column:is_enabled"`
	Description   string    `gorm:"column:description"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (*ValidationSetting) TableName() string {
	return common.POSTGRES_TABLE_NAME_VALIDTION_SETTINGS
}
