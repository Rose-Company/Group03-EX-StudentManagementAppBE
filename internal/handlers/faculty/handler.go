package faculty

import (
	"Group03-EX-StudentManagementAppBE/common"
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
	}
}

func (h *Handler) GetList(c *gin.Context) {
	result, err := h.Service.Faculty.GetList(c.Request.Context())
	if err != nil {
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Failed to get faculty list",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "ok!",
		Data:    result,
	})
}
