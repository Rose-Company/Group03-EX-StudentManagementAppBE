package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"
	"Group03-EX-StudentManagementAppBE/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Invalid student ID format",
		})
		return
	}

	// Get student details from service
	student, err := h.Service.studentService.GetByID(c, id)
	if err != nil {
		c.JSON(common.NOT_FOUND_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Student not found",
		})
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.Response{
		Code:    common.REQUEST_SUCCESS,
		Message: "ok!",
		Data:    student,
	})
}
