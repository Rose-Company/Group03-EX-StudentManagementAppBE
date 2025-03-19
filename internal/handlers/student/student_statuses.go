package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) CreateStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var studentStatus models.StudentStatus
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	createdStudentStatus, err := h.Service.Student.CreateStudentStatus(c.Request.Context(), &studentStatus)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status created successfully",
		Data:    createdStudentStatus,
	})
}

func (h *Handler) UpdateStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}


	id := c.Param("id")
	var studentStatus models.StudentStatus
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updatedStudentStatus, err := h.Service.Student.UpdateStudentStatus(c.Request.Context(),id, &studentStatus)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status updated successfully",
		Data:    updatedStudentStatus,
	})
}

func (h *Handler) DeleteStudentStatus(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Student.DeleteStudentStatus(c.Request.Context(), id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status deleted successfully",
	})
}