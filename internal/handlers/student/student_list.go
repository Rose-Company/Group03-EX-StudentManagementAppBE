package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"

	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentList(c *gin.Context) {
	log.Println("Handling request: GetStudentList - Fetching student list")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in GetStudentList")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListStudentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding query params in GetStudentList: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
	}

	// Get student list from service
	students, err := h.Service.Student.GetStudentList(c, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, students)

}
