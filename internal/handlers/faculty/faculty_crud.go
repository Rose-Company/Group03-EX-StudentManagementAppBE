package faculty

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/faculty"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetList(c *gin.Context) {

	sort := c.Query("sort")

	result, err := h.Service.Faculty.GetList(c.Request.Context(), sort)
	if err != nil {
		c.JSON(common.BAD_REQUEST_STATUS, common.Response{
			Code:    common.REQUEST_FAILED,
			Message: "Failed to get faculty list",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "ok!",
		Data:    result,
	})
}

func (h *Handler) CreateAFaculty(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)

	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var faculty models.Faculty
	if err := c.ShouldBindJSON(&faculty); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	createdFaculty, err := h.Service.Faculty.CreateAFaculty(c.Request.Context(), &faculty)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty created successfully",
		Data:    createdFaculty,
	})

}

func (h *Handler) UpdateFaculty(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)

	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var faculty models.Faculty
	if err := c.ShouldBindJSON(&faculty); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	UpdateFaculty, err := h.Service.Faculty.UpdateFaculty(c.Request.Context(),id, &faculty)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty updated successfully",
		Data:    UpdateFaculty,
	})
}

func (h *Handler) DeleteFaculty(c *gin.Context) {
	ok, _ := common.ProfileFromJwt(c)

	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Faculty.DeleteFaculty(c.Request.Context(), id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code:    common.SUCCESS_STATUS,
		Message: "Faculty deleted successfully",
	})
}