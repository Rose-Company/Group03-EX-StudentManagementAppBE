package student

import (
	"Group03-EX-StudentManagementAppBE/common"

	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentByID(c *gin.Context) {
	log.Println("Handling request: GetStudentByID - Fetching student details")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in GetStudentByID")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		log.Println("Student ID is required in GetStudentByID")
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Student ID is required",
		})
		return
	}

	student, err := h.Service.Student.GetByID(c, idStr)
	if err != nil {
		log.Printf("Error fetching student with ID %s: %v", idStr, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Successfully fetched student details for ID: %s", idStr)
	c.JSON(common.SUCCESS_STATUS, student)
}
