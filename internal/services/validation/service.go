package validation

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/validation"
	"Group03-EX-StudentManagementAppBE/internal/repositories/validation"
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Service interface {
	UpdateValidationSetting(ctx context.Context,id string, req *models.ValidationSettingUpdateRequest) (*models.ValidationSetting, error)
}

type validationSettingService struct {
	validationSettingRepo validation.ValidationSettingRepository
}

func NewValidationService(validationSettingRepo validation.ValidationSettingRepository) Service {
	return &validationSettingService{
		validationSettingRepo: validationSettingRepo,
	}
}

func (s *validationSettingService) UpdateValidationSetting(ctx context.Context, idStr string, req *models.ValidationSettingUpdateRequest) (*models.ValidationSetting, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "UpdateValidationSetting",
		"id":       idStr,
	})

	logger.Info("Updating validation setting")

	// Chuyển đổi id từ string sang int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WithError(err).Error("Invalid ID format")
		return nil, common.ErrInvalidFormat
	}

	// Prepare update data
	updateData := &models.ValidationSetting{
		ValidationKey: req.ValidationKey,
		IsEnabled:     req.IsEnabled,
	}

	// Update validation setting
	updatedSetting, err := s.validationSettingRepo.Update(
		ctx, 
		id, // Truyền vào int thay vì string
		updateData,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to update validation setting")
		return nil, err
	}

	logger.WithField("id", id).Info("Successfully updated validation setting")
	return updatedSetting, nil
}
