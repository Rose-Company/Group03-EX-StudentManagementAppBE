package services

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	facultyService "Group03-EX-StudentManagementAppBE/internal/services/faculty"
	studentService "Group03-EX-StudentManagementAppBE/internal/services/student"
	studentStatusService "Group03-EX-StudentManagementAppBE/internal/services/student_status"
)

// Service is a container for all services
type Service struct {
	Auth    auth.Service
	Student studentService.Service
	Faculty facultyService.Service
	StudentStatus studentStatusService.Service
}

// NewService creates a new service container with all dependencies
func NewService(userRepo user.Repository, studentRepo student.Repository, facultyRepo faculty.Repository, studentStatusRepo student_status.Repository) *Service {
	// Load config for JWT secret
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Create individual service implementations
	authSvc := auth.NewAuthService(userRepo, cfg.JWTSecret)
	studentSvc := studentService.NewStudentService(studentRepo)
	facultySvc := facultyService.NewFalcutyService(facultyRepo)
	studentStatusSvc := studentStatusService.NewStudentService(studentStatusRepo)

	return &Service{
		Auth:    authSvc,
		Student: studentSvc,
		Faculty: facultySvc,
		StudentStatus: studentStatusSvc,
	}
}
