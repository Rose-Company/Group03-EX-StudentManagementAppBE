package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"context"

	"gorm.io/gorm"
)

type Service interface {
	// Define the methods that the service layer should implement
	GetByID(ctx context.Context, id string) (*models.StudentResponse, error)
	GetList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error)
	CreateAStudent(ctx context.Context, req *models.Student) (*models.StudentResponse, error)
	UpdateStudent(ctx context.Context, id string, req *models.Student) (*models.StudentResponse, error)
	DeleteByID(ctx context.Context, id string) error
	GetStatuses(ctx context.Context, sort string) ([]*models.StudentStatus, error)
	CreateStudentStatus(ctx context.Context, studentStatus *models.StudentStatus) (*models.StudentStatus, error)
	UpdateStudentStatus(ctx context.Context,id string, studentStatus *models.StudentStatus) (*models.StudentStatus, error)
	DeleteStudentStatus(ctx context.Context, id string) error
}

type studentService struct {
	studentRepo       student.Repository
	studentStatusRepo student_status.Repository
}

func NewStudentService(studentRepo student.Repository, studentStatusRepo student_status.Repository) Service {
	return &studentService{
		studentRepo:       studentRepo,
		studentStatusRepo: studentStatusRepo,
	}
}

func (s *studentService) GetByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return student.ToResponse(), nil
}

func (s *studentService) GetList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error) {
	if req.Sort == "" {
		req.Sort = "student_code.asc"
	}

	totalCount, err := s.studentRepo.Count(ctx, models2.QueryParams{}, func(tx *gorm.DB) {
		// Apply student_code filter if provided
		if req.StudentCode != "" {
			tx.Where("CAST(student_code AS TEXT) LIKE ?", req.StudentCode+"%")
		}

		// Apply fullname filter if provided
		if req.Fullname != "" {
			tx.Where("LOWER(fullname) LIKE LOWER(?)", "%"+req.Fullname+"%")
		}
	})
	if err != nil {
		return nil, err
	}
	// Calculate page number from offset and page size

	if req.PageSize < 0 {
		return nil, common.ErrInvalidInput
	}
	offset := (req.Page - 1) * req.PageSize

	// Get the paginated list of students
	students, err := s.studentRepo.List(ctx, models2.QueryParams{
		Offset: offset,
		Limit:  req.PageSize,
		QuerySort: models2.QuerySort{
			Origin: req.Sort,
		},
	}, func(tx *gorm.DB) {
		if req.StudentCode != "" {
			tx.Where("CAST(student_code AS TEXT) LIKE ?", req.StudentCode+"%")
		}

		if req.Fullname != "" {
			tx.Where("LOWER(fullname) LIKE LOWER(?)", "%"+req.Fullname+"%")
		}
	})

	if err != nil {
		return nil, err
	}

	// Convert students to response DTOs
	var studentResponses []*models.StudentResponse
	for _, student := range students {
		studentResponses = append(studentResponses, student.ToResponse())
	}

	// Create the paginated response
	response := &models2.BaseListResponse{
		Total:    int(totalCount),
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    studentResponses,
		Extra:    nil, // Add extra data if needed
	}

	return response, nil
}

func (s *studentService) CreateAStudent(ctx context.Context, student *models.Student) (*models.StudentResponse, error) {
	if student == nil {
		return nil, common.ErrInvalidInput
	}
	createdStudent, err := s.studentRepo.Create(ctx, student)
	if err != nil {
		return nil, err
	}
	return createdStudent.ToResponse(), nil
}

func (s *studentService) UpdateStudent(ctx context.Context, id string, student *models.Student) (*models.StudentResponse, error) {
	if student == nil {
		return nil, common.ErrInvalidInput
	}
	updatedStudent, err := s.studentRepo.Update(ctx, id, student)
	if err != nil {
		return nil, err
	}
	return updatedStudent.ToResponse(), nil
}

func (s *studentService) DeleteByID(ctx context.Context, id string) error {
	return s.studentRepo.DeleteByID(ctx, id)
}


func (s *studentService) GetStatuses(ctx context.Context, sort string) ([]*models.StudentStatus, error) {
	// Cấu hình query params để truyền sort
	queryParams := models2.QueryParams{
		QuerySort: models2.QuerySort{
			Sort: sort, // Truyền sort vào QuerySort
		},
	}

	// Gọi repository với query params đã có sort
	studentStatus, err := s.studentStatusRepo.List(ctx, queryParams)
	if err != nil {
		return nil, err
	}

	return studentStatus, nil
}


func (s *studentService) CreateStudentStatus(ctx context.Context, studentStatus *models.StudentStatus) (*models.StudentStatus, error) {
	createdStudentStatus, err := s.studentStatusRepo.Create(ctx, studentStatus)
	if err != nil {
		return nil, err
	}
	return createdStudentStatus, nil
}

func (s *studentService) UpdateStudentStatus(ctx context.Context,id string, studentStatus *models.StudentStatus) (*models.StudentStatus, error) {
	updatedStudentStatus, err := s.studentStatusRepo.Update(ctx,id, studentStatus)
	if err != nil {
		return nil, err
	}
	return updatedStudentStatus, nil
}

func (s *studentService) DeleteStudentStatus(ctx context.Context, id string) error {
	return s.studentStatusRepo.DeleteByID(ctx, id)
}
