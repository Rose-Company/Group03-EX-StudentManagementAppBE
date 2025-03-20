package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	admin_models "Group03-EX-StudentManagementAppBE/internal/models/admin"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/services/gdrive"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	GetStudentStatuses(ctx context.Context) ([]*models.StudentStatus, error)
	ImportStudentsFromFile(ctx context.Context, userID string, fileURL string) (*admin_models.ImportResult, error)
	ExportStudentsToCSV(ctx context.Context) (string, error)
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
	fileType := determineFileTypeFromContent(contentReader)

	// If couldn't determine from content, try URL and headers
	if fileType == "unknown" {
		fileType = determineFileTypeFromMetadata(fileURL, resp.Header)
	}

	logger.WithField("fileType", fileType).Info("Determined file type")

	// Reset the reader to the beginning for processing
	contentReader.Seek(0, io.SeekStart)

	// Process the file based on its type
	var result *admin_models.ImportResult
	var processErr error

	switch fileType {
	case "csv":
		result, processErr = s.processCSVFile(ctx, userID, contentReader, logger)
	case "json":
		result, processErr = s.processJSONFile(ctx, userID, contentReader, logger)
	default:
		// For Google Drive, make one more attempt with CSV since it's common
		if strings.Contains(fileURL, "drive.google.com") || strings.Contains(fileURL, "docs.google.com") {
			contentReader.Seek(0, io.SeekStart)
			logger.Info("Trying CSV processing for Google Drive file")
			result, processErr = s.processCSVFile(ctx, userID, contentReader, logger)
		} else {
			logger.Error("Unsupported file type")
			return nil, common.ErrInvalidFormat
		}
	}

	if processErr != nil {
		return nil, processErr
	}

	return result, nil
}

// determineFileTypeFromContent analyzes the content to determine file type
func determineFileTypeFromContent(reader io.ReadSeeker) string {
	// Save original position
	currentPosition, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return "unknown"
	}
	defer reader.Seek(currentPosition, io.SeekStart) // Restore position afterward

	// Read first 1024 bytes to analyze
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return "unknown"
	}

	sample := buf[:n]

	// Check if it's JSON by looking for JSON structure
	trimmedSample := bytes.TrimSpace(sample)
	if len(trimmedSample) > 0 {
		firstChar := trimmedSample[0]
		if (firstChar == '{' && bytes.Contains(trimmedSample, []byte{':'})) ||
			(firstChar == '[' && bytes.Contains(trimmedSample, []byte{'{'})) {
			return "json"
		}
	}

	// Check if it looks like CSV by looking for comma-separated values and newlines
	// Count commas and newlines to ensure it's consistent with CSV format
	commaCount := bytes.Count(sample, []byte{','})
	newlineCount := bytes.Count(sample, []byte{'\n'})

	// Only consider it CSV if there are reasonable comma counts per line (at least one comma)
	// and more than one line
	if newlineCount > 0 && commaCount > 0 && commaCount/newlineCount >= 1 {
		// Check if first line looks like a header (no numeric values)
		lines := bytes.Split(sample, []byte{'\n'})
		if len(lines) > 0 {
			firstLine := lines[0]
			// Check if first line has commas and doesn't look like numeric data
			if bytes.Contains(firstLine, []byte{','}) &&
				!bytes.ContainsAny(firstLine, "0123456789") {
				return "csv"
			}
		}

		// Even without clear header, if it has commas and lines, likely CSV
		return "csv"
	}

	// Default
	return "unknown"
}

// determineFileTypeFromMetadata examines URL and headers for file type clues
func determineFileTypeFromMetadata(fileURL string, headers http.Header) string {
	// Try to determine from URL
	lowerURL := strings.ToLower(fileURL)
	if strings.HasSuffix(lowerURL, ".csv") {
		return "csv"
	} else if strings.HasSuffix(lowerURL, ".json") {
		return "json"
	}

	// Look for content disposition header which might have filename
	contentDisposition := headers.Get("Content-Disposition")
	if contentDisposition != "" {
		if strings.Contains(strings.ToLower(contentDisposition), ".csv") {
			return "csv"
		} else if strings.Contains(strings.ToLower(contentDisposition), ".json") {
			return "json"
		}
	}

	// Check content type
	contentType := headers.Get("Content-Type")
	if strings.Contains(contentType, "csv") || strings.Contains(contentType, "text/comma-separated-values") {
		return "csv"
	} else if strings.Contains(contentType, "json") || strings.Contains(contentType, "application/json") {
		return "json"
	} else if strings.Contains(contentType, "text/plain") {
		// Many CSVs are served as text/plain
		return "csv"
	}

	// Default
	return "unknown"
}

// Process CSV files with goroutines for concurrent processing
func (s *studentService) processCSVFile(ctx context.Context, userID string, reader io.Reader, logger *log.Entry) (*admin_models.ImportResult, error) {
	// Create CSV reader
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

	// Verify required columns - only check essential columns
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

	// Setup concurrency controls
	totalRows := len(rows)
	logger.WithField("totalRows", totalRows).Info("Starting concurrent processing of student records")

	// Variables for tracking results
	var (
		successCount int32 = 0
		errorCount   int32 = 0
		mu           sync.Mutex
		wg           sync.WaitGroup
		maxWorkers   = 10 // Adjust based on your system capabilities
		rowChan      = make(chan struct {
			rowNum int
			row    []string
		}, maxWorkers)
		failedRecords = make([]admin_models.FailedRecordDetail, 0)
	)

	// Create worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			workerLogger := logger.WithField("workerID", workerID)

			for job := range rowChan {
				rowNum := job.rowNum
				row := job.row

				// Create student request from row
				studentReq, err := s.createStudentRequestFromCSV(row, headerMap)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: rowNum,
						Error:     fmt.Sprintf("Failed to parse row: %v", err),
					})
					mu.Unlock()
					workerLogger.WithError(err).WithField("rowNum", rowNum).Warn("Error creating student from CSV row, skipping")
					continue
				}

				// Create student
				err = s.CreateAStudent(ctx, userID, studentReq)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber:   rowNum,
						StudentCode: fmt.Sprintf("%d", *studentReq.StudentCode),
						Email:       *studentReq.Email,
						Error:       fmt.Sprintf("Database error: %v", err),
					})
					mu.Unlock()
					workerLogger.WithError(err).WithFields(log.Fields{
						"rowNum":      rowNum,
						"studentCode": *studentReq.StudentCode,
						"email":       *studentReq.Email,
					}).Warn("Error creating student, skipping")
					continue
				}

				atomic.AddInt32(&successCount, 1)
				workerLogger.WithFields(log.Fields{
					"rowNum":      rowNum,
					"studentCode": *studentReq.StudentCode,
					"email":       *studentReq.Email,
				}).Info("Successfully created student")
			}
		}(i)
	}

	// Send rows to workers
	for i, row := range rows {
		rowNum := i + 2 // +2 because row 0 is header and we're 1-indexed for human readability
		rowChan <- struct {
			rowNum int
			row    []string
		}{rowNum: rowNum, row: row}

		// Log progress periodically
		if (i+1)%100 == 0 || i+1 == totalRows {
			logger.WithFields(log.Fields{
				"progress": fmt.Sprintf("%d/%d", i+1, totalRows),
				"percent":  fmt.Sprintf("%.1f%%", float64(i+1)/float64(totalRows)*100),
			}).Info("Import progress")
		}
	}

	// Close channel when all rows are sent
	close(rowChan)

	// Wait for all workers to finish
	wg.Wait()

	// Log completion
	logger.WithFields(log.Fields{
		"totalRows":    totalRows,
		"successCount": successCount,
		"errorCount":   errorCount,
	}).Info("Completed processing CSV file")

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

// Process JSON files with goroutines for concurrent processing
// Refactored JSON processing function with improved error handling
func (s *studentService) processJSONFile(ctx context.Context, userID string, reader io.Reader, logger *log.Entry) (*admin_models.ImportResult, error) {
	// Read JSON data
	var studentsData []map[string]interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&studentsData)
	if err != nil {
		logger.WithError(err).Error("Failed to decode JSON data")
		return nil, fmt.Errorf("failed to decode JSON data: %w", err)
	}

	totalRecords := len(studentsData)
	logger.WithField("totalRecords", totalRecords).Info("Starting concurrent processing of JSON records")

	// Variables for tracking results
	var (
		successCount int32 = 0
		errorCount   int32 = 0
		mu           sync.Mutex
		wg           sync.WaitGroup
		maxWorkers   = 10 // Adjust based on your system capabilities
		dataChan     = make(chan struct {
			index int
			data  map[string]interface{}
		}, maxWorkers)
		failedRecords = make([]admin_models.FailedRecordDetail, 0)
	)

	// Create worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				// Recover from panics in goroutines
				if r := recover(); r != nil {
					logger.WithField("recover", r).Error("Recovered from panic in worker goroutine")
					atomic.AddInt32(&errorCount, 1)
				}
				wg.Done()
			}()

			workerLogger := logger.WithField("workerID", workerID)

			for job := range dataChan {
				index := job.index
				data := job.data

				// Skip nil or empty data
				if data == nil || len(data) == 0 {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: index + 1, // +1 for human-readable indexing
						Error:     "Empty or nil record data",
					})
					mu.Unlock()
					continue
				}

				// Convert JSON record to CreateStudentRequest with panic protection
				var studentReq *models.CreateStudentRequest
				var err error

				func() {
					defer func() {
						if r := recover(); r != nil {
							err = fmt.Errorf("panic while parsing JSON data: %v", r)
						}
					}()
					studentReq, err = s.createStudentRequestFromJSON(data)
				}()

				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: index + 1, // +1 for human-readable indexing
						Error:     fmt.Sprintf("Failed to parse data: %v", err),
					})
					mu.Unlock()
					workerLogger.WithError(err).WithField("recordIndex", index).Warn("Error creating student from JSON data, skipping")
					continue
				}

				// Double-check required fields before proceeding
				if studentReq == nil || studentReq.StudentCode == nil || studentReq.Fullname == nil || studentReq.Email == nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber: index + 1,
						Error:     "Missing required fields after parsing",
					})
					mu.Unlock()
					workerLogger.WithField("recordIndex", index).Warn("Missing required fields after parsing JSON data, skipping")
					continue
				}

				// Create student with panic protection
				func() {
					defer func() {
						if r := recover(); r != nil {
							err = fmt.Errorf("panic while creating student: %v", r)
						}
					}()
					err = s.CreateAStudent(ctx, userID, studentReq)
				}()

				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					mu.Lock()
					failedRecords = append(failedRecords, admin_models.FailedRecordDetail{
						RowNumber:   index + 1,
						StudentCode: fmt.Sprintf("%d", *studentReq.StudentCode),
						Email:       *studentReq.Email,
						Error:       fmt.Sprintf("Database error: %v", err),
					})
					mu.Unlock()
					workerLogger.WithError(err).WithFields(log.Fields{
						"recordIndex": index,
						"studentCode": *studentReq.StudentCode,
						"email":       *studentReq.Email,
					}).Warn("Error creating student, skipping")
					continue
				}

				atomic.AddInt32(&successCount, 1)
				workerLogger.WithFields(log.Fields{
					"recordIndex": index,
					"studentCode": *studentReq.StudentCode,
					"email":       *studentReq.Email,
				}).Info("Successfully created student")
			}
		}(i)
	}

	// Send data to workers
	for i, data := range studentsData {
		dataChan <- struct {
			index int
			data  map[string]interface{}
		}{index: i, data: data}

		// Log progress periodically
		if (i+1)%100 == 0 || i+1 == totalRecords {
			logger.WithFields(log.Fields{
				"progress": fmt.Sprintf("%d/%d", i+1, totalRecords),
				"percent":  fmt.Sprintf("%.1f%%", float64(i+1)/float64(totalRecords)*100),
			}).Info("Import progress")
		}
	}

	// Close channel when all data is sent
	close(dataChan)

	// Wait for all workers to finish
	wg.Wait()

	// Log completion
	logger.WithFields(log.Fields{
		"totalRecords": totalRecords,
		"successCount": successCount,
		"errorCount":   errorCount,
	}).Info("Completed processing JSON file")

	// Sort failed records by index for easier reading
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

// Helper function to create student request from CSV row
func (s *studentService) createStudentRequestFromCSV(row []string, headerMap map[string]int) (*models.CreateStudentRequest, error) {
	// Helper function to safely get column value
	getCol := func(name string) *string {
		if idx, ok := headerMap[name]; ok && idx < len(row) {
			value := strings.TrimSpace(row[idx])
			if value != "" {
				return &value
			}
		}
		return nil
	}

	// Helper function to parse date
	parseDate := func(name string) *time.Time {
		if idx, ok := headerMap[name]; ok && idx < len(row) {
			dateStr := strings.TrimSpace(row[idx])
			if dateStr != "" {
				// Try different date formats
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

	// Helper function to parse int
	parseInt := func(name string) *int {
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

	// Helper function to parse bool
	parseBool := func(name string) *bool {
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

	// Create student request
	req := &models.CreateStudentRequest{
		StudentCode: parseInt("studentcode"),
		Fullname:    getCol("fullname"),
		Email:       getCol("email"),
		DateOfBirth: parseDate("dateofbirth"),
		Gender:      getCol("gender"),
		FacultyID:   parseInt("facultyid"),
		Batch:       getCol("batch"),
		Program:     getCol("program"),
		Address:     getCol("address"),
		Phone:       getCol("phone"),
		StatusID:    parseInt("statusid"),
		ProgramID:   parseInt("programid"),
		Nationality: getCol("nationality"),
		Addresses:   nil,
		Documents:   nil,
	}

	// Process addresses if columns exist
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

			address := &models.AddressRequest{
				AddressType: addressType,
				Street: func() string {
					if val := getCol(prefix + "street"); val != nil {
						return *val
					}
					return ""
				}(),
				Ward: func() string {
					if val := getCol(prefix + "ward"); val != nil {
						return *val
					}
					return ""
				}(),
				District: func() string {
					if val := getCol(prefix + "district"); val != nil {
						return *val
					}
					return ""
				}(),
				City: func() string {
					if val := getCol(prefix + "city"); val != nil {
						return *val
					}
					return ""
				}(),
				Country: func() string {
					if val := getCol(prefix + "country"); val != nil {
						return *val
					}
					return "Vietnam" // Default country
				}(),
			}

			// Only add if we have at least street and city
			if address.Street != "" && address.City != "" {
				req.Addresses = append(req.Addresses, address)
			}
		}
	}

	// Process documents if columns exist
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
			documentType := strings.ToUpper(prefix)

			documentNumber := getCol(prefix + "number")
			if documentNumber != nil {
				document := &models.DocumentRequest{
					DocumentType:   documentType,
					DocumentNumber: *documentNumber,
					IssueDate: func() time.Time {
						if date := parseDate(prefix + "issuedate"); date != nil {
							return *date
						}
						return time.Time{}
					}(),
					ExpiryDate: func() time.Time {
						if date := parseDate(prefix + "expirydate"); date != nil {
							return *date
						}
						return time.Time{}
					}(),
					IssuePlace: func() string {
						if place := getCol(prefix + "issueplace"); place != nil {
							return *place
						}
						return ""
					}(),
					CountryOfIssue: func() string {
						if country := getCol(prefix + "country"); country != nil {
							return *country
						}
						return "Vietnam" // Default country
					}(),
					HasChip: func() bool {
						if chip := parseBool(prefix + "haschip"); chip != nil {
							return *chip
						}
						return false
					}(),
				}

				// Add notes if available
				notes := getCol(prefix + "notes")
				if notes != nil {
					document.Notes = notes
				}

				req.Documents = append(req.Documents, document)
			}
		}
	}

	// Validate required fields
	if req.StudentCode == nil || req.Fullname == nil || req.Email == nil {
		return nil, fmt.Errorf("missing required fields: studentCode, fullname, or email")
	}

	// Set default status if not provided
	if req.StatusID == nil {
		defaultStatus := 1 // Assuming 1 is "Active" or the default status
		req.StatusID = &defaultStatus
	}

	return req, nil
}

// Helper function to create student request from JSON data
func (s *studentService) createStudentRequestFromJSON(data map[string]interface{}) (*models.CreateStudentRequest, error) {
	// Debug logging
	jsonBytes, _ := json.Marshal(data)
	log.WithField("data", string(jsonBytes)).Debug("Processing JSON record")

	// Helper function to get string value with case-insensitive key matching and nil protection
	getString := func(key string) *string {
		// Try various key formats with nil protection
		possibleKeys := []string{
			key,                                // Original
			strings.ToUpper(key[:1]) + key[1:], // PascalCase
			strings.ToLower(key),               // lowercase
		}

		for _, k := range possibleKeys {
			if val, ok := data[k]; ok && val != nil {
				if strVal, ok := val.(string); ok && strVal != "" {
					return &strVal
				} else if numVal, ok := val.(float64); ok {
					// Convert numeric to string if needed
					strVal := fmt.Sprintf("%v", numVal)
					return &strVal
				}
			}
		}
		return nil
	}

	// Helper function to get int value with case-insensitive key matching and nil protection
	getInt := func(key string) *int {
		// Try various key formats with nil protection
		possibleKeys := []string{
			key,                                // Original
			strings.ToUpper(key[:1]) + key[1:], // PascalCase
			strings.ToLower(key),               // lowercase
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

		// Return default value (1) for status ID if not provided
		if strings.ToLower(key) == "statusid" {
			defaultStatus := 1
			return &defaultStatus
		}

		return nil
	}

	// Helper function to get date value with case-insensitive key matching and nil protection
	getDate := func(key string) *time.Time {
		// Try various key formats with nil protection
		possibleKeys := []string{
			key,                                // Original
			strings.ToUpper(key[:1]) + key[1:], // PascalCase
			strings.ToLower(key),               // lowercase
		}

		for _, k := range possibleKeys {
			if val, ok := data[k]; ok && val != nil {
				if strVal, ok := val.(string); ok && strVal != "" {
					// Try different date formats
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

		// Return current date as fallback for required date fields
		if strings.ToLower(key) == "dateofbirth" {
			now := time.Now()
			return &now
		}

		return nil
	}

	// Create student request with nil checks and default values
	req := &models.CreateStudentRequest{
		StudentCode: getInt("studentCode"),
		Fullname:    getString("fullname"),
		Email:       getString("email"),
		DateOfBirth: getDate("dateOfBirth"),
		Gender:      getString("gender"),
		FacultyID:   getInt("facultyId"),
		Batch:       getString("batch"),
		Program:     getString("program"),
		Address:     getString("address"),
		Phone:       getString("phone"),
		StatusID:    getInt("statusId"),
		ProgramID:   getInt("programId"),
		Nationality: getString("nationality"),
		Addresses:   make([]*models.AddressRequest, 0),
		Documents:   make([]*models.DocumentRequest, 0),
	}

	// Ensure required fields have default values if not provided
	if req.StudentCode == nil {
		return nil, fmt.Errorf("missing required field: studentCode")
	}

	if req.Fullname == nil {
		return nil, fmt.Errorf("missing required field: fullname")
	}

	if req.Email == nil {
		return nil, fmt.Errorf("missing required field: email")
	}

	// Ensure gender has default value if not provided
	if req.Gender == nil {
		defaultGender := "Other"
		req.Gender = &defaultGender
	}

	// Ensure batch has default value if not provided
	if req.Batch == nil {
		currentYear := fmt.Sprintf("%d", time.Now().Year())
		req.Batch = &currentYear
	}

	// Ensure program has default value if not provided
	if req.Program == nil {
		defaultProgram := "Unknown"
		req.Program = &defaultProgram
	}

	// Ensure address has default value if not provided
	if req.Address == nil {
		defaultAddress := ""
		req.Address = &defaultAddress
	}

	// Ensure phone has default value if not provided
	if req.Phone == nil {
		defaultPhone := ""
		req.Phone = &defaultPhone
	}

	// Ensure programID has default value if not provided
	if req.ProgramID == nil {
		defaultProgramID := 1
		req.ProgramID = &defaultProgramID
	}

	// Ensure facultyID has default value if not provided
	if req.FacultyID == nil {
		defaultFacultyID := 1
		req.FacultyID = &defaultFacultyID
	}

	// Ensure nationality has default value if not provided
	if req.Nationality == nil {
		defaultNationality := "Vietnam"
		req.Nationality = &defaultNationality
	}

	// Process addresses from flat structure in JSON
	// Map fields like PermanentStreet, PermanentWard, etc. to address objects
	addressTypes := map[string]struct{}{
		"Permanent": {},
		"Temporary": {},
		"Mailing":   {},
	}

	addressMap := make(map[string]*models.AddressRequest)

	// Check for address fields in the flat structure with nil protection
	for key, value := range data {
		if value == nil {
			continue // Skip nil values
		}

		// Look for fields like "PermanentStreet", "TemporaryCity", etc.
		for prefix := range addressTypes {
			if strings.HasPrefix(key, prefix) && len(key) > len(prefix) {
				field := key[len(prefix):]

				// Create the address object if it doesn't exist
				if _, exists := addressMap[prefix]; !exists {
					addressMap[prefix] = &models.AddressRequest{
						AddressType: prefix,
						Country:     "Vietnam", // Default
					}
				}

				// Set the appropriate field with nil check
				if value != nil {
					if strings.EqualFold(field, "Street") {
						if strVal, ok := value.(string); ok && strVal != "" {
							addressMap[prefix].Street = strVal
						}
					} else if strings.EqualFold(field, "Ward") {
						if strVal, ok := value.(string); ok && strVal != "" {
							addressMap[prefix].Ward = strVal
						}
					} else if strings.EqualFold(field, "District") {
						if strVal, ok := value.(string); ok && strVal != "" {
							addressMap[prefix].District = strVal
						}
					} else if strings.EqualFold(field, "City") {
						if strVal, ok := value.(string); ok && strVal != "" {
							addressMap[prefix].City = strVal
						}
					} else if strings.EqualFold(field, "Country") {
						if strVal, ok := value.(string); ok && strVal != "" {
							addressMap[prefix].Country = strVal
						}
					}
				}
			}
		}
	}

	// Add valid addresses to the request with nil protection
	for _, addr := range addressMap {
		// Only add if we have at least street and city and they're not nil
		if addr != nil && addr.Street != "" && addr.City != "" {
			req.Addresses = append(req.Addresses, addr)
		}
	}

	// Process documents from flat structure in JSON with nil protection
	// Map fields like CCCDNumber, CCCDIssueDate, etc. to document objects
	docPrefixes := []string{"CCCD", "CMND", "Passport"}
	docMap := make(map[string]*models.DocumentRequest)

	// Check for document fields in the flat structure with nil protection
	for key, value := range data {
		if value == nil {
			continue // Skip nil values
		}

		for _, prefix := range docPrefixes {
			if strings.HasPrefix(key, prefix) && len(key) > len(prefix) {
				field := key[len(prefix):]

				// Create the document object if it doesn't exist
				if _, exists := docMap[prefix]; !exists {
					docMap[prefix] = &models.DocumentRequest{
						DocumentType:   prefix,
						CountryOfIssue: "Vietnam", // Default
					}
				}

				// Set the appropriate field with nil protection
				if strings.EqualFold(field, "Number") {
					if numVal, ok := value.(float64); ok {
						docMap[prefix].DocumentNumber = fmt.Sprintf("%v", int(numVal))
					} else if strVal, ok := value.(string); ok && strVal != "" {
						docMap[prefix].DocumentNumber = strVal
					}
				} else if strings.EqualFold(field, "IssueDate") && value != nil {
					if strVal, ok := value.(string); ok && strVal != "" {
						// Try different date formats
						formats := []string{
							"2006-01-02",
							"01/02/2006",
							"02/01/2006",
							"2006/01/02",
						}

						for _, format := range formats {
							date, err := time.Parse(format, strVal)
							if err == nil {
								docMap[prefix].IssueDate = date
								break
							}
						}
					}
				} else if strings.EqualFold(field, "IssuePlace") && value != nil {
					if strVal, ok := value.(string); ok && strVal != "" {
						docMap[prefix].IssuePlace = strVal
					}
				} else if strings.EqualFold(field, "ExpiryDate") && value != nil {
					if strVal, ok := value.(string); ok && strVal != "" {
						// Try different date formats
						formats := []string{
							"2006-01-02",
							"01/02/2006",
							"02/01/2006",
							"2006/01/02",
						}

						for _, format := range formats {
							date, err := time.Parse(format, strVal)
							if err == nil {
								docMap[prefix].ExpiryDate = date
								break
							}
						}
					}
				} else if strings.EqualFold(field, "Country") && value != nil {
					if strVal, ok := value.(string); ok && strVal != "" {
						docMap[prefix].CountryOfIssue = strVal
					}
				} else if strings.EqualFold(field, "HasChip") && value != nil {
					if boolVal, ok := value.(bool); ok {
						docMap[prefix].HasChip = boolVal
					}
				} else if strings.EqualFold(field, "Notes") {
					if strVal, ok := value.(string); ok && strVal != "" {
						docMap[prefix].Notes = &strVal
					} else if value == nil {
						// Handle null notes explicitly
						emptyNote := ""
						docMap[prefix].Notes = &emptyNote
					}
				}
			}
		}
	}

	// Add valid documents to the request with nil protection
	for _, doc := range docMap {
		// Only add if we have at least document number and it's not nil
		if doc != nil && doc.DocumentNumber != "" {
			req.Documents = append(req.Documents, doc)
		}
	}

	return req, nil
}

// customMultipartFile implements multipart.File using a bytes.Reader
type customMultipartFile struct {
	reader *bytes.Reader
	size   int64
}

func (f *customMultipartFile) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *customMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader.ReadAt(p, off)
}

func (f *customMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *customMultipartFile) Close() error {
	// Nothing to close for a bytes.Reader
	return nil
}

// customFileHeader extends multipart.FileHeader with a custom Open method
type customFileHeader struct {
	*multipart.FileHeader
	openFunc func() (multipart.File, error)
}

// Open implements the Open method for the custom file header
func (cfh *customFileHeader) Open() (multipart.File, error) {
	return cfh.openFunc()
}

func (s *studentService) ExportStudentsToCSV(ctx context.Context) (string, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"function": "ExportStudentsToCSV",
	})
	logger.Info("Starting student export to CSV")

	// Get all students from the database with associated faculty names
	students, err := s.studentRepo.List(ctx, models2.QueryParams{
		Limit: -1, // Get all students
	}, func(tx *gorm.DB) {
		tx.Joins(`LEFT JOIN "PUBLIC"."faculties" ON students.faculty_id = "PUBLIC"."faculties".id`)
		tx.Select(`students.*, "PUBLIC"."faculties".name as faculty_name`)
		tx.Preload("Addresses") // Preload addresses
	})

	if err != nil {
		logger.WithError(err).Error("Failed to retrieve students from database")
		return "", fmt.Errorf("failed to retrieve students: %w", err)
	}

	// Create a buffer to hold the CSV data in memory
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
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Loop through students and write each as a row
	for _, student := range students {
		permanentAddress := ""
		temporaryAddress := ""

		// Get addresses if available
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

		// Get faculty ID and name
		facultyID := 0
		if student.FacultyID != 0 {
			facultyID = student.FacultyID
		}

		// Get status name
		statusName := "Active" // Default
		if student.StatusID != 1 {
			// Optionally, get the actual status name from the database
			status, err := s.studentStatusRepo.GetByID(ctx, fmt.Sprintf("%d", student.StatusID))
			if err == nil && status != nil {
				statusName = status.Name
			}
		}

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
		return "", fmt.Errorf("error flushing CSV data: %w", err)
	}

	// Get the CSV data as bytes
	csvData := csvBuffer.Bytes()
	logger.WithField("csvSize", len(csvData)).Info("Generated CSV data in memory")

	// Create a custom implementation of multipart.File using a bytes.Reader
	csvFile := &customMultipartFile{
		reader: bytes.NewReader(csvData),
		size:   int64(len(csvData)),
	}

	// Create a multipart file header
	fileHeader := &multipart.FileHeader{
		Filename: "students-export.csv",
		Size:     int64(len(csvData)),
	}

	// Create the custom file header with our overridden Open method
	myHeader := &customFileHeader{
		FileHeader: fileHeader,
		openFunc: func() (multipart.File, error) {
			return csvFile, nil
		},
	}

	// Upload the file to Google Drive using the gdrive service
	driveFileInfo, err := s.driveService.UploadFile(ctx, myHeader.FileHeader, "students-export.csv")
	if err != nil {
		logger.WithError(err).Error("Failed to upload CSV to Google Drive")
		return "", fmt.Errorf("failed to upload to Google Drive: %w", err)
	}

	// Get the download URL from the drive file info
	downloadURL := driveFileInfo.DownloadURL
	logger.WithField("downloadURL", downloadURL).Info("Successfully exported students to CSV")

	return downloadURL, nil
}
