package repositories

import (
	"Group03-EX-StudentManagementAppBE/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error)
}

type studentRepository struct {
	BaseRepository[models.Student]
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{
		BaseRepository: NewBaseRepository[models.Student](db),
	}
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	student, err := r.BaseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return student, nil
}
