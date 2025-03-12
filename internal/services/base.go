package services

import (
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services/auth"
)

type Service struct {
	Auth auth.Service
}

func NewService(userRepo user.Repository) *Service {
	return &Service{
		Auth: auth.NewService(userRepo),
	}
}
