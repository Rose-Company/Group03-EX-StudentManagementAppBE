// internal/handlers/auth/login.go
package auth

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/auth"

	"github.com/gin-gonic/gin"
)

// LogIn handles user authentication
func (h *Handler) LogIn(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	loginResp, err := h.Service.Auth.LoginUser(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	// Trả về ID cùng với Token
	c.JSON(common.SUCCESS_STATUS, gin.H{
		"code":  loginResp.Code,
		"id":    loginResp.ID,  
		"token": loginResp.Token,
	})
}
