package validation

import (
	"Group03-EX-StudentManagementAppBE/common"
	validationSetting_models "Group03-EX-StudentManagementAppBE/internal/models/validation"

	"github.com/gin-gonic/gin"
)




func (h *Handler) UpdateValidationSetting(c *gin.Context) {
	
	// Verify user authentication
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	// Parse update request
	var updateReq validationSetting_models.ValidationSettingUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get validation setting ID from path parameter
	validationSettingId := c.Param("id")

	// Call service to update validation setting
	result, err := h.Service.ValidationSetting.UpdateValidationSetting(
		c.Request.Context(), 
		validationSettingId, 
		&updateReq,
	)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	// Return successful response
	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(result))
}