package program

import (
	"Group03-EX-StudentManagementAppBE/internal/models"
	models_program "Group03-EX-StudentManagementAppBE/internal/models/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"context"
	"log"
)

type Service interface {
	ListPrograms(ctx context.Context, req *models_program.ListProgramRequest) ([]*models_program.Program, error)
	CreateProgram(ctx context.Context, userID string, program *models_program.Program) (*models_program.Program, error)
	UpdateProgram(ctx context.Context, userID string, id string, program *models_program.Program) (*models_program.Program, error)
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

func (s *service) ListPrograms(ctx context.Context, req *models_program.ListProgramRequest) ([]*models_program.Program, error) {
	// Set default sort if not provided
	if req.Sort == "" {
		req.Sort = "name.asc"
	}

	params := models.QueryParams{
		QuerySort: models.QuerySort{
			Origin: req.Sort,
		},
	}
	return s.programRepo.List(ctx, params)
}

func (s *service) CreateProgram(ctx context.Context, userID string, program *models_program.Program) (*models_program.Program, error) {
	createdProgram, err := s.programRepo.Create(ctx, program)
	if err != nil {
		return nil, err
	}
	log.Printf("User ID: %s created program with ID: %d", userID, createdProgram.ID)
	return createdProgram, nil
}

func (s *service) UpdateProgram(ctx context.Context, userID string, id string, program *models_program.Program) (*models_program.Program, error) {
	updatedProgram, err := s.programRepo.Update(ctx, id, program)
	if err != nil {
		return nil, err
	}
	log.Printf("User ID: %s updated program with ID: %s", userID, id)
	return updatedProgram, nil
}

func (s *service) DeleteProgram(ctx context.Context, userID string, id string) error {
	err := s.programRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	log.Printf("User ID: %s deleted program with ID: %s", userID, id)
	return nil
}
