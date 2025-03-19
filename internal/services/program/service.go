package program

import (
	"Group03-EX-StudentManagementAppBE/internal/models"
	models_program "Group03-EX-StudentManagementAppBE/internal/models/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"context"
)

type Service interface {
	ListPrograms(ctx context.Context, req *models_program.ListProgramRequest) ([]*models_program.Program, error)
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
