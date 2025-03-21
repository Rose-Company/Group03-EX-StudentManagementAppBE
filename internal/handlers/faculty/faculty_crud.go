package faculty

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetList(c *gin.Context) {
	log.Println("Handling request: GetList - Fetching faculty list")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in GetList")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListFacultyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding query params in GetList: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	result, err := h.Service.Faculty.GetList(c, &req)
	if err != nil {
		log.Printf("Error fetching faculty list: %v", err)
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Failed to get faculty list",
			Data:    err.Error(),
		})
		return
	}

	log.Println("Successfully fetched faculty list")
	c.JSON(common.SUCCESS_STATUS, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "ok!",
		Data:    result,
	})
}

func (h *Handler) CreateAFaculty(c *gin.Context) {
	log.Println("Handling request: CreateAFaculty - Creating a new faculty")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in CreateAFaculty")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var faculty models.CreateFacultyRequest
	if err := c.ShouldBindJSON(&faculty); err != nil {
		log.Printf("Error binding request body in CreateAFaculty: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Faculty.CreateAFaculty(c.Request.Context(), &faculty)
	if err != nil {
		log.Printf("Error creating faculty: %v", err)
		common.AbortWithError(c, err)
		return
	}

	log.Println("Faculty created successfully")
	c.JSON(http.StatusCreated, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty created successfully",
	})
}

func (h *Handler) UpdateFaculty(c *gin.Context) {
	log.Println("Handling request: UpdateFaculty - Updating faculty")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in UpdateFaculty")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var faculty models.UpdateFacultyRequest
	if err := c.ShouldBindJSON(&faculty); err != nil {
		log.Printf("Error binding request body in UpdateFaculty: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updateFaculty, err := h.Service.Faculty.UpdateFaculty(c.Request.Context(), id, &faculty)
	if err != nil {
		log.Printf("Error updating faculty with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Faculty updated successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty updated successfully",
		Data:    updateFaculty,
	})
}

func (h *Handler) DeleteFaculty(c *gin.Context) {
	log.Println("Handling request: DeleteFaculty - Deleting faculty")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in DeleteFaculty")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Faculty.DeleteFaculty(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error deleting faculty with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Faculty deleted successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty deleted successfully",
	})
}
