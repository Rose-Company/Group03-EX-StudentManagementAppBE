// internal/services/auth/service.go
package auth

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/config"
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JWTSecret = []byte("your_secret_key")

type authService struct {
	userRepo  user.Repository
	jwtSecret []byte
}

func NewService(userRepo user.Repository) Service {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(cfg.JWTSecret),
	}
}

// LoginUser handles user login
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

// refreshJWTToken refreshes a JWT token
// func (s *authService) refreshJWTToken(tokenString string) (*string, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return s.jwtSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok || !token.Valid {
// 		return nil, common.ErrInvalidToken
// 	}

// 	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":    claims["id"],
// 		"email": claims["email"],
// 		"role":  claims["role"],
// 		"exp":   time.Now().Add(24 * time.Hour).Unix(),
// 	})

// 	newTokenString, err := newToken.SignedString(s.jwtSecret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &newTokenString, nil
// }