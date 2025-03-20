package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models"
)


type StudentStatus struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:text;not null"`
}

type ListStudentStatusRequest struct {
	models.BaseRequestParamsUri
	Name string `form:"name"`
	Sort string `form:"sort"`
}

type CreateStudentStatusRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateStudentStatusRequest struct {
	Name string `json:"name"`
}

func (s *StudentStatus) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENTS_STATUSES
}

type ListStudentStatusResponse struct {
	Items []*StudentStatus `json:"items"`
}