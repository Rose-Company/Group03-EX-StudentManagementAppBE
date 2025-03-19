package faculty

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
	facultyRoutes := rg.Group("/v1/faculties")
	{
		facultyRoutes.GET("", middleware.UserAuthentication, h.GetList)
		facultyRoutes.POST("", middleware.UserAuthentication, h.CreateAFaculty)
		facultyRoutes.PUT("/:id", middleware.UserAuthentication, h.UpdateFaculty)
		facultyRoutes.DELETE("/:id", middleware.UserAuthentication, h.DeleteFaculty)
	}
}
