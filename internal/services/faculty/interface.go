package faculty

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"context"
)

type Service interface {
	GetList(ctx context.Context) (*models.ListFacultyResponse, error)
}
