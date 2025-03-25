package gdrive

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models/gdrive"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// HTTP implementation of Drive Service
type httpDriveService struct {
	config *jwt.Config
	client *http.Client
}

type Service interface {
	// UploadFile uploads a file to Google Drive and returns the file info
	UploadFile(ctx context.Context, file *multipart.FileHeader, fileName string) (*gdrive.DriveFileInfo, error)
	// ValidateFileFormat checks if the file has an allowed format
	ValidateFileFormat(fileName string) bool
	// GetMimeType returns the appropriate MIME type based on file extension
	GetMimeType(fileName string) string
}

// DriveAPI endpoints
const (
	uploadEndpoint         = "https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart"
	getFileEndpoint        = "https://www.googleapis.com/drive/v3/files/%s?fields=id,name,mimeType,webViewLink,size"
	createPermissionFormat = "https://www.googleapis.com/drive/v3/files/%s/permissions"
)

// Allowed file extensions
var allowedExtensions = []string{common.CSV_FILE_EXTENSION, common.CSV_FILE_EXTENSION}

// NewHTTPDriveService creates a new Google Drive service using direct HTTP calls
func NewHTTPDriveService(credentialsFile string) (Service, error) {
	logger := log.WithFields(log.Fields{
		"component": "HTTP Drive Service",
		"action":    "Initialize",
	})

	logger.Info("Initializing Google Drive HTTP service")

	// Read credentials file
	credBytes, err := os.ReadFile(credentialsFile)
	if err != nil {
		logger.WithError(err).Error("Unable to read credentials file")
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	// Parse credentials into JWT config
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		logger.WithError(err).Error("Unable to parse credentials")
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	// Create HTTP client with JWT auth
	client := config.Client(context.Background())

	logger.Info("Google Drive HTTP service initialized successfully")

	return &httpDriveService{
		config: config,
		client: client,
	}, nil
}

// UploadFile uploads a file to Google Drive using HTTP API
func (s *httpDriveService) UploadFile(ctx context.Context, file *multipart.FileHeader, fileName string) (*gdrive.DriveFileInfo, error) {
	logger := log.WithFields(log.Fields{
		"component": "HTTP Drive Service",
		"action":    "Upload File",
		"fileName":  fileName,
	})

	logger.Info("Starting file upload to Google Drive via HTTP API")

	// Open the file
	src, err := file.Open()
	if err != nil {
		logger.WithError(err).Error("Failed to open file")
		return nil, err
	}
	defer src.Close()

	// Read file content
	fileContent, err := io.ReadAll(src)
	if err != nil {
		logger.WithError(err).Error("Failed to read file content")
		return nil, err
	}

	// Get MIME type from file extension
	mimeType := s.GetMimeType(fileName)

	// Create metadata part
	metadata := map[string]interface{}{
		"name":     fileName,
		"mimeType": mimeType,
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal metadata")
		return nil, err
	}

	// Create multipart request body
	boundary := "foo_bar_baz"
	body := new(bytes.Buffer)

	// Add metadata part
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/json; charset=UTF-8\r\n\r\n")
	body.Write(metadataBytes)
	body.WriteString("\r\n")

	// Add file content part
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString(fmt.Sprintf("Content-Type: %s\r\n\r\n", mimeType))
	body.Write(fileContent)
	body.WriteString(fmt.Sprintf("\r\n--%s--", boundary))

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", uploadEndpoint, body)
	if err != nil {
		logger.WithError(err).Error("Failed to create request")
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", boundary))

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		logger.WithError(err).Error("Failed to execute upload request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		logger.WithFields(log.Fields{
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Upload request failed")
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Extract file ID
	fileID, ok := result["id"].(string)
	if !ok || fileID == "" {
		logger.Error("Failed to get file ID from response")
		return nil, fmt.Errorf("invalid response: missing file ID")
	}

	// Set file permission (anyone with link can view)
	err = s.setFilePermission(ctx, fileID)
	if err != nil {
		logger.WithError(err).Warn("Failed to set file permissions, but upload succeeded")
		// Continue even if permission update fails
	}

	// Get file details
	fileInfo, err := s.getFileInfo(ctx, fileID)
	if err != nil {
		logger.WithError(err).Warn("Failed to get file details, returning partial information")
		// Create minimal file info if get details fails
		downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
		fileInfo = &gdrive.DriveFileInfo{
			ID:          fileID,
			Name:        fileName,
			MimeType:    mimeType,
			ViewLink:    fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID),
			DownloadURL: downloadURL,
			CreatedAt:   time.Now(),
		}
	}

	logger.WithFields(log.Fields{
		"fileId":   fileID,
		"viewLink": fileInfo.ViewLink,
	}).Info("Successfully uploaded file to Google Drive")

	return fileInfo, nil
}

// setFilePermission sets the file permission to be accessible by anyone with the link
func (s *httpDriveService) setFilePermission(ctx context.Context, fileID string) error {
	url := fmt.Sprintf(createPermissionFormat, fileID)

	// Create permission payload
	permission := map[string]string{
		"role": "reader",
		"type": "anyone",
	}

	payloadBytes, err := json.Marshal(permission)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("set permission failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// getFileInfo gets the file metadata from Google Drive
func (s *httpDriveService) getFileInfo(ctx context.Context, fileID string) (*gdrive.DriveFileInfo, error) {
	url := fmt.Sprintf(getFileEndpoint, fileID)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get file info failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Create download URL
	downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)

	// Extract fields
	name, _ := result["name"].(string)
	mimeType, _ := result["mimeType"].(string)
	webViewLink, _ := result["webViewLink"].(string)

	// Extract size (might be float64 from JSON)
	var size int64
	if sizeFloat, ok := result["size"].(float64); ok {
		size = int64(sizeFloat)
	}

	return &gdrive.DriveFileInfo{
		ID:          fileID,
		Name:        name,
		MimeType:    mimeType,
		ViewLink:    webViewLink,
		DownloadURL: downloadURL,
		CreatedAt:   time.Now(),
		Size:        size,
	}, nil
}

// ValidateFileFormat checks if the file has an allowed format
func (s *httpDriveService) ValidateFileFormat(fileName string) bool {
	// Get file extension
	ext := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			ext = fileName[i:]
			break
		}
	}

	// Convert to lowercase
	ext = strings.ToLower(ext)

	// Check if extension is allowed
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

// GetMimeType returns the appropriate MIME type based on file extension
func (s *httpDriveService) GetMimeType(fileName string) string {
	// Get file extension
	ext := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			ext = fileName[i:]
			break
		}
	}

	// Convert to lowercase
	ext = strings.ToLower(ext)

	// Return appropriate MIME type
	switch ext {
	case common.XLSX_FILE_EXTENSION:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case common.XLS_FILE_EXTENSION:
		return "application/vnd.ms-excel"
	case common.CSV_FILE_EXTENSION:
		return "text/csv"
	case common.JSON_FILE_EXTENSION:
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
