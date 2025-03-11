package handlers

import (
	"Group03-EX-StudentManagementAppBE/internal/services"

	"github.com/gin-gonic/gin"
)

type FeatureHandler interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

type BaseHandler struct {
	Service *services.Service
}

func NewBaseHandler(service *services.Service) BaseHandler {
	return BaseHandler{
		Service: service,
	}
}
