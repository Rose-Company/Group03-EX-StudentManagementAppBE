// internal/services/auth/interface.go
package student

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	// Define the methods that the service layer should implement
	GetByID(ctx context.Context, id uuid.UUID) (*models.StudentResponse, error)
}
