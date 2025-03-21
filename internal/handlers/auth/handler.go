// internal/handlers/auth_handler.go
package auth

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"
	"Group03-EX-StudentManagementAppBE/internal/services"

	"github.com/gin-gonic/gin"
)

// Handler handles authentication routes
type Handler struct {
	base.Handler
}

// NewHandler creates a new auth handler
func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Handler: base.NewHandler(service),
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	authRoutes := rg.Group("/users")
	{
		authRoutes.POST("/login", h.LogIn)
	}
}

// LogIn handles user authentication
func (h *Handler) LogIn(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	loginResp, err := h.Service.Auth.LoginUser(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}


	c.JSON(common.SUCCESS_STATUS, gin.H{
		"code":  loginResp.Code,
		"id":    loginResp.ID,
		"token": loginResp.Token,
	})
}
