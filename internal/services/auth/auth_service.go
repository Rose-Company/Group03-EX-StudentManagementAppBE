// internal/services/auth/service.go
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

var JWTSecret = []byte("your_secret_key")

type authService struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) Service {
	return &authService{
		userRepo: userRepo,
	}
}

// LoginUser handles user login
func (s *authService) LoginUser(ctx context.Context, req models.LoginRequest) (*string, error) {
	var user *models.User
	var err error
	user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", *req.Email)
	})

	//  SELECT * FROM "public"."users" WHERE email = 'mnn27@gmail.com' ORDER BY "users"."id" LIMIT 1
	if err != nil {
		return nil, err
	}

	if user.IsBanned {
		return nil, common.ErrUserBanned
	}

	if user.Provider == common.USER_PROVIDER_GOOGLE {
		return nil, common.ErrGoogleAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password)); err != nil {
		return nil, common.ErrInvalidEmailOrPassWord
	}

	return s.generateJWTToken(user)
}

// generateJWTToken generates a JWT token for authentication
func (s *authService) generateJWTToken(user *models.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.RoleID,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
