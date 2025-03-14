package services

import (
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	facultyService "Group03-EX-StudentManagementAppBE/internal/services/faculty"
	studentService "Group03-EX-StudentManagementAppBE/internal/services/student"
)

type Service struct {
	Auth    auth.Service
	Student studentService.Service
	Faculty facultyService.Service
}

func NewService(userRepo user.Repository, studentRepo student.Repository, facultyRepo faculty.Repository, studentStatusRepo student_status.Repository) *Service {
	return &Service{
		Auth:    auth.NewService(userRepo),
		Student: studentService.NewService(studentRepo, studentStatusRepo),
		Faculty: facultyService.NewService(facultyRepo),
	}
}
