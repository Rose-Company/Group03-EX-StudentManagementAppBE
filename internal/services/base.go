package services

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/repositories/admin"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	adminService "Group03-EX-StudentManagementAppBE/internal/services/admin"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
	facultyService "Group03-EX-StudentManagementAppBE/internal/services/faculty"
	gdriveService "Group03-EX-StudentManagementAppBE/internal/services/gdrive"
	programService "Group03-EX-StudentManagementAppBE/internal/services/program"
	studentService "Group03-EX-StudentManagementAppBE/internal/services/student"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service is a container for all services
type Service struct {
	Auth    auth.Service
	Student studentService.Service
	Faculty facultyService.Service
	Program programService.Service
	GDrive  gdriveService.Service
	Admin   adminService.Service
}

// NewService creates a new service container with all dependencies
func NewService(userRepo user.Repository,
	studentRepo student.Repository,
	facultyRepo faculty.Repository,
	studentStatusRepo student_status.Repository,
	studentAddressRepo student_addresses.Repository,
	studentDocumentRepo student_identity_documents.Repository,
	adminRepo admin.Repository,
	db *gorm.DB,
	program program.Repository) *Service {
	// Load config for JWT secret
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize service container with placeholder for Student service
	service := &Service{
		Auth:    auth.NewAuthService(userRepo, cfg.JWTSecret),
		Faculty: facultyService.NewFalcutyService(facultyRepo),
	}

	// Initialize Google Drive service if credentials are configured
	var driveSvc gdriveService.Service
	if cfg.GoogleDriveCredentialsFile != "" {
		var err error
		driveSvc, err = gdriveService.NewHTTPDriveService(cfg.GoogleDriveCredentialsFile)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to initialize Google Drive service")
		} else {
			log.Info("Google Drive service initialized successfully")
			service.GDrive = driveSvc

			// Initialize admin service with the drive service
			service.Admin = adminService.NewAdminService(adminRepo, driveSvc)
		}
	} else {
		log.Warn("Google Drive credentials not configured. Drive integration will be disabled.")
	}

	// Now initialize Student service with the drive service
	service.Student = studentService.NewStudentService(
		studentRepo,
		studentStatusRepo,
		studentAddressRepo,
		studentDocumentRepo,
		driveSvc, // Pass the drive service here
	)

	return service
}
