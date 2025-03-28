package validation

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/validation"
	"Group03-EX-StudentManagementAppBE/internal/repositories"

	"gorm.io/gorm"
)

// ValidationSettingRepository defines the repository interface for ValidationSetting
type ValidationSettingRepository interface {
	repositories.BaseRepository[models.ValidationSetting]
}

// NewValidationSettingRepository creates a new repository for ValidationSetting
func NewValidationSettingRepository(db *gorm.DB) ValidationSettingRepository {
	return repositories.NewBaseRepository[models.ValidationSetting](db)
}
