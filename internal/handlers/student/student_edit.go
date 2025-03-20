package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudent(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var createReq models.CreateStudentRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Student.CreateAStudent(c.Request.Context(), profile.Id, &createReq)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Code:    200,
		Message: "Student created successfully",
	})
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	studentId := c.Param("id")
	fmt.Println(studentId)

	var updateReq models.UpdateStudentRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}
	err := h.Service.Student.UpdateStudent(c.Request.Context(), profile.Id, studentId, &updateReq)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student updated successfully",
	})
}

func (h *Handler) DeleteStudentByID(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	studentID := c.Param("id")
	err := h.Service.Student.DeleteStudentByID(c.Request.Context(), profile.Id, studentID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student deleted successfully",
	})
}
