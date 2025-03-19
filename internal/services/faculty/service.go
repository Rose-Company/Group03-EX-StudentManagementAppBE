package faculty

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"context"
)

type Service interface {
	GetList(ctx context.Context) (*models.ListFacultyResponse, error)
}

type facultyService struct {
	facultyRepo faculty.Repository
}

func NewFalcutyService(facultyRepo faculty.Repository) Service {
	return &facultyService{
		facultyRepo: facultyRepo,
	}
}

func (s *facultyService) GetList(ctx context.Context) (*models.ListFacultyResponse, error) {
	faculties, err := s.facultyRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return &models.ListFacultyResponse{
		Items: faculties,
	}, nil
}
