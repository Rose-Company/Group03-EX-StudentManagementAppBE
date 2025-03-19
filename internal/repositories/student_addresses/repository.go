// internal/repositories/student_address/repository.go
package student_address

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.StudentAddress]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.StudentAddress](db)
}
