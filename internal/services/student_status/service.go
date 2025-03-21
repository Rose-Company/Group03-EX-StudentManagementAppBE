package student_status

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"context"
)

type Service interface {
	GetStatuses(ctx context.Context, req *models.ListStudentStatusRequest) ([]*models.StudentStatus, error)
	CreateStudentStatus(ctx context.Context, studentStatus *models.CreateStudentStatusRequest) error
	UpdateStudentStatus(ctx context.Context, id string, req *models.UpdateStudentStatusRequest) (*models.StudentStatus, error)
	DeleteStudentStatus(ctx context.Context, id string) error
}

type studentStatusService struct {
	studentStatusRepo student_status.Repository
}

func NewStudentService(studentStatusRepo student_status.Repository) Service {
	return &studentStatusService{
		studentStatusRepo: studentStatusRepo,
	}
}

func (s *studentStatusService) GetStatuses(ctx context.Context, req *models.ListStudentStatusRequest) ([]*models.StudentStatus, error) {
	// Cấu hình query params để truyền sort

	sort := common.ParseSortString(req.Sort)

	// Gọi repository với query params đã có sort
	studentStatus, err := s.studentStatusRepo.List(ctx, models2.QueryParams{
		QuerySort: models2.QuerySort{
			Sort: sort,
		},
	},
	)
	if err != nil {
		return nil, err
	}

	return studentStatus, nil
}

func (s *studentStatusService) CreateStudentStatus(ctx context.Context, req *models.CreateStudentStatusRequest) error {

	studentStatus := &models.StudentStatus{
		Name: req.Name,
	}

	_, err := s.studentStatusRepo.Create(ctx, studentStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *studentStatusService) UpdateStudentStatus(ctx context.Context, id string, req *models.UpdateStudentStatusRequest) (*models.StudentStatus, error) {
	studentStatus := &models.StudentStatus{
		Name: req.Name,
	}
	updatedStudentStatus, err := s.studentStatusRepo.Update(ctx, id, studentStatus)
	if err != nil {
		return nil, err
	}
	return updatedStudentStatus, nil
}

func (s *studentStatusService) DeleteStudentStatus(ctx context.Context, id string) error {
	return s.studentStatusRepo.DeleteByID(ctx, id)
}
