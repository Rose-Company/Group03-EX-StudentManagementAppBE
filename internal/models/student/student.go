package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StudentCode int       `json:"student_code" gorm:"unique;not null"`
	Fullname    string    `json:"fullname" gorm:"type:text;not null"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date;not null"`
	Gender      string    `json:"gender" gorm:"type:text;check:gender IN ('Male', 'Female', 'Other')"`
	FacultyID   int       `json:"faculty_id" gorm:"type:integer;references:faculties(id)"`
	Batch       string    `json:"batch" gorm:"type:text;not null"`
	Program     string    `json:"program" gorm:"type:text;not null"`
	Address     string    `json:"address" gorm:"type:text"`
	Email       string    `json:"email" gorm:"type:text;unique"`
	Phone       string    `json:"phone" gorm:"type:text"`
	StatusID    int       `json:"status_id" gorm:"type:integer;references:student_statuses(id)"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;references:users(id)"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

func (s *Student) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENTS
}

type StudentResponse struct {
	ID          uuid.UUID `json:"id"`
	StudentCode int       `json:"student_code"`
	Fullname    string    `json:"fullname"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	FacultyID   int       `json:"faculty_id"`
	Batch       string    `json:"batch"`
	Program     string    `json:"program"`
	Address     string    `json:"address"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	StatusID    int       `json:"status_id"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Student) ToResponse() *StudentResponse {
	return &StudentResponse{
		ID:          s.ID,
		StudentCode: s.StudentCode,
		Fullname:    s.Fullname,
		DateOfBirth: s.DateOfBirth,
		Gender:      s.Gender,
		FacultyID:   s.FacultyID,
		Batch:       s.Batch,
		Program:     s.Program,
		Address:     s.Address,
		Email:       s.Email,
		Phone:       s.Phone,
		StatusID:    s.StatusID,
		UserID:      s.UserID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

type ListStudentRequest struct {
	models.BaseRequestParamsUri
}
