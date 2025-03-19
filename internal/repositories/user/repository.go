package user

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"
	repositories "Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

// Repository defines the user repository interface
type Repository interface {
	repositories.BaseRepository[models.User]
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.User](db)
}
