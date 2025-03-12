package student

import (
	"Group03-EX-StudentManagementAppBE/common"
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
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	// Get and validate student ID from request
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Student ID is required",
		})
		return
	}

	// Get student details from service
	student, err := h.Service.Student.GetByID(c, idStr)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	c.JSON(common.SUCCESS_STATUS, student)
}
