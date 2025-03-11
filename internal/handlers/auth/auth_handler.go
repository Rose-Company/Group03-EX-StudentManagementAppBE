// internal/handlers/auth/handler.go
package auth

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers"
	"Group03-EX-StudentManagementAppBE/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	handlers.BaseHandler
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		BaseHandler: handlers.NewBaseHandler(service),
	}
}

// RegisterRoutes registers all auth routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	userRoutes := rg.Group("/users")
	{
		userRoutes.POST("/login", h.LogIn)
		// userRoutes.POST("/register", func(c *gin.Context) { middleware.UserAuthentication(c)
		// 	h.LogIn(c)
		// })
	}

}
