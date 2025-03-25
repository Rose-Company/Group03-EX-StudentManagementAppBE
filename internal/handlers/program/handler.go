package program

import (
	"Group03-EX-StudentManagementAppBE/common"

	models "Group03-EX-StudentManagementAppBE/internal/models/program"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListPrograms(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListProgramRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
	}

	programs, err := h.Service.Program.ListPrograms(c, profile.Id, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(programs))
}

func (h *Handler) CreateProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := h.Service.Program.CreateProgram(c.Request.Context(), profile.Id, &program); err != nil {
		common.AbortWithError(c, err)
		return
	}
	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) UpdateProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		common.AbortWithError(c, err)
		return
	}
	err := h.Service.Program.UpdateProgram(c, profile.Id, id, &program)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) DeleteProgram(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok || profile == nil {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	if err := h.Service.Program.DeleteProgram(c, profile.Id, id); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}
