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

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Service interface {
	GetStudentByID(ctx context.Context, id string) (*models.StudentResponse, error)
	GetStudentList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error)
	CreateAStudent(ctx context.Context, userID string, req *models.CreateStudentRequest) error
	UpdateStudent(ctx context.Context, userID string, studentId string, req *models.UpdateStudentRequest) error
	DeleteStudentByID(ctx context.Context, userID string, studentID string) error
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

func (s *studentService) CreateAStudent(ctx context.Context, userID string, request *models.CreateStudentRequest) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function":   "Create A Student",
		"email":      request.Email,
		"created_by": userID,
	})

	logger.Info("Creating new student")

	if request == nil {
		return common.ErrInvalidInput
	}
	studentModel := &models.Student{
		StudentCode: *request.StudentCode,
		Fullname:    *request.Fullname,
		DateOfBirth: *request.DateOfBirth,
		Gender:      *request.Gender,
		FacultyID:   *request.FacultyID,
		Batch:       *request.Batch,
		Program:     *request.Program,
		Address:     *request.Address,
		Email:       *request.Email,
		Phone:       *request.Phone,
		StatusID:    *request.StatusID,
		ProgramID:   *request.ProgramID,
		Nationality: *request.Nationality,
	}
	createdStudent, err := s.studentRepo.Create(ctx, studentModel)
	if err != nil {
		return err
	}

	if len(request.Addresses) > 0 {
		for _, addr := range request.Addresses {
			studentAddr := &models.StudentAddress{
				StudentID:   createdStudent.ID,
				AddressType: addr.AddressType,
				Street:      addr.Street,
				Ward:        addr.Ward,
				District:    addr.District,
				City:        addr.City,
				Country:     addr.Country,
			}
			_, err := s.studentAddressRepo.Create(ctx, studentAddr)
			if err != nil {
				logger.Error("Failed to create student address", log.Fields{
					"error":        err.Error(),
					"student_id":   createdStudent.ID,
					"address_type": addr.AddressType,
					"created_by":   userID,
				})
				return err
			}
		}
	}

	if len(request.Documents) > 0 {
		for _, doc := range request.Documents {
			studentDoc := &models.StudentDocument{
				StudentID:      createdStudent.ID,
				DocumentType:   doc.DocumentType,
				DocumentNumber: doc.DocumentNumber,
				IssueDate:      doc.IssueDate,
				IssuePlace:     doc.IssuePlace,
				ExpiryDate:     doc.ExpiryDate,
				CountryOfIssue: doc.CountryOfIssue,
				HasChip:        doc.HasChip,
				Notes:          doc.Notes,
			}
			_, err := s.studentDocumentRepo.Create(ctx, studentDoc)
			if err != nil {
				logger.Error("Failed to create student document", log.Fields{
					"error":           err.Error(),
					"student_id":      createdStudent.ID,
					"document_type":   doc.DocumentType,
					"document_number": doc.DocumentNumber,
					"created_by":      userID,
				})
				return err
			}
		}
	}

	return nil
}

func (s *studentService) UpdateStudent(ctx context.Context, userID string, studentID string, request *models.UpdateStudentRequest) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function":   "Update Student",
		"student_id": studentID,
		"updated_by": userID,
	})

	logger.Info("Updating student")

	if request == nil {
		return common.ErrInvalidInput
	}
	updatedStudent := &models.Student{}

	// Apply updates to the student model
	if request.StudentCode != nil {
		updatedStudent.StudentCode = *request.StudentCode
	}
	if request.Fullname != nil {
		updatedStudent.Fullname = *request.Fullname
	}
	if request.DateOfBirth != nil {
		updatedStudent.DateOfBirth = *request.DateOfBirth
	}
	if request.Gender != nil {
		updatedStudent.Gender = *request.Gender
	}
	if request.FacultyID != nil {
		updatedStudent.FacultyID = *request.FacultyID
	}
	if request.Batch != nil {
		updatedStudent.Batch = *request.Batch
	}
	if request.Program != nil {
		updatedStudent.Program = *request.Program
	}
	if request.Address != nil {
		updatedStudent.Address = *request.Address
	}
	if request.Email != nil {
		updatedStudent.Email = *request.Email
	}
	if request.Phone != nil {
		updatedStudent.Phone = *request.Phone
	}
	if request.StatusID != nil {
		updatedStudent.StatusID = *request.StatusID
	}
	if request.ProgramID != nil {
		updatedStudent.ProgramID = *request.ProgramID
	}
	if request.Nationality != nil {
		updatedStudent.Nationality = *request.Nationality
	}

	updatedStudent, err := s.studentRepo.Update(ctx, studentID, updatedStudent)
	if err != nil {
		logger.Error("Failed to update student", log.Fields{
			"error":      err.Error(),
			"updated_by": userID,
		})
		return err
	}

	if request.Addresses != nil {
		for _, addr := range request.Addresses {
			studentAddr := &models.StudentAddress{
				StudentID:   updatedStudent.ID,
				AddressType: addr.AddressType,
				Street:      addr.Street,
				Ward:        addr.Ward,
				District:    addr.District,
				City:        addr.City,
				Country:     addr.Country,
			}
			err = s.studentAddressRepo.UpdatesByConditions(ctx, studentAddr, func(tx *gorm.DB) {
				tx.Where("student_id = ?", updatedStudent.ID)
			})
			if err != nil {
				logger.Error("Failed to update student address ", log.Fields{
					"error":      err.Error(),
					"updated_by": userID,
				})
				return err
			}
		}
	}

	// Handle documents if provided
	if request.Documents != nil {
		// Create new documents
		for _, doc := range request.Documents {
			studentDoc := &models.StudentDocument{
				StudentID:      updatedStudent.ID,
				DocumentType:   doc.DocumentType,
				DocumentNumber: doc.DocumentNumber,
				IssueDate:      doc.IssueDate,
				IssuePlace:     doc.IssuePlace,
				ExpiryDate:     doc.ExpiryDate,
				CountryOfIssue: doc.CountryOfIssue,
				HasChip:        doc.HasChip,
				Notes:          doc.Notes,
			}
			err = s.studentDocumentRepo.UpdatesByConditions(ctx, studentDoc, func(tx *gorm.DB) {
				tx.Where("student_id = ?", updatedStudent.ID)
			})
			if err != nil {
				logger.Error("Failed to update student document ", log.Fields{
					"error":      err.Error(),
					"updated_by": userID,
				})
				return err
			}
		}
	}

	return nil
}

func (s *studentService) DeleteStudentByID(ctx context.Context, userID string, studentID string) error {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function":   "DeleteStudentByID",
		"student_id": studentID,
		"deleted_by": userID,
	})

	logger.Info("Deleting student")

	err := s.studentRepo.DeleteByID(ctx, studentID)
	if err != nil {
		logger.Error("Failed to delete student", log.Fields{
			"error": err.Error(),
		})
		return err
	}

	logger.Info("Student deleted successfully")
	return nil
}

func (s *studentService) GetStudentStatuses(ctx context.Context) ([]*models.StudentStatus, error) {
	studentStatus, err := s.studentStatusRepo.List(ctx, models2.QueryParams{}, func(tx *gorm.DB) {

	})
	if err != nil {
		return nil, err
	}
	return studentStatus, nil
}
