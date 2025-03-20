package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
)

type Faculty struct {
	ID   uint   `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}

type ListFacultyRequest struct {
	models.BaseRequestParamsUri
	Name string `form:"name"`
	Sort string `form:"sort"` 
}

type CreateFacultyRequest struct {
	ID   uint   `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}

type UpdateFacultyRequest struct {
	ID   uint   `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}
// Ensure Faculty implements Model interface
var _ repositories.Model = (*Faculty)(nil)

func (Faculty) TableName() string {
	return common.POSTGRES_TABLE_NAME_FACULTY
}

type ListFacultyResponse struct {
	Items []Faculty `json:"items"`
}
