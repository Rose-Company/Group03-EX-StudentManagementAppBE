package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateAStudent(c *gin.Context) {
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

	createdStudent, err := h.Service.Student.CreateAStudent(c.Request.Context(), &student)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Code:    common.REQUEST_SUCCESS,
		Message: "Student created successfully",
		Data:    createdStudent,
	})
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	fmt.Println(id)
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
		Code:    common.REQUEST_SUCCESS,
		Message: "Student updated successfully",
		Data:    updatedStudent,
	})
}

func (h *Handler) DeleteByID(c *gin.Context) {
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
		Code:    common.REQUEST_SUCCESS,
		Message: "Student deleted successfully",
	})
}