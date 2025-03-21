package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentStatuses(c *gin.Context) {
	log.Println("Handling request: GetStudentStatuses - Fetching student statuses")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in GetStudentStatuses")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListStudentStatusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding query params in GetStudentStatuses: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	studentStatuses, err := h.Service.Student.GetStudentStatuses(c, &req)
	if err != nil {
		log.Printf("Error fetching student statuses: %v", err)
		common.AbortWithError(c, err)
		return
	}

	log.Println("Successfully fetched student statuses")
	c.JSON(common.SUCCESS_STATUS, studentStatuses)
}

func (h *Handler) CreateStudentStatus(c *gin.Context) {
	log.Println("Handling request: CreateStudentStatus - Creating student status")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in CreateStudentStatus")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var studentStatus models.CreateStudentStatusRequest
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		log.Printf("Error binding request body in CreateStudentStatus: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Student.CreateStudentStatus(c.Request.Context(), &studentStatus)
	if err != nil {
		log.Printf("Error creating student status: %v", err)
		common.AbortWithError(c, err)
		return
	}

	log.Println("Student status created successfully")
	c.JSON(http.StatusCreated, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status created successfully",
	})
}

func (h *Handler) UpdateStudentStatus(c *gin.Context) {
	log.Println("Handling request: UpdateStudentStatus - Updating student status")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in UpdateStudentStatus")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var studentStatus models.UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&studentStatus); err != nil {
		log.Printf("Error binding request body in UpdateStudentStatus: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updatedStudentStatus, err := h.Service.Student.UpdateStudentStatus(c.Request.Context(), id, &studentStatus)
	if err != nil {
		log.Printf("Error updating student status with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Student status updated successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status updated successfully",
		Data:    updatedStudentStatus,
	})
}

func (h *Handler) DeleteStudentStatus(c *gin.Context) {
	log.Println("Handling request: DeleteStudentStatus - Deleting student status")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in DeleteStudentStatus")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	if id == "" {
		log.Println("Error: Missing student status ID in DeleteStudentStatus")
		c.JSON(http.StatusBadRequest, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Student status ID is required",
		})
		return
	}

	log.Printf("Attempting to delete student status with ID: %s", id)

	err := h.Service.Student.DeleteStudentStatus(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error deleting student status with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Student status deleted successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Student status deleted successfully",
	})
}
