package status_transition_rule

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/status_transition_rule"

	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[models.StatusTransitionRule]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[models.StatusTransitionRule](db)
}
