package student_status

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"context"

	log "github.com/sirupsen/logrus"
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
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "GetStatuses",
	})

	logger.Info("Fetching student statuses")
	sort := common.ParseSortString(req.Sort)

	studentStatus, err := s.studentStatusRepo.List(ctx, models2.QueryParams{
		QuerySort: models2.QuerySort{
			Sort: sort,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to fetch student statuses")
		return nil, err
	}

	logger.WithField("count", len(studentStatus)).Info("Successfully fetched student statuses")
	return studentStatus, nil
}

func (s *studentStatusService) CreateStudentStatus(ctx context.Context, req *models.CreateStudentStatusRequest) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "CreateStudentStatus",
	})

	logger.Info("Creating new student status")
	studentStatus := &models.StudentStatus{
		Name: req.Name,
	}

	_, err := s.studentStatusRepo.Create(ctx, studentStatus)
	if err != nil {
		logger.WithError(err).Error("Failed to create student status")
		return err
	}

	logger.Info("Student status created successfully")
	return nil
}

func (s *studentStatusService) UpdateStudentStatus(ctx context.Context, id string, req *models.UpdateStudentStatusRequest) (*models.StudentStatus, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "UpdateStudentStatus",
		"id":       id,
	})

	logger.Info("Updating student status")
	studentStatus := &models.StudentStatus{
		Name: req.Name,
	}

	updatedStudentStatus, err := s.studentStatusRepo.Update(ctx, id, studentStatus)
	if err != nil {
		logger.WithError(err).Error("Failed to update student status")
		return nil, err
	}

	logger.Info("Student status updated successfully")
	return updatedStudentStatus, nil
}

func (s *studentStatusService) DeleteStudentStatus(ctx context.Context, id string) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "DeleteStudentStatus",
		"id":       id,
	})

	logger.Info("Deleting student status")
	err := s.studentStatusRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("Failed to delete student status")
		return err
	}

	logger.Info("Student status deleted successfully")
	return nil
}
