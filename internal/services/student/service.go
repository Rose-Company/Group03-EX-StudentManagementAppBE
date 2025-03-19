package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
	GetStudentByID(ctx context.Context, id string) (*models.StudentResponse, error)
	GetStudentList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error)
	CreateAStudent(ctx context.Context, req *models.Student) (*models.StudentResponse, error)
	UpdateStudent(ctx context.Context, id string, req *models.Student) (*models.StudentResponse, error)
	DeleteStudentByID(ctx context.Context, id string) error
	GetStudentStatuses(ctx context.Context) ([]*models.StudentStatus, error)
}

type studentService struct {
	studentRepo         student.Repository
	studentStatusRepo   student_status.Repository
	studentAddressRepo  student_addresses.Repository
	studentDocumentRepo student_identity_documents.Repository
}

func NewStudentService(
	studentRepo student.Repository,
	studentStatusRepo student_status.Repository,
	studentAddressRepo student_addresses.Repository,
	studentDocumentRepo student_identity_documents.Repository) Service {
	return &studentService{
		studentRepo:         studentRepo,
		studentStatusRepo:   studentStatusRepo,
		studentAddressRepo:  studentAddressRepo,
		studentDocumentRepo: studentDocumentRepo,
	}
}

func (s *studentService) GetStudentByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	var clauses []repositories.Clause
	clauses = append(clauses, func(tx *gorm.DB) {
		tx.Preload("Addresses", func(db *gorm.DB) *gorm.DB {
			return db.Order("address_type")
		})
	})

	// Preload documents with filtering
	clauses = append(clauses, func(tx *gorm.DB) {
		tx.Preload("Documents", func(db *gorm.DB) *gorm.DB {
			return db.Where("(document_type LIKE ? OR document_type LIKE ? OR document_type LIKE ?) AND student_id = ?",
				"CCCD%", "CMND%", "Passport%", id)
		})
	})

	clauses = append(clauses, func(tx *gorm.DB) {
		tx.Joins(`LEFT JOIN "PUBLIC"."faculties" ON students.faculty_id = "PUBLIC"."faculties".id`)
		tx.Select(`students.*, "PUBLIC"."faculties".name as faculty_name`).Where("students.id = ?", id)
	})

	combinedClause := func(tx *gorm.DB) {
		for _, clause := range clauses {
			clause(tx)
		}
	}

	// Get student with all related data
	student, err := s.studentRepo.GetDetailByConditions(ctx, combinedClause)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}

	return student.ToResponse(), nil
}

func (s *studentService) GetStudentList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error) {
	if req.Sort == "" {
		req.Sort = "student_code.asc"
	}

	var clauses []repositories.Clause

	if req.StudentCode != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Where("CAST(student_code AS TEXT) LIKE ?", req.StudentCode+"%")
		})
	}

	if req.Fullname != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Where("LOWER(fullname) LIKE LOWER(?)", "%"+req.Fullname+"%")
		})
	}

	if req.FacultyName != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Joins(`JOIN "PUBLIC".faculties ON students.faculty_id = faculties.id`).Where("LOWER(faculties.name) LIKE LOWER(?)", "%"+req.FacultyName+"%")
		})
	}

	combinedClause := func(tx *gorm.DB) {
		for _, clause := range clauses {
			clause(tx)
		}
	}

	totalCount, err := s.studentRepo.Count(ctx, models2.QueryParams{}, combinedClause)
	if err != nil {
		return nil, err
	}

	if req.PageSize < 0 {
		return nil, common.ErrInvalidInput
	}
	offset := (req.Page - 1) * req.PageSize

	students, err := s.studentRepo.List(ctx, models2.QueryParams{
		Offset: offset,
		Limit:  req.PageSize,
		QuerySort: models2.QuerySort{
			Origin: req.Sort,
		},
	}, combinedClause)

	if err != nil {
		return nil, err
	}

	var studentResponses []*models.StudentListResponse
	for _, student := range students {
		studentResponses = append(studentResponses, student.ToListResponse())
	}

	response := &models2.BaseListResponse{
		Total:    int(totalCount),
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    studentResponses,
		Extra:    nil,
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

func (s *studentService) DeleteStudentByID(ctx context.Context, id string) error {
	return s.studentRepo.DeleteByID(ctx, id)
}

func (s *studentService) GetStudentStatuses(ctx context.Context) ([]*models.StudentStatus, error) {
	studentStatus, err := s.studentStatusRepo.List(ctx, models2.QueryParams{}, func(tx *gorm.DB) {

	})
	if err != nil {
		return nil, err
	}
	return studentStatus, nil
}
