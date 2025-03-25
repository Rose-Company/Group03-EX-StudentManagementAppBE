package program

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	base.Handler
}

func NewHandler(service *services.Service) Handler {
	return Handler{
		Handler: base.NewHandler(service),
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	programGroup := router.Group("/v1/students")
	programGroup.GET("/programs", middleware.UserAuthentication, h.ListPrograms)
	programGroup.POST("/programs", middleware.UserAuthentication, h.CreateProgram)
	programGroup.PATCH("/programs/:id", middleware.UserAuthentication, h.UpdateProgram)
	programGroup.DELETE("/programs/:id", middleware.UserAuthentication, h.DeleteProgram)
}
