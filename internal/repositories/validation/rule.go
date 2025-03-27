package validation

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/validation"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

// ValidationRuleRepository defines the repository interface for ValidationRule
type ValidationRuleRepository interface {
	repositories.BaseRepository[models.ValidationRule]
}

// NewValidationRuleRepository creates a new repository for ValidationRule
func NewValidationRuleRepository(db *gorm.DB) ValidationRuleRepository {
	return repositories.NewBaseRepository[models.ValidationRule](db)
}
