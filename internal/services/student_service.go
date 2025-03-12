package services

import (
	"Group03-EX-StudentManagementAppBE/internal/models"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"context"

	"github.com/google/uuid"
)

type StudentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.StudentResponse, error)
}

type studentService struct {
	studentRepo repositories.StudentRepository
}

func NewStudentService(studentRepo repositories.StudentRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
	}
}

func (s *studentService) GetByID(ctx context.Context, id uuid.UUID) (*models.StudentResponse, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return student.ToResponse(), nil
}
