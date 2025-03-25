package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentStatuses(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListStudentStatusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	studentStatuses, err := h.Service.Student.GetStudentStatuses(c, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(studentStatuses))
}

func (h *Handler) CreateStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var studentStatus models.CreateStudentStatusRequest
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Student.CreateStudentStatus(c.Request.Context(), &studentStatus)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) UpdateStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var studentStatus models.UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updatedStudentStatus, err := h.Service.Student.UpdateStudentStatus(c.Request.Context(), id, &studentStatus)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(updatedStudentStatus))
}

func (h *Handler) DeleteStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Mã trạng thái sinh viên không được để trống",
		})
		return
	}

	err := h.Service.Student.DeleteStudentStatus(c.Request.Context(), id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}
