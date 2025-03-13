// internal/services/auth/interface.go
package student

import (
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"context"
)

type Service interface {
	// Define the methods that the service layer should implement
	GetByID(ctx context.Context, id string) (*models.StudentResponse, error)
	GetList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error)
	CreateAStudent(ctx context.Context, req *models.Student) (*models.StudentResponse, error)
	UpdateStudent(ctx context.Context, id string, req *models.Student) (*models.StudentResponse, error)
	DeleteByID(ctx context.Context, id string) error
	GetStatuses(ctx context.Context) ([]*models.StudentStatus, error)
}
