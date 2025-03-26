package gdrive

import (
	"time"
)

type DriveFileInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	MimeType    string    `json:"mimeType"`
	ViewLink    string    `json:"viewLink"`
	DownloadURL string    `json:"downloadUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	Size        int64     `json:"size"`
}

type UploadRequest struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
}

type UploadResponse struct {
	FileInfo *DriveFileInfo `json:"fileInfo"`
	Success  bool           `json:"success"`
	Message  string         `json:"message,omitempty"`
}
