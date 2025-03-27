package models

import "time"

// ValidationRule represents the validation_rules table
type ValidationRule struct {
	ID              int       `gorm:"primaryKey;column:id"`
	SettingID       int       `gorm:"column:setting_id"`
	RuleValue       string    `gorm:"column:rule_value"`
	RuleDescription string    `gorm:"column:rule_description"`
	IsEnabled       bool      `gorm:"column:is_enabled"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}
