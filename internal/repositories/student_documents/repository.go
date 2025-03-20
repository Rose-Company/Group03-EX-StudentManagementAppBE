package student_identity_documents

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.StudentDocument]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.StudentDocument](db)
}
