package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/models"
	student_models "Group03-EX-StudentManagementAppBE/internal/models/student"
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

	var createReq student_models.CreateStudentRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Student.CreateAStudent(c.Request.Context(), profile.Id, &createReq)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	studentId := c.Param("id")
	var updateReq student_models.UpdateStudentRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}
	err := h.Service.Student.UpdateStudent(c.Request.Context(), profile.Id, studentId, &updateReq)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
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

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
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

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) ExportStudentsToFile(c *gin.Context) {
	// Parse and validate query parameters
	var req models.FileTypeQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	if !req.IsValidFileType() {
		common.AbortWithError(c, common.ErrInvalidFileFormat)
		return
	}

	fileType := req.GetFileType()

	var data []byte
	var exportErr error

	if fileType == "json" {
		data, exportErr = h.Service.Student.ExportStudentsToJSON(c.Request.Context())
		if exportErr != nil {
			common.AbortWithError(c, exportErr)
			return
		}

		filename := common.FILENAME_JSON
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Type", "application/json")
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
		c.Header("Expires", "0")
		c.Header("Cache-Control", "must-revalidate")
		c.Header("Pragma", "public")
		c.Data(http.StatusOK, "application/json", data)
	} else {
		data, exportErr = h.Service.Student.ExportStudentsToCSV(c.Request.Context())
		if exportErr != nil {
			common.AbortWithError(c, exportErr)
			return
		}

		filename := common.FILENAME_CSV
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Expires", "0")
		c.Header("Cache-Control", "must-revalidate")
		c.Header("Pragma", "public")
		c.Data(http.StatusOK, "text/csv", data)
	}
}
