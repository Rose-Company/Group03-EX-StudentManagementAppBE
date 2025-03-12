package student

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	studentGroup := rg.Group("/v1/students")
	{
		studentGroup.GET("/:id", middleware.UserAuthentication, h.GetByID)
		studentGroup.GET("", middleware.UserAuthentication, h.GetList)
	}
}
