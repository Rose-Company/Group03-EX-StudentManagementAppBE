package services

import (
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	studentService "Group03-EX-StudentManagementAppBE/internal/services/student"
)

type Service struct {
	Auth    auth.Service
	Student studentService.Service // Changed from studentService to Student for consistency
}

func NewService(userRepo user.Repository, studentRepo student.Repository) *Service {
	return &Service{
		Auth:    auth.NewService(userRepo),
		Student: studentService.NewService(studentRepo), // Changed to match the field name
	}
}
