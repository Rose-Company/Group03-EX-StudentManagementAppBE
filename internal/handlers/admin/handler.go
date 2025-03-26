package admin

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	"Group03-EX-StudentManagementAppBE/internal/models/admin"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize = 10 << 20
)

type Handler struct {
	base.Handler
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Handler: base.NewHandler(service),
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	adminGroup := rg.Group("/v1/admins")
	adminGroup.Use(middleware.UserAuthentication)
	{
		adminGroup.POST("/imported-file", h.ImportStudentFile)
		adminGroup.GET("/imported-file/:id", h.GetImportedFile)
		adminGroup.GET("/imported-files", h.ListImportedFiles)
	}
}

func (h *Handler) ImportStudentFile(c *gin.Context) {
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		common.AbortWithError(c, err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		common.AbortWithError(c, common.ErrFileTooLarge)
		return
	}

	response, err := h.Service.Admin.ImportStudentFile(c.Request.Context(), profile.Id, header)
	if err != nil {
		common.AbortWithError(c, err)
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(response))
}

func (h *Handler) GetImportedFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		common.AbortWithError(c, common.ErrInvalidFormat)
		return
	}

	file, err := h.Service.Admin.GetImportedFile(c.Request.Context(), fileID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(file))
}

func (h *Handler) ListImportedFiles(c *gin.Context) {
	var req admin.ImportedFileListRequest

	req.Page = 1
	req.PageSize = 10

	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	files, err := h.Service.Admin.ListImportedFiles(c.Request.Context(), &req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(files))
}
