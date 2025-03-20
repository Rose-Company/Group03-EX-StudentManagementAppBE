package admin

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/admin"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.File]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.File](db)
}
