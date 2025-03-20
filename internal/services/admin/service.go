package admin

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models"
	adminModels "Group03-EX-StudentManagementAppBE/internal/models/admin"
	"Group03-EX-StudentManagementAppBE/internal/repositories/admin"
	"Group03-EX-StudentManagementAppBE/internal/services/gdrive"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service defines the interface for admin operations
type Service interface {
	ImportStudentFile(ctx context.Context, userID string, file *multipart.FileHeader) (*adminModels.ImportedFileResponse, error)
	GetImportedFile(ctx context.Context, fileID string) (*adminModels.ImportedFileResponse, error)
	ListImportedFiles(ctx context.Context, req *adminModels.ImportedFileListRequest) (*models.BaseListResponse, error)
}

type adminService struct {
	adminRepo    admin.Repository
	driveService gdrive.Service
}

// NewAdminService creates a new admin service
func NewAdminService(adminRepo admin.Repository, driveService gdrive.Service) Service {
	return &adminService{
		adminRepo:    adminRepo,
		driveService: driveService,
	}
}

// ImportStudentFile handles importing a student file
func (s *adminService) ImportStudentFile(ctx context.Context, userID string, file *multipart.FileHeader) (*adminModels.ImportedFileResponse, error) {
	logger := log.WithFields(log.Fields{
		"function": "Import Student File",
		"fileName": file.Filename,
		"fileSize": file.Size,
		"userId":   userID,
	})

	logger.Info("Importing student file")

	// Validate file format
	if !s.driveService.ValidateFileFormat(file.Filename) {
		logger.Error("Invalid file format")
		return nil, common.ErrInvalidFormat
	}

	// Generate unique filename
	fileExt := filepath.Ext(file.Filename)
	uniqueFileName := fmt.Sprintf("student_import_%s%s", uuid.New().String(), fileExt)

	// Upload file to Google Drive
	driveFile, err := s.driveService.UploadFile(ctx, file, uniqueFileName)
	if err != nil {
		logger.WithError(err).Error("Failed to upload file to Google Drive")
		return nil, err
	}

	// Create record in database
	importedFile := &adminModels.File{
		DriveFileID: driveFile.ID,
		FileName:    file.Filename,
		FileSize:    file.Size,
		FileType:    strings.TrimPrefix(fileExt, "."),
		ViewLink:    driveFile.ViewLink,
		DownloadURL: driveFile.DownloadURL,
		ImportedBy:  userID,
		ImportedAt:  time.Now(),
		Status:      "pending",
	}

	// Save to database
	createdFile, err := s.adminRepo.Create(ctx, importedFile)
	if err != nil {
		logger.WithError(err).Error("Failed to create imported file record")
		return nil, common.ErrorWrapper("failed to create imported file record", err)
	}

	logger.WithFields(log.Fields{
		"fileId":      createdFile.ID,
		"driveFileId": driveFile.ID,
	}).Info("Student file imported successfully")

	return createdFile.ToResponse(), nil
}

// GetImportedFile retrieves an imported file by ID
func (s *adminService) GetImportedFile(ctx context.Context, fileID string) (*adminModels.ImportedFileResponse, error) {
	logger := log.WithFields(log.Fields{
		"function": "GetImportedFile",
		"fileId":   fileID,
	})

	logger.Info("Getting imported file")

	// Parse ID
	id, err := uuid.Parse(fileID)
	if err != nil {
		return nil, common.ErrInvalidFormat
	}

	// Get file from database
	file, err := s.adminRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("Imported file not found")
			return nil, err
		}
		logger.WithError(err).Error("Failed to get imported file")
		return nil, err
	}

	return file.ToResponse(), nil
}

// ListImportedFiles retrieves a list of imported files
func (s *adminService) ListImportedFiles(ctx context.Context, req *adminModels.ImportedFileListRequest) (*models.BaseListResponse, error) {
	logger := log.WithFields(log.Fields{
		"function": "ListImportedFiles",
		"page":     req.Page,
		"pageSize": req.PageSize,
	})

	logger.Info("Listing imported files")

	// Prepare query parameters
	queryParams := models.QueryParams{
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
		QuerySort: models.QuerySort{
			Origin: req.Sort,
		},
	}

	// Get files from database

	files, err := s.adminRepo.List(ctx, queryParams)
	if err != nil {
		logger.WithError(err).Error("Failed to list imported files")
		return nil, err
	}

	// Get total count
	totalCount, err := s.adminRepo.Count(ctx, models.QueryParams{})
	if err != nil {
		logger.WithError(err).Error("Failed to get total count")
		return nil, err
	}

	// Convert to response format
	var responseItems []interface{}
	for _, file := range files {
		responseItems = append(responseItems, file.ToResponse())
	}

	return &models.BaseListResponse{
		Total:    int(totalCount),
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    responseItems,
	}, nil
}
