package student

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"context"
)

type studentService struct {
	studentRepo student.Repository
}

func NewService(studentRepo student.Repository) Service {
	return &studentService{
		studentRepo: studentRepo,
	}
}

func (s *studentService) GetByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return student.ToResponse(), nil
}
