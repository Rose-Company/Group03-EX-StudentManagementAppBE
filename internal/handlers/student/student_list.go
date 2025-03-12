package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetList(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListStudentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
	}

	// Get student list from service
	students, err := h.Service.Student.GetList(c, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	c.JSON(common.SUCCESS_STATUS, students)
}
