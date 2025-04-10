package program

import (
	"Group03-EX-StudentManagementAppBE/common"
)

type Program struct {
	ID   int    `json:"id" gorm:"type:serial;primary_key"`
	Name string `json:"name" gorm:"type:text;not null"`
}

func (p *Program) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENT_PROGRAMS
}

type ListProgramRequest struct {
	Sort string `form:"sort"`
}
