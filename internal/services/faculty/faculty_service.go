package faculty

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"context"
)

type facultyService struct {
	facultyRepo faculty.Repository
}

func NewService(facultyRepo faculty.Repository) Service {
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
