package student

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"

	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.Student]
}

func NewStudentRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.Student](db)
}
