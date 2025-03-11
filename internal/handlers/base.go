package handlers

import (
	services "ielts-web-api/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Define API route here
func (h *Handler) RegisterRoutes(c *gin.Engine) {

	health := c.Group("api/health")
	{
		health.GET("/status", h.CheckStatusHealth)
	}

}
