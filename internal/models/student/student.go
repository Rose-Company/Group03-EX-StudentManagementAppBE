package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID          uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StudentCode int                `json:"student_code" gorm:"unique;not null"`
	Fullname    string             `json:"fullname" gorm:"type:text;not null"`
	DateOfBirth time.Time          `json:"date_of_birth" gorm:"type:date;not null"`
	Gender      string             `json:"gender" gorm:"type:text;check:gender IN ('Male', 'Female', 'Other')"`
	FacultyID   int                `json:"faculty_id" gorm:"type:integer;references:faculties(id)"`
	Batch       string             `json:"batch" gorm:"type:text;not null"`
	Program     string             `json:"program" gorm:"type:text;not null"`
	Address     string             `json:"address" gorm:"type:text"`
	Email       string             `json:"email" gorm:"type:text;unique"`
	Phone       string             `json:"phone" gorm:"type:text"`
	StatusID    int                `json:"status_id" gorm:"type:integer;references:student_statuses(id)"`
	Addresses   []*StudentAddress  `json:"addresses,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	Documents   []*StudentDocument `json:"documents,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	CreatedAt   time.Time          `json:"created_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time          `json:"updated_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	ProgramID   int                `json:"program_id" gorm:"type:integer;references:programs(id)"`
	Nationality string             `json:"nationality" gorm:"type:text"`
	FacultyName string             `json:"faculty_name,omitempty" gorm:"-"`
}

func (s *Student) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENTS
}

type StudentAddress struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StudentID   uuid.UUID `json:"student_id" gorm:"type:uuid;not null"`
	AddressType string    `json:"address_type" gorm:"type:text"`
	Street      string    `json:"street" gorm:"type:text"`
	Ward        string    `json:"ward" gorm:"type:text"`
	District    string    `json:"district" gorm:"type:text"`
	City        string    `json:"city" gorm:"type:text"`
	Country     string    `json:"country" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

func (a *StudentAddress) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENT_ADDRESSES
}

type StudentDocument struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StudentID      uuid.UUID `json:"student_id" gorm:"type:uuid;not null"`
	DocumentType   string    `json:"document_type" gorm:"type:text;not null"`
	DocumentNumber string    `json:"document_number" gorm:"type:text;not null"`
	IssueDate      time.Time `json:"issue_date" gorm:"type:date"`
	IssuePlace     string    `json:"issue_place" gorm:"type:text"`
	ExpiryDate     time.Time `json:"expiry_date" gorm:"type:date"`
	CountryOfIssue string    `json:"country_of_issue" gorm:"type:text"`
	HasChip        bool      `json:"has_chip" gorm:"type:boolean"`
	Notes          *string   `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

func (d *StudentDocument) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENT_IDENTITY_DOCUMENTS
}

type StudentStatus struct {
	ID   int    `json:"id" gorm:"type:serial;primary_key"`
	Name string `json:"name" gorm:"type:text;not null"`
}

func (s *StudentStatus) TableName() string {
	return common.POSTGRES_TABLE_NAME_STUDENTS_STATUSES
}

type StudentResponse struct {
	ID               uuid.UUID            `json:"id"`
	StudentCode      int                  `json:"student_code"`
	Fullname         string               `json:"fullname"`
	DateOfBirth      time.Time            `json:"date_of_birth"`
	Gender           string               `json:"gender"`
	FacultyID        int                  `json:"faculty_id"`
	FacultyName      string               `json:"faculty_name,omitempty"`
	Batch            string               `json:"batch"`
	Program          string               `json:"program"`
	Nationality      string               `json:"nationality,omitempty"`
	Email            string               `json:"email"`
	Phone            string               `json:"phone"`
	StatusID         int                  `json:"status_id"`
	PermanentAddress *AddressResponse     `json:"permanent_address,omitempty"`
	TempAddress      *AddressResponse     `json:"temp_address,omitempty"`
	MailingAddress   *AddressResponse     `json:"mailing_address,omitempty"`
	IDDocuments      []IDDocumentResponse `json:"id_documents,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	ProgramID        int                  `json:"program_id,omitempty"`
}

type AddressResponse struct {
	Street   string `json:"street"`
	Ward     string `json:"ward"`
	District string `json:"district"`
	City     string `json:"city"`
	Country  string `json:"country"`
}

type IDDocumentResponse struct {
	ID             uuid.UUID `json:"id"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
	IssueDate      time.Time `json:"issue_date"`
	IssuePlace     string    `json:"issue_place"`
	ExpiryDate     time.Time `json:"expiry_date"`
	CountryOfIssue string    `json:"country_of_issue"`
	HasChip        bool      `json:"has_chip"`
	Notes          *string   `json:"notes"`
}

type StudentAddressResponse struct {
	ID          uuid.UUID `json:"id"`
	StudentID   uuid.UUID `json:"student_id"`
	AddressType string    `json:"address_type"`
	Street      string    `json:"street"`
	Ward        string    `json:"ward"`
	District    string    `json:"district"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
}

type StudentDocumentResponse struct {
	ID             uuid.UUID `json:"id"`
	StudentID      uuid.UUID `json:"student_id"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
	IssueDate      time.Time `json:"issue_date"`
	IssuePlace     string    `json:"issue_place"`
	ExpiryDate     time.Time `json:"expiry_date"`
	CountryOfIssue string    `json:"country_of_issue"`
	HasChip        bool      `json:"has_chip"`
	Notes          *string   `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StudentListResponse struct {
	ID          uuid.UUID `json:"id"`
	Fullname    string    `json:"fullname"`
	StudentCode int       `json:"student_code"`
	Email       string    `json:"email"`
	FacultyID   int       `json:"faculty_id"`
	FacultyName string    `json:"faculty_name"`
	Gender      string    `json:"gender"`
}

type ListStudentRequest struct {
	models.BaseRequestParamsUri
	StudentCode string `form:"student_code"`
	Fullname    string `form:"fullname"`
	FacultyID   int    `form:"faculty_id"`
	FacultyName string `form:"faculty_name"`
}

func (s *Student) ToResponse() *StudentResponse {
	response := &StudentResponse{
		ID:          s.ID,
		StudentCode: s.StudentCode,
		Fullname:    s.Fullname,
		DateOfBirth: s.DateOfBirth,
		Gender:      s.Gender,
		FacultyID:   s.FacultyID,
		Batch:       s.Batch,
		Program:     s.Program,
		Email:       s.Email,
		Phone:       s.Phone,
		StatusID:    s.StatusID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		ProgramID:   s.ProgramID,
		Nationality: s.Nationality,
	}

	if len(s.Addresses) > 0 {
		for _, addr := range s.Addresses {
			addrResponse := &AddressResponse{
				Street:   addr.Street,
				Ward:     addr.Ward,
				District: addr.District,
				City:     addr.City,
				Country:  addr.Country,
			}

			switch addr.AddressType {
			case "Permanent":
				response.PermanentAddress = addrResponse
			case "Temporary":
				response.TempAddress = addrResponse
			case "Mailing":
				response.MailingAddress = addrResponse
			}
		}
	}

	if len(s.Documents) > 0 {
		for _, doc := range s.Documents {
			response.IDDocuments = append(response.IDDocuments, IDDocumentResponse{
				ID:             doc.ID,
				DocumentType:   doc.DocumentType,
				DocumentNumber: doc.DocumentNumber,
				IssueDate:      doc.IssueDate,
				IssuePlace:     doc.IssuePlace,
				ExpiryDate:     doc.ExpiryDate,
				CountryOfIssue: doc.CountryOfIssue,
				HasChip:        doc.HasChip,
				Notes:          doc.Notes,
			})
		}
	}

	return response
}

func (s *Student) ToListResponse() *StudentListResponse {
	return &StudentListResponse{
		ID:          s.ID,
		Fullname:    s.Fullname,
		StudentCode: s.StudentCode,
		Email:       s.Email,
		FacultyID:   s.FacultyID,
		Gender:      s.Gender,
	}
}

func (a *StudentAddress) ToResponse() *StudentAddressResponse {
	return &StudentAddressResponse{
		ID:          a.ID,
		StudentID:   a.StudentID,
		AddressType: a.AddressType,
		Street:      a.Street,
		Ward:        a.Ward,
		District:    a.District,
		City:        a.City,
		Country:     a.Country,
	}
}

func (d *StudentDocument) ToResponse() *StudentDocumentResponse {
	return &StudentDocumentResponse{
		ID:             d.ID,
		StudentID:      d.StudentID,
		DocumentType:   d.DocumentType,
		DocumentNumber: d.DocumentNumber,
		IssueDate:      d.IssueDate,
		IssuePlace:     d.IssuePlace,
		ExpiryDate:     d.ExpiryDate,
		CountryOfIssue: d.CountryOfIssue,
		HasChip:        d.HasChip,
		Notes:          d.Notes,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

type StudentRequest struct {
	StudentCode *int               `json:"student_code"`
	Fullname    *string            `json:"fullname"`
	DateOfBirth *time.Time         `json:"date_of_birth"`
	Gender      *string            `json:"gender"`
	FacultyID   *int               `json:"faculty_id"`
	Batch       *string            `json:"batch"`
	Program     *string            `json:"program"`
	Address     *string            `json:"address"`
	Email       *string            `json:"email"`
	Phone       *string            `json:"phone"`
	StatusID    *int               `json:"status_id"`
	ProgramID   *int               `json:"program_id"`
	Nationality *string            `json:"nationality"`
	Addresses   []*AddressRequest  `json:"addresses"`
	Documents   []*DocumentRequest `json:"documents"`
}

type AddressRequest struct {
	ID          uuid.UUID `json:"id,omitempty"`
	AddressType string    `json:"address_type" binding:"required,oneof=Permanent Temporary Mailing"`
	Street      string    `json:"street" binding:"required"`
	Ward        string    `json:"ward"`
	District    string    `json:"district"`
	City        string    `json:"city" binding:"required"`
	Country     string    `json:"country" binding:"required"`
}

type DocumentRequest struct {
	ID             uuid.UUID `json:"id,omitempty"`
	DocumentType   string    `json:"document_type" binding:"required"`
	DocumentNumber string    `json:"document_number" binding:"required"`
	IssueDate      time.Time `json:"issue_date"`
	IssuePlace     string    `json:"issue_place"`
	ExpiryDate     time.Time `json:"expiry_date"`
	CountryOfIssue string    `json:"country_of_issue"`
	HasChip        bool      `json:"has_chip"`
	Notes          *string   `json:"notes"`
}

type CreateStudentRequest StudentRequest

type UpdateStudentRequest StudentRequest

// ImportRecord represents a single record from an import file
type ImportRecord struct {
	Index int
	Data  *CreateStudentRequest
	Err   error
}

// StudentExport represents the structure for JSON export
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

// ImportResult represents the result of an import operation
type ImportResult struct {
	SuccessCount  int
	ErrorCount    int
	FailedRecords []FailedRecordDetail
}

// FailedRecordDetail represents details about a failed import record
type FailedRecordDetail struct {
	RowNumber   int
	StudentCode string
	Email       string
	Error       string
}
