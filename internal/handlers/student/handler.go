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
		studentGroup.DELETE("/:id", middleware.UserAuthentication, h.DeleteStudentByID)
		studentGroup.POST("/create", middleware.UserAuthentication, h.CreateStudent)
		studentGroup.PATCH("/:id", middleware.UserAuthentication, h.UpdateStudent)
		studentGroup.GET("/:id", middleware.UserAuthentication, h.GetStudentByID)
		studentGroup.GET("/statuses", middleware.UserAuthentication, h.GetStudentStatuses)
		studentGroup.POST("/statuses", middleware.UserAuthentication, h.CreateStudentStatus)
		studentGroup.PATCH("/statuses/:id", middleware.UserAuthentication, h.UpdateStudentStatus)
		studentGroup.DELETE("/statuses/:id", middleware.UserAuthentication, h.DeleteStudentStatus)
		studentGroup.POST("/import-from-file", middleware.UserAuthentication, h.ImportStudentsFromFile)
		studentGroup.GET("/exported-file", middleware.UserAuthentication, h.ExportStudentsToFile)
		studentGroup.GET("", middleware.UserAuthentication, h.GetStudentList)

	}

}
