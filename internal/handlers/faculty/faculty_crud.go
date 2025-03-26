package faculty

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetList(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var req models.ListFacultyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	result, err := h.Service.Faculty.GetList(c, &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(result))
}

func (h *Handler) CreateAFaculty(c *gin.Context) {

	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var faculty models.CreateFacultyRequest
	if err := c.ShouldBindJSON(&faculty); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Faculty.CreateAFaculty(c.Request.Context(), profile.Id, &faculty)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) UpdateFaculty(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var faculty models.UpdateFacultyRequest
	if err := c.ShouldBindJSON(&faculty); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	if err := h.Service.Faculty.UpdateFaculty(c.Request.Context(), profile.Id, id, &faculty); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}

func (h *Handler) DeleteFaculty(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Faculty.DeleteFaculty(c.Request.Context(), profile.Id, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
}
