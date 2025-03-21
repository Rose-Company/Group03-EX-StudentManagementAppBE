package faculty

import (
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"context"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Service interface {
	GetList(ctx context.Context, req *models.ListFacultyRequest) (*models.ListFacultyResponse, error)
	CreateAFaculty(ctx context.Context, faculty *models.CreateFacultyRequest) error
	UpdateFaculty(ctx context.Context, id string, faculty *models.UpdateFacultyRequest) (*models.Faculty, error)
	DeleteFaculty(ctx context.Context, id string) error
}

type facultyService struct {
	facultyRepo faculty.Repository
}

func NewFalcutyService(facultyRepo faculty.Repository) Service {
	return &facultyService{
		facultyRepo: facultyRepo,
	}
}

func (s *facultyService) GetList(ctx context.Context, req *models.ListFacultyRequest) (*models.ListFacultyResponse, error) {
	// Set default values if not provided
	if req.Sort == "" {
		req.Sort = "name.asc" // Default sort by name ascending
	}

	if req.PageSize <= 0 {
		req.PageSize = 10 // Default page size
	}

	// Calculate offset based on page and page size
	offset := (req.Page - 1) * req.PageSize

	// Create query parameters
	queryParams := models2.QueryParams{
		Limit:  req.PageSize,
		Offset: offset,
		QuerySort: models2.QuerySort{
			Origin: req.Sort, // Use Origin instead of Sort for proper handling
		},
	}

	// Create filter clauses if needed
	var clauses []repositories.Clause
	if req.Name != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Name+"%")
		})
	}

	// Combine clauses if any
	combinedClause := func(tx *gorm.DB) {
		for _, clause := range clauses {
			clause(tx)
		}
	}

	// Get faculties with filters
	faculties, err := s.facultyRepo.List(ctx, queryParams, combinedClause)
	if err != nil {
		log.WithError(err).Error("Failed to list faculties")
		return nil, err
	}

	// Convert to response model
	var facultyList []models.Faculty
	for _, faculty := range faculties {
		facultyList = append(facultyList, *faculty)
	}

	// Create and return response
	return &models.ListFacultyResponse{
		Items: facultyList,
	}, nil
}

func (s *facultyService) CreateAFaculty(ctx context.Context, req *models.CreateFacultyRequest) error {
	// Chuyển đổi từ CreateFacultyRequest sang Faculty
	faculty := &models.Faculty{
		Name: req.Name,
	}

	_, err := s.facultyRepo.Create(ctx, faculty)
	if err != nil {
		return err
	}

	return nil
}

func (s *facultyService) UpdateFaculty(ctx context.Context, id string, req *models.UpdateFacultyRequest) (*models.Faculty, error) {

	faculty := &models.Faculty{
		Name: req.Name,
	}
	updatedFaculty, err := s.facultyRepo.Update(ctx, id, faculty)
	if err != nil {
		return nil, err
	}

	return updatedFaculty, nil
}

func (s *facultyService) DeleteFaculty(ctx context.Context, id string) error {
	return s.facultyRepo.DeleteByID(ctx, id)
}
