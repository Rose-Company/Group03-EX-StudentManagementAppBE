package faculty

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"context"
)

type Service interface {
	GetList(ctx context.Context, sort string) (*models.ListFacultyResponse, error)
	CreateAFaculty(ctx context.Context, faculty *models.Faculty) (*models.Faculty, error)
	UpdateFaculty(ctx context.Context,id string, faculty *models.Faculty) (*models.Faculty, error)
	DeleteFaculty(ctx context.Context, id string) error
}

type facultyService struct {
	facultyRepo faculty.Repository
}

func NewFalcutyService(facultyRepo faculty.Repository) Service {
	return &facultyService{
		facultyRepo: facultyRepo,
	}
}

func (s *facultyService) GetList(ctx context.Context, sort string) (*models.ListFacultyResponse, error) {
    faculties, err := s.facultyRepo.GetListFaculty(ctx, sort)
    if err != nil {
        return nil, err
    }

    return &models.ListFacultyResponse{
        Items: faculties,
    }, nil
}

func (s *facultyService) CreateAFaculty(ctx context.Context, faculty *models.Faculty) (*models.Faculty, error) {
	createdFaculty, err := s.facultyRepo.Create(ctx, faculty)
	if err != nil {
		return nil, err
	}

	return createdFaculty, nil
}

func (s *facultyService) UpdateFaculty(ctx context.Context,id string, faculty *models.Faculty) (*models.Faculty, error) {
	updatedFaculty, err := s.facultyRepo.Update(ctx, id, faculty)
	if err != nil {
		return nil, err
	}

	return updatedFaculty, nil
}

func (s *facultyService) DeleteFaculty(ctx context.Context, id string) error {
	err := s.facultyRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

