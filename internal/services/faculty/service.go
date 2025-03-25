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
	CreateAFaculty(ctx context.Context, userID string, faculty *models.CreateFacultyRequest) error
	UpdateFaculty(ctx context.Context, userID string, id string, faculty *models.UpdateFacultyRequest) error
	DeleteFaculty(ctx context.Context, userID string, id string) error
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
	if req.Sort == "" {
		req.Sort = "name.asc"
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	queryParams := models2.QueryParams{
		Limit:  req.PageSize,
		Offset: offset,
		QuerySort: models2.QuerySort{
			Origin: req.Sort,
		},
	}

	var clauses []repositories.Clause
	if req.Name != "" {
		clauses = append(clauses, func(tx *gorm.DB) {
			tx.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Name+"%")
		})
	}

	combinedClause := func(tx *gorm.DB) {
		for _, clause := range clauses {
			clause(tx)
		}
	}

	faculties, err := s.facultyRepo.List(ctx, queryParams, combinedClause)
	if err != nil {
		log.WithError(err).Error("Failed to list faculties")
		return nil, err
	}

	var facultyList []models.Faculty
	for _, faculty := range faculties {
		facultyList = append(facultyList, *faculty)
	}

	return &models.ListFacultyResponse{
		Items: facultyList,
	}, nil
}

func (s *facultyService) CreateAFaculty(ctx context.Context, userID string, req *models.CreateFacultyRequest) error {
	// Change from CreateFacultyRequest to Faculty
	faculty := &models.Faculty{
		Name: req.Name,
	}

	_, err := s.facultyRepo.Create(ctx, faculty)
	if err != nil {
		return err
	}
	log.Printf("User ID: %s Faculty create successfully with ID: %d", userID, req.ID)
	return nil
}

func (s *facultyService) UpdateFaculty(ctx context.Context, userID string, id string, req *models.UpdateFacultyRequest) error {

	faculty := &models.Faculty{
		Name: req.Name,
	}
	_, err := s.facultyRepo.Update(ctx, id, faculty)
	if err != nil {
		return err
	}
	log.Printf("User ID: %s Faculty updated successfully with ID: %s", userID, id)
	return nil
}

func (s *facultyService) DeleteFaculty(ctx context.Context, userID string, id string) error {
	err := s.facultyRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	log.Printf("User ID: %s deleted faculty with ID: %s", userID, id)
	return nil
}
