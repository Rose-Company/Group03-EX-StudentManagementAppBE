package gdrive

import (
	"time"
)

// DriveFileInfo contains information about a file stored in Google Drive
type DriveFileInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	MimeType    string    `json:"mimeType"`
	ViewLink    string    `json:"viewLink"`
	DownloadURL string    `json:"downloadUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	Size        int64     `json:"size"`
}

// UploadRequest represents a request to upload a file to Google Drive
type UploadRequest struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
}

// UploadResponse represents a response from uploading a file to Google Drive
type UploadResponse struct {
	FileInfo *DriveFileInfo `json:"fileInfo"`
	Success  bool           `json:"success"`
	Message  string         `json:"message,omitempty"`
}
