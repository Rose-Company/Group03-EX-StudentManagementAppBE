package admin

import (
	"Group03-EX-StudentManagementAppBE/common"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DriveFileID string    `json:"drive_file_id" gorm:"type:text;not null"`
	FileName    string    `json:"file_name" gorm:"type:text;not null"`
	FileSize    int64     `json:"file_size" gorm:"type:bigint;not null"`
	FileType    string    `json:"file_type" gorm:"type:text;not null"`
	ViewLink    string    `json:"view_link" gorm:"type:text"`
	DownloadURL string    `json:"download_url" gorm:"type:text"`
	ImportedBy  string    `json:"imported_by" gorm:"type:uuid;not null"`
	ImportedAt  time.Time `json:"imported_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	Status      string    `json:"status" gorm:"type:text;default:'pending'"`
	ProcessedAt time.Time `json:"processed_at" gorm:"type:timestamptz"`
	Notes       string    `json:"notes" gorm:"type:text"`
}

func (f *File) TableName() string {
	return common.POSTGRES_TABLE_NAME_FILES
}

func (f *File) ToResponse() *ImportedFileResponse {
	return &ImportedFileResponse{
		ID:          f.ID.String(),
		FileName:    f.FileName,
		FileSize:    f.FileSize,
		FileType:    f.FileType,
		ViewLink:    f.ViewLink,
		DownloadURL: f.DownloadURL,
		ImportedBy:  f.ImportedBy,
		ImportedAt:  f.ImportedAt,
		Status:      f.Status,
	}
}

type ImportedFileResponse struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	FileType    string    `json:"file_type"`
	ViewLink    string    `json:"view_link"`
	DownloadURL string    `json:"download_url"`
	ImportedBy  string    `json:"imported_by"`
	ImportedAt  time.Time `json:"imported_at"`
	Status      string    `json:"status"`
}

type ImportedFileListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Sort     string `form:"sort"`
}

type ImportedFileCreateRequest struct {
	File []byte `json:"-"`
}

type ImportResult struct {
	SuccessCount  int                  `json:"success_count"`
	ErrorCount    int                  `json:"error_count"`
	FailedRecords []FailedRecordDetail `json:"failed_records,omitempty"`
}

type FailedRecordDetail struct {
	RowNumber   int    `json:"row_number,omitempty"`
	StudentCode string `json:"student_code,omitempty"`
	Email       string `json:"email,omitempty"`
	Error       string `json:"error"`
}

type FileHeaderWrapper struct {
	Header *multipart.FileHeader
	File   multipart.File
}

func (f *FileHeaderWrapper) Open() (multipart.File, error) {
	return f.File, nil
}

func (f *FileHeaderWrapper) Filename() string {
	return f.Header.Filename
}

func (f *FileHeaderWrapper) Size() int64 {
	return f.Header.Size
}
