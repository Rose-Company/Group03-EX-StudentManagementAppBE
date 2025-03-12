package services

import (
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
)

type Service struct {
	Auth auth.Service
	studentService student.Service
}

func NewService(userRepo user.Repository, studentRepo student.Repository) *Service {
	return &Service{
		Auth: auth.NewService(userRepo),
		studentService: student.NewStudentService(studentRepo),
	}
}
