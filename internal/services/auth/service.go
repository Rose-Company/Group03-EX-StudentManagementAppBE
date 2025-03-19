package auth

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service defines the auth service interface
type Service interface {
	LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

// authService is the implementation of the auth service
type authService struct {
	userRepo  user.Repository
	jwtSecret []byte
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo user.Repository, jwtSecret string) Service {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// LoginUser handles user authentication
func (s *authService) LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	var err error
	user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", *req.Email)
	})

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password)); err != nil {
		return nil, common.ErrInvalidEmailOrPassWord
	}

	tokenString, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Code:  200,
		ID:    user.ID,
		Token: *tokenString,
	}, nil
}

// generateJWTToken generates a JWT token for authentication
func (s *authService) generateJWTToken(user *models.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.RoleID,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
