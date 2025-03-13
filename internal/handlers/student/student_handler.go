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
		studentGroup.GET("/test", middleware.UserAuthentication, h.GetList)
		studentGroup.DELETE("/:id", middleware.UserAuthentication, h.DeleteByID)
		studentGroup.POST("/create", middleware.UserAuthentication, h.CreateAStudent)
		studentGroup.PUT("/:id", middleware.UserAuthentication, h.UpdateStudent)
		studentGroup.GET("/:id", middleware.UserAuthentication, h.GetByID)

	}
}
