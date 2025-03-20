package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudent(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	 err := h.Service.Student.CreateAStudent(c.Request.Context(), &student)
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
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updatedStudent, err := h.Service.Student.UpdateStudent(c.Request.Context(), id, &student)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student updated successfully",
		Data:    updatedStudent,
	})
}

func (h *Handler) DeleteStudentByID(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Student.DeleteByID(c.Request.Context(), id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student deleted successfully",
	})
}
