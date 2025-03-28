package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"time"
)

// StatusTransitionRule represents a rule for transitioning between statuses
type StatusTransitionRule struct {
    ID          int       ` gorm:"type:serial;primary_key" json:"id"`
    FromStatusID int       `gorm:"column:from_status_id;not null" json:"from_status_id"`
    ToStatusID   int       `gorm:"column:to_status_id;not null" json:"to_status_id"`
    IsEnabled    bool      `gorm:"column:is_enabled;not null;default:true" json:"is_enabled"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}


// ListStatusTransitionRuleRequest represents a request to list status transition rules
type ListStatusTransitionRuleRequest struct {
    Page     int `form:"page"`
    PageSize int `form:"page_size"`
    Sort     string `form:"sort"`
    FromStatusID int `form:"from_status_id"`
    ToStatusID   int `form:"to_status_id"`
    IsEnabled    bool `form:"is_enabled"`
    CreatedAt    time.Time `form:"created_at"`
    UpdatedAt    time.Time `form:"updated_at"`
}


// CreateStatusTransitionRuleRequest represents a request to create a status transition rule
type CreateStatusTransitionRuleRequest struct {
    FromStatusID int `json:"from_status_id"`
    ToStatusID   int `json:"to_status_id"`
    IsEnabled    bool `json:"is_enabled"`
}

// UpdateStatusTransitionRuleRequest represents a request to update a status transition rule
type UpdateStatusTransitionRuleRequest struct {
    FromStatusID int `json:"from_status_id"`
    ToStatusID   int `json:"to_status_id"`
    IsEnabled    bool `json:"is_enabled"`
}

func (s *StatusTransitionRule) TableName() string {
	return common.POSTGRES_TABLE_NAME_STATUS_TRANSITION_RULES
}
