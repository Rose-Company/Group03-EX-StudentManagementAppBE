package status_transition_rule

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/status_transition_rule"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/repositories/status_transition_rule"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
    OnOffTransition(ctx context.Context,  req *models.UpdateStatusTransitionRuleRequest) (*models.StatusTransitionRule, error)
}

type statusTransitionRuleService struct {
    statusTransitionRule status_transition_rule.Repository
}


func NewService(ruleRepo status_transition_rule.Repository) Service {
    return &statusTransitionRuleService{
        statusTransitionRule: ruleRepo,
    }
}

func (s *statusTransitionRuleService) OnOffTransition(ctx context.Context, req *models.UpdateStatusTransitionRuleRequest) (*models.StatusTransitionRule, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid input: request cannot be nil")
	}

	findClauses := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("from_status_id = ? AND to_status_id = ?", req.FromStatusID, req.ToStatusID)
		},
	}

	existingRule, err := s.statusTransitionRule.GetDetailByConditions(ctx, findClauses...)
	if err != nil {
		return nil, fmt.Errorf("status transition rule not found: %w", err)
	}

	updateColumns := map[string]interface{}{
		"is_enabled": req.IsEnabled,
	}

	updatedRule, err := s.statusTransitionRule.UpdateColumns(
		ctx, 
		existingRule.ID, 
		updateColumns,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update status transition rule: %w", err)
	}

	return updatedRule, nil
}