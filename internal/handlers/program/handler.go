package program

import (
	"Group03-EX-StudentManagementAppBE/common"

	models "Group03-EX-StudentManagementAppBE/internal/models/program"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListPrograms(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListProgramRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
	}

	programs, err := h.Service.Program.ListPrograms(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.ResponseCustom(common.REQUEST_SUCCESS, programs, "Success"))
}

func (h *Handler) CreateProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		c.JSON(http.StatusUnauthorized, common.ResponseCustom(common.REQUEST_FAILED, nil, "Unauthorized"))
		return
	}

	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		log.Printf("Error binding create program request: %v", err)
		c.JSON(http.StatusBadRequest, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	createdProgram, err := h.Service.Program.CreateProgram(c, &program)
	if err != nil {
		log.Printf("Error creating program: %v", err)
		c.JSON(http.StatusInternalServerError, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	log.Printf("Program created successfully with ID: %d", createdProgram.ID)
	c.JSON(http.StatusCreated, common.ResponseCustom(common.REQUEST_SUCCESS, createdProgram, "Program created successfully"))
}

func (h *Handler) UpdateProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		c.JSON(http.StatusUnauthorized, common.ResponseCustom(common.REQUEST_FAILED, nil, "Unauthorized"))
		return
	}

	id := c.Param("id")
	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		log.Printf("Error binding update program request: %v", err)
		c.JSON(http.StatusBadRequest, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	updatedProgram, err := h.Service.Program.UpdateProgram(c, id, &program)
	if err != nil {
		log.Printf("Error updating program: %v", err)
		c.JSON(http.StatusInternalServerError, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	log.Printf("Program updated successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.ResponseCustom(common.REQUEST_SUCCESS, updatedProgram, "Program updated successfully"))
}

func (h *Handler) DeleteProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		c.JSON(http.StatusUnauthorized, common.ResponseCustom(common.REQUEST_FAILED, nil, "Unauthorized"))
		return
	}

	id := c.Param("id")
	if err := h.Service.Program.DeleteProgram(c, id); err != nil {
		log.Printf("Error deleting program: %v", err)
		c.JSON(http.StatusInternalServerError, common.ResponseCustom(common.REQUEST_FAILED, nil, err.Error()))
		return
	}

	log.Printf("Program deleted successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.ResponseCustom(common.REQUEST_SUCCESS, nil, "Program deleted successfully"))
}
