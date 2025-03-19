package program

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	models "Group03-EX-StudentManagementAppBE/internal/models/program"
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
}

// bổ sung
func (h *Handler) ListPrograms(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListProgramRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
	}

	programs, err := h.Service.Program.ListPrograms(c, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, programs)
}
