package program

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.Program]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.Program](db)
}
