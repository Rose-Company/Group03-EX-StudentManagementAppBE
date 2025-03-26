package program

import (
	"Group03-EX-StudentManagementAppBE/internal/models"
	models_program "Group03-EX-StudentManagementAppBE/internal/models/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"context"

	log "github.com/sirupsen/logrus"
)

type Service interface {
	ListPrograms(ctx context.Context, userID string, req *models_program.ListProgramRequest) ([]*models_program.Program, error)
	CreateProgram(ctx context.Context, userID string, program *models_program.Program) error
	UpdateProgram(ctx context.Context, userID string, id string, program *models_program.Program) error
	DeleteProgram(ctx context.Context, userID string, id string) error
}

type service struct {
	programRepo program.Repository
}

func NewProgramService(programRepo program.Repository) Service {
	return &service{
		programRepo: programRepo,
	}
}

func (s *service) ListPrograms(ctx context.Context, userID string, req *models_program.ListProgramRequest) ([]*models_program.Program, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "ListPrograms",
		"userID":   userID,
	})

	logger.Info("Fetching program list with filters")
	// Set default sort if not provided
	if req.Sort == "" {
		req.Sort = "name.asc"
	}

	params := models.QueryParams{
		QuerySort: models.QuerySort{
			Origin: req.Sort,
		},
	}

	programs, err := s.programRepo.List(ctx, params)
	if err != nil {
		logger.WithError(err).Error("Failed to fetch program list")
		return nil, err
	}
	logger.WithField("count", len(programs)).Info("Successfully fetched programs")
	return programs, nil
}

func (s *service) CreateProgram(ctx context.Context, userID string, program *models_program.Program) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "CreateProgram",
		"userID":   userID,
	})

	logger.Info("Creating new program")
	if _, err := s.programRepo.Create(ctx, program); err != nil {
		logger.WithError(err).Error("Failed to create program")
		return err
	}
	logger.WithField("programID", program.ID).Info("Program created successfully")
	return nil
}

func (s *service) UpdateProgram(ctx context.Context, userID string, id string, program *models_program.Program) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "UpdateProgram",
		"userID":   userID,
		"id":       id,
	})

	logger.Info("Updating program")
	if _, err := s.programRepo.Update(ctx, id, program); err != nil {
		logger.WithError(err).Error("Failed to update program")
		return err
	}
	logger.Info("Program updated successfully")
	return nil
}

func (s *service) DeleteProgram(ctx context.Context, userID string, id string) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "DeleteProgram",
		"userID":   userID,
		"id":       id,
	})

	logger.Info("Deleting program")
	err := s.programRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("Failed to delete program")
		return err
	}
	logger.Info("Program deleted successfully")
	return nil
}
