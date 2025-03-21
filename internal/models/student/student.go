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

// StudentAddress updated to match the database schema
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

// StudentDocument updated to match the database schema
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

// Response DTOs
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

// Simplified address response for the StudentResponse
type AddressResponse struct {
	Street   string `json:"street"`
	Ward     string `json:"ward"`
	District string `json:"district"`
	City     string `json:"city"`
	Country  string `json:"country"`
}

// ID document response for the StudentResponse
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

// For backward compatibility
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

// For backward compatibility
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

// Used for list views
type StudentListResponse struct {
	Fullname    string `json:"fullname"`
	StudentCode int    `json:"student_code"`
	Email       string `json:"email"`
	FacultyID   int    `json:"faculty_id"`
	FacultyName string `json:"faculty_name"`
	Gender      string `json:"gender"`
}

type ListStudentRequest struct {
	models.BaseRequestParamsUri
	StudentCode string `form:"student_code"`
	Fullname    string `form:"fullname"`
	FacultyID   int    `form:"faculty_id"`
	FacultyName string `form:"faculty_name"`
}

// ToResponse converts Student model to the detailed response
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

	// Process addresses by type if present
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

	// Process documents if present
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

// ToListResponse for simplified list view
func (s *Student) ToListResponse() *StudentListResponse {
	return &StudentListResponse{
		Fullname:    s.Fullname,
		StudentCode: s.StudentCode,
		Email:       s.Email,
		FacultyID:   s.FacultyID,
		Gender:      s.Gender,
	}
}

// For backward compatibility
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

// For backward compatibility
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

// StudentRequest represents the shared fields between create and update operations
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

// AddressRequest represents address fields for API operations
type AddressRequest struct {
	ID          uuid.UUID `json:"id,omitempty"`
	AddressType string    `json:"address_type" binding:"required,oneof=Permanent Temporary Mailing"`
	Street      string    `json:"street" binding:"required"`
	Ward        string    `json:"ward"`
	District    string    `json:"district"`
	City        string    `json:"city" binding:"required"`
	Country     string    `json:"country" binding:"required"`
}

// DocumentRequest represents document fields for API operations
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

// CreateStudentRequest for student creation API - same as StudentRequest
type CreateStudentRequest StudentRequest

// UpdateStudentRequest for student update API - same as StudentRequest
// But when handling in the service, only update fields that are provided
type UpdateStudentRequest StudentRequest
