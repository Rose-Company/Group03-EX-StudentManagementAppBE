package faculty

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.Faculty]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.Faculty](db)
}
