package student

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	base.Handler
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Handler: base.NewHandler(service),
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	studentGroup := rg.Group("/v1/students")
	{
		studentGroup.GET("", middleware.UserAuthentication, h.GetStudentList)
		studentGroup.DELETE("/:id", middleware.UserAuthentication, h.DeleteStudentByID)
		studentGroup.POST("/create", middleware.UserAuthentication, h.CreateStudent)
		studentGroup.PUT("/:id", middleware.UserAuthentication, h.UpdateStudent)
		studentGroup.GET("/:id", middleware.UserAuthentication, h.GetStudentByID)
		// /student/statuses
		studentGroup.GET("/statuses", middleware.UserAuthentication, h.GetStudentStatuses)
		studentGroup.POST("/statuses", middleware.UserAuthentication, h.CreateStudentStatus)
		studentGroup.PUT("/statuses/:id", middleware.UserAuthentication, h.UpdateStudentStatus)
		studentGroup.DELETE("/statuses/:id", middleware.UserAuthentication, h.DeleteStudentStatus)
	}

	
}
