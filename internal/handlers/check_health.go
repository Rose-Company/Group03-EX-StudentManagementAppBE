package handlers

import (
	"ielts-web-api/common"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckStatusHealth(c *gin.Context) {
	c.JSON(common.SUCCESS_STATUS, common.ResponseSuccess(common.REQUEST_SUCCESS, "", "success"))
}
