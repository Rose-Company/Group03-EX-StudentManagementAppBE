package student_status

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.StudentStatus]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.StudentStatus](db)
}
