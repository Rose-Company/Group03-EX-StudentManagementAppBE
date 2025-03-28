package status_transition_rule

import (
	"Group03-EX-StudentManagementAppBE/common"
	status_transition_rule_models "Group03-EX-StudentManagementAppBE/internal/models/status_transition_rule"

	"github.com/gin-gonic/gin"
)


func (h *Handler) OnOffTransition(c *gin.Context) {
	// Verify user authentication
	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	// Parse and validate update request
	var updateReq status_transition_rule_models.UpdateStatusTransitionRuleRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Call service method to update status transition rule
	result, err := h.Service.StatusTransition.OnOffTransition(c.Request.Context(), &updateReq)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	// Return successful response
	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(result))
}