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

func (h *Handler) ImportStudentsFromFile(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req struct {
		File string `json:"file" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	if req.File == "" {
		common.AbortWithError(c, common.ErrLinkNotFound)
		return
	}

	result, err := h.Service.Student.ImportStudentsFromFile(c.Request.Context(), profile.Id, req.File)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if result.ErrorCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":       1,
			"error_code": "import_partial_failure",
			"message":    fmt.Sprintf("Import completed with %d successful and %d failed records", result.SuccessCount, result.ErrorCount),
			"data": gin.H{
				"successful_count": result.SuccessCount,
				"failed_count":     result.ErrorCount,
				"failed_records":   result.FailedRecords,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Students imported successfully"})
}
