package student

import (
	"Group03-EX-StudentManagementAppBE/common"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentByID(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Mã sinh viên không được để trống",
		})
		return
	}

	student, err := h.Service.Student.GetStudentByID(c, idStr)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(student))
}
