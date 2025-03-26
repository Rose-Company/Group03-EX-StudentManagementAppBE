package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	admin_models "Group03-EX-StudentManagementAppBE/internal/models/admin"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	student_status_models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"
	student_status "Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/services/gdrive"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Service interface {
	GetStudentByID(ctx context.Context, id string) (*models.StudentResponse, error)
	GetStudentList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error)
	CreateAStudent(ctx context.Context, userID string, req *models.CreateStudentRequest) error
	UpdateStudent(ctx context.Context, userID string, studentId string, req *models.UpdateStudentRequest) error
	DeleteStudentByID(ctx context.Context, userID string, studentID string) error
	GetStudentStatuses(ctx context.Context, req *student_status_models.ListStudentStatusRequest) ([]*models.StudentStatus, error)
	
	
	CreateStudentStatus(ctx context.Context, studentStatus *student_status_models.CreateStudentStatusRequest) error
	UpdateStudentStatus(ctx context.Context, id string, req *student_status_models.UpdateStudentStatusRequest) (*models.StudentStatus, error)
	DeleteStudentStatus(ctx context.Context, id string) error
	ImportStudentsFromFile(ctx context.Context, userID string, fileURL string) (*admin_models.ImportResult, error)
	ExportStudentsToCSV(ctx context.Context) ([]byte, error)
	ExportStudentsToJSON(ctx context.Context) ([]byte, error)
}

type studentService struct {
	studentRepo         student.Repository
	studentStatusRepo   student_status.Repository
	studentAddressRepo  student_addresses.Repository
	studentDocumentRepo student_identity_documents.Repository
	driveService        gdrive.Service
}

func NewStudentService(
	studentRepo student.Repository,
	studentStatusRepo student_status.Repository,
	studentAddressRepo student_addresses.Repository,
	studentDocumentRepo student_identity_documents.Repository,
	driveService gdrive.Service) Service {
	return &studentService{
		studentRepo:         studentRepo,
		studentStatusRepo:   studentStatusRepo,
		studentAddressRepo:  studentAddressRepo,
		studentDocumentRepo: studentDocumentRepo,
		driveService:        driveService,
	}
}

func (s *studentService) GetStudentByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	log.Printf("Fetching student details for ID: %s", id)
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
		log.Printf("Error fetching student with ID %s: %v", id, err)
		return nil, err
	}
	log.Printf("Successfully fetched student details for ID: %s", id)
	return student.ToResponse(), nil
}

func (s *studentService) GetStudentList(ctx context.Context, req *models.ListStudentRequest) (*models2.BaseListResponse, error) {
	log.Printf("Fetching student list with filters: %+v", req)
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
	log.Printf("Successfully fetched student list with %d records", len(students))
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
			createdStudent, err := s.studentAddressRepo.Create(ctx, studentAddr)
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
	log.Printf("Successfully fetched student list ")
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
	log.Printf("Student created successfully by user ID: %s", userID)
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
	log.Printf("Student deleted successfully with ID: %s by user ID: %s", studentID, userID)
	return nil
}

func (s *studentService) GetStudentStatuses(ctx context.Context, req *student_status_models.ListStudentStatusRequest) ([]*models.StudentStatus, error) {
	log.Printf("Fetching student statuses with filters: %+v", req)

	// Pass any filter conditions if needed
	var clauses []repositories.Clause
	if req.Name != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Name+"%")
		})
	}

	if req.Sort == "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Order("id ASC")
		})
	}

	// Combine clauses
	combinedClause := func(tx *gorm.DB) {
		for _, clause := range clauses {
			clause(tx)
		}
	}

	// Query with proper parameters
	studentStatus, err := s.studentStatusRepo.List(ctx, models2.QueryParams{}, combinedClause)
	if err != nil {
		log.Printf("Error fetching student statuses: %v", err)
		return nil, err
	}
	log.Printf("Successfully fetched %d student statuses", len(studentStatus))

	// Map studentStatus to the expected type
	var result []*models.StudentStatus
	for _, status := range studentStatus {
		result = append(result, &models.StudentStatus{
			ID:   status.ID,
			Name: status.Name,
		})
	}
	return result, nil
}

func (s *studentService) CreateStudentStatus(ctx context.Context, req *student_status_models.CreateStudentStatusRequest) error {

	studentStatus := &student_status_models.StudentStatus{
		Name: req.Name,
	}

	_, err := s.studentStatusRepo.Create(ctx, studentStatus)
	if err != nil {
		return err
	}
	log.Printf("Student status created successfully")
	return nil
}

func (s *studentService) UpdateStudentStatus(ctx context.Context, id string, req *student_status_models.UpdateStudentStatusRequest) (*models.StudentStatus, error) {
	log.Printf("Updating student status with ID: %s", id)
	studentStatus := &models.StudentStatus{
		Name: req.Name,
	}
	convertedStudentStatus := &student_status_models.StudentStatus{
		Name: studentStatus.Name,
	}
	updatedStudentStatus, err := s.studentStatusRepo.Update(ctx, id, convertedStudentStatus)
	if err != nil {
		log.Printf("Error updating student status with ID %s: %v", id, err)
		return nil, err
	}
	log.Printf("Student status updated successfully with ID: %s", id)
	return &models.StudentStatus{
		ID:   updatedStudentStatus.ID,
		Name: updatedStudentStatus.Name,
	}, nil
}

func (s *studentService) DeleteStudentStatus(ctx context.Context, id string) error {
	log.Printf("Deleting student status with ID: %s", id)
	if err := s.studentStatusRepo.DeleteByID(ctx, id); err != nil {
		log.Printf("Error deleting student status with ID %s: %v", id, err)
		return err
	}
	log.Printf("Student status deleted successfully with ID: %s", id)
	return nil
}

// Main import function that handles different file types
func (s *studentService) ImportStudentsFromFile(ctx context.Context, userID string, fileURL string) (*admin_models.ImportResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "ImportStudentsFromFile",
		"fileURL":  fileURL,
		"userID":   userID,
	})

	logger.Info("Starting student import from file")

	// Special handling for Google Drive URLs
	if strings.Contains(fileURL, "drive.google.com") || strings.Contains(fileURL, "docs.google.com") {
		logger.Info("Detected Google Drive URL, applying special handling")
	}

	// Download the file
	resp, err := http.Get(fileURL)
	if err != nil {
		logger.WithError(err).Error("Failed to download file from URL")
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.WithFields(log.Fields{
			"statusCode": resp.StatusCode,
		}).Error("Failed to download file, non-OK status code")
		return nil, fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	// Read the response body into a buffer for examination
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Error("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(bodyBytes) > 100 {
		logger.WithField("firstBytes", string(bodyBytes[:100])).Debug("First few bytes of file")
	}

	// Create a reader from the bytes for content analysis
	contentReader := bytes.NewReader(bodyBytes)

	// First try to determine file type from content
	fileType := common.DetermineFileTypeFromContent(contentReader)

	// If couldn't determine from content, try URL and headers
	if fileType == common.FILE_TYPE_UNKNOWN {
		fileType = common.DetermineFileTypeFromMetadata(fileURL, resp.Header)
	}

	logger.WithField("fileType", fileType).Info("Determined file type")

	// Reset the reader to the beginning for processing
	contentReader.Seek(0, io.SeekStart)

	// Process the file based on its type
	var result *admin_models.ImportResult
	var processErr error

	switch fileType {
	case common.FILE_TYPE_CSV:
		result, processErr = s.processImportFile(ctx, userID, contentReader, fileType, logger)
	case common.FILE_TYPE_JSON:
		result, processErr = s.processImportFile(ctx, userID, contentReader, fileType, logger)
	default:
		// For Google Drive, make one more attempt with CSV since it's common
		if strings.Contains(fileURL, "drive.google.com") || strings.Contains(fileURL, "docs.google.com") {
			contentReader.Seek(0, io.SeekStart)
			logger.Info("Trying CSV processing for Google Drive file")
			result, processErr = s.processImportFile(ctx, userID, contentReader, fileType, logger)
		} else {
			logger.Error("Unsupported file type")
			return nil, common.ErrInvalidFileFormat
		}
	}

	if processErr != nil {
		return nil, processErr
	}

	return result, nil
}

// Process files with goroutines for concurrent processing
func (s *studentService) processImportFile(ctx context.Context, userID string, reader io.Reader, fileType string, logger *log.Entry) (*admin_models.ImportResult, error) {
	// Parse input data based on file type
	var records []models.ImportRecord
	var err error

	switch fileType {
	case common.FILE_TYPE_CSV:
		records, err = s.parseCSVData(reader, logger)
	case common.FILE_TYPE_JSON:
		records, err = s.parseJSONData(reader, logger)
	default:
		return nil, common.ErrInvalidFileFormat
	}

	if err != nil {
		return nil, err
	}

	totalRecords := len(records)
	logger.WithFields(log.Fields{
		"fileType":     fileType,
		"totalRecords": totalRecords,
	}).Info("Starting concurrent processing of records")

	// Variables for tracking results
	var (
		successCount  int32 = 0
		errorCount    int32 = 0
		mu            sync.Mutex
		wg            sync.WaitGroup
		maxWorkers    = 10 // Adjust based on your system capabilities
		recordChan    = make(chan models.ImportRecord, maxWorkers)
		failedRecords = make([]admin_models.FailedRecordDetail, 0)
	)

	// Create worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					logger.WithField("recover", r).Error("Recovered from panic in worker goroutine")
					atomic.AddInt32(&errorCount, 1)
				}
				wg.Done()
			}()

			workerLogger := logger.WithField("workerID", workerID)

			for record := range recordChan {
				// If record already has an error from parsing, log it and continue
				if record.Err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: record.Index + 1,
						Error:     record.Err.Error(),
					})
					mu.Unlock()
					continue
				}

				// Double-check required fields
				if record.Data == nil || record.Data.StudentCode == nil || record.Data.Fullname == nil || record.Data.Email == nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: record.Index + 1,
						Error:     "Missing required fields after parsing",
					})
					mu.Unlock()
					workerLogger.WithField("recordIndex", record.Index).Warn("Missing required fields after parsing, skipping")
					continue
				}

				// Create student with panic protection
				var err error
				func() {
					defer func() {
						if r := recover(); r != nil {
							err = fmt.Errorf("panic while creating student: %v", r)
						}
					}()
					err = s.CreateAStudent(ctx, userID, record.Data)
				}()

				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber:   record.Index + 1,
						StudentCode: fmt.Sprintf("%d", *record.Data.StudentCode),
						Email:       *record.Data.Email,
						Error:       fmt.Sprintf("Database error: %v", err),
					})
					mu.Unlock()
					workerLogger.WithError(err).WithFields(log.Fields{
						"recordIndex": record.Index,
						"studentCode": *record.Data.StudentCode,
						"email":       *record.Data.Email,
					}).Warn("Error creating student, skipping")
					continue
				}

				atomic.AddInt32(&successCount, 1)
				workerLogger.WithFields(log.Fields{
					"recordIndex": record.Index,
					"studentCode": *record.Data.StudentCode,
					"email":       *record.Data.Email,
				}).Info("Successfully created student")
			}
		}(i)
	}

	// Send records to workers
	for i, record := range records {
		recordChan <- record

		// Log progress periodically
		if (i+1)%100 == 0 || i+1 == totalRecords {
			logger.WithFields(log.Fields{
				"progress": fmt.Sprintf("%d/%d", i+1, totalRecords),
				"percent":  fmt.Sprintf("%.1f%%", float64(i+1)/float64(totalRecords)*100),
			}).Info("Import progress")
		}
	}

	// Close channel when all records are sent
	close(recordChan)

	// Wait for all workers to finish
	wg.Wait()

	// Log completion
	logger.WithFields(log.Fields{
		"totalRecords": totalRecords,
		"successCount": successCount,
		"errorCount":   errorCount,
	}).Info("Completed processing file")

	// Sort failed records by row number for easier reading
	sort.Slice(failedRecords, func(i, j int) bool {
		return failedRecords[i].RowNumber < failedRecords[j].RowNumber
	})

	result := &admin_models.ImportResult{
		SuccessCount:  int(successCount),
		ErrorCount:    int(errorCount),
		FailedRecords: failedRecords,
	}

	return result, nil
}

// Update the parseCSVData function to use models.ImportRecord
func (s *studentService) parseCSVData(reader io.Reader, logger *log.Entry) ([]models.ImportRecord, error) {
	csvReader := csv.NewReader(reader)

	// Read headers
	headers, err := csvReader.Read()
	if err != nil {
		logger.WithError(err).Error("Failed to read CSV headers")
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Map column indices
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	// Verify required columns
	requiredColumns := []string{"studentcode", "fullname", "email"}
	for _, col := range requiredColumns {
		if _, ok := headerMap[col]; !ok {
			logger.WithField("missingColumn", col).Error("CSV is missing required column")
			return nil, fmt.Errorf("CSV is missing required column: %s", col)
		}
	}

	// Read all rows
	rows, err := csvReader.ReadAll()
	if err != nil {
		logger.WithError(err).Error("Failed to read CSV rows")
		return nil, fmt.Errorf("failed to read CSV rows: %w", err)
	}

	records := make([]models.ImportRecord, len(rows))
	for i, row := range rows {
		studentReq, err := s.createStudentRequestFromCSV(row, headerMap)
		records[i] = models.ImportRecord{
			Index: i + 1, // +1 for human-readable indexing
			Data:  studentReq,
			Err:   err,
		}
	}

	return records, nil
}

// Update the parseJSONData function to use models.ImportRecord
func (s *studentService) parseJSONData(reader io.Reader, logger *log.Entry) ([]models.ImportRecord, error) {
	var studentsData []map[string]interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&studentsData)
	if err != nil {
		logger.WithError(err).Error("Failed to decode JSON data")
		return nil, fmt.Errorf("failed to decode JSON data: %w", err)
	}

	records := make([]models.ImportRecord, len(studentsData))
	for i, data := range studentsData {
		if data == nil || len(data) == 0 {
			records[i] = models.ImportRecord{
				Index: i,
				Err:   fmt.Errorf("empty or nil record data"),
			}
			continue
		}

		studentReq, err := s.createStudentRequestFromJSON(data)
		records[i] = models.ImportRecord{
			Index: i,
			Data:  studentReq,
			Err:   err,
		}
	}

	return records, nil
}

// Helper functions for CSV parsing
func (s *studentService) getCSVColumnValue(row []string, headerMap map[string]int, name string) *string {
	if idx, ok := headerMap[name]; ok && idx < len(row) {
		value := strings.TrimSpace(row[idx])
		if value != "" {
			return &value
		}
	}
	return nil
}

func (s *studentService) parseCSVDate(row []string, headerMap map[string]int, name string) *time.Time {
	if idx, ok := headerMap[name]; ok && idx < len(row) {
		dateStr := strings.TrimSpace(row[idx])
		if dateStr != "" {
			formats := []string{
				"2006-01-02",
				"01/02/2006",
				"02/01/2006",
				"2006/01/02",
			}

			for _, format := range formats {
				date, err := time.Parse(format, dateStr)
				if err == nil {
					return &date
				}
			}
		}
	}
	return nil
}

func (s *studentService) parseCSVInt(row []string, headerMap map[string]int, name string) *int {
	if idx, ok := headerMap[name]; ok && idx < len(row) {
		numStr := strings.TrimSpace(row[idx])
		if numStr != "" {
			num, err := strconv.Atoi(numStr)
			if err == nil {
				return &num
			}
		}
	}
	return nil
}

func (s *studentService) parseCSVBool(row []string, headerMap map[string]int, name string) *bool {
	if idx, ok := headerMap[name]; ok && idx < len(row) {
		boolStr := strings.ToLower(strings.TrimSpace(row[idx]))
		if boolStr != "" {
			var boolVal bool
			if boolStr == "true" || boolStr == "yes" || boolStr == "1" {
				boolVal = true
				return &boolVal
			} else if boolStr == "false" || boolStr == "no" || boolStr == "0" {
				boolVal = false
				return &boolVal
			}
		}
	}
	return nil
}

func (s *studentService) processCSVAddresses(row []string, headerMap map[string]int) []*models.AddressRequest {
	var addresses []*models.AddressRequest
	addressPrefixes := []string{"permanent", "temporary", "mailing"}

	for _, prefix := range addressPrefixes {
		// Check if columns for this address type exist
		hasType := false
		for col := range headerMap {
			if strings.HasPrefix(col, prefix) {
				hasType = true
				break
			}
		}

		if hasType {
			addressType := strings.Title(prefix)
			address := s.createAddressFromCSV(row, headerMap, prefix, addressType)
			if address != nil && address.Street != "" && address.City != "" {
				addresses = append(addresses, address)
			}
		}
	}

	return addresses
}

func (s *studentService) createAddressFromCSV(row []string, headerMap map[string]int, prefix string, addressType string) *models.AddressRequest {
	getVal := func(field string) string {
		if val := s.getCSVColumnValue(row, headerMap, prefix+field); val != nil {
			return *val
		}
		return ""
	}

	address := &models.AddressRequest{
		AddressType: addressType,
		Street:      getVal("street"),
		Ward:        getVal("ward"),
		District:    getVal("district"),
		City:        getVal("city"),
		Country:     getVal("country"),
	}

	if address.Country == "" {
		address.Country = "Vietnam" // Default country
	}

	return address
}

func (s *studentService) processCSVDocuments(row []string, headerMap map[string]int) []*models.DocumentRequest {
	var documents []*models.DocumentRequest
	documentPrefixes := []string{"cccd", "cmnd", "passport"}

	for _, prefix := range documentPrefixes {
		// Check if columns for this document type exist
		hasType := false
		for col := range headerMap {
			if strings.HasPrefix(col, prefix) {
				hasType = true
				break
			}
		}

		if hasType {
			document := s.createDocumentFromCSV(row, headerMap, prefix)
			if document != nil && document.DocumentNumber != "" {
				documents = append(documents, document)
			}
		}
	}

	return documents
}

func (s *studentService) createDocumentFromCSV(row []string, headerMap map[string]int, prefix string) *models.DocumentRequest {
	documentNumber := s.getCSVColumnValue(row, headerMap, prefix+"number")
	if documentNumber == nil {
		return nil
	}

	document := &models.DocumentRequest{
		DocumentType:   strings.ToUpper(prefix),
		DocumentNumber: *documentNumber,
		IssueDate:      time.Time{},
		ExpiryDate:     time.Time{},
		CountryOfIssue: func() string {
			if val := s.getCSVColumnValue(row, headerMap, prefix+"country"); val != nil {
				return *val
			}
			return "Vietnam"
		}(),
		HasChip: func() bool {
			if val := s.parseCSVBool(row, headerMap, prefix+"haschip"); val != nil {
				return *val
			}
			return false
		}(),
	}

	if date := s.parseCSVDate(row, headerMap, prefix+"issuedate"); date != nil {
		document.IssueDate = *date
	}
	if date := s.parseCSVDate(row, headerMap, prefix+"expirydate"); date != nil {
		document.ExpiryDate = *date
	}
	if place := s.getCSVColumnValue(row, headerMap, prefix+"issueplace"); place != nil {
		document.IssuePlace = *place
	}
	document.Notes = s.getCSVColumnValue(row, headerMap, prefix+"notes")

	return document
}

// Helper functions for JSON parsing
func (s *studentService) getJSONString(data map[string]interface{}, key string) *string {
	possibleKeys := []string{
		key,
		strings.ToUpper(key[:1]) + key[1:],
		strings.ToLower(key),
	}

	for _, k := range possibleKeys {
		if val, ok := data[k]; ok && val != nil {
			if strVal, ok := val.(string); ok && strVal != "" {
				return &strVal
			} else if numVal, ok := val.(float64); ok {
				strVal := fmt.Sprintf("%v", numVal)
				return &strVal
			}
		}
	}
	return nil
}

func (s *studentService) getJSONInt(data map[string]interface{}, key string) *int {
	possibleKeys := []string{
		key,
		strings.ToUpper(key[:1]) + key[1:],
		strings.ToLower(key),
	}

	for _, k := range possibleKeys {
		if val, ok := data[k]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				intVal := int(v)
				return &intVal
			case int:
				return &v
			case string:
				intVal, err := strconv.Atoi(v)
				if err == nil {
					return &intVal
				}
			}
		}
	}

	if strings.ToLower(key) == "statusid" {
		defaultStatus := 1
		return &defaultStatus
	}

	return nil
}

func (s *studentService) getJSONDate(data map[string]interface{}, key string) *time.Time {
	possibleKeys := []string{
		key,
		strings.ToUpper(key[:1]) + key[1:],
		strings.ToLower(key),
	}

	for _, k := range possibleKeys {
		if val, ok := data[k]; ok && val != nil {
			if strVal, ok := val.(string); ok && strVal != "" {
				formats := []string{
					"2006-01-02",
					"01/02/2006",
					"02/01/2006",
					"2006/01/02",
				}

				for _, format := range formats {
					date, err := time.Parse(format, strVal)
					if err == nil {
						return &date
					}
				}
			}
		}
	}

	if strings.ToLower(key) == "dateofbirth" {
		now := time.Now()
		return &now
	}

	return nil
}

func (s *studentService) processJSONAddresses(data map[string]interface{}) []*models.AddressRequest {
	addressTypes := map[string]struct{}{
		"Permanent": {},
		"Temporary": {},
		"Mailing":   {},
	}

	addressMap := make(map[string]*models.AddressRequest)

	for key, value := range data {
		if value == nil {
			continue
		}

		for prefix := range addressTypes {
			if strings.HasPrefix(key, prefix) && len(key) > len(prefix) {
				field := key[len(prefix):]
				s.updateAddressField(addressMap, prefix, field, value)
			}
		}
	}

	var addresses []*models.AddressRequest
	for _, addr := range addressMap {
		if addr != nil && addr.Street != "" && addr.City != "" {
			addresses = append(addresses, addr)
		}
	}

	return addresses
}

func (s *studentService) updateAddressField(addressMap map[string]*models.AddressRequest, prefix, field string, value interface{}) {
	if _, exists := addressMap[prefix]; !exists {
		addressMap[prefix] = &models.AddressRequest{
			AddressType: prefix,
			Country:     "Vietnam", // Default
		}
	}

	if value == nil {
		return
	}

	strVal, ok := value.(string)
	if !ok || strVal == "" {
		return
	}

	addr := addressMap[prefix]
	switch {
	case strings.EqualFold(field, "Street"):
		addr.Street = strVal
	case strings.EqualFold(field, "Ward"):
		addr.Ward = strVal
	case strings.EqualFold(field, "District"):
		addr.District = strVal
	case strings.EqualFold(field, "City"):
		addr.City = strVal
	case strings.EqualFold(field, "Country"):
		addr.Country = strVal
	}
}

func (s *studentService) processJSONDocuments(data map[string]interface{}) []*models.DocumentRequest {
	docPrefixes := []string{"CCCD", "CMND", "Passport"}
	docMap := make(map[string]*models.DocumentRequest)

	for key, value := range data {
		if value == nil {
			continue
		}

		for _, prefix := range docPrefixes {
			if strings.HasPrefix(key, prefix) && len(key) > len(prefix) {
				field := key[len(prefix):]
				s.updateDocumentField(docMap, prefix, field, value)
			}
		}
	}

	var documents []*models.DocumentRequest
	for _, doc := range docMap {
		if doc != nil && doc.DocumentNumber != "" {
			documents = append(documents, doc)
		}
	}

	return documents
}

func (s *studentService) updateDocumentField(docMap map[string]*models.DocumentRequest, prefix, field string, value interface{}) {
	if _, exists := docMap[prefix]; !exists {
		docMap[prefix] = &models.DocumentRequest{
			DocumentType:   prefix,
			CountryOfIssue: "Vietnam", // Default
		}
	}

	if value == nil {
		return
	}

	doc := docMap[prefix]
	switch {
	case strings.EqualFold(field, "Number"):
		if numVal, ok := value.(float64); ok {
			doc.DocumentNumber = fmt.Sprintf("%v", int(numVal))
		} else if strVal, ok := value.(string); ok && strVal != "" {
			doc.DocumentNumber = strVal
		}
	case strings.EqualFold(field, "IssueDate"), strings.EqualFold(field, "ExpiryDate"):
		if strVal, ok := value.(string); ok && strVal != "" {
			formats := []string{"2006-01-02", "01/02/2006", "02/01/2006", "2006/01/02"}
			for _, format := range formats {
				date, err := time.Parse(format, strVal)
				if err == nil {
					if strings.EqualFold(field, "IssueDate") {
						doc.IssueDate = date
					} else {
						doc.ExpiryDate = date
					}
					break
				}
			}
		}
	case strings.EqualFold(field, "IssuePlace"):
		if strVal, ok := value.(string); ok && strVal != "" {
			doc.IssuePlace = strVal
		}
	case strings.EqualFold(field, "Country"):
		if strVal, ok := value.(string); ok && strVal != "" {
			doc.CountryOfIssue = strVal
		}
	case strings.EqualFold(field, "HasChip"):
		if boolVal, ok := value.(bool); ok {
			doc.HasChip = boolVal
		}
	case strings.EqualFold(field, "Notes"):
		if strVal, ok := value.(string); ok && strVal != "" {
			doc.Notes = &strVal
		} else if value == nil {
			emptyNote := ""
			doc.Notes = &emptyNote
		}
	}
}

// Main functions that use the helpers
func (s *studentService) createStudentRequestFromCSV(row []string, headerMap map[string]int) (*models.CreateStudentRequest, error) {
	req := &models.CreateStudentRequest{
		StudentCode: s.parseCSVInt(row, headerMap, "studentcode"),
		Fullname:    s.getCSVColumnValue(row, headerMap, "fullname"),
		Email:       s.getCSVColumnValue(row, headerMap, "email"),
		DateOfBirth: s.parseCSVDate(row, headerMap, "dateofbirth"),
		Gender:      s.getCSVColumnValue(row, headerMap, "gender"),
		FacultyID:   s.parseCSVInt(row, headerMap, "facultyid"),
		Batch:       s.getCSVColumnValue(row, headerMap, "batch"),
		Program:     s.getCSVColumnValue(row, headerMap, "program"),
		Address:     s.getCSVColumnValue(row, headerMap, "address"),
		Phone:       s.getCSVColumnValue(row, headerMap, "phone"),
		StatusID:    s.parseCSVInt(row, headerMap, "statusid"),
		ProgramID:   s.parseCSVInt(row, headerMap, "programid"),
		Nationality: s.getCSVColumnValue(row, headerMap, "nationality"),
	}

	// Process addresses and documents
	req.Addresses = s.processCSVAddresses(row, headerMap)
	req.Documents = s.processCSVDocuments(row, headerMap)

	// Validate required fields
	if req.StudentCode == nil || req.Fullname == nil || req.Email == nil {
		return nil, fmt.Errorf("missing required fields: studentCode, fullname, or email")
	}

	// Set default status if not provided
	if req.StatusID == nil {
		defaultStatus := 1
		req.StatusID = &defaultStatus
	}

	return req, nil
}

func (s *studentService) createStudentRequestFromJSON(data map[string]interface{}) (*models.CreateStudentRequest, error) {
	// Debug logging
	jsonBytes, _ := json.Marshal(data)
	log.WithField("data", string(jsonBytes)).Debug("Processing JSON record")

	req := &models.CreateStudentRequest{
		StudentCode: s.getJSONInt(data, "studentCode"),
		Fullname:    s.getJSONString(data, "fullname"),
		Email:       s.getJSONString(data, "email"),
		DateOfBirth: s.getJSONDate(data, "dateOfBirth"),
		Gender:      s.getJSONString(data, "gender"),
		FacultyID:   s.getJSONInt(data, "facultyId"),
		Batch:       s.getJSONString(data, "batch"),
		Program:     s.getJSONString(data, "program"),
		Address:     s.getJSONString(data, "address"),
		Phone:       s.getJSONString(data, "phone"),
		StatusID:    s.getJSONInt(data, "statusId"),
		ProgramID:   s.getJSONInt(data, "programId"),
		Nationality: s.getJSONString(data, "nationality"),
	}

	// Process addresses and documents
	req.Addresses = s.processJSONAddresses(data)
	req.Documents = s.processJSONDocuments(data)

	// Validate and set defaults
	if err := s.validateAndSetDefaults(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *studentService) validateAndSetDefaults(req *models.CreateStudentRequest) error {
	// Validate required fields
	if req.StudentCode == nil {
		return fmt.Errorf("missing required field: studentCode")
	}
	if req.Fullname == nil {
		return fmt.Errorf("missing required field: fullname")
	}
	if req.Email == nil {
		return fmt.Errorf("missing required field: email")
	}

	// Set defaults for optional fields
	if req.Gender == nil {
		defaultGender := "Other"
		req.Gender = &defaultGender
	}
	if req.Batch == nil {
		currentYear := fmt.Sprintf("%d", time.Now().Year())
		req.Batch = &currentYear
	}
	if req.Program == nil {
		defaultProgram := "Unknown"
		req.Program = &defaultProgram
	}
	if req.Address == nil {
		defaultAddress := ""
		req.Address = &defaultAddress
	}
	if req.Phone == nil {
		defaultPhone := ""
		req.Phone = &defaultPhone
	}
	if req.ProgramID == nil {
		defaultProgramID := 1
		req.ProgramID = &defaultProgramID
	}
	if req.FacultyID == nil {
		defaultFacultyID := 1
		req.FacultyID = &defaultFacultyID
	}
	if req.Nationality == nil {
		defaultNationality := "Vietnam"
		req.Nationality = &defaultNationality
	}

	return nil
}

// ExportStudentsToCSV exports student data as CSV
func (s *studentService) ExportStudentsToCSV(ctx context.Context) ([]byte, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "ExportStudentsToCSV",
		"format":   "CSV",
	})
	logger.Info("Starting student export to CSV")

	// Get students from the database
	students, err := s.fetchStudentsWithDetails(ctx)
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the CSV data
	var csvBuffer bytes.Buffer
	writer := csv.NewWriter(&csvBuffer)

	// Write the header row
	header := []string{
		"Student Code", "Full Name", "Email", "Date of Birth", "Gender",
		"Faculty ID", "Faculty Name", "Batch", "Program", "Address", "Phone", "Status",
		"Nationality", "Permanent Address", "Temporary Address",
	}
	if err := writer.Write(header); err != nil {
		logger.WithError(err).Error("Failed to write CSV header")
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Loop through students and write each as a row
	for _, student := range students {
		// Format addresses
		permanentAddress, temporaryAddress := s.formatStudentAddresses(student)

		// Get faculty information
		facultyID, facultyName := s.getStudentFacultyInfo(student)

		// Get status name
		statusName := s.getStudentStatusName(ctx, student)

		// Format the date of birth
		dob := ""
		if !student.DateOfBirth.IsZero() {
			dob = student.DateOfBirth.Format("2006-01-02")
		}

		// Write the student data row
		row := []string{
			fmt.Sprintf("%d", student.StudentCode),
			student.Fullname,
			student.Email,
			dob,
			student.Gender,
			fmt.Sprintf("%d", facultyID),
			facultyName,
			student.Batch,
			student.Program,
			student.Address,
			student.Phone,
			statusName,
			student.Nationality,
			permanentAddress,
			temporaryAddress,
		}

		if err := writer.Write(row); err != nil {
			logger.WithError(err).Error("Failed to write student row to CSV")
			continue // Skip this student and continue with others
		}
	}

	// Flush data to ensure it's written to the buffer
	writer.Flush()
	if err := writer.Error(); err != nil {
		logger.WithError(err).Error("Error flushing CSV data")
		return nil, fmt.Errorf("error flushing CSV data: %w", err)
	}

	// Get the CSV data as bytes
	csvData := csvBuffer.Bytes()
	logger.WithField("csvSize", len(csvData)).Info("Generated CSV data in memory")

	return csvData, nil
}

// ExportStudentsToJSON exports student data as JSON
func (s *studentService) ExportStudentsToJSON(ctx context.Context) ([]byte, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "ExportStudentsToJSON",
		"format":   "JSON",
	})
	logger.Info("Starting student export to JSON")

	// Get students from the database
	students, err := s.fetchStudentsWithDetails(ctx)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the student data in a format suitable for JSON
	type StudentExport struct {
		StudentCode      int    `json:"student_code"`
		FullName         string `json:"full_name"`
		Email            string `json:"email"`
		DateOfBirth      string `json:"date_of_birth"`
		Gender           string `json:"gender"`
		FacultyID        int    `json:"faculty_id"`
		FacultyName      string `json:"faculty_name"`
		Batch            string `json:"batch"`
		Program          string `json:"program"`
		Address          string `json:"address"`
		Phone            string `json:"phone"`
		Status           string `json:"status"`
		Nationality      string `json:"nationality"`
		PermanentAddress string `json:"permanent_address"`
		TemporaryAddress string `json:"temporary_address"`
	}

	var studentsData []StudentExport

	// Format each student's data
	for _, student := range students {
		// Format addresses
		permanentAddress, temporaryAddress := s.formatStudentAddresses(student)

		// Get faculty information
		facultyID, facultyName := s.getStudentFacultyInfo(student)

		// Get status name
		statusName := s.getStudentStatusName(ctx, student)

		// Format the date of birth
		dob := ""
		if !student.DateOfBirth.IsZero() {
			dob = student.DateOfBirth.Format("2006-01-02")
		}

		// Add student data to the slice
		studentsData = append(studentsData, StudentExport{
			StudentCode:      student.StudentCode,
			FullName:         student.Fullname,
			Email:            student.Email,
			DateOfBirth:      dob,
			Gender:           student.Gender,
			FacultyID:        facultyID,
			FacultyName:      facultyName,
			Batch:            student.Batch,
			Program:          student.Program,
			Address:          student.Address,
			Phone:            student.Phone,
			Status:           statusName,
			Nationality:      student.Nationality,
			PermanentAddress: permanentAddress,
			TemporaryAddress: temporaryAddress,
		})
	}

	// Create the response structure
	response := map[string]interface{}{
		"success": true,
		"data":    studentsData,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal JSON data")
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	logger.WithField("jsonSize", len(jsonData)).Info("Generated JSON data in memory")
	return jsonData, nil
}

// Helper functions to reduce code duplication

// fetchStudentsWithDetails gets all students with their related data
func (s *studentService) fetchStudentsWithDetails(ctx context.Context) ([]*models.Student, error) {
	logger := log.WithContext(ctx)

	students, err := s.studentRepo.List(ctx, models2.QueryParams{
		Limit: -1, // Get all students
	}, func(tx *gorm.DB) {
		tx.Joins(`LEFT JOIN "PUBLIC"."faculties" ON students.faculty_id = "PUBLIC"."faculties".id`)
		tx.Select(`students.*, "PUBLIC"."faculties".name as faculty_name`)
		tx.Preload("Addresses") // Preload addresses
	})

	if err != nil {
		logger.WithError(err).Error("Failed to retrieve students from database")
		return nil, fmt.Errorf("failed to retrieve students: %w", err)
	}

	return students, nil
}

// formatStudentAddresses extracts permanent and temporary addresses
func (s *studentService) formatStudentAddresses(student *models.Student) (string, string) {
	permanentAddress := ""
	temporaryAddress := ""

	if student.Addresses != nil {
		for _, addr := range student.Addresses {
			addrStr := fmt.Sprintf("%s, %s, %s, %s, %s",
				addr.Street, addr.Ward, addr.District, addr.City, addr.Country)

			if addr.AddressType == "Permanent" {
				permanentAddress = addrStr
			} else if addr.AddressType == "Temporary" {
				temporaryAddress = addrStr
			}
		}
	}

	return permanentAddress, temporaryAddress
}

// getStudentFacultyInfo extracts faculty ID and name
func (s *studentService) getStudentFacultyInfo(student *models.Student) (int, string) {
	facultyID := 0
	if student.FacultyID != 0 {
		facultyID = student.FacultyID
	}

	facultyName := student.FacultyName
	if facultyName == "" {
		facultyName = "Unknown" // Default if not available
	}

	return facultyID, facultyName
}

func (s *studentService) getStudentStatusName(ctx context.Context, student *models.Student) string {
	statusName := "Active"
	if student.StatusID != 1 {
		status, err := s.studentStatusRepo.GetByID(ctx, fmt.Sprintf("%d", student.StatusID))
		if err == nil && status != nil {
			statusName = status.Name
		}
	}
	return statusName
}
