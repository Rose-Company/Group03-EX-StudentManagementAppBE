package faculty

import (
	"Group03-EX-StudentManagementAppBE/common"
	models2 "Group03-EX-StudentManagementAppBE/internal/models"
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"context"
)

type Service interface {
	GetList(ctx context.Context, req *models.ListFacultyRequest) (*models.ListFacultyResponse, error)
	CreateAFaculty(ctx context.Context, faculty *models.CreateFacultyRequest) (error)
	UpdateFaculty(ctx context.Context,id string, faculty *models.UpdateFacultyRequest) (*models.Faculty, error)
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
    if req.Sort == "" {
		req.Sort = "NAME.ASC"
	}
	sort := common.ParseSortString(req.Sort)
	if req.PageSize < 0	{
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	faculties, err := s.facultyRepo.List(ctx, models2.QueryParams{
		Limit:  req.PageSize,
		Offset: offset,
		QuerySort: models2.QuerySort{
			Sort: sort,
		}},)	
	if err != nil {
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


func (s *facultyService) CreateAFaculty(ctx context.Context, req *models.CreateFacultyRequest) ( error) {
    // Chuyển đổi từ CreateFacultyRequest sang Faculty
    faculty := &models.Faculty{
        Name: req.Name,
    }

     err := s.facultyRepo.Create(ctx, faculty)
    if err != nil {
        return err
    }

    return  nil
}

func (s *facultyService) UpdateFaculty(ctx context.Context,id string, req *models.UpdateFacultyRequest) (*models.Faculty, error) {

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

