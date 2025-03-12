package student

import (
	"Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"context"

	"github.com/google/uuid"
)

type StudentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.StudentResponse, error)
}

type studentService struct {
	studentRepo student.Repository
}

func NewStudentService(studentRepo student.Repository) Service {
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
