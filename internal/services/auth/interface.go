// internal/services/auth/interface.go
package auth

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"
	"context"
)

type Service interface {
	// Define the methods that the service layer should implement
	LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}
