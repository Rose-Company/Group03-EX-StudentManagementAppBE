package services

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"

	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	facultyService "Group03-EX-StudentManagementAppBE/internal/services/faculty"
	programService "Group03-EX-StudentManagementAppBE/internal/services/program"
	studentService "Group03-EX-StudentManagementAppBE/internal/services/student"
)

// Service is a container for all services
type Service struct {
	Auth    auth.Service
	Student studentService.Service
	Faculty facultyService.Service
	Program programService.Service
}

// NewService creates a new service container with all dependencies
func NewService(userRepo user.Repository, studentRepo student.Repository, facultyRepo faculty.Repository, studentStatusRepo student_status.Repository, studentAddressRepo student_addresses.Repository, studentDocumentRepo student_identity_documents.Repository, programRepo program.Repository) *Service {
	// Load config for JWT secret
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Create individual service implementations
	authSvc := auth.NewAuthService(userRepo, cfg.JWTSecret)
	studentSvc := studentService.NewStudentService(studentRepo, studentStatusRepo, studentAddressRepo, studentDocumentRepo)
	facultySvc := facultyService.NewFalcutyService(facultyRepo)
	programSvc := programService.NewProgramService(programRepo)

	return &Service{
		Auth:    authSvc,
		Student: studentSvc,
		Faculty: facultySvc,
		Program: programSvc,
	}
}
